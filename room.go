package main

import "fmt"

type Room struct {
	name    string
	clients []*Client
}

func (r *Room) AddClient(client *Client) bool {
	fmt.Println("ROOM CLIENTS", r.clients)
	for idx, rClient := range r.clients {
		if rClient.id == client.id {
			fmt.Println("UPDATING")
			r.clients[idx] = client
			return true
		}
	}

	fmt.Println("new client", client)
	r.clients = append(r.clients, client)
	return true
}
