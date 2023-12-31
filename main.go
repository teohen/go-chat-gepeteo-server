package main

import (
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

/*

server
	ls-clients ROOM_NAME ok
	ls-clients ok
	ls-rooms ok
	cr_room ok

client
	set_name ok
	join ROOM_NAME ok
	leave ROOM_NAME ok
	ls-rooms ok

msg
	c CLIENT_NAME ok
	r ROOM_NAME ok

help
	help (server, client, room)
*/

func handleWS(ws *websocket.Conn) {
	server := NewServer()

	client := Client{
		id: uuid.New().String(),
		ws: ws,
	}

	generalRoom := server.GetRoomByName("general")
	if generalRoom == nil {
		generalRoom = &Room{
			name: "general",
		}
		server.addRoom(generalRoom)
	}

	generalRoom.AddClient(&client)
	server.readLoop(client)
}

func main() {
	http.Handle("/ws", websocket.Handler(handleWS))
	http.ListenAndServe(":8000", nil)
}
