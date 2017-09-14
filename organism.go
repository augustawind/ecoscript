package main

import (
	"log"
)

const (
	baseActionCost int = -5
	baseTimeUnits  int = 10
)

type OrganismID int

var (
	oid            OrganismID  = -1
	lastOrganismID *OrganismID = &oid
)

type Organism struct {
	id         OrganismID
	Attrs      *Attributes `mapstructure:"attributes"`
	Classes    []Class     `mapstructure:"classes"`
	Abilities  []*Ability  `mapstructure:"abilities"`
}

type Attributes struct {
	Name     string `mapstructure:"name"`
	Symbol   string `mapstructure:"symbol"`
	Walkable bool   `mapstructure:"walkable"`
	Energy   int    `mapstructure:"energy"`
	Size     int    `mapstructure:"size"`
	Mass     int    `mapstructure:"mass"`
}

type Class string

func NewOrganism(attrs *Attributes) *Organism {
	abilities := make([]*Ability, 0)
	classes := make([]Class, 0)
	*lastOrganismID++
	return &Organism{
		id:         *lastOrganismID,
		Attrs:      attrs,
		Classes:    classes,
		Abilities:  abilities,
	}
}

func (o *Organism) AddAbilities(abilities ...*Ability) *Organism {
	for i := range abilities {
		ability := abilities[i]
		o.Abilities = append(o.Abilities, ability)
	}
	return o
}

func (o *Organism) AddClasses(classes ...Class) *Organism {
	o.Classes = append(o.Classes, classes...)
	return o
}

func (o *Organism) Act(world *World, vec Vector) {
	t := baseTimeUnits
	timeUnits := &t

	var unusedAbilities []*Ability
	copy(unusedAbilities, o.Abilities)
	for _, ability := range o.Abilities {
		unusedAbilities = append(unusedAbilities, ability)
	}

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
		done := o.NextMove(world, vec, timeUnits, unusedAbilities)
		if done {
			break
		}
	}
}

// TODO: maybe make Organism an interface so this can be more flexible?
func (o *Organism) NextMove(world *World, vec Vector, timeUnits *int, unusedAbilities []*Ability) (done bool) {
	// TODO: make this more interesting. This just cycles through each Ability.
	if len(unusedAbilities) == 0 {
		done = true
		return
	}
	ability := unusedAbilities[0]

	// Attempt to perform ability.
	delay, exec := ability.Execute(world, o, vec)

	// Skip ability if not enough time.
	if *timeUnits-delay < 0 {
		log.Printf("ability '%s' has delay '%d' but there are only '%d' time units left", ability.Name, delay, timeUnits)
		return
	}

	// Make this the last action if time is out.
	*timeUnits -= delay
	if *timeUnits == 0 {
		done = true
	}
	exec()
	unusedAbilities = unusedAbilities[1:]
	return
}

// ---------------------------------------------------------------------
// Behavior API.

func (o *Organism) ID() OrganismID {
	return o.id
}

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
