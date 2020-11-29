package models

import (
	"github.com/gorilla/websocket"
)

type ActiveUserInfo struct {
	Conn           *websocket.Conn
	CurrentTrack   uint64
	RadioSessionID string
}
