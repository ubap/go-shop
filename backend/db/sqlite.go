package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite" // Import the driver
)

type Basket struct {
	ID       int64     `db:"id"`
	Key      string    `db:"key"`
	CreatdAt time.Time `db:"created_at"`
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
func (s *Store) AddItem(title string) (int64, error) {
	res, err := s.conn.Exec("INSERT INTO shopping_items (title) VALUES (?)", title)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) GetAllItems() ([]Item, error) {
	rows, err := s.conn.Query("SELECT id, title, completed FROM shopping_items ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var i Item
		if err := rows.Scan(&i.ID, &i.Title, &i.Completed); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func (s *Store) ToggleItem(id int64, completed bool) error {
	_, err := s.conn.Exec("UPDATE shopping_items SET completed = ? WHERE id = ?", completed, id)
	return err
}

func (s *Store) DeleteItem(id int64) error {
	_, err := s.conn.Exec("DELETE FROM shopping_items WHERE id = ?", id)
	return err
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
