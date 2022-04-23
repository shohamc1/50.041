package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	var managers []*Manager

	for i := 0; i < 3; i++ {
		tempManager := Manager{ID: i + 1, Pages: []PageManager{}, up: true}
		managers = append(managers, &tempManager)
	}

	var nodes []*Node

	for i := 0; i < 10; i++ {
		tempNode := Node{ID: int(time.Now().UnixNano()), Manager: managers[0]}
		nodes = append(nodes, &tempNode)
	}

	for _, manager := range managers {
		manager.DeclareReplica(managers)
		manager.UpdateNodes(nodes)
	}

	var selectedNode int
	var operation string
	var key int
	var data string

	fmt.Println("Default central server ID is 0\nInput format: [node ID] [read/write] [key] [value]\nOR [manager ID] [up/down]")

	for {
		selectedNode, key = -1, -1
		operation, data = "", ""
		fmt.Scanf("%d %s %d %s", &selectedNode, &operation, &key, &data)

		if selectedNode != -1 {
			switch strings.ToLower(operation) {
			case "read":
				fmt.Println(nodes[selectedNode].Read(key))
			case "write":
				if data == "" {
					fmt.Println("No data provided")
				} else {
					fmt.Println(nodes[selectedNode].Write(key, data))
				}
			case "up":
				managers[selectedNode].Up()
			case "down":
				managers[selectedNode].TakeDown()
			default:
				fmt.Println("Invalid operation, please use \"read\" or \"write\"")
			}
		} else {
			fmt.Println("No node ID or key provided.")
		}
	}
}
