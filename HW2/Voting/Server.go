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

func (s Server) PassRequestToAll(message Request, source *Client) {
	// pass message to all other clients
	for _, v := range s.clients {
		if v.ID != source.ID {
			v.ReceiveRequest(message)
		}
	}
}

func (s Server) PassRequestToOne(message Request, recipient int) {
	// send vote to client
	for _, v := range s.clients {
		if v.ID == recipient {
			v.ReceiveRequest(message)
		}
	}
}

func (s Server) PassVoteToOne(message Vote, recipient int) {
	// send vote to client
	for _, v := range s.clients {
		if v.ID == recipient {
			v.ReceiveVote(message)
		}
	}
}

func (s Server) PassRescindToOne(message Vote, recipient int) {
	// send vote to client
	for _, v := range s.clients {
		if v.ID == recipient {
			v.ReceiveRescind(message)
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
