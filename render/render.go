package render

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"math"
	"os"
	"runtime"
	"strings"
)

const (
	windowWidth  = 800
	windowHeight = 800

	vertexShader = `
#version 330

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vert;
in vec2 vertTexCoord;
out vec2 fragTex;

void main() {
    fragTex = vertTexCoord;
    gl_Position = projection * camera * model * vec4(vert, 1);
}
` + "\x00"

	fragmentShader = `
#version 330

uniform sampler2D tex;
uniform vec3 color;
uniform bool useColor;

in vec2 fragTex;

out vec4 outputColor;

void main() {
	if(useColor){
		outputColor = vec4(color,0);
	}else{
		outputColor = texture(tex, fragTex);
	}
}
` + "\x00"
)

var (
	Scene *Object
)

type RenderList struct {
	Item    *Object
	Texture *uint32
	Vbo     uint32
	Len     int32
}

func init() {
	runtime.LockOSThread()
}

func SetScene(scene *Object) {
	Scene = scene
}

func Run() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Test", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	// Configure the vertex and fragment shaders
	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 500.0)
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
	window.SetSizeCallback(func(w *glfw.Window, width int, height int) {
		projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/float32(height), 0.1, 500.0)
		gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	func() {
		from := mgl32.Vec3{15, 15, 15}
		target := mgl32.Vec3{0, 0, 0}
		camera := mgl32.LookAtV(from, target, mgl32.Vec3{0, 1, 0})
		gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
		window.SetScrollCallback(func(w *glfw.Window, xoff float64, yoff float64) {
			if (from.Sub(target).Len() < 3 && yoff < 0) || (from.Sub(target).Len() > 100 && yoff > 0) {
				return
			}
			from = from.Add(from.Sub(target).Normalize().Mul(float32(yoff)))
			camera := mgl32.LookAtV(from, target, mgl32.Vec3{0, 1, 0})
			gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
		})
		var isRotateing bool
		window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
			if button != glfw.MouseButtonRight {
				return
			}
			if action == glfw.Press {
				isRotateing = true
			}
			if action == glfw.Release {
				isRotateing = false
			}
		})
		var prePosX, prePosY float64
		window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
			if !isRotateing {
				return
			}
			x, y := xpos-prePosX, ypos-prePosY
			prePosX, prePosY = xpos, ypos
			if math.Abs(x) > 15 || math.Abs(y) > 15 {
				return
			}
			from = mgl32.Vec3{
				from[0]*float32(math.Cos(x/360*math.Pi)) - from[2]*float32(math.Sin(x/360*math.Pi)),
				from[1] + float32(y),
				from[0]*float32(math.Sin(x/360*math.Pi)) + from[2]*float32(math.Cos(x/360*math.Pi)),
			}
			camera := mgl32.LookAtV(from, target, mgl32.Vec3{0, 1, 0})
			gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
		})
	}()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)
	colorUniform := gl.GetUniformLocation(program, gl.Str("color\x00"))
	useColorUniform := gl.GetUniformLocation(program, gl.Str("useColor\x00"))
	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))
	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)

	//var currentVbo uint32
	items := Scene.GetAll()
	//gl.BufferData(gl.ARRAY_BUFFER, len(Scene.Data)*4, gl.Ptr(Scene.Data), gl.STATIC_DRAW)
	var lists []RenderList
	for _, it := range items {
		if it.Data == nil {
			continue
		}
		var list RenderList
		if len(it.Image) != 0 {
			texture, err := newTexture(it.Image)
			if err != nil {
				log.Fatalln(err)
			}
			list.Texture = &texture
		}
		var vbo uint32
		gl.GenBuffers(1, &vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(*it.Data)*4, gl.Ptr(*it.Data), gl.STATIC_DRAW)
		list.Item = it
		list.Vbo = vbo
		if it.Type == LINE {
			list.Len = int32(len(*it.Data) / 3)
		} else {
			list.Len = int32(len(*it.Data) / 5)
		}
		lists = append(lists, list)
	}
	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Render
		gl.UseProgram(program)
		gl.BindVertexArray(vao)

		for _, it := range items {
			if it.OnRender != nil {
				it.OnRender()
			}
		}
		for _, it := range lists {
			model := it.Item.Model
			if it.Item.Parent != nil {
				model = it.Item.Model.Mul4(it.Item.Parent.Model)
			}
			if it.Texture != nil {
				gl.ActiveTexture(gl.TEXTURE0)
				gl.BindTexture(gl.TEXTURE_2D, *it.Texture)
			}
			if it.Item.Color != nil {
				gl.Uniform1i(useColorUniform, 1)
				gl.Uniform3fv(colorUniform, 1, &it.Item.Color[0])
			} else {
				gl.Uniform1i(useColorUniform, 0)
			}

			gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
			gl.BindBuffer(gl.ARRAY_BUFFER, it.Vbo)
			if it.Item.Type == LINE {
				gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
				gl.DrawArrays(gl.LINES, 0, it.Len)
			} else {
				gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
				gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
				gl.DrawArrays(gl.TRIANGLES, 0, it.Len)
			}
		}
		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func newTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, nil
}
