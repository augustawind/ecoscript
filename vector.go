package main

import (
	"math/rand"
)

type Vector struct {
	X int
	Y int
	Z int
}

func Vec2D(x, y int) Vector {
	return Vector{X: x, Y: y, Z: -1}
}

// Is3D returns whether the Vector has a Z dimension.
func (v Vector) Is3D() bool {
	return v.Z >= 0
}

// Equals returns whether the given Vector is identical.
func (v Vector) Equals(a Vector) bool {
	return v.X == a.X && v.Y == a.Y
}

// Compare compares the Vector with another by flattening each Vector and then
// comparing ordinality. It returns 1, 0, or -1 if the Vector is greater than,
// equal to, or less than the subject, respectively.
func (v Vector) Compare(a Vector) int {
	sumV := v.Flatten(1)
	sumA := a.Flatten(1)
	if sumV < sumA {
		return -1
	} else if sumV > sumA {
		return 1
	}
	return 0
}

// Plus creates a new Vector by adding X, Y, and Z values.
func (v Vector) Plus(a Vector) Vector {
	return Vector{v.X + a.X, v.Y + a.Y, v.Z + a.Z}
}

// Minus creates a new Vector by subtracting each of the subject's X, Y, and
// Z values from each of its X, Y, and Z values, respectively.
func (v Vector) Minus(a Vector) Vector {
	return Vector{v.X - a.X, v.Y - a.Y, v.Z - a.Z}
}

// Map creates a new Vector by applying the given function to each of its
// X, Y, and Z values, respectively.
func (v Vector) Map(f func(int) int) Vector {
	return Vector{f(v.X), f(v.Y), f(v.Z)}
}

// Dir creates a new Vector by converting each of its X, Y, and Z values to
// its sign. Negatives are converted to -1, positives are converted to 1,
// and 0 is left as 0.
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

// Flatten returns the index of the Vector as if its XY grid were flattened
// into a single row, given the total number of rows in the grid.
func (v Vector) Flatten(nRows int) int {
	return v.X + (v.Y * nRows)
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

// VecFilter transforms a list of Vectors by applying a predicate function
// to each and discarding those for which it returns false.
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
