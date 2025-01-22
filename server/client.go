package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// Expected close codes
var closeCodes = []int{websocket.CloseGoingAway, websocket.CloseNormalClosure}

type Client struct {
	// Websocket connection
	conn *websocket.Conn
	// Client ID
	id string
	// Game ID, 0 if unassigned
	game uint64
	// Write queue
	out chan MessageRaw
}

const NO_GAME uint64 = 0
const NO_ID string = ""

func NewClient(conn *websocket.Conn) *Client {
	client := &Client{
		conn,
		NO_ID,
		NO_GAME,
		make(chan MessageRaw),
	}

	clientsLock.Lock()
	clients[client.id] = client
	clientsLock.Unlock()

	return client
}

func (client *Client) readPump() {
	for {
		mtype, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, closeCodes...) {
				log.Println("Received close message",)
			} else {
				log.Print("websocket read: ", err)
			}

			// Maybe close channel?
			clientsLock.Lock()
			delete(clients, client.id)
			clientsLock.Unlock()
			return
		}

		log.Printf("received (%d): %s", mtype, message)

		// TODO add deadlines
		switch mtype {
		case websocket.TextMessage:
			client.handleTextMessage(message)
		case websocket.BinaryMessage:
			client.handleBinaryMessage(message)
		default:
			log.Println("WARNING: Bad message type")
		}
	}
}

func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.out:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				log.Println("Client channel closed")
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			messageStr, err := json.Marshal(message)
			if err != nil {
				log.Print("Marshal: ", err)
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Print("NextWriter: ", err)
				return
			}
			w.Write(messageStr)
			log.Printf("sent: %s", messageStr)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) write(message MessageRaw) {
	client.out <- message
}

func TryWrite(id string, message MessageRaw) bool {
	clientsLock.RLock()
	client := clients[id]
	clientsLock.RUnlock()

	if client == nil {
		return false
	}
	client.write(message)
	return true
}

func (client *Client) Close() {
	log.Print("Closing client")
	close(client.out)

	clientsLock.Lock()
	delete(clients, client.id)
	clientsLock.Unlock()
}

func (client *Client) isValid(id string) bool {
	if client.id == "" {
		return false
	} else if client.id != id {
		client.Close()
		return false
	}
	return true
}
