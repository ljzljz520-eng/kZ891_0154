package main

import (
	"testing"
	"warrantyservice/internal/config"
)

func TestBuildHandler(t *testing.T) {
	cfg := config.Default()
	cfg.DatabasePath = ":memory:"
	handler, store, err := buildHandler(cfg)
	if err != nil || handler == nil || store == nil {
		t.Fatalf("handler=%v store=%v err=%v", handler, store, err)
	}
	store.Close()
}
