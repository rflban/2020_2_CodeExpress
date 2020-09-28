package models

type PasswordForm struct {
	Password         string `json:"password"`
	RepeatedPassword string `json:"repeated_password"`
}
