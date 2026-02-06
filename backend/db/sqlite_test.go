package db

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// region Setup
type TestStore struct {
	*Store
	t *testing.T
}

func (ts *TestStore) AddItemToBasket(basketKey, title string) int64 {
	id, err := ts.Store.AddItemToBasket(basketKey, title)
	require.NoError(ts.t, err)
	time.Sleep(1001 * time.Microsecond)
	return id
}

func (ts *TestStore) SetItemCompletion(basketKey string, id int64, completed bool) {
	err := ts.Store.SetItemCompletion(basketKey, id, completed)
	require.NoError(ts.t, err)
	time.Sleep(1001 * time.Microsecond)
}

func setup(t *testing.T) *TestStore {
	s, err := NewStore(":memory:")
	require.NoError(t, err)
	return &TestStore{Store: s, t: t}
}

// endregion

func titles(items []Item) []string {
	res := make([]string, len(items))
	for i, itm := range items {
		res[i] = itm.Title
	}
	return res
}

const basketKey = "sample-basket"

func Test_Sorting(t *testing.T) {
	t.Run("New items appear at the top (touched_at DESC)", func(t *testing.T) {
		store := setup(t)

		store.AddItemToBasket(basketKey, "Item 1")
		store.AddItemToBasket(basketKey, "Item 2")

		items, _ := store.GetItemsForBasket(basketKey)
		require.Equal(t, []string{"Item 2", "Item 1"}, titles(items))
	})

	t.Run("Re-adding an item moves it to the top (touching)", func(t *testing.T) {
		store := setup(t)

		store.AddItemToBasket(basketKey, "Item 1")
		store.AddItemToBasket(basketKey, "Item 2")
		store.AddItemToBasket(basketKey, "Item 1") // Touched Item 1

		items, _ := store.GetItemsForBasket(basketKey)
		require.Equal(t, []string{"Item 1", "Item 2"}, titles(items))
	})

	t.Run("Updating completion moves item to the top", func(t *testing.T) {
		store := setup(t)

		id1 := store.AddItemToBasket(basketKey, "Item 1")
		store.AddItemToBasket(basketKey, "Item 2")
		id3 := store.AddItemToBasket(basketKey, "Item 3")
		id4 := store.AddItemToBasket(basketKey, "Item 4")

		store.SetItemCompletion(basketKey, id1, false)
		store.SetItemCompletion(basketKey, id4, true)

		items, _ := store.GetItemsForBasket(basketKey)
		require.Equal(t, []string{"Item 1", "Item 3", "Item 2", "Item 4"}, titles(items))

		store.SetItemCompletion(basketKey, id3, true)
		items, _ = store.GetItemsForBasket(basketKey)
		require.Equal(t, []string{"Item 1", "Item 2", "Item 3", "Item 4"}, titles(items))
	})
}

func Test_ValidateInput(t *testing.T) {
	t.Run("Adding item with empty title fails", func(t *testing.T) {
		store := setup(t)

		_, err := store.Store.AddItemToBasket(basketKey, "")
		require.Error(t, err)
	})

	t.Run("Too long item title fails", func(t *testing.T) {
		store := setup(t)

		longTitle := strings.Repeat("a", 256)

		_, err := store.Store.AddItemToBasket(basketKey, longTitle)
		require.Error(t, err)
	})
}

func Test_IsolatedBaskets(t *testing.T) {
	t.Run("Items in different baskets do not interfere", func(t *testing.T) {
		store := setup(t)

		id1 := store.AddItemToBasket("basket-1", "Item A")
		store.AddItemToBasket("basket-2", "Item B")

		// Completing item in basket-1 should not affect basket-2
		err := store.Store.SetItemCompletion("basket-1", id1, true)
		require.NoError(t, err)

		items1, _ := store.GetItemsForBasket("basket-1")
		require.Equal(t, []string{"Item A"}, titles(items1))
		require.True(t, items1[0].Completed)
	})

	t.Run("Completing item with wrong basket key fails", func(t *testing.T) {
		store := setup(t)

		id1 := store.AddItemToBasket("basket-1", "Item A")
		store.AddItemToBasket("basket-2", "Item B")

		// Attempting to complete item in basket-1 using basket-2 key should fail
		err := store.Store.SetItemCompletion("basket-2", id1, true)
		require.Error(t, err)
	})
}
