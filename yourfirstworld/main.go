package main

import (
	"github.com/celskeggs/mediator/platform"
	"github.com/celskeggs/mediator/platform/datum"
	"github.com/celskeggs/mediator/platform/framework"
	"github.com/celskeggs/mediator/platform/icon"
	"github.com/celskeggs/mediator/util"
)

type YourFirstWorld struct {
	platform.BaseTreeDefiner
}

type IMobPlayer interface {
	platform.IMob
}

type MobPlayer struct {
	platform.IMob
}

var _ IMobPlayer = &MobPlayer{}

func (d MobPlayer) RawClone() datum.IDatum {
	d.IMob = d.IMob.RawClone().(platform.IMob)
	return &d
}

///// ***** CustomArea

type ICustomArea interface {
	platform.IArea
	AsCustomArea() *CustomArea
}

type CustomArea struct {
	platform.IArea
	Music string
}

var _ ICustomArea = &CustomArea{}

func (d CustomArea) RawClone() datum.IDatum {
	d.IArea = d.IArea.RawClone().(platform.IArea)
	return &d
}

func (d *CustomArea) AsCustomArea() *CustomArea {
	return d
}

func (d *CustomArea) Entered(atom platform.IAtomMovable, oldloc platform.IAtom) {
	if mob, ismob := atom.(platform.IMob); ismob {
		mob.OutputString(d.AsAtom().Appearance.Desc)
		util.FIXME("use actual sounds, not just strings")
		mob.OutputSound(d.World().Sound(d.Music, true, false, 1, 100))
	}
}

func (d *CustomArea) NextOverride() (this datum.IDatum, next datum.IDatum) {
	return d, d.IArea
}

func ExtractCustomArea(area platform.IArea) ICustomArea {
	datum.AssertConsistent(area)

	var iter datum.IDatum = area
	for {
		cur, next := iter.NextOverride()
		if ca, ok := cur.(ICustomArea); ok {
			return ca
		}
		if iter == nil {
			panic("area does not implement CustomArea")
		}
		iter = next
	}
}

func (y YourFirstWorld) AreaTemplate(parent platform.IAtom) platform.IArea {
	area := y.BaseTreeDefiner.AreaTemplate(parent)
	return &CustomArea{
		IArea: area,
		Music: "",
	}
}

func (YourFirstWorld) ElaborateTree(tree *datum.TypeTree, icons *icon.IconCache) {
	mobPlayer := &MobPlayer{
		IMob: tree.DeriveNew("/mob").(platform.IMob),
	}
	mobPlayer.AsAtom().Appearance.Icon = icons.LoadOrPanic("player.dmi")
	tree.RegisterStruct("/mob/player", mobPlayer)

	mobRat := tree.Derive("/mob", "/mob/rat").(platform.IMob)
	mobRat.AsAtom().Appearance.Icon = icons.LoadOrPanic("rat.dmi")

	turfFloor := tree.Derive("/turf", "/turf/floor").(platform.ITurf)
	turfFloor.AsAtom().Appearance.Icon = icons.LoadOrPanic("floor.dmi")

	turfWall := tree.Derive("/turf", "/turf/wall").(platform.ITurf)
	turfWall.AsAtom().Appearance.Icon = icons.LoadOrPanic("wall.dmi")
	turfWall.AsAtom().Density = true
	turfWall.AsAtom().Opacity = true

	objCheese := tree.Derive("/obj", "/obj/cheese").(platform.IObj)
	objCheese.AsAtom().Appearance.Icon = icons.LoadOrPanic("cheese.dmi")

	objScroll := tree.Derive("/obj", "/obj/scroll").(platform.IObj)
	objScroll.AsAtom().Appearance.Icon = icons.LoadOrPanic("scroll.dmi")

	areaOutside := tree.Derive("/area", "/area/outside").(platform.IArea)
	areaOutside.AsAtom().Appearance.Desc = "Nice and jazzy, here..."
	ExtractCustomArea(areaOutside).AsCustomArea().Music = "jazzy.mid"

	areaCave := tree.Derive("/area", "/area/cave").(platform.IArea)
	areaCave.AsAtom().Appearance.Desc = "Watch out for the giant rat!"
	ExtractCustomArea(areaCave).AsCustomArea().Music = "cavern.mid"
}

func (YourFirstWorld) BeforeMap(world *platform.World) {
	world.Name = "Your First World"
	world.Mob = "/mob/player"
}

func (y YourFirstWorld) Definer() platform.TreeDefiner {
	return y
}

func main() {
	framework.Launch(YourFirstWorld{}, framework.ResourceDefaults{
		CoreResourcesDir: "../resources",
		IconsDir:         "resources",
		MapPath:          "map.dmm",
	})
}
