package db

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// region Setup
type TestStore struct {
	SqliteStore *SqliteStore
	t           *testing.T
}

func (ts *TestStore) AddItemToBasket(basketKey string, title string) (int64, error) {
	id, err := ts.SqliteStore.AddItemToBasket(basketKey, title)
	time.Sleep(1001 * time.Microsecond)
	return id, err
}

func (ts *TestStore) SetItemCompletion(basketKey string, id int64, completed bool) error {
	err := ts.SqliteStore.SetItemCompletion(basketKey, id, completed)
	time.Sleep(1001 * time.Microsecond)
	return err
}

func (ts *TestStore) GetItemsForBasket(basketKey string) ([]Item, error) {
	return ts.SqliteStore.GetItemsForBasket(basketKey)
}

func (ts *TestStore) Close() error {
	//TODO implement me
	panic("implement me")
}

func setup(t *testing.T) Store {
	s, err := NewSqliteStore(":memory:")
	require.NoError(t, err)
	return &TestStore{SqliteStore: s, t: t}
}

// endregion

func titles(items []Item) []string {
	res := make([]string, len(items))
	for i, itm := range items {
		res[i] = itm.Title
	}
	return res
}

const basket1Key = "7d74de4b-dd45-4aa8-b6ba-11484bc92a93"
const basket2Key = "afde3d43-b023-42dd-80bb-9bf36d09d83c"

func Test_Sorting(t *testing.T) {
	t.Run("New items appear at the top (touched_at DESC)", func(t *testing.T) {
		store := setup(t)

		store.AddItemToBasket(basket1Key, "Item 1")
		store.AddItemToBasket(basket1Key, "Item 2")

		items, _ := store.GetItemsForBasket(basket1Key)
		require.Equal(t, []string{"Item 2", "Item 1"}, titles(items))
	})

	t.Run("Re-adding an item moves it to the top (touching)", func(t *testing.T) {
		store := setup(t)

		store.AddItemToBasket(basket1Key, "Item 1")
		store.AddItemToBasket(basket1Key, "Item 2")
		store.AddItemToBasket(basket1Key, "Item 1") // Touched Item 1

		items, _ := store.GetItemsForBasket(basket1Key)
		require.Equal(t, []string{"Item 1", "Item 2"}, titles(items))
	})

	t.Run("Updating completion moves item to the top", func(t *testing.T) {
		store := setup(t)

		id1, _ := store.AddItemToBasket(basket1Key, "Item 1")
		store.AddItemToBasket(basket1Key, "Item 2")
		id3, _ := store.AddItemToBasket(basket1Key, "Item 3")
		id4, _ := store.AddItemToBasket(basket1Key, "Item 4")

		store.SetItemCompletion(basket1Key, id1, false)
		store.SetItemCompletion(basket1Key, id4, true)

		items, _ := store.GetItemsForBasket(basket1Key)
		require.Equal(t, []string{"Item 1", "Item 3", "Item 2", "Item 4"}, titles(items))

		store.SetItemCompletion(basket1Key, id3, true)
		items, _ = store.GetItemsForBasket(basket1Key)
		require.Equal(t, []string{"Item 1", "Item 2", "Item 3", "Item 4"}, titles(items))
	})
}

func Test_ValidateInput(t *testing.T) {
	t.Run("Adding item with empty title fails", func(t *testing.T) {
		store := setup(t)

		_, err := store.AddItemToBasket(basket1Key, "")
		require.Error(t, err)
	})

	t.Run("Too long item title fails", func(t *testing.T) {
		store := setup(t)

		longTitle := strings.Repeat("a", 256)
		_, err := store.AddItemToBasket(basket1Key, longTitle)
		require.Error(t, err)
	})

	t.Run("UUID v1 basket key fails", func(t *testing.T) {
		store := setup(t)

		uuidV1Key := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

		_, err := store.AddItemToBasket(uuidV1Key, "Valid Title")
		require.Error(t, err)

		err = store.SetItemCompletion(uuidV1Key, 1, false)
		require.Error(t, err)

		_, err = store.GetItemsForBasket(uuidV1Key)
		require.Error(t, err)
	})
}

func Test_IsolatedBaskets(t *testing.T) {
	t.Run("Items in different baskets do not interfere", func(t *testing.T) {
		store := setup(t)

		id1, _ := store.AddItemToBasket(basket1Key, "Item A")
		store.AddItemToBasket(basket2Key, "Item B")

		// Completing item in basket-1 should not affect basket-2
		err := store.SetItemCompletion(basket1Key, id1, true)
		require.NoError(t, err)

		items1, _ := store.GetItemsForBasket(basket1Key)
		require.Equal(t, []string{"Item A"}, titles(items1))
		require.True(t, items1[0].Completed)
	})

	t.Run("Completing item with wrong basket key fails", func(t *testing.T) {
		store := setup(t)

		id1, _ := store.AddItemToBasket(basket1Key, "Item A")
		store.AddItemToBasket(basket2Key, "Item B")

		// Attempting to complete item in basket-1 using basket-2 key should fail
		err := store.SetItemCompletion(basket2Key, id1, true)
		require.Error(t, err)
	})
}
