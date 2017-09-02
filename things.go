package main

type Organism struct {
	behaviors []Behavior
	classes   []Class
	energy    int
	size      int
	mass      int
}

type Class string

func (o *Organism) transfer(energy int) {
	o.energy += energy
}

func (o *Organism) Biomass() int {
	return o.size * o.mass
}

func (o *Organism) Alive() bool {
	return o.energy > 0
}

func (o *Organism) EndLife() {
	o.energy = 0
}
