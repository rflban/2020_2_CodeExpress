package models

type Track struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Duration    int    `json:"duration"`
	AlbumPoster string `json:"album_poster"`
	AlbumID     uint64 `json:"album_id"`
	Index       uint8  `json:"index"`
	Audio       string `json:"audio"`
	Artist      string `json:"artist"`
	ArtistID    uint64 `json:"artist_id"`
	IsFavorite  bool   `json:"is_favorite"`
}

//easyjson:json
type Tracks []*Track
