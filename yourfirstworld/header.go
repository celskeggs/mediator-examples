package main

import (
	_ "github.com/celskeggs/mediator/platform/atoms"
	_ "github.com/celskeggs/mediator/platform/datum"
	"github.com/celskeggs/mediator/platform/framework"
	_ "github.com/celskeggs/mediator/platform/world"
)

//go:generate go run github.com/celskeggs/mediator/autocoder main-a.dm main-b.dm main.go
//go:generate go run github.com/celskeggs/mediator/boilerplate

func main() {
	framework.Launch(Tree, BeforeMap, framework.ResourceDefaults{
		CoreResourcesDir: "../resources",
		IconsDir:         "resources",
		MapPath:          "map.dmm",
	})
}
