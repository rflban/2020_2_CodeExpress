package usecase

import (
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/search"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
)

type SearchUsecase struct {
	searchRep search.SearchRep
}

func NewSearchUsecase(searchRep search.SearchRep) *SearchUsecase {
	return &SearchUsecase{
		searchRep: searchRep,
	}
}

func (sUc *SearchUsecase) SearchAlbums(query string, offset uint64, limit uint64) ([]*models.Album, *ErrorResponse) {
	albums, err := sUc.searchRep.SelectAlbums(query, offset, limit)
	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return albums, nil
}

func (sUc *SearchUsecase) SearchArtists(query string, offset uint64, limit uint64) ([]*models.Artist, *ErrorResponse) {
	artists, err := sUc.searchRep.SelectArtists(query, offset, limit)
	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return artists, nil
}

func (sUc *SearchUsecase) SearchTracks(query string, offset uint64, limit uint64) ([]*models.Track, *ErrorResponse) {
	tracks, err := sUc.searchRep.SelectTracks(query, offset, limit)
	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return tracks, nil
}
