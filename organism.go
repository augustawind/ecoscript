package main

import (
	"log"
)

const (
	baseActionCost int = -5
	baseTimeUnits  int = 10
)

type Attributes struct {
	Name      string     `mapstructure:"name"`
	Symbol    rune       `mapstructure:"symbol"`
	Walkable  bool       `mapstructure:"walkable"`
	Energy    int        `mapstructure:"energy"`
	Size      int        `mapstructure:"size"`
	Mass      int        `mapstructure:"mass"`
	Classes   []string   `mapstructure:"classes"`
	Behaviors []Behavior `mapstructure:"behaviors"`
}

func (a Attributes) init() {
	if a.Behaviors == nil {
		a.Behaviors = make([]Behavior, 0)
	}
	if a.Classes == nil {
		a.Classes = make([]string, 0)
	}
}

type Organism struct {
	Attributes
	id OrganismID
}

type OrganismID int

var nextOrganismID OrganismID = 0

func NewOrganism(attrs Attributes) *Organism {
	attrs.init()
	organism := &Organism{id: nextOrganismID, Attributes: attrs}
	nextOrganismID++
	return organism
}

func (o *Organism) AddBehaviors(behaviors ...Behavior) {
	o.Behaviors = append(o.Behaviors, behaviors...)
}

func (o *Organism) Init() {
	for i := range o.Behaviors {
		behavior := o.Behaviors[i]
		behavior.Init()
	}
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
				log.Panicf("organism '%s' died, but Kill() failed unexpectedly", o.Name)
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
	o.Energy += energy
	return o.Alive()
}

func (o *Organism) Biomass() int {
	return o.Size * o.Mass
}

func (o *Organism) Alive() bool {
	return o.Energy > 0
}

func (o *Organism) EndLife() {
	o.Energy = 0
}
