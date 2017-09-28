package ecoscript

import (
	"math/rand"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/go-playground/validator.v9"
)

type Behavior interface {
	Define(Properties) Behavior
	Execute(*World, *Entity, Vector) (delay int, exec func())
}

type Properties map[string]interface{}

var behaviorValidator = validator.New()

func DefineBehavior(behavior Behavior, properties Properties) Behavior {
	var err error

	// Set custom properties.
	err = mapstructure.Decode(behavior, properties)
	Guard(err)

	// Validate Behavior.
	err = behaviorValidator.Struct(behavior)
	Guard(err)

	return behavior
}

// ---------------------------------------------------------------------
// Behavior: Grow

// Grow increases the subject's energy by its growth rate.
type Grow struct {
	Rate int `mapstructure:"rate";validate:"min=1,max=10"`
}

func (b *Grow) Define(props Properties) Behavior {
	b.Rate = 5
	return DefineBehavior(b, props)
}

func (b *Grow) Execute(wld *World, ent *Entity, vec Vector) (delay int, exec func()) {
	delay = 10
	exec = func() {
		ent.Transfer(b.rateToEnergy())
	}
	return
}

func (b *Grow) rateToEnergy() int {
	return b.Rate * 2
}

// ---------------------------------------------------------------------
// Behavior: Consume

// Consume attempts to consume an adjacent entity. If successful, the subject
// gains energy from the consumed entity.
type Consume struct {
	Diet []Trait `mapstructure:"diet"`
}

func (b *Consume) Define(props Properties) Behavior {
	b.Diet = make([]Trait, 0)
	return DefineBehavior(b, props)
}

func (b *Consume) Execute(wld *World, ent *Entity, vec Vector) (delay int, exec func()) {
	vectors := wld.View(vec, 1)

	for i := range vectors {
		vec := vectors[i]
		if !wld.InBounds(vec) {
			continue
		}
		cell := wld.Cell(vec)

		ents := cell.Shuffled()
		for j := range ents {
			entity := ents[j]
			if b.isEdible(ent) {
				execDestroy, ok := wld.Destroy(entity, vec)
				if ok {
					energy := b.biomassToEnergy(entity.Biomass())
					delay = 15
					exec = func() {
						execDestroy()
						entity.Transfer(energy)
					}
				}
				return
			}
		}
	}
	return
}

func (b *Consume) isEdible(ent *Entity) bool {
	for i := range ent.Traits {
		trait := ent.Traits[i]
		for _, subjectClass := range b.Diet {
			if trait == subjectClass {
				return true
			}
		}
	}
	return false
}

func (b *Consume) biomassToEnergy(biomass int) int {
	return -biomass
}

// ---------------------------------------------------------------------
// Behavior: Move

type Move struct {
	Dir        Vector  `mapstructure:"dir"`
	Delay      int     `mapstructure:"speed";validate:"min=1,max=30"`
	MoveRate   float32 `mapstructure:"moveRate";validate:"min=0,max=1"`
	SwitchRate float32 `mapstructure:"switchRate";validate:"min=0,max=1"`
}

func (b *Move) Define(props Properties) Behavior {
	b.Dir = b.randomDir()
	b.Delay = 10
	b.MoveRate = 1
	b.SwitchRate = 1
	return DefineBehavior(b, props)
}

func (b *Move) Execute(wld *World, ent *Entity, vec Vector) (delay int, exec func()) {
	// TODO: totally redo this to match spec
	dest := vec.Plus(b.Dir)

	if !wld.Walkable(dest) {
		dest = wld.RandWalkable(vec, 1)
		if !wld.Walkable(dest) {
			return
		}
	}

	delay = 10
	exec = func() {
		b.Dir = dest.Minus(vec)
		wld.Move(ent, vec, dest)
		ent.Transfer(10)
	}
	return
}

func (b *Move) randomDir() Vector {
	i := rand.Intn(len(directions))
	return directions[i]
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
