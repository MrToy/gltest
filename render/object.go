package render

import (
	//"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Object struct {
	Controller
	Color        mgl32.Vec4
	colorUniform int32
	vao          uint32
	vbo          uint32
	total        int
	texture      uint32
}

func (this *Render) CreateObject(data []float32) *Object {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(this.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	texCoordAttrib := uint32(gl.GetAttribLocation(this.program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	color := mgl32.Vec4{0, 0, 0, 0}
	colorUniform := gl.GetUniformLocation(this.program, gl.Str("color\x00"))
	gl.Uniform4fv(colorUniform, 1, &color[0])

	return &Object{
		vao:          vao,
		vbo:          vbo,
		total:        len(data) / 5,
		colorUniform: colorUniform,
		Controller:   *this.CreateController(),
	}
}

func (this *Object) Render() {
	this.Controller.Render()
	gl.BindVertexArray(this.vao)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, this.texture)
	gl.UniformMatrix4fv(this.modelUniform, 1, false, &this.Model[0])
	gl.Uniform4fv(this.colorUniform, 1, &this.Color[0])
	gl.DrawArrays(gl.TRIANGLES, 0, int32(this.total))
}
func (this *Object) SetTexture(file string) {
	this.texture, _ = NewTexture(file)
}
