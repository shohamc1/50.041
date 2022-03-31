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

	var clients []Client
	for i := 1; i <= 15; i++ {
		tempClient := CreateClient()
		tempClient.RegisterServer(&server)
		server.AddClient(&tempClient)
		clients = append(clients, tempClient)
	}

	for i, _ := range clients {
		waitGroupFunction(clients[i].EnterCS, &wg)
	}

	wg.Wait()
	fmt.Println(server.Counter)

	//waitGroupFunction(client1.SendMessage, &wg)
	//waitGroupFunction(client1.LocalEvent, &wg)
	//waitGroupFunction(client2.LocalEvent, &wg)
	//waitGroupFunction(client2.LocalEvent, &wg)
	//waitGroupFunction(client2.SendMessage, &wg)
	//
	//wg.Wait()
	//
	//fmt.Println(client1.GetTotalOrder())
	//fmt.Println(client2.GetTotalOrder())
}
