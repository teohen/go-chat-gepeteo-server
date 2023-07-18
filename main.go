package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

/*

server
	ls-clients
	ls-rooms
	cr_room

client
	set_name

msg
	c CLIENT_NAME
	r ROOM_NAME
*/

type Server struct {
	rooms   []*Room
	clients []*Client
}

func NewServer() *Server {
	server := &Server{
		rooms:   []*Room{},
		clients: []*Client{},
	}

	generalRoom := Room{
		name: "general",
	}

	server.addRoom(&generalRoom)
	return server
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("new incoming conn", ws.RemoteAddr())
	client := Client{
		id: uuid.New().String(),
		ws: ws,
	}

	generalRoom := s.getRoom("general")

	if generalRoom != nil {
		generalRoom.AddClient(&client)
	}

	s.readLoop(client)
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
		sClient := s.getClient(&client)
		s.handleCommands(sClient, string(msg))
	}
}

func (s *Server) getClient(client *Client) *Client {
	for _, sClient := range s.clients {
		if sClient.id == client.id {
			return sClient
		}
	}
	return client
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
	fmt.Println(client)
	if client.name == "" && msg[0] != "client" && msg[1] != "set_name" {
		returnMessage(client, "Welcome! To start chatting, please register your name handle with the command 'client set_name YOUR_NICKNAME'")
		return false
	}
	return true
}

func (s *Server) handleCommands(client *Client, msg string) {
	clientMsg := strings.Fields(msg)

	if ok := validateBeforeMsg(client, clientMsg); ok != true {
		return
	}

	if clientMsg[0] == "server" {
		s.serverCommands(client, clientMsg)
	} else if clientMsg[0] == "client" {
		s.clientCommands(client, clientMsg)
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
	}

	if msg[1] == "ls_rooms" {
		listRooms := "Rooms:\n"
		fmt.Println(s.rooms)
		for _, sRoom := range s.rooms {
			listRooms += sRoom.name + "\n"
		}
		client.ws.Write([]byte(fmt.Sprintf("%s\n", listRooms)))
		return
	}

	if msg[1] == "cr_room" {
		if len(msg) < 3 {
			returnMessage(client, fmt.Sprintf("room name is required! send 'server cr_room ROOM_NAME'"))
			return
		}

		roomName := msg[2]

		for _, sRoom := range s.rooms {
			if sRoom.name == roomName {
				returnMessage(client, fmt.Sprintf("Room name: %s already exists", roomName))
				return
			}
			room := Room{
				name: roomName,
			}

			room.AddClient(client)
			s.addRoom(&room)
		}
	}
}

func (s *Server) clientCommands(client *Client, msg []string) {
	if msg[1] == "set_name" {
		client.setName(msg[2])

		generalRoom := s.getRoom("general")
		generalRoom.AddClient(client)
		returnMessage(client, fmt.Sprintf("name set: %s", msg[2]))
	}
}

func returnMessage(client *Client, msg string) {
	client.ws.Write([]byte(msg))
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.ListenAndServe(":8000", nil)

}
