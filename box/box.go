package box

import (
	"math"

	"github.com/aligator/ControllerCase/board"

	"github.com/go-gl/mathgl/mgl64"
	p "github.com/ljanyst/ghostscad/primitive"
)

type Box struct {
	BoxPrimitive   p.Primitive
	CoverPrimitive p.Primitive

	Boards         []board.Board
	Wall           float64
	StandoffHeight float64
	CoverInsert    float64
	HeightPadding  float64

	MountingHolesRadius float64

	CoverHoles bool
}

func NewBox(
	wall float64,
	standoffHeight float64,
	coverInsert float64,
	mountingHolesRadius float64,
	boards ...board.Board,
) *Box {
	return &Box{
		Boards:              boards,
		Wall:                wall,
		StandoffHeight:      standoffHeight,
		CoverInsert:         coverInsert,
		MountingHolesRadius: mountingHolesRadius,
	}
}

func (o *Box) SetCoverHoles() *Box {
	o.CoverHoles = true
	return o
}

func (o *Box) SetHeightPadding(padding float64) *Box {
	o.HeightPadding = padding
	return o
}

func (o *Box) GetDimensions(withWalls bool) (x, y, height float64) {
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

func (o *Box) GetCoverDimensions(withWalls bool) (x, y, height float64) {
	for _, board := range o.Boards {
		x += board.X
		y = math.Max(y, board.Y)
	}

	if withWalls {
		x += 2 * o.Wall
		y += 2 * o.Wall
	}

	return x, y, o.Wall * 2
}

func (o *Box) buildStandoff(x float64, hole board.Hole) p.Primitive {
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

func (o *Box) BuildCover() p.Primitive {
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

func (o *Box) applyCutouts(box p.Primitive) p.Primitive {
	cuts := []p.Primitive{}

	_, dimY, _ := o.GetDimensions(false)

	x := 0.0
	for _, b := range o.Boards {
		for _, cutout := range b.Cutouts {
			var cut p.Primitive = p.NewCube(mgl64.Vec3{cutout.Width, o.Wall + 2, cutout.Height}).SetCenter(false)

			switch cutout.Side {
			case board.SideTop:
				cut = p.NewTranslation(mgl64.Vec3{
					cutout.X + b.Padding[board.SideLeft] + o.Wall,
					dimY + o.Wall - 1,
					o.Wall + cutout.Y,
				}, cut)
			case board.SideRight:
				cut = p.NewRotation(mgl64.Vec3{0, 0, 90}, cut)
				cut = p.NewTranslation(mgl64.Vec3{
					o.Wall*2 + 1 + b.X,
					o.Wall + cutout.X + b.Padding[board.SideBottom],
					o.Wall + cutout.Y,
				}, cut)
			case board.SideBottom:
				cut = p.NewTranslation(mgl64.Vec3{
					o.Wall + cutout.X + b.Padding[board.SideLeft],
					-1,
					o.Wall + cutout.Y,
				}, cut)
			case board.SideLeft:
				cut = p.NewRotation(mgl64.Vec3{0, 0, 90}, cut)
				cut = p.NewTranslation(mgl64.Vec3{
					o.Wall + 1,
					-cutout.Width + b.Y + o.Wall - b.Padding[board.SideTop] - cutout.X,
					o.Wall + cutout.Y,
				}, cut)
			}

			cut = p.NewTranslation(mgl64.Vec3{x, 0, o.StandoffHeight}, cut)

			cuts = append(cuts, cut)
		}

		x += b.X
	}

	return p.NewDifference(append(
		[]p.Primitive{box},
		cuts...,
	)...)
}

func (o *Box) buildMountingHole() p.Primitive {
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

func (o *Box) addMountingHoles(box p.Primitive) p.Primitive {
	x, y, _ := o.GetDimensions(true)

	// The 0.01 is just to make sure the mount is connected with the box.
	// Otherwise I had issues with PrusaSlicer.

	mount1 := p.NewTranslation(
		mgl64.Vec3{x - 0.01, y / 2, o.Wall / 2},
		o.buildMountingHole(),
	)

	mount2 := o.buildMountingHole()
	mount2 = p.NewMirror(mgl64.Vec3{1, 0, 0}, mount2)
	mount2 = p.NewTranslation(
		mgl64.Vec3{0.01, y / 2, o.Wall / 2},
		mount2,
	)

	box = p.NewUnion(
		box,
		mount1,
		mount2,
	)

	return box
}

func (o *Box) BuildBox() p.Primitive {
	standoffs := []p.Primitive{}

	_, fullY, _ := o.GetDimensions(false)

	// Build the base-block.
	var x float64
	for _, b := range o.Boards {
		// Add hole standoffs.
		for _, hole := range b.Holes {
			newStandoff := o.buildStandoff(x, hole)
			switch b.Position {
			case board.PositionTop:
				newStandoff = p.NewTranslation(mgl64.Vec3{0, (fullY - b.Y), 0}, newStandoff)
			case board.PositionCenter:
				newStandoff = p.NewTranslation(mgl64.Vec3{0, (fullY - b.Y) / 2, 0}, newStandoff)
			}

			standoffs = append(standoffs, newStandoff)
		}

		x += b.X
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
