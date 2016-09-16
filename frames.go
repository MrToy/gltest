package render

import (
	//"fmt"
	//"github.com/go-gl/gl/v4.1-core/gl"
	//"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"sort"
)

type Frame struct {
	Time  float64
	Model mgl32.Mat4
}

type Frames struct {
	Frames   []Frame
	BaseTime float64
	//AutoRise *mgl32.Mat4
	Replay bool
}

func (this *Frames) Len() int {
	return len(this.Frames)
}
func (this *Frames) Less(i, j int) bool {
	return this.Frames[i].Time < this.Frames[j].Time
}
func (this *Frames) Swap(i, j int) {
	this.Frames[i], this.Frames[j] = this.Frames[j], this.Frames[i]
}

func (this *Frames) Add(model mgl32.Mat4, t float64) {
	frame := Frame{t, model}
	this.Frames = append(this.Frames, frame)
	sort.Sort(this)
}

func (this *Frames) Calc() *mgl32.Mat4 {
	if this.Len() < 2 {
		return nil
	}
	current := glfw.GetTime()
	elapsed := current - this.BaseTime
	i := sort.Search(this.Len(), func(i int) bool {
		return this.Frames[i].Time > elapsed
	})
	if i >= this.Len() {
		if this.Replay {
			this.BaseTime = glfw.GetTime()
		}
		return &this.Frames[i-1].Model
	}
	rate := float32(elapsed-this.Frames[i-1].Time) / float32(this.Frames[i].Time-this.Frames[i-1].Time)
	model := this.Frames[i-1].Model.Mul(1 - rate).Add(this.Frames[i].Model.Mul(rate))
	return &model
}
