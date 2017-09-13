package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	var world = NewWorld(25, 25, 1)
	world.NewLayer("The Overworld")
	tree := mkTree()
	world.Add(tree, Vector{5, 5, 1})
	fmt.Println(world.Layer(1).Display())
	spew.Dump(BehaviorIndex)
}

func mkTree() *Organism {
	o := NewOrganism(Attributes{
		Name:     "Tree",
		Symbol:   '$',
		Walkable: false,
		Energy:   10,
		Size:     20,
		Mass:     20,
		Classes: []string{
			"passive",
			"producer",
		},
	})
	o.AddBehaviors(
		&Grow{Rate: 10},
	)
	o.Init()
	return o
}
