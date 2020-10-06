package main

import (
	"github.com/go-park-mail-ru/2020_2_CodeExpress/server"
	"os"
	"fmt"
)

func main() {
	if len(os.Args) <  3 {
		fmt.Println("Usage: go run main.go $url $port")
		return
	}

	url := os.Args[1]
	port := os.Args[2]

	server.ServerStart(url, port)
}
