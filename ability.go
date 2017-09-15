package main

type Ability struct {
	Name       string
	Properties Properties
}

func NewAbility(bhv Behavior, customProps Properties) *Ability {
	props := make(Properties)
	for key, defaultVal := range bhv.Defaults() {
		if customVal, ok := customProps[key]; ok {
			props[key] = customVal
		} else {
			props[key] = defaultVal
		}
	}
	return &Ability{bhv.Name(), props}
}

func (abl *Ability) Get(key string) interface{} {
	return abl.Properties[key]
}

func (abl *Ability) Set(key string, value interface{}) {
	abl.Properties[key] = value
}

func (abl *Ability) Execute(wld *World, org *Organism, vec Vector) (delay int, exec func()) {
	behavior := Behaviors[abl.Name]
	return behavior.Execute(abl, wld, org, vec)
}

type BehaviorIndex map[string]Behavior

func NewBehaviorIndex(behaviors ...Behavior) BehaviorIndex {
	index := make(BehaviorIndex)
	for i := range behaviors {
		behavior := behaviors[i]
		index[behavior.Name()] = behavior
	}
	return index
}

var Behaviors = NewBehaviorIndex(
	new(Grow),
	new(Consume),
	new(Move),
	new(Wander),
)
