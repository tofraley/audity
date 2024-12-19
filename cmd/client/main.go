package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/tofraley/audity/rpc/auditer"
)

func main() {
	// Load audit_results.json
	auditResultsBytes, err := ioutil.ReadFile("audit_results.json")
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

	fmt.Printf("audit: %+v\n", &npmAuditResult)

	vulnerabilities := npmAuditResult.GetVulnerabilities()

	// Get the first value
	var firstVulnerability *auditer.Vulnerability
	for _, v := range vulnerabilities {
		firstVulnerability = v
		break
	}

	if firstVulnerability != nil {
		fmt.Printf("First fix: %+v\n", firstVulnerability.FixAvailable)
	} else {
		fmt.Println("No vulnerabilities found")
	}

	// client := auditer.NewAuditerProtobufClient("http://localhost:8080", &http.Client{})

	// result, err := client.RecordNpmAudit(context.Background(), &auditer.NpmAuditRequest{
	// 	ProjectName: "TestProject",
	// 	Result:      &npmAuditResult,
	// })
	// if err != nil {
	// 	fmt.Printf("oh no: %v", err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("Result: %+v", result)
}
