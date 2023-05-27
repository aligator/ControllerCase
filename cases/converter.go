package cases

import (
	"github.com/aligator/ControllerCase/board"
	"github.com/aligator/ControllerCase/box"
	"github.com/ljanyst/ghostscad/sys"
)

func LM317TConverterCase() []sys.Shape {
	converterCase := box.NewBox(
		1.4, // wall
		3,   // standoff height
		2,   // cover insert
		2.5, // mounting holes radius
		board.LM317TConverter().WithPadding(3, board.SideAll),
	).
		SetCoverHoles()

	return Finish(converterCase)
}
