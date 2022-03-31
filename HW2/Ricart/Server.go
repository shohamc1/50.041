package main

type Server struct {
	clients []*Client
	Counter int
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

func (s Server) PassMessageToAll(message Message, source *Client) {
	// pass message to all other clients
	for _, v := range s.clients {
		if v.ID != source.ID {
			v.ReceiveCS(message, *source)
		}
	}
}

func (s Server) PassReleaseToOne(message Release, recipient int) {
	// pass message to other client
	for _, v := range s.clients {
		if v.ID == recipient {
			v.ReceiveRelease(message)
		}
	}
}
