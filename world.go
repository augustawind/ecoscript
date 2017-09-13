package main

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

func (w *World) Width() int {
	return w.width
}

func (w *World) Height() int {
	return w.height
}

func (w *World) Cell(vec Vector) *Cell {
	index := vec.Flatten(w.height)
	return w.layers[vec.Z].cells[index]
}

func (w *World) InBounds(vec Vector) bool {
	return SpaceInBounds(w, vec)
}

func (w *World) Walkable(vec Vector) bool {
	return SpaceWalkable(w, vec)
}

func (w *World) View(origin Vector, radius int) []Vector {
	return SpaceView(w, origin, radius)
}

func (w *World) ViewR(origin Vector, radius int) []Vector {
	return SpaceViewR(w, origin, radius)
}

func (w *World) ViewWalkable(origin Vector, radius int) []Vector {
	return SpaceViewWalkable(w, origin, radius)
}

func (w *World) ViewWalkableR(origin Vector, radius int) []Vector {
	return SpaceViewWalkableR(w, origin, radius)
}

func (w *World) RandWalkable(origin Vector, radius int) Vector {
	return SpaceRandWalkable(w, origin, radius)
}

func (w *World) Add(organism *Organism, vec Vector) (exec func(), ok bool) {
	return SpaceAdd(w, organism, vec)
}

func (w *World) Remove(organism *Organism, vec Vector) (exec func(), ok bool) {
	return SpaceRemove(w, organism, vec)
}

func (w *World) Move(organism *Organism, src Vector, dst Vector) (exec func(), ok bool) {
	return SpaceMove(w, organism, src, dst)
}

func (w *World) Kill(organism *Organism, vec Vector) (exec func(), ok bool) {
	return SpaceKill(w, organism, vec)
}

// ---------------------------------------------------------------------
// Layer

func (l *Layer) Width() int {
	return l.width
}

func (l *Layer) Height() int {
	return l.height
}

func (l *Layer) Cell(vec Vector) *Cell {
	index := vec.Flatten(l.height)
	return l.cells[index]
}

func (l *Layer) InBounds(vec Vector) bool {
	return SpaceInBounds(l, vec)
}

func (l *Layer) Walkable(vec Vector) bool {
	return SpaceWalkable(l, vec)
}

func (l *Layer) View(origin Vector, radius int) []Vector {
	return SpaceView(l, origin, radius)
}

func (l *Layer) ViewR(origin Vector, radius int) []Vector {
	return SpaceViewR(l, origin, radius)
}

func (l *Layer) ViewWalkable(origin Vector, radius int) []Vector {
	return SpaceViewWalkable(l, origin, radius)
}

func (l *Layer) ViewWalkableR(origin Vector, radius int) []Vector {
	return SpaceViewWalkableR(l, origin, radius)
}

func (l *Layer) RandWalkable(origin Vector, radius int) Vector {
	return SpaceRandWalkable(l, origin, radius)
}

func (l *Layer) Add(organism *Organism, vec Vector) (exec func(), ok bool) {
	return SpaceAdd(l, organism, vec)
}

func (l *Layer) Remove(organism *Organism, vec Vector) (exec func(), ok bool) {
	return SpaceRemove(l, organism, vec)
}

func (l *Layer) Move(organism *Organism, src Vector, dst Vector) (exec func(), ok bool) {
	return SpaceMove(l, organism, src, dst)
}

func (l *Layer) Kill(organism *Organism, vec Vector) (exec func(), ok bool) {
	return SpaceKill(l, organism, vec)
}
