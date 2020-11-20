package playlist

import (
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
)

type PlaylistUsecase interface {
	CreatePlaylist(playlist *models.Playlist) *ErrorResponse
	UpdatePlaylist(playlist *models.Playlist) *ErrorResponse
	DeletePlaylist(id uint64) *ErrorResponse
	GetByID(id uint64) (*models.Playlist, *ErrorResponse)
	GetByUserID(userID uint64) ([]*models.Playlist, *ErrorResponse)
	AddTrack(trackID uint64, playlistID uint64) *ErrorResponse
	DeleteTrack(trackID uint64, playlistID uint64) *ErrorResponse
}
