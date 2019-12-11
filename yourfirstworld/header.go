package main

import (
	_ "github.com/celskeggs/mediator/platform/atoms"
	_ "github.com/celskeggs/mediator/platform/datum"
	_ "github.com/celskeggs/mediator/platform/world"
)

//go:generate go run github.com/celskeggs/mediator/autocoder "Your First World.dme"
//go:generate go run github.com/celskeggs/mediator/boilerplate
