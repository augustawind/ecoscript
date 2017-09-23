package main

import (
	"log"
)

const (
	baseActionCost int = -5
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

func (o *Organism) Act(world *World, vec Vector) {
	// Apply universal action energy cost.
	if alive := o.Transfer(baseActionCost); !alive {
		execKill, ok := world.Kill(o, vec)
		if !ok {
			// TODO: figure out how to handle ok=false here
			log.Panicf("organism '%s' died, but Kill() failed unexpectedly", o.Name)
		}
		execKill()
	}
	done := o.nextMove(world, vec)
	if done {
		return
	}
}

func (o *Organism) nextMove(world *World, vec Vector) (done bool) {
	if o.activity.InProgress() {
		// Continue activity if in progress.
		done = o.activity.Continue()
		if done {
			return true
		}
	} else {
		// Apply universal action energy cost.
		if alive := o.Transfer(baseActionCost); !alive {
			execKill, ok := world.Kill(o, vec)
			if !ok {
				// TODO: figure out how to handle ok=false here
				log.Panicf("organism '%s' died, but Kill() failed unexpectedly", o.Name)
			}
			execKill()
			return true
		}

		// Start new activity.
		ability := o.nextAbility()
		delay, exec := ability.Execute(world, o, vec)
		done = o.activity.Begin(delay, exec)
		if done {
			return true
		}
	}
	return false
}

func (o *Organism) nextAbility() *Ability {
	n := len(o.Abilities) - 1
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
