package websocket

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan Message
	UserID uint
}

type Message struct {
	UserID  uint   `json:"user_id"`
	Content string `json:"content"`
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			break
		}
		msg.UserID = c.UserID
		c.Hub.broadcast <- msg
	}
}

// Sends messages from Hub to WebSocket
func (c *Client) WritePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		err := c.Conn.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}
