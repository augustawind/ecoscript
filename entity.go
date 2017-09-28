package ecoscript

import "reflect"

type (
	// Entity represents an entity in the world.
	Entity struct {
		id EntityID

		Name   string      `mapstructure:"name"`
		Symbol string      `mapstructure:"symbol"`
		Attrs  *Attributes `mapstructure:"attributes"`
		Traits []Trait     `mapstructure:"traits"`

		Behaviors      Behaviors
		ChooseBehavior Strategy

		currentAbility int
		activity       *Activity
	}

	EntityID int

	Attributes struct {
		Walkable bool `mapstructure:"walkable"`
		Energy   int  `mapstructure:"energy"`
		Size     int  `mapstructure:"size"`
		Mass     int  `mapstructure:"mass"`
	}

	Trait string

	Behaviors map[string]Behavior

	Strategy func() string
)

var (
	oid          EntityID  = -1
	lastEntityID *EntityID = &oid
)

func NewEntity(name, symbol string) *Entity {
	traits := make([]Trait, 0)
	attrs := new(Attributes)
	behaviors := make(map[string]Behavior)
	activity := NewActivity()
	*lastEntityID++
	return &Entity{
		id:        *lastEntityID,
		Name:      name,
		Symbol:    symbol,
		Attrs:     attrs,
		Traits:    traits,
		Behaviors: behaviors,
		activity:  activity,
	}
}

func (e *Entity) AddAttributes(attrs *Attributes) *Entity {
	e.Attrs = attrs
	return e
}

func (e *Entity) AddTraits(traits ...Trait) *Entity {
	e.Traits = append(e.Traits, traits...)
	return e
}

func (e *Entity) AddBehaviors(behaviors ...Behavior) *Entity {
	for i := range behaviors {
		behavior := behaviors[i]
		name := reflect.TypeOf(behavior).Elem().Name()
		e.Behaviors[name] = behavior
	}
	return e
}

func (e *Entity) AddStrategy(fn Strategy) *Entity {
	e.ChooseBehavior = fn
	return e
}

func (e *Entity) Tick(world *World, vec Vector) {
	// If activity in progress, continue it. Otherwise, start a new activity.
	if e.activity.InProgress() {
		e.activity.Continue()
	} else {
		// Start new activity.
		behaviorKey := e.ChooseBehavior()
		behavior := e.Behaviors[behaviorKey]
		delay, exec := behavior.Execute(world, e, vec)
		e.activity.Begin(delay, exec)
	}
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
