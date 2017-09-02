package main

import "math/rand"

type World struct {
	things []*Thing
	height int
	width  int
}

func (w *World) Get(v Vector) (thing *Thing, ok bool) {
	i := v.X + (v.Y * w.height)
	if i >= len(w.things) {
		panic("testing")
		return
	}
	thing = w.things[i]
	if thing != nil {
		ok = true
	}
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

func (w *World) EndLifeAt(vector Vector) (ok bool) {
	thing, ok := w.Get(vector)
	if ok {
		thing.EndLife()
		w.Remove(vector)
	}
	return
}
