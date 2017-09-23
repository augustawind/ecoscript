package main

import (
	"log"
)

const (
	//baseActionCost int = 5
	baseActionCost int = 0
)

type OrganismID int

var (
	oid            OrganismID  = -1
	lastOrganismID *OrganismID = &oid
)

type Organism struct {
	id OrganismID

	Name      string      `mapstructure:"name"`
	Symbol    string      `mapstructure:"symbol"`
	Attrs     *Attributes `mapstructure:"attributes"`
	Traits    []Trait     `mapstructure:"traits"`
	Abilities []*Ability  `mapstructure:"abilities"`

	currentAbility int
	activity       *Activity
}

type Attributes struct {
	Walkable bool `mapstructure:"walkable"`
	Energy   int  `mapstructure:"energy"`
	Size     int  `mapstructure:"size"`
	Mass     int  `mapstructure:"mass"`
}

type Trait string

func NewOrganism(name, symbol string, attrs *Attributes) *Organism {
	abilities := make([]*Ability, 0)
	traits := make([]Trait, 0)
	activity := NewActivity()
	*lastOrganismID++
	return &Organism{
		id:        *lastOrganismID,
		Name:      name,
		Symbol:    symbol,
		Attrs:     attrs,
		Traits:    traits,
		Abilities: abilities,
		activity:  activity,
	}
}

func (o *Organism) AddAbilities(abilities ...*Ability) *Organism {
	for i := range abilities {
		ability := abilities[i]
		o.Abilities = append(o.Abilities, ability)
	}
	return o
}

func (o *Organism) AddClasses(traits ...Trait) *Organism {
	o.Traits = append(o.Traits, traits...)
	return o
}

func (o *Organism) Tick(world *World, vec Vector) {
	// If activity in progress, continue it. Otherwise, start a new activity.
	if o.activity.InProgress() {
		o.activity.Continue()
	} else {
		// Apply universal action energy cost.
		if alive := o.Transfer(-baseActionCost); !alive {
			// If energy depleted, kill and remove organism.
			execKill, ok := world.Kill(o, vec)
			if !ok {
				// TODO: figure out how to handle ok=false here
				log.Panicf("organism '%s' died, but Kill() failed unexpectedly", o.Name)
			}
			execKill()
		} else {
			// Start new activity.
			ability := o.nextAbility()
			delay, exec := ability.Execute(world, o, vec)
			o.activity.Begin(delay, exec)
		}
	}
}

func (o *Organism) nextAbility() *Ability {
	n := len(o.Abilities) - 1
	if n == 0 {
		n = 1
	}
	ability := o.Abilities[o.currentAbility%n]
	o.currentAbility++
	return ability
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
