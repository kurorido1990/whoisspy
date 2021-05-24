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
	Players    []*Player
	Citizens   []*Player
	Spy        []*Player
}

func createRoom(maxLimit int) string {
	room := &Room{
		ID:         node.Generate().String(),
		Status:     RoomStatusPrepare,
		TopicIndex: getTopicIndex(),
		MaxLimit:   maxLimit,
		SpyNum:     getSpyNum(maxLimit),
		Players:    make([]*Player, 0),
		Citizens:   make([]*Player, 0),
		Spy:        make([]*Player, 0),
	}

	roomList.Store(room.ID, room)

	return room.ID
}

func (r *Room) start() {
	r.Status = RoomStatusStart
}

func (r *Room) addPlayer(player *Player) error {
	if r.isMax() {
		return fmt.Errorf("RoomID : %d 人數已滿", r.ID)
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

func (r *Room) kickPlayer(playerID string) {
	for _, player := range r.Players {
		if playerID == player.ID {
			player.Dead = true
			return
		}
	}
}

func (r *Room) settlement() int {
	if r.isNum(SPY) < 1 {
		r.Status = RoomStatusEnd
		return Result_CITIZEN_WIN
	}

	if r.isNum(0) <= winNum {
		r.Status = RoomStatusEnd
		return Result_SPY_WIN
	}

	return Result_Continue
}
