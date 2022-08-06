package main

import (
	"github.com/celskeggs/mediator/common"
	"github.com/celskeggs/mediator/platform/atoms"
	"github.com/celskeggs/mediator/platform/datum"
	"github.com/celskeggs/mediator/platform/types"
	"github.com/celskeggs/mediator/platform/world"
	"github.com/celskeggs/mediator/util"
	"github.com/celskeggs/mediator/webclient"
	"github.com/celskeggs/mediator/webclient/sprite"
	"github.com/celskeggs/mediator/websession"
	"github.com/stretchr/testify/assert"
	"runtime"
	"strings"
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
	gameworld.ServerAPI().AddPlayer("")
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
	gameworld.ServerAPI().AddPlayer("").Remove()
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
	player := api.AddPlayer("")
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
	playerAPI := gameworld.ServerAPI().AddPlayer("")

	player := gameworld.FindOneType("/mob/player")
	assert.NotNil(t, player)
	player.SetVar("loc", gameworld.LocateXYZ(11, 4, 1))
	assert.Equal(t, types.TypePath("/area/outside"), atoms.ContainingArea(player).(*types.Datum).Type())

	lines, sounds, _ := playerAPI.PullRequests()
	assert.Equal(t, 1, len(lines))
	assert.Contains(t, lines, "Nice and jazzy, here...")
	assert.Equal(t, 1, len(sounds))
	if len(sounds) >= 1 {
		assert.Equal(t, "jazzy.ogg", sounds[0].File)
	}

	playerAPI.Command(webclient.Command{Verb: ".west"})
	lines, sounds, _ = playerAPI.PullRequests()
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
	playerAPI := gameworld.ServerAPI().AddPlayer("")

	player := gameworld.FindOneType("/mob/player")
	assert.NotNil(t, player)
	assert.Equal(t, types.TypePath("/area/outside"), atoms.ContainingArea(player).(*types.Datum).Type())
	cx, cy := world.XY(player)
	assert.Equal(t, uint(1), cx)
	assert.Equal(t, uint(1), cy)

	lines, sounds, _ := playerAPI.PullRequests()
	assert.Equal(t, 1, len(lines))
	assert.Contains(t, lines, "Nice and jazzy, here...")
	assert.Equal(t, 1, len(sounds))
	if len(sounds) >= 1 {
		assert.Equal(t, "jazzy.ogg", sounds[0].File)
	}

	playerAPI.Command(webclient.Command{Verb: ".south"})
	lines, sounds, _ = playerAPI.PullRequests()
	assert.Equal(t, 0, len(lines))
	assert.Equal(t, 0, len(sounds))
	nx, ny := world.XY(player)
	assert.Equal(t, uint(1), nx)
	assert.Equal(t, uint(1), ny)
}

func TestWalkIntoWalls(t *testing.T) {
	gameworld := BuildWorld()
	playerAPI := gameworld.ServerAPI().AddPlayer("")

	player := gameworld.FindOneType("/mob/player")
	assert.NotNil(t, player)
	_, _, _ = playerAPI.PullRequests()

	player.SetVar("loc", gameworld.LocateXYZ(6, 4, 1))
	assert.Equal(t, types.TypePath("/area/outside"), atoms.ContainingArea(player).(*types.Datum).Type())

	playerAPI.Command(webclient.Command{Verb: ".north"})
	lines, sounds, _ := playerAPI.PullRequests()
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
		player := serverAPI.AddPlayer("")
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
	playerAPI := gameworld.ServerAPI().AddPlayer("")

	player := gameworld.FindOneType("/mob/player")
	assert.NotNil(t, player)
	_, _, _ = playerAPI.PullRequests()

	playerAPI.Command(webclient.Command{Verb: "look"})
	lines, _, _ := playerAPI.PullRequests()
	assert.Equal(t, 2, len(lines))
	assert.Contains(t, lines, "You see...")
	assert.Contains(t, lines, "The scroll.  It looks to be rather old.")
}

func iconCounts(playerAPI websession.PlayerAPI) map[string]int {
	view := playerAPI.Render()
	iconCount := map[string]int{}
	for _, sprite := range view.Sprites {
		iconCount[sprite.Icon] += 1
	}
	return iconCount
}

