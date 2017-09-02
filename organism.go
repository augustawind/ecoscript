package main

type Organism struct {
	id        OrganismID
	display   rune
	behaviors []Behavior
	classes   []string
	walkable  bool
	energy    int
	size      int
	mass      int
}

type OrganismID int

var nextOrganismID OrganismID = 0

func NewOrganism() *Organism {
	organism := &Organism{id: nextOrganismID}
	nextOrganismID++
	return organism
}

func (o *Organism) Init() {
	for i := range o.behaviors {
		behavior := o.behaviors[i]
		behavior.Init(o)
	}
}

func (o *Organism) Act(world *World, origin Vector) {
	for i := range o.behaviors {
		behavior := o.behaviors[i]
		behavior.Act(world, origin)
	}
}

// ---------------------------------------------------------------------
// Behavior API.

func (o *Organism) Display() string {
	return string(o.display)
}

func (o *Organism) Biomass() int {
	return o.size * o.mass
}

func (o *Organism) Alive() bool {
	return o.energy > 0
}

func (o *Organism) Walkable() bool {
	return o.walkable
}

func (o *Organism) EndLife() {
	o.energy = 0
}

// ---------------------------------------------------------------------
// Internal API.

func (o *Organism) transfer(energy int) {
	o.energy += energy
}
