package main

import (
	"github.com/MrToy/gltest/render"
	//"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	//"math/rand"
	//"time"
	"fmt"
)

func main() {
	r, err := render.New()
	if err != nil {
		panic(err)
	}

	data := []float32{
		-1.0, 0.0, -1.0, 0.0, 0.0,
		1.0, 0.0, 0.0, 1.0, 0.0,
		-1.0, 0.0, 1.0, 0.0, 1.0,

		0.0, 1.0, 0.0, 1.0, 0.0,
		1.0, 0.0, 0.0, 0.0, 0.0,
		-1.0, 0.0, 1.0, 0.0, 1.0,

		-1.0, 0.0, -1.0, 0.0, 1.0,
		0.0, 1.0, 0.0, 1.0, 0.0,
		-1.0, 0.0, 1.0, 0.0, 0.0,

		-1.0, 0.0, -1.0, 1.0, 0.0,
		1.0, 0.0, 0.0, 0.0, 1.0,
		0.0, 1.0, 0.0, 0.0, 0.0,
	}

	o1 := r.CreateObject(data)
	o1.SetTexture("square.png")
	o1.SelfModel = mgl32.Scale3D(0.3, 0.3, 0.3)
	frames := render.Frames{Replay: true}
	frames.Add(mgl32.Translate3D(10, 0, 0), 0)
	frames.Add(mgl32.Translate3D(1, 0, 0), 2)
	frames.Add(mgl32.Translate3D(1, 0, 5), 5)
	frames.Add(mgl32.Translate3D(-5, 0, 5), 9)

	frames2 := render.Frames{Replay: true}
	frames2.Add(mgl32.HomogRotate3DY(0), 0)
	frames2.Add(mgl32.HomogRotate3DY(360), 2)

	x1 := r.CreateController()
	x1.Frames = frames2
	//o1.Frames = frames

	render.AddToTree(x1, o1)
	render.AddToTree(r.Scene, x1)

	var lineData []float32
	for z := -15; z < 16; z++ {
		it := []float32{-15, 0.0, float32(z), 15, 0.0, float32(z)}
		lineData = append(lineData, it...)
	}
	for x := -15; x < 16; x++ {
		it := []float32{float32(x), 0.0, -15.0, float32(x), 0.0, 15.0}
		lineData = append(lineData, it...)
	}
	o3 := r.CreateLine(lineData)
	render.AddToTree(r.Scene, o3)
	for it := range r.Scene.WalkTree() {
		if item, ok := it.(render.Object); ok {
			fmt.Println(item)
		}
	}
	r.Run()
}
