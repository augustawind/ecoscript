package main

import (
	"log"
)

const (
	baseActionCost int = -5
	baseTimeUnits  int = 10
)

type Organism struct {
	id        OrganismID
	name      string
	display   rune
	behaviors []Behavior
	classes   []string
	walkable  bool
	delay     int
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
	t := baseTimeUnits
	timeUnits := &t
	prevTurns := make([]*Behavior, 0)

	for {
		// Apply universal action energy cost.
		if alive := o.Transfer(baseActionCost); !alive {
			execKill, ok := world.Kill(o, origin)
			if !ok {
				// TODO: figure out how to handle ok=false here
				log.Panicf("organism '%s' died, but Kill() failed unexpectedly", o.name)
			}
			execKill()
			break
		}
		done := o.NextMove(world, origin, timeUnits, prevTurns)
		if done {
			break
		}
	}
}

// TODO: maybe make Organism an interface so this can be more flexible?
func (o *Organism) NextMove(world *World, origin Vector, timeUnits *int, prevTurns []*Behavior) (done bool) {
	// TODO: make this more interesting. This just cycles through each behavior in order.
	if len(prevTurns) == len(o.behaviors) {
		done = true
		return
	}
	behavior := o.behaviors[len(prevTurns)]

	// Attempt to perform behavior.
	delay, exec := behavior.Act(world, origin)

	// Skip behavior if not enough time.
	if *timeUnits-delay < 0 {
		log.Printf("behavior '%s' has delay '%d' but there are only '%d' time units left", behavior.Name(), delay, timeUnits)
		return
	}

	// Make this the last action if time is out.
	*timeUnits -= delay
	if *timeUnits == 0 {
		done = true
	}
	exec()
	prevTurns = append(prevTurns, &behavior)
	return
}

// ---------------------------------------------------------------------
// Behavior API.

func (o *Organism) Transfer(energy int) bool {
	o.energy += energy
	return o.Alive()
}

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
