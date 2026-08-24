package config

import "testing"

func TestConfigDefaults(t *testing.T) {
	cfg := Default()
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
	if cfg.IsMemoryDatabase() {
		t.Fatal("default should be file backed")
	}
	if cfg.HotlineMessage() == "" {
		t.Fatal("missing hotline")
	}
}

func TestBusinessPolicy(t *testing.T) {
	policy := DefaultPolicy()
	if !policy.SupportsCategory("电视") || policy.SupportsCategory("未知") {
		t.Fatal("category policy failed")
	}
	if err := policy.ValidatePointCount(99); err == nil {
		t.Fatal("expected count error")
	}
}
