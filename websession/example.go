package main

import (
	"github.com/celskeggs/mediator/webclient"
	"github.com/celskeggs/mediator/webclient/sprite"
	"github.com/celskeggs/mediator/websession"
	"time"
)

type ExampleAPI struct {
	Cheese  sprite.GameSprite
	Updates chan struct{}
}

func (e *ExampleAPI) AddPlayer() websession.PlayerAPI {
	return &ExamplePlayer{
		API: e,
	}
}

func (e *ExampleAPI) SubscribeToUpdates() <-chan struct{} {
	return e.Updates
}

func (e *ExampleAPI) MoveTheCheese() {
	for {
		time.Sleep(time.Second / 50)
		e.Cheese.X += 1
		if e.Cheese.X >= 270 {
			e.Cheese.X = 40
		}
		e.Updates <- struct{}{}
	}
}

type ExamplePlayer struct {
	API *ExampleAPI
}

func (e ExamplePlayer) IsValid() bool {
	return true
}

func (e ExamplePlayer) Remove() {
	// nothing to do
}

func (e ExamplePlayer) Command(cmd webclient.Command) {
	// nothing to do
}

func (e ExamplePlayer) Render() sprite.SpriteView {
	return sprite.SpriteView{
		ViewPortWidth:  320,
		ViewPortHeight: 240,
		Sprites: []sprite.GameSprite{
			e.API.Cheese,
		},
	}
}

func main() {
	api := &ExampleAPI{
		Cheese: sprite.GameSprite{
			Icon: "cheese.dmi",
			X:    150,
			Y:    50,
		},
		Updates: make(chan struct{}),
	}
	go api.MoveTheCheese()
	err := websession.LaunchServer(api, "../resources", "../yourfirstworld/icons")
	// should not get here
	panic(err)
}
