package main

import "github.com/celskeggs/mediator/platform/framework"

//go:generate go run github.com/celskeggs/mediator/autocoder main-a.dm main.go

func main() {
	framework.Launch(DefinedWorld{}, framework.ResourceDefaults{
		CoreResourcesDir: "../resources",
		IconsDir:         "resources",
		MapPath:          "map.dmm",
	})
}
