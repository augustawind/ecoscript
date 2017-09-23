package main

import (
	"fmt"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	mapfile, err := ParseMapfile("examples/Mapfile")
	if err != nil {
		log.Fatal(err)
	}

	world := mapfile.ToWorld()

	for {
		fmt.Println(world.Layer(0).Display())
		world.Tick()
		time.Sleep(500 * time.Millisecond)
	}
}

func show(s ...interface{}) {
	spew.Dump(fmt.Sprint(s...))
}

func showf(format string, a ...interface{}) {
	spew.Dump(fmt.Sprintf(format, a...))
}
