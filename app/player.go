package app

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Player struct {
	ID             string
	Name           string
	Identity       int
	Topic          string
	Dead           bool
	Ticket         int
	Vote           bool
	reconnectQueue [][]byte
	ws             *websocket.Conn
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

func (p *Player) pushLoseMsg() {
	success := 0
	for _, data := range p.reconnectQueue {
		if p.ws != nil {
			success++
			p.ws.WriteMessage(websocket.TextMessage, data)
		}
	}

	p.reconnectQueue = p.reconnectQueue[success:len(p.reconnectQueue)]
}

func (p *Player) startGambling(playerList []*Player) {
	data, _ := json.Marshal(WsData{
		Cmd:  "startGambling",
		Data: playerList,
	})

	if p.ws != nil {
		p.ws.WriteMessage(websocket.TextMessage, data)
	} else {
		p.reconnectQueue = append(p.reconnectQueue, data)
	}
}

func (p *Player) kickPlayer(kickPlayerName string) {
	data, _ := json.Marshal(WsData{
		Cmd:  "kickPlayer",
		Data: kickPlayerName,
	})

	if p.ws != nil {
		p.ws.WriteMessage(websocket.TextMessage, data)
	} else {
		p.reconnectQueue = append(p.reconnectQueue, data)
	}
}

func (p *Player) settlement(winner int) {
	data, _ := json.Marshal(WsData{
		Cmd:  "settlement",
		Data: winner,
	})

	if p.ws != nil {
		p.ws.WriteMessage(websocket.TextMessage, data)
	} else {
		p.reconnectQueue = append(p.reconnectQueue, data)
	}
}
