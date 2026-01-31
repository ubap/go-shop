package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Store(t *testing.T) {
	store, err := NewStore(":memory:")
	require.NoError(t, err)

	itemId, err := store.AddItem("Test Item")
	require.NoError(t, err)

	itemId2, err := store.AddItem("Test Item 2")
	require.NoError(t, err)

	err = store.ToggleItem(itemId2, true)
	require.NoError(t, err)

	items, err := store.GetAllItems()
	require.NoError(t, err)
	require.ElementsMatch(t, []Item{
		{
			ID:        itemId,
			Title:     "Test Item",
			Completed: false,
		},
		{
			ID:        itemId2,
			Title:     "Test Item 2",
			Completed: true,
		},
	}, items)
}
