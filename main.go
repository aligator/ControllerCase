package main

import (
	"github.com/aligator/ControllerCase/cases"
	"github.com/ljanyst/ghostscad/sys"
)

func main() {
	sys.Initialize()

	sys.SetFs(0.2)
	sys.SetFa(12)
	sys.RenderMultiple(cases.HumiditySensorCase())
}
