package main

import (
	"github.com/celskeggs/mediator/platform"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

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
