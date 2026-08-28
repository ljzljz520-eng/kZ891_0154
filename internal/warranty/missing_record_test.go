package warranty

import (
	"context"
	"testing"
)

func TestWarrantyMissingRecord(t *testing.T) {
	service, store := testService(t)
	defer store.Close()
	result, err := service.Query(context.Background(), "13600136000", "NO-SUCH-DEVICE")
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "missing" {
		t.Fatalf("expected missing status, got %+v", result)
	}
	if result.Advice == nil {
		t.Fatal("expected support advice")
	}
}

func TestMissingRecordAdviceWorkflow(t *testing.T) {
	service, store := testService(t)
	defer store.Close()
	advice, err := service.Advice(context.Background(), "13600136000", "NO-SUCH-DEVICE")
	if err != nil {
		t.Fatal(err)
	}
	if advice == nil || advice.Hotline == "" || advice.Message == "" {
		t.Fatalf("bad advice: %+v", advice)
	}
}
