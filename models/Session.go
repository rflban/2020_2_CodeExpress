package models

import "time"

type Session struct {
	ID     string
	Expire time.Time
}
