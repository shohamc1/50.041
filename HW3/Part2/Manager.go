package main

import (
	"fmt"
)

type Manager struct {
	ID       int
	Pages    []PageManager
	Replicas []*Manager
	Nodes    []*Node
	up       bool
}

var localPages = make([]PageManager, 1)

func (m *Manager) UpdateNodes(newNodes []*Node) {
	m.Nodes = newNodes
}

func (m *Manager) DeclareReplica(managerReplica []*Manager) {
	if !m.up {
		return
	}

	m.Replicas = append(m.Replicas, managerReplica...)
}

func (m *Manager) UpdateReplicaSet(newReplicas []*Manager) bool {
	// update local copy, send ACK
	m.Replicas = newReplicas
	return true
}

func (m *Manager) GetPagesUpdate(pages []PageManager) {
	if !m.up {
		return
	}

	fmt.Printf("[Manager %d] Got page update %v\n", m.ID, pages)
	localPages = pages
}

func (m *Manager) UpdateReplicaPages() {
	if !m.up {
		return
	}

	for _, replica := range m.Replicas {
		if replica.ID != m.ID {
			replica.GetPagesUpdate(localPages)
		}
	}
}

func (m *Manager) TakeDown() {
	found := false
	m.up = false

	// initiate election
	for _, replica := range m.Replicas {
		if replica.ID != m.ID && replica.up {
			replica.InitiateElection(m)
			found = true
			break
		}
	}

	if !found {
		panic("No more replicas available!")
	}
}

func (m *Manager) Up() {
	m.up = true
	currentMaster := m.Nodes[0].Manager
	currentMaster.AddReplica(m)
}

func (m *Manager) AddReplica(newManager *Manager) {
	m.Replicas = append(m.Replicas, newManager)
	fmt.Printf("[Manager %d] Adding %v to replica set. New replica set: %v\n", m.ID, newManager.ID, m.Replicas)

	for _, replica := range m.Replicas {
		if replica.ID != m.ID {
			replica.UpdateReplicaSet(m.Replicas)
		}
	}
}

func (m *Manager) InitiateElection(initiator *Manager) {
	fmt.Printf("[Manager %d] Initiating election\n", m.ID)

	removeIndex := -1

	// remove existing entry
	for idx, replica := range m.Replicas {
		if replica == initiator {
			removeIndex = idx
		}
	}
	if removeIndex == -1 {
		panic("Election initiator not a replica!")
	}

	m.Replicas = removeWithManagerIndex(m.Replicas, removeIndex)
	fmt.Printf("[Manager %d] New replica set: %v\n", m.ID, m.Replicas)

	// update other nodes
	allNodesOk := true
	for _, replica := range m.Replicas {
		if replica.ID != m.ID {
			allNodesOk = replica.UpdateReplicaSet(m.Replicas)
		}
	}

	if !allNodesOk {
		panic("Error updating replica sets")
	}

	fmt.Printf("[Manager %d] Election succeeded, updating nodes\n", m.ID)
	for _, node := range m.Nodes {
		node.UpdateManager(m)
	}
}

func (m *Manager) ReadRequest(pageID int, senderNode *Node) NodePage {
	fmt.Printf("[Manager] Current set: %v\n", localPages)
	reqPage := &PageManager{ID: -1}

	for _, page := range localPages {
		if page.ID == pageID {
			reqPage = &page
		}
	}

	if reqPage.ID == -1 {
		// not found in central server, throw error
		fmt.Println("[Manager] Page does not exist")
		return NodePage{PageID: -1, Access: NONE}
	} else {
		// add to copy set
		fmt.Printf("[Manager] Added %v to CopySet\n", senderNode.ID)
		reqPage.CopySet = append(reqPage.CopySet, senderNode)
		m.UpdateReplicaPages()

		// forward read request
		fmt.Printf("[Manager] Forwarding read request to %v\n", reqPage.Owner.ID)
		return reqPage.Owner.ReadRequest(pageID)
	}
}

func (m *Manager) WriteRequest(nodePage *NodePage, senderNode *Node) bool {
	isSuccess := false
	reqPage := &PageManager{ID: -1}

	for _, page := range localPages {
		if page.ID == nodePage.PageID {
			reqPage = &page
		}
	}

	if reqPage.ID == -1 {
		// new page, add record to page manager
		fmt.Printf("[Manager] Page not exist, adding to local set and sending write request back to sender %v\n", senderNode.ID)
		newPage := PageManager{ID: nodePage.PageID, Owner: senderNode, CopySet: []*Node{}}
		localPages = append(localPages, newPage)
		fmt.Printf("[Manager] New set: %v\n", localPages)
		m.UpdateReplicaPages()
		isSuccess = senderNode.WriteRequest(*nodePage, true)
	} else {
		// invalidate copies
		fmt.Println("[Manager] Invalidating copy sets")
		for _, node := range reqPage.CopySet {
			node.Invalidate(nodePage.PageID)
		}

		reqPage.CopySet = make([]*Node, 0)
		m.UpdateReplicaPages()

		// send write to owner
		fmt.Printf("[Manager] Send write to %v", reqPage.Owner)
		isSuccess = reqPage.Owner.WriteRequest(*nodePage, false)
	}
	return isSuccess
}

func removeWithManagerIndex(array []*Manager, i int) []*Manager {
	array[i] = array[len(array)-1]
	return array[:len(array)-1]
}
