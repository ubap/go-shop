package db

import (
	"database/sql"

	_ "modernc.org/sqlite" // Import the driver
)

type Item struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type Store struct {
	conn *sql.DB
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
