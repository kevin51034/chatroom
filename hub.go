package main

import (
	"time"
)
// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	//broadcast chan []byte
	broadcast chan formatMessage


	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

type formatMessage struct {
	Username string
	Message string
	Room string
    Time string
}

func newHub() *Hub {
	return &Hub{
		//broadcast:  make(chan []byte),
		broadcast:  make(chan formatMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			if message.Message == "welcome" || message.Message == "leave" {
				for client := range h.clients {
					if client.room == message.Room {
						if client.username == message.Username {
							welcomemsg := formatMessage{Username:"ChatBot", Room: client.room, 
							Message: "Welcome to the chat room " + message.Username, Time: time.Now().Format("3:04 pm")}
							select {
							case client.send <- welcomemsg:
							default:
								close(client.send)
								delete(h.clients, client)
							}
						} else {
							var msg formatMessage
							if message.Message == "welcome" {
								msg = formatMessage{Username:"ChatBot", Room: client.room, 
								Message: message.Username+ " has entered the chat", Time: time.Now().Format("3:04 pm")}
							} else if message.Message == "leave" {
								msg = formatMessage{Username:"ChatBot", Room: client.room, 
								Message: message.Username+ " has left the chat", Time: time.Now().Format("3:04 pm")}
							}
							select {
							case client.send <- msg:
							default:
								close(client.send)
								delete(h.clients, client)
							}
						}
					}
				}
			} else {
				for client := range h.clients {
					if client.room == message.Room {
						select {
						case client.send <- message:
						default:
							close(client.send)
							delete(h.clients, client)
						}
					}
				}
			}
		}
	}
}