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

func (r *Room) GetClientsNames() []string {
	var clientsNames []string

	for _, client := range r.clients {
		clientsNames = append(clientsNames, client.name)
	}
	return clientsNames
}

func (r *Room) GetClientByName(name string) *Client {
	for _, rClient := range r.clients {
		if rClient.name == name {
			return rClient
		}
	}
	return nil
}

func (r *Room) RemoveClient(client *Client) bool {
	for idx, rClient := range r.clients {
		if rClient.id == client.id {
			r.clients[idx] = r.clients[len(r.clients)-1]
			r.clients = r.clients[:len(r.clients)-1]
			return true
		}
	}

	return false
}
