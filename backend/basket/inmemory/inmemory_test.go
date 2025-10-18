package inmemory

import (
	"go-shop/backend/basket"
	"testing"
)

func Test_AddNewItem(t *testing.T) {
	sut := basket.NewInMemoryBasket()

	want := basket.Item{ItemID: "itemId", Name: "Water", LastModified: "2023-10-27T10:00:01Z"}
	sut.UpsertItem(want)

	items := sut.GetAllItems()

	if len(items) != 1 {
		t.Fatalf("GetAllItems() returned wrong number of items; want 1, got %d", len(items))
	}
	got := items[0]
	if got != want {
		t.Errorf("GetAllItems() returned different want;\nwant %+v\ngot  %+v", want, got)
	}
}
