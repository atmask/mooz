package models

import (
  "github.com/gorilla/websocket"
)

type BroadcastMsg struct {
	Message map[string]interface{}
	RoomID  string
	Client  *websocket.Conn
}
