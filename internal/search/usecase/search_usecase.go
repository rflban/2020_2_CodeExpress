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

func (sUc *SearchUsecase) Search(query string, offset uint64, limit uint64, userId uint64) (*models.Search, *ErrorResponse) {
	search := &models.Search{}
	var err error
	search.Albums, err = sUc.searchRep.SelectAlbums(query, offset, limit)
	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	search.Artists, err = sUc.searchRep.SelectArtists(query, offset, limit)
	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	search.Tracks, err = sUc.searchRep.SelectTracks(query, offset, limit, userId)
	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return search, nil
}
