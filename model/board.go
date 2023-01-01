package model

// AmicaNodeMCU https://www.az-delivery.de/products/nodemcu
func AmicaNodeMCU() Board {
	x := 26.0
	y := 48.0
	radius := 0.85

	hole := 2.50

	return Board{
		Holes: []Hole{
			{hole, hole, radius},
			{x - hole, hole, radius},
			{x - hole, y - hole, radius},
			{hole, y - hole, radius},
		},
		X:      x,
		Y:      y,
		Height: 13,
	}
}

func SensirionSCD30() Board {
	x := 23.0
	y := 35.0
	radius := 0.85

	return Board{
		Holes: []Hole{
			{5.30, 1.70, radius},
			{x - 1.70, y - 1.75, radius},
		},
		X:      x,
		Y:      y,
		Height: 7,
	}
}
