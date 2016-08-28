package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Object struct {
	vao   uint32
	vbo   uint32
	total int
}

var (
	vao uint32
)

func (this *Render) CreateObject(data []float32) *Object {
	if vao == 0 {
		gl.GenVertexArrays(1, &vao)
		gl.BindVertexArray(vao)
	}
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

	return &Object{
		vao:   vao,
		vbo:   vbo,
		total: len(data) / 5,
	}
}

func (this *Object) Render() {
	gl.BindVertexArray(this.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(this.total))
}
