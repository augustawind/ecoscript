package main

import (
	"fmt"
)

func main() {
	var world = NewWorld(25, 25, 1)
	world.NewLayer("The Overworld")
	tree := mkTree()
	world.Add(tree, Vector{5, 5, 1})
	fmt.Println(world.Layer(1).Display())
}

func mkTree() *Organism {
	return NewOrganism(&Attributes{
		Name:     "Tree",
		Symbol:   '$',
		Walkable: false,
		Energy:   10,
		Size:     20,
		Mass:     20,
	}).AddClasses(
		"passive",
		"producer",
	).AddAbilities(
		&Ability{"grow", Properties{"rate": 10}},
	).Init()
}
