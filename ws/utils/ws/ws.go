package ws

import "github.com/gorilla/websocket"

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var Clients = make(map[*websocket.Conn]bool)

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

var Broadcast = make(chan Message)
