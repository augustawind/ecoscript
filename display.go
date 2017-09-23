package main

const (
	blankSymbol = " "
)

func (l *Layer) Display() string {
	var result string
	for y := 0; y < l.Height(); y++ {
		for x := 0; x < l.Width(); x++ {
			vec := Vec2D(x, y)
			cell := l.Cell(vec)
			result += cell.Display()
		}
		result += "\n"
	}
	return result
}

func (c *Cell) Display() string {
	orgs := c.Entities()
	if len(orgs) > 0 {
		ent := orgs[len(orgs)-1]
		return ent.Display()
	}
	return blankSymbol
}

func (e *Entity) Display() string {
	return e.Symbol
}
