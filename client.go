package main

import (
	"bytes"
	"log"
	"net/http"
	"time"
	"fmt"
	"encoding/json"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	// send chan []byte
	send chan formatMessage

	username string
	room 	 string
}


// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	//c.hub.broadcast <- formatMessage{Username:c.username, Room: c.room, Message: "welcome", Time: time.Now()}
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		// broadcast the message to other users
		c.hub.broadcast <- formatMessage{Username:c.username, Room: c.room, Message: string(message), Time: time.Now().Format("3:04 pm")}
		fmt.Println(message)
		fmt.Println(c.hub.clients)
		// try response
		/*
		reply := string(message) + " receive message from client and send it back!"
		if err := c.conn.WriteMessage(mt, []byte(reply)); err != nil {
			log.Println(err)
			return
		}
		*/
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.hub.broadcast <- formatMessage{Username:c.username, Room: c.room, Message: "leave", Time: time.Now().Format("3:04 pm")}

		ticker.Stop()
		c.conn.Close()
	}()
	c.hub.broadcast <- formatMessage{Username:c.username, Room: c.room, Message: "welcome", Time: time.Now().Format("3:04 pm")}

	for {
		select {
		// formatMessage send from hub.broadcast
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			b, err := json.Marshal(msg) 
			fmt.Println(msg)
			//fmt.Println(b)

			w.Write(b)

			// Add queued chat messages to the current websocket message.
			/*n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				//w.Write([]byte(<-c.send.message))
			}*/

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}


// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	var username, room string
	if _, ok := r.URL.Query()["username"]; ok {
		fmt.Println(r.URL.Query()["username"])
		username = r.URL.Query()["username"][0]
	}
	if _, ok := r.URL.Query()["room"]; ok {
		fmt.Println(r.URL.Query()["room"])
		room = r.URL.Query()["room"][0]
	}

	client := &Client{hub: hub, conn: conn, send: make(chan formatMessage), username: username, room: room}
	//fmt.Println(client)

	client.hub.register <- client
	
	//msg := formatMessage{Username:username, Room: room, Message: "welcome", Time: time.Now()}
	//client.hub.broadcast <- msg

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()

}