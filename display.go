package ecoscript

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
	ents := c.Entities()
	if len(ents) > 0 {
		ent := ents[len(ents)-1]
		return ent.Display()
	}
	return blankSymbol
}

func (e *Entity) Display() string {
	return e.Symbol
}
