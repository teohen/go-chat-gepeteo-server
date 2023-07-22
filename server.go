package main

import (
	"fmt"
	"io"
	"strings"
)

type Server struct {
	rooms   []*Room
	clients []*Client
}

var server *Server

func NewServer() *Server {
	if server == nil {
		server = &Server{
			rooms: []*Room{},
		}
	}

	return server
}

func (s *Server) readLoop(client Client) {
	buf := make([]byte, 1024)

	for {
		n, err := client.ws.Read(buf)

		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error: ", err)
			continue
		}
		msg := buf[:n]
		generalRoom := s.GetRoomByName("general")
		sClient := generalRoom.GetClient(&client)
		s.handleCommands(sClient, string(msg))
	}
}

func (s *Server) GetRoomByName(name string) *Room {
	for _, sRoom := range s.rooms {
		if sRoom.name == name {
			return sRoom
		}
	}
	return nil
}

func (s *Server) addRoom(room *Room) {
	s.rooms = append(s.rooms, room)
}

func validateBeforeMsg(client *Client, msg []string) bool {
	if client.name == "" {
		if len(msg) < 2 || msg[0] != "client" || msg[1] != "set_name" {
			sendMessageTo(client, "Welcome! To start chatting, please register your name handle with the command 'client set_name YOUR_NICKNAME'")
			return false
		}
	}
	return true
}

// TODO: MAYBE CREATE A SEPARATED COMMANDS MODULE
func (s *Server) handleCommands(client *Client, msg string) {
	clientMsg := strings.Fields(msg)

	if ok := validateBeforeMsg(client, clientMsg); ok != true {
		return
	}

	if clientMsg[0] == "server" {
		s.serverCommands(client, clientMsg)
	} else if clientMsg[0] == "client" {
		s.clientCommands(client, clientMsg)
	} else {
		sendMessageTo(client, "unknown command. Send 'help' for a list of commands")
		return
	}
}

func (s *Server) serverCommands(client *Client, msg []string) {
	generalRoom := s.GetRoomByName("general")

	if msg[1] == "ls_clients" && len(msg) < 3 {
		listClients := "Clients:\n"

		names := generalRoom.GetClientsNames()

		for _, name := range names {
			listClients += name + "\n"
		}

		sendMessageTo(client, listClients)
		return
	} else if msg[1] == "ls_clients" && len(msg) > 2 {
		roomName := msg[2]
		room := s.GetRoomByName(roomName)

		if room == nil {
			sendMessageTo(client, "room does not exist. See list of rooms with 'server ls_rooms' or create one with 'server cr_room ROOM_NAME")
			return
		}

		names := room.GetClientsNames()

		listClientsNames := fmt.Sprintf("Clients in room %s: \n", roomName)

		for _, name := range names {
			listClientsNames += name + "\n"
		}

		sendMessageTo(client, listClientsNames)

	} else if msg[1] == "ls_rooms" {
		listRooms := "Rooms:\n"

		for _, sRoom := range s.rooms {
			listRooms += sRoom.name + "\n"
		}

		sendMessageTo(client, listRooms)
		return

	} else if msg[1] == "cr_room" {
		if len(msg) < 3 {
			sendMessageTo(client, fmt.Sprintf("room name is required! send 'server cr_room ROOM_NAME'"))
			return
		}

		newRoomName := msg[2]
		existentRoom := s.GetRoomByName(newRoomName)

		if existentRoom != nil {
			sendMessageTo(client, "room name already exists. Pick another one!")
			return
		}

		newRoom := Room{
			name: newRoomName,
		}

		s.addRoom(&newRoom)
		newRoom.AddClient(client)
		sendMessageTo(client, fmt.Sprintf("Room %s created and client %s has entered it", newRoomName, client.name))
	} else {
		sendMessageTo(client, "unknown command. Send 'help' for a list of commands")
		return
	}
}

func (s *Server) clientCommands(client *Client, msg []string) {
	if msg[1] == "set_name" {
		if len(msg) < 3 {
			sendMessageTo(client, "client name is required! send 'client set_name CLIENT_NAME'")
			return
		}

		generalRoom := s.GetRoomByName("general")
		existClient := generalRoom.GetClientByName(msg[2])

		if existClient != nil {
			sendMessageTo(client, "client name already taken. Pick another one!")
			return
		}

		client.setName(msg[2])
		sendMessageTo(client, fmt.Sprintf("name set: %s", msg[2]))
	} else if msg[1] == "join" {
		if len(msg) < 3 {
			sendMessageTo(client, "room name is required for the join command! send 'client join ROOM_NAME'")
			return
		}

		room := s.GetRoomByName(msg[2])
		if room == nil {
			sendMessageTo(client, "room does not exist. See list of rooms with 'server ls_rooms' or create one with 'server cr_room ROOM_NAME")
			return
		}

		room.AddClient(client)
		sendMessageTo(client, fmt.Sprintf("Client %s has entered room %s", client.name, room.name))
	} else if msg[1] == "leave" {
		if len(msg) < 3 {
			sendMessageTo(client, "room name is required! send 'client leave ROOM_NAME'")
			return
		}

		room := s.GetRoomByName(msg[2])

		if room == nil {
			sendMessageTo(client, "room was not found. Send a room name that exists")
			return
		}

		clientExistRoom := room.GetClient(client)

		if clientExistRoom == nil {
			sendMessageTo(client, fmt.Sprintf("client is not inside the room %s. See list of rooms that you're in with 'client ls_rooms'", msg[2]))
			return
		}

		room.RemoveClient(client)
		sendMessageTo(client, fmt.Sprintf("client %s left room %s", client.name, room.name))
	} else {
		sendMessageTo(client, "unknown command. Send 'help' for a list of commands")
		return
	}

}

func sendMessageTo(client *Client, msg string) {
	client.ws.Write([]byte(msg))
}
