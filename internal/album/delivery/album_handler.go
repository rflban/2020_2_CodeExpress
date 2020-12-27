package delivery

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/user"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/mwares"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/artist"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/album"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/photo_uploader"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/responser"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type AlbumHandler struct {
	albumUsecase   album.AlbumUsecase
	artistUsecase  artist.ArtistUsecase
	trackUsecase   track.TrackUsecase
	sessionUsecase session.SessionUsecase
	userUsecase    user.UserUsecase
}

func NewAlbumHandler(albumUsecase album.AlbumUsecase, artistUsecase artist.ArtistUsecase,
	trackUsecase track.TrackUsecase, sessionUsecase session.SessionUsecase, userUsecase user.UserUsecase) *AlbumHandler {
	return &AlbumHandler{
		albumUsecase:   albumUsecase,
		artistUsecase:  artistUsecase,
		trackUsecase:   trackUsecase,
		sessionUsecase: sessionUsecase,
		userUsecase:    userUsecase,
	}
}

func (ah *AlbumHandler) Configure(e *echo.Echo, mm *mwares.MiddlewareManager) {
	e.GET("/api/v1/artists/:id/albums", ah.HandlerAlbumsByArtist())
	e.GET("/api/v1/albums", ah.HandlerAlbumsByParams())
	e.GET("/api/v1/albums/top", ah.HandlerTopAlbumsByParams())
	e.GET("/api/v1/albums/:id", ah.HandlerAlbumTracks())
	e.POST("api/v1/albums", ah.HandlerCreateAlbum(), mm.CheckCSRF)
	e.PUT("/api/v1/albums/:id", ah.HandlerUpdateAlbum(), mm.CheckCSRF)
	e.DELETE("/api/v1/albums/:id", ah.HandlerDeleteAlbum(), mm.CheckCSRF)
	e.POST("/api/v1/albums/:id/photo", ah.HandlerUploadAlbumPhoto(), middleware.BodyLimit("10M"), mm.CheckCSRF)
}

func (ah *AlbumHandler) HandlerAlbumsByArtist() echo.HandlerFunc {
	type Response struct {
		Artist *models.Artist  `json:"artist"`
		Albums []*models.Album `json:"albums"`
	}

	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		artist, errResp := ah.artistUsecase.GetByID(uint64(id))

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		albums, errResp := ah.albumUsecase.GetByArtistID(uint64(id))

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		res := &Response{
			Artist: artist,
			Albums: albums,
		}

		return ctx.JSON(http.StatusOK, res)
	}
}

func (ah *AlbumHandler) HandlerCreateAlbum() echo.HandlerFunc {
	type Request struct {
		Title    string `json:"title" validate:"required"`
		ArtistID uint64 `json:"artist_id" validate:"required"`
	}

	return func(ctx echo.Context) error {
		req := &Request{}

		if err := validator.NewRequestValidator(ctx).Validate(req); err != nil {
			if err.Error != nil {
				logrus.Info(err.Error)
				validator.GetValidationError(err)
			}
			return ctx.JSON(err.StatusCode, err.UserError)
		}

		artist, err := ah.artistUsecase.GetByID(req.ArtistID)

		if err != nil {
			return RespondWithError(err, ctx)
		}

		album := &models.Album{
			Title:      req.Title,
			ArtistID:   req.ArtistID,
			ArtistName: artist.Name,
		}

		if err := ah.albumUsecase.CreateAlbum(album); err != nil {
			return RespondWithError(err, ctx)
		}

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := json.Marshal(album)
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}

func (ah *AlbumHandler) HandlerUpdateAlbum() echo.HandlerFunc {
	type Request struct {
		Title    string `json:"title" validate:"required"`
		ArtistID uint64 `json:"artist_id" validate:"required"`
	}

	return func(ctx echo.Context) error {
		req := &Request{}

		if err := validator.NewRequestValidator(ctx).Validate(req); err != nil {
			if err.Error != nil {
				logrus.Info(err.Error)
				validator.GetValidationError(err)
			}
			return ctx.JSON(err.StatusCode, err.UserError)
		}

		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		artist, errResp := ah.artistUsecase.GetByID(req.ArtistID)

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		album := &models.Album{
			ID:         uint64(id),
			Title:      req.Title,
			ArtistID:   req.ArtistID,
			ArtistName: artist.Name,
		}

		if err := ah.albumUsecase.UpdateAlbum(album); err != nil {
			return RespondWithError(err, ctx)
		}

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := json.Marshal(album)
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}

func (ah *AlbumHandler) HandlerDeleteAlbum() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		errResp := ah.albumUsecase.DeleteAlbum(uint64(id))

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		return ctx.JSON(http.StatusOK, OKResponse)
	}
}

func (ah *AlbumHandler) HandlerUploadAlbumPhoto() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		photoUploader := &PhotoUploader{}

		path, err := photoUploader.UploadPhoto(ctx, "poster", "./album_posters/")

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, err), ctx)
		}

		album, errResp := ah.albumUsecase.GetByID(uint64(id))

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		album.Poster = path

		if errResp := ah.albumUsecase.UpdateAlbumPoster(album); errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := json.Marshal(album)
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}

func (ah *AlbumHandler) HandlerAlbumsByParams() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		params := ctx.QueryParams()
		count, err := strconv.Atoi(params.Get("count"))
		if err != nil || count < 0 {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		from, err := strconv.Atoi(params.Get("from"))
		if err != nil || from < 0 {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		albums, errResp := ah.albumUsecase.GetByParams(uint64(count), uint64(from))

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := json.Marshal(models.Albums(albums))
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}

func (ah *AlbumHandler) HandlerAlbumTracks() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		album, errResp := ah.albumUsecase.GetByID(uint64(id))

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		var userId uint64
		if user := ah.tryGetUser(ctx); user != nil {
			userId = user.ID
		}

		tracks, errResp := ah.trackUsecase.GetByAlbumID(uint64(id), userId)

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		album.Tracks = tracks

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := json.Marshal(album)
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}

func (ah *AlbumHandler) tryGetUser(ctx echo.Context) *models.User {
	cookie, err := ctx.Cookie(ConstSessionName)
	if err != nil {
		return nil
	}

	userSession, errResp := ah.sessionUsecase.GetByID(cookie.Value)
	if errResp != nil {
		return nil
	}

	user, errNoUser := ah.userUsecase.GetById(userSession.UserID)
	if errNoUser != nil {
		return nil
	}

	return user
}

func (ah *AlbumHandler) HandlerTopAlbumsByParams() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		params := ctx.QueryParams()
		count, err := strconv.Atoi(params.Get("count"))
		if err != nil || count < 0 {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		from, err := strconv.Atoi(params.Get("from"))
		if err != nil || from < 0 {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		albums, errResp := ah.albumUsecase.GetTopByParams(uint64(count), uint64(from))

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := json.Marshal(models.Albums(albums))
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}
