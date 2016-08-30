package render

import (
	//"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Controller struct {
	Model        mgl32.Mat4
	modelUniform int32
}

func (this *Render) CreateController() *Controller {
	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(this.program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
	return &Controller{
		Model:        model,
		modelUniform: modelUniform,
	}
}

func (this *Controller) Render() {

}
