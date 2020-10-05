package main

import (
	"github.com/go-park-mail-ru/2020_2_CodeExpress/server"
)

func main() {
	url := "127.0.0.1"
	port := ":8080"
	server.ServerStart(url, port)
}
