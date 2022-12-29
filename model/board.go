package model

// AmicaNodeMCU https://www.az-delivery.de/products/nodemcu
func AmicaNodeMCU() Board {
	return Board{
		Holes:  []Hole{},
		X:      26,
		Y:      48,
		Height: 13,
	}
}

func SensirionSCD30() Board {
	return Board{
		Holes:  []Hole{},
		X:      23,
		Y:      35,
		Height: 7,
	}
}
