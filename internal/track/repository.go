package track

import "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"

type TrackRep interface {
	Insert(track *models.Track) error
	InsertTrackToUser(userID, trackID uint64) error
	DeleteTrackFromUsersTracks(userID, trackID uint64) error
	Update(track *models.Track) error
	UpdateAudio(track *models.Track) error
	Delete(id uint64) error
	SelectByArtistID(artistID uint64) ([]*models.Track, error)
	SelectByID(id uint64) (*models.Track, error)
	SelectByParam(count uint64, from uint64) ([]*models.Track, error)
	SelectByAlbumID(albumID uint64) ([]*models.Track, error)
	SelectFavoritesByUserID(userID uint64) ([]*models.Track, error)
}
