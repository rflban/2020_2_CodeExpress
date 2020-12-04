package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/user"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/mwares"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/duration_counter"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/photo_uploader"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/responser"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/validator"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type TrackHandler struct {
	trackUsecase   track.TrackUsecase
	sessionUsecase session.SessionUsecase
	userUsecase    user.UserUsecase
}

func NewTrackHandler(trackUsecase track.TrackUsecase, sessionUsecase session.SessionUsecase,
	userUsecase user.UserUsecase) *TrackHandler {
	return &TrackHandler{
		trackUsecase:   trackUsecase,
		sessionUsecase: sessionUsecase,
		userUsecase:    userUsecase,
	}
}

func (ah *TrackHandler) Configure(e *echo.Echo, mm *mwares.MiddlewareManager) {
	e.GET("/api/v1/tracks", ah.HandlerTracksByParams())
	e.POST("/api/v1/tracks", ah.HandlerCreateTrack(), mm.CheckCSRF)
	e.PUT("/api/v1/tracks/:id", ah.HandlerUpdateTrack(), mm.CheckCSRF)
	e.DELETE("/api/v1/tracks/:id", ah.HandlerDeleteTrack(), mm.CheckCSRF)
	e.POST("/api/v1/tracks/:id/audio", ah.HandlerUploadTrackAudio(), middleware.BodyLimit("10M"), mm.CheckCSRF)
	e.GET("/api/v1/artists/:id/tracks", ah.HandlerTracksByArtistID())
	e.GET("/api/v1/favorite/tracks", ah.HandlerFavouritesByUser(), mm.CheckAuth)
	e.POST("/api/v1/favorite/track/:id", ah.HandlerAddToUsersFavourites(), mm.CheckAuth, mm.CheckCSRF)
	e.DELETE("/api/v1/favorite/track/:id", ah.HandlerDeleteFromUsersFavourites(), mm.CheckAuth, mm.CheckCSRF)
}

func (ah *TrackHandler) HandlerDeleteFromUsersFavourites() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user_id := ctx.Get(ConstAuthedUserParam)

		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		if err := ah.trackUsecase.DeleteFromFavourites(user_id.(uint64), uint64(id)); err != nil {
			return RespondWithError(err, ctx)
		}

		return ctx.JSON(http.StatusOK, OKResponse)
	}
}

func (ah *TrackHandler) HandlerAddToUsersFavourites() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user_id := ctx.Get(ConstAuthedUserParam)

		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		if err := ah.trackUsecase.AddToFavourites(user_id.(uint64), uint64(id)); err != nil {
			return RespondWithError(err, ctx)
		}

		return ctx.JSON(http.StatusOK, OKResponse)
	}
}

func (ah *TrackHandler) HandlerFavouritesByUser() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user_id := ctx.Get(ConstAuthedUserParam)

		tracks, errResp := ah.trackUsecase.GetFavoritesByUserID(user_id.(uint64))

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := json.Marshal(models.Tracks(tracks))
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}

func (ah *TrackHandler) HandlerTracksByArtistID() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		var tracks []*models.Track
		var errResp *ErrorResponse
		if user := ah.tryGetUser(ctx); user != nil {
			tracks, errResp = ah.trackUsecase.GetByArtistId(uint64(id), user.ID)
		} else {
			tracks, errResp = ah.trackUsecase.GetByArtistId(uint64(id), 0)
		}
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := json.Marshal(models.Tracks(tracks))
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}

func (ah *TrackHandler) HandlerTracksByParams() echo.HandlerFunc {
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

		var tracks []*models.Track
		var errResp *ErrorResponse
		if user := ah.tryGetUser(ctx); user != nil {
			tracks, errResp = ah.trackUsecase.GetByParams(uint64(count), uint64(from), user.ID)
		} else {
			tracks, errResp = ah.trackUsecase.GetByParams(uint64(count), uint64(from), 0)
		}
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := json.Marshal(models.Tracks(tracks))
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}

func (ah *TrackHandler) HandlerCreateTrack() echo.HandlerFunc {
	type Request struct {
		Title   string `json:"title" validate:"required"`
		AlbumID uint64 `json:"album_id" validate:"required"`
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

		track := &models.Track{
			Title:   req.Title,
			AlbumID: req.AlbumID,
		}

		if err := ah.trackUsecase.CreateTrack(track); err != nil {
			return RespondWithError(err, ctx)
		}

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := json.Marshal(track)
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}

func (ah *TrackHandler) HandlerUpdateTrack() echo.HandlerFunc {
	type Request struct {
		Title   string `json:"title" validate:"required"`
		AlbumID uint64 `json:"album_id" validate:"required"`
		Index   uint8  `json:"index" validate:"required"`
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

		track := &models.Track{
			ID:      uint64(id),
			Title:   req.Title,
			AlbumID: req.AlbumID,
			Index:   req.Index,
		}

		if err := ah.trackUsecase.UpdateTrack(track); err != nil {
			return RespondWithError(err, ctx)
		}

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := json.Marshal(track)
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}

func (ah *TrackHandler) HandlerDeleteTrack() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		errResp := ah.trackUsecase.DeleteTrack(uint64(id))

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		return ctx.JSON(http.StatusOK, OKResponse)
	}
}

func (ah *TrackHandler) HandlerUploadTrackAudio() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		photoUploader := &PhotoUploader{}

		path, err := photoUploader.UploadPhoto(ctx, "audio", "./track_audio/")

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, err), ctx)
		}

		track, errResp := ah.trackUsecase.GetByID(uint64(id))

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		track.Audio = path
		track.Duration, err = CountDuration(path)

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, err), ctx)
		}

		if errResp := ah.trackUsecase.UpdateTrackAudio(track); errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := json.Marshal(track)
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}

func (ah *TrackHandler) tryGetUser(ctx echo.Context) *models.User {
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
