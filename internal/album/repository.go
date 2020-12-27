package album

import "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"

//go:generate mockgen -destination mock_album/mock_album.go github.com/go-park-mail-ru/2020_2_CodeExpress/internal/album AlbumRep,AlbumUsecase

type AlbumRep interface {
	Insert(album *models.Album) error
	Update(album *models.Album) error
	UpdatePoster(album *models.Album) error
	Delete(id uint64) error
	SelectByID(id uint64) (*models.Album, error)
	SelectByParam(count uint64, from uint64) ([]*models.Album, error)
	SelectTopByParam(count uint64, from uint64) ([]*models.Album, error)
	SelectByArtistID(artistID uint64) ([]*models.Album, error)
}
