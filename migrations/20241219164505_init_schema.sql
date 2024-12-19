-- +goose Up
-- +goose StatementBegin

CREATE TABLE projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE npm_audits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    audit_report_version INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

CREATE TABLE vulnerabilities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    npm_audit_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    severity TEXT NOT NULL,
    is_direct BOOLEAN NOT NULL,
    via TEXT,
    effects TEXT,
    range TEXT,
    nodes TEXT,
    fix_name TEXT,
    fix_version TEXT,
    fix_is_sem_ver_major BOOLEAN,
    FOREIGN KEY (npm_audit_id) REFERENCES npm_audits(id)
);

CREATE TABLE vulnerability_summaries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    npm_audit_id INTEGER NOT NULL,
    info INTEGER NOT NULL,
    low INTEGER NOT NULL,
    moderate INTEGER NOT NULL,
    high INTEGER NOT NULL,
    critical INTEGER NOT NULL,
    total INTEGER NOT NULL,
    FOREIGN KEY (npm_audit_id) REFERENCES npm_audits(id)
);

CREATE TABLE dependency_summaries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    npm_audit_id INTEGER NOT NULL,
    prod INTEGER NOT NULL,
    dev INTEGER NOT NULL,
    optional INTEGER NOT NULL,
    peer INTEGER NOT NULL,
    peer_optional INTEGER NOT NULL,
    total INTEGER NOT NULL,
    FOREIGN KEY (npm_audit_id) REFERENCES npm_audits(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE dependency_summaries;
DROP TABLE vulnerability_summaries;
DROP TABLE vulnerabilities;
DROP TABLE npm_audits;
DROP TABLE projects;

-- +goose StatementEnd