package board

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
			{2.8, 2.8, radius, standoffRadius, false},
			{2.8, y - 2.8, radius, standoffRadius, false},
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
