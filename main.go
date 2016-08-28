package main

import (
	"github.com/MrToy/gltest/renderer"
	"github.com/go-gl/mathgl/mgl32"
	//"math/rand"
	//"time"
)

func main() {
	render := renderer.New()
	defer render.Close()
	world := renderer.NewObject()
	tree := renderer.NewObject()
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
	tree.Type = renderer.LINE
	tree.Model = mgl32.Translate3D(0.5, 0, 0.5)
	tree.Color = &[3]float32{0.8, 0.8, 0.8}

	axisX := renderer.NewObject()
	axisX.Data = &[]float32{
		-2.0, 0.0, 0.0,
		2.0, 0.0, 0.0,
	}
	axisX.Type = renderer.LINE
	axisX.Color = &[3]float32{0.8, 0, 0}

	axisZ := renderer.NewObject()
	axisZ.Data = &[]float32{
		0.0, -2.0, 0.0,
		0.0, 2.0, 0.0,
	}
	axisZ.Type = renderer.LINE
	axisZ.Color = &[3]float32{0, 0.8, 0}

	axisY := renderer.NewObject()
	axisY.Data = &[]float32{
		0.0, 0.0, -2.0,
		0.0, 0.0, 2.0,
	}
	axisY.Type = renderer.LINE
	axisY.Color = &[3]float32{0, 0, 0.8}
	// world.AddChild(axisY)

	person := renderer.NewObject()
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
	//person.Image = "square.png"

	world.AddChild(person)
	world.AddChild(tree)
	world.AddChild(axisX)
	world.AddChild(axisY)
	world.AddChild(axisZ)

	render.SetScene(world)
	render.Run()
}
