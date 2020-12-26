package album

import (
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
)

type AlbumUsecase interface {
	CreateAlbum(album *models.Album) *ErrorResponse
	DeleteAlbum(id uint64) *ErrorResponse
	GetByID(id uint64) (*models.Album, *ErrorResponse)
	GetByArtistID(artistID uint64) ([]*models.Album, *ErrorResponse)
	GetByParams(count uint64, from uint64) ([]*models.Album, *ErrorResponse)
	GetTopByParams(count uint64, from uint64) ([]*models.Album, *ErrorResponse)
	UpdateAlbum(album *models.Album) *ErrorResponse
	UpdateAlbumPoster(album *models.Album) *ErrorResponse
}
