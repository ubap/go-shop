package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite" // Import the driver
)

type SqliteStore struct {
	conn *sqlx.DB
}

var _ Store = (*SqliteStore)(nil)

func (s *SqliteStore) Close() error {
	return s.conn.Close()
}

func (s *SqliteStore) AddItemToBasket(basketKey string, title string) (int64, error) {
	if !s.validateItemTitle(title) {
		return 0, fmt.Errorf("invalid item title")
	}
	if !s.validateBasketKey(basketKey) {
		return 0, fmt.Errorf("invalid basket key")
	}

	tx, err := s.conn.BeginTxx(context.Background(), nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT OR IGNORE INTO baskets (key) VALUES (?)", basketKey)
	if err != nil {
		return 0, err
	}

	query := `
		INSERT INTO basket_items (basket_key, title, status, completed)
		VALUES (?, ?, ?, 0)
		ON CONFLICT(basket_key, title COLLATE NOCASE)
		DO UPDATE SET
		              title = excluded.title,
		              completed = excluded.completed,
					  status = excluded.status
		RETURNING id
	`
	var id int64
	err = tx.QueryRow(query, basketKey, title, StatusActive).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("upsert item: %w", err)
	}

	// 5. Finalize
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *SqliteStore) SetItemCompletion(basketKey string, id int64, completed bool) error {
	if !s.validateBasketKey(basketKey) {
		return fmt.Errorf("invalid basket key")
	}

	query := `
		UPDATE basket_items 
		SET completed = ?
		WHERE id = ? AND basket_key = ?
	`
	res, err := s.conn.Exec(query, completed, id, basketKey)
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		// Either the ID doesn't exist, OR it exists but belongs to a
		// different basket_key. We return the same error for both
		// to avoid "leaking" information about other users' IDs.
		return fmt.Errorf("item not found in this basket")
	}

	return nil
}

func (s *SqliteStore) DeleteItem(basketKey string, itemId int64) error {
	if !s.validateBasketKey(basketKey) {
		return fmt.Errorf("invalid basket key")
	}

	query := `
        UPDATE basket_items 
        SET status = ? 
        WHERE basket_key = ? AND id = ? AND status != ?
    `

	res, err := s.conn.Exec(query, StatusDeleted, basketKey, itemId, StatusDeleted)
	if err != nil {
		return fmt.Errorf("soft delete item: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("item not found or already deleted")
	}

	return nil
}

func (s *SqliteStore) RestoreItem(basketKey string, itemId int64) error {
	if !s.validateBasketKey(basketKey) {
		return fmt.Errorf("invalid basket key")
	}

	query := `
        UPDATE basket_items 
        SET status = ? 
        WHERE basket_key = ? AND id = ? AND status != ?
    `

	res, err := s.conn.Exec(query, StatusActive, basketKey, itemId, StatusActive)
	if err != nil {
		return fmt.Errorf("undelete item: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("item not found or already active")
	}

	return nil
}

func (s *SqliteStore) GetItemsForBasket(basketKey string) ([]Item, error) {
	if !s.validateBasketKey(basketKey) {
		return nil, fmt.Errorf("invalid basket key")
	}
	// Initialize an empty slice so we return [] instead of null in JSON
	var items []Item

	query := `
		SELECT id, title, completed, status
		FROM basket_items 
		WHERE basket_key = ? and status = ?
		ORDER BY  completed ASC, touched_at DESC, id DESC`

	err := s.conn.Select(&items, query, basketKey, StatusActive)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch items for basket %s: %w", basketKey, err)
	}

	return items, nil
}

func (s *SqliteStore) validateItemTitle(title string) bool {
	if len(title) == 0 {
		return false
	}
	if len(title) > 255 {
		return false
	}
	return true
}

func (s *SqliteStore) validateBasketKey(key string) bool {
	parsed, err := uuid.Parse(key)
	if err != nil {
		return false
	}
	// Specifically check that it is UUID v4
	return parsed.Version() == 4
}
