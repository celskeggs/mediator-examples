package main

import (
	"github.com/celskeggs/mediator/platform"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
	"github.com/celskeggs/mediator/platform/framework"
	"github.com/celskeggs/mediator/platform/datum"
	"github.com/celskeggs/mediator/webclient"
)

func BuildWorld() *platform.World {
	return framework.BuildWorld(YourFirstWorld{}, framework.ResourceDefaults{
		CoreResourcesDir: "../resources",
		IconsDir:         "resources",
		MapPath:          "map.dmm",
	}, false)
}

func TestTurfLocation(t *testing.T) {
	world := BuildWorld()
	turf := world.FindOne(func(atom platform.IAtom) bool {
		_, isturf := atom.(platform.ITurf)
		return isturf
	})
	assert.NotNil(t, turf)
	if turf != nil {
		area := turf.(platform.ITurf).Location()
		assert.NotNil(t, area)
		if area != nil {
			turfs := area.(platform.IArea).Turfs()
			assert.True(t, len(turfs) > 0)
			assert.Contains(t, turfs, turf)
		}
	}
}

func TestPlayerLocation(t *testing.T) {
	world := BuildWorld()
	world.ServerAPI().AddPlayer()
	player := world.FindOne(func(atom platform.IAtom) bool {
		return atom.AsDatum().Type == "/mob/player"
	})
	assert.NotNil(t, player)
	if player != nil {
		turf := player.(platform.IMob).Location()
		assert.NotNil(t, turf)
		if turf != nil {
			contents := turf.(platform.ITurf).Contents()
			assert.True(t, len(contents) > 0)
			assert.Contains(t, contents, player)
		}
	}
	x, y, z := player.XYZ()
	assert.Equal(t, uint(1), x)
	assert.Equal(t, uint(1), y)
	assert.Equal(t, uint(1), z)
}

func TestPlayerContinuedExistence(t *testing.T) {
	world := BuildWorld()
	world.ServerAPI().AddPlayer().Remove()
	player := world.FindOne(func(atom platform.IAtom) bool {
		return atom.AsDatum().Type == "/mob/player"
	})
	assert.NotNil(t, player)
	if player == nil {
		return
	}
	turf := player.(platform.IMob).Location()
	assert.NotNil(t, turf)
	if turf == nil {
		return
	}
	contents := turf.(platform.ITurf).Contents()
	assert.True(t, len(contents) > 0)
	assert.Contains(t, contents, player)

	runtime.GC()
	runtime.GC()
	playerAgain := world.FindOne(func(atom platform.IAtom) bool {
		return atom.AsDatum().Type == "/mob/player"
	})
	assert.NotNil(t, playerAgain)
	assert.Equal(t, player, playerAgain)
	contents2 := turf.(platform.ITurf).Contents()
	assert.Equal(t, len(contents), len(contents2))
	assert.Contains(t, contents2, player)
}

func TestWorldRender(t *testing.T) {
	world := BuildWorld()
	api := world.ServerAPI()
	player := api.AddPlayer()
	view := player.Render()
	assert.True(t, len(view.Sprites) > 0)
	runtime.GC()
	view2 := player.Render()
	assert.True(t, len(view2.Sprites) == len(view.Sprites))
}

func TestSingletons(t *testing.T) {
	tree := platform.NewAtomicTree(platform.BaseTreeDefiner{})
	areaOne := tree.New("/area")
	areaTwo := tree.New("/area")
	assert.True(t, areaOne == areaTwo)
	assert.True(t, areaOne.Impl() == areaTwo.Impl())
}

func TestSingletonsInWorld(t *testing.T) {
	world := BuildWorld()
	areaOne := world.Tree.New("/area")
	areaTwo := world.Tree.New("/area")
	assert.True(t, areaOne == areaTwo)
	assert.True(t, areaOne.Impl() == areaTwo.Impl())
}

func TestSingletonAreas(t *testing.T) {
	world := BuildWorld()
	areas := world.FindAll(func(atom platform.IAtom) bool {
		_, isarea := atom.(platform.IArea)
		return isarea
	})
	pathCount := map[datum.TypePath]int{}
	for _, area := range areas {
		pathCount[area.AsDatum().Type] += 1
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

	world := BuildWorld()
	playerAPI := world.ServerAPI().AddPlayer()

	player := world.FindOne(func(atom platform.IAtom) bool {
		return atom.AsDatum().Type == "/mob/player"
	})
	assert.NotNil(t, player)
	player.SetLocation(world.LocateXYZ(11, 4, 1))
	assert.Equal(t, datum.TypePath("/area/outside"), player.ContainingArea().AsDatum().Type)

	playerAPI.Command(webclient.Command{Verb: ".west"})
	lines, sounds := playerAPI.PullRequests()
	assert.Contains(t, lines, "Watch out for the giant rat!")
	assert.Equal(t, 1, len(sounds))
	assert.Equal(t, "cavern.ogg", sounds[0].File)

	assert.Equal(t, datum.TypePath("/area/cave"), player.ContainingArea().AsDatum().Type)
}

func TestFillUpEntireMap(t *testing.T) {
	world := BuildWorld()
	found := false
	serverAPI := world.ServerAPI()
	for i := 0; i < 1000; i++ {
		player := serverAPI.AddPlayer()
		player.Render()
		player.PullRequests()
		player.Remove()
		playerNowhere := world.FindOne(func(atom platform.IAtom) bool {
			return atom.AsDatum().Type == "/mob/player" && atom.Location() == nil
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
