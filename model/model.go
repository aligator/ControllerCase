package model

import (
	"github.com/go-gl/mathgl/mgl64"
	p "github.com/ljanyst/ghostscad/primitive"
	"github.com/ljanyst/ghostscad/sys"
)

func Finish(controllerCase *Case) []sys.Shape {
	_, _, height := controllerCase.GetDimensions(false)

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
				p.NewTranslation(mgl64.Vec3{0, 0, height + 20}, controllerCase.BuildCover()),
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

func CO2SensorCase() []sys.Shape {
	co2SensorCase := NewCase(
		1.4, // wall
		10,  // standoff height
		2,   // cover insert
		1.5, // mounting holes radius
		AmicaNodeMCU().WithPadding(0.6).WithPosition(PositionTop),
		SensirionSCD30().WithPadding(15).WithPosition(PositionCenter).
			// Add cutout for better air ventilation.
			WithCutout(Cutout{
				X:      (23.0 - 15) / 2,
				Y:      -5,
				Width:  15,
				Height: 10,
				Side:   SideTop,
			}).
			WithCutout(Cutout{
				X:      (23.0 - 15) / 2,
				Y:      -4,
				Width:  15,
				Height: 10,
				Side:   SideBottom,
			}).
			WithCutout(Cutout{
				X:      (35.0 - 20) / 2,
				Y:      -5,
				Width:  20,
				Height: 10,
				Side:   SideRight,
			}),
	).
		SetCoverHoles().
		SetHeightPadding(1)

	return Finish(co2SensorCase)
}

func LM317TConverterCase() []sys.Shape {
	converterCase := NewCase(
		1.4, // wall
		3,   // standoff height
		2,   // cover insert
		1.5, // mounting holes radius#
		LM317TConverter().WithPadding(1),
	).
		SetCoverHoles()

	return Finish(converterCase)
}
