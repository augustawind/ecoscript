package main

import "math/rand"

type Behavior interface {
	Init(organism Organism)
	Act(world World, origin Vector) (energy int)
}

type baseBehavior struct {
	organism *Organism
}

func (b *baseBehavior) Init(organism *Organism) {
	b.organism = organism
}

func (b *baseBehavior) Act(world World, origin Vector) int {
	return 0
}

// ---------------------------------------------------------------------
// Grow

// Grow increases the subject's energy by its growth rate.
type Grow struct {
	*baseBehavior
	Rate int
}

func (b *Grow) Act(world World, origin Vector) (energy int) {
	energy = b.Rate
	b.organism.transfer(energy)
	return
}

// ---------------------------------------------------------------------
// Eat

// Eat attempts to consume an adjacent organism.
// If successful, the subject gains energy from the consumed organism.
type Eat struct {
	*baseBehavior
	Diet []Class
}

func (b *Eat) Act(world World, origin Vector) (energy int) {
	var vectors = world.ViewShuffled(origin, 1)

	for i := range vectors {
		vector := vectors[i]
		if !world.InBounds(vector) {
			continue
		}
		cell := world.GetCell(vector)

		orgs := cell.Shuffled()
		for j := range orgs {
			organism := orgs[j]
			if b.isEdible(organism) {
				ok = world.KillOrganism(organism, vector)
				if ok {
					energy = b.consumeBiomass(organism.Biomass())
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

func (b *Eat) consumeBiomass(biomass int) (energy int) {
	b.organism.transfer(biomass)
	return biomass
}

// ---------------------------------------------------------------------
// Move

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

// TODO: implement this on Organisms and use Behaviors to change the Delta.
type Move struct {
	*baseBehavior
	Delta  Vector
	Effort int
}

func (b *Move) Init(organism *Organism) {
	b.baseBehavior.Init(organism)
	b.randomizeDelta()
}

func (b *Move) randomizeDelta() {
	i := rand.Intn(len(directions))
	b.Delta = directions[i]
}

func (b *Move) Act(world World, origin Vector) (energy int) {
	dest := origin.Plus(b.Delta)

	if !world.Walkable(dest) {
		dest, ok := world.RandWalkable(origin, b.Delta)
		if !ok {
			return 0
		}
		b.Delta = dest.Minus(origin)
	}

	world.MoveOrganism(b.organism, origin, dest)
	energy = -b.Effort
	b.organism.transfer(energy)
	return
}

// ---------------------------------------------------------------------
// Wander

type Wander struct {
	*Move
}

func (b *Wander) Act(world World, origin Vector) (energy int) {
	energy = b.Move.Act(world, origin)
	b.randomizeDelta()
	return
}
