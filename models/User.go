package models

type User struct {
	ID       uint64 `json:"id"`
	Name     string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Avatar   string `json:"avatar"`
}
