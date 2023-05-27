package cases

import (
	"github.com/aligator/ControllerCase/board"
	"github.com/aligator/ControllerCase/box"
	"github.com/ljanyst/ghostscad/sys"
)

// CO2SensorCase contains a NodeMCU and a SCD30 sensor.
func CO2SensorCase() []sys.Shape {
	co2SensorCase := box.NewBox(
		1.4, // wall
		10,  // standoff height
		2,   // cover insert
		1.5, // mounting holes radius
		board.AmicaNodeMCU().WithPadding(0.6, board.SideAll).WithPosition(board.PositionTop),
		board.SensirionSCD30().WithPadding(15, board.SideAll).WithPosition(board.PositionCenter).
			// Add cutout for better air ventilation.
			WithCutout(board.Cutout{
				X:      (23.0 - 15) / 2,
				Y:      -5,
				Width:  15,
				Height: 10,
				Side:   board.SideTop,
			}).
			WithCutout(board.Cutout{
				X:      (23.0 - 15) / 2,
				Y:      -4,
				Width:  15,
				Height: 10,
				Side:   board.SideBottom,
			}).
			WithCutout(board.Cutout{
				X:      (35.0 - 20) / 2,
				Y:      -5,
				Width:  20,
				Height: 10,
				Side:   board.SideRight,
			}),
	).
		SetCoverHoles().
		SetHeightPadding(1)

	return Finish(co2SensorCase)
}

// KY001JoyItCase is just a container for a KY001 board.
func KY001JoyItCase() []sys.Shape {
	converterCase := box.NewBox(
		1.4, // wall
		3,   // standoff height
		2,   // cover insert
		2.5, // mounting holes radius
		board.KY001JoyIt().
			WithPadding(1, board.SideAll).
			WithPadding(20, board.SideTop),
	).
		SetCoverHoles()

	return Finish(converterCase)
}

// HumiditySenorCase contains a NodeMCU, Joyit BMP280 air pressure sensor and an AM2302 humidity sensor.
func HumiditySensorCase() []sys.Shape {
	humiditySensorCase := box.NewBox(
		1.4, // wall
		5,   // standoff height
		2,   // cover insert
		1.5, // mounting holes radius
		board.LolinV3NodeMCU().
			WithPadding(0.6, board.SideAll).
			WithPosition(board.PositionTop).
			WithCutout(board.Cutout{
				X:      (23.0 - 15) / 2,
				Y:      0.0,
				Width:  30,
				Height: 10,
				Side:   board.SideLeft,
			}),
		board.RelaisMakerfactory().
			WithPosition(board.PositionCenter).
			WithPadding(1, board.SideAll).
			WithAdditionalHeight(2).
			WithoutCutouts(),
		board.BMP280JoyIt().
			WithPadding(3, board.SideAll).
			WithoutCutouts(),
		board.AM2302_DHT22().
			WithPadding(3, board.SideAll).
			WithoutCutouts().
			WithCutout(board.Cutout{
				X:      (23.0 - 15) / 2,
				Y:      0.0,
				Width:  30,
				Height: 10,
				Side:   board.SideRight,
			}),
	).
		SetCoverHoles().
		SetHeightPadding(3)

	return Finish(humiditySensorCase)
}
