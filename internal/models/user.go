package models

type User struct {
	ID           uint64 `json:"id"`
	Name         string `json:"username"`
	Email        string `json:"email,omitempty"`
	Password     string `json:"-"`
	Avatar       string `json:"avatar"`
	IsSubscribed bool   `json:"is_subscribed,omitempty"`
}
