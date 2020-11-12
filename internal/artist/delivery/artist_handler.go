package delivery

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/artist"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/mwares"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/photo_uploader"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/responser"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type ArtistHandler struct {
	artistUsecase artist.ArtistUsecase
}

func NewArtistHandler(artistUsecase artist.ArtistUsecase) *ArtistHandler {
	return &ArtistHandler{
		artistUsecase: artistUsecase,
	}
}

func (ah *ArtistHandler) Configure(e *echo.Echo, mm *mwares.MiddlewareManager) {
	e.GET("/api/v1/artists/:id", ah.HandlerArtistByID())
	e.GET("/api/v1/artists", ah.HandlerArtistsByParams())
	e.POST("api/v1/artists", ah.HandlerCreateArtist(), mm.CheckCSRF)
	e.PUT("/api/v1/artists/:id", ah.HandlerUpdateArtist(), mm.CheckCSRF)
	e.DELETE("/api/v1/artists/:id", ah.HandlerDeleteArtist(), mm.CheckCSRF)
	e.POST("/api/v1/artists/:id/photo", ah.HandlerUploadArtistPhoto(), middleware.BodyLimit("10M"), mm.CheckCSRF)
}

func (ah *ArtistHandler) HandlerArtistByID() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		artist, errResp := ah.artistUsecase.GetByID(uint64(id))

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		return ctx.JSON(http.StatusOK, artist)
	}
}

func (ah *ArtistHandler) HandlerArtistsByParams() echo.HandlerFunc {
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

		artists, errResp := ah.artistUsecase.GetByParams(uint64(count), uint64(from))

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		return ctx.JSON(http.StatusOK, artists)
	}
}

func (ah *ArtistHandler) HandlerCreateArtist() echo.HandlerFunc {
	type Request struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
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

		artist := &models.Artist{
			Name:        req.Name,
			Description: req.Description,
		}

		if err := ah.artistUsecase.CreateArtist(artist); err != nil {
			return RespondWithError(err, ctx)
		}

		return ctx.JSON(http.StatusOK, artist)
	}
}

func (ah *ArtistHandler) HandlerUpdateArtist() echo.HandlerFunc {
	type Request struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
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

		artist, errResp := ah.artistUsecase.GetByName(req.Name)

		if errResp == nil && artist.ID != uint64(id) {
			return RespondWithError(NewErrorResponse(ErrNameAlreadyExist, nil), ctx)
		}

		artist = &models.Artist{
			ID:          uint64(id),
			Name:        req.Name,
			Description: req.Description,
		}

		if err := ah.artistUsecase.UpdateArtist(artist); err != nil {
			return RespondWithError(err, ctx)
		}

		return ctx.JSON(http.StatusOK, artist)
	}
}

func (ah *ArtistHandler) HandlerDeleteArtist() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		errResp := ah.artistUsecase.DeleteArtist(uint64(id))

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		return ctx.JSON(http.StatusOK, OKResponse)
	}
}

func (ah *ArtistHandler) HandlerUploadArtistPhoto() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		artist, errResp := ah.artistUsecase.GetByID(uint64(id))

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		photoUploader := &PhotoUploader{}
		posterField, avatarField := "poster", "avatar"
		isChanged := false

		if _, err := ctx.FormFile(posterField); err == nil {
			path, err := photoUploader.UploadPhoto(ctx, posterField, "./artist_posters/")

			if err != nil {
				return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
			}
			isChanged = true
			artist.Poster = path
		}

		if _, err := ctx.FormFile(avatarField); err == nil {
			path, err := photoUploader.UploadPhoto(ctx, avatarField, "./artist_avatars/")

			if err != nil {
				return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
			}
			isChanged = true
			artist.Avatar = path
		}

		if isChanged {
			if err := ah.artistUsecase.UpdateArtist(artist); err != nil {
				return RespondWithError(err, ctx)
			}
		} else {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		return ctx.JSON(http.StatusOK, artist)
	}
}
