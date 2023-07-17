package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type Server struct {
	room    []Room
	clients []Client
}

func NewServer() *Server {
	return &Server{
		room:    []Room{},
		clients: []Client{},
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("new incoming conn", ws.RemoteAddr())
	client := Client{
		id: uuid.New().String(),
		ws: ws,
	}
	s.clients = append(s.clients, client)

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
		s.handleCommands(client, string(msg))
	}
}

func validateBeforeMsg(client Client) bool {
	fmt.Println("ID: ", client)
	if client.name == "" {
		client.ws.Write([]byte("Welcome! To start chatting, please register your name handle with the command: set_name YOUR_NICKNAME"))
		return false
	}

	return true
}

func (s *Server) handleCommands(client Client, msg string) {
	// validateBeforeMsg(client)
	fmt.Println("client", client)
	fmt.Println("msg", msg)

	clientMsg := strings.Fields(msg)

	if clientMsg[0] == "set_name" {
		fmt.Println("setting your name to: ", clientMsg[1])
		for idx, sClients := range s.clients {
			if sClients.id == client.id {
				fmt.Println("changing names")
				s.clients[idx].setName(clientMsg[1])
			}
		}
	}

	if clientMsg[0] == "server" {
		fmt.Println("server commands")

		if clientMsg[1] == "ls_clients" {
			fmt.Println("listing clients", client)

			listClients := ""
			for _, sClients := range s.clients {
				fmt.Println("banco", sClients)
				listClients += sClients.name + "\n"
			}
			fmt.Println("clients: ", listClients)

			client.ws.Write([]byte(listClients))
		}
	}
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.ListenAndServe(":8000", nil)

}
