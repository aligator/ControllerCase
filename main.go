package main

import (
	"co2sensor/model"

	"github.com/ljanyst/ghostscad/sys"
)

func main() {
	sys.Initialize()

	sys.SetFs(0.2)
	sys.SetFa(12)
	sys.RenderMultiple(model.CO2SensorCase())
}
