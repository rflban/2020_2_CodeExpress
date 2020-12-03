package delivery

import (
	"encoding/json"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/search"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/responser"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

type SearchHandler struct {
	searchUsecase search.SearchUsecase
}

func NewSearchHandler(searchUsecase search.SearchUsecase) *SearchHandler {
	return &SearchHandler{
		searchUsecase: searchUsecase,
	}
}

func (sh *SearchHandler) Configure(e *echo.Echo) {
	e.GET("/api/v1/search", sh.HandlerSearch())
}

func (sh *SearchHandler) HandlerSearch() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		query := strings.Trim(ctx.QueryParam("query"), " ")
		if len(query) == 0 {
			return RespondWithError(NewErrorResponse(ErrEmptySearchQuery, nil), ctx)
		}

		offset, err := strconv.ParseUint(ctx.QueryParam("offset"), 10, 64)
		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		limit, err := strconv.ParseUint(ctx.QueryParam("limit"), 10, 64)
		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		search, errResp := sh.searchUsecase.Search(query, offset, limit)
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := json.Marshal(search)
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}
