package render

import (
	//"fmt"
	//"github.com/go-gl/gl/v4.1-core/gl"
	//"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"sort"
)

type Frame struct {
	Time  float64
	Model mgl32.Mat4
}

type Frames struct {
	Frames   []Frame
	Current  float64
	AutoRise *mgl32.Mat4
	Replay   bool
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

func (this *Frames) Calc(t float64) *mgl32.Mat4 {
	i := sort.Search(this.Len(), func(i int) bool {
		return this.Frames[i].Time > t
	})
	if i <= 0 || i >= this.Len() {
		return nil
	}
	rate := float32(this.Frames[i].Time-t) / float32(this.Frames[i].Time-this.Frames[i-1].Time)
	m1 := this.Frames[i-1].Model.Mul(1 - rate)
	m2 := this.Frames[i].Model.Mul(rate)
	model := m1.Add(m2)
	return &model
}
