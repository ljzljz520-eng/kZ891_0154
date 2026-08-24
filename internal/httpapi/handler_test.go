package httpapi

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"warrantyservice/internal/adapters"
	"warrantyservice/internal/config"
	"warrantyservice/internal/persistence"
	"warrantyservice/internal/support"
	"warrantyservice/internal/warranty"
)

func testHandler(t *testing.T) (*Handler, *persistence.Store) {
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
	service := warranty.NewService(adapters.NewWarrantyAdapter(store), adapters.NewPointDirectory(store), adapters.NewAuditWriter(store), support.NewAdvisor(store, "400-800-1234", "电话、网点", "2025-06-01"), "2025-06-01")
	return NewHandler(service, config.Default()), store
}

func TestQueryHandlerJSON(t *testing.T) {
	handler, store := testHandler(t)
	defer store.Close()
	req := httptest.NewRequest(http.MethodGet, "/warranty/query?phone=13800138000&serial=TV-2025-0001&format=json", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "延保有效") {
		t.Fatalf("body=%s", rec.Body.String())
	}
}

func TestAdminUpsertHandler(t *testing.T) {
	handler, store := testHandler(t)
	defer store.Close()
	body := `{"id":22,"phone":"13600136000","serial_number":"FR-2025-2222","expiry_date":"2027-12-31","purchase_channel":"官方商城","category":"冰箱","updated_at":"2025-06-01","active":true,"service_points":[]}`
	req := httptest.NewRequest(http.MethodPost, "/admin/warranties", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+config.Default().AdminToken)
	req.Header.Set("X-Operator", "test-admin")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	data, _ := io.ReadAll(rec.Result().Body)
	if !strings.Contains(string(data), "saved") {
		t.Fatalf("body=%s", data)
	}
}
