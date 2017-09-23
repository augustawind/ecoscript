package main

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

func (o *Entity) AddAbilities(abilities ...*Ability) *Entity {
	for i := range abilities {
		ability := abilities[i]
		o.Abilities = append(o.Abilities, ability)
	}
	return o
}

func (o *Entity) AddClasses(traits ...Trait) *Entity {
	o.Traits = append(o.Traits, traits...)
	return o
}

func (o *Entity) Tick(world *World, vec Vector) {
	// If activity in progress, continue it. Otherwise, start a new activity.
	if o.activity.InProgress() {
		o.activity.Continue()
	} else {
		// Start new activity.
		ability := o.nextAbility()
		delay, exec := ability.Execute(world, o, vec)
		o.activity.Begin(delay, exec)
	}
}

func (o *Entity) nextAbility() *Ability {
	n := len(o.Abilities) - 1
	if n == 0 {
		n = 1
	}
	ability := o.Abilities[o.currentAbility%n]
	o.currentAbility++
	return ability
}

// ---------------------------------------------------------------------
// Behavior API.

func (o *Entity) ID() EntityID {
	return o.id
}

func (o *Entity) Transfer(energy int) bool {
	o.Attrs.Energy += energy
	return o.Alive()
}

func (o *Entity) Biomass() int {
	return o.Attrs.Size * o.Attrs.Mass
}

func (o *Entity) Alive() bool {
	return o.Attrs.Energy > 0
}

func (o *Entity) Walkable() bool {
	return o.Attrs.Walkable
}

func (o *Entity) EndLife() {
	o.Attrs.Energy = 0
}
