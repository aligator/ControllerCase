package board

const (
	usbWidth  = 10.0
	usbHeight = 4.0
)

type Side int

const (
	SideTop Side = iota
	SideRight
	SideBottom
	SideLeft
	SideX
	SideY
	SideAll
)

type Position int

const (
	PositionBottom Position = iota
	PositionCenter
	PositionTop
)

type Hole struct {
	X, Y            float64
	R               float64
	StandoffRadius  float64
	WithThickBottom bool
}

// Cutout defines a hole in the wall.
// X and Y are always relative from the side, not from the global axes.
// So if Side is "Left", you have to look from the left side onto the box
// and then x is from bottom-left to bottom-right,
// and y is from bottom to top.
type Cutout struct {
	X, Y          float64
	Width, Height float64
	Side          Side
}

type Board struct {
	Holes    []Hole
	Cutouts  []Cutout
	Position Position
	X, Y     float64
	Height   float64
	Padding  [4]float64
}

// WithPadding returns a copy of the board, with the given padding added to all sides.
// THe holes get adapted accordingly
func (b Board) WithPadding(padding float64, side Side) Board {
	withPadding := b

	newPadding := [4]float64{0, 0, 0, 0}

	switch side {
	case SideTop:
		fallthrough
	case SideRight:
		fallthrough
	case SideBottom:
		fallthrough
	case SideLeft:
		newPadding[side] = padding
	case SideX:
		newPadding[SideLeft] = padding
		newPadding[SideRight] = padding
	case SideY:
		newPadding[SideTop] = padding
		newPadding[SideBottom] = padding
	default:
		newPadding[SideTop] = padding
		newPadding[SideRight] = padding
		newPadding[SideBottom] = padding
		newPadding[SideLeft] = padding
	}

	withPadding.X += newPadding[SideLeft] + newPadding[SideRight]
	withPadding.Y += newPadding[SideTop] + newPadding[SideBottom]

	// Modify holes.
	for i := range withPadding.Holes {
		withPadding.Holes[i].X += newPadding[SideLeft]
		withPadding.Holes[i].Y += newPadding[SideBottom]
	}

	// Persist the padding for further calculations.
	withPadding.Padding[SideTop] += newPadding[SideTop]
	withPadding.Padding[SideRight] += newPadding[SideRight]
	withPadding.Padding[SideBottom] += newPadding[SideBottom]
	withPadding.Padding[SideLeft] += newPadding[SideLeft]

	return withPadding
}

func (b Board) WithCutout(cutout Cutout) Board {
	b.Cutouts = append(b.Cutouts, cutout)
	return b
}

func (b Board) WithPosition(position Position) Board {
	b.Position = position
	return b
}

func (b Board) WithAdditionalHeight(height float64) Board {
	b.Height += height
	return b
}

func (b Board) WithoutCutouts() Board {
	b.Cutouts = []Cutout{}
	return b
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
