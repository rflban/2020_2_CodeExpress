package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/admin/grpc_admin"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/admin/proto_admin"

	"google.golang.org/grpc"

	"github.com/sirupsen/logrus"

	artistRepository "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/artist/repository"

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

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.AdminMicroservice.Port))
	if err != nil {
		logrus.Fatal(err)
	}

	if err := dbConn.Ping(); err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("DB connected on %s", conf.GetDbConnString())

	artistRep := artistRepository.NewArtistRep(dbConn)

	adminServer := grpc.NewServer()
	proto_admin.RegisterAdminServiceServer(adminServer, grpc_admin.NewAdminGRPCUsecase(artistRep))
	logrus.Fatal(adminServer.Serve(lis))
}
