package models

type Artist struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Poster string `json:"poster"`
}
