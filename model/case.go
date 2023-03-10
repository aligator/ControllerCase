package model

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
	p "github.com/ljanyst/ghostscad/primitive"
)

type Side int

const (
	Top Side = iota
	Right
	Bottom
	Left
)

type Hole struct {
	X, Y             float64
	R                float64
	StandoffRadius   float64
	WithtThickBottom bool
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
	Holes     []Hole
	Cutouts   []Cutout
	X, Y      float64
	Height    float64
	tolerance float64
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

	// The cutouts need to be calculated when building.
	// The the tolerance has to be persisted.
	withTolerance.tolerance = tolerance

	return withTolerance
}

type Case struct {
	BoxPrimitive   p.Primitive
	CoverPrimitive p.Primitive

	Boards         []Board
	Wall           float64
	StandoffHeight float64
	CoverInsert    float64

	BoxHolesRadius float64

	CoverHoles bool
}

func NewCase(
	wall float64,
	standoffHeight float64,
	coverInsert float64,
	boxHolesRadius float64,
	boards ...Board,
) *Case {
	return &Case{
		Boards:         boards,
		Wall:           wall,
		StandoffHeight: standoffHeight,
		CoverInsert:    coverInsert,
		BoxHolesRadius: boxHolesRadius,
	}
}

func (o *Case) WithCoverHoles() Case {
	o.CoverHoles = true
	return *o
}

func (o *Case) GetDimensions() (x, y, height float64) {
	for _, board := range o.Boards {
		x += board.X
		y = math.Max(y, board.Y)
		height = math.Max(height, board.Height+o.StandoffHeight) // Take into account the StandoffHeight.
	}

	return x, y, height + o.CoverInsert // Add cover insert to sure the boards fit.
}

func (o *Case) buildStandoff(x float64, hole Hole) p.Primitive {
	var standoff p.Primitive = p.NewCylinder(o.StandoffHeight+o.Wall, hole.StandoffRadius).SetCenter(false)
	if hole.WithtThickBottom {
		standoff = p.NewUnion(
			standoff,
			p.NewCylinder(o.StandoffHeight/2+o.Wall, hole.StandoffRadius*2).SetCenter(false),
		)
	}
	standoff = p.NewDifference(
		standoff,
		p.NewCylinder(o.StandoffHeight+o.Wall+1, hole.R).SetCenter(false),
	)

	standoff = p.NewTranslation(mgl64.Vec3{x + hole.X + o.Wall, hole.Y + o.Wall, 0}, standoff)
	return standoff
}

func (o *Case) BuildCover() p.Primitive {
	x, y, _ := o.GetDimensions()
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
			p.NewCube(mgl64.Vec3{x - o.Wall*2, y - o.Wall*2, o.CoverInsert + 1}).SetCenter(false),
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

	cover = p.NewTranslation(mgl64.Vec3{0, 0, -o.CoverInsert + o.Wall}, cover)

	if o.CoverHoles {
		// Make the holes 2 times the wall thick.
		holeWidth := o.Wall * 2
		var count int = int(x/holeWidth) - 1

		holes := []p.Primitive{}

		for i := 0; i < count; i += 2 {
			holes = append(
				holes,
				p.NewTranslation(
					mgl64.Vec3{float64(i)*holeWidth + holeWidth, o.Wall * 2, -1},
					p.NewCube(mgl64.Vec3{holeWidth, y - 2*o.Wall, o.Wall*2 + 2}).SetCenter(false),
				),
			)
		}

		cover = p.NewDifference(
			append([]p.Primitive{cover}, holes...)...,
		)
	}

	o.CoverPrimitive = cover

	return o.CoverPrimitive
}

func (o *Case) applyCutouts(box p.Primitive) p.Primitive {
	cuts := []p.Primitive{}

	x := 0.0
	for _, board := range o.Boards {
		for _, cutout := range board.Cutouts {
			var cut p.Primitive = p.NewCube(mgl64.Vec3{cutout.Width, o.Wall + 2, cutout.Height}).SetCenter(false)

			switch cutout.Side {
			case Top:
				cut = p.NewTranslation(mgl64.Vec3{
					cutout.X + board.tolerance + o.Wall,
					board.Y + o.Wall - 1,
					o.Wall + cutout.Y,
				}, cut)
			case Right:
				cut = p.NewRotation(mgl64.Vec3{0, 0, 90}, cut)
				cut = p.NewTranslation(mgl64.Vec3{
					o.Wall*2 + 1 + board.X,
					o.Wall + cutout.X + board.tolerance,
					o.Wall + cutout.Y,
				}, cut)
			case Bottom:
				cut = p.NewTranslation(mgl64.Vec3{
					o.Wall + cutout.X + board.tolerance,
					-1,
					o.Wall + cutout.Y,
				}, cut)
			case Left:
				cut = p.NewRotation(mgl64.Vec3{0, 0, 90}, cut)
				cut = p.NewTranslation(mgl64.Vec3{
					o.Wall + 1,
					-cutout.Width + board.Y + o.Wall - board.tolerance - cutout.X,
					o.Wall + cutout.Y,
				}, cut)
			}

			cut = p.NewTranslation(mgl64.Vec3{x, 0, o.StandoffHeight}, cut)

			cuts = append(cuts, cut)
		}

		x += board.X
	}

	return p.NewDifference(append(
		[]p.Primitive{box},
		cuts...,
	)...)
}

func (o *Case) addBoxHoles(box p.Primitive) p.Primitive {
	size := o.BoxHolesRadius * 3
	var hole p.Primitive = p.NewCube(mgl64.Vec3{size, size, o.Wall}).SetCenter(false)

	hole = p.NewDifference(
		hole,
		p.NewTranslation(
			mgl64.Vec3{size / 2, size / 2},
			p.NewCylinder(o.Wall*3, o.BoxHolesRadius),
		),
	)

	_, y, _ := o.GetDimensions()

	hole = p.NewTranslation(
		mgl64.Vec3{-size, o.Wall + y/2},
		hole,
	)

	box = p.NewUnion(
		box,
		hole.Highlight(),
	)

	return box
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

	x, y, height := o.GetDimensions()
	xWithWall := x + 2*o.Wall
	yWithWall := y + 2*o.Wall
	heightWithWall := height + o.Wall

	baseBlock := p.NewCube(mgl64.Vec3{xWithWall, yWithWall, heightWithWall}).SetCenter(false)

	// Cut out the space
	cutout := p.NewTranslation(
		mgl64.Vec3{o.Wall, o.Wall, o.Wall},
		p.NewCube(mgl64.Vec3{x, y, height + 1}).SetCenter(false),
	)

	var box p.Primitive = p.NewDifference(
		baseBlock,
		cutout,
	)

	box = p.NewUnion(
		o.applyCutouts(box),
		p.NewUnion(holes...),
	)

	box = o.addBoxHoles(box)

	o.BoxPrimitive = box
	return o.BoxPrimitive
}
