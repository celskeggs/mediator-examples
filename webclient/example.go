package main

import (
	"github.com/celskeggs/mediator/webclient"
	"github.com/celskeggs/mediator/webclient/sprite"
	"log"
	"time"
)

type ExampleServer struct {
	LastID uint
}

func (e ExampleServer) ListResources() (map[string]string, error) {
	return map[string]string{
		"cheese.dmi": "../yourfirstworld/icons/cheese.dmi",
	}, nil
}

func (e ExampleServer) CoreResourcePath() string {
	return "../resources"
}

func (e ExampleServer) Connect() webclient.ServerSession {
	e.LastID += 1
	log.Println("opened session: ", e.LastID)
	return &ExampleSession{
		ID: e.LastID,
	}
}

type ExampleSession struct {
	ID uint
}

func (e ExampleSession) Close() {
	log.Println("closed session: ", e.ID)
}

type MessageHolder struct {
	Test string
}

func (e ExampleSession) NewMessageHolder() interface{} {
	return &MessageHolder{}
}

func (e ExampleSession) OnMessage(message interface{}) {
	mh := message.(*MessageHolder)
	log.Println("got message holder for: ", e.ID, " with message ", mh.Test)
}

func (e ExampleSession) BeginSend(send func(*sprite.SpriteView) error) {
	go func() {
		defer func() {
			_ = send(nil)
		}()
		for {
			log.Println("sending first message for: ", e.ID)
			if send(&sprite.SpriteView{
				Sprites: []sprite.GameSprite{
					{
						Icon: "cheese.dmi",
						X:    128,
						Y:    128,
					},
				},
			}) != nil {
				break
			}
			log.Println("sent message for: ", e.ID)
			time.Sleep(time.Second / 2)
			if send(&sprite.SpriteView{
				Sprites: []sprite.GameSprite{
					{
						Icon: "cheese.dmi",
						X:    128,
						Y:    192,
					},
				},
			}) != nil {
				break
			}
			log.Println("sent another message for: ", e.ID)
			time.Sleep(time.Second / 2)
		}
	}()
}

func main() {
	es := ExampleServer{}
	err := webclient.LaunchHTTP(es)
	panic(err)
}
