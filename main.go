package main

import (
	"log"
)

func main() {
	mapfile, err := ParseMapfile("examples/Mapfile")
	if err != nil {
		log.Fatal(err)
	}

	world := mapfile.ToWorld()
	log.Println(world.Layer(0).Display())
}

