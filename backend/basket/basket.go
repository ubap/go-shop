package basket

import "sync"

type BasketItem struct {
	ItemID       string `json:"id"`
	Name         string `json:"name"`
	LastModified string `json:"lastModified"`
}

type Basket interface {
	UpsertItem(basketItem BasketItem)
	GetAllItems() []BasketItem
}

type InMemoryBasket struct {
	mu          sync.Mutex
	basketItems map[string]BasketItem
}

func NewInMemoryBasket() *InMemoryBasket {
	return &InMemoryBasket{
		basketItems: make(map[string]BasketItem),
	}
}

func (i *InMemoryBasket) UpsertItem(basketItem BasketItem) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.basketItems[basketItem.ItemID] = basketItem
}

func (i *InMemoryBasket) GetAllItems() []BasketItem {
	i.mu.Lock()
	defer i.mu.Unlock()

	basketItemsSlice := make([]BasketItem, 0, len(i.basketItems))

	for _, value := range i.basketItems {
		basketItemsSlice = append(basketItemsSlice, value)
	}
	return basketItemsSlice
}

// global variable
var basketItems = make(map[string]BasketItem)

func UpdateItem(basketItem BasketItem) {
	basketItems[basketItem.ItemID] = basketItem
}

func GetAllItems() []BasketItem {
	basketItemsSlice := make([]BasketItem, 0, len(basketItems))

	for _, value := range basketItems {
		basketItemsSlice = append(basketItemsSlice, value)
	}
	return basketItemsSlice
}
