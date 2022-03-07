package main

import "time"

type Event struct {
	source int
	clock  int
}

type Client struct {
	ID         int
	server     *Server
	clock      int
	totalOrder []Event
}

func CreateClient() Client {
	return Client{ID: int(time.Now().UnixNano())}
}

func (c Client) GetTotalOrder() []Event {
	return c.totalOrder
}

func (c *Client) RegisterServer(server *Server) {
	c.server = server
}

func (c *Client) LocalEvent() {
	c.clock++
	c.totalOrder = append(c.totalOrder, Event{clock: c.clock, source: c.ID})
}

func (c *Client) SendMessage() {
	c.clock++
	c.server.PassMessage(c.clock, c)
	c.totalOrder = append(c.totalOrder, Event{clock: c.clock, source: c.ID})
}

func (c *Client) ReceiveMessage(message int, source Client) {
	if message > c.clock {
		c.clock = message
		c.totalOrder = append(c.totalOrder, Event{clock: message, source: source.ID})
		c.totalOrder = append(c.totalOrder, Event{clock: message, source: c.ID})
	} else if message == c.clock { // use PID to break tie
		if message*source.ID > message*c.ID {
			c.totalOrder = append(c.totalOrder, Event{clock: message, source: source.ID})
			c.totalOrder = append(c.totalOrder, Event{clock: message, source: c.ID})
		} else {
			c.totalOrder = append(c.totalOrder, Event{clock: message, source: c.ID})
			c.totalOrder = append(c.totalOrder, Event{clock: message, source: source.ID})
		}
	} else {
		c.clock++
		c.totalOrder = append(c.totalOrder, Event{clock: message, source: c.ID})
		c.totalOrder = append(c.totalOrder, Event{clock: message, source: source.ID})
	}
}
