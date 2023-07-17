package main

import "golang.org/x/net/websocket"

type Client struct {
	id   string
	name string
	ws   *websocket.Conn
}

func (c *Client) setName(name string) {
	c.name = name
}
