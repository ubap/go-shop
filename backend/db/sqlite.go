package db

import (
	"context"
	"fmt"
	"time"

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
func (s *Store) AddItemToBasket(basketKey string, title string) error {
	// TODO: Input validation
	tx, err := s.conn.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT OR IGNORE INTO baskets (key) VALUES (?)", basketKey)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT OR IGNORE INTO basket_items (basket_key, title) VALUES (?, ?)", basketKey, title)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Store) GetItemsForBasket(basketKey string) ([]Item, error) {
	// TODO: Input validation
	// Initialize an empty slice so we return [] instead of null in JSON
	var items []Item

	query := `
		SELECT id, title, completed 
		FROM basket_items 
		WHERE basket_key = ? 
		ORDER BY id ASC`

	err := s.conn.Select(&items, query, basketKey)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch items for basket %s: %w", basketKey, err)
	}

	return items, nil
}
