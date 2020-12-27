package track

import "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"

type TrackRep interface {
	Insert(track *models.Track) error
	InsertTrackToUser(userID, trackID uint64) error
	DeleteTrackFromUsersTracks(userID, trackID uint64) error
	Update(track *models.Track) error
	UpdateAudio(track *models.Track) error
	Delete(id uint64) error
	SelectByArtistId(artistId, userId uint64) ([]*models.Track, error)
	SelectByID(id, userId uint64) (*models.Track, error)
	SelectByParams(count, from, userId uint64) ([]*models.Track, error)
	SelectTopByParams(count, from, userId uint64) ([]*models.Track, error)
	SelectByAlbumID(albumID, userId uint64) ([]*models.Track, error)
	SelectFavoritesByUserID(userID uint64) ([]*models.Track, error)
	SelectByPlaylistID(playlistID, userId uint64) ([]*models.Track, error)
	LikeTrack(userId, trackId uint64) error
	DislikeTrack(userId, trackId uint64) error
	SelectRandomByArtistId(artistId, userId, count uint64) ([]*models.Track, error)
}
