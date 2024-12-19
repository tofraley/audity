package auditerdb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

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

	// Insert npm_audits
	auditStmt, err := tx.Prepare(`
		INSERT INTO npm_audits (project_id, audit_report_version)
		VALUES (?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare npm_audits statement: %v", err)
	}
	defer auditStmt.Close()

	res, err := auditStmt.Exec(projectID, req.Result.AuditReportVersion)
	if err != nil {
		return fmt.Errorf("failed to insert npm_audits: %v", err)
	}

	auditID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %v", err)
	}

	// Insert vulnerabilities
	vulnStmt, err := tx.Prepare(`
		INSERT INTO vulnerabilities (
			npm_audit_id, name, severity, is_direct, via, effects, range,
			nodes, fix_name, fix_version, fix_is_sem_ver_major
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare vulnerabilities statement: %v", err)
	}
	defer vulnStmt.Close()

	for _, vuln := range req.Result.Vulnerabilities {
		via, _ := json.Marshal(vuln.Via)
		effects, _ := json.Marshal(vuln.Effects)
		nodes, _ := json.Marshal(vuln.Nodes)

		_, err = vulnStmt.Exec(
			auditID, vuln.Name, vuln.Severity, vuln.IsDirect, string(via),
			string(effects), vuln.Range, string(nodes),
			vuln.FixAvailable.Name, vuln.FixAvailable.Version, vuln.FixAvailable.IsSemVerMajor,
		)
		if err != nil {
			return fmt.Errorf("failed to insert vulnerability: %v", err)
		}
	}

	// Insert vulnerability summary
	vulnSummaryStmt, err := tx.Prepare(`
		INSERT INTO vulnerability_summaries (
			npm_audit_id, info, low, moderate, high, critical, total
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare vulnerability_summaries statement: %v", err)
	}
	defer vulnSummaryStmt.Close()

	_, err = vulnSummaryStmt.Exec(
		auditID, req.Result.Metadata.Vulnerabilities.Info,
		req.Result.Metadata.Vulnerabilities.Low,
		req.Result.Metadata.Vulnerabilities.Moderate,
		req.Result.Metadata.Vulnerabilities.High,
		req.Result.Metadata.Vulnerabilities.Critical,
		req.Result.Metadata.Vulnerabilities.Total,
	)
	if err != nil {
		return fmt.Errorf("failed to insert vulnerability summary: %v", err)
	}

	// Insert dependency summary
	depSummaryStmt, err := tx.Prepare(`
		INSERT INTO dependency_summaries (
			npm_audit_id, prod, dev, optional, peer, peer_optional, total
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare dependency_summaries statement: %v", err)
	}
	defer depSummaryStmt.Close()

	_, err = depSummaryStmt.Exec(
		auditID, req.Result.Metadata.Dependencies.Prod,
		req.Result.Metadata.Dependencies.Dev,
		req.Result.Metadata.Dependencies.Optional,
		req.Result.Metadata.Dependencies.Peer,
		req.Result.Metadata.Dependencies.PeerOptional,
		req.Result.Metadata.Dependencies.Total,
	)
	if err != nil {
		return fmt.Errorf("failed to insert dependency summary: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("Successfully recorded npm audit for project: %s", req.ProjectName)
	return nil
}
