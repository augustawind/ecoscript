package main

import (
	"fmt"
	"log"
)

func main() {
	mapfile, err := ParseMapfile("examples/Mapfile")
	if err != nil {
		log.Fatal(err)
	}

	world := mapfile.ToWorld()
	fmt.Println(world.Layer(0).Display())
}
