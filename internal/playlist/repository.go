package playlist

import (
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
)

type PlaylistRep interface {
	Insert(playlist *models.Playlist) error
	Update(playlist *models.Playlist) error
	Delete(id uint64) error
	SelectByID(id uint64) (*models.Playlist, error)
	SelectByUserID(userID uint64) ([]*models.Playlist, error)
	InsertTrack(trackID uint64, playlistID uint64) error
	DeleteTrack(trackID uint64, playlistID uint64) error
	SelectPublicByUserID(userID uint64) ([]*models.Playlist, error)
}
