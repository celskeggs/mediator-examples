// Code generated by mediator boilerplate; DO NOT EDIT.
package main

import (
	"github.com/celskeggs/mediator/common"
	"github.com/celskeggs/mediator/platform/atoms"
	"github.com/celskeggs/mediator/platform/datum"
	"github.com/celskeggs/mediator/platform/types"
)

type MobPlayerImpl struct {
	MobPlayerData
	atoms.MobData
	atoms.AtomMovableData
	atoms.AtomData
	datum.DatumData
}

func NewMobPlayer(realm *types.Realm, params ...types.Value) *types.Datum {
	i := &MobPlayerImpl{}
	d := realm.NewDatum(i)
	datum.NewDatumData(d, &i.DatumData, params...)
	atoms.NewAtomData(d, &i.AtomData, params...)
	atoms.NewAtomMovableData(d, &i.AtomMovableData, params...)
	atoms.NewMobData(d, &i.MobData, params...)
	NewMobPlayerData(d, &i.MobPlayerData, params...)
	return d
}

func (t *MobPlayerImpl) Type() types.TypePath {
	return "/mob/player"
}

func (t *MobPlayerImpl) Var(src *types.Datum, name string) (types.Value, bool) {
	switch name {
	case "type":
		return types.TypePath("/mob/player"), true
	case "parent_type":
		return types.TypePath("/mob"), true
	case "appearance":
		return t.AtomData.VarAppearance, true
	case "density":
		return types.Int(t.AtomData.VarDensity), true
	case "dir":
		return t.AtomData.VarDir, true
	case "opacity":
		return types.Int(t.AtomData.VarOpacity), true
	case "verbs":
		return datum.NewListFromSlice(t.AtomData.VarVerbs), true
	case "client":
		return t.MobData.GetClient(src), true
	case "contents":
		return t.AtomData.GetContents(src), true
	case "desc":
		return t.AtomData.GetDesc(src), true
	case "icon":
		return t.AtomData.GetIcon(src), true
	case "icon_state":
		return t.AtomData.GetIconState(src), true
	case "key":
		return t.MobData.GetKey(src), true
	case "layer":
		return t.AtomData.GetLayer(src), true
	case "loc":
		return t.AtomData.GetLoc(src), true
	case "name":
		return t.AtomData.GetName(src), true
	case "suffix":
		return t.AtomData.GetSuffix(src), true
	case "x":
		return t.AtomData.GetX(src), true
	case "y":
		return t.AtomData.GetY(src), true
	case "z":
		return t.AtomData.GetZ(src), true
	default:
		return nil, false
	}
}

func (t *MobPlayerImpl) SetVar(src *types.Datum, name string, value types.Value) types.SetResult {
	switch name {
	case "type":
		return types.SetResultReadOnly
	case "parent_type":
		return types.SetResultReadOnly
	case "appearance":
		t.AtomData.VarAppearance = value.(atoms.Appearance)
		return types.SetResultOk
	case "density":
		t.AtomData.VarDensity = types.Unint(value)
		return types.SetResultOk
	case "dir":
		t.AtomData.VarDir = value.(common.Direction)
		return types.SetResultOk
	case "opacity":
		t.AtomData.VarOpacity = types.Unint(value)
		return types.SetResultOk
	case "verbs":
		t.AtomData.VarVerbs = datum.ElementsAsType([]atoms.Verb{}, value).([]atoms.Verb)
		return types.SetResultOk
	case "client":
		return types.SetResultReadOnly
	case "contents":
		return types.SetResultReadOnly
	case "desc":
		t.AtomData.SetDesc(src, value)
		return types.SetResultOk
	case "icon":
		t.AtomData.SetIcon(src, value)
		return types.SetResultOk
	case "icon_state":
		t.AtomData.SetIconState(src, value)
		return types.SetResultOk
	case "key":
		return types.SetResultReadOnly
	case "layer":
		t.AtomData.SetLayer(src, value)
		return types.SetResultOk
	case "loc":
		t.AtomData.SetLoc(src, value)
		return types.SetResultOk
	case "name":
		t.AtomData.SetName(src, value)
		return types.SetResultOk
	case "suffix":
		t.AtomData.SetSuffix(src, value)
		return types.SetResultOk
	case "x":
		return types.SetResultReadOnly
	case "y":
		return types.SetResultReadOnly
	case "z":
		return types.SetResultReadOnly
	default:
		return types.SetResultNonexistent
	}
}

