package persistence

import (
	"context"
	"path/filepath"
	"testing"

	"warrantyservice/internal/config"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "warranty.db")
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	record := config.FixtureWarranties()[0]
	if err := store.SaveWarranty(context.Background(), record, config.FixturePoints()[:2]); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	store, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	found, err := store.FindWarranty(context.Background(), record.Phone, record.SerialNumber)
	if err != nil || found == nil {
		t.Fatalf("found=%+v err=%v", found, err)
	}
	if len(found.ServicePoints) != 2 {
		t.Fatalf("points=%v", found.ServicePoints)
	}
}

func TestStoreIntegrity(t *testing.T) {
	store, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	health := store.Health(context.Background())
	if !health.Ready {
		t.Fatalf("health=%+v", health)
	}
	if err := store.EnsureReady(context.Background()); err != nil {
		t.Fatal(err)
	}
}
