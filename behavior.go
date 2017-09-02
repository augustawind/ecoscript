package main

type Behavior interface {
	Init(thing Thing)
	Act(world World, origin Vector) (energy int)
}

type baseBehavior struct {
	thing Thing
}

func (b *baseBehavior) Init(thing Thing) {
	b.thing = thing
}

func (b *baseBehavior) Act(world World, origin Vector) int {
	return 0
}

// ---------------------------------------------------------------------
// Grow

// Grow increases a Thing's energy over time.
type Grow struct {
	*baseBehavior
	Rate int
}

func (b *Grow) Act(world World, origin Vector) (energy int) {
	energy = b.Rate
	b.thing.transfer(energy)
	return
}

// ---------------------------------------------------------------------
// Eat

type Eat struct {
	*baseBehavior
	Diet []Class
}

func (b *Eat) Act(world World, origin Vector) (energy int) {
	var vectors = world.ViewShuffled(origin, 1)

	for i := range vectors {
		target, ok := world.Get(vectors[i])

		if ok && b.isEdible(target) {
			energy = b.transferBiomass(target.Biomass())
			world.Kill(target)
			return
		}
	}
	return
}

func (b *Eat) isEdible(target *Thing) bool {
	var targetClass Class
	var subjectClass Class

	for i := range target.classes {
		targetClass = target.classes[i]
		for j := range b.Diet {
			subjectClass = b.Diet[j]
			if targetClass == subjectClass {
				return true
			}
		}
	}
	return false
}

func (b *Eat) transferBiomass(biomass int) (energy int) {
	b.thing.transfer(biomass)
	return biomass
}