func TestGetDropVerbs(t *testing.T) {
	gameworld := BuildWorld()
	playerAPI := gameworld.ServerAPI().AddPlayer("")

	player := gameworld.FindOneType("/mob/player")
	assert.NotNil(t, player)
	scroll := gameworld.FindOneType("/obj/scroll")
	assert.NotNil(t, scroll)
	_, _, _ = playerAPI.PullRequests()

	for x := types.Unuint(player.Var("x")); x < types.Unuint(scroll.Var("x")); x++ {
		playerAPI.Command(webclient.Command{Verb: ".east"})
	}
	for y := types.Unuint(player.Var("y")); y < types.Unuint(scroll.Var("y")); y++ {
		playerAPI.Command(webclient.Command{Verb: ".north"})
	}

	playerAPI.Command(webclient.Command{Verb: "get"})
	lines, _, _ := playerAPI.PullRequests()
	assert.Equal(t, 1, len(lines))
	assert.Contains(t, lines, "You get the scroll.")

	iconCount := iconCounts(playerAPI)
	// should only be wall, floor, and player, not scroll
	assert.Equal(t, 3, len(iconCount))
	assert.Equal(t, 0, iconCount["scroll.dmi"])
	assert.Contains(t, iconCount, "wall.dmi")
	assert.Contains(t, iconCount, "floor.dmi")
	assert.Equal(t, 1, iconCount["player.dmi"])

	assert.Equal(t, player, scroll.Var("loc"))

	playerAPI.Command(webclient.Command{Verb: "drop"})
	lines, _, _ = playerAPI.PullRequests()
	assert.Equal(t, 1, len(lines))
	assert.Contains(t, lines, "You drop the scroll.")

	assert.Equal(t, player.Var("loc"), scroll.Var("loc"))

	playerAPI.Command(webclient.Command{Verb: "get scroll"})
	lines, _, _ = playerAPI.PullRequests()
	assert.Equal(t, 1, len(lines))
	assert.Contains(t, lines, "You get the scroll.")

	assert.Equal(t, player, scroll.Var("loc"))

	playerAPI.Command(webclient.Command{Verb: "drop cheese"})
	lines, _, _ = playerAPI.PullRequests()
	assert.Equal(t, 1, len(lines))
	assert.Contains(t, lines, "Not a known verb: \"drop\"")

	// location should stay the same
	assert.Equal(t, player, scroll.Var("loc"))

	playerAPI.Command(webclient.Command{Verb: "drop scroll"})
	lines, _, _ = playerAPI.PullRequests()
	assert.Equal(t, 1, len(lines))
	assert.Contains(t, lines, "You drop the scroll.")

	assert.Equal(t, player.Var("loc"), scroll.Var("loc"))

	iconCount = iconCounts(playerAPI)
	assert.Equal(t, 1, iconCount["scroll.dmi"])
}

func TestEatVerb(t *testing.T) {
	gameworld := BuildWorld()
	playerAPI := gameworld.ServerAPI().AddPlayer("")

	player := gameworld.FindOneType("/mob/player")
	assert.NotNil(t, player)
	cheese := gameworld.FindOneType("/obj/cheese")
	assert.NotNil(t, cheese)
	ok := types.Unint(player.Invoke(nil, "Move", cheese.Var("loc")))
	assert.Equal(t, 1, ok)
	_, _, _ = playerAPI.PullRequests()

	playerAPI.Command(webclient.Command{Verb: "eat cheese"})
	lines, _, _ := playerAPI.PullRequests()
	assert.Equal(t, 1, len(lines))
	assert.Contains(t, lines, "Not a known verb: \"eat\"")

	playerAPI.Command(webclient.Command{Verb: "get cheese"})
	lines, _, _ = playerAPI.PullRequests()
	assert.Equal(t, 1, len(lines))
	assert.Contains(t, lines, "You get the cheese.")

	assert.Equal(t, "", types.Unstring(cheese.Var("suffix")))

	playerAPI.Command(webclient.Command{Verb: "eat cheese"})
	lines, _, _ = playerAPI.PullRequests()
	assert.Equal(t, 1, len(lines))
	assert.Contains(t, lines, "You take a bite of the cheese. Bleck!")

	assert.Equal(t, "(nibbled)", types.Unstring(cheese.Var("suffix")))

	playerAPI.Command(webclient.Command{Verb: "drop cheese"})
	lines, _, _ = playerAPI.PullRequests()
	assert.Equal(t, 1, len(lines))
	assert.Contains(t, lines, "You drop the cheese.")

	assert.Equal(t, "(nibbled)", types.Unstring(cheese.Var("suffix")))

	util.FIXME("test that suffix shows up in stat panel")
}

