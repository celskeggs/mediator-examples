package main

import (
	"github.com/celskeggs/mediator/platform"
	"github.com/celskeggs/mediator/platform/datum"
	"github.com/celskeggs/mediator/platform/worldmap"
	"github.com/celskeggs/mediator/websession"
	"path"
)

func BuildTree() *datum.TypeTree {
	tree := platform.NewAtomicTree()

	mobPlayer := tree.Derive("/mob", "/mob/player").(*platform.Mob)
	mobPlayer.Appearance.Icon = "player.dmi"

	mobRat := tree.Derive("/mob", "/mob/rat").(*platform.Mob)
	mobRat.Appearance.Icon = "rat.dmi"

	turfFloor := tree.Derive("/turf", "/turf/floor").(*platform.Turf)
	turfFloor.Appearance.Icon = "floor.dmi"

	turfWall := tree.Derive("/turf", "/turf/wall").(*platform.Turf)
	turfWall.Appearance.Icon = "wall.dmi"
	turfWall.Density = true
	turfWall.Opacity = true

	objCheese := tree.Derive("/obj", "/obj/cheese").(*platform.Obj)
	objCheese.Appearance.Icon = "cheese.dmi"

	objScroll := tree.Derive("/obj", "/obj/scroll").(*platform.Obj)
	objScroll.Appearance.Icon = "scroll.dmi"

	tree.Derive("/area", "/area/outside")
	tree.Derive("/area", "/area/cave")

	return tree
}

func BuildWorld() *platform.World {
	_, resources := websession.ParseFlags()
	world := platform.NewWorld(BuildTree(), "Your First World", "/mob/player", "/client")
	err := worldmap.LoadMapFromFile(world, path.Join(resources, "../map.dmm"))
	if err != nil {
		panic("cannot load world: " + err.Error())
	}
	return world
}

func main() {
	websession.SetDefaultFlags("../resources", "icons")

	world := BuildWorld()

	websession.LaunchServerFromFlags(world.ServerAPI())
}
