package main

import (
	"math/rand"
)

// ---------------------------------------------------------------------
// World

type World struct {
	width  int
	height int
	depth  int
	layers []*Layer
}

func NewWorld(width, height int, layerNames []string) *World {
	depth := len(layerNames)
	layers := make([]*Layer, depth)

	world := &World{
		width:  width,
		height: height,
		depth:  depth,
		layers: layers,
	}
	for z, name := range layerNames {
		world.addLayer(z, name)
	}
	return world
}

func (w *World) Tick() {
	// For each layer...
	for z := 0; z < w.Depth(); z++ {
		layer := w.Layer(z)

		// For each cell...
		for _, y := range rand.Perm(layer.Height()) {
			for _, x := range rand.Perm(layer.Width()) {
				vec := Vec(x, y, z)
				cell := layer.Cell(vec)
				organisms := cell.Shuffled()

				for i := range organisms {
					// Check index each iteration to account for organisms that were removed.
					if i <= cell.Population() - 1 {
						break
					}
					// Tick organism.
					org := organisms[i]
					// TODO: MAYBE (?) suspend actions until end to resolve conflicts (?)
					org.Tick(w, vec)
				}
			}
		}

		//<<< sequential iteration through world >>>
		//------------------------------------------
		//for y := 0; y < layer.Height(); y++ {
		//	for x := 0; x < layer.Width(); x++ {
		//		vec := To2D(x, y)
		//		cell := layer.Cell(vec)
		//		organisms := cell.Organisms()
		//
		//		for i := range organisms {
		//			org := organisms[i]
		//			org.Tick(w, vec)
		//		}
		//	}
		//}
		//------------------------------------------
	}
}

func (w *World) addLayer(z int, name string) *Layer {
	width := w.Width()
	height := w.Height()

	nCells := width * height
	cells := make([]*Cell, nCells)
	for i := range cells {
		cells[i] = newCell()
	}

	layer := &Layer{
		name:   name,
		width:  width,
		height: height,
		depth:  w.depth,
		cells:  cells,
	}
	w.layers[z] = layer
	return layer
}

func (w *World) Layer(z int) *Layer {
	return w.layers[z]
}

func (w *World) Width() int {
	return w.width
}

func (w *World) Height() int {
	return w.height
}

func (w *World) Depth() int {
	return w.depth
}

func (w *World) Cell(vec Vector) *Cell {
	index := vec.Flatten(w.Height())
	return w.layers[vec.Z].cells[index]
}

func (w *World) InBounds(vec Vector) bool {
	inBounds := SpaceInBounds(w, vec)
	if vec.Is3D() {
		inBounds = inBounds && vec.Z >= 0 && vec.Z < w.Depth()
	}
	return inBounds
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

func (w *World) Add(organism *Organism, vec Vector) (exec action, ok bool) {
	return SpaceAdd(w, organism, vec)
}

func (w *World) Remove(organism *Organism, vec Vector) (exec action, ok bool) {
	return SpaceRemove(w, organism, vec)
}

func (w *World) Move(organism *Organism, src Vector, dst Vector) (exec action, ok bool) {
	return SpaceMove(w, organism, src, dst)
}

func (w *World) Kill(organism *Organism, vec Vector) (exec action, ok bool) {
	return SpaceKill(w, organism, vec)
}

// ---------------------------------------------------------------------
// Layer

type Layer struct {
	width  int
	height int
	depth  int
	name   string
	cells  []*Cell
}

func (l *Layer) Width() int {
	return l.width
}

func (l *Layer) Height() int {
	return l.height
}

func (l *Layer) Cell(vec Vector) *Cell {
	index := vec.Flatten(l.Height())
	return l.cells[index]
}

func (l *Layer) Cells() []*Cell {
	return l.cells
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

func (l *Layer) Add(organism *Organism, vec Vector) (exec action, ok bool) {
	return SpaceAdd(l, organism, vec)
}

func (l *Layer) Remove(organism *Organism, vec Vector) (exec action, ok bool) {
	return SpaceRemove(l, organism, vec)
}

func (l *Layer) Move(organism *Organism, src Vector, dst Vector) (exec action, ok bool) {
	return SpaceMove(l, organism, src, dst)
}

func (l *Layer) Kill(organism *Organism, vec Vector) (exec action, ok bool) {
	return SpaceKill(l, organism, vec)
}
