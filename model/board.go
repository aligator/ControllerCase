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

// LM317TConverter exists in many different breakout boards.
// This one may not necessarily match yours.
// Mine matches this
// https://www.amazon.de/DollaTek-Einstellbare-Netzteil-Converter-4-2V-40V/dp/B07DK73B7C/ref=sr_1_32?keywords=LM317&qid=1681031471&sr=8-32
func LM317TConverter() Board {
	x := 16.8
	y := 35.6
	radius := 1.6
	standoffRadius := 3.0

	cutoutWidth := 10.0
	cutoutHeight := 10.0

	// It has two holes on the long side, with a diameter of 2.74mm anda distance from the border of 1,64mm.
	return Board{
		Holes: []Hole{
			{1.64 + radius, 1.64 + radius, radius, standoffRadius, false},
			{1.64 + radius, y - 1.64 - radius, radius, standoffRadius, false},
		},
		X:      x,
		Y:      y,
		Height: 23.1,

		// It has cutouts on both short sides (top / bottom) big enough for the input and output wires.
		// The cutouts are centered.
		Cutouts: []Cutout{
			{X: x/2 - cutoutWidth/2 + 3, Y: 0, Width: cutoutWidth, Height: cutoutHeight, Side: SideTop},
			{X: x/2 - cutoutWidth/2 + 3, Y: 0, Width: cutoutWidth, Height: cutoutHeight, Side: SideBottom},
		},
	}
}

func KY001JoyIt() Board {
	x := 15.0
	y := 21.0
	radius := 1.4
	standoffRadius := 3.0

	cutoutWidth := 5.0
	cutoutHeight := 2.1

	// It has two holes on the long side, with a diameter of 2.74mm anda distance from the border of 1,64mm.
	return Board{
		Holes: []Hole{
			{2.5, 4, radius, standoffRadius, false},
			{x - 2.5, 4, radius, standoffRadius, false},
		},
		X:      x,
		Y:      y,
		Height: 11.0,

		// It has cutouts on both short sides (top / bottom) big enough for the input and output wires.
		// The cutouts are centered.
		Cutouts: []Cutout{
			{X: x/2 - cutoutWidth/2, Y: 11, Width: cutoutWidth, Height: cutoutHeight, Side: SideTop},
			{X: x/2 - cutoutWidth/2, Y: 11, Width: cutoutWidth, Height: cutoutHeight, Side: SideLeft},
			{X: x/2 - cutoutWidth/2, Y: 11, Width: cutoutWidth, Height: cutoutHeight, Side: SideRight},
			{X: x/2 - cutoutWidth/2, Y: 11, Width: cutoutWidth, Height: cutoutHeight, Side: SideBottom},
		},
	}
}
