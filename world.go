package main

import "math/rand"

// ---------------------------------------------------------------------
// World

type World struct {
	width  int
	height int
	layers []*Layer
}

type Layer struct {
	width  int
	height int
	depth  int
	name   string
	cells  []*Cell
}

type Cell struct {
	// main
	occupier *Organism
	// over
	cover     *Organism
	walkables []*Organism
	depths    map[OrganismID]int
}

//func (w *World) Display() string {
//	for y := range w.height {
//		for x := range w.width {
//
//		}
//	}
//}

func (w *World) Layer(layer int) *Layer {
	return w.layers[layer]
}

func (w *World) Cell(vec Vector) *Cell {
	index := vec.Flatten(w.height)
	return w.layers[vec.Z].cells[index]
}

func (w *World) Width() int {
	return w.width
}

func (w *World) Height() int {
	return w.height
}

func (w *World) InBounds(vec Vector) bool {
	return vec.Flatten(w.height) < w.width*w.height
}

func (w *World) Walkable(vec Vector) bool {
	return w.InBounds(vec) && !w.Cell(vec).Occupied()
}

func (w *World) View(origin Vector, radius int) []Vector {
	vectors := origin.Radius(radius)
	return VecFilter(vectors, w.InBounds)
}

func (w *World) ViewR(origin Vector, radius int) []Vector {
	vectors := origin.RadiusR(radius)
	return VecFilter(vectors, w.InBounds)
}

func (w *World) ViewWalkable(origin Vector, radius int) []Vector {
	vectors := origin.Radius(radius)
	return VecFilter(vectors, w.Walkable)
}

func (w *World) ViewWalkableR(origin Vector, radius int) []Vector {
	vectors := origin.RadiusR(radius)
	return VecFilter(vectors, w.Walkable)
}

func (w *World) RandWalkable(origin Vector, radius int) Vector {
	vectors := w.ViewWalkable(origin, radius)
	index := rand.Intn(len(vectors))
	return vectors[index]
}

func (w *World) Add(organism *Organism, vec Vector) (exec func(), ok bool) {
	cell := w.Cell(vec)
	return cell.Add(organism)
}

func (w *World) Remove(organism *Organism, vec Vector) (exec func(), ok bool) {
	cell := w.Cell(vec)
	return cell.Remove(organism)
}

func (w *World) Move(organism *Organism, src Vector, dst Vector) (exec func(), ok bool) {
	oldCell := w.Cell(src)
	newCell := w.Cell(dst)
	execAdd, okAdd := newCell.Add(organism)
	execRm, okRm := oldCell.Remove(organism)

	ok = okAdd && okRm
	if ok {
		exec = chain(execAdd, execRm)
	}
	return
}

func (w *World) Kill(organism *Organism, vec Vector) (exec func(), ok bool) {
	// TODO: implement corpses
	execRm, ok := w.Remove(organism, vec)
	if ok {
		exec = chain(execRm, organism.EndLife)
	}
	return
}

// ---------------------------------------------------------------------
// Layer

func (l *Layer) Cell(vec Vector) *Cell {
	index := vec.Flatten(l.height)
	return l.cells[index]
}

func (l *Layer) Width() int {
	return l.width
}

func (l *Layer) Height() int {
	return l.height
}

func (l *Layer) InBounds(vec Vector) bool {
	return vec.Flatten(l.height) < l.width*l.height
}

func (l *Layer) Walkable(vec Vector) bool {
	return l.InBounds(vec) && !l.Cell(vec).Occupied()
}

func (l *Layer) View(origin Vector, radius int) []Vector {
	vectors := origin.Radius(radius)
	return VecFilter(vectors, l.InBounds)
}

func (l *Layer) ViewR(origin Vector, radius int) []Vector {
	vectors := origin.RadiusR(radius)
	return VecFilter(vectors, l.InBounds)
}

func (l *Layer) ViewWalkable(origin Vector, radius int) []Vector {
	vectors := origin.Radius(radius)
	return VecFilter(vectors, l.Walkable)
}

func (l *Layer) ViewWalkableR(origin Vector, radius int) []Vector {
	vectors := origin.RadiusR(radius)
	return VecFilter(vectors, l.Walkable)
}

func (l *Layer) RandWalkable(origin Vector, radius int) Vector {
	vectors := l.ViewWalkable(origin, radius)
	index := rand.Intn(len(vectors))
	return vectors[index]
}

func (l *Layer) Add(organism *Organism, vec Vector) (exec func(), ok bool) {
	cell := l.Cell(vec)
	return cell.Add(organism)
}

func (l *Layer) Remove(organism *Organism, vec Vector) (exec func(), ok bool) {
	cell := l.Cell(vec)
	return cell.Remove(organism)
}

func (l *Layer) Move(organism *Organism, src Vector, dst Vector) (exec func(), ok bool) {
	oldCell := l.Cell(src)
	nelCell := l.Cell(dst)
	execAdd, okAdd := nelCell.Add(organism)
	execRm, okRm := oldCell.Remove(organism)

	ok = okAdd && okRm
	if ok {
		exec = chain(execAdd, execRm)
	}
	return
}

func (l *Layer) Kill(organism *Organism, vec Vector) (exec func(), ok bool) {
	// TODO: implement corpses
	execRm, ok := l.Remove(organism, vec)
	if ok {
		exec = chain(execRm, organism.EndLife)
	}
	return
}

// ---------------------------------------------------------------------
// Cell

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
	if organism.Walkable() {
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
	if organism.Walkable() {
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
	if organism.Walkable() {
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
