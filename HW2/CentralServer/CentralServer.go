package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	server := Server{}
	var wg sync.WaitGroup

	numClients, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return
	}

	var clients []Client
	for i := 1; i <= numClients; i++ {
		tempClient := CreateClient()
		tempClient.RegisterServer(&server)
		server.AddClient(&tempClient)
		clients = append(clients, tempClient)
	}

	timer := time.Now().UnixNano()
	for i, _ := range clients {
		wg.Add(1)

		go func() {
			defer wg.Done()
			clients[i].SendRequest()
		}()
	}

	wg.Wait()
	fmt.Println((time.Now().UnixNano() - timer) / int64(time.Microsecond))
	fmt.Println(server.Counter)
}
