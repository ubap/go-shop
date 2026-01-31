package db

import (
	"database/sql"
	"fmt"
)

func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 2. Safety Settings
	db.SetMaxOpenConns(1) // The Golden Rule for SQLite

	// Enable WAL mode for performance and crash resilience
	if _, err := db.Exec(`PRAGMA journal_mode = WAL; PRAGMA synchronous = NORMAL;`); err != nil {
		return nil, fmt.Errorf("failed to set pragmas: %w", err)
	}

	// 3. Simple Migration (Create table)
	query := `
	CREATE TABLE IF NOT EXISTS shopping_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		completed BOOLEAN DEFAULT 0
	);`
	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &Store{conn: db}, nil
}
