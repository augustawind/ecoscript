package main

const (
	emptyTileDisplay = " "
)

func (l *Layer) Display() string {
	display := ""
	for y := 0; y < l.Height(); y++ {
		for x := 0; x < l.Width(); x++ {
			vec := Vec2D(x, y)
			cell := l.Cell(vec)
			display += cell.Display()
		}
		display += "\n"
	}
	return display
}

func (c *Cell) Display() string {
	if c.cover != nil {
		return c.cover.Display()
	}
	if c.occupier != nil {
		return c.occupier.Display()
	}
	if len(c.walkables) > 0 {
		return c.walkables[len(c.walkables)-1].Display()
	}
	return emptyTileDisplay
}

func (o *Organism) Display() string {
	return o.Stats.Symbol
}
