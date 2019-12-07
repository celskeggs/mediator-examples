package world

import (
	"github.com/celskeggs/mediator/platform/atoms"
	"github.com/celskeggs/mediator/platform/format"
	"github.com/celskeggs/mediator/platform/procs"
	"github.com/celskeggs/mediator/platform/types"
	"github.com/celskeggs/mediator/platform/world"
	"github.com/celskeggs/mediator/util"
	"github.com/celskeggs/mediator/webclient/sprite"
)

//mediator:declare MobPlayerData /mob/player /mob
type MobPlayerData struct{}

func NewMobPlayerData(src *types.Datum, _ ...types.Value) MobPlayerData {
	src.SetVar("name", types.String("player"))
	src.SetVar("icon", atoms.WorldOf(src).Icon("player.dmi"))
	return MobPlayerData{}
}

func (d *MobPlayerData) ProcBump(src *types.Datum, obstacle types.Value) types.Value {
	src.Invoke("<<", types.String(format.Format("You bump into [].", obstacle)))
	src.Invoke("<<", procs.NewSound("ouch.wav"))
	return nil
}

//mediator:declare MobRatData /mob/rat /mob
type MobRatData struct{}

func NewMobRatData(src *types.Datum, _ ...types.Value) MobRatData {
	src.SetVar("name", types.String("rat"))
	src.SetVar("icon", atoms.WorldOf(src).Icon("rat.dmi"))
	return MobRatData{}
}

//mediator:declare TurfFloorData /turf/floor /turf
type TurfFloorData struct{}

func NewTurfFloorData(src *types.Datum, _ ...types.Value) TurfFloorData {
	src.SetVar("name", types.String("floor"))
	src.SetVar("icon", atoms.WorldOf(src).Icon("floor.dmi"))
	return TurfFloorData{}
}

//mediator:declare TurfWallData /turf/wall /turf
type TurfWallData struct{}

func NewTurfWallData(src *types.Datum, _ ...types.Value) TurfWallData {
	src.SetVar("name", types.String("wall"))
	src.SetVar("icon", atoms.WorldOf(src).Icon("wall.dmi"))
	src.SetVar("density", types.Bool(true))
	src.SetVar("opacity", types.Bool(true))
	return TurfWallData{}
}

//mediator:declare ObjCheeseData /obj/cheese /obj
type ObjCheeseData struct{}

func NewObjCheeseData(src *types.Datum, _ ...types.Value) ObjCheeseData {
	src.SetVar("name", types.String("cheese"))
	src.SetVar("icon", atoms.WorldOf(src).Icon("cheese.dmi"))
	return ObjCheeseData{}
}

//mediator:declare ObjScrollData /obj/scroll /obj
type ObjScrollData struct{}

func NewObjScrollData(src *types.Datum, _ ...types.Value) ObjScrollData {
	src.SetVar("name", types.String("scroll"))
	src.SetVar("icon", atoms.WorldOf(src).Icon("scroll.dmi"))
	return ObjScrollData{}
}

//mediator:extend ExtAreaData /area
type ExtAreaData struct {
	VarMusic sprite.Sound
}

func NewExtAreaData(src *types.Datum, _ ...types.Value) ExtAreaData {
	return ExtAreaData{}
}

func (d *ExtAreaData) ProcEntered(src *types.Datum, atom types.Value) types.Value {
	if types.IsType(atom, "/mob") {
		atom.Invoke("<<", src.Var("desc"))
		atom.Invoke("<<", procs.NewSoundFrom(d.VarMusic, true, false, 1, 100))
	}
	return nil
}

//mediator:declare AreaOutsideData /area/outside /area
type AreaOutsideData struct{}

func NewAreaOutsideData(src *types.Datum, _ ...types.Value) AreaOutsideData {
	src.SetVar("name", types.String("outside"))
	src.SetVar("desc", types.String("Nice and jazzy, here..."))
	src.SetVar("music", procs.NewSound("jazzy.ogg"))
	return AreaOutsideData{}
}

//mediator:declare AreaCaveData /area/cave /area
type AreaCaveData struct{}

func NewAreaCaveData(src *types.Datum, _ ...types.Value) AreaCaveData {
	src.SetVar("name", types.String("cave"))
	src.SetVar("desc", types.String("Watch out for the giant rat!"))
	src.SetVar("music", procs.NewSound("cavern.ogg"))
	return AreaCaveData{}
}

func BeforeMap(world *world.World) {
	util.FIXME("eliminate this vestigal interface")
	world.Name = "Your First World"
	world.Mob = "/mob/player"
}
