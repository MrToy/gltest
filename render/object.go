package render

import (
	"github.com/go-gl/mathgl/mgl32"
)

const (
	XYZUV = iota
	LINE
)

type Object struct {
	Data     *[]float32
	Type     int
	Model    mgl32.Mat4
	Image    string
	Color    *[3]float32
	Children []*Object
	Parent   *Object
	Frames   Frames
	OnRender func()
}

func NewObject() *Object {
	return &Object{
		Model: mgl32.Ident4(),
	}
}

func recur(objs []*Object, items []*Object) []*Object {
	for _, item := range items {
		objs = append(objs, item)
		objs = recur(objs, item.Children)
	}
	return objs
}

func (this *Object) GetAll() []*Object {
	objs := []*Object{}
	objs = append(objs, this)
	objs = recur(objs, this.Children)
	return objs
}

func (this *Object) AddChild(child *Object) {
	this.Children = append(this.Children, child)
	child.Parent = this
}

//var previousTime = time.Now()

// func (this *Object) Render() {
// 	angle := 0.0
// 	elapsed := time.Since(previousTime).Seconds()
// 	angle += elapsed
// 	//fmt.Println(elapsed)
// 	this.Model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
// 	// gl.UniformMatrix4fv(this.modelUniform, 1, false, &this.Model[0])
// 	// gl.VertexAttribPointer(this.vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
// 	// gl.BindBuffer(gl.ARRAY_BUFFER, this.Vbo)
// 	// gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
// }
