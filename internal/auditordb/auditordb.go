package auditerdb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"

	pb "github.com/tofraley/audity/rpc/auditer"
)

type AuditerService struct {
	db *sql.DB
}

func NewAuditerService(dbPath string) (*AuditerService, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	return &AuditerService{db: db}, nil
}

func (s *AuditerService) Close() error {
	return s.db.Close()
}

func (s *AuditerService) RecordNpmAudit(req *pb.NpmAuditRequest) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Insert project if not exists
	var projectID int64
	err = tx.QueryRow("INSERT OR IGNORE INTO projects (name) VALUES (?) RETURNING id", req.ProjectName).Scan(&projectID)
	if err != nil {
		if err != sql.ErrNoRows {
			return fmt.Errorf("failed to insert project: %v", err)
		}
		err = tx.QueryRow("SELECT id FROM projects WHERE name = ?", req.ProjectName).Scan(&projectID)
		if err != nil {
			return fmt.Errorf("failed to get project id: %v", err)
		}
	}

	// Insert npm_audit_results
	resultStmt, err := tx.Prepare(`
		INSERT INTO npm_audit_results (
			project_id, audit_time, total_dependencies, total_dev_dependencies, 
			total_optional_dependencies, total_vulnerabilities
		) VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare npm_audit_results statement: %v", err)
	}
	defer resultStmt.Close()

	res, err := resultStmt.Exec(
		projectID,
		time.Unix(req.Result.AuditTime, 0),
		req.Result.TotalDependencies,
		req.Result.TotalDevDependencies,
		req.Result.TotalOptionalDependencies,
		req.Result.TotalVulnerabilities,
	)
	if err != nil {
		return fmt.Errorf("failed to insert npm_audit_results: %v", err)
	}

	resultID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %v", err)
	}

	// Insert metadata
	metadataStmt, err := tx.Prepare("INSERT INTO metadata (result_id, key, value) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare metadata statement: %v", err)
	}
	defer metadataStmt.Close()

	for key, metadata := range req.Result.Metadata {
		value, err := json.Marshal(metadata.Values)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata values: %v", err)
		}
		_, err = metadataStmt.Exec(resultID, key, string(value))
		if err != nil {
			return fmt.Errorf("failed to insert metadata: %v", err)
		}
	}

	// Insert vulnerabilities
	vulnStmt, err := tx.Prepare(`
		INSERT OR IGNORE INTO vulnerabilities (
			id, url, title, severity, module_name, vulnerable_functions, 
			access, patched_versions, cwe, updated, recommendation
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare vulnerabilities statement: %v", err)
	}
	defer vulnStmt.Close()

	vulnAuditStmt, err := tx.Prepare("INSERT INTO vulnerability_audit_results (vulnerability_id, result_id) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare vulnerability_audit_results statement: %v", err)
	}
	defer vulnAuditStmt.Close()

	for _, vuln := range req.Result.Vulnerabilities {
		_, err = vulnStmt.Exec(
			vuln.Id, vuln.Url, vuln.Title, vuln.Severity, vuln.ModuleName,
			vuln.VulnerableFunctions, vuln.Access, vuln.PatchedVersions,
			vuln.Cwe, vuln.Updated, vuln.Recommendation,
		)
		if err != nil {
			return fmt.Errorf("failed to insert vulnerability: %v", err)
		}

		_, err = vulnAuditStmt.Exec(vuln.Id, resultID)
		if err != nil {
			return fmt.Errorf("failed to insert vulnerability_audit_result: %v", err)
		}

		// Insert vulnerable versions
		for _, version := range vuln.VulnerableVersions {
			_, err = tx.Exec("INSERT OR IGNORE INTO vulnerable_versions (vulnerability_id, version) VALUES (?, ?)", vuln.Id, version)
			if err != nil {
				return fmt.Errorf("failed to insert vulnerable version: %v", err)
			}
		}

		// Insert CVEs
		for _, cve := range vuln.Cves {
			_, err = tx.Exec("INSERT OR IGNORE INTO cves (vulnerability_id, cve) VALUES (?, ?)", vuln.Id, cve)
			if err != nil {
				return fmt.Errorf("failed to insert CVE: %v", err)
			}
		}

		// Insert findings
		for _, finding := range vuln.Findings {
			var findingID int64
			err = tx.QueryRow("INSERT INTO findings (vulnerability_id, version) VALUES (?, ?) RETURNING id", vuln.Id, finding.Version).Scan(&findingID)
			if err != nil {
				return fmt.Errorf("failed to insert finding: %v", err)
			}

			for _, path := range finding.Paths {
				_, err = tx.Exec("INSERT INTO finding_paths (finding_id, path) VALUES (?, ?)", findingID, path)
				if err != nil {
					return fmt.Errorf("failed to insert finding path: %v", err)
				}
			}
		}
	}

	// Insert audit summary
	_, err = tx.Exec(`
		INSERT INTO audit_summaries (result_id, info, low, moderate, high, critical)
		VALUES (?, ?, ?, ?, ?, ?)
	`, resultID, req.Result.Summary.Info, req.Result.Summary.Low, req.Result.Summary.Moderate, req.Result.Summary.High, req.Result.Summary.Critical)
	if err != nil {
		return fmt.Errorf("failed to insert audit summary: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("Successfully recorded npm audit for project: %s", req.ProjectName)
	return nil
}
