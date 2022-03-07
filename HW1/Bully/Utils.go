package main

import (
	"fmt"
	"sync"
)

type NodeIDType = int

type PriorityType = int

type Node struct {
	NodeID            NodeIDType
	CoordinatorNodeID NodeIDType
	Priority          PriorityType
	IsCoordinator     bool
	IsFailed          bool
}

func CreateNode(id NodeIDType, priority int) Node {
	return Node{
		id,
		-1,
		priority,
		false,
		false,
	}
}

type Network struct {
	sync.WaitGroup
	sync.Mutex
	Nodes []Node
	Stop  chan bool
}

func CreateNetwork() Network {
	return Network{
		Nodes: make([]Node, 0),
		Stop:  make(chan bool),
	}
}

func (n *Network) State() {
	failedNodes := make([]NodeIDType, 0)
	currentActiveNode := -1
	for _, node := range n.Nodes {
		if node.IsFailed {
			failedNodes = append(failedNodes, node.NodeID)
			continue
		}
		if node.IsCoordinator {
			currentActiveNode = node.NodeID
		}
	}

	fmt.Printf("Failed nodes: %v\nCoordinator: %v\n", failedNodes, currentActiveNode)
}

func (n *Network) MakeCoordinator(nodeID NodeIDType) {
	// make one of node coordinator
	fmt.Printf("Making Node %v coordinator\n", nodeID)
	for idx := range n.Nodes {
		if n.Nodes[idx].IsFailed {
			continue
		}
		n.Nodes[idx].CoordinatorNodeID = nodeID
		if n.Nodes[idx].NodeID == nodeID {
			n.Nodes[idx].IsCoordinator = true
		} else if n.Nodes[idx].IsCoordinator {
			n.Nodes[idx].IsCoordinator = false
		}
	}
}

func (n *Network) InsertNode(node Node) {
	// insertion sort based on node priority
	totalNodes := len(n.Nodes)
	if totalNodes == 0 {
		n.Nodes = []Node{node}
		n.MakeCoordinator(node.NodeID)
		return
	}

	n.Nodes = append(n.Nodes, node)
	itr := len(n.Nodes) - 2
	for itr >= 0 && n.Nodes[itr].Priority < node.Priority {
		n.Nodes[itr+1] = n.Nodes[itr]
		itr--
	}
	n.Nodes[itr+1] = node
	n.MakeCoordinator(n.Nodes[0].NodeID)
}

func (n *Network) IsCoordinatorFailed() bool {
	// health check
	for _, node := range n.Nodes {
		if node.IsCoordinator && node.IsFailed {
			return true
		}
	}
	return false
}
