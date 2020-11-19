package mwares

import (
	"errors"
	"time"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/user"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/csrf"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/responser"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type MiddlewareManager struct {
	sessionUsecase session.SessionUsecase
	userUsecase    user.UserUsecase
}

func NewMiddlewareManager(sessionUsecase session.SessionUsecase, userUsecase user.UserUsecase) *MiddlewareManager {
	return &MiddlewareManager{
		sessionUsecase: sessionUsecase,
		userUsecase:    userUsecase,
	}
}

func (mm *MiddlewareManager) PanicRecovering(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				logrus.Warn(err)
			}
		}()

		return next(ctx)
	}
}

func (mm *MiddlewareManager) AccessLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		logrus.Info(ctx.Request().RemoteAddr, " ", ctx.Request().Method, " ", ctx.Request().URL)

		start := time.Now()
		err := next(ctx)
		end := time.Now()

		logrus.Info("Status: ", ctx.Response().Status, " Work time: ", end.Sub(start))
		logrus.Println()

		return err
	}
}

func (mm *MiddlewareManager) CORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowHeaders:     ConstAllowedHeaders,
		AllowMethods:     ConstAllowedMethods,
		AllowOrigins:     ConstAllowedOrigins,
		ExposeHeaders:    ConstAllowedExpose,
	})
}

func (mm *MiddlewareManager) CheckCSRF(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		csrfErr := errors.New("Bad csrf token received")
		sessionCookie, err := ctx.Cookie(ConstSessionName)

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrNotAuthorized, err), ctx)
		}

		session, errResp := mm.sessionUsecase.GetByID(sessionCookie.Value)
		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		token := ctx.Request().Header.Get(ConstCSRFTokenName)

		if token == "" {
			return RespondWithError(NewErrorResponse(ErrBadRequest, csrfErr), ctx)
		}

		errResp = csrf.ValidateCSRFToken(session, token)

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		return next(ctx)
	}
}

func (mm *MiddlewareManager) XSS() echo.MiddlewareFunc {
	return middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection: "1; mode=block",
	})
}

func (mm *MiddlewareManager) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cookie, err := ctx.Cookie(ConstSessionName)

		if err != nil {
			return RespondWithError(NewErrorResponse(ErrNotAuthorized, err), ctx)
		}

		session, errResp := mm.sessionUsecase.GetByID(cookie.Value)

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		_, errResp = mm.userUsecase.GetById(session.UserID)

		if errResp != nil {
			return RespondWithError(errResp, ctx)
		}

		ctx.Set(ConstAuthedUserParam, session.UserID)

		return next(ctx)
	}
}
