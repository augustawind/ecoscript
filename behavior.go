package main

import "math/rand"

type Behavior interface {
	Name() string
	Init()
	Act(world *World, organism *Organism, vec Vector) (delay int, exec func())
}

type base struct{}

func (b *base) Init() {}

// ---------------------------------------------------------------------
// Behavior: Grow

// Grow increases the subject's energy by its growth rate.
type Grow struct {
	*base
	Rate int
}

func (b *Grow) Name() string {
	return "Grow"
}

func (b *Grow) Act(world *World, organism *Organism, vec Vector) (delay int, exec func()) {
	delay = 10
	exec = func() {
		organism.Transfer(b.Rate)
	}
	return
}

// ---------------------------------------------------------------------
// Behavior: Eat

// Eat attempts to consume an adjacent organism. If successful, the subject
// gains energy from the consumed organism.
type Eat struct {
	*base
	Diet []string
}

func (b *Eat) Name() string {
	return "Eat"
}

func (b *Eat) Act(world *World, organism *Organism, vec Vector) (delay int, exec func()) {
	vectors := world.View(vec, 1)

	for i := range vectors {
		vec := vectors[i]
		if !world.InBounds(vec) {
			continue
		}
		cell := world.Cell(vec)

		orgs := cell.Shuffled()
		for j := range orgs {
			organism := orgs[j]
			if b.isEdible(organism) {
				execKill, ok := world.Kill(organism, vec)
				if ok {
					energy := b.consumeBiomass(organism.Biomass())
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

func (b *Eat) isEdible(organism *Organism) bool {
	for i := range organism.classes {
		class := organism.classes[i]
		for j := range b.Diet {
			subjectClass := b.Diet[j]
			if class == subjectClass {
				return true
			}
		}
	}
	return false
}

func (b *Eat) consumeBiomass(biomass int) int {
	return -biomass
}

// ---------------------------------------------------------------------
// Behavior: Flow

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

type Flow struct {
	*base
	Delta  Vector
	Speed  int
	Effort int
}

func (b *Flow) Name() string {
	return "Flow"
}

func (b *Flow) Init() {
	b.randomizeDelta()
}

func (b *Flow) Act(world *World, organism *Organism, vec Vector) (delay int, exec func()) {
	dest := vec.Plus(b.Delta)

	if !world.Walkable(dest) {
		dest = world.RandWalkable(vec, b.Speed)
		if !world.Walkable(dest) {
			return
		}
	}

	delay = 10
	exec = func() {
		b.Delta = dest.Minus(vec)
		world.Move(organism, vec, dest)
		organism.Transfer(b.Effort)
	}
	return
}

func (b *Flow) randomizeDelta() {
	i := rand.Intn(len(directions))
	b.Delta = directions[i].Plus(Vec2D(b.Speed, b.Speed))
}

// ---------------------------------------------------------------------
// Behavior: Wander(Flow)

type Wander struct {
	*Flow
}

func (b *Wander) Name() string {
	return "Wander"
}

func (b *Wander) Act(world *World, organism *Organism, vec Vector) (delay int, exec func()) {
	b.randomizeDelta()
	return b.Flow.Act(world, organism, vec)
}
