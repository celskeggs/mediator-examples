package main

import (
	"github.com/celskeggs/mediator/platform"
	"github.com/celskeggs/mediator/platform/datum"
	"github.com/celskeggs/mediator/platform/framework"
	"github.com/celskeggs/mediator/platform/icon"
)

type YourFirstWorld struct{}

func (YourFirstWorld) ElaborateTree(tree *datum.TypeTree, icons *icon.IconCache) {
	mobPlayer := tree.Derive("/mob", "/mob/player").(*platform.Mob)
	mobPlayer.AsAtom().Appearance.Icon = icons.LoadOrPanic("player.dmi")

	mobRat := tree.Derive("/mob", "/mob/rat").(*platform.Mob)
	mobRat.AsAtom().Appearance.Icon = icons.LoadOrPanic("rat.dmi")

	turfFloor := tree.Derive("/turf", "/turf/floor").(*platform.Turf)
	turfFloor.AsAtom().Appearance.Icon = icons.LoadOrPanic("floor.dmi")

	turfWall := tree.Derive("/turf", "/turf/wall").(*platform.Turf)
	turfWall.AsAtom().Appearance.Icon = icons.LoadOrPanic("wall.dmi")
	turfWall.AsAtom().Density = true
	turfWall.AsAtom().Opacity = true

	objCheese := tree.Derive("/obj", "/obj/cheese").(*platform.Obj)
	objCheese.AsAtom().Appearance.Icon = icons.LoadOrPanic("cheese.dmi")

	objScroll := tree.Derive("/obj", "/obj/scroll").(*platform.Obj)
	objScroll.AsAtom().Appearance.Icon = icons.LoadOrPanic("scroll.dmi")

	tree.Derive("/area", "/area/outside")
	tree.Derive("/area", "/area/cave")
}

func (YourFirstWorld) BeforeMap(world *platform.World) {
	world.Name = "Your First World"
	world.Mob = "/mob/player"
}

func main() {
	framework.Launch(YourFirstWorld{}, framework.ResourceDefaults{
		CoreResourcesDir: "../resources",
		IconsDir:         "icons",
		MapPath:          "map.dmm",
	})
}
