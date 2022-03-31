package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Client struct {
	ID            int   // server PID
	replies       []int // list of PIDs who have sent votes to this node
	server        *Server
	hasVoted      bool // has this node already voted
	voteTimestamp int64
	voteHolder    int
	votingQueue   []Request // list of pending votes, FIFO
	quorum        int       // number of votes needed to enter critical section
}

var replies = make([]int, 0)

func CreateClient() Client {
	return Client{ID: int(time.Now().UnixNano())}
}

func (c *Client) RegisterServer(server *Server) {
	c.server = server
}

// SendRequest wants to enter critical section
func (c *Client) SendRequest() {
	c.quorum = int(math.Ceil(float64(len(c.server.clients) / 2)))

	// broadcast requests
	fmt.Printf("[%d] Requesting to enter CS.\n", c.ID)
	c.server.PassRequestToAll(Request{source: c.ID, timestamp: time.Now().UnixNano()}, c)

	// wait until have the majority of votes
	for {
		if len(replies) >= c.quorum {
			break
		}

		fmt.Printf("[%d] Only have %d/%d replies.\n", c.ID, len(replies), c.quorum)
		fmt.Printf("Replies: %v\n", replies)

		// random backoff
		rand.Seed(time.Now().UnixNano())
		delay := rand.Intn(5)
		time.Sleep(time.Duration(delay) * time.Second)
	}

	// enter CS
	c.server.Counter++

	// broadcast releases
	fmt.Printf("[%d] CS done, releasing votes.\n", c.ID)
	for _, nodeID := range replies {
		c.server.PassReleaseToOne(Release(c.ID), nodeID)
		fmt.Printf("[%d] Release sent to %d\n", c.ID, nodeID)
	}
}

// ReceiveVote receives vote from other nodes
func (c *Client) ReceiveVote(vote Vote) {
	// add vote to replies
	fmt.Printf("[%d] Received vote from %d.\n", c.ID, vote)
	replies = append(replies, int(vote))
	fmt.Printf("Replies: %v\n", replies)
}

// ReceiveRescind other node rescinded vote
func (c *Client) ReceiveRescind(vote Vote) {
	// check if node has reached quorum, if we above quorum then we are in CS and do nothing
	fmt.Printf("[%d] %d is rescinding vote.\n", c.ID, vote)
	if len(replies) > c.quorum {
		fmt.Printf("[%d] Already at quorum, do nothing.\n", c.ID)
		return
	}

	// not above quorum, remove vote from replies
	for idx, nodeID := range replies {
		if nodeID == int(vote) {
			replies[idx] = replies[len(replies)-1]
			replies = replies[:len(replies)-1]
		}
	}

	// send release
	fmt.Printf("[%d] Sending release to %d\n", c.ID, vote)
	c.server.PassReleaseToOne(Release(c.ID), int(vote))

	// re-request vote
	fmt.Printf("[%d] Sending request to %d\n", c.ID, vote)
	c.server.PassRequestToOne(Request{source: c.ID, timestamp: time.Now().UnixNano()}, int(vote))
}

// ReceiveRequest other node wants to enter critical section
func (c *Client) ReceiveRequest(request Request) {
	fmt.Printf("[%d] Recieved request from %d with timestamp %d\n", c.ID, request.source, request.timestamp)
	if c.hasVoted {
		// if node has voted, check timestamp
		fmt.Printf("[%d] Has already voted.\n", c.ID)
		if request.timestamp < c.voteTimestamp {
			// rescind vote
			fmt.Printf("[%d] Incoming timestamp is smaller, rescinding vote from %d\n", c.ID, c.voteHolder)
			c.server.PassRescindToOne(Vote(c.ID), c.voteHolder)

			// send vote to incoming request
			fmt.Printf("[%d] Sending vote to %d\n", c.ID, request.source)
			c.voteHolder = request.source
			c.voteTimestamp = request.timestamp
			c.server.PassVoteToOne(Vote(c.ID), request.source)
		} else {
			// add to queue
			fmt.Printf("[%d] Added %d to vote queue.\n", c.ID, request.source)
			c.votingQueue = append(c.votingQueue, request)
		}
	} else {
		// else send vote
		fmt.Printf("[%d] Sending vote to %d\n", c.ID, request.source)
		c.hasVoted = true
		c.voteHolder = request.source
		c.voteTimestamp = request.timestamp
		c.server.PassVoteToOne(Vote(c.ID), request.source)
	}
}

// ReceiveRelease node has finished critical section or acknowledged vote rescind
func (c *Client) ReceiveRelease(release Release) {
	fmt.Printf("[%d] Recieved release from %d\n", c.ID, release)
	if len(c.votingQueue) != 0 {
		// check for pending votes in queue, send vote if pending request exists
		c.voteHolder = c.votingQueue[0].source
		c.voteTimestamp = c.votingQueue[0].timestamp
		c.server.PassVoteToOne(Vote(c.ID), c.votingQueue[0].source)
		if len(c.votingQueue) <= 1 {
			c.votingQueue = make([]Request, 0)
		} else {
			c.votingQueue = c.votingQueue[1:]
		}
		fmt.Printf("[%d] Sending vote to %d\n", c.ID, c.voteHolder)
	} else {
		// else change hasVoted to false
		c.voteHolder = 0
		c.voteTimestamp = 0
		c.hasVoted = false
		fmt.Printf("[%d] No one in the voting queue.\n", c.ID)
	}
}
