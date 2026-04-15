package resources

import (
	"context"
	"strings"
	"testing"
)

func TestConfigure_ApiKeyFromField(t *testing.T) {
	t.Parallel()
	cfg := &Config{ApiKey: "test-key"}
	if err := cfg.Configure(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.GetClient() == nil {
		t.Fatal("expected client to be created")
	}
	if cfg.GetClient().APIKey != "test-key" {
		t.Errorf("expected APIKey 'test-key', got %q", cfg.GetClient().APIKey)
	}
}

func TestConfigure_ApiKeyFromEnv(t *testing.T) {
	t.Setenv("ONEUPTIME_API_KEY", "env-key")
	cfg := &Config{}
	if err := cfg.Configure(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.GetClient().APIKey != "env-key" {
		t.Errorf("expected APIKey 'env-key', got %q", cfg.GetClient().APIKey)
	}
}

func TestConfigure_ApiKeyMissing(t *testing.T) {
	t.Setenv("ONEUPTIME_API_KEY", "")
	cfg := &Config{}
	err := cfg.Configure(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "apiKey is required") {
		t.Errorf("error %q does not contain 'apiKey is required'", err.Error())
	}
}

func TestConfigure_BaseURLDefault(t *testing.T) {
	t.Setenv("ONEUPTIME_BASE_URL", "")
	cfg := &Config{ApiKey: "test-key"}
	if err := cfg.Configure(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.GetClient().BaseURL != "https://oneuptime.com" {
		t.Errorf("expected default BaseURL, got %q", cfg.GetClient().BaseURL)
	}
}

func TestConfigure_BaseURLFromField(t *testing.T) {
	t.Parallel()
	url := "https://custom.example.com/"
	cfg := &Config{ApiKey: "test-key", BaseURL: &url}
	if err := cfg.Configure(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.GetClient().BaseURL != "https://custom.example.com" {
		t.Errorf("expected trailing slash stripped, got %q", cfg.GetClient().BaseURL)
	}
}

func TestConfigure_BaseURLFromEnv(t *testing.T) {
	t.Setenv("ONEUPTIME_BASE_URL", "https://env.example.com/")
	cfg := &Config{ApiKey: "test-key"}
	if err := cfg.Configure(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.GetClient().BaseURL != "https://env.example.com" {
		t.Errorf("expected env BaseURL with trailing slash stripped, got %q", cfg.GetClient().BaseURL)
	}
}

func TestConfigure_BaseURLFieldOverridesEnv(t *testing.T) {
	t.Setenv("ONEUPTIME_BASE_URL", "https://env.example.com")
	url := "https://field.example.com"
	cfg := &Config{ApiKey: "test-key", BaseURL: &url}
	if err := cfg.Configure(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.GetClient().BaseURL != "https://field.example.com" {
		t.Errorf("expected field to override env, got %q", cfg.GetClient().BaseURL)
	}
}

func TestConfigure_ProjectIDFromEnv(t *testing.T) {
	t.Setenv("ONEUPTIME_PROJECT_ID", "env-proj-id")
	cfg := &Config{ApiKey: "test-key"}
	if err := cfg.Configure(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ProjectID == nil {
		t.Fatal("expected ProjectID to be set from env")
	}
	if *cfg.ProjectID != "env-proj-id" {
		t.Errorf("expected ProjectID 'env-proj-id', got %q", *cfg.ProjectID)
	}
}

func TestConfigure_ProjectIDPreservesField(t *testing.T) {
	t.Setenv("ONEUPTIME_PROJECT_ID", "env-proj-id")
	projID := "field-proj-id"
	cfg := &Config{ApiKey: "test-key", ProjectID: &projID}
	if err := cfg.Configure(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *cfg.ProjectID != "field-proj-id" {
		t.Errorf("expected field ProjectID preserved, got %q", *cfg.ProjectID)
	}
}

func TestConfigure_ClientCreated(t *testing.T) {
	t.Parallel()
	cfg := &Config{ApiKey: "test-key"}
	if cfg.GetClient() != nil {
		t.Fatal("expected nil client before Configure")
	}
	if err := cfg.Configure(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c := cfg.GetClient()
	if c == nil {
		t.Fatal("expected non-nil client after Configure")
	}
	if c.HTTPClient == nil {
		t.Fatal("expected HTTPClient to be set")
	}
}
