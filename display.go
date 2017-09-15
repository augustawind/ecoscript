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
	orgs := c.Organisms()
	if len(orgs) > 0 {
		org := orgs[len(orgs)-1]
		return org.Display()
	}
	return blankSymbol
}

func (o *Organism) Display() string {
	return o.Symbol
}
