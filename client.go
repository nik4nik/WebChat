package main

import (
	"github.com/gorilla/websocket"
	"strconv"
)

type client struct {
	socket *websocket.Conn
	send   chan []byte // channel on which messages are sent to the user's browser
	room   *room
	ID     int
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- []byte("[" + strconv.Itoa(c.ID) + "] " + string(msg))
		} else {
			break
		}
	}
	c.socket.Close()
}
func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
