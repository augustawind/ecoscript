package main

var nothing = NewOrganism(&Attributes{
	Symbol: " ",
})

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
	var org *Organism
	stack := c.All()
	if len(stack) > 0 {
		org = stack[len(stack)-1]
	} else {
		org = nothing
	}
	return org.Display()
}

func (o *Organism) Display() string {
	return o.Attrs.Symbol
}
