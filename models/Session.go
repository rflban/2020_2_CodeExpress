package models

import "time"

type Session struct {
	Name   string
	ID     string
	UserID uint64
	Expire time.Time
}
