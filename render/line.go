package render

import (
	//"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Line struct {
	Model        mgl32.Mat4
	Color        mgl32.Vec4
	modelUniform int32
	colorUniform int32
	vao          uint32
	vbo          uint32
	total        int
}

func (this *Render) CreateLine(data []float32) *Line {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STATIC_DRAW)
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	vertAttrib := uint32(gl.GetAttribLocation(this.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(this.program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	color := mgl32.Vec4{0, 0, 0, 0}
	colorUniform := gl.GetUniformLocation(this.program, gl.Str("color\x00"))
	gl.Uniform4fv(colorUniform, 1, &color[0])
	return &Line{
		vao:          vao,
		vbo:          vbo,
		total:        len(data) / 3,
		Model:        model,
		modelUniform: modelUniform,
		colorUniform: colorUniform,
	}
}

func (this *Line) Render() {
	gl.BindVertexArray(this.vao)
	gl.UniformMatrix4fv(this.modelUniform, 1, false, &this.Model[0])
	gl.Uniform4fv(this.colorUniform, 1, &this.Color[0])
	gl.DrawArrays(gl.LINES, 0, int32(this.total))
}