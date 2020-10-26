package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/config"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/server"
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

	config, err := config.LoadConfig(configFileName)
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := sql.Open(dbName, config.GetDbConnString())
	if err != nil {
		log.Fatal(err)
	}

	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Printf("DB connected on %s", config.GetDbConnString())

	server.ServerStart(config.GetServerConnString(), dbConn)
}
