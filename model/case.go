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
	Holes   []Hole
	Cutouts []Cutout
	X, Y    float64
	Height  float64
	padding float64
}

// WithPadding returns a copy of the board, with the given padding added to all sides.
// THe holes get adapted accordingly
func (b Board) WithPadding(padding float64) Board {
	withPadding := b
	withPadding.X += padding * 2
	withPadding.Y += padding * 2

	// Modify holes.
	for i := range withPadding.Holes {
		withPadding.Holes[i].X += padding
		withPadding.Holes[i].Y += padding
	}

	// The cutouts need to be calculated when building.
	// The the padding has to be persisted for this.
	withPadding.padding = padding

	return withPadding
}

func (b Board) WithCutout(cutout Cutout) Board {
	b.Cutouts = append(b.Cutouts, cutout)
	return b
}

type Case struct {
	BoxPrimitive   p.Primitive
	CoverPrimitive p.Primitive

	Boards         []Board
	Wall           float64
	StandoffHeight float64
	CoverInsert    float64
	HeightPadding  float64

	MountingHolesRadius float64

	CoverHoles bool
}

func NewCase(
	wall float64,
	standoffHeight float64,
	coverInsert float64,
	mountingHolesRadius float64,
	boards ...Board,
) *Case {
	return &Case{
		Boards:              boards,
		Wall:                wall,
		StandoffHeight:      standoffHeight,
		CoverInsert:         coverInsert,
		MountingHolesRadius: mountingHolesRadius,
	}
}

func (o *Case) SetCoverHoles() *Case {
	o.CoverHoles = true
	return o
}

func (o *Case) SetHeightPadding(padding float64) *Case {
	o.HeightPadding = padding
	return o
}

func (o *Case) GetDimensions(withWalls bool) (x, y, height float64) {
	for _, board := range o.Boards {
		x += board.X
		y = math.Max(y, board.Y)
		height = math.Max(height, board.Height+o.StandoffHeight) // Take into account the StandoffHeight.
	}

	if withWalls {
		x += 2 * o.Wall
		y += 2 * o.Wall
		height += o.Wall
	}

	return x, y, height + o.CoverInsert + o.HeightPadding // Add cover insert to be sure the boards fits.
}

func (o *Case) buildStandoff(x float64, hole Hole) p.Primitive {
	var standoff p.Primitive = p.NewCylinder(o.StandoffHeight+o.Wall, hole.StandoffRadius).SetCenter(false)
	if hole.WithThickBottom {
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
	x, y, _ := o.GetDimensions(false)
	xWithWall, yWithWall, _ := o.GetDimensions(true)

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

	_, dimY, _ := o.GetDimensions(false)

	x := 0.0
	for _, board := range o.Boards {
		for _, cutout := range board.Cutouts {
			var cut p.Primitive = p.NewCube(mgl64.Vec3{cutout.Width, o.Wall + 2, cutout.Height}).SetCenter(false)

			switch cutout.Side {
			case Top:
				cut = p.NewTranslation(mgl64.Vec3{
					cutout.X + board.padding + o.Wall,
					dimY + o.Wall - 1,
					o.Wall + cutout.Y,
				}, cut)
			case Right:
				cut = p.NewRotation(mgl64.Vec3{0, 0, 90}, cut)
				cut = p.NewTranslation(mgl64.Vec3{
					o.Wall*2 + 1 + board.X,
					o.Wall + cutout.X + board.padding,
					o.Wall + cutout.Y,
				}, cut)
			case Bottom:
				cut = p.NewTranslation(mgl64.Vec3{
					o.Wall + cutout.X + board.padding,
					-1,
					o.Wall + cutout.Y,
				}, cut)
			case Left:
				cut = p.NewRotation(mgl64.Vec3{0, 0, 90}, cut)
				cut = p.NewTranslation(mgl64.Vec3{
					o.Wall + 1,
					-cutout.Width + board.Y + o.Wall - board.padding - cutout.X,
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

func (o *Case) buildMountingHole() p.Primitive {
	size := o.MountingHolesRadius*2 + o.Wall*2
	var mount p.Primitive = p.NewHull(
		p.NewTranslation(
			mgl64.Vec3{o.Wall / 2, 0, 0},
			p.NewCube(mgl64.Vec3{o.Wall, size, o.Wall}).SetCenter(true),
		),
		p.NewTranslation(
			mgl64.Vec3{size / 2, 0, 0},
			p.NewCylinder(o.Wall, size/2),
		),
	)

	mount = p.NewDifference(
		mount,
		p.NewTranslation(
			mgl64.Vec3{o.MountingHolesRadius + o.Wall, 0, 0},
			p.NewCylinder(o.Wall*3, o.MountingHolesRadius),
		),
	)

	return mount
}

func (o *Case) addMountingHoles(box p.Primitive) p.Primitive {
	x, y, _ := o.GetDimensions(true)

	mount1 := p.NewTranslation(
		mgl64.Vec3{x - 0.1, y / 2, o.Wall / 2},
		o.buildMountingHole(),
	)

	mount2 := o.buildMountingHole()
	mount2 = p.NewMirror(mgl64.Vec3{1, 0, 0}, mount2)
	mount2 = p.NewTranslation(
		mgl64.Vec3{0.1, y / 2, o.Wall / 2},
		mount2,
	)

	box = p.NewUnion(
		box,
		mount1,
		mount2,
	)

	return box
}

func (o *Case) BuildBox() p.Primitive {
	standoffs := []p.Primitive{}

	_, fullY, _ := o.GetDimensions(false)

	// Build the base-block.
	var x float64
	for _, board := range o.Boards {
		// Add hole standoffs.
		for _, hole := range board.Holes {
			newStandoff := o.buildStandoff(x, hole)
			// Move the board to the y center.
			newStandoff = p.NewTranslation(mgl64.Vec3{0, (fullY - board.Y) / 2, 0}, newStandoff)
			standoffs = append(standoffs, newStandoff)
		}

		x += board.X
	}

	x, y, height := o.GetDimensions(false)
	xWithWall, yWithWall, heightWithWall := o.GetDimensions(true)

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
		p.NewUnion(standoffs...),
	)

	box = o.addMountingHoles(box)

	o.BoxPrimitive = box
	return o.BoxPrimitive
}
