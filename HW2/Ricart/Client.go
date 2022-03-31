package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Client struct {
	ID           int
	replies      []int
	server       *Server
	requests     []Message
	selfRequest  Message
	isExecuting  bool
	isRequesting bool
}

func CreateClient() Client {
	return Client{ID: int(time.Now().UnixNano())}
}

func (c *Client) RegisterServer(server *Server) {
	c.server = server
}

func (c *Client) EnterCS() {
	// send timestamped message to all other clients
	fmt.Printf("[%d] Trying to enter CS\n", c.ID)
	c.isRequesting = true
	c.selfRequest = Message{source: c.ID, timestamp: time.Now().UnixNano()}
	c.server.PassMessageToAll(c.selfRequest, c)

	for {
		if len(c.replies) != len(c.server.clients)-1 {
			break
		}
		fmt.Printf("[%d] Only have %d/%d replies.\n", c.ID, c.replies, len(c.server.clients)-1)
		fmt.Println(c.replies)

		// random backoff
		rand.Seed(time.Now().UnixNano())
		delay := rand.Intn(2)
		time.Sleep(time.Duration(delay) * time.Second)
	}
	fmt.Printf("[%d] Got all replies, entering CS.\n", c.ID)
	c.isExecuting = true

	// critical section
	c.server.Counter++

	c.isExecuting = false
	c.isRequesting = false

	fmt.Printf("[%d] CS complete, sending replies to all deferred requests.\n", c.ID)
	for _, req := range c.requests {
		c.server.PassReleaseToOne(Release(c.ID), req.source)
		fmt.Printf("[%d] Replied to %d\n", c.ID, req.source)
	}
	c.replies = make([]int, 0)
}

func (c *Client) ReceiveRelease(message Release) {
	fmt.Printf("[%d] Recieved reply from %d\n", c.ID, message)
	c.replies = append(c.replies, int(message))
}

func (c *Client) ReceiveCS(message Message, source Client) {
	fmt.Printf("[%d] Recieved request from %d\n", c.ID, source.ID)
	for c.isExecuting {
	}

	fmt.Printf("[%d] Is not executing. \n", c.ID)
	if !c.isRequesting {
		fmt.Printf("[%d] Is not requesting, sending reply. \n", c.ID)
		c.server.PassReleaseToOne(Release(c.ID), source.ID)
	} else if c.isRequesting && message.timestamp < c.selfRequest.timestamp {
		fmt.Printf("[%d] Is requesting, but incoming timestamp is lower. \n", c.ID)
		c.server.PassReleaseToOne(Release(c.ID), source.ID)
	} else {
		fmt.Printf("[%d] Deferred response. \n", c.ID)
		c.requests = append(c.requests, message)
	}
}
