package inmemory

import (
	"go-shop/backend/basket"
	"sort"
	"sync"
)

type Basket struct {
	mu          sync.Mutex
	basketItems map[string]basket.Item
}

func NewBasket() *Basket {
	return &Basket{
		basketItems: make(map[string]basket.Item),
	}
}

func (i *Basket) UpsertItem(basketItem basket.Item) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.basketItems[basketItem.ItemID] = basketItem
}

func (i *Basket) GetAllItems() []basket.Item {
	i.mu.Lock()
	defer i.mu.Unlock()

	basketItemsSlice := make([]basket.Item, 0, len(i.basketItems))

	for _, value := range i.basketItems {
		basketItemsSlice = append(basketItemsSlice, value)
	}
	sort.Slice(basketItemsSlice, func(i, j int) bool {
		return basketItemsSlice[i].LastModified.Before(basketItemsSlice[j].LastModified)
	})
	return basketItemsSlice
}
