package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"

	userDelivery "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/user/delivery"
	userRepository "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/user/repository"
	userUsecase "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/user/usecase"

	sessionDelivery "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/delivery"
	sessionRepository "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/repository"
	sessionUsecase "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/usecase"

	artistDelivery "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/artist/delivery"
	artistRepository "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/artist/repository"
	artistUsecase "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/artist/usecase"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/config"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/mwares"
	"github.com/labstack/echo/v4"

	_ "github.com/lib/pq"
)

const (
	argsCount = 2
	argsUsage = "Usage: go run main.go $config_file"
	dbName    = "postgres"
)

func main() {
	if len(os.Args) != argsCount {
		fmt.Println(argsUsage)
		return
	}

	configFileName := os.Args[1]

	conf, err := config.LoadConfig(configFileName)
	if err != nil {
		logrus.Fatal(err)
	}

	dbConn, err := sql.Open(dbName, conf.GetDbConnString())
	if err != nil {
		logrus.Fatal(err)
	}

	defer func() {
		_ = dbConn.Close()
	}()

	if err := dbConn.Ping(); err != nil {
		logrus.Fatal(err)
	}
	log.Printf("DB connected on %s", conf.GetDbConnString())

	e := echo.New()
	mm := &mwares.MiddlewareManager{}

	e.Use(mm.AccessLog, mm.PanicRecovering, mm.CORS()) //TODO: Сделать нормальный CORS
	e.Static("/avatars", "avatars")
	e.Static("/artist_posters", "artist_posters")

	userRep := userRepository.NewUserRep(dbConn)
	sessionRep := sessionRepository.NewSessionRep(dbConn)
	artistRep := artistRepository.NewArtistRep(dbConn)

	userUsecase := userUsecase.NewUserUsecase(userRep)
	sessionUsecase := sessionUsecase.NewSessionUsecase(sessionRep)
	artistUsecase := artistUsecase.NewArtistUsecase(artistRep)

	userHandler := userDelivery.NewUserHandler(userUsecase, sessionUsecase)
	sessionHandler := sessionDelivery.NewSessionHandler(sessionUsecase, userUsecase)
	artistHandler := artistDelivery.NewArtistHandler(artistUsecase)

	userHandler.Configure(e)
	sessionHandler.Configure(e)
	artistHandler.Configure(e)

	e.Logger.Fatal(e.Start(conf.GetServerConnString()))
}
