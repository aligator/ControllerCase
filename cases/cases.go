package cases

import (
	"github.com/aligator/ControllerCase/box"
	"github.com/go-gl/mathgl/mgl64"
	p "github.com/ljanyst/ghostscad/primitive"
	"github.com/ljanyst/ghostscad/sys"
)

func Finish(controllerCase *box.Box) []sys.Shape {
	x, y, _ := controllerCase.GetDimensions(true)
	_, _, coverHeight := controllerCase.GetCoverDimensions(true)

	return []sys.Shape{
		{
			Name:      "all",
			Primitive: controllerCase.BuildBox(),
			Flags:     sys.Default,
		},

		{
			Name: "all",
			Primitive: p.NewUnion(
				controllerCase.BuildBox(),
				p.NewRotationByAxis(180, mgl64.Vec3{1}, p.NewTranslation(mgl64.Vec3{x + 20, -y, -coverHeight}, controllerCase.BuildCover())),
			),
			Flags: sys.Default,
		},
		{
			Name:      "case",
			Primitive: controllerCase.BuildBox(),
		},
		{
			Name:      "cover",
			Primitive: controllerCase.BuildCover(),
		},
	}
}
