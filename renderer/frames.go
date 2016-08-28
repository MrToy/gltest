package renderer

import (
	"github.com/go-gl/mathgl/mgl32"
	"time"
)

type Frame struct {
	Model mgl32.Mat4
	Time  time.Time
}
type Frames struct {
	All     []Frame
	Current time.Time
	Loop    bool
}

func (this *Frames) Add(frame Frame) {
	this.All = append(this.All, frame)
}
