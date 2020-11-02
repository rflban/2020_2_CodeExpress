package delivery

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/artist"
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

type ArtistHandler struct {
	artistUsecase artist.ArtistUsecase
}

func NewArtistHandler(artistUsecase artist.ArtistUsecase) *ArtistHandler {
	return &ArtistHandler{
		artistUsecase: artistUsecase,
	}
}

func (ah *ArtistHandler) Configure(e *echo.Echo) {
	e.GET("/api/v1/artists/:id", ah.handlerArtistByID())
	e.GET("/api/v1/artists", ah.handlerArtistsByParams())
	e.POST("api/v1/artists", ah.handlerCreateArtist())
	e.PUT("/api/v1/artists/:id", ah.handlerUpdateArtist())
	e.DELETE("/api/v1/artists/:id", ah.handlerDeleteArtist())
	e.POST("/api/v1/artists/:id/photo", ah.handlerUploadArtistPhoto(), middleware.BodyLimit("10M"))
}

func (ah *ArtistHandler) handlerArtistByID() echo.HandlerFunc {
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

func (ah *ArtistHandler) handlerArtistsByParams() echo.HandlerFunc {
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

func (ah *ArtistHandler) handlerCreateArtist() echo.HandlerFunc {
	type Request struct {
		Name string `json:"name" validate:"required"`
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
			Name: req.Name,
		}

		if err := ah.artistUsecase.CreateArtist(artist); err != nil {
			return RespondWithError(err, ctx)
		}

		return ctx.JSON(http.StatusOK, artist)
	}
}

func (ah *ArtistHandler) handlerUpdateArtist() echo.HandlerFunc {
	type Request struct {
		Name string `json:"name" validate:"required"`
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

		artist := &models.Artist{
			ID:   uint64(id),
			Name: req.Name,
		}

		if err := ah.artistUsecase.UpdateArtistName(artist); err != nil {
			return RespondWithError(err, ctx)
		}

		return ctx.JSON(http.StatusOK, artist)
	}
}

func (ah *ArtistHandler) handlerDeleteArtist() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		errResp := ah.artistUsecase.DeleteArtist(uint64(id))

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		return ctx.JSON(http.StatusOK, OKRespose)
	}
}

func (ah *ArtistHandler) handlerUploadArtistPhoto() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrBadRequest, err), ctx)
		}

		photoUploader := &PhotoUploader{}

		path, err := photoUploader.UploadPhoto(ctx, "poster", "./artist_posters/")

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, err), ctx)
		}

		artist, errResp := ah.artistUsecase.GetByID(uint64(id))

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		artist.Poster = path

		if errResp := ah.artistUsecase.UpdateArtistPoster(artist); errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		return ctx.JSON(http.StatusOK, artist)
	}
}
