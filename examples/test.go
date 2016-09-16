package main

import (
	"fmt"
	collada "github.com/GlenKelley/go-collada"
)

func main() {
	f, _ := collada.LoadDocument("xx.dae")
	for _, it := range f.LibraryGeometries[0].Geometry {
		fmt.Println(it.Mesh.Source[1].FloatArray.F())
	}
	fmt.Println(len(f.LibraryGeometries[0].Geometry[0].Mesh.Polylist[0].P.I()))
}
