package render

import (
	//"fmt"
	collada "github.com/GlenKelley/go-collada"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type DaeObject struct {
	Controller
	vao   uint32
	vbo   uint32
	total int
}

func (this *Render) CreateDaeObject(vert []float32) *DaeObject {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vert)*4, gl.Ptr(vert), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(this.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	//fmt.Println(vert)
	return &DaeObject{
		vao:        vao,
		vbo:        vbo,
		total:      len(vert) / 6,
		Controller: *this.CreateController(),
	}
}

func (this *DaeObject) Render() {
	gl.BindVertexArray(this.vao)
	this.Controller.Render()
	gl.DrawArrays(gl.TRIANGLES, 0, int32(this.total))
}
func (this *Render) ImportDae(filename string) (*TreeNode, error) {
	f, err := collada.LoadDocument(filename)
	if err != nil {
		return nil, err
	}
	tree := &TreeNode{}
	for _, it := range f.LibraryVisualScenes[0].VisualScene[0].Node {
		node := &TreeNode{}
		if it.HasGeometry() {
			for _, geometry := range f.LibraryGeometries[0].Geometry {
				if geometry.Name == it.InstanceGeometry[0].Name {
					var data []float32
					p := append(geometry.Mesh.Source[0].FloatArray.F32(), geometry.Mesh.Source[1].FloatArray.F32()...)
					m := geometry.Mesh.Polylist[0].P.I()
					for _, i := range m {
						data = append(data, p[i*3])
						data = append(data, p[i*3+1])
						data = append(data, p[i*3+2])
					}
					obj := this.CreateDaeObject(data)
					model := mgl32.Mat4{}
					for i, v := range it.Matrix[0].F32() {
						model[i] = v
					}
					obj.SelfModel = model.Transpose()
					node.Data = obj
				}
			}
		}
		tree.Add(node)
	}
	// for _, it := range f.LibraryGeometries[0].Geometry {
	// 	var data []float32
	// 	p := append(it.Mesh.Source[0].FloatArray.F32(), it.Mesh.Source[1].FloatArray.F32()...)
	// 	m := it.Mesh.Polylist[0].P.I()
	// 	for _, i := range m {
	// 		data = append(data, p[i*3])
	// 		data = append(data, p[i*3+1])
	// 		data = append(data, p[i*3+2])
	// 	}
	// 	obj := this.CreateDaeObject(data)
	// 	tree.Add(&TreeNode{Data: obj})
	// }
	//fmt.Println(f.LibraryGeometries[0].Geometry[0].Mesh.Source[0].FloatArray.F())
	return tree, nil
}
