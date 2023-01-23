package model

import (
	"github.com/go-gl/mathgl/mgl64"
	p "github.com/ljanyst/ghostscad/primitive"
	"github.com/ljanyst/ghostscad/sys"
)

func CO2SensorCase() []sys.Shape {
	co2SensorCase := NewCase(
		0.8,
		2,
		AmicaNodeMCU().WithTolerance(0.6),
		SensirionSCD30().WithTolerance(0.6),
	)

	_, _, height := co2SensorCase.GetDimensions()

	return []sys.Shape{
		{
			Name: "all",
			Primitive: p.NewUnion(
				co2SensorCase.BuildBox(),
				p.NewTranslation(mgl64.Vec3{0, 0, height + 20}, co2SensorCase.BuildCover()),
			),
			Flags: sys.Default,
		},
		{
			Name:      "case",
			Primitive: co2SensorCase.BuildBox(),
		},
		{
			Name:      "cover",
			Primitive: co2SensorCase.BuildCover(),
		},
	}
}
