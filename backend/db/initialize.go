package db

import (
	_ "embed"
	"fmt"

	"github.com/jmoiron/sqlx"
)

//go:embed schema.sql
var schemaSQL string

func NewSqliteStore(dbPath string) (*SqliteStore, error) {
	db, err := sqlx.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// second+ connection will wait
	db.SetMaxOpenConns(1)

	// Enable WAL mode for performance and crash resilience
	if _, err := db.Exec(`PRAGMA journal_mode = WAL; PRAGMA synchronous = NORMAL;`); err != nil {
		return nil, fmt.Errorf("failed to set pragmas: %w", err)
	}

	if _, err := db.Exec(schemaSQL); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return &SqliteStore{conn: db}, nil
}
