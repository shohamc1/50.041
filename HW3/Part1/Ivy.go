package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	manager := Manager{Pages: []PageManager{}}

	var nodes []Node
	var tempNode Node

	for i := 0; i < 10; i++ {
		tempNode = Node{ID: int(time.Now().UnixNano()), Manager: manager}
		nodes = append(nodes, tempNode)
	}

	var selectedNode int
	var operation string
	var key int
	var data string

	fmt.Println("Input format: [node ID] [read/write] [key] [value]")

	for {
		selectedNode, key = -1, -1
		operation, data = "", ""
		fmt.Scanf("%d %s %d %s", &selectedNode, &operation, &key, &data)

		if selectedNode != -1 && key != -1 {
			switch strings.ToLower(operation) {
			case "read":
				fmt.Println(nodes[selectedNode].Read(key))
			case "write":
				if data == "" {
					fmt.Println("No data provided")
				} else {
					fmt.Println(nodes[selectedNode].Write(key, data))
				}
			default:
				fmt.Println("Invalid operation, please use \"read\" or \"write\"")
			}
		} else {
			fmt.Println("No node ID or key provided.")
		}
	}
}

func PerformanceEval() {
	manager := Manager{Pages: []PageManager{}}

	var nodes []Node
	var tempNode Node

	for i := 0; i < 10; i++ {
		tempNode = Node{ID: int(time.Now().UnixNano()), Manager: manager}
		nodes = append(nodes, tempNode)
	}

	nodes[0].Write(1, "Hello")
	timer := time.Now().UnixNano()

	for _, node := range nodes {
		node.Read(1)
	}

	fmt.Printf("Time to read: %v", (time.Now().UnixNano()-timer)/int64(time.Microsecond))
}
