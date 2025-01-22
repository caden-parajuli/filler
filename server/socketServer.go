package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	Subprotocols: []string{"JSON-v1"},

}

// Handles HTTP upgrade and listens for WebSocket messages
func wsHandler(writer http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	var client *Client = NewClient(conn)

	go client.writePump()
	go client.readPump()
}
