package renderer

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"runtime"
)

const (
	windowWidth  = 800
	windowHeight = 800
)

var vertexShader = `
#version 330

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vert;
in vec2 vertTexCoord;

out vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * camera * model * vec4(vert, 1);
}
` + "\x00"

var fragmentShader = `
#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;

out vec4 outputColor;

void main() {
    outputColor = texture(tex, fragTexCoord);
}
` + "\x00"

type RenderList struct {
	Item    *Object2
	Texture *uint32
	Vbo     uint32
	Len     int32
}

type Render struct {
	window         *glfw.Window
	program        uint32
	vao            uint32
	items          []*Object2
	lists          []RenderList
	colorUniform   int32
	modelUniform   int32
	vertAttrib     uint32
	texCoordAttrib uint32
	scene          *Object2
	tree           *TreeNode
}

func (this *Render) Close() {
	glfw.Terminate()
}

func (this *Render) SetScene(scene *Object2) {
	items := scene.GetAll()
	//gl.BufferData(gl.ARRAY_BUFFER, len(Scene.Data)*4, gl.Ptr(Scene.Data), gl.STATIC_DRAW)
	var lists []RenderList
	for _, it := range items {
		if it.Data == nil {
			continue
		}
		var list RenderList
		if len(it.Image) != 0 {
			texture, err := NewTexture(it.Image)
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
	this.lists = lists
	this.items = items
}

func New() *Render {
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
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
	program, err := NewProgram(vertexShader, fragmentShader)
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
		var pitch float64 = 20
		var yaw float64 = 10
		var length float64 = 20
		camera := mgl32.LookAtV(GetRotateVec3(pitch, yaw, length), mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
		gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
		window.SetScrollCallback(func(w *glfw.Window, xoff float64, yoff float64) {
			length -= yoff
			if length < 2 {
				length = 2
			}
			if length > 50 {
				length = 50
			}
			camera := mgl32.LookAtV(GetRotateVec3(pitch, yaw, length), mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
			gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
		})
		var isRotateing bool
		var prePosX, prePosY float64
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
			prePosX, prePosY = w.GetCursorPos()
		})
		window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
			if !isRotateing {
				return
			}
			x, y := xpos-prePosX, ypos-prePosY
			prePosX, prePosY = xpos, ypos
			pitch += y * 0.4
			yaw += x * 0.4
			if pitch > 90 {
				pitch = 90
			}
			if pitch < -45 {
				pitch = -45
			}
			camera := mgl32.LookAtV(GetRotateVec3(pitch, yaw, length), mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
			gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
		})
	}()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)
	colorUniform := gl.GetUniformLocation(program, gl.Str("color\x00"))
	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	//var currentVbo uint32
	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	return &Render{
		window:         window,
		program:        program,
		vao:            vao,
		colorUniform:   colorUniform,
		modelUniform:   modelUniform,
		vertAttrib:     vertAttrib,
		texCoordAttrib: texCoordAttrib,
	}
}

func (this *Render) Run() {
	for !this.window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Render
		gl.UseProgram(this.program)
		gl.BindVertexArray(this.vao)

		for _, it := range this.items {
			if it.OnRender != nil {
				it.OnRender()
			}
		}
		for _, it := range this.lists {
			model := it.Item.Model
			// if it.Item.Parent != nil {
			// 	model = it.Item.Model.Mul4(it.Item.Parent.Model)
			// }
			if it.Texture != nil {
				gl.ActiveTexture(gl.TEXTURE0)
				gl.BindTexture(gl.TEXTURE_2D, *it.Texture)
			} else {
				gl.Disable(gl.TEXTURE0)
			}
			if it.Item.Color != nil {
				gl.Uniform3fv(this.colorUniform, 1, &it.Item.Color[0])
			}

			gl.UniformMatrix4fv(this.modelUniform, 1, false, &model[0])
			gl.BindBuffer(gl.ARRAY_BUFFER, it.Vbo)
			if it.Item.Type == LINE {
				gl.VertexAttribPointer(this.vertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
				gl.DrawArrays(gl.LINES, 0, it.Len)
			} else {
				gl.VertexAttribPointer(this.vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
				gl.VertexAttribPointer(this.texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
				gl.DrawArrays(gl.TRIANGLES, 0, it.Len)
			}
		}
		// Maintenance
		this.window.SwapBuffers()
		glfw.PollEvents()
	}
}
