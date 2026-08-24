package warranty

import (
	"context"
	"testing"

	"warrantyservice/internal/adapters"
	"warrantyservice/internal/config"
	"warrantyservice/internal/persistence"
	"warrantyservice/internal/support"
)

func testService(t *testing.T) (*Service, *persistence.Store) {
	t.Helper()
	store, err := persistence.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	for _, point := range config.FixturePoints() {
		if err := store.UpsertServicePoint(context.Background(), point); err != nil {
			t.Fatal(err)
		}
	}
	for _, record := range config.FixtureWarranties() {
		points, err := store.FindServicePoints(context.Background(), record.ServicePoints)
		if err != nil {
			t.Fatal(err)
		}
		if err := store.SaveWarranty(context.Background(), record, points); err != nil {
			t.Fatal(err)
		}
	}
	reader := adapters.NewWarrantyAdapter(store)
	service := NewService(reader, adapters.NewPointDirectory(store), adapters.NewAuditWriter(store), support.NewAdvisor(store, "400-800-1234", "电话、网点", "2025-06-01"), "2025-06-01")
	return service, store
}

func TestQueryWarrantyWorkflow(t *testing.T) {
	service, store := testService(t)
	defer store.Close()
	result, err := service.Query(context.Background(), "13800138000", "tv-2025-0001")
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "valid" || result.Record == nil {
		t.Fatalf("unexpected result: %+v", result)
	}
	if len(result.ServicePoints) != 2 {
		t.Fatalf("points=%d", len(result.ServicePoints))
	}
}

func TestAdminUpdateWorkflow(t *testing.T) {
	service, store := testService(t)
	defer store.Close()
	record, err := service.Upsert(context.Background(), UpsertInput{ID: 10, Phone: "13600136000", SerialNumber: "FR-2025-1010", ExpiryDate: "2027-12-31", PurchaseChannel: "官方商城", Category: "冰箱", UpdatedAt: "2025-06-01", Active: true, ServicePoints: config.FixturePoints()[:1]}, "operator")
	if err != nil {
		t.Fatal(err)
	}
	if record.ID != 10 {
		t.Fatalf("id=%d", record.ID)
	}
	result, err := service.Query(context.Background(), record.Phone, record.SerialNumber)
	if err != nil || result.Record == nil {
		t.Fatalf("query=%+v err=%v", result, err)
	}
}
