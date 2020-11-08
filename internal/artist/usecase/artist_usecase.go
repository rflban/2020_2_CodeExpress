package usecase

import (
	"database/sql"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/artist"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
)

type ArtistUsecase struct {
	artistRep artist.ArtistRep
}

func NewArtistUsecase(artistRep artist.ArtistRep) *ArtistUsecase {
	return &ArtistUsecase{
		artistRep: artistRep,
	}
}

func (aUc *ArtistUsecase) CreateArtist(artist *models.Artist) *ErrorResponse {
	exists, err := aUc.checkNameExists(artist.Name)

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}
	if exists {
		return NewErrorResponse(ErrNameAlreadyExist, nil)
	}

	if err := aUc.artistRep.Insert(artist); err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	return nil
}

func (aUc *ArtistUsecase) DeleteArtist(id uint64) *ErrorResponse {
	err := aUc.artistRep.Delete(id)

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrArtistNotExist, err)
	}

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	return nil
}

func (aUc *ArtistUsecase) GetByID(id uint64) (*models.Artist, *ErrorResponse) {
	artist, err := aUc.artistRep.SelectByID(id)

	if err == sql.ErrNoRows {
		return nil, NewErrorResponse(ErrArtistNotExist, err)
	}

	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return artist, nil
}

func (aUc *ArtistUsecase) GetByParams(count uint64, from uint64) ([]*models.Artist, *ErrorResponse) {
	artists, err := aUc.artistRep.SelectByParam(count, from)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return artists, nil
}

func (aUc *ArtistUsecase) UpdateArtist(artist *models.Artist) *ErrorResponse {
	err := aUc.artistRep.Update(artist)

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrArtistNotExist, err)
	}

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	return nil
}

func (aUc *ArtistUsecase) GetByName(name string) (*models.Artist, *ErrorResponse) {
	artist, err := aUc.artistRep.SelectByName(name)

	if err == sql.ErrNoRows {
		return nil, NewErrorResponse(ErrArtistNotExist, err)
	}

	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return artist, nil
}

func (aUc *ArtistUsecase) checkNameExists(name string) (bool, error) {
	_, err := aUc.artistRep.SelectByName(name)
	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
