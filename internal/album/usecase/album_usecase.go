package usecase

import (
	"database/sql"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/album"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
)

type AlbumUsecase struct {
	albumRep album.AlbumRep
}

func NewAlbumUsecase(albumRep album.AlbumRep) *AlbumUsecase {
	return &AlbumUsecase{
		albumRep: albumRep,
	}
}

func (aUc *AlbumUsecase) CreateAlbum(album *models.Album) *ErrorResponse {
	if err := aUc.albumRep.Insert(album); err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	return nil
}

func (aUc *AlbumUsecase) DeleteAlbum(id uint64) *ErrorResponse {
	err := aUc.albumRep.Delete(id)

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrAlbumNotExist, err)
	}

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	return nil
}

func (aUc *AlbumUsecase) GetByID(id uint64) (*models.Album, *ErrorResponse) {
	album, err := aUc.albumRep.SelectByID(id)

	if err == sql.ErrNoRows {
		return nil, NewErrorResponse(ErrAlbumNotExist, err)
	}

	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return album, nil
}

func (aUc *AlbumUsecase) GetByArtistID(artistID uint64) ([]*models.Album, *ErrorResponse) {
	albums, err := aUc.albumRep.SelectByArtistID(artistID)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return albums, nil
}

func (aUc *AlbumUsecase) UpdateAlbum(album *models.Album) *ErrorResponse {
	err := aUc.albumRep.Update(album)

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrAlbumNotExist, err)
	}

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	return nil
}

func (aUc *AlbumUsecase) UpdateAlbumPoster(album *models.Album) *ErrorResponse {
	err := aUc.albumRep.UpdatePoster(album)

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrAlbumNotExist, err)
	}

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	return nil
}

func (aUc *AlbumUsecase) GetByParams(count uint64, from uint64) ([]*models.Album, *ErrorResponse) {
	albums, err := aUc.albumRep.SelectByParam(count, from)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return albums, nil
}

func (aUc *AlbumUsecase) GetTopByParams(count uint64, from uint64) ([]*models.Album, *ErrorResponse) {
	albums, err := aUc.albumRep.SelectTopByParam(count, from)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return albums, nil
}
