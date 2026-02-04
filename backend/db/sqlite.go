package db

import (
	"context"
	"database/sql"
	"errors"
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
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type Store struct {
	conn *sqlx.DB
}

func (s *Store) Close() error {
	return s.conn.Close()
}

// AddItemToBasket links an item (by title) to a specific basket (by key).
// Creates the item if it doesn't exist, and
// ensures the item is linked to the basket exactly once.
func (s *Store) AddItemToBasket(basketKey string, title string) error {
	tx, err := s.conn.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	basketID, err := s.ensureBasket(tx, basketKey)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT OR IGNORE INTO items (title) VALUES (?)", title)
	if err != nil {
		return err
	}

	var itemID int64
	err = tx.Get(&itemID, "SELECT id FROM items WHERE title = ?", title)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT OR IGNORE INTO basket_items (basket_id, item_id) VALUES (?, ?)", basketID, itemID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// ensureBasket handles the JIT creation and returns the internal ID.
func (s *Store) ensureBasket(tx *sqlx.Tx, key string) (int64, error) {
	_, err := tx.Exec("INSERT OR IGNORE INTO baskets (key) VALUES (?)", key)
	if err != nil {
		return 0, fmt.Errorf("auto-create basket: %w", err)
	}

	var id int64
	err = tx.Get(&id, "SELECT id FROM baskets WHERE key = ?", key)
	if err != nil {
		return 0, fmt.Errorf("lookup basket id: %w", err)
	}

	return id, nil
}

func (s *Store) GetItemsForBasket(basketKey string) ([]Item, error) {
	var items []Item
	query := `
	SELECT i.id, i.title, i.completed
	FROM items i
	JOIN basket_items bi ON i.id = bi.item_id
	JOIN baskets b ON bi.basket_id = b.id
	WHERE b.key = ?;
	`
	err := s.conn.Select(&items, query, basketKey)
	if err != nil {
		return nil, fmt.Errorf("get items for basket: %w", err)
	}
	return items, nil
}

// CreateBasketIfNotExists ensures a basket with the given key exists in the database.
// It returns true if a new basket was created, or false if the basket already existed.
// If an error occurs during the database operation, it returns false and the error.
func (s *Store) CreateBasketIfNotExists(key string) (bool, error) {
	res, err := s.conn.Exec("INSERT OR IGNORE INTO baskets (key) VALUES (?)", key)
	if err != nil {
		return false, fmt.Errorf("create basket: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		// Some drivers don't support RowsAffected, but SQLite does.
		return false, fmt.Errorf("rows affected: %w", err)
	}

	return count > 0, nil
}

func (s *Store) GetBasket(key string) (Basket, error) {
	var b Basket
	err := s.conn.Get(&b, "SELECT * FROM baskets WHERE key = ?", key)
	if errors.Is(err, sql.ErrNoRows) {
		return Basket{}, nil
	}
	return b, err
}
