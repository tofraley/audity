package main

import (
	"fmt"
	"net/http"

	"github.com/tofraley/audity/internal/auditerserver"
	"github.com/tofraley/audity/rpc/auditer"
)

func main() {
	server := &auditerserver.Server{}
	twirpHandler := auditer.NewAuditerServer(server)

	fmt.Printf("Starting server\n")
	http.ListenAndServe(":8080", twirpHandler)
}
