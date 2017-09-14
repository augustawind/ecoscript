package main

import (
	"math/rand"
)

type BehaviorIndex map[string]Behavior

func NewBehaviorIndex(behaviors ...Behavior) BehaviorIndex {
	index := make(BehaviorIndex)
	for i := range behaviors {
		behavior := behaviors[i]
		index[behavior.Name()] = behavior
	}
	return index
}

var Behaviors = NewBehaviorIndex(
	new(Grow),
	new(Eat),
	new(Flow),
	new(Wander),
)

type Behavior interface {
	Name() string
	Defaults() Properties
	Ability(Properties) *Ability
	Execute(*Ability, *World, *Organism, Vector) (delay int, exec func())
}

type Ability struct {
	Name       string
	Properties Properties
}

type Properties map[string]interface{}

func NewAbility(bhv Behavior, customProps Properties) *Ability {
	props := make(Properties)
	for key, defaultVal := range bhv.Defaults() {
		if customVal, ok := customProps[key]; ok {
			props[key] = customVal
		} else {
			props[key] = defaultVal
		}
	}
	return &Ability{bhv.Name(), props}
}

func (abl *Ability) Get(key string) interface{} {
	return abl.Properties[key]
}

func (abl *Ability) Set(key string, value interface{}) {
	abl.Properties[key] = value
}

func (abl *Ability) Execute(wld *World, org *Organism, vec Vector) (delay int, exec func()) {
	behavior := Behaviors[abl.Name]
	return behavior.Execute(abl, wld, org, vec)
}

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
// Behavior: Eat

// Eat attempts to consume an adjacent organism. If successful, the subject
// gains energy from the consumed organism.
type Eat struct{}

func (bhv *Eat) Name() string {
	return "eat"
}

func (bhv *Eat) Defaults() Properties {
	return Properties{}
}

func (bhv *Eat) Ability(props Properties) *Ability {
	return NewAbility(bhv, props)
}

func (bhv *Eat) Execute(abl *Ability, wld *World, org *Organism, vec Vector) (delay int, exec func()) {
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
				execKill, ok := wld.Kill(organism, vec)
				if ok {
					energy := bhv.consumeBiomass(organism.Biomass())
					delay = 15
					exec = func() {
						execKill()
						organism.Transfer(energy)
					}
				}
				return
			}
		}
	}
	return
}

func (bhv *Eat) isEdible(abl *Ability, org *Organism) bool {
	for i := range org.Classes {
		class := org.Classes[i]
		for _, subjectClass := range bhv.diet(abl) {
			if class == subjectClass {
				return true
			}
		}
	}
	return false
}

func (bhv *Eat) diet(abl *Ability) []Class {
	return abl.Get("diet").([]Class)
}

func (bhv *Eat) consumeBiomass(biomass int) int {
	return -biomass
}

// ---------------------------------------------------------------------
// Behavior: Flow

type Flow struct{}

func (bhv *Flow) Name() string {
	return "flow"
}

func (bhv *Flow) Defaults() Properties {
	speed := 1
	return Properties{
		"delta":  randomDelta(speed),
		"speed":  speed,
		"effort": 5,
	}
}

func (bhv *Flow) Ability(props Properties) *Ability {
	return NewAbility(bhv, props)
}

func (bhv *Flow) Execute(abl *Ability, wld *World, org *Organism, vec Vector) (delay int, exec func()) {
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

func (bhv *Flow) randomizeDelta(abl *Ability) {
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
// Behavior: Wander(Flow)

type Wander struct {
	*Flow
}

func (bhv *Wander) Name() string {
	return "wander"
}

func (bhv *Wander) Ability(props Properties) *Ability {
	return NewAbility(bhv, props)
}

func (bhv *Wander) Execute(abl *Ability, wld *World, org *Organism, vec Vector) (delay int, exec func()) {
	bhv.randomizeDelta(abl)
	return bhv.Flow.Execute(abl, wld, org, vec)
}
