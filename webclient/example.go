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

var _ webclient.ServerSession = &ExampleSession{}

func (e ExampleSession) Close() {
	log.Println("closed session: ", e.ID)
}

func (e ExampleSession) OnMessage(message webclient.Command) {
	log.Println("got verb:", message.Verb)
}

func (e ExampleSession) BeginSend(send func(update *sprite.ViewUpdate) error) {
	go func() {
		defer func() {
			_ = send(nil)
		}()
		for {
			log.Println("sending first message for: ", e.ID)
			if send(&sprite.ViewUpdate{
				NewState: &sprite.SpriteView{
					Sprites: []sprite.GameSprite{
						{
							Icon: "cheese.dmi",
							X:    128,
							Y:    128,
						},
					},
				},
			}) != nil {
				break
			}
			log.Println("sent message for: ", e.ID)
			time.Sleep(time.Second / 2)
			if send(&sprite.ViewUpdate{
				NewState: &sprite.SpriteView{
					Sprites: []sprite.GameSprite{
						{
							Icon: "cheese.dmi",
							X:    128,
							Y:    192,
						},
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
