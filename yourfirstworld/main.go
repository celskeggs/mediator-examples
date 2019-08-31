package main

import (
	"github.com/celskeggs/mediator/platform"
	"github.com/celskeggs/mediator/platform/datum"
	"github.com/celskeggs/mediator/platform/framework"
	"github.com/celskeggs/mediator/platform/icon"
	"github.com/celskeggs/mediator/platform/format"
)

type DefinedWorld struct {
	platform.BaseTreeDefiner
}

///// ***** MobPlayer

type IMobPlayer interface {
	platform.IMob
	AsMobPlayer() *MobPlayer
}

type MobPlayer struct {
	platform.IMob
}

var _ IMobPlayer = &MobPlayer{}

func (d MobPlayer) RawClone() datum.IDatum {
	d.IMob = d.IMob.RawClone().(platform.IMob)
	return &d
}

func (d *MobPlayer) AsMobPlayer() *MobPlayer {
	return d
}

func (this *MobPlayer) Bump(obstacle platform.IAtom) {
	this.OutputString(format.Format("You bump into [].", obstacle))
	this.OutputSound(this.World().Sound("ouch.wav", false, false, 0, 100))
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

func (d *CustomArea) NextOverride() (this datum.IDatum, next datum.IDatum) {
	return d, d.IArea
}

func CastCustomArea(base platform.IArea) ICustomArea {
	datum.AssertConsistent(base)

	var iter datum.IDatum = base
	for {
		cur, next := iter.NextOverride()
		if ca, ok := cur.(ICustomArea); ok {
			return ca
		}
		iter = next
		if iter == nil {
			panic("type does not implement CustomArea")
		}
	}
}

func (d DefinedWorld) AreaTemplate(parent platform.IAtom) platform.IArea {
	base := d.BaseTreeDefiner.AreaTemplate(parent)
	return &CustomArea{
		IArea: base,
		Music: "",
	}
}

func (this *CustomArea) Entered(atom platform.IAtomMovable, oldloc platform.IAtom) {
	if mob, ismob := atom.(platform.IMob); ismob {
		mob.OutputString(this.AsAtom().Appearance.Desc)
		mob.OutputSound(this.World().Sound(this.Music, true, false, 1, 100))
	}
}

func (DefinedWorld) ElaborateTree(tree *datum.TypeTree, icons *icon.IconCache) {
	prototypeMobPlayer := &MobPlayer{
		IMob: tree.DeriveNew("/mob").(platform.IMob),
	}
	prototypeMobPlayer.AsAtom().Appearance.Name = "player"
	prototypeMobPlayer.AsAtom().Appearance.Icon = icons.LoadOrPanic("player.dmi")
	tree.RegisterStruct("/mob/player", prototypeMobPlayer)

	prototypeMobRat := tree.Derive("/mob", "/mob/rat").(platform.IMob)
	prototypeMobRat.AsAtom().Appearance.Name = "rat"
	prototypeMobRat.AsAtom().Appearance.Icon = icons.LoadOrPanic("rat.dmi")

	prototypeTurfFloor := tree.Derive("/turf", "/turf/floor").(platform.ITurf)
	prototypeTurfFloor.AsAtom().Appearance.Name = "floor"
	prototypeTurfFloor.AsAtom().Appearance.Icon = icons.LoadOrPanic("floor.dmi")

	prototypeTurfWall := tree.Derive("/turf", "/turf/wall").(platform.ITurf)
	prototypeTurfWall.AsAtom().Appearance.Name = "wall"
	prototypeTurfWall.AsAtom().Appearance.Icon = icons.LoadOrPanic("wall.dmi")
	prototypeTurfWall.AsAtom().Density = true
	prototypeTurfWall.AsAtom().Opacity = true

	prototypeObjCheese := tree.Derive("/obj", "/obj/cheese").(platform.IObj)
	prototypeObjCheese.AsAtom().Appearance.Name = "cheese"
	prototypeObjCheese.AsAtom().Appearance.Icon = icons.LoadOrPanic("cheese.dmi")

	prototypeObjScroll := tree.Derive("/obj", "/obj/scroll").(platform.IObj)
	prototypeObjScroll.AsAtom().Appearance.Name = "scroll"
	prototypeObjScroll.AsAtom().Appearance.Icon = icons.LoadOrPanic("scroll.dmi")

	prototypeAreaOutside := tree.Derive("/area", "/area/outside").(platform.IArea)
	prototypeAreaOutside.AsAtom().Appearance.Name = "outside"
	prototypeAreaOutside.AsAtom().Appearance.Desc = "Nice and jazzy, here..."
	CastCustomArea(prototypeAreaOutside).AsCustomArea().Music = "jazzy.ogg"

	prototypeAreaCave := tree.Derive("/area", "/area/cave").(platform.IArea)
	prototypeAreaCave.AsAtom().Appearance.Name = "cave"
	prototypeAreaCave.AsAtom().Appearance.Desc = "Watch out for the giant rat!"
	CastCustomArea(prototypeAreaCave).AsCustomArea().Music = "cavern.ogg"
}

func (DefinedWorld) BeforeMap(world *platform.World) {
	world.Name = "Your First World"
	world.Mob = "/mob/player"
}

func (d DefinedWorld) Definer() platform.TreeDefiner {
	return d
}

func main() {
	framework.Launch(DefinedWorld{}, framework.ResourceDefaults{
		CoreResourcesDir: "../resources",
		IconsDir:         "resources",
		MapPath:          "map.dmm",
	})
}
