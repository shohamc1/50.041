package main

type Server struct {
	clients    []*Client
	grantOrder []Client
	hasGranted bool
	Counter    int
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

func (s *Server) AddRequest(requester Client) {
	// incoming request to enter CS
	if s.hasGranted {
		// some node already in CS, add to queue
		s.grantOrder = append(s.grantOrder, requester)
	} else {
		// no node in CS, grant permission
		requester.ReceiveGrant()
	}
}

func (s *Server) AcceptRelease() {
	if len(s.grantOrder) > 0 {
		// grant permission to next node in order
		s.grantOrder[0].ReceiveGrant()
		if len(s.grantOrder) == 0 {
			s.grantOrder = make([]Client, 0)
		} else {
			s.grantOrder = s.grantOrder[1:]
		}
	} else {
		s.hasGranted = false
	}
}
