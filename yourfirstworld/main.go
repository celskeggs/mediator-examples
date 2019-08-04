package main

import (
	"github.com/celskeggs/mediator/platform"
	"github.com/celskeggs/mediator/platform/datum"
	"github.com/celskeggs/mediator/platform/icon"
	"github.com/celskeggs/mediator/platform/worldmap"
	"github.com/celskeggs/mediator/websession"
	"path"
)

func BuildTree(resourceDir string) *datum.TypeTree {
	tree := platform.NewAtomicTree()
	icons := icon.NewIconCache(resourceDir)

	mobPlayer := tree.Derive("/mob", "/mob/player").(*platform.Mob)
	mobPlayer.Appearance.Icon = icons.LoadOrPanic("player.dmi")

	mobRat := tree.Derive("/mob", "/mob/rat").(*platform.Mob)
	mobRat.Appearance.Icon = icons.LoadOrPanic("rat.dmi")

	turfFloor := tree.Derive("/turf", "/turf/floor").(*platform.Turf)
	turfFloor.Appearance.Icon = icons.LoadOrPanic("floor.dmi")

	turfWall := tree.Derive("/turf", "/turf/wall").(*platform.Turf)
	turfWall.Appearance.Icon = icons.LoadOrPanic("wall.dmi")
	turfWall.Density = true
	turfWall.Opacity = true

	objCheese := tree.Derive("/obj", "/obj/cheese").(*platform.Obj)
	objCheese.Appearance.Icon = icons.LoadOrPanic("cheese.dmi")

	objScroll := tree.Derive("/obj", "/obj/scroll").(*platform.Obj)
	objScroll.Appearance.Icon = icons.LoadOrPanic("scroll.dmi")

	tree.Derive("/area", "/area/outside")
	tree.Derive("/area", "/area/cave")

	return tree
}

func BuildWorld() *platform.World {
	_, resources := websession.ParseFlags()
	world := platform.NewWorld(BuildTree(resources))
	world.Name = "Your First World"
	world.Mob = "/mob/player"
	err := worldmap.LoadMapFromFile(world, path.Join(resources, "../map.dmm"))
	if err != nil {
		panic("cannot load world: " + err.Error())
	}
	world.UpdateDefaultViewDistance()
	return world
}

func main() {
	websession.SetDefaultFlags("../resources", "icons")

	world := BuildWorld()

	websession.LaunchServerFromFlags(world.ServerAPI())
}
