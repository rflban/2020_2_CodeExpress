package main

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/grpc_session"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/proto_session"
	"net"
	"os"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	sessionRepository "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/session/repository"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/config"

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

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.SessionMicroservice.Port))
	if err != nil {
		logrus.Fatal(err)
	}

	if err := dbConn.Ping(); err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("DB connected on %s", conf.GetDbConnString())

	sessionRep := sessionRepository.NewSessionRep(dbConn)

	sessionServer := grpc.NewServer()
	proto_session.RegisterSessionServiceServer(sessionServer, grpc_session.NewSessionGRPCUsecase(sessionRep))

	logrus.Fatal(sessionServer.Serve(lis))
}
