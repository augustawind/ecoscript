package main

import "math/rand"

// Space is a container that has width and height and contains Cells.
type Space interface {
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

	// Add attempts to add an Organism at the given Vector.
	// It returns true if it succeeded or false if it wasn't found.
	Add(org *Organism, vec Vector) (ok bool)

	// Remove attempts to remove an Organism at the given Vector.
	// It returns true if it succeeded or false if it wasn't found.
	Remove(org *Organism, vec Vector) (ok bool)

	// Remove attempts to remove and kill an Organism at the given Vector.
	// It returns true if it succeeded or false if it wasn't found.
	Kill(org *Organism, vec Vector) (ok bool)

	// Move attempts to move an Organism from one Vector to another
	// It returns true if it succeeded or false if it wasn't found.
	Move(org *Organism, src Vector, dst Vector) (ok bool)
}

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
	occupier  *Organism
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

func (w *World) FilterWalkable(vectors []Vector) []Vector {
	walkables := make([]Vector, 0)
	for i := range vectors {
		vec := vectors[i]
		if w.Walkable(vec) {
			walkables = append(walkables, vec)
		}
	}
	return walkables
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

func (w *World) Add(organism *Organism, vec Vector) (ok bool) {
	cell := w.Cell(vec)
	return cell.Add(organism)
}

func (w *World) Remove(organism *Organism, vec Vector) (ok bool) {
	cell := w.Cell(vec)
	ok = cell.Remove(organism)
	return
}

func (w *World) Move(organism *Organism, src Vector, dest Vector) (ok bool) {
	oldCell := w.Cell(src)
	newCell := w.Cell(dest)

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

func (l *Layer) Add(organism *Organism, vec Vector) (ok bool) {
	cell := l.Cell(vec)
	return cell.Add(organism)
}

func (l *Layer) Remove(organism *Organism, vec Vector) (ok bool) {
	cell := l.Cell(vec)
	ok = cell.Remove(organism)
	return
}

func (l *Layer) Move(organism *Organism, src Vector, dest Vector) (ok bool) {
	oldCell := l.Cell(src)
	newCell := l.Cell(dest)

	newCell.Add(organism)
	ok = oldCell.Remove(organism)
	return
}

func (l *Layer) Kill(organism *Organism, vec Vector) (ok bool) {
	// TODO: implement corpses
	ok = l.Remove(organism, vec)
	if ok {
		organism.EndLife()
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

func (c *Cell) Add(organism *Organism) (ok bool) {
	if organism.Walkable() {
		if !c.Occupied() {
			c.occupier = organism
			return true
		}
	} else {
		c.depths[organism.id] = len(c.walkables)
		c.walkables = append(c.walkables, organism)
		return true
	}
	return false
}

func (c *Cell) Remove(organism *Organism) (ok bool) {
	if organism.Walkable() {
		if c.occupier.id == organism.id {
			c.occupier = nil
			return true
		}
	} else {
		for id, depth := range c.depths {
			if id == organism.id {
				c.delWalkable(depth)
				return true
			}
		}
	}
	return false
}

func (c *Cell) delWalkable(depth int) {
	copy(c.walkables[depth:], c.walkables[depth+1:])
	z := len(c.walkables) - 1
	c.walkables[z] = nil
	c.walkables = c.walkables[:z]
}