func TestReadVerb(t *testing.T) {
	gameworld := BuildWorld()
	playerAPI := gameworld.ServerAPI().AddPlayer("")

	player := gameworld.FindOneType("/mob/player")
	assert.NotNil(t, player)
	scroll := gameworld.FindOneType("/obj/scroll")
	assert.NotNil(t, scroll)
	ok := types.Unint(player.Invoke(nil, "Move", scroll.Var("loc")))
	assert.Equal(t, 1, ok)
	_, _, _ = playerAPI.PullRequests()

	rats := gameworld.FindAllType("/mob/rat")
	assert.Equal(t, 1, len(rats))

	playerAPI.Command(webclient.Command{Verb: "get scroll"})
	lines, _, _ := playerAPI.PullRequests()
	assert.Equal(t, 1, len(lines))
	assert.Contains(t, lines, "You get the scroll.")

	assert.Equal(t, player, scroll.Var("loc"))
	countScrollsInPlayer := 0
	for _, element := range datum.Elements(player.Var("contents")) {
		if types.IsType(element, "/obj/scroll") {
			countScrollsInPlayer++
		}
	}
	assert.Equal(t, 1, countScrollsInPlayer)
	assert.Equal(t, 1, len(gameworld.FindAllType("/obj/scroll")))

	playerAPI.Command(webclient.Command{Verb: "read scroll"})
	lines, _, _ = playerAPI.PullRequests()
	assert.Equal(t, 2, len(lines))
	assert.Contains(t, lines, "You utter the phrase written on the scroll: \"Knuth\"!")
	assert.Contains(t, lines, "A giant rat appears!")

	nrats := gameworld.FindAllType("/mob/rat")
	assert.Equal(t, 2, len(nrats))
	if len(rats) >= 1 {
		var newrat types.Value
		for _, rat := range nrats {
			if rat != rats[0] {
				newrat = rat
			}
		}
		assert.NotNil(t, newrat)
		if newrat != nil {
			assert.Equal(t, player.Var("loc"), newrat.Var("loc"))
		}
	}

	countScrollsInPlayer = 0
	for _, element := range datum.Elements(player.Var("contents")) {
		if types.IsType(element, "/obj/scroll") {
			countScrollsInPlayer++
		}
	}
	assert.Equal(t, 0, countScrollsInPlayer)
	assert.Equal(t, 0, len(gameworld.FindAllType("/obj/scroll")))

	util.FIXME("confirm that scroll is no longer in stat panel")
}

func TestStatPanel(t *testing.T) {
	gameworld := BuildWorld()
	serverAPI := gameworld.ServerAPI()
	playerAPI := serverAPI.AddPlayer("")

	player := gameworld.FindOneType("/mob/player")
	assert.NotNil(t, player)

	assert.Equal(t, player, player.Var("client").Var("statobj"))

	view := playerAPI.Render()
	assert.Equal(t, sprite.StatDisplay{}, view.Stats)

	serverAPI.Tick()

	view = playerAPI.Render()
	assert.NotEqual(t, sprite.StatDisplay{}, view.Stats)
}

