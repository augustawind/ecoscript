package main

import "math/rand"

type Behavior interface {
	Init(organism *Organism)
	Act(world *World, origin Vector) (delay int)
}

type baseBehavior struct {
	organism *Organism
}

func (b *baseBehavior) Init(organism *Organism) {
	b.organism = organism
}

// ---------------------------------------------------------------------
// Behavior: Grow

// Grow increases the subject's energy by its growth rate.
type Grow struct {
	*baseBehavior
	Rate int
}

func (b *Grow) Act(world *World, origin Vector) int {
	b.organism.Transfer(b.Rate)
	return 10
}

// ---------------------------------------------------------------------
// Behavior: Eat

// Eat attempts to consume an adjacent organism. If successful, the subject
// gains energy from the consumed organism.
type Eat struct {
	*baseBehavior
	Diet []string
}

func (b *Eat) Act(world *World, origin Vector) int {
	var vectors = world.View(origin, 1)

	for i := range vectors {
		vec := vectors[i]
		if !world.InBounds(vec) {
			continue
		}
		cell := world.GetCell(vec)

		orgs := cell.Shuffled()
		for j := range orgs {
			organism := orgs[j]
			if b.isEdible(organism) {
				ok := world.Kill(organism, vec)
				if ok {
					energy := b.consumeBiomass(organism.Biomass())
					b.organism.Transfer(energy)
					return 15
				}
				return 0
			}
		}
	}
	return 0
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
// Behavior: Move

var directions = []Vector{
	Vector{0, -1},
	Vector{1, -1},
	Vector{1, 0},
	Vector{1, 1},
	Vector{0, 1},
	Vector{-1, 1},
	Vector{-1, 0},
	Vector{-1, -1},
}

type Move struct {
	*baseBehavior
	Delta  Vector
	Speed  int
	Effort int
}

func (b *Move) Init(organism *Organism) {
	b.baseBehavior.Init(organism)
	b.randomizeDelta()
}

func (b *Move) Act(world *World, origin Vector) int {
	dest := origin.Plus(b.Delta)

	if !world.Walkable(dest) {
		dest = world.RandWalkable(origin, b.Speed)
		if !world.Walkable(dest) {
			return 0
		}
	}
	b.Delta = dest.Minus(origin)

	world.Move(b.organism, origin, dest)
	b.organism.Transfer(b.Effort)
	return 10
}

func (b *Move) randomizeDelta() {
	i := rand.Intn(len(directions))
	b.Delta = directions[i].Plus(Vector{b.Speed, b.Speed})
}

// ---------------------------------------------------------------------
// Behavior: Wander(Move)

type Wander struct {
	*Move
}

func (b *Wander) Act(world *World, origin Vector) int {
	b.randomizeDelta()
	return b.Move.Act(world, origin)
}
