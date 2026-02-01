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
