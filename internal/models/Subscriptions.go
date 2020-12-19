package models

type Subscriptions struct {
	Subscribers   []*User `json:"subscribers"`
	Subscriptions []*User `json:"subscriptions"`
}
