package main

import (
	"fmt"
)

func main() {
	fmt.Println("vim-go")
}

var Directions = []Vector{
	Vector{0, -1},
	Vector{1, -1},
	Vector{1, 0},
	Vector{1, 1},
	Vector{0, 1},
	Vector{-1, 1},
	Vector{-1, 0},
	Vector{-1, -1},
}

type Vector struct {
	X int
	Y int
}

func (v Vector) Equals(u Vector) bool {
	return v.X == u.X && v.Y == u.Y
}

func (v Vector) Compare(u Vector) int {
	nV := v.X + v.Y
	nU := u.X + u.Y
	if nV < nU {
		return -1
	} else if nV > nU {
		return 1
	}
	return 0
}

func (v Vector) Plus(u Vector) Vector {
	return Vector{v.X + u.X, v.Y + u.Y}
}

func (v Vector) Minus(u Vector) Vector {
	return Vector{v.X - u.X, v.Y - u.Y}
}

func (v Vector) Map(f func(int) int) Vector {
	return Vector{f(v.X), f(v.Y)}
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
