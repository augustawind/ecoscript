package main

import (
	"math/rand"
)

// Space is a container that has width and height and contains Cells.
type Space interface {
	// Width returns the length of the x-axis.
	Width() int

	// Height returns the length of the y-axis.
	Height() int

	// Cell returns the Cell at the given vector.
	Cell(vec Vector) *Cell

	// InBounds returns true if the given Vector is in bounds.
	InBounds(vec Vector) bool

	// Walkable returns true if the given Vector is walkable.
	Walkable(vec Vector) bool

	// View returns all Vectors that are in bounds and within a radius.
	View(origin Vector, radius int) []Vector

	// ViewR is like View but randomizes the returned Vectors.
	ViewR(origin Vector, radius int) []Vector

	// ViewWalkable returns all walkable Vectors within a radius.
	ViewWalkable(origin Vector, radius int) []Vector

	// ViewWalkableR is like ViewWalkable but randomizes the returned Vectors.
	ViewWalkableR(origin Vector, radius int) []Vector

	// RandWalkable finds a random walkable Vector within a radius.
	RandWalkable(origin Vector, radius int) Vector

	// Add attempts to add an Entity at the given Vector.
	// It returns true if it succeeded or false if it wasn't found.
	Add(ent *Entity, vec Vector) (action, bool)

	// Remove attempts to remove an Entity at the given Vector.
	// It returns true if it succeeded or false if it wasn't found.
	Remove(ent *Entity, vec Vector) (action, bool)

	// Destroy attempts to remove and destroy an Entity at the given Vector.
	// It returns true if it succeeded or false if it wasn't found.
	Destroy(ent *Entity, vec Vector) (action, bool)

	// Move attempts to move an Entity from one Vector to another
	// It returns true if it succeeded or false if it wasn't found.
	Move(ent *Entity, src Vector, dst Vector) (action, bool)
}

func SpaceInBounds(s Space, vec Vector) bool {
	return vec.Flatten(s.Height()) < s.Width()*s.Height()
}

func SpaceWalkable(s Space, vec Vector) bool {
	return s.InBounds(vec) && !s.Cell(vec).Occupied()
}

func SpaceView(s Space, origin Vector, radius int) []Vector {
	vectors := origin.Radius(radius)
	return VecFilter(vectors, s.InBounds)
}

func SpaceViewR(s Space, origin Vector, radius int) []Vector {
	vectors := origin.RadiusR(radius)
	return VecFilter(vectors, s.InBounds)
}

func SpaceViewWalkable(s Space, origin Vector, radius int) []Vector {
	vectors := origin.Radius(radius)
	return VecFilter(vectors, s.Walkable)
}

func SpaceViewWalkableR(s Space, origin Vector, radius int) []Vector {
	vectors := origin.RadiusR(radius)
	return VecFilter(vectors, s.Walkable)
}

func SpaceRandWalkable(s Space, origin Vector, radius int) Vector {
	vectors := s.ViewWalkable(origin, radius)
	index := rand.Intn(len(vectors))
	return vectors[index]
}

func SpaceAdd(s Space, entity *Entity, vec Vector) (exec action, ok bool) {
	cell := s.Cell(vec)
	return cell.Add(entity)
}

func SpaceRemove(s Space, entity *Entity, vec Vector) (exec action, ok bool) {
	cell := s.Cell(vec)
	return cell.Remove(entity)
}

func SpaceMove(s Space, entity *Entity, src Vector, dst Vector) (exec action, ok bool) {
	oldCell := s.Cell(src)
	newCell := s.Cell(dst)
	execAdd, okAdd := newCell.Add(entity)
	execRm, okRm := oldCell.Remove(entity)

	ok = okAdd && okRm
	if ok {
		exec = chain(execAdd, execRm)
	}
	return
}

func SpaceDestroy(s Space, entity *Entity, vec Vector) (exec action, ok bool) {
	// TODO: implement corpses
	exec, ok = s.Remove(entity, vec)
	if ok {
		exec = chain(exec, entity.EndLife)
	}
	return
}

type action func()

func chain(actions ...action) action {
	return func() {
		for _, action := range actions {
			if action != nil {
				action()
			}
		}
	}
}
