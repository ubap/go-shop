package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite" // Import the driver
)

type Basket struct {
	ID        int64     `db:"id"`
	Key       string    `db:"key"`
	CreatedAt time.Time `db:"created_at"`
}

type Item struct {
	ID        int64  `db:"id" json:"id"`
	Title     string `db:"title" json:"title"`
	Completed bool   `db:"completed" json:"completed"`
}

type Store struct {
	conn *sqlx.DB
}

func (s *Store) Close() error {
	return s.conn.Close()
}

// AddItemToBasket links an item (by title) to a specific basket (by key).
func (s *Store) AddItemToBasket(basketKey string, title string) (int64, error) {
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
		INSERT INTO basket_items (basket_key, title)
		VALUES (?, ?)
		ON CONFLICT(basket_key, title COLLATE NOCASE)
		DO UPDATE SET title = excluded.title
		RETURNING id
	`
	var id int64
	err = tx.QueryRow(query, basketKey, title).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("upsert item: %w", err)
	}

	// 5. Finalize
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Store) validateItemTitle(title string) bool {
	if len(title) == 0 {
		return false
	}
	if len(title) > 255 {
		return false
	}
	return true
}

func (s *Store) validateBasketKey(key string) bool {
	parsed, err := uuid.Parse(key)
	if err != nil {
		return false
	}
	// Specifically check that it is UUID v4
	return parsed.Version() == 4
}

// SetItemCompletion updates an item but ONLY if it belongs to the specified basket.
// This ensures that knowing an item ID isn't enough to modify someone else's list.
func (s *Store) SetItemCompletion(basketKey string, id int64, completed bool) error {
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

func (s *Store) GetItemsForBasket(basketKey string) ([]Item, error) {
	if !s.validateBasketKey(basketKey) {
		return nil, fmt.Errorf("invalid basket key")
	}
	// Initialize an empty slice so we return [] instead of null in JSON
	var items []Item

	query := `
		SELECT id, title, completed 
		FROM basket_items 
		WHERE basket_key = ?
		ORDER BY  completed ASC, touched_at DESC, id DESC`

	err := s.conn.Select(&items, query, basketKey)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch items for basket %s: %w", basketKey, err)
	}

	return items, nil
}
