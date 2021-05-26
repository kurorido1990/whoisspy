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
	Speak          bool
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
		Speak:    false,
		Ticket:   0,
		Vote:     false,
	}

	return player
}

func (p *Player) reset() {
	p.Dead = false
	p.Topic = ""
	p.Identity = 0
	p.resetTicket()
	p.resetSpeak()

	data, _ := json.Marshal(WsData{
		Cmd:  "reset",
		Data: nil,
	})

	p.ws.WriteMessage(websocket.TextMessage, data)
	p.reconnectQueue = append(p.reconnectQueue, data)
}

func (p *Player) alive() bool {
	return !p.Dead
}

func (p *Player) resetSpeak() {
	p.Speak = false

	data, _ := json.Marshal(WsData{
		Cmd:  "resetSpeak",
		Data: nil,
	})

	p.ws.WriteMessage(websocket.TextMessage, data)
	p.reconnectQueue = append(p.reconnectQueue, data)
}

func (p *Player) resetTicket() {
	p.Vote = false
	p.Ticket = 0
}

func (p *Player) clearLoseMsg() {
	p.reconnectQueue = [][]byte{}
}

func (p *Player) pushLoseMsg() {
	success := 0
	for _, data := range p.reconnectQueue {
		success++
		p.ws.WriteMessage(websocket.TextMessage, data)
	}

	p.clearLoseMsg()
}

func (p *Player) startGambling(playerList []*Player) {
	data, _ := json.Marshal(WsData{
		Cmd:  "startGambling",
		Data: playerList,
	})

	p.ws.WriteMessage(websocket.TextMessage, data)
	p.reconnectQueue = append(p.reconnectQueue, data)
}

func (p *Player) kickPlayer(kickPlayerName string) {
	data, _ := json.Marshal(WsData{
		Cmd:  "kickPlayer",
		Data: kickPlayerName,
	})

	p.ws.WriteMessage(websocket.TextMessage, data)
	p.reconnectQueue = append(p.reconnectQueue, data)
}

func (p *Player) endVote() {
	data, _ := json.Marshal(WsData{
		Cmd:  "endVote",
		Data: nil,
	})

	p.ws.WriteMessage(websocket.TextMessage, data)
	p.reconnectQueue = append(p.reconnectQueue, data)
}

func (p *Player) settlement(winner int) {
	data, _ := json.Marshal(WsData{
		Cmd:  "settlement",
		Data: winner,
	})

	p.ws.WriteMessage(websocket.TextMessage, data)
	p.reconnectQueue = append(p.reconnectQueue, data)
}
