package app

type IndexData struct {
	Title   string
	Content string
}

type WsData struct {
	Cmd  string
	Data interface{}
}
