package main

import (
	"fmt"
	collada "github.com/GlenKelley/go-collada"
)

func main() {
	f, _ := collada.LoadDocument("test.dae")
	fmt.Println(f.LibraryGeometries[0].Geometry[0].Mesh.Source[0].FloatArray.F())
}
