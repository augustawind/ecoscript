package main

import (
	"math/rand"
)

type Cell struct {
	occupier *Entity
	stack    *entStack
}

type entStack struct {
	entities []*Entity
	indexes  map[EntityID]int
}

func newCell() *Cell {
	cell := new(Cell)
	cell.stack = new(entStack)
	cell.stack.entities = make([]*Entity, 0)
	cell.stack.indexes = make(map[EntityID]int)
	return cell
}

func (c *Cell) Population() int {
	return len(c.stack.entities)
}

func (c *Cell) Occupied() bool {
	return c.occupier != nil
}

func (c *Cell) Occupier() *Entity {
	return c.occupier
}

func (c *Cell) Entities() []*Entity {
	return c.stack.entities
}

func (c *Cell) Shuffled() []*Entity {
	ents := c.Entities()
	shuffled := make([]*Entity, len(ents))
	for i, j := range rand.Perm(len(ents)) {
		shuffled[i] = ents[j]
	}
	return shuffled
}

func (c *Cell) Exists(ent *Entity) bool {
	if !ent.Walkable() {
		if c.occupier.ID() == ent.ID() {
			return true
		}
	} else {
		for id := range c.stack.indexes {
			if id == ent.ID() {
				return true
			}
		}
	}
	return false
}

func (c *Cell) Add(ent *Entity) (exec action, ok bool) {
	if !ent.Walkable() {
		if c.Occupied() {
			return
		}

		index := c.stack.indexes[ent.ID()]
		exec, ok = c.setAt(index, ent)
		if !ok {
			return
		}
		c.occupier = ent
	}

	exec = chain(exec, func() {
		index := len(c.stack.entities)
		c.stack.entities = append(c.stack.entities, ent)
		c.stack.indexes[ent.ID()] = index
	})
	ok = true
	return
}

//func (c *Cell) setOccupier(ent *Entity) (exec action, ok bool) {
//	prevOrg := c.occupier
//	index := c.stack.indexes[prevOrg.ID()]
//	exec, ok = c.removeIndex(index)
//	if !ok {
//		return
//	}
//
//	exec = chain(exec, func() {
//		c.stack.indexes[ent.ID()] = index
//		delete(c.stack.indexes, prevOrg.ID())
//		c.occupier = ent
//	})
//	ok = true
//	return
//}
//
//func (c *Cell) setOccupier(ent *Entity) (exec action, ok bool) {
//	prevOrg := c.occupier
//	index := c.stack.indexes[prevOrg.ID()]
//	exec, ok = c.setAt(index, ent)
//	if !ok {
//		return
//	}
//
//	exec = chain(exec, func() {
//		c.stack.indexes[ent.ID()] = index
//		delete(c.stack.indexes, prevOrg.ID())
//		c.occupier = ent
//	})
//	ok = true
//	return
//}

func (c *Cell) Remove(ent *Entity) (exec action, ok bool) {
	index := c.stack.indexes[ent.ID()]
	exec, ok = c.removeIndex(index)
	if !ok {
		return
	}

	if !ent.Walkable() && c.Occupied() && c.occupier.ID() == ent.ID() {
		exec = chain(exec, func() {
			c.occupier = nil
		})
		ok = true
	}
	return
}

func (c *Cell) setAt(i int, ent *Entity) (exec action, ok bool) {
	exec, ok = c.removeIndex(i)
	if !ok {
		return
	}

	exec = chain(exec, func() {
		c.stack.indexes[ent.ID()] = i
	})
	ok = true
	return

}

func (c *Cell) removeEnt(id EntityID) (exec action, ok bool) {
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

	if i >= len(c.stack.entities) {
		return
	}

	exec = func() {
		copy(c.stack.entities[i:], c.stack.entities[i+1:])
		z := len(c.stack.entities) - 1
		c.stack.entities[z] = nil
		c.stack.entities = c.stack.entities[:z]
	}
	return
}
