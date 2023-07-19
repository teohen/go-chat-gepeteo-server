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
		generalRoom := s.getRoom("general")
		sClient := generalRoom.GetClient(&client)
		s.handleCommands(sClient, string(msg))
	}
}

func (s *Server) getRoom(name string) *Room {
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
	fmt.Println("server commands")
	generalRoom := s.getRoom("general")

	if msg[1] == "ls_clients" {
		listClients := "Clients:\n"

		for _, gClient := range generalRoom.clients {
			fmt.Println("clients", gClient.name)
			listClients += gClient.name + "\n"
		}

		client.ws.Write([]byte(listClients))
		return
	} else if msg[1] == "ls_rooms" {
		listRooms := "Rooms:\n"
		for _, sRoom := range s.rooms {
			listRooms += sRoom.name + "\n"
		}
		client.ws.Write([]byte(fmt.Sprintf("%s\n", listRooms)))
		return
	} else if msg[1] == "cr_room" {
		if len(msg) < 3 {
			sendMessageTo(client, fmt.Sprintf("room name is required! send 'server cr_room ROOM_NAME'"))
			return
		}

		roomName := msg[2]

		for _, sRoom := range s.rooms {
			if sRoom.name == roomName {
				sendMessageTo(client, fmt.Sprintf("Room name: %s already exists", roomName))
				return
			}
			room := Room{
				name: roomName,
			}

			room.AddClient(client)
			s.addRoom(&room)
		}
	} else {
		sendMessageTo(client, "unknown command. Send 'help' for a list of commands")
		return
	}
}

func (s *Server) clientCommands(client *Client, msg []string) {
	if msg[1] == "set_name" {
		if len(msg) < 3 {
			sendMessageTo(client, fmt.Sprintf("client name is required! send 'client set_name CLIENT_NAME'"))
			return
		}

		client.setName(msg[2])

		generalRoom := s.getRoom("general")
		generalRoom.AddClient(client)
		sendMessageTo(client, fmt.Sprintf("name set: %s", msg[2]))
	} else {
		sendMessageTo(client, "unknown command. Send 'help' for a list of commands")
		return
	}

}

func sendMessageTo(client *Client, msg string) {
	client.ws.Write([]byte(msg))
}
