package main

type Room struct {
	name    string
	clients []*Client
}

func (r *Room) AddClient(client *Client) bool {
	for idx, rClient := range r.clients {
		if rClient.id == client.id {
			r.clients[idx] = client
			return true
		}
	}

	r.clients = append(r.clients, client)
	return true
}

func (r *Room) GetClient(client *Client) *Client {
	for _, rClient := range r.clients {
		if rClient.id == client.id {
			return rClient
		}
	}
	return nil
}
