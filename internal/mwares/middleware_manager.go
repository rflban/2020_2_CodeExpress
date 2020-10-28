package mwares

import (
	"time"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type MiddlewareManager struct{}

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
	})
}
