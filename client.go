package main

import (
	"crypto/rand"
	"sync"

	"github.com/gorilla/websocket"
)

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
