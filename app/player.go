package app

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Player struct {
	ID       string
	Name     string
	Identity int
	Topic    string
	Dead     bool
	ws       *websocket.Conn
}

func CreatePlayer(name string) *Player {
	player := &Player{
		ID:       node.Generate().String(),
		Name:     name,
		Identity: 0,
		Topic:    "",
		Dead:     false,
	}

	return player
}

func (p *Player) reset() {
	p.Dead = false
	p.Topic = ""
	p.Identity = 0

	if p.ws != nil {
		data, _ := json.Marshal(WsData{
			Cmd:  "reset",
			Data: nil,
		})

		p.ws.WriteMessage(websocket.TextMessage, data)
	}
}
