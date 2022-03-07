package main

type Server struct {
	clients []int
}

func (s *Server) AddClient(clientId int) {
	s.clients = append(s.clients, clientId)
}

func (s *Server) RemoveClient(clientId int) {
	for i, v := range s.clients {
		if v == clientId {
			s.clients[i] = s.clients[len(s.clients)-1]
			s.clients = s.clients[:len(s.clients)-1]
		}
	}
}

func (s Server) GetClients() []int {
	return s.clients
}
