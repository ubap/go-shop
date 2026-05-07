package db

import (
	"embed"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

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

	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, fmt.Errorf("failed to set goose dialect: %w", err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.Up(db.DB, "migrations"); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &SqliteStore{conn: db}, nil
}
