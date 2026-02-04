package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CreateBasketIfNotExists(t *testing.T) {
	store, err := NewStore(":memory:")
	require.NoError(t, err)

	created, err := store.CreateBasketIfNotExists("test-basket")
	require.NoError(t, err)
	require.True(t, created)
	basket, err := store.GetBasket("test-basket")
	require.NoError(t, err)
	require.Equal(t, "test-basket", basket.Key)

	created, err = store.CreateBasketIfNotExists("test-basket")
	require.NoError(t, err)
	require.False(t, created)
}

func Test_AddItem(t *testing.T) {
	store, _ := NewStore(":memory:")

	store.CreateBasketIfNotExists("test-basket")

}

func TestStore_AddItemToBasket(t *testing.T) {
	store, _ := NewStore(":memory:")

	store.AddItemToBasket("test-basket", "Test Item")

	basket, err := store.GetItemsForBasket("test-basket")
	require.NoError(t, err)
	require.Len(t, basket, 1)

}
