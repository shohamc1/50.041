package main

import (
	"fmt"
	"time"
)

type Client struct {
	ID     int // server PID
	server *Server
}

var hasGrant = false

var replies = make([]int, 0)

func CreateClient() Client {
	return Client{ID: int(time.Now().UnixNano())}
}

func (c *Client) RegisterServer(server *Server) {
	c.server = server
}

// SendRequest wants to enter critical section
func (c *Client) SendRequest() {
	fmt.Printf("[%d] Trying to enter CS.\n", c.ID)

	// send request to server
	fmt.Printf("[%d] Sent request to server.\n", c.ID)
	c.server.AddRequest(*c)

	// wait until receive grant
	for {
		if hasGrant {
			break
		}
	}
	fmt.Printf("[%d] Entering CS.\n", c.ID)
	// execute critical section
	c.server.Counter++

	// send release
	fmt.Printf("[%d] Sending release.\n", c.ID)
	c.server.AcceptRelease()
}

// ReceiveGrant receives vote from other nodes
func (c *Client) ReceiveGrant() {
	fmt.Printf("[%d] Recieved grant.\n", c.ID)
	hasGrant = true
}
