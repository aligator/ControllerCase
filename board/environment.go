package board

// K001JoyIt is a temperature sensor breakout board. https://joy-it.net/de/products/SEN-KY001TS
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
		},
	}
}

// K001JoyIt is a temperature sensor breakout board. https://joy-it.net/de/products/SEN-KY001TS
func BMP280JoyIt() Board {
	x := 15.5
	y := 11.8
	radius := 1.5
	standoffRadius := 3.0
	holeCenterToCenter := 10.0

	leftHoleX := x/2 - holeCenterToCenter/2
	rightHoleX := leftHoleX + holeCenterToCenter
	holeY := leftHoleX

	cutoutWidth := 23.0
	cutoutHeight := 5.0

	return Board{
		Holes: []Hole{
			{leftHoleX, holeY, radius, standoffRadius, false},
			{rightHoleX, holeY, radius, standoffRadius, false},
		},
		Cutouts: []Cutout{
			{X: x/2 - cutoutWidth/2, Y: 1.5, Width: cutoutWidth, Height: cutoutHeight, Side: SideTop},
		},
		X:      x,
		Y:      y,
		Height: 11.0,
	}
}

// AM2302_DHT22 is a common humidity and temperature sensor. https://www.az-delivery.de/products/dht22-temperatursensor-modul
func AM2302_DHT22() Board {
	x := 16.0
	y := 38.0
	radius := 1.5
	standoffRadius := 3.0

	holeX := x / 2
	holeY := 6.8

	cutoutWidth := 10.0
	cutoutHeight := 5.0

	return Board{
		Holes: []Hole{
			{holeX, holeY, radius, standoffRadius, false},
		},
		Cutouts: []Cutout{
			{X: x/2 - cutoutWidth/2, Y: 1.5, Width: cutoutWidth, Height: cutoutHeight, Side: SideTop},
		},
		X:      x,
		Y:      y,
		Height: 11.0,
	}
}
