package main

import (
	"fmt"
)

func main() {
	server := Server{}
	server.AddClient(100)
	server.AddClient(69)
	server.AddClient(420)
	fmt.Println(server.GetClients())

	server.RemoveClient(69)
	fmt.Println(server.GetClients())

	server.RemoveClient(100)
	fmt.Println(server.GetClients())
}
