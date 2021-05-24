package app

type Player struct {
	ID       string
	Name     string
	Identity int
	Topic    string
	Dead     bool
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
