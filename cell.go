package main

import (
	"math/rand"
)

type Cell struct {
	occupier *Organism
	stack    *orgStack
}

type orgStack struct {
	organisms []*Organism
	indexes   map[OrganismID]int
}

func newCell() *Cell {
	cell := new(Cell)
	cell.stack = new(orgStack)
	cell.stack.organisms = make([]*Organism, 0)
	cell.stack.indexes = make(map[OrganismID]int)
	return cell
}

func (c *Cell) Population() int {
	return len(c.stack.organisms)
}

func (c *Cell) Occupied() bool {
	return c.occupier != nil
}

func (c *Cell) Occupier() *Organism {
	return c.occupier
}

func (c *Cell) Organisms() []*Organism {
	return c.stack.organisms
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
	if !org.Walkable() {
		if c.occupier.ID() == org.ID() {
			return true
		}
	} else {
		for id := range c.stack.indexes {
			if id == org.ID() {
				return true
			}
		}
	}
	return false
}

func (c *Cell) Add(org *Organism) (exec action, ok bool) {
	if !org.Walkable() {
		if c.Occupied() {
			return
		}

		index := c.stack.indexes[org.ID()]
		exec, ok = c.setAt(index, org)
		if !ok {
			return
		}
		c.occupier = org
	}

	exec = chain(exec, func() {
		index := len(c.stack.organisms)
		c.stack.organisms = append(c.stack.organisms, org)
		c.stack.indexes[org.ID()] = index
	})
	ok = true
	return
}

//func (c *Cell) setOccupier(org *Organism) (exec action, ok bool) {
//	prevOrg := c.occupier
//	index := c.stack.indexes[prevOrg.ID()]
//	exec, ok = c.removeIndex(index)
//	if !ok {
//		return
//	}
//
//	exec = chain(exec, func() {
//		c.stack.indexes[org.ID()] = index
//		delete(c.stack.indexes, prevOrg.ID())
//		c.occupier = org
//	})
//	ok = true
//	return
//}
//
//func (c *Cell) setOccupier(org *Organism) (exec action, ok bool) {
//	prevOrg := c.occupier
//	index := c.stack.indexes[prevOrg.ID()]
//	exec, ok = c.setAt(index, org)
//	if !ok {
//		return
//	}
//
//	exec = chain(exec, func() {
//		c.stack.indexes[org.ID()] = index
//		delete(c.stack.indexes, prevOrg.ID())
//		c.occupier = org
//	})
//	ok = true
//	return
//}

func (c *Cell) Remove(org *Organism) (exec action, ok bool) {
	if org.Walkable() {
		if c.occupier.ID() == org.ID() {
			exec = func() { c.occupier = nil }
			ok = true
		}
	} else {
		for id, index := range c.stack.indexes {
			if id == org.ID() {
				exec, ok = c.removeIndex(index)
				break
			}
		}
	}
	return
}

func (c *Cell) setAt(i int, org *Organism) (exec action, ok bool) {
	exec, ok = c.removeIndex(i)
	if !ok {
		return
	}

	exec = chain(exec, func() {
		c.stack.indexes[org.ID()] = i
	})
	ok = true
	return

}

func (c *Cell) removeOrg(id OrganismID) (exec action, ok bool) {
	index, ok := c.stack.indexes[id]
	if !ok {
		return
	}
	exec, ok = c.removeIndex(index)
	if !ok {
		return
	}
	exec = chain(exec, func() {
		if id == c.occupier.ID() {
			c.occupier = nil
		}
		delete(c.stack.indexes, id)
	})
	ok = true
	return
}

func (c *Cell) removeIndex(i int) (exec action, ok bool) {
	ok = true

	if i >= len(c.stack.organisms) {
		return
	}

	exec = func() {
		copy(c.stack.organisms[i:], c.stack.organisms[i+1:])
		z := len(c.stack.organisms) - 1
		c.stack.organisms[z] = nil
		c.stack.organisms = c.stack.organisms[:z]
	}
	return
}
