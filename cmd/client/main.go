package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/tofraley/audity/rpc/auditer"
)

func main() {
	// Load audit_results.json
	auditResultsBytes, err := ioutil.ReadFile("audit_results-long.json")
	if err != nil {
		fmt.Printf("Error reading audit_results.json: %v\n", err)
		os.Exit(1)
	}

	var npmAuditResult auditer.NpmAuditResult
	err = json.Unmarshal(auditResultsBytes, &npmAuditResult)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		os.Exit(1)
	}

	client := auditer.NewAuditerProtobufClient("http://localhost:8080", &http.Client{})

	result, err := client.RecordNpmAudit(context.Background(), &auditer.NpmAuditRequest{
		ProjectName: "TestProject",
		Result:      &npmAuditResult,
	})
	if err != nil {
		fmt.Printf("oh no: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Result: %+v", result)
}
