package main

import (
	"fmt"
)

func (n *Network) Start() {
	n.Add(1)
	go n.bully()
}

func (n *Network) ping(nodeID int) bool {
	for _, node := range n.Nodes {
		if node.NodeID == nodeID {
			if node.IsFailed {
				return false
			}
			return true
		}
	}
	return false
}

func (n *Network) election(nodeIndex int) {
	fmt.Printf("Node %v is holding election\n", n.Nodes[nodeIndex].NodeID)
	itr := nodeIndex - 1
	for itr >= 0 {
		// to get the feeling of distribution, I intentionally implemented verbose ping
		OK := n.ping(n.Nodes[itr].NodeID)
		if OK {
			fmt.Printf("Node %v with high priority is up.\n", n.Nodes[itr].NodeID)
			// it's now upto Node[itr]
			n.election(itr)
			return
		}
		itr--
	}

	// if no greater priority node are active
	n.MakeCoordinator(n.Nodes[nodeIndex].NodeID)
}

func (n *Network) bully() {
	defer n.Done()
	totalNodes := len(n.Nodes)
	i := 0

	for {
		n.Lock()
		if !n.Nodes[i].IsFailed {
			fmt.Printf("Node %v is in process", n.Nodes[i].NodeID)

			if n.IsCoordinatorFailed() {
				fmt.Println("Coordinator node failed")
				n.election(i)
			}
		} else {
			fmt.Printf("Node %v is failed. Skipping...", n.Nodes[i].NodeID)
		}
		n.Unlock()
		i = (i + 1) % totalNodes

		n.State()

		var in string
		fmt.Println("Press i for input mode, c for continue...")
		fmt.Scanf("%s", &in)
		switch in {
		case "i":
			n.Controll()
		case "s":
			continue
		}
	}
}

// Controll is used to make up and down nodes to feel the simulation and bully algorithm
func (n *Network) Controll() {
	var nodeID NodeIDType
	var operation int
	fmt.Printf("\nEnter node Id and 0/1 to take that node down/up: ")
	fmt.Scanf("%d %d", &nodeID, &operation)

	for idx := range n.Nodes {
		if n.Nodes[idx].NodeID == nodeID {
			if operation == 0 {
				if !n.Nodes[idx].IsFailed {
					n.Nodes[idx].IsFailed = true
				}
			} else {
				if n.Nodes[idx].IsFailed {
					n.Nodes[idx].IsFailed = false
					n.election(idx)
				}
			}
			return
		}
	}
}

func main() {
	network := CreateNetwork()

	nodes := []Node{
		CreateNode(1, 1),
		CreateNode(2, 2),
		CreateNode(3, 3),
		CreateNode(4, 4),
	}

	for _, node := range nodes {
		network.InsertNode(node)
	}

	network.Start()
	network.Wait()
}
