package main

import "fmt"

type Manager struct {
	Pages []PageManager
}

var localPages = make([]PageManager, 1)

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
		isSuccess = senderNode.WriteRequest(*nodePage, true)
	} else {
		// invalidate copies
		fmt.Println("[Manager] Invalidating copy sets")
		for _, node := range reqPage.CopySet {
			node.Invalidate(nodePage.PageID)
		}

		// send write to owner
		fmt.Printf("[Manager] Send write to %v", reqPage.Owner)
		isSuccess = reqPage.Owner.WriteRequest(*nodePage, false)
	}
	return isSuccess
}
