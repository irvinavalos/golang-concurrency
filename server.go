package main

import "sync"

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
