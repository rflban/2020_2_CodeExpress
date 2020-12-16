package search

import "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"

type SearchRep interface {
	SelectAlbums(query string, offset uint64, limit uint64) ([]*models.Album, error)
	SelectArtists(query string, offset uint64, limit uint64) ([]*models.Artist, error)
	SelectTracks(query string, offset uint64, limit uint64, userId uint64) ([]*models.Track, error)
}
