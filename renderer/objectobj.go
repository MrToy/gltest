package renderer

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

var (
	vao uint32
)

type ObjectObj struct {
	vbo    uint32
	length int32
}

func setProgram() {
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
}

func NewObjectObj(data []float32) ObjectObj {
	if vao == 0 {
		setVao()
	}
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STATIC_DRAW)
	return ObjectObj{
		vbo:    vbo,
		length: int32(len(data) / 5),
	}
}

func (this *ObjectObj) Render() {
	gl.VertexAttribPointer(this.vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.VertexAttribPointer(this.texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.DrawArrays(gl.TRIANGLES, 0, it.Len)
}
