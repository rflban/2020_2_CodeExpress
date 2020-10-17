package server

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/handlers"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/repositories"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
)

func ServerStart(addr string, dbConn *sql.DB) {
	UserRep := repositories.NewUserRepPGImpl(dbConn)
	SessionRep := repositories.NewSessionRepPGImpl(dbConn)

	UserHandler := handlers.NewUserHandler(UserRep, SessionRep)

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Access-Control-Allow-Headers", "Authorization", "X-Requested-With"},
		AllowMethods:     []string{http.MethodDelete, http.MethodGet, http.MethodOptions, http.MethodPost},
		AllowOrigins:     []string{"http://musicexpress.sarafa2n.ru"},
	}))

	e.POST("/api/v1/user/register", UserHandler.HandleCreateUser)

	e.POST("/api/v1/user/login", UserHandler.HandleLogInUser)

	e.DELETE("/api/v1/user/logout", UserHandler.HandleLogOutUser)

	e.POST("/api/v1/user/change/profile", UserHandler.HandleUpdateProfile)

	e.POST("/api/v1/user/change/password", UserHandler.HandleUpdatePassword)

	e.POST("/api/v1/user/change/avatar", UserHandler.HandleUpdateAvatar)

	e.GET("/api/v1/user/current", UserHandler.HandleCurrentUser)

	e.Static("/avatars", "avatars")

	log.Println("Server listening on ", addr)
	e.Logger.Fatal(e.Start(addr))
}
