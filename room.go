package main

type Room struct {
	name    string
	clients []Client
}

func (r *Room) AddClient(client Client) (bool, error) {
	r.clients = append(r.clients, client)

	return true, nil
}
