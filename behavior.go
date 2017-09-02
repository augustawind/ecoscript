package main

type Behavior interface {
	Init(organism Organism)
	Act(world World, origin Vector) (energy int)
}

type baseBehavior struct {
	organism Organism
}

func (b *baseBehavior) Init(organism Organism) {
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
		cell, ok := world.GetCell(vector)
		if !ok {
			continue
		}

		orgs := cell.Shuffled()
		for j := range orgs {
			organism := orgs[j]
			if b.isEdible(organism) {
				ok = world.KillOrganism(vector, organism)
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
