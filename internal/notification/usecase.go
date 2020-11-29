package notification

import (
	"github.com/gorilla/websocket"
)

type NotificationUsecase interface {
	InitWSConnection(userID uint64, ws *websocket.Conn)
	GetActiveUserTrack(userID uint64) uint64
	TrackChanged(userID, trackID uint64)
	NotifyUserEmoji(userID uint64, emoji rune)
	SetUserNoTrack(userID uint64)
}
