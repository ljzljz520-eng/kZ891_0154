package support

import (
	"context"
	"testing"

	"warrantyservice/internal/persistence"
)

func TestAdvicePersistence(t *testing.T) {
	store, err := persistence.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	advisor := NewAdvisor(store, "400-800-1234", "官方小程序", "2025-06-01")
	advice, err := advisor.PersistMissing(context.Background(), "13600136000", "UNKNOWN-1")
	if err != nil || advice == nil {
		t.Fatalf("advice=%+v err=%v", advice, err)
	}
	latest, err := advisor.GetLatest(context.Background(), advice.Phone, advice.SerialNumber)
	if err != nil || latest == nil {
		t.Fatalf("latest=%+v err=%v", latest, err)
	}
}

func TestEscalationClassification(t *testing.T) {
	event := ClassifyMissing("13800138000", "TV-UNKNOWN")
	if event.Level != "L2" {
		t.Fatalf("event=%+v", event)
	}
}
