package client

import (
	"bytes"
	"io"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

// WebsocketClient is the base client to consume the websockets
type WebsocketClient struct {
	conn *websocket.Conn
}

// Consume creates a channel that returns messages
func (c *WebsocketClient) Consume() chan interface{} {
	messages := make(chan interface{})
	go func() {
		for {
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				return
			}
			messages <- message
		}
	}()
	return messages
}

// WriteToDisk writes new messages to disk
func (c *WebsocketClient) WriteToDisk(filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	go func() {
		for {
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				return
			}
			io.Copy(out, bytes.NewReader(message))
		}
	}()
	return nil
}

// NewClient creates a new websocket client
func NewClient(endpoint string) (*WebsocketClient, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return &WebsocketClient{
		conn: c,
	}, nil
}
