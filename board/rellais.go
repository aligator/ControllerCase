package board

// RelaisMakerfactory is a single 5V relais
// https://www.voelkner.de/products/2327822/MAKERFACTORY-MF-6402384-Relais-Modul-1-St.-Passend-fuer-Entwicklungskits-Arduino.html
func RelaisMakerfactory() Board {
	x := 26.5
	y := 34.0
	radius := 1.3
	standoffRadius := 3.0

	holeX, holeY := 2.6, 2.6

	return Board{
		Holes: []Hole{
			{holeX, holeY, radius, standoffRadius, false},
			{x - holeX, holeY, radius, standoffRadius, false},
			{x - holeX, y - holeY, radius, standoffRadius, false},
			{holeX, y - holeY, radius, standoffRadius, false},
		},
		Cutouts: []Cutout{
			{X: 5, Y: 1.5, Width: 10, Height: 4, Side: SideTop},
		},
		X:      x,
		Y:      y,
		Height: 17.3,
	}
}
