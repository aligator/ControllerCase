package main

import (
	"co2sensor/model"

	"github.com/ljanyst/ghostscad/sys"
)

func main() {
	sys.Initialize()
	sys.RenderMultiple(model.CO2SensorCase())
}
