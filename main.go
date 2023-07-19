package main

import (
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

/*

server
	ls-clients ROOM_NAME
	ls-clients ok
	ls-rooms ok
	cr_room ok

client
	set_name ok
	join ROOM_NAME
	leave ROOM_NAME

msg
	c CLIENT_NAME
	r ROOM_NAME

help
	help (server, client, room)
*/

func handleWS(ws *websocket.Conn) {
	server := NewServer()

	client := Client{
		id: uuid.New().String(),
		ws: ws,
	}

	generalRoom := server.getRoom("general")
	if generalRoom == nil {
		generalRoom = &Room{
			name: "general",
		}
	}

	generalRoom.AddClient(&client)
	server.addRoom(generalRoom)
	server.readLoop(client)
}

func main() {
	http.Handle("/ws", websocket.Handler(handleWS))
	http.ListenAndServe(":8000", nil)
}
