package main

import "fmt"

type Node struct {
	ID      int
	Manager Manager
	Pages   []NodePage
}

func (n *Node) Read(pageID int) string {
	reqPage := NodePage{PageID: -1, Access: NONE}

	for _, page := range n.Pages {
		if page.PageID == pageID {
			reqPage = page
		}
	}

	if reqPage.Access == NONE {
		// read fault, contact central manager
		fmt.Printf("[%v] Page not found locally, requesting from central server\n", n.ID)
		return n.Manager.ReadRequest(pageID, n).Data
	} else {
		return reqPage.Data
	}
}

func (n *Node) ReadRequest(pageID int) NodePage {
	reqPage := NodePage{PageID: -1, Access: NONE}

	for _, page := range n.Pages {
		if page.PageID == pageID {
			reqPage = page
		}
	}

	if reqPage.PageID == -1 || reqPage.Access == NONE {
		panic("Error reading page from owner")
	}

	// change access type
	reqPage.Access = READ

	return reqPage
}

func (n *Node) Invalidate(pageID int) bool {
	isSuccess := false
	for _, page := range n.Pages {
		if page.PageID == pageID {
			page.Access = NONE
			isSuccess = true
			break
		}
	}

	return isSuccess
}

func (n *Node) Write(pageID int, data string) bool {
	reqPage := NodePage{PageID: pageID, Access: NONE, Data: data}

	return n.Manager.WriteRequest(&reqPage, n)
}

func (n *Node) WriteRequest(nodePage NodePage, newEntry bool) bool {
	if !newEntry {
		removeIndex := -1
		// remove existing entry
		for idx, page := range n.Pages {
			if page.PageID == nodePage.PageID {
				removeIndex = idx
			}
		}

		if removeIndex == -1 {
			panic("New entry specified but entry not found")
		}

		removeWithIndex(n.Pages, removeIndex)
	}

	nodePage.Access = WRITE
	n.Pages = append(n.Pages, nodePage)

	return true
}

func removeWithIndex(array []NodePage, i int) []NodePage {
	array[i] = array[len(array)-1]
	return array[:len(array)-1]
}
