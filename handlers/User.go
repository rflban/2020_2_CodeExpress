package handlers

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	ID       uint64 `json:"id"`
}
