package basket

import (
	"time"
)

type Item struct {
	ItemID       string    `json:"id"`
	Name         string    `json:"name"`
	LastModified time.Time `json:"lastModified"`
	ToBuy        bool      `json:"toBuy"`
}

type Basket interface {
	UpsertItem(basketItem Item)
	GetAllItems() []Item
}
