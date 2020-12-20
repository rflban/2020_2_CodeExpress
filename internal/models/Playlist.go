package models

type Playlist struct {
	ID     uint64   `json:"id"`
	UserID uint64   `json:"user_id"`
	Title  string   `json:"title"`
	Poster string   `json:"poster"`
	Tracks []*Track `json:"tracks"`
	IsPublic bool   `json:"is_public"`
}

//easyjson:json
type Playlists []*Playlist
