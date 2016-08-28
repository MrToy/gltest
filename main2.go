package main

import (
	"github.com/MrToy/gltest/render"
	//"github.com/go-gl/mathgl/mgl32"
	//"math/rand"
	//"time"
)

func main() {
	r, err := render.New()
	if err != nil {
		panic(err)
	}
	r.Run()
}
