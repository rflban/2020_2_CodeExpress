package models

type NewUser struct {
	Name     string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RepeatedPassword string `json:"repeated_password"`
}
