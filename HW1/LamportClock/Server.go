package main

import (
	"math/rand"
	"time"
)

type Server struct {
	clients []*Client
}

func (s *Server) AddClient(client *Client) {
	s.clients = append(s.clients, client)
}

func (s *Server) RemoveClient(client *Client) {
	for i, v := range s.clients {
		if v == client {
			s.clients[i] = s.clients[len(s.clients)-1]
			s.clients = s.clients[:len(s.clients)-1]
		}
	}
}

func (s Server) GetClients() []Client {
	var resolvedClients []Client
	for _, v := range s.clients {
		resolvedClients = append(resolvedClients, *v)
	}

	return resolvedClients
}

func (s Server) PassMessage(message int, source *Client) {
	// random time backoff

	rand.Seed(time.Now().UnixNano())
	delay := rand.Intn(3)
	time.Sleep(time.Duration(delay) * time.Second)

	// pass message to all other clients
	for _, v := range s.clients {
		if v != source {
			v.ReceiveMessage(message, *source)
		}
	}
}
