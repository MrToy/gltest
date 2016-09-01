package render

import (
	//"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	//"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Controller struct {
	TreeNode
	SelfModel     mgl32.Mat4
	Model         mgl32.Mat4
	modelUniform  int32
	parentUniform int32
	Frames
}

type ModelGetter interface {
	GetModel() mgl32.Mat4
}

func (this *Render) CreateController() *Controller {
	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(this.program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
	parentUniform := gl.GetUniformLocation(this.program, gl.Str("parent\x00"))
	gl.UniformMatrix4fv(parentUniform, 1, false, &model[0])
	return &Controller{
		SelfModel:     model,
		Model:         model,
		modelUniform:  modelUniform,
		parentUniform: parentUniform,
	}
}

func (this *Controller) Render() {
	if it := this.Frames.Calc(); it != nil {
		this.Model = this.SelfModel.Mul4(*it)
	} else {
		this.Model = this.SelfModel
	}
	model := mgl32.Ident4()
	if it, ok := this.Parent.(ModelGetter); ok {
		model = it.GetModel()
	}
	gl.UniformMatrix4fv(this.parentUniform, 1, false, &model[0])
	gl.UniformMatrix4fv(this.modelUniform, 1, false, &this.Model[0])
}

func (this *Controller) GetModel() mgl32.Mat4 {
	return this.Model
}
