package models

import "time"

type Session struct {
	Name   string
	ID     string
	Expire time.Time
}
