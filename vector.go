package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println("vim-go")
}

type Vector struct {
	X int
	Y int
	Z int
}

func Vec2D(x, y int) Vector {
	return Vector{X: x, Y: y, Z: -1}
}

func (v Vector) Equals(a Vector) bool {
	return v.X == a.X && v.Y == a.Y
}

func (v Vector) Compare(a Vector) int {
	sumV := v.X + v.Y
	sumA := a.X + a.Y
	if sumV < sumA {
		return -1
	} else if sumV > sumA {
		return 1
	}
	return 0
}

func (v Vector) Plus(a Vector) Vector {
	return Vector{v.X + a.X, v.Y + a.Y, v.Z + a.Z}
}

func (v Vector) Minus(a Vector) Vector {
	return Vector{v.X - a.X, v.Y - a.Y, v.Z - a.Z}
}

func (v Vector) Map(f func(int) int) Vector {
	return Vector{f(v.X), f(v.Y), f(v.Z)}
}

func (v Vector) Dir() Vector {
	return v.Map(func(n int) int {
		if n > 0 {
			return 1
		} else if n < 0 {
			return -1
		}
		return 0
	})
}

func (v Vector) Flatten(n int) int {
	return v.X + (v.Y * n)
}

// Radius returns the surrounding Vectors by the given radius.
func (v Vector) Radius(radius int) []Vector {
	n := (2*radius + 1) ^ 2 - 1
	vectors := make([]Vector, n)

	i := 0
	for y := -radius; y < radius; y++ {
		for x := -radius; x < radius; x++ {
			vec := v.Plus(Vec2D(x, y))
			if !vec.Equals(v) {
				vectors[i] = vec
				i++
			}
		}
	}
	return vectors
}

// RadiusR is like Radius but randomizes the returned Vectors.
func (v Vector) RadiusR(radius int) []Vector {
	vectors := v.Radius(radius)

	shuffled := make([]Vector, len(vectors))
	for i, j := range rand.Perm(len(vectors)) {
		shuffled[i] = vectors[j]
	}
	return shuffled
}

func VecFilter(vectors []Vector, pred func(Vector) bool) []Vector {
	result := make([]Vector, 0)
	for i := range vectors {
		vec := vectors[i]
		if pred(vec) {
			result = append(result, vec)
		}
	}
	return result
}
