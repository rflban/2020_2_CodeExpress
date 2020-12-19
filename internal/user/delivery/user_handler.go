package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/mwares"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/csrf"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/builder"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/photo_uploader"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/responser"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/validator"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	userUsecase    user.UserUsecase
	sessionUsecase session.SessionUsecase
}

func NewUserHandler(userUsecase user.UserUsecase, sessionUsecase session.SessionUsecase) *UserHandler {
	return &UserHandler{
		userUsecase:    userUsecase,
		sessionUsecase: sessionUsecase,
	}
}

func (uh *UserHandler) Configure(e *echo.Echo, mm *mwares.MiddlewareManager) {
	e.POST("/api/v1/user", uh.HandlerRegisterUser())
	e.GET("/api/v1/user", uh.HandlerCurrentUserInfo(), mm.CheckAuth)
	e.PUT("/api/v1/user/profile", uh.HandlerUpdateProfile(), mm.CheckCSRF, mm.CheckAuth)
	e.PUT("/api/v1/user/password", uh.HandlerUpdatePassword(), mm.CheckCSRF, mm.CheckAuth)
	e.PUT("/api/v1/user/photo", uh.HandlerUpdateAvatar(), mm.CheckCSRF, mm.CheckAuth,
		middleware.BodyLimit("10M"))
	e.GET("/api/v1/user/:name/profile", uh.HandlerGetProfile(), mm.CheckAuth)
	e.POST("/api/v1/user/:name/subscription", uh.HandlerAddSubscription(), mm.CheckAuth)
	e.DELETE("/api/v1/user/:name/subscription", uh.HandlerRemoveSubscription(), mm.CheckAuth)
	e.GET("/api/v1/user/:name/subscriptions", uh.HandlerGetSubscriptions(), mm.CheckAuth)
}

func (uh *UserHandler) HandlerRegisterUser() echo.HandlerFunc {
	type Request struct {
		Name             string `json:"username" validate:"required"`
		Email            string `json:"email" validate:"required,email"`
		Password         string `json:"password" validate:"required,gte=8"`
		RepeatedPassword string `json:"repeated_password" validate:"required,eqfield=Password"`
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

		user, err := uh.userUsecase.Create(req.Name, req.Email, req.Password)
		if err != nil {
			return RespondWithError(err, ctx)
		}

		session := models.NewSession(user.ID)
		if err := uh.sessionUsecase.CreateSession(session); err != nil {
			return RespondWithError(err, ctx)
		}

		token, e := csrf.NewCSRFToken(session)

		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		ctx.Response().Header().Set("X-CSRF-TOKEN", token)

		cookie := builder.BuildCookie(session)
		ctx.SetCookie(cookie)
		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := json.Marshal(user)
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}

func (uh *UserHandler) HandlerCurrentUserInfo() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user_id := ctx.Get(ConstAuthedUserParam)

		if user_id == nil {
			return RespondWithError(NewErrorResponse(ErrNotAuthorized, nil), ctx)
		}

		user, errResp := uh.userUsecase.GetById(user_id.(uint64))
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := json.Marshal(user)
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}

func (uh *UserHandler) HandlerUpdateProfile() echo.HandlerFunc {
	type Request struct {
		Name  string `json:"username" validate:"required"`
		Email string `json:"email" validate:"required,email"`
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

		user_id := ctx.Get(ConstAuthedUserParam)

		user, errResp := uh.userUsecase.UpdateProfile(user_id.(uint64), req.Name, req.Email)
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := json.Marshal(user)
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}

func (uh *UserHandler) HandlerUpdatePassword() echo.HandlerFunc {
	type Request struct {
		OldPassword      string `json:"old_password" validate:"required"`
		Password         string `json:"password" validate:"required,gte=8"`
		RepeatedPassword string `json:"repeated_password" validate:"required,eqfield=Password"`
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

		user_id := ctx.Get(ConstAuthedUserParam)

		errResp := uh.userUsecase.UpdatePassword(user_id.(uint64), req.OldPassword, req.Password)
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		return ctx.JSON(http.StatusOK, OKResponse)
	}
}

func (uh *UserHandler) HandlerUpdateAvatar() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user_id := ctx.Get(ConstAuthedUserParam)

		user, errResp := uh.userUsecase.GetById(user_id.(uint64))
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		photoUploader := photo_uploader.PhotoUploader{}
		avatarPath, err := photoUploader.UploadPhoto(ctx, "avatar", "./avatars/")
		if err != nil {
			errResp := NewErrorResponse(ErrBadRequest, err)
			return RespondWithError(errResp, ctx)
		}

		user, errResp = uh.userUsecase.UpdateAvatar(user_id.(uint64), avatarPath)
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := json.Marshal(user)
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}

func (uh *UserHandler) HandlerGetProfile() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authUserId := ctx.Get(ConstAuthedUserParam).(uint64)

		user, errResp := uh.userUsecase.GetByName(ctx.Param("name"), authUserId)
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := json.Marshal(user)
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}

func (uh *UserHandler) HandlerAddSubscription() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userId := ctx.Get(ConstAuthedUserParam).(uint64)

		if errResp := uh.userUsecase.AddSubscription(userId, ctx.Param("name")); errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		return ctx.JSON(http.StatusOK, OKResponse)
	}
}

func (uh *UserHandler) HandlerRemoveSubscription() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userId := ctx.Get(ConstAuthedUserParam).(uint64)

		if errResp := uh.userUsecase.RemoveSubscription(userId, ctx.Param("name")); errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		return ctx.JSON(http.StatusOK, OKResponse)
	}
}

func (uh *UserHandler) HandlerGetSubscriptions() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authUserId := ctx.Get(ConstAuthedUserParam).(uint64)

		user, errResp := uh.userUsecase.GetByName(ctx.Param("name"), authUserId)
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		subscriptions, errResp := uh.userUsecase.GetSubscriptions(user.ID, authUserId)
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)

		resp, e := json.Marshal(subscriptions)
		if e != nil {
			return RespondWithError(NewErrorResponse(ErrInternal, e), ctx)
		}

		_, e = ctx.Response().Write(resp)
		return e
	}
}
