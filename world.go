package main

import "math/rand"

type World struct {
	width  int
	height int
	cells  []*Cell
}

type Cell struct {
	organisms map[OrganismID]*Organism
}

func (w *World) index(v Vector) (i int, ok bool) {
	i = v.X + (v.Y * w.height)
	if i <= len(w.cells) {
		ok = true
	}
	return
}

func (w *World) GetCell(vector Vector) (cell *Cell, ok bool) {
	i, ok := w.index(vector)
	if !ok {
		// TODO: figure out if we should crash the program here
		panic("seeing if this triggers; may just need to fail silently")
		return
	}
	cell = w.cells[i]
	return
}

func (w *World) View(origin Vector, distance int) []Vector {
	n := (2*distance + 1) ^ 2 - 1
	vectors := make([]Vector, n)

	i := 0
	for y := -distance; y < distance; y++ {
		for x := -distance; x < distance; x++ {
			vec := origin.Plus(Vector{x, y})
			if !vec.Equals(origin) {
				vectors[i] = vec
				i++
			}
		}
	}
	return vectors
}

func (w *World) ViewShuffled(origin Vector, distance int) []Vector {
	vectors := w.View(origin, distance)
	n := len(vectors)

	shuffled := make([]Vector, n)
	for i, j := range rand.Perm(n) {
		shuffled[i] = vectors[j]
	}
	return shuffled
}

func (w *World) KillOrganism(vector Vector, organism *Organism) (ok bool) {
	cell, ok := w.GetCell(vector)
	if ok {
		ok = cell.Kill(organism)
	}
	return
}

func (c *Cell) Shuffled() []*Organism {
	n := len(c.organisms)
	keys := make([]OrganismID, n)

	i := 0
	for k := range c.organisms {
		keys[i] = k
		i++
	}

	shuffled := make([]*Organism, n)
	for i, j := range rand.Perm(n) {
		shuffled[i] = c.organisms[keys[j]]
	}
	return shuffled
}

func (c *Cell) Remove(organism *Organism) (ok bool) {
	for id := range c.organisms {
		if id == organism.id {
			delete(c.organisms, id)
			return true
		}
	}
	return false
}

func (c *Cell) Kill(organism *Organism) (ok bool) {
	ok = c.Remove(organism)
	if ok {
		organism.EndLife()
	}
	return
}
