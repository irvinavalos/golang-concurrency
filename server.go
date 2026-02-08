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
	Clients map[string]*Client
	mu      *sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		Clients: map[string]*Client{},
		mu:      new(sync.RWMutex),
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
	s.Clients[client.ID] = client
}

func startServer() {
	server := NewServer()
	http.HandleFunc("/", server.handleWS)
    log.Println("Starting server")
	log.Fatal(http.ListenAndServe(PORT, nil))
}
