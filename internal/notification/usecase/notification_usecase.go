package usecase

import (
	"fmt"
	"sync"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"

	"github.com/gorilla/websocket"
)

type NotificationUsecase struct {
	activeUsers map[uint64]*models.ActiveUserInfo
	mu          *sync.RWMutex
}

func NewNotificationUsecase() *NotificationUsecase {
	return &NotificationUsecase{
		activeUsers: make(map[uint64]*models.ActiveUserInfo),
		mu:          &sync.RWMutex{},
	}
}

func (nu *NotificationUsecase) InitWSConnection(userID uint64, ws *websocket.Conn) {
	nu.mu.Lock()
	defer nu.mu.Unlock()
	nu.activeUsers[userID].Conn = ws
}

func (nu *NotificationUsecase) NotifyUserEmoji(userID uint64, emoji rune) {
	nu.mu.Lock()
	defer nu.mu.Unlock()

	err := nu.activeUsers[userID].Conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d", emoji)))
	if err != nil {
		nu.activeUsers[userID].Conn.Close()
		delete(nu.activeUsers, userID)
	}
}

func (nu *NotificationUsecase) TrackChanged(userID, trackID uint64) {
	nu.mu.Lock()
	defer nu.mu.Unlock()
	nu.activeUsers[userID].CurrentTrack = trackID
}

func (nu *NotificationUsecase) GetActiveUserTrack(userID uint64) uint64 {
	nu.mu.Lock()
	defer nu.mu.Unlock()
	return nu.activeUsers[userID].CurrentTrack
}

func (nu *NotificationUsecase) SetUserNoTrack(userID uint64) {
	nu.mu.Lock()
	defer nu.mu.Unlock()
	nu.activeUsers[userID].CurrentTrack = 0
}
