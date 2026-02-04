package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_AddItem(t *testing.T) {
	store, _ := NewStore(":memory:")

	basketKey := "sample-basket"

	store.AddItemToBasket(basketKey, "Item 1")
	store.AddItemToBasket(basketKey, "Item 2")

	items, err := store.GetItemsForBasket(basketKey)
	require.NoError(t, err)
	require.Len(t, items, 2)

}
