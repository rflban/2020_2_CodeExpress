package track

import (
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
)

type TrackUsecase interface {
	CreateTrack(track *models.Track, userId uint64) *ErrorResponse
	DeleteTrack(id uint64) *ErrorResponse
	GetByArtistId(artistId, userId uint64) ([]*models.Track, *ErrorResponse)
	GetByAlbumID(albumId, userId uint64) ([]*models.Track, *ErrorResponse)
	GetByID(id, userId uint64) (*models.Track, *ErrorResponse)
	GetByParams(count uint64, from uint64, userId uint64) ([]*models.Track, *ErrorResponse)
	GetFavoritesByUserID(userID uint64) ([]*models.Track, *ErrorResponse)
	UpdateTrack(track *models.Track, userId uint64) *ErrorResponse
	UpdateTrackAudio(track *models.Track, userId uint64) *ErrorResponse
	AddToFavourites(userID, trackID uint64) *ErrorResponse
	DeleteFromFavourites(userID, trackID uint64) *ErrorResponse
	GetByPlaylistID(playlistID, userId uint64) ([]*models.Track, *ErrorResponse)
	LikeTrack(userId, trackId uint64) *ErrorResponse
	DislikeTrack(userId, trackId uint64) *ErrorResponse
}