func TestDisappearingRat(t *testing.T) {
	gameworld := BuildWorld()
	serverAPI := gameworld.ServerAPI()
	playerAPI := serverAPI.AddPlayer("")

	player := gameworld.FindOneType("/mob/player")
	assert.NotNil(t, player)
	cheese := gameworld.FindOneType("/obj/cheese")
	assert.NotNil(t, cheese)
	ok := types.Unint(player.Invoke(nil, "Move", cheese.Var("loc")))
	assert.Equal(t, 1, ok)
	_, _, _ = playerAPI.PullRequests()

	playerAPI.Command(webclient.Command{Verb: ".south"})
	playerAPI.Command(webclient.Command{Verb: "look"})
	lines, _, _ := playerAPI.PullRequests()
	assert.Equal(t, 2, len(lines))
	// Weirdly enough, the player *can* see the rat from here in their rendered viewport, but the oview proc says they
	// can't. This is because the view area is messed with in the case where the map is small, but that doesn't do
	// anything for oview, so the two end up being different. This is expected behavior! Huh.
	assert.Contains(t, lines, "You see...")
	assert.Contains(t, lines, "The cheese.  It is quite smelly.")

	playerAPI.Command(webclient.Command{Verb: ".west"})
	playerAPI.Command(webclient.Command{Verb: "look"})
	lines, _, _ = playerAPI.PullRequests()
	assert.Equal(t, 3, len(lines))
	assert.Contains(t, lines, "You see...")
	assert.Contains(t, lines, "The cheese.  It is quite smelly.")
	assert.Contains(t, lines, "The rat.  It's quite large.")
}

func TestBumpAngryRat(t *testing.T) {
	gameworld := BuildWorld()
	serverAPI := gameworld.ServerAPI()
	playerAPI := serverAPI.AddPlayer("")

	player := gameworld.FindOneType("/mob/player")
	assert.NotNil(t, player)
	rat := gameworld.FindOneType("/mob/rat")
	assert.NotNil(t, rat)
	x, y, z := world.XYZ(rat)
	loc := gameworld.LocateXYZ(x+1, y, z)
	assert.NotNil(t, loc)
	ok := types.Unint(player.Invoke(nil, "Move", loc))
	assert.Equal(t, 1, ok)
	_, _, flicks := playerAPI.PullRequests()
	assert.Equal(t, 0, len(flicks))

	playerAPI.Command(webclient.Command{Verb: ".west"})
	lines, sounds, flicks := playerAPI.PullRequests()
	assert.Equal(t, 2, len(lines))
	if len(lines) >= 1 {
		assert.Contains(t, lines, "You bump into the rat.")
		assert.Contains(t, lines, "The giant rat defends its territory ferociously!")
	}
	assert.Equal(t, 1, len(sounds))
	if len(sounds) >= 1 {
		assert.Equal(t, "ouch.wav", sounds[0].File)
	}
	assert.Equal(t, 1, len(flicks))
	if len(flicks) >= 1 {
		assert.Equal(t, "rat.dmi", flicks[0].Icon)
		assert.Equal(t, 3, len(flicks[0].Frames))
		assert.Equal(t, rat.(*types.Datum).UID(), flicks[0].UID)
	}

	assert.Equal(t, common.East, rat.Var("dir"))
}

func TestRatScurryingTowardsCheese(t *testing.T) {
	gameworld := BuildWorld()
	serverAPI := gameworld.ServerAPI()

	cheese := gameworld.FindOneType("/obj/cheese")
	assert.NotNil(t, cheese)
	rat := gameworld.FindOneType("/mob/rat")
	assert.NotNil(t, rat)
	x, y, z := world.XYZ(rat)
	l0 := gameworld.LocateXYZ(x, y, z)
	l1 := gameworld.LocateXYZ(x+1, y, z)
	l2 := gameworld.LocateXYZ(x+2, y, z)
	l3 := gameworld.LocateXYZ(x+3, y, z)
	assert.NotNil(t, l0)
	assert.NotNil(t, l1)
	assert.NotNil(t, l2)
	assert.NotNil(t, l3)
	assert.Equal(t, l0, rat.Var("loc"))
	ok := types.Unint(cheese.Invoke(nil, "Move", l3))
	assert.Equal(t, 1, ok)

	ws := atoms.GetWalkState(rat)
	assert.Equal(t, l3, ws.WalkTarget)

	for i := 0; i < 5; i++ {
		assert.Equal(t, l0, rat.Var("loc"))
		assert.Equal(t, common.South, rat.Var("dir"))
		serverAPI.Tick()
	}
	for i := 0; i < 5; i++ {
		assert.Equal(t, l1, rat.Var("loc"))
		assert.Equal(t, common.East, rat.Var("dir"))
		serverAPI.Tick()
	}
	for i := 0; i < 100; i++ {
		assert.Equal(t, l2, rat.Var("loc"))
		assert.Equal(t, common.East, rat.Var("dir"))
		serverAPI.Tick()
	}
}

