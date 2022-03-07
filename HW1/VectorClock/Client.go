package main

import "time"

type Event struct {
	source int
	clock  map[int]int
}

type Client struct {
	ID          int
	server      *Server
	vectorClock map[int]int
	totalOrder  []Event
}

func CreateClient() Client {
	return Client{ID: int(time.Now().UnixNano()), vectorClock: make(map[int]int)}
}

func (c Client) GetTotalOrder() []Event {
	return c.totalOrder
}

func (c *Client) RegisterServer(server *Server) {
	c.server = server
}

func (c *Client) LocalEvent() {
	c.vectorClock[c.ID]++
	c.totalOrder = append(c.totalOrder, Event{clock: c.vectorClock, source: c.ID})
}

func (c *Client) SendMessage() {
	c.vectorClock[c.ID]++
	c.server.PassMessage(c.vectorClock, c)
	c.totalOrder = append(c.totalOrder, Event{clock: c.vectorClock, source: c.ID})
}

func (c *Client) ReceiveMessage(message map[int]int, source Client) {
	incomingBefore := true // is the incoming message before the current event
	for PID, mClock := range message {
		if mClock > c.vectorClock[PID] {
			c.vectorClock[PID] = mClock
		} else if mClock == c.vectorClock[PID] {
			if PID > c.ID {
				incomingBefore = false
			}
		} else {
		}
	}

	if incomingBefore {
		c.totalOrder = append(c.totalOrder, Event{clock: message, source: source.ID})
		c.totalOrder = append(c.totalOrder, Event{clock: message, source: c.ID})
	} else {
		c.totalOrder = append(c.totalOrder, Event{clock: message, source: c.ID})
		c.totalOrder = append(c.totalOrder, Event{clock: message, source: source.ID})
	}
}
