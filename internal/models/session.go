package models

import (
	"time"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	uuid "github.com/satori/go.uuid"
)

type Session struct {
	ID     string
	Name   string
	UserID uint64
	Expire time.Time
}

func NewSession(userID uint64) *Session {
	return &Session{
		ID:     uuid.NewV4().String(),
		Name:   ConstSessionName,
		UserID: userID,
		Expire: time.Now().AddDate(0, 0, ConstDaysSession),
	}
}
