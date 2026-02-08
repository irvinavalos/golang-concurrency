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
}

func NewServer() *Server {
	return &Server{
		Clients:       map[string]*Client{},
		mu:            new(sync.RWMutex),
		joinServerCh:  make(chan *Client),
		leaveServerCh: make(chan *Client),
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
}

func (s *Server) joinServer(c *Client) {
	s.Clients[c.ID] = c
	log.Printf("ClientID: %s joined\n", c.ID)
}

func (s *Server) leaveServer(c *Client) {
	delete(s.Clients, c.ID)
	log.Printf("ClientID: %s left\n", c.ID)
}

func (s *Server) AcceptLoop() {
	for {
		select {
		case c := <-s.joinServerCh:
			s.joinServer(c)
		case c := <-s.leaveServerCh:
			s.leaveServer(c)
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
