package artist

import "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"

type ArtistRep interface {
	Insert(artist *models.Artist) error
	UpdateName(artist *models.Artist) error
	UpdatePoster(artist *models.Artist) error
	Delete(id uint64) error
	SelectByID(id uint64) (*models.Artist, error)
	SelectByParam(count uint64, from uint64) ([]*models.Artist, error)
	SelectByName(name string) (*models.Artist, error)
}
