package main

const (
	baseActionCost int = -5
)

type Organism struct {
	id        OrganismID
	display   rune
	behaviors []Behavior
	classes   []string
	walkable  bool
	delay     int
	energy    int
	size      int
	mass      int
}

type OrganismID int

var nextOrganismID OrganismID = 0

func NewOrganism() *Organism {
	organism := &Organism{id: nextOrganismID}
	nextOrganismID++
	return organism
}

func (o *Organism) Init() {
	for i := range o.behaviors {
		behavior := o.behaviors[i]
		behavior.Init(o)
	}
}

func (o *Organism) Act(world *World, origin Vector) {
	timeUnits := 10

	for i := range o.behaviors {
		// TODO: come up with better way to decide which behavior(s) to use
		behavior := o.behaviors[i]

		// Apply universal action energy cost.
		if alive := o.Transfer(baseActionCost); !alive {
			execKill, ok := world.Kill(o, origin)
			if !ok {
				// TODO: figure out how to handle ok=false here
				panic(o)
			}
			execKill()
			break
		}

		// Act out behavior.
		delay, exec := behavior.Act(world, origin)
	}
}

// ---------------------------------------------------------------------
// Behavior API.

func (o *Organism) Transfer(energy int) bool {
	o.energy += energy
	return o.Alive()
}

func (o *Organism) Display() string {
	return string(o.display)
}

func (o *Organism) Biomass() int {
	return o.size * o.mass
}

func (o *Organism) Alive() bool {
	return o.energy > 0
}

func (o *Organism) Walkable() bool {
	return o.walkable
}

func (o *Organism) EndLife() {
	o.energy = 0
}
