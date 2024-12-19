package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/tofraley/audity/rpc/auditer"
)

func main() {
	client := auditer.NewAuditerProtobufClient("http://localhost:8080", &http.Client{})

	result, err := client.RecordNpmAudit(context.Background(), &auditer.NpmAuditRequest{})
	if err != nil {
		fmt.Printf("oh no: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Result: %+v", result)
}
