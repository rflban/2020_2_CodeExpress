package models

type Album struct {
	ID         uint64   `json:"id"`
	Title      string   `json:"title"`
	ArtistID   uint64   `json:"artist_id"`
	ArtistName string   `json:"artist_name"`
	Poster     string   `json:"poster"`
	Tracks     []*Track `json:"tracks"`
}
