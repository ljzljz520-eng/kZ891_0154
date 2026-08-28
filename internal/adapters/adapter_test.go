package adapters

import (
	"context"
	"testing"

	"warrantyservice/internal/persistence"
)

func TestWarrantyAdapterMissingReturnsNil(t *testing.T) {
	store, err := persistence.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	adapter := NewWarrantyAdapter(store)
	record, err := adapter.Lookup(context.Background(), "13600136000", "UNKNOWN-1")
	if err != nil {
		t.Fatal(err)
	}
	if record != nil {
		t.Fatalf("record=%+v", record)
	}
}

func TestCanonicalIdentifiers(t *testing.T) {
	if CanonicalPhone("+86 138-0013-8000") != "13800138000" {
		t.Fatal("phone normalization failed")
	}
	if CanonicalSerial(" tv 2025 0001 ") != "TV20250001" {
		t.Fatal("serial normalization failed")
	}
}
