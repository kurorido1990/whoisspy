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
	Ticket   int
	Vote     bool
	ws       *websocket.Conn
}

func CreatePlayer(name string) *Player {
	player := &Player{
		ID:       node.Generate().String(),
		Name:     name,
		Identity: 0,
		Topic:    "",
		Dead:     false,

		Ticket: 0,
		Vote:   false,
	}

	return player
}

func (p *Player) reset() {
	p.Dead = false
	p.Topic = ""
	p.Identity = 0
	p.resetTicket()

	if p.ws != nil {
		data, _ := json.Marshal(WsData{
			Cmd:  "reset",
			Data: nil,
		})

		p.ws.WriteMessage(websocket.TextMessage, data)
	}
}

func (p *Player) alive() bool {
	return !p.Dead
}

func (p *Player) resetTicket() {
	p.Vote = false
	p.Ticket = 0
}

func (p *Player) startGambling(playerList []*Player) {
	if p.ws != nil {
		data, _ := json.Marshal(WsData{
			Cmd:  "startGambling",
			Data: playerList,
		})

		p.ws.WriteMessage(websocket.TextMessage, data)
	}
}

func (p *Player) kickPlayer(kickPlayerName string) {
	if p.ws != nil {
		data, _ := json.Marshal(WsData{
			Cmd:  "kickPlayer",
			Data: kickPlayerName,
		})

		p.ws.WriteMessage(websocket.TextMessage, data)
	}
}

func (p *Player) settlement(winner int) {
	if p.ws != nil {
		data, _ := json.Marshal(WsData{
			Cmd:  "settlement",
			Data: winner,
		})

		p.ws.WriteMessage(websocket.TextMessage, data)
	}
}
