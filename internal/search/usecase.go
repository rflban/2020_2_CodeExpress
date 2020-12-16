package search

import (
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
)

type SearchUsecase interface {
	Search(query string, offset, limit, userId uint64) (*models.Search, *ErrorResponse)
}
