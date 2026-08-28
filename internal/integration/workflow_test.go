package integration

import (
	"context"
	"testing"

	"warrantyservice/internal/adapters"
	"warrantyservice/internal/config"
	"warrantyservice/internal/persistence"
	"warrantyservice/internal/support"
	"warrantyservice/internal/warranty"
)

func integrationService(t *testing.T) (*warranty.Service, *persistence.Store) {
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
		points, _ := store.FindServicePoints(context.Background(), record.ServicePoints)
		if err := store.SaveWarranty(context.Background(), record, points); err != nil {
			t.Fatal(err)
		}
	}
	reader := adapters.NewWarrantyAdapter(store)
	return warranty.NewService(reader, adapters.NewPointDirectory(store), adapters.NewAuditWriter(store), support.NewAdvisor(store, "400-800-1234", "电话、网点", "2025-06-01"), "2025-06-01"), store
}

func TestPrimaryWorkflow(t *testing.T) {
	service, store := integrationService(t)
	defer store.Close()
	result, err := service.Query(context.Background(), "13800138000", "TV-2025-0001")
	if err != nil || result.Record == nil {
		t.Fatalf("result=%+v err=%v", result, err)
	}
}

func TestSecondaryWorkflow(t *testing.T) {
	service, store := integrationService(t)
	defer store.Close()
	_, err := service.Upsert(context.Background(), warranty.UpsertInput{ID: 44, Phone: "13600136000", SerialNumber: "OV-2025-0044", ExpiryDate: "2028-01-01", PurchaseChannel: "授权门店", Category: "厨房电器", UpdatedAt: "2025-06-01", Active: true}, "operator")
	if err != nil {
		t.Fatal(err)
	}
	result, err := service.Query(context.Background(), "13600136000", "OV-2025-0044")
	if err != nil || result.Record == nil {
		t.Fatalf("result=%+v err=%v", result, err)
	}
}

func TestTertiaryWorkflow(t *testing.T) {
	service, store := integrationService(t)
	defer store.Close()
	advice, err := service.Advice(context.Background(), "13600136000", "UNKNOWN-TERTIARY")
	if err != nil || advice == nil {
		t.Fatalf("advice=%+v err=%v", advice, err)
	}
}
