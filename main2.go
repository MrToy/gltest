package main

import (
	"github.com/MrToy/gltest/render"
	"github.com/go-gl/mathgl/mgl32"
	//"math/rand"
	//"time"
)

func main() {
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
	r, err := render.New()
	if err != nil {
		panic(err)
	}
	o1 := r.CreateObject(data)
	o1.Color = mgl32.Vec4{1, 0, 0, 1}
	r.Scene.Add(o1)
	o2 := r.CreateObject(data)
	o2.Model = mgl32.Translate3D(2, 0, 0)
	o2.SetTexture("square.png")
	r.Scene.Add(o2)

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
	r.Scene.Add(o3)
	r.Run()
}
