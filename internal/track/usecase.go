package track

import (
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
)

type TrackUsecase interface {
	CreateTrack(track *models.Track) *ErrorResponse
	DeleteTrack(id uint64) *ErrorResponse
	GetByArtistID(artistID uint64) ([]*models.Track, *ErrorResponse)
	GetByAlbumID(albumID uint64) ([]*models.Track, *ErrorResponse)
	GetByID(id uint64) (*models.Track, *ErrorResponse)
	GetByParams(count uint64, from uint64) ([]*models.Track, *ErrorResponse)
	GetFavoritesByUserID(userID uint64) ([]*models.Track, *ErrorResponse)
	UpdateTrack(track *models.Track) *ErrorResponse
	UpdateTrackAudio(track *models.Track) *ErrorResponse
	AddToFavourites(userID, trackID uint64) *ErrorResponse
	DeleteFromFavourites(userID, trackID uint64) *ErrorResponse
	GetByPlaylistID(playlistID uint64) ([]*models.Track, *ErrorResponse)
}
