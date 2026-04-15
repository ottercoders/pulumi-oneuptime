package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateResource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/team" {
			t.Errorf("expected /api/team, got %s", r.URL.Path)
		}
		if r.Header.Get("APIKey") != "test-key" {
			t.Errorf("expected APIKey header 'test-key', got '%s'", r.Header.Get("APIKey"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got '%s'", r.Header.Get("Content-Type"))
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		data, ok := body["data"].(map[string]interface{})
		if !ok {
			t.Fatal("expected body wrapped in 'data' key")
		}
		if data["name"] != "Engineering" {
			t.Errorf("expected name 'Engineering', got '%v'", data["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"_id":       "abc123",
			"name":      "Engineering",
			"projectId": "proj-1",
			"createdAt": "2024-01-01T00:00:00Z",
		})
	}))
	defer server.Close()

	c := &Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: http.DefaultClient,
	}

	result, err := c.CreateResource(context.Background(), "team", map[string]interface{}{
		"name":      "Engineering",
		"projectId": "proj-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["_id"] != "abc123" {
		t.Errorf("expected _id 'abc123', got '%v'", result["_id"])
	}
}

func TestReadResource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/team/abc123/get-item" {
			t.Errorf("expected /api/team/abc123/get-item, got %s", r.URL.Path)
		}

		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		selectFields, ok := body["select"].(map[string]interface{})
		if !ok {
			t.Fatal("expected select in body")
		}
		if selectFields["name"] != true {
			t.Error("expected name in select fields")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"_id":  "abc123",
			"name": "Engineering",
		})
	}))
	defer server.Close()

	c := &Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: http.DefaultClient,
	}

	result, err := c.ReadResource(context.Background(), "team", "abc123", map[string]bool{"name": true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["name"] != "Engineering" {
		t.Errorf("expected name 'Engineering', got '%v'", result["name"])
	}
}

func TestUpdateResource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/team/abc123" {
			t.Errorf("expected /api/team/abc123, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := &Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: http.DefaultClient,
	}

	err := c.UpdateResource(context.Background(), "team", "abc123", map[string]interface{}{
		"name": "Platform",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteResource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/team/abc123" {
			t.Errorf("expected /api/team/abc123, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := &Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: http.DefaultClient,
	}

	err := c.DeleteResource(context.Background(), "team", "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteResource_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "not found"})
	}))
	defer server.Close()

	c := &Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: http.DefaultClient,
	}

	// Delete on a 404 should succeed (resource already gone)
	err := c.DeleteResource(context.Background(), "team", "abc123")
	if err != nil {
		t.Fatalf("expected no error for 404 delete, got: %v", err)
	}
}

func TestAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid api key"})
	}))
	defer server.Close()

	c := &Client{
		BaseURL:    server.URL,
		APIKey:     "bad-key",
		HTTPClient: http.DefaultClient,
	}

	_, err := c.CreateResource(context.Background(), "team", map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *APIError
	if !isAPIError(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != 401 {
		t.Errorf("expected status 401, got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "invalid api key" {
		t.Errorf("expected message 'invalid api key', got '%s'", apiErr.Message)
	}
}

func TestIsNotFound(t *testing.T) {
	err404 := &APIError{StatusCode: 404, Message: "not found"}
	err500 := &APIError{StatusCode: 500, Message: "internal error"}

	if !IsNotFound(err404) {
		t.Error("expected IsNotFound to return true for 404")
	}
	if IsNotFound(err500) {
		t.Error("expected IsNotFound to return false for 500")
	}
	if IsNotFound(nil) {
		t.Error("expected IsNotFound to return false for nil")
	}
}

func isAPIError(err error, target **APIError) bool {
	if err == nil {
		return false
	}
	apiErr, ok := err.(*APIError)
	if ok {
		*target = apiErr
	}
	return ok
}
