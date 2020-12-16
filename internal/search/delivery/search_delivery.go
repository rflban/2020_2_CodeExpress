package delivery

import (
	"encoding/json"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/search"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/responser"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/user"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

type SearchHandler struct {
	searchUsecase  search.SearchUsecase
	sessionUsecase session.SessionUsecase
	userUsecase    user.UserUsecase
}

func NewSearchHandler(searchUsecase search.SearchUsecase, sessionUsecase session.SessionUsecase,
	userUsecase user.UserUsecase) *SearchHandler {
	return &SearchHandler{
		searchUsecase:  searchUsecase,
		sessionUsecase: sessionUsecase,
		userUsecase:    userUsecase,
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

		var search *models.Search
		var errResp *ErrorResponse
		if user := sh.tryGetUser(ctx); user != nil {
			search, errResp = sh.searchUsecase.Search(query, offset, limit, user.ID)
		} else {
			search, errResp = sh.searchUsecase.Search(query, offset, limit, 0)
		}
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

func (sh *SearchHandler) tryGetUser(ctx echo.Context) *models.User {
	cookie, err := ctx.Cookie(ConstSessionName)
	if err != nil {
		return nil
	}

	userSession, errResp := sh.sessionUsecase.GetByID(cookie.Value)
	if errResp != nil {
		return nil
	}

	user, errNoUser := sh.userUsecase.GetById(userSession.UserID)
	if errNoUser != nil {
		return nil
	}

	return user
}
