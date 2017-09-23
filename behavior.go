package main

import (
	"math/rand"
)

type Behavior interface {
	Name() string
	Defaults() Properties
	Ability(Properties) *Ability
	Execute(*Ability, *World, *Organism, Vector) (delay int, exec func())
}

type Properties map[string]interface{}

// ---------------------------------------------------------------------
// Behavior: Grow

// Grow increases the subject's energy by its growth rate.
type Grow struct{}

func (bhv *Grow) Name() string {
	return "grow"
}

func (bhv *Grow) Defaults() Properties {
	return Properties{
		"rate": 5,
	}
}

func (bhv *Grow) Ability(props Properties) *Ability {
	return NewAbility(bhv, props)
}

func (bhv *Grow) Execute(abl *Ability, wld *World, org *Organism, vec Vector) (delay int, exec func()) {
	delay = 10
	exec = func() {
		org.Transfer(abl.Get("rate").(int))
	}
	return
}

// ---------------------------------------------------------------------
// Behavior: Consume

// Consume attempts to consume an adjacent organism. If successful, the subject
// gains energy from the consumed organism.
type Consume struct{}

func (bhv *Consume) Name() string {
	return "eat"
}

func (bhv *Consume) Defaults() Properties {
	return Properties{
		"diet": make([]Trait, 0),
	}
}

func (bhv *Consume) Ability(props Properties) *Ability {
	return NewAbility(bhv, props)
}

func (bhv *Consume) Execute(abl *Ability, wld *World, org *Organism, vec Vector) (delay int, exec func()) {
	vectors := wld.View(vec, 1)

	for i := range vectors {
		vec := vectors[i]
		if !wld.InBounds(vec) {
			continue
		}
		cell := wld.Cell(vec)

		orgs := cell.Shuffled()
		for j := range orgs {
			organism := orgs[j]
			if bhv.isEdible(abl, org) {
				execDestroy, ok := wld.Destroy(organism, vec)
				if ok {
					energy := bhv.biomassToEnergy(organism.Biomass())
					delay = 15
					exec = func() {
						execDestroy()
						organism.Transfer(energy)
					}
				}
				return
			}
		}
	}
	return
}

func (bhv *Consume) isEdible(abl *Ability, org *Organism) bool {
	for i := range org.Traits {
		trait := org.Traits[i]
		for _, subjectClass := range bhv.diet(abl) {
			if trait == subjectClass {
				return true
			}
		}
	}
	return false
}

func (bhv *Consume) diet(abl *Ability) []Trait {
	return abl.Get("diet").([]Trait)
}

func (bhv *Consume) biomassToEnergy(biomass int) int {
	return -biomass
}

// ---------------------------------------------------------------------
// Behavior: Move

type Move struct{}

func (bhv *Move) Name() string {
	return "move"
}

func (bhv *Move) Defaults() Properties {
	speed := 1
	return Properties{
		"delta":  randomDelta(speed),
		"speed":  speed,
		"effort": 5,
	}
}

func (bhv *Move) Ability(props Properties) *Ability {
	return NewAbility(bhv, props)
}

func (bhv *Move) Execute(abl *Ability, wld *World, org *Organism, vec Vector) (delay int, exec func()) {
	dest := vec.Plus(abl.Get("delta").(Vector))

	if !wld.Walkable(dest) {
		dest = wld.RandWalkable(vec, abl.Get("speed").(int))
		if !wld.Walkable(dest) {
			return
		}
	}

	delay = 10
	exec = func() {
		abl.Set("delta", dest.Minus(vec))
		wld.Move(org, vec, dest)
		org.Transfer(abl.Get("effort").(int))
	}
	return
}

func (bhv *Move) randomizeDelta(abl *Ability) {
	abl.Set("delta", randomDelta(abl.Get("speed").(int)))
}

func randomDelta(speed int) Vector {
	i := rand.Intn(len(directions))
	return directions[i].Plus(Vec2D(speed, speed))
}

var directions = []Vector{
	Vec2D(0, -1),
	Vec2D(1, -1),
	Vec2D(1, 0),
	Vec2D(1, 1),
	Vec2D(0, 1),
	Vec2D(-1, 1),
	Vec2D(-1, 0),
	Vec2D(-1, -1),
}

// ---------------------------------------------------------------------
// Behavior: Wander(Move)

type Wander struct {
	*Move
}

func (bhv *Wander) Name() string {
	return "wander"
}

func (bhv *Wander) Ability(props Properties) *Ability {
	return NewAbility(bhv, props)
}

func (bhv *Wander) Execute(abl *Ability, wld *World, org *Organism, vec Vector) (delay int, exec func()) {
	bhv.randomizeDelta(abl)
	return bhv.Move.Execute(abl, wld, org, vec)
}
