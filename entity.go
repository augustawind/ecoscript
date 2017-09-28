package ecoscript

type EntityID int

var (
	oid          EntityID  = -1
	lastEntityID *EntityID = &oid
)

// Entity represents an entity in the world.
type Entity struct {
	id EntityID

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

func NewEntity(name, symbol string, attrs *Attributes) *Entity {
	abilities := make([]*Ability, 0)
	traits := make([]Trait, 0)
	activity := NewActivity()
	*lastEntityID++
	return &Entity{
		id:        *lastEntityID,
		Name:      name,
		Symbol:    symbol,
		Attrs:     attrs,
		Traits:    traits,
		Abilities: abilities,
		activity:  activity,
	}
}

func (e *Entity) AddAbilities(abilities ...*Ability) *Entity {
	for i := range abilities {
		ability := abilities[i]
		e.Abilities = append(e.Abilities, ability)
	}
	return e
}

func (e *Entity) AddClasses(traits ...Trait) *Entity {
	e.Traits = append(e.Traits, traits...)
	return e
}

func (e *Entity) Tick(world *World, vec Vector) {
	// If activity in progress, continue it. Otherwise, start a new activity.
	if e.activity.InProgress() {
		e.activity.Continue()
	} else {
		// Start new activity.
		ability := e.nextAbility()
		delay, exec := ability.Execute(world, e, vec)
		e.activity.Begin(delay, exec)
	}
}

func (e *Entity) nextAbility() *Ability {
	n := len(e.Abilities) - 1
	if n == 0 {
		n = 1
	}
	ability := e.Abilities[e.currentAbility%n]
	e.currentAbility++
	return ability
}

// ---------------------------------------------------------------------
// Behavior API.

func (e *Entity) ID() EntityID {
	return e.id
}

func (e *Entity) Transfer(energy int) bool {
	e.Attrs.Energy += energy
	return e.Alive()
}

func (e *Entity) Biomass() int {
	return e.Attrs.Size * e.Attrs.Mass
}

func (e *Entity) Alive() bool {
	return e.Attrs.Energy > 0
}

func (e *Entity) Walkable() bool {
	return e.Attrs.Walkable
}

func (e *Entity) EndLife() {
	e.Attrs.Energy = 0
}
