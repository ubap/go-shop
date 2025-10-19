package basket

import (
	"sync"
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

type InMemoryBasket struct {
	mu          sync.Mutex
	basketItems map[string]Item
}

func NewInMemoryBasket() *InMemoryBasket {
	return &InMemoryBasket{
		basketItems: make(map[string]Item),
	}
}

func (i *InMemoryBasket) UpsertItem(basketItem Item) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.basketItems[basketItem.ItemID] = basketItem
}

func (i *InMemoryBasket) GetAllItems() []Item {
	i.mu.Lock()
	defer i.mu.Unlock()

	basketItemsSlice := make([]Item, 0, len(i.basketItems))

	for _, value := range i.basketItems {
		basketItemsSlice = append(basketItemsSlice, value)
	}
	return basketItemsSlice
}
