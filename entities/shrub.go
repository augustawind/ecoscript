package entities

import (
	es "github.com/dustinrohde/ecoscript"
)

func NewShrub() *es.Entity {
	return es.NewEntity(
		"shrub",
		"*",
	).AddAttributes(&es.Attributes{
		Walkable: true,
		Energy:   100,
		Size:     1,
		Mass:     3,
	}).AddBehaviors(
		new(es.Grow).Define(es.Properties{
			"rate": 3,
		}),
	).AddStrategy(func() string {
		return "Grow"
	})
}
