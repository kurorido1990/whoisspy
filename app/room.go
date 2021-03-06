package app

import (
	"fmt"
)

type Room struct {
	ID         string
	Status     int
	TopicIndex int
	MaxLimit   int
	SpyNum     int
	Gambling   bool
	Ticket     int
	Speak      int
	Players    []*Player
	Citizens   []*Player
	Spy        []*Player
	Round      int
}

func createRoom(maxLimit int) string {
	room := &Room{
		ID:         node.Generate().String(),
		Status:     RoomStatusPrepare,
		Gambling:   false,
		Ticket:     0,
		Speak:      0,
		TopicIndex: getTopicIndex(),
		MaxLimit:   maxLimit,
		SpyNum:     getSpyNum(maxLimit),
		Players:    make([]*Player, 0),
		Citizens:   make([]*Player, 0),
		Spy:        make([]*Player, 0),
		Round:      0,
	}

	roomList.Store(room.ID, room)

	return room.ID
}

func (r *Room) start() {
	r.Status = RoomStatusStart
}

func (r *Room) addPlayer(player *Player) error {
	if r.isMax() {
		return fmt.Errorf("人數已滿")
	}

	if player.Name == "" {
		return fmt.Errorf("名字不能為空")
	}

	for _, p := range r.Players {
		if p.Name == player.Name {
			return fmt.Errorf("名字重複")
		}
	}

	topicList := getTopic(r.TopicIndex)
	switch gen.Identity() {
	case CITIZEN:
		if r.MaxLimit-r.isNum(CITIZEN) > r.SpyNum {
			r.Citizens = append(r.Citizens, player)
			player.Identity = CITIZEN
			player.Topic = topicList[CITIZEN-1]
		} else {
			r.Spy = append(r.Spy, player)
			player.Identity = SPY
			player.Topic = topicList[SPY-1]
		}
	case SPY:
		if r.SpyNum == r.isNum(SPY) {
			r.Citizens = append(r.Citizens, player)
			player.Identity = CITIZEN
			player.Topic = topicList[CITIZEN-1]
		} else {
			r.Spy = append(r.Spy, player)
			player.Identity = SPY
			player.Topic = topicList[SPY-1]
		}
	}

	r.Players = append(r.Players, player)

	return nil
}

func (r *Room) isNum(identity int) int {
	switch identity {
	case CITIZEN:
		return len(r.Citizens)
	case SPY:
		return len(r.Spy)
	default:
		return len(r.Players)
	}
}

func (r *Room) isMax() bool {
	if r.isNum(0) == r.MaxLimit {
		return true
	} else {
		return false
	}
}

func (r *Room) startGambling() {
	r.Gambling = true

	alivePlayer := r.getAlivePlayer()
	for _, player := range alivePlayer {
		if !player.Vote {
			player.startGambling(alivePlayer)
		}
	}
}

func (r *Room) playerSpeak() {
	r.Speak++

	if r.Speak >= len(r.getAlivePlayer()) {
		r.startGambling()
	}
}

func (r *Room) resetPlayerSpeak() {
	for _, player := range r.getAlivePlayer() {
		player.resetSpeak()
	}
}

func (r *Room) votePlayer(playerID string) {
	for _, player := range r.Players {
		if playerID == player.ID {
			player.Ticket++
			r.Ticket++
			break
		}
	}

	if r.Ticket >= len(r.getAlivePlayer()) {
		r.stopGambling()
	}
}

func (r *Room) stopGambling() {
	if !r.Gambling {
		return
	}

	r.Gambling = false
	r.Ticket = 0
	playerID := ""
	topTicket := 0
	for _, player := range r.getAlivePlayer() {
		if player.Ticket > 0 {
			if player.Ticket > topTicket {
				playerID = player.ID
				topTicket = player.Ticket
			}
		}
	}

	if playerID != "" {
		r.kickPlayer(playerID)
	}
}

func (r *Room) getAlivePlayer() []*Player {
	alivePlayer := make([]*Player, 0)
	for _, player := range r.Players {
		if player.alive() {
			alivePlayer = append(alivePlayer, player)
		}
	}

	return alivePlayer
}

func (r *Room) kickPlayer(playerID string) {
	var kickPlayerName string
	for _, player := range r.Players {
		if playerID == player.ID {
			player.Dead = true
			kickPlayerName = player.Name
			break
		}
	}

	for _, player := range r.Players {
		player.kickPlayer(kickPlayerName)
	}

	r.settlement()
}

func (r *Room) resetPlayerVote() {
	for _, player := range r.Players {
		player.resetTicket()
	}
}

func (r *Room) getAliveSpy() int {
	alive := 0
	for _, player := range r.getAlivePlayer() {
		if player.Identity == SPY {
			alive++
		}
	}

	return alive
}

func (r *Room) settlement() {
	winner := 0
	if r.getAliveSpy() < 1 {
		r.Status = RoomStatusEnd
		winner = Result_CITIZEN_WIN
	} else if len(r.getAlivePlayer()) <= winNum {
		r.Status = RoomStatusEnd
		winner = Result_SPY_WIN
	} else if r.Round > 2 {
		r.Status = RoomStatusEnd
		winner = Result_SPY_WIN
	}

	if winner < 1 {
		r.Round++
		r.resetPlayerSpeak()
		r.resetPlayerVote()
	} else {
		for _, player := range r.Players {
			player.settlement(winner)
		}

		r.resetGame()
	}
}

func (r *Room) resetGame() {
	playerList := r.Players
	r.TopicIndex = getTopicIndex()
	r.Status = Result_Continue
	r.Gambling = false
	r.Ticket = 0
	r.Speak = 0
	r.Round = 0

	r.Players = make([]*Player, 0)
	r.Spy = make([]*Player, 0)
	r.Citizens = make([]*Player, 0)

	for _, player := range playerList {
		player.reset()
		r.addPlayer(player)
	}
}
