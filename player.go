package whoisspy

type Player struct {
	ID       int64
	Name     string
	Identity int
	Topic    string
	Dead     bool
}

func CreatePlayer(name string) *Player {
	player := &Player{
		ID:       node.Generate().Int64(),
		Name:     name,
		Identity: 0,
		Topic:    "",
		Dead:     false,
	}

	return player
}
