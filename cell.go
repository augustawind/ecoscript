package main

import (
	"math/rand"
)

type Cell struct {
	// main
	occupier *Organism
	// over
	cover     *Organism
	walkables []*Organism
	depths    map[OrganismID]int
}

func (c *Cell) Occupied() bool {
	return c.occupier != nil
}

func (c *Cell) Occupier() *Organism {
	return c.occupier
}

func (c *Cell) Organisms() []*Organism {
	organisms := make([]*Organism, len(c.walkables)+1)
	copy(organisms, c.walkables)
	organisms = append(organisms, c.occupier)
	return organisms
}

func (c *Cell) Shuffled() []*Organism {
	organisms := c.Organisms()
	n := len(organisms)
	shuffled := make([]*Organism, n)

	for i, j := range rand.Perm(n) {
		shuffled[i] = organisms[j]
	}
	return shuffled
}

func (c *Cell) Exists(organism *Organism) bool {
	if organism.Walkable {
		if c.occupier.id == organism.id {
			return true
		}
	} else {
		for id := range c.depths {
			if id == organism.id {
				return true
			}
		}
	}
	return false
}

func (c *Cell) Add(organism *Organism) (exec func(), ok bool) {
	if organism.Walkable {
		if !c.Occupied() {
			exec = func() { c.occupier = organism }
			ok = true
		}
	} else {
		exec = func() {
			c.depths[organism.id] = len(c.walkables)
			c.walkables = append(c.walkables, organism)
		}
		ok = true
	}
	return
}

func (c *Cell) Remove(organism *Organism) (exec func(), ok bool) {
	if organism.Walkable {
		if c.occupier.id == organism.id {
			exec = func() { c.occupier = nil }
			ok = true
		}
	} else {
		for id, depth := range c.depths {
			if id == organism.id {
				exec = func() { c.delWalkable(depth) }
				ok = true
				break
			}
		}
	}
	return
}

func (c *Cell) delWalkable(depth int) {
	copy(c.walkables[depth:], c.walkables[depth+1:])
	z := len(c.walkables) - 1
	c.walkables[z] = nil
	c.walkables = c.walkables[:z]
}
