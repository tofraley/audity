package main

import (
	"fmt"
	"net/http"

	"github.com/tofraley/audity/internal/haberdasherserver"
	"github.com/tofraley/audity/rpc/haberdasher"
)

func main() {
	server := &haberdasherserver.Server{} // implements Haberdasher interface
	twirpHandler := haberdasher.NewHaberdasherServer(server)

	fmt.Printf("Starting server\n")
	http.ListenAndServe(":8080", twirpHandler)
}
