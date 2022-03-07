package main

import (
	"fmt"
	"sync"
)

func waitGroupFunction(fn func(), wg *sync.WaitGroup) {
	defer wg.Done()
	wg.Add(1)
	fn()
}

func main() {
	server := Server{}
	var wg sync.WaitGroup

	client1 := CreateClient()
	client1.RegisterServer(&server)
	server.AddClient(&client1)

	client2 := CreateClient()
	client2.RegisterServer(&server)
	server.AddClient(&client2)

	waitGroupFunction(client1.SendMessage, &wg)
	waitGroupFunction(client1.LocalEvent, &wg)
	waitGroupFunction(client2.LocalEvent, &wg)
	waitGroupFunction(client2.LocalEvent, &wg)
	waitGroupFunction(client2.SendMessage, &wg)

	wg.Wait()

	fmt.Printf("Client 1: %v\n", client1.ID)
	fmt.Printf("Client 2: %v\n", client2.ID)

	fmt.Println("CLIENT 1")
	totalOrder := client1.GetTotalOrder()
	for _, v := range totalOrder {
		fmt.Println(v.source)
		fmt.Println(v.clock)
		fmt.Println()
	}

	fmt.Println("CLIENT 2")
	totalOrder = client2.GetTotalOrder()
	for _, v := range totalOrder {
		fmt.Println(v.source)
		fmt.Println(v.clock)
		fmt.Println()
	}
}
