package main

// global variable
var basketItems = make(map[string]BasketItem)

func UpdateItem(basketItem BasketItem) {
	basketItems[basketItem.ItemID] = basketItem
}

func GetAllItems() []BasketItem {
	basketItemsSlice := make([]BasketItem, 0, len(basketItems))

	// 2. The loop remains the same.
	for _, value := range basketItems {
		basketItemsSlice = append(basketItemsSlice, value)
	}
	return basketItemsSlice
}
