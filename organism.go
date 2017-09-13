package main

import (
	"log"
)

const (
	baseActionCost int = -5
	baseTimeUnits  int = 10
)

type OrganismID int

var nextOrganismID OrganismID = 0

type Organism struct {
	id        OrganismID
	Attrs 	  *Attributes
	Behaviors []Behavior
	Classes   []Class
}

type Attributes struct {
	Name     string   `mapstructure:"name"`
	Symbol   rune     `mapstructure:"symbol"`
	Walkable bool     `mapstructure:"walkable"`
	Energy   int      `mapstructure:"energy"`
	Size     int      `mapstructure:"size"`
	Mass     int      `mapstructure:"mass"`
}

type Class string

func NewOrganism(attrs *Attributes) *Organism {
	behaviors := make([]Behavior, 0)
	classes := make([]Class, 0)
	organism := &Organism{
		id: nextOrganismID,
		Attrs: attrs,
		Behaviors: behaviors,
		Classes: classes,
	}
	nextOrganismID++
	return organism
}

func (o *Organism) AddBehaviors(behaviors ...Behavior) *Organism {
	o.Behaviors = append(o.Behaviors, behaviors...)
	return o
}

func (o *Organism) AddClasses(classes ...Class) *Organism {
	o.Classes = append(o.Classes, classes...)
	return o
}

func (o *Organism) Init() *Organism {
	for i := range o.Behaviors {
		behavior := o.Behaviors[i]
		behavior.Init()
	}
	return o
}

func (o *Organism) Act(world *World, vec Vector) {
	t := baseTimeUnits
	timeUnits := &t
	prevTurns := make([]Behavior, 0)

	for {
		// Apply universal action energy cost.
		if alive := o.Transfer(baseActionCost); !alive {
			execKill, ok := world.Kill(o, vec)
			if !ok {
				// TODO: figure out how to handle ok=false here
				log.Panicf("organism '%s' died, but Kill() failed unexpectedly", o.Attrs.Name)
			}
			execKill()
			break
		}
		done := o.NextMove(world, vec, timeUnits, prevTurns)
		if done {
			break
		}
	}
}

// TODO: maybe make Organism an interface so this can be more flexible?
func (o *Organism) NextMove(world *World, vec Vector, timeUnits *int, prevTurns []Behavior) (done bool) {
	// TODO: make this more interesting. This just cycles through each behavior in order.
	if len(prevTurns) == len(o.Behaviors) {
		done = true
		return
	}
	behavior := o.Behaviors[len(prevTurns)]

	// Attempt to perform behavior.
	delay, exec := behavior.Act(world, o, vec)

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
	prevTurns = append(prevTurns, behavior)
	return
}

// ---------------------------------------------------------------------
// Behavior API.

func (o *Organism) Transfer(energy int) bool {
	o.Attrs.Energy += energy
	return o.Alive()
}

func (o *Organism) Biomass() int {
	return o.Attrs.Size * o.Attrs.Mass
}

func (o *Organism) Alive() bool {
	return o.Attrs.Energy > 0
}

func (o *Organism) Walkable() bool {
	return o.Attrs.Walkable
}

func (o *Organism) EndLife() {
	o.Attrs.Energy = 0
}
