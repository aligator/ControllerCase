package board

// AmicaNodeMCU https://www.az-delivery.de/products/nodemcu
func AmicaNodeMCU() Board {
	x := 25.5
	y := 48.0
	radius := 1.5
	standoffRadius := 2.5

	holeX := 2.50
	holeY := 2.30

	return Board{
		Holes: []Hole{
			{holeX, holeY, radius, standoffRadius, false},
			{x - holeX, holeY, radius, standoffRadius, false},
			{x - holeX, y - holeY, radius, standoffRadius, false},
			{holeX, y - holeY, radius, standoffRadius, false},
		},
		Cutouts: []Cutout{
			{X: x/2 - usbWidth/2, Y: 1.5, Width: usbWidth, Height: usbHeight, Side: SideTop},
		},
		X:      x,
		Y:      y,
		Height: 13,
	}
}

// LolinV3NodeMCU https://www.azdelivery.de/products/nodemcu-lolin-v3-modul-mit-esp8266
func LolinV3NodeMCU() Board {
	x := 31.0
	y := 58.0
	radius := 1.6
	standoffRadius := 2.6

	holeX := 2.50
	holeY := 2.30

	return Board{
		Holes: []Hole{
			{holeX, holeY, radius, standoffRadius, false},
			{x - holeX, holeY, radius, standoffRadius, false},
			{x - holeX, y - holeY, radius, standoffRadius, false},
			{holeX, y - holeY, radius, standoffRadius, false},
		},
		Cutouts: []Cutout{
			{X: x/2 - usbWidth/2, Y: 1.5, Width: usbWidth, Height: usbHeight, Side: SideTop},
		},
		X:      x,
		Y:      y,
		Height: 13,
	}
}
