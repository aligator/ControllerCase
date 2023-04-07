package model

// AmicaNodeMCU https://www.az-delivery.de/products/nodemcu
func AmicaNodeMCU() Board {
	x := 25.5
	y := 48.0
	radius := 1.5
	standoffRadius := 2.5

	holeX := 2.50
	holeY := 2.30

	usbWidth := 10.0

	return Board{
		Holes: []Hole{
			{holeX, holeY, radius, standoffRadius, false},
			{x - holeX, holeY, radius, standoffRadius, false},
			{x - holeX, y - holeY, radius, standoffRadius, false},
			{holeX, y - holeY, radius, standoffRadius, false},
		},
		Cutouts: []Cutout{
			{X: x/2 - usbWidth/2, Y: 1.5, Width: usbWidth, Height: 4, Side: SideTop},
		},
		X:      x,
		Y:      y,
		Height: 13,
	}
}

func SensirionSCD30() Board {
	x := 23.0
	y := 35.0
	radius := 1.0
	standoffRadius := 1.7

	return Board{
		Holes: []Hole{
			{5.30, 1.70, radius, standoffRadius, true},
			{x - 1.70, y - 1.75, radius, standoffRadius, true},
		},
		X:      x,
		Y:      y,
		Height: 7,
	}
}
