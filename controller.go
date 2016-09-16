package render

import (
	//"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	//"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Controller struct {
	SelfModel     mgl32.Mat4
	ParentModel   mgl32.Mat4
	Model         mgl32.Mat4
	modelUniform  int32
	parentUniform int32
	Frames
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
		ParentModel:   model,
		modelUniform:  modelUniform,
		parentUniform: parentUniform,
	}
}

func (this *Controller) SetParentModel(it mgl32.Mat4) {
	this.ParentModel = it
}

func (this *Controller) GetModel() mgl32.Mat4 {
	return this.Model
}

func (this *Controller) Render() {
	if it := this.Frames.Calc(); it != nil {
		this.Model = this.SelfModel.Mul4(*it)
	} else {
		this.Model = this.SelfModel
	}
	gl.UniformMatrix4fv(this.parentUniform, 1, false, &this.ParentModel[0])
	gl.UniformMatrix4fv(this.modelUniform, 1, false, &this.Model[0])
}