func (t *MobPlayerImpl) Proc(src *types.Datum, usr *types.Datum, name string, params ...types.Value) (types.Value, bool) {
	switch name {
	case "<<":
		return t.MobData.OperatorWrite(src, usr, types.Param(params, 0)), true
	case "Bump":
		return t.MobPlayerData.ProcBump(src, usr, params), true
	case "Enter":
		return t.AtomData.ProcEnter(src, usr, types.Param(params, 0), types.Param(params, 1)), true
	case "Entered":
		return t.AtomData.ProcEntered(src, usr, types.Param(params, 0), types.Param(params, 1)), true
	case "Exit":
		return t.AtomData.ProcExit(src, usr, types.Param(params, 0), types.Param(params, 1)), true
	case "Exited":
		return t.AtomData.ProcExited(src, usr, types.Param(params, 0), types.Param(params, 1)), true
	case "Login":
		return t.MobData.ProcLogin(src, usr), true
	case "Move":
		return t.AtomData.ProcMove(src, usr, types.Param(params, 0), types.Param(params, 1)), true
	case "New":
		return t.DatumData.ProcNew(src, usr), true
	case "Stat":
		return t.MobPlayerData.ProcStat(src, usr, params), true
	case "look":
		return t.MobPlayerData.Proclook(src, usr, params), true
	default:
		return nil, false
	}
}

func (t *MobPlayerImpl) SuperProc(src *types.Datum, usr *types.Datum, chunk string, name string, params ...types.Value) (types.Value, bool) {
	switch chunk {
	case "github.com/celskeggs/mediator-examples/yourfirstworld.MobPlayerData":
		switch name {
		case "Bump":
			return t.AtomData.ProcBump(src, usr, types.Param(params, 0)), true
		case "Stat":
			return t.AtomData.ProcStat(src, usr), true
		}
	}
	return nil, false
}

func (t *MobPlayerImpl) ProcSettings(name string) (types.ProcSettings, bool) {
	switch name {
	case "<<":
		return types.ProcSettings{}, true
	case "Bump":
		return types.ProcSettings{}, true
	case "Enter":
		return types.ProcSettings{}, true
	case "Entered":
		return types.ProcSettings{}, true
	case "Exit":
		return types.ProcSettings{}, true
	case "Exited":
		return types.ProcSettings{}, true
	case "Login":
		return types.ProcSettings{}, true
	case "Move":
		return types.ProcSettings{}, true
	case "New":
		return types.ProcSettings{}, true
	case "Stat":
		return types.ProcSettings{}, true
	case "look":
		return types.ProcSettings{}, true
	default:
		return types.ProcSettings{}, false
	}
}

func (t *MobPlayerImpl) Chunk(ref string) interface{} {
	switch ref {
	case "github.com/celskeggs/mediator-examples/yourfirstworld.MobPlayerData":
		return &t.MobPlayerData
	case "github.com/celskeggs/mediator/platform/atoms.MobData":
		return &t.MobData
	case "github.com/celskeggs/mediator/platform/atoms.AtomMovableData":
		return &t.AtomMovableData
	case "github.com/celskeggs/mediator/platform/atoms.AtomData":
		return &t.AtomData
	case "github.com/celskeggs/mediator/platform/datum.DatumData":
		return &t.DatumData
	default:
		return nil
	}
}
