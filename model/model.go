package model

import (
	"github.com/ljanyst/ghostscad/sys"
)

func CO2SensorCase() []sys.Shape {
	co2SensorCase := NewCase(
		0.8,
		2,
		AmicaNodeMCU().WithTolerance(0.6),
		SensirionSCD30().WithTolerance(0.6),
	)

	return []sys.Shape{
		{
			Name:      "case",
			Primitive: co2SensorCase.BuildBox(),
			Flags:     sys.Default,
		},
		{
			Name:      "cover",
			Primitive: co2SensorCase.BuildCover(),
		},
	}
}