func TestPlayerIsGuest(t *testing.T) {
	gameworld := BuildWorld()

	gameworld.ServerAPI().AddPlayer("")
	player := gameworld.FindOneType("/mob/player")
	assert.NotNil(t, player)
	client := player.Var("client")
	key := types.Unstring(client.Var("key"))
	assert.True(t, strings.HasPrefix(key, "Guest-"), "expected key to start with \"Guest-\", but was %q", key)
	name := types.Unstring(player.Var("name"))
	assert.Equal(t, key, name)
}

func TestPlayerIsNamed(t *testing.T) {
	gameworld := BuildWorld()

	expectedName := "MyPlayerName"

	gameworld.ServerAPI().AddPlayer(expectedName)
	player := gameworld.FindOneType("/mob/player")
	assert.NotNil(t, player)
	client := player.Var("client")
	key := types.Unstring(client.Var("key"))
	assert.Equal(t, expectedName, key)
	name := types.Unstring(player.Var("name"))
	assert.Equal(t, expectedName, name)
}

func TestSay(t *testing.T) {
	gameworld := BuildWorld()
	serverAPI := gameworld.ServerAPI()

	speakername := "SamplePlayerName"

	speakerPlayerAPI := serverAPI.AddPlayer(speakername)
	speakerPlayer := gameworld.FindOneType("/mob/player")
	assert.NotNil(t, speakerPlayer)

	hearerPlayerAPI := serverAPI.AddPlayer("")
	hearerPlayer := gameworld.FindOne(func(t *types.Datum) bool {
		return types.IsType(t, "/mob/player") && t != speakerPlayer
	})
	assert.NotNil(t, hearerPlayer)

	distantPlayerAPI := serverAPI.AddPlayer("")
	distantPlayer := gameworld.FindOne(func(t *types.Datum) bool {
		return types.IsType(t, "/mob/player") && t != speakerPlayer && t != hearerPlayer
	})
	assert.NotNil(t, distantPlayer)

	x, y, z := world.XYZ(speakerPlayer)

	ok := types.Unint(hearerPlayer.Invoke(nil, "Move", gameworld.LocateXYZ(x+5, y, z)))
	assert.Equal(t, 1, ok)
	ok = types.Unint(distantPlayer.Invoke(nil, "Move", gameworld.LocateXYZ(x+6, y, z)))
	assert.Equal(t, 1, ok)

	// ignore any messages so far
	speakerPlayerAPI.PullRequests()
	hearerPlayerAPI.PullRequests()
	distantPlayerAPI.PullRequests()

	util.FIXME("The 'say' command with no argument should actually pop up a dialog for the player to enter a string")

	for verb, display := range map[string]string {
		"say": speakername + " says, \"\"",
		"say \"": speakername + " says, \"\"",
		"say \"test123\"": speakername + " says, \"test123\"",
		"say \"hello world\"": speakername + " says, \"hello world\"",
		"say \"one two three": speakername + " says, \"one two three\"",
	} {
		speakerPlayerAPI.Command(webclient.Command{
			Verb: verb,
		})

		for who, playerAPI := range map[string]websession.PlayerAPI{ "speaker": speakerPlayerAPI, "hearer": hearerPlayerAPI } {
			lines, sounds, flicks := playerAPI.PullRequests()
			assert.Equal(t, 1, len(lines), "during verb %q for %s, actual was: %v", verb, who, lines)
			assert.Contains(t, lines, display, "during verb %q for %s", verb, who)
			assert.Equal(t, 0, len(sounds))
			assert.Equal(t, 0, len(flicks))
		}

		lines, sounds, flicks := distantPlayerAPI.PullRequests()
		assert.Equal(t, 0, len(lines), "during verb %q", verb)
		assert.Equal(t, 0, len(sounds))
		assert.Equal(t, 0, len(flicks))
	}
}

func TestSectionE(t *testing.T) {
	util.FIXME("uncomment the rest of the things in Section E")
	t.Error("unimplemented")
}
