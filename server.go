package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	PORT = ":8000"
)

type Server struct {
	Clients       map[string]*Client
	mu            *sync.RWMutex
	joinServerCh  chan *Client
	leaveServerCh chan *Client
	broadcastCh   chan *RequestMessage
}

func NewServer() *Server {
	return &Server{
		Clients:       map[string]*Client{},
		mu:            new(sync.RWMutex),
		joinServerCh:  make(chan *Client),
		leaveServerCh: make(chan *Client),
		broadcastCh:   make(chan *RequestMessage),
	}
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Errror on HTTP connection upgrade: %v\n", err)
		return
	}

	client := NewClient(conn)
	s.joinServerCh <- client

	go client.readMessageLoop(s.leaveServerCh, s.broadcastCh)
}

func (s *Server) joinServer(c *Client) {
	s.Clients[c.ID] = c
	log.Printf("ClientID: %s joined\n", c.ID)
}

func (s *Server) leaveServer(c *Client) {
	delete(s.Clients, c.ID)
	log.Printf("ClientID: %s left\n", c.ID)
}

func (s *Server) broadcast(msg *RequestMessage) {
	clients := []*Client{}

	s.mu.Lock()
	for _, c := range s.Clients {
		if c.ID != msg.client.ID {
			clients = append(clients, c)
		}
	}
	s.mu.Unlock()

	response := NewResponseMessage(msg)

	for _, c := range clients {
		err := c.conn.WriteJSON(response)
		if err != nil {
			log.Printf("Error sending message to ClientID: %s", c.ID)
			continue
		}
	}

	log.Println("Broadcast was sent...")
}

func (s *Server) AcceptLoop() {
	for {
		select {
		case client := <-s.joinServerCh:
			s.joinServer(client)
		case client := <-s.leaveServerCh:
			s.leaveServer(client)
		case msg := <-s.broadcastCh:
			go s.broadcast(msg)
		}
	}
}

func startServer() {
	server := NewServer()

	go server.AcceptLoop()

	http.HandleFunc("/", server.handleWS)

	log.Println("Starting server")
	log.Fatal(http.ListenAndServe(PORT, nil))
}
