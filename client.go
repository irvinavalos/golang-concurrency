package main

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type MessageType string

const (
	MessageType_Broadcast MessageType = "broadcast"
)

type RequestMessage struct {
	messageType MessageType
	client      *Client
	data        string
}

type ResponseMessage struct {
	messageType MessageType
	data        string
	senderID    string
}

func NewResponseMessage(msg *RequestMessage) *ResponseMessage {
	return &ResponseMessage{
		messageType: msg.messageType,
		data:        msg.data,
		senderID:    msg.client.ID,
	}
}

type Client struct {
	ID   string
	mu   *sync.RWMutex
	conn *websocket.Conn
}

func NewClient(conn *websocket.Conn) *Client {
	id := rand.Text()[:10]
	return &Client{
		ID:   id,
		mu:   new(sync.RWMutex),
		conn: conn,
	}
}

func (c *Client) readMessageLoop(leaveServerCH chan<- *Client, broadcastCH chan<- *RequestMessage) {
	defer func() {
		c.conn.Close()
		leaveServerCH <- c
	}()

	for {
		_, msgBytes, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		msg := new(RequestMessage)
		err = json.Unmarshal(msgBytes, msg)
		if err != nil {
			log.Printf("Unable to unmarshal message: %v\n", err)
			continue
		}

		broadcastCH <- msg
	}
}
