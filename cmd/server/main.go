package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/tofraley/audity/internal/auditerserver"
	"github.com/tofraley/audity/rpc/auditer"
)

func main() {
	server, err := auditerserver.NewServer("./test.db")
	if err != nil {
		fmt.Printf("Failed to create server: %v\n", err)
		os.Exit(1)
	}

	twirpHandler := auditer.NewAuditerServer(server)
	http.ListenAndServe(":8080", twirpHandler)
}
