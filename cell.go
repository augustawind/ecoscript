package main

import (
	"math/rand"
)

type Cell struct {
	occupier  *Organism
	cover     *Organism
	walkables []*Organism
	indexes   map[OrganismID]int
}

func newCell() *Cell {
	cell := new(Cell)
	cell.indexes = make(map[OrganismID]int)
	return cell
}

func (c *Cell) Occupied() bool {
	return c.occupier != nil
}

func (c *Cell) Occupier() *Organism {
	return c.occupier
}

func (c *Cell) Organisms() []*Organism {
	orgs := make([]*Organism, len(c.walkables)+1)
	copy(orgs, c.walkables)
	orgs = append(orgs, c.occupier)
	return orgs
}

func (c *Cell) Shuffled() []*Organism {
	orgs := c.Organisms()
	shuffled := make([]*Organism, len(orgs))

	for i, j := range rand.Perm(len(orgs)) {
		shuffled[i] = orgs[j]
	}
	return shuffled
}

func (c *Cell) Exists(org *Organism) bool {
	if org.Walkable() {
		if c.occupier.id == org.id {
			return true
		}
	} else {
		for id := range c.indexes {
			if id == org.id {
				return true
			}
		}
	}
	return false
}

func (c *Cell) Add(org *Organism) (exec func(), ok bool) {
	if org.Walkable() {
		if !c.Occupied() {
			exec = func() { c.occupier = org }
			ok = true
		}
	} else {
		exec = func() {
			c.indexes[org.id] = len(c.walkables)
			c.walkables = append(c.walkables, org)
		}
		ok = true
	}
	return
}

func (c *Cell) Remove(org *Organism) (exec func(), ok bool) {
	if org.Walkable() {
		if c.occupier.id == org.id {
			exec = func() { c.occupier = nil }
			ok = true
		}
	} else {
		for id, index := range c.indexes {
			if id == org.id {
				exec = func() { c.delWalkable(index) }
				ok = true
				break
			}
		}
	}
	return
}

func (c *Cell) delWalkable(i int) {
	copy(c.walkables[i:], c.walkables[i+1:])
	z := len(c.walkables) - 1
	c.walkables[z] = nil
	c.walkables = c.walkables[:z]
}
