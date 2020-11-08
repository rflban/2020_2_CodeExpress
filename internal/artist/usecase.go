package artist

import (
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
)

type ArtistUsecase interface {
	CreateArtist(artist *models.Artist) *ErrorResponse
	DeleteArtist(id uint64) *ErrorResponse
	GetByID(id uint64) (*models.Artist, *ErrorResponse)
	GetByName(name string) (*models.Artist, *ErrorResponse)
	GetByParams(count uint64, from uint64) ([]*models.Artist, *ErrorResponse)
	UpdateArtist(artist *models.Artist) *ErrorResponse
}
