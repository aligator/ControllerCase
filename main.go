package main

import (
	p "github.com/ljanyst/ghostscad/primitive"
	"github.com/ljanyst/ghostscad/sys"
)

func main() {
	sys.Initialize()
	sys.RenderOne(p.NewSphere(20))
}
