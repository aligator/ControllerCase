package model

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
	p "github.com/ljanyst/ghostscad/primitive"
)

type Hole struct {
	X, Y           float64
	R              float64
	StandoffRadius float64
}

type Board struct {
	Holes  []Hole
	X, Y   float64
	Height float64
}

// WithTolerance returns a copy of the board, with the given tolerance added to all sides.
// THe holes get adapted accordingly
func (b Board) WithTolerance(tolerance float64) Board {
	withTolerance := b
	withTolerance.Height += tolerance * 2
	withTolerance.X += tolerance * 2
	withTolerance.Y += tolerance * 2

	// Modify holes.
	for i := range withTolerance.Holes {
		withTolerance.Holes[i].X += tolerance
		withTolerance.Holes[i].Y += tolerance
	}

	return withTolerance
}

type Case struct {
	Primitive p.Primitive

	Boards         []Board
	Wall           float64
	StandoffHeight float64
	CoverInsert    float64
}

func NewCase(wall float64, coverInsert float64, boards ...Board) *Case {
	return &Case{
		Boards:         boards,
		Wall:           wall,
		StandoffHeight: 1,
		CoverInsert:    coverInsert,
	}
}

func (o *Case) getDimensions() (x, y, height float64) {
	for _, board := range o.Boards {
		x += board.X
		y = math.Max(y, board.Y)
		height = math.Max(height, board.Height+o.StandoffHeight) // Take into account the StandoffHeight.
	}

	return x, y, height
}

func (o *Case) buildStandoff(x float64, hole Hole) p.Primitive {
	var standoff p.Primitive = p.NewCylinder(o.StandoffHeight, hole.StandoffRadius).SetCenter(false)

	standoff = p.NewDifference(
		standoff,
		p.NewCylinder(o.StandoffHeight+1, hole.R).SetCenter(false),
	)

	standoff = p.NewTranslation(mgl64.Vec3{x + hole.X + o.Wall, hole.Y + o.Wall, o.Wall}, standoff)
	return standoff
}

func (o *Case) BuildCover() p.Primitive {
	x, y, _ := o.getDimensions()
	xWithWall := x + 2*o.Wall
	yWithWall := y + 2*o.Wall

	// The height is just the wall thickness*2 and then one wall thickness is cut out again.
	heightWithWall := o.CoverInsert + o.Wall

	var cover p.Primitive = p.NewCube(mgl64.Vec3{xWithWall, yWithWall, heightWithWall}).SetCenter(false)
	cover = p.NewDifference(
		cover,

		// Cut out the inner part.
		p.NewTranslation(
			mgl64.Vec3{o.Wall * 2, o.Wall * 2, -1},
			p.NewCube(mgl64.Vec3{x - o.Wall*2, y - o.Wall*2, o.Wall + 1}).SetCenter(false),
		),
	)

	// Cut out the outer part.
	// 1. Design the negative.
	var outerNegative p.Primitive = p.NewCube(mgl64.Vec3{xWithWall + 2, yWithWall + 2, o.CoverInsert + 1}).SetCenter(false)
	outerNegative = p.NewDifference(
		p.NewTranslation(
			mgl64.Vec3{-1, -1, -1},
			outerNegative,
		),

		// Cut out the inner part.
		p.NewTranslation(
			mgl64.Vec3{o.Wall, o.Wall},
			p.NewCube(mgl64.Vec3{x, y, o.Wall + 2}).SetCenter(false),
		),
	)

	// 2. Cut the negative from the cover.
	cover = p.NewDifference(
		cover,
		outerNegative,
	)

	return cover
}

func (o *Case) BuildBox() p.Primitive {
	holes := []p.Primitive{}

	// Build the base-block.
	var x float64
	for _, board := range o.Boards {
		// Add hole standoffs.
		for _, hole := range board.Holes {
			holes = append(holes, o.buildStandoff(x, hole))
		}

		x += board.X
	}

	x, y, height := o.getDimensions()
	xWithWall := x + 2*o.Wall
	yWithWall := y + 2*o.Wall
	heightWithWall := height + o.Wall

	baseBlock := p.NewCube(mgl64.Vec3{xWithWall, yWithWall, heightWithWall}).SetCenter(false)

	// Cut out the space
	cutout := p.NewTranslation(
		mgl64.Vec3{o.Wall, o.Wall, o.Wall},
		p.NewCube(mgl64.Vec3{x, y, height + 1}).SetCenter(false),
	)

	board := p.NewDifference(
		baseBlock,
		cutout,
	)

	board = p.NewUnion(
		board,
		p.NewUnion(holes...),
	)

	o.Primitive = board
	return o.Primitive
}
