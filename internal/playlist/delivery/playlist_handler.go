package delivery

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/mwares"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/playlist"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/photo_uploader"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/responser"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/validator"
)

type PlaylistHandler struct {
	playlistUsecase playlist.PlaylistUsecase
	trackUsecase    track.TrackUsecase
}

func NewPlaylistHandler(playlistUsecase playlist.PlaylistUsecase, trackUsecase track.TrackUsecase) *PlaylistHandler {
	return &PlaylistHandler{
		playlistUsecase: playlistUsecase,
		trackUsecase:    trackUsecase,
	}
}

func (ph *PlaylistHandler) Configure(e *echo.Echo, mm *mwares.MiddlewareManager) {
	e.POST("/api/v1/playlists", ph.HandlerCreatePlaylist(), mm.CheckAuth)
	e.GET("/api/v1/playlists", ph.HandlerUserPlaylists(), mm.CheckAuth)
	e.GET("/api/v1/playlists/:id", ph.HandlerConcretePlaylist())
	e.PUT("/api/v1/playlists/:id", ph.HandlerUpdatePlaylist(), mm.CheckAuth)
	e.DELETE("/api/v1/playlists/:id", ph.HandlerDeletePlaylist(), mm.CheckAuth)
	e.POST("/api/v1/playlists/:id/photo", ph.HandlerUploadPlaylistPhoto(), middleware.BodyLimit("10M"), mm.CheckAuth)
	e.POST("/api/v1/playlists/:id/tracks", ph.HandlerAddTrackToPlaylist(), mm.CheckAuth)
	e.DELETE("/api/v1/playlists/:id/tracks/:track_id", ph.HandlerDeleteTrackFromPlaylist(), mm.CheckAuth)
}

func (ph *PlaylistHandler) HandlerCreatePlaylist() echo.HandlerFunc {
	type Request struct {
		Title string `json:"title"`
	}

	return func(ctx echo.Context) error {
		userID := ctx.Get(ConstAuthedUserParam)

		req := &Request{}

		if err := validator.NewRequestValidator(ctx).Validate(req); err != nil {
			if err.Error != nil {
				logrus.Info(err.Error)
				validator.GetValidationError(err)
			}
			return ctx.JSON(err.StatusCode, err.UserError)
		}

		playlist := &models.Playlist{
			UserID: userID.(uint64),
			Title:  req.Title,
		}

		if errResp := ph.playlistUsecase.CreatePlaylist(playlist); errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := playlist.MarshalJSON()
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}

func (ph *PlaylistHandler) HandlerUpdatePlaylist() echo.HandlerFunc {
	type Request struct {
		Title string `json:"title"`
	}

	return func(ctx echo.Context) error {
		userID := ctx.Get(ConstAuthedUserParam)

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

		playlist := &models.Playlist{
			ID:     uint64(id),
			UserID: userID.(uint64),
			Title:  req.Title,
		}

		if errResp := ph.playlistUsecase.UpdatePlaylist(playlist); errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := playlist.MarshalJSON()
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}

func (ph *PlaylistHandler) HandlerDeletePlaylist() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		if errResp := ph.playlistUsecase.DeletePlaylist(uint64(id)); errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		return ctx.JSON(http.StatusOK, OKResponse)
	}
}

func (ph *PlaylistHandler) HandlerConcretePlaylist() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		playlist, errResp := ph.playlistUsecase.GetByID(uint64(id))

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		tracks, errResp := ph.trackUsecase.GetByPlaylistID(playlist.ID)

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		playlist.Tracks = tracks

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := playlist.MarshalJSON()
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}

func (ph *PlaylistHandler) HandlerUserPlaylists() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID := ctx.Get(ConstAuthedUserParam)

		playlists, errResp := ph.playlistUsecase.GetByUserID(userID.(uint64))

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := models.Playlists(playlists).MarshalJSON()
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}

func (ph *PlaylistHandler) HandlerUploadPlaylistPhoto() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		photoUploader := &PhotoUploader{}

		path, err := photoUploader.UploadPhoto(ctx, "poster", "./playlist_posters/")

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, err), ctx)
		}

		playlist, errResp := ph.playlistUsecase.GetByID(uint64(id))

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		playlist.Poster = path

		if errResp := ph.playlistUsecase.UpdatePlaylist(playlist); errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := playlist.MarshalJSON()
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}

func (ph *PlaylistHandler) HandlerAddTrackToPlaylist() echo.HandlerFunc {
	type Request struct {
		TrackID uint64 `json:"track_id"`
	}

	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		req := &Request{}

		if err := validator.NewRequestValidator(ctx).Validate(req); err != nil {
			if err.Error != nil {
				logrus.Info(err.Error)
				validator.GetValidationError(err)
			}
			return ctx.JSON(err.StatusCode, err.UserError)
		}

		if errResp := ph.playlistUsecase.AddTrack(req.TrackID, uint64(id)); errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		return ctx.JSON(http.StatusOK, OKResponse)
	}
}

func (ph *PlaylistHandler) HandlerDeleteTrackFromPlaylist() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		trackID, err := strconv.Atoi(ctx.Param("track_id"))
		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		if errResp := ph.playlistUsecase.DeleteTrack(uint64(trackID), uint64(id)); errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		return ctx.JSON(http.StatusOK, OKResponse)
	}
}
