package main

import (
	"github.com/celskeggs/mediator/platform/atoms"
	"github.com/celskeggs/mediator/platform/datum"
	"github.com/celskeggs/mediator/platform/types"
	"github.com/celskeggs/mediator/platform/world"
	"github.com/celskeggs/mediator/util"
	"github.com/celskeggs/mediator/webclient"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestTurfLocation(t *testing.T) {
	gameworld := BuildWorld()
	turf := gameworld.FindOneType("/turf")
	assert.NotNil(t, turf)
	if turf != nil {
		types.AssertType(turf, "/turf")

		area := turf.Var("loc")
		assert.NotNil(t, area)
		if area != nil {
			types.AssertType(area, "/area")

			turfs := atoms.TurfsInArea(area)
			assert.True(t, len(turfs) > 0)
			assert.Contains(t, turfs, turf)
		}
	}
}

func TestPlayerLocation(t *testing.T) {
	gameworld := BuildWorld()
	gameworld.ServerAPI().AddPlayer()
	player := gameworld.FindOneType("/mob/player")
	assert.NotNil(t, player)
	if player != nil {
		types.AssertType(player, "/mob/player")

		turf := player.Var("loc")
		assert.NotNil(t, turf)
		if turf != nil {
			types.AssertType(turf, "/turf")

			contents := datum.Elements(turf.Var("contents"))
			assert.True(t, len(contents) > 0)
			assert.Contains(t, contents, player)
		}
	}
	x, y, z := world.XYZ(player)
	assert.Equal(t, uint(1), x)
	assert.Equal(t, uint(1), y)
	assert.Equal(t, uint(1), z)
}

func TestPlayerContinuedExistence(t *testing.T) {
	gameworld := BuildWorld()
	gameworld.ServerAPI().AddPlayer().Remove()
	player := gameworld.FindOneType("/mob/player")
	assert.NotNil(t, player)
	if player == nil {
		return
	}
	turf := player.Var("loc")
	assert.NotNil(t, turf)
	if turf == nil {
		return
	}
	types.AssertType(turf, "/turf")
	contents := datum.Elements(turf.Var("contents"))
	assert.True(t, len(contents) > 0)
	assert.Contains(t, contents, player)

	runtime.GC()
	runtime.GC()
	playerAgain := gameworld.FindOneType("/mob/player")
	assert.NotNil(t, playerAgain)
	assert.Equal(t, player, playerAgain)
	contents2 := datum.Elements(turf.Var("contents"))
	assert.Equal(t, len(contents), len(contents2))
	assert.Contains(t, contents2, player)
}

func TestWorldRender(t *testing.T) {
	gameworld := BuildWorld()
	api := gameworld.ServerAPI()
	player := api.AddPlayer()
	view := player.Render()
	assert.True(t, len(view.Sprites) > 0)
	runtime.GC()
	view2 := player.Render()
	assert.True(t, len(view2.Sprites) == len(view.Sprites))
}

func TestSingletons(t *testing.T) {
	realm := types.NewRealm(Tree)
	areaOne := realm.New("/area", nil)
	areaTwo := realm.New("/area", nil)
	assert.True(t, areaOne == areaTwo)

	turfOne := realm.New("/turf", nil)
	turfTwo := realm.New("/turf", nil)
	assert.True(t, turfOne != turfTwo)
}

func TestSingletonsInWorld(t *testing.T) {
	gameworld := BuildWorld()
	areaOne := gameworld.Realm().New("/area", nil)
	areaTwo := gameworld.Realm().New("/area", nil)
	assert.True(t, areaOne == areaTwo)

	turfOne := gameworld.Realm().New("/turf", nil)
	turfTwo := gameworld.Realm().New("/turf", nil)
	assert.True(t, turfOne != turfTwo)
}

func TestSingletonAreas(t *testing.T) {
	gameworld := BuildWorld()
	areas := gameworld.FindAllType("/area")
	pathCount := map[types.TypePath]int{}
	for _, area := range areas {
		pathCount[area.(*types.Datum).Type()] += 1
	}
	assert.Equal(t, 2, len(pathCount))
	assert.Equal(t, 1, pathCount["/area/outside"])
	assert.Equal(t, 1, pathCount["/area/cave"])
	for _, count := range pathCount {
		assert.Equal(t, 1, count)
	}
}

func TestWalkBetweenAreas(t *testing.T) {
	// 11, 4

	gameworld := BuildWorld()
	playerAPI := gameworld.ServerAPI().AddPlayer()

	player := gameworld.FindOneType("/mob/player")
	assert.NotNil(t, player)
	player.SetVar("loc", gameworld.LocateXYZ(11, 4, 1))
	assert.Equal(t, types.TypePath("/area/outside"), atoms.ContainingArea(player).(*types.Datum).Type())

	lines, sounds := playerAPI.PullRequests()
	assert.Equal(t, 1, len(lines))
	assert.Contains(t, lines, "Nice and jazzy, here...")
	assert.Equal(t, 1, len(sounds))
	if len(sounds) >= 1 {
		assert.Equal(t, "jazzy.ogg", sounds[0].File)
	}

	playerAPI.Command(webclient.Command{Verb: ".west"})
	lines, sounds = playerAPI.PullRequests()
	assert.Equal(t, 1, len(lines))
	assert.Contains(t, lines, "Watch out for the giant rat!")
	assert.Equal(t, 1, len(sounds))
	if len(sounds) >= 1 {
		assert.Equal(t, "cavern.ogg", sounds[0].File)
	}

	assert.Equal(t, types.TypePath("/area/cave"), atoms.ContainingArea(player).(*types.Datum).Type())
}

func TestStepOffMap(t *testing.T) {
	gameworld := BuildWorld()
	playerAPI := gameworld.ServerAPI().AddPlayer()

	player := gameworld.FindOneType("/mob/player")
	assert.NotNil(t, player)
	assert.Equal(t, types.TypePath("/area/outside"), atoms.ContainingArea(player).(*types.Datum).Type())
	cx, cy := world.XY(player)
	assert.Equal(t, uint(1), cx)
	assert.Equal(t, uint(1), cy)

	lines, sounds := playerAPI.PullRequests()
	assert.Equal(t, 1, len(lines))
	assert.Contains(t, lines, "Nice and jazzy, here...")
	assert.Equal(t, 1, len(sounds))
	if len(sounds) >= 1 {
		assert.Equal(t, "jazzy.ogg", sounds[0].File)
	}

	playerAPI.Command(webclient.Command{Verb: ".south"})
	lines, sounds = playerAPI.PullRequests()
	assert.Equal(t, 0, len(lines))
	assert.Equal(t, 0, len(sounds))
	nx, ny := world.XY(player)
	assert.Equal(t, uint(1), nx)
	assert.Equal(t, uint(1), ny)
}

func TestWalkIntoWalls(t *testing.T) {
	gameworld := BuildWorld()
	playerAPI := gameworld.ServerAPI().AddPlayer()

	player := gameworld.FindOneType("/mob/player")
	assert.NotNil(t, player)
	_, _ = playerAPI.PullRequests()

	player.SetVar("loc", gameworld.LocateXYZ(6, 4, 1))
	assert.Equal(t, types.TypePath("/area/outside"), atoms.ContainingArea(player).(*types.Datum).Type())

	playerAPI.Command(webclient.Command{Verb: ".north"})
	lines, sounds := playerAPI.PullRequests()
	assert.Equal(t, 1, len(lines))
	if len(lines) >= 1 {
		assert.Contains(t, lines, "You bump into the wall.")
	}
	assert.Equal(t, 1, len(sounds))
	if len(sounds) >= 1 {
		assert.Equal(t, "ouch.wav", sounds[0].File)
	}
	nx, ny, nz := world.XYZ(player)
	assert.Equal(t, uint(6), nx)
	assert.Equal(t, uint(4), ny)
	assert.Equal(t, uint(1), nz)
}

func TestFillUpEntireMap(t *testing.T) {
	gameworld := BuildWorld()
	found := false
	serverAPI := gameworld.ServerAPI()
	for i := 0; i < 1000; i++ {
		player := serverAPI.AddPlayer()
		player.Render()
		player.PullRequests()
		player.Remove()
		playerNowhere := gameworld.FindOne(func(atom *types.Datum) bool {
			return atom.Type() == "/mob/player" && atom.Var("loc") == nil
		})
		if playerNowhere != nil {
			// should be able to add a whole bunch of players
			assert.True(t, i >= 50)
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestLookVerb(t *testing.T) {
	gameworld := BuildWorld()
	playerAPI := gameworld.ServerAPI().AddPlayer()

	player := gameworld.FindOneType("/mob/player")
	assert.NotNil(t, player)
	_, _ = playerAPI.PullRequests()

	playerAPI.Command(webclient.Command{Verb: "look"})
	lines, _ := playerAPI.PullRequests()
	assert.Equal(t, 2, len(lines))
	assert.Contains(t, lines, "You see...")
	assert.Contains(t, lines, "The scroll.  It looks to be rather old.")
}

func TestGetVerb(t *testing.T) {
	gameworld := BuildWorld()
	playerAPI := gameworld.ServerAPI().AddPlayer()

	player := gameworld.FindOneType("/mob/player")
	assert.NotNil(t, player)
	scroll := gameworld.FindOneType("/obj/scroll")
	assert.NotNil(t, scroll)
	_, _ = playerAPI.PullRequests()

	for x := types.Unuint(player.Var("x")); x < types.Unuint(scroll.Var("x")); x++ {
		playerAPI.Command(webclient.Command{Verb: ".east"})
	}
	for y := types.Unuint(player.Var("y")); y < types.Unuint(scroll.Var("y")); y++ {
		playerAPI.Command(webclient.Command{Verb: ".north"})
	}
	playerAPI.Command(webclient.Command{Verb: "get"})
	lines, _ := playerAPI.PullRequests()
	assert.Equal(t, 1, len(lines))
	assert.Contains(t, lines, "You get the scroll.")

	view := playerAPI.Render()
	iconCount := map[string]int{}
	for _, sprite := range view.Sprites {
		iconCount[sprite.Icon] += 1
	}
	// should only be wall, floor, and player, not scroll
	assert.Equal(t, 3, len(iconCount))
	assert.Equal(t, 0, iconCount["scroll.dmi"])
	assert.Contains(t, iconCount, "wall.dmi")
	assert.Contains(t, iconCount, "floor.dmi")
	assert.Equal(t, 1, iconCount["player.dmi"])

	assert.Equal(t, player, scroll.Var("loc"))

	util.FIXME("try dropping")
	util.FIXME("try picking up with 'get scroll' instead of just 'get'")
}
