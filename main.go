package main

import "fmt"

func main() {
	var world = NewWorld(25, 25, 1)
	world.NewLayer("The Overworld")
	tree := mkTree()
	world.Add(tree, Vector{5, 5, 1})
	fmt.Println(world.Layer(1).Display())
}

func mkTree() *Organism {
	o := NewOrganism(Attributes{
		name:     "Tree",
		display:  '$',
		walkable: false,
		energy:   10,
		size:     20,
		mass:     20,
		classes: []string{
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
