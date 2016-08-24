package main

import (
	"github.com/MrToy/gltest/render"
	"github.com/go-gl/mathgl/mgl32"
	//"math/rand"
	//"time"
)

func main() {
	world := render.NewObject()
	tree := render.NewObject()
	var lineData []float32
	for z := -15; z < 16; z++ {
		it := []float32{-15, 0.0, float32(z), 15, 0.0, float32(z)}
		lineData = append(lineData, it...)
	}
	for x := -15; x < 16; x++ {
		it := []float32{float32(x), 0.0, -15.0, float32(x), 0.0, 15.0}
		lineData = append(lineData, it...)
	}
	tree.Data = &lineData
	tree.Type = render.LINE
	tree.Model = mgl32.Translate3D(0.5, 0, 0.5)
	tree.Color = &[3]float32{0.8, 0.8, 0.8}

	axis := render.NewObject()
	axis.Data = &[]float32{
		-10.0, 0.0, 0.0,
		15.0, 0.0, 0.0,
		0.0, 0.0, -10.0,
		0.0, 0.0, 15.0,
	}
	axis.Type = render.LINE
	axis.Color = &[3]float32{0.8, 0, 0}

	person := render.NewObject()
	person.Data = &[]float32{
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
		0.0, 1.0, 0.0, 0.0, 1.0,
		-1.0, 0.0, 1.0, 0.0, 0.0,
	}
	person.Model = mgl32.Scale3D(0.4, 0.4, 0.4)
	person.Image = "square.png"

	world.AddChild(person)
	world.AddChild(tree)
	world.AddChild(axis)
	render.SetScene(world)
	render.Run()
}
