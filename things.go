package main

type Thing struct {
	behaviors []Behavior
	classes   []Class
	energy    int
	size      int
	mass      int
}

type Class string

func (t *Thing) transfer(energy int) {
	t.energy += energy
}

func (t *Thing) Biomass() int {
	return t.size * t.mass
}
