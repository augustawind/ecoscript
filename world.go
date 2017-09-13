package main

import "math/rand"

// spaces have width and height and contain Cells.
type space interface {
	// Width returns the length of the x-axis.
	Width() int

	// Height returns the length of the y-axis.
	Height() int

	// InBounds returns true if the given Vector is in bounds.
	InBounds(vec Vector) bool

	// Cell returns the Cell at the given vector.
	Cell(vec Vector) *Cell

	// View returns all Vectors that are in-bounds within a radius.
	View(origin Vector, radius int) []Vector

	// ViewR is like View but randomizes the returned Vectors.
	ViewR(origin Vector, radius int) []Vector

	// ViewWalkable returns all walkable Vectors within a radius.
	ViewWalkable(origin Vector, radius int) []Vector

	// ViewWalkableR is like ViewWalkable but randomizes the returned Vectors.
	ViewWalkableR(origin Vector, radius int) []Vector

	// RandWalkable finds a random walkable Vector within a radius.
	RandWalkable(origin Vector, radius int) Vector

	// Remove attempts to remove an Organism at the given Vector.
	// It returns true if it succeeded or false if it wasn't found.
	Remove(org *Organism, src Vector) (ok bool)

	// Move attempts to move an Organism from one Vector to another
	// It returns true if it succeeded or false if it wasn't found.
	Move(org *Organism, src Vector, dst Vector) (ok bool)

	// Remove attempts to remove and kill an Organism at the given Vector.
	// It returns true if it succeeded or false if it wasn't found.
	Kill(org *Organism, src Vector) (ok bool)
}

// ---------------------------------------------------------------------
// World

type World struct {
	width  int
	height int
	cells  []*Cell
}

func (w *World) GetCell(vec Vector) *Cell {
	idx := w.getIndex(vec)
	return w.cells[idx]
}

func (w *World) InBounds(vec Vector) bool {
	return w.getIndex(vec) < len(w.cells)
}

func (w *World) Walkable(vec Vector) bool {
	return w.InBounds(vec) && w.GetCell(vec).AllP(func(o *Organism) bool {
		return o.Walkable()
	})
}

// View returns all Vectors in a radius around a given origin in random order.
func (w *World) View(origin Vector, radius int) []Vector {
	vectors := w.view(origin, radius)

	shuffled := make([]Vector, len(vectors))
	for i, j := range rand.Perm(len(vectors)) {
		shuffled[i] = vectors[j]
	}
	return shuffled
}

func (w *World) view(origin Vector, radius int) []Vector {
	n := (2*radius + 1) ^ 2 - 1
	vectors := make([]Vector, n)

	i := 0
	for y := -radius; y < radius; y++ {
		for x := -radius; x < radius; x++ {
			vec := origin.Plus(Vector{x, y})
			if !vec.Equals(origin) {
				vectors[i] = vec
				i++
			}
		}
	}
	return vectors
}

// ViewWalkable is like View except it only returns walkable tiles.
func (w *World) ViewWalkable(origin Vector, radius int) []Vector {
	vectors := w.View(origin, radius)
	walkables := make([]Vector, 0)
	for i := range vectors {
		vec := vectors[i]
		if w.Walkable(vec) {
			walkables = append(walkables, vec)
		}
	}
	return walkables
}

func (w *World) RandWalkable(origin Vector, radius int) Vector {
	vectors := w.ViewWalkable(origin, radius)
	idx := rand.Intn(len(vectors))
	return vectors[idx]
}

func (w *World) Remove(organism *Organism, vec Vector) (ok bool) {
	cell := w.GetCell(vec)
	ok = cell.Remove(organism)
	return
}

func (w *World) Move(organism *Organism, src Vector, dest Vector) (ok bool) {
	oldCell := w.GetCell(src)
	newCell := w.GetCell(dest)

	newCell.Add(organism)
	ok = oldCell.Remove(organism)
	return
}

func (w *World) Kill(organism *Organism, vec Vector) (ok bool) {
	// TODO: implement corpses
	ok = w.Remove(organism, vec)
	if ok {
		organism.EndLife()
	}
	return
}

func (w *World) getIndex(vec Vector) int {
	return vec.X + (vec.Y * w.height)
}

// ---------------------------------------------------------------------
// Cell

type Cell struct {
	indexes   map[OrganismID]int
	organisms []*Organism
}

func (c *Cell) AllP(p func(o *Organism) bool) bool {
	for i := range c.organisms {
		organism := c.organisms[i]
		if !p(organism) {
			return false
		}
	}
	return true
}

func (c *Cell) Shuffled() []*Organism {
	n := len(c.organisms)
	shuffled := make([]*Organism, n)

	for i, j := range rand.Perm(n) {
		shuffled[i] = c.organisms[j]
	}
	return shuffled
}

func (c *Cell) Add(organism *Organism) {
	c.indexes[organism.id] = len(c.organisms)
	c.organisms = append(c.organisms, organism)
}

func (c *Cell) Remove(organism *Organism) (ok bool) {
	for id, index := range c.indexes {
		if id == organism.id {
			c.delOrganismByIndex(index)
			return true
		}
	}
	return false
}

func (c *Cell) delOrganismByIndex(i int) {
	copy(c.organisms[i:], c.organisms[i+1:])
	c.organisms[len(c.organisms)-1] = nil
	c.organisms = c.organisms[:len(c.organisms)-1]
}
