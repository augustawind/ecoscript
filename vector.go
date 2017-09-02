package main

import (
	"fmt"
)

func main() {
	fmt.Println("vim-go")
}

type Vector struct {
	X int
	Y int
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
	return Vector{v.X + a.X, v.Y + a.Y}
}

func (v Vector) Minus(a Vector) Vector {
	return Vector{v.X - a.X, v.Y - a.Y}
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
