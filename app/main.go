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

	trackDelivery "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track/delivery"
	trackRepository "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track/repository"
	trackUsecase "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track/usecase"

	albumDelivery "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/album/delivery"
	albumRepository "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/album/repository"
	albumUsecase "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/album/usecase"

	playlistDelivery "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/playlist/delivery"
	playlistRepository "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/playlist/repository"
	playlistUsecase "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/playlist/usecase"

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

	userRep := userRepository.NewUserRep(dbConn)
	sessionRep := sessionRepository.NewSessionRep(dbConn)
	artistRep := artistRepository.NewArtistRep(dbConn)
	trackRep := trackRepository.NewTrackRep(dbConn)
	albumRep := albumRepository.NewAlbumRep(dbConn)
	playlistRep := playlistRepository.NewPlaylistRep(dbConn)

	userUsecase := userUsecase.NewUserUsecase(userRep)
	sessionUsecase := sessionUsecase.NewSessionUsecase(sessionRep)
	artistUsecase := artistUsecase.NewArtistUsecase(artistRep)
	trackUsecase := trackUsecase.NewTrackUsecase(trackRep)
	albumUsecase := albumUsecase.NewAlbumUsecase(albumRep)
	playlistUsecase := playlistUsecase.NewPlaylistUsecase(playlistRep)

	userHandler := userDelivery.NewUserHandler(userUsecase, sessionUsecase)
	sessionHandler := sessionDelivery.NewSessionHandler(sessionUsecase, userUsecase)
	artistHandler := artistDelivery.NewArtistHandler(artistUsecase)
	albumHandler := albumDelivery.NewAlbumHandler(albumUsecase, artistUsecase, trackUsecase)
	trackHandler := trackDelivery.NewTrackHandler(trackUsecase, sessionUsecase, userUsecase)
	playlistHandler := playlistDelivery.NewPlaylistHandler(playlistUsecase, trackUsecase)

	e := echo.New()
	mm := mwares.NewMiddlewareManager(sessionUsecase, userUsecase)

	e.Use(mm.AccessLog, mm.PanicRecovering, mm.CORS(), mm.XSS())
	e.Static("/avatars", "avatars")
	e.Static("/artist_posters", "artist_posters")
	e.Static("/track_audio", "track_audio")
	e.Static("/album_posters", "album_posters")
	e.Static("/artist_avatars", "artist_avatars")
	e.Static("/playlist_posters", "playlist_posters")

	userHandler.Configure(e, mm)
	sessionHandler.Configure(e)
	artistHandler.Configure(e, mm)
	trackHandler.Configure(e, mm)
	albumHandler.Configure(e, mm)
	playlistHandler.Configure(e, mm)

	e.Logger.Fatal(e.Start(conf.GetServerConnString()))
}
