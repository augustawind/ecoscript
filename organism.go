package main

import (
	"log"

	"github.com/pkg/errors"
)

const (
	baseActionCost int = -5
	baseTimeUnits  int = 10
)

type Attributes struct {
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

func (a Attributes) validateNonempty(name string, attr string) {
	if len(attr) == 0 {
		err := errors.Errorf("attribute '%s' must not be blank", name)
		log.Panic(err)
	}
}

func (a Attributes) validatePositive(name string, attr int) {
	if attr == 0 {
		err := errors.Errorf("attribute '%s' must be greater than 0", name)
		log.Panic(err)
	}
}

func (a Attributes) init() {
	a.validateNonempty("name", a.name)
	a.validatePositive("display", int(a.display))
	a.validatePositive("delay", a.delay)
	a.validatePositive("energy", a.energy)
	a.validatePositive("size", a.size)
	a.validatePositive("mass", a.mass)

	if a.behaviors == nil {
		a.behaviors = make([]Behavior, 0)
	}
	if a.classes == nil {
		a.classes = make([]string, 0)
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
	o.behaviors = append(o.behaviors, behaviors...)
}

func (o *Organism) Init() {
	for i := range o.behaviors {
		behavior := o.behaviors[i]
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
				log.Panicf("organism '%s' died, but Kill() failed unexpectedly", o.name)
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
	if len(prevTurns) == len(o.behaviors) {
		done = true
		return
	}
	behavior := o.behaviors[len(prevTurns)]

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
	o.energy += energy
	return o.Alive()
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
