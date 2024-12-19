-- +goose Up
-- +goose StatementBegin
-- Projects table
CREATE TABLE projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

-- NpmAuditResults table
CREATE TABLE npm_audit_results (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    audit_time INTEGER NOT NULL,
    total_dependencies INTEGER NOT NULL,
    total_dev_dependencies INTEGER NOT NULL,
    total_optional_dependencies INTEGER NOT NULL,
    total_vulnerabilities INTEGER NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

-- Metadata table (for the metadata map in NpmAuditResult)
CREATE TABLE metadata (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    result_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    FOREIGN KEY (result_id) REFERENCES npm_audit_results(id)
);

-- Vulnerabilities table
CREATE TABLE vulnerabilities (
    id TEXT PRIMARY KEY,
    url TEXT,
    title TEXT NOT NULL,
    severity TEXT NOT NULL,
    module_name TEXT NOT NULL,
    vulnerable_functions TEXT,
    access TEXT,
    patched_versions TEXT,
    cwe TEXT,
    updated TEXT,
    recommendation TEXT
);

-- VulnerabilityAuditResults table (to correlate vulnerabilities with audit results)
CREATE TABLE vulnerability_audit_results (
    vulnerability_id TEXT NOT NULL,
    result_id INTEGER NOT NULL,
    PRIMARY KEY (vulnerability_id, result_id),
    FOREIGN KEY (vulnerability_id) REFERENCES vulnerabilities(id),
    FOREIGN KEY (result_id) REFERENCES npm_audit_results(id)
);

-- VulnerableVersions table
CREATE TABLE vulnerable_versions (
    vulnerability_id TEXT NOT NULL,
    version TEXT NOT NULL,
    PRIMARY KEY (vulnerability_id, version),
    FOREIGN KEY (vulnerability_id) REFERENCES vulnerabilities(id)
);

-- CVEs table
CREATE TABLE cves (
    vulnerability_id TEXT NOT NULL,
    cve TEXT NOT NULL,
    PRIMARY KEY (vulnerability_id, cve),
    FOREIGN KEY (vulnerability_id) REFERENCES vulnerabilities(id)
);

-- Findings table
CREATE TABLE findings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    vulnerability_id TEXT NOT NULL,
    version TEXT NOT NULL,
    FOREIGN KEY (vulnerability_id) REFERENCES vulnerabilities(id)
);

-- FindingPaths table
CREATE TABLE finding_paths (
    finding_id INTEGER NOT NULL,
    path TEXT NOT NULL,
    PRIMARY KEY (finding_id, path),
    FOREIGN KEY (finding_id) REFERENCES findings(id)
);

-- AuditSummary table
CREATE TABLE audit_summaries (
    result_id INTEGER PRIMARY KEY,
    info INTEGER NOT NULL,
    low INTEGER NOT NULL,
    moderate INTEGER NOT NULL,
    high INTEGER NOT NULL,
    critical INTEGER NOT NULL,
    FOREIGN KEY (result_id) REFERENCES npm_audit_results(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS audit_summaries;
DROP TABLE IF EXISTS finding_paths;
DROP TABLE IF EXISTS findings;
DROP TABLE IF EXISTS cves;
DROP TABLE IF EXISTS vulnerable_versions;
DROP TABLE IF EXISTS vulnerability_audit_results;
DROP TABLE IF EXISTS vulnerabilities;
DROP TABLE IF EXISTS metadata;
DROP TABLE IF EXISTS npm_audit_results;
DROP TABLE IF EXISTS projects;
-- +goose StatementEnd
