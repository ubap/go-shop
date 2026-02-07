package db

import "time"

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

type Store interface {
	// AddItemToBasket links an item (by title) to a specific basket (by key).
	AddItemToBasket(basketKey string, title string) (int64, error)
	// SetItemCompletion updates an item but ONLY if it belongs to the specified basket.
	// This ensures that knowing an item ID isn't enough to modify someone else's list.
	SetItemCompletion(basketKey string, id int64, completed bool) error
	// GetItemsForBasket retrieves all items for a given basket, sorted by completion status and recency.
	GetItemsForBasket(basketKey string) ([]Item, error)

	Close() error
}
