package main

const (
	NONE  = 0
	READ  = 1
	WRITE = 2
)

type PageManager struct {
	ID      int
	CopySet []*Node
	Owner   *Node
}

type NodePage struct {
	PageID int
	Access int
	Data   string
}
