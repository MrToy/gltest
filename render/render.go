package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	//"log"
	"runtime"
)

const (
	vertexShader = `
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

	fragmentShader = `
#version 330

uniform sampler2D tex;
uniform vec4 color;
in vec2 fragTexCoord;

out vec4 outputColor;

void main() {
    outputColor = vec4(vec3(color)*color[3]+vec3(texture(tex, fragTexCoord)*(1-color[3])),1);
}
` + "\x00"
)

type Render struct {
	Scene   *TreeNode
	window  *glfw.Window
	program uint32
}

func (this *Render) Close() {
	glfw.Terminate()
}

func New() (*Render, error) {
	windowWidth, windowHeight := 800, 800
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		return nil, err
	}
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Test", nil, nil)
	if err != nil {
		return nil, err
	}
	window.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		return nil, err
	}
	program, err := NewProgram(vertexShader, fragmentShader)
	if err != nil {
		return nil, err
	}

	gl.UseProgram(program)

	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), 0.1, 500.0)
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
	window.SetSizeCallback(func(w *glfw.Window, width int, height int) {
		projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/float32(height), 0.1, 500.0)
		gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
		gl.Viewport(0, 0, int32(width), int32(height))
	})
	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

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
	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	render := &Render{
		window:  window,
		program: program,
		Scene:   &TreeNode{},
	}
	return render, nil
}

func (this *Render) Run() {
	for !this.window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(this.program)
		for _, it := range this.Scene.GetAll() {
			if it.Object != nil {
				it.Object.Render()
			}
		}
		this.window.SwapBuffers()
		glfw.PollEvents()
	}
}
