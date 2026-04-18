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

func TestCreateResource_TenantFromBody(t *testing.T) {
	var gotTenant string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotTenant = r.Header.Get("tenantid")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"_id": "abc"})
	}))
	defer server.Close()

	// Client has no TenantID configured; the body's projectId must supply it.
	c := &Client{BaseURL: server.URL, APIKey: "k", HTTPClient: http.DefaultClient}
	_, err := c.CreateResource(context.Background(), "monitor", map[string]interface{}{
		"name":      "m",
		"projectId": "proj-from-arg",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotTenant != "proj-from-arg" {
		t.Errorf("expected tenantid header 'proj-from-arg', got %q", gotTenant)
	}
}

func TestCreateResource_BodyTenantOverridesConfig(t *testing.T) {
	var gotTenant string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotTenant = r.Header.Get("tenantid")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"_id": "abc"})
	}))
	defer server.Close()

	c := &Client{BaseURL: server.URL, APIKey: "k", TenantID: "config-proj", HTTPClient: http.DefaultClient}
	_, err := c.CreateResource(context.Background(), "monitor", map[string]interface{}{
		"projectId": "arg-proj",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotTenant != "arg-proj" {
		t.Errorf("expected body projectId to win: got tenantid %q", gotTenant)
	}
}

func TestCreateResource_FallsBackToConfigTenant(t *testing.T) {
	var gotTenant string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotTenant = r.Header.Get("tenantid")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"_id": "abc"})
	}))
	defer server.Close()

	c := &Client{BaseURL: server.URL, APIKey: "k", TenantID: "config-proj", HTTPClient: http.DefaultClient}
	_, err := c.CreateResource(context.Background(), "project", map[string]interface{}{
		"name": "p",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotTenant != "config-proj" {
		t.Errorf("expected tenantid header 'config-proj', got %q", gotTenant)
	}
}

func TestUpdateResource_TenantFromBody(t *testing.T) {
	var gotTenant string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotTenant = r.Header.Get("tenantid")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := &Client{BaseURL: server.URL, APIKey: "k", HTTPClient: http.DefaultClient}
	err := c.UpdateResource(context.Background(), "monitor", "abc", map[string]interface{}{
		"projectId": "proj-from-arg",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotTenant != "proj-from-arg" {
		t.Errorf("expected tenantid header 'proj-from-arg', got %q", gotTenant)
	}
}

func TestReadResource_RetriesOnUnknownColumn(t *testing.T) {
	var attempts int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		sel, _ := body["select"].(map[string]interface{})

		// First request asks for the missing "description" column — reject it.
		if _, hasDesc := sel["description"]; hasDesc {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "TableColumnMetadata not found for description column",
			})
			return
		}

		// Second request omits description — succeed.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"_id":  "abc",
			"name": "Engineering",
		})
	}))
	defer server.Close()

	c := &Client{BaseURL: server.URL, APIKey: "k", HTTPClient: http.DefaultClient}
	result, err := c.ReadResource(context.Background(), "team", "abc", map[string]bool{
		"_id":         true,
		"name":        true,
		"description": true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["name"] != "Engineering" {
		t.Errorf("expected name Engineering, got %v", result["name"])
	}
	if attempts != 2 {
		t.Errorf("expected 2 attempts (initial + retry), got %d", attempts)
	}
}

func TestReadResource_CachesUnknownColumn(t *testing.T) {
	var saw []bool // whether each call included "description"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		sel, _ := body["select"].(map[string]interface{})
		_, hasDesc := sel["description"]
		saw = append(saw, hasDesc)

		if hasDesc {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "TableColumnMetadata not found for description column",
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"_id": "abc"})
	}))
	defer server.Close()

	c := &Client{BaseURL: server.URL, APIKey: "k", HTTPClient: http.DefaultClient}
	sel := map[string]bool{"_id": true, "description": true}

	// First read: server rejects, client retries without description.
	if _, err := c.ReadResource(context.Background(), "team", "abc", sel); err != nil {
		t.Fatalf("first read: %v", err)
	}
	// Second read: client should remember description is bad and skip it upfront.
	if _, err := c.ReadResource(context.Background(), "team", "abc", sel); err != nil {
		t.Fatalf("second read: %v", err)
	}

	// Expect: first call had description (rejected), second call didn't (retry succeeded),
	// third call (fresh ReadResource) skipped it upfront.
	if len(saw) != 3 {
		t.Fatalf("expected 3 total server calls, got %d (saw=%v)", len(saw), saw)
	}
	if !saw[0] || saw[1] || saw[2] {
		t.Errorf("expected description in call 1 only; saw=%v", saw)
	}
}

func TestListResources_RetriesOnUnknownColumn(t *testing.T) {
	var attempts int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		sel, _ := body["select"].(map[string]interface{})
		if _, bad := sel["isDegradedState"]; bad {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "TableColumnMetadata not found for isDegradedState column",
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []interface{}{map[string]interface{}{"_id": "abc", "name": "Operational"}},
		})
	}))
	defer server.Close()

	c := &Client{BaseURL: server.URL, APIKey: "k", HTTPClient: http.DefaultClient}
	results, err := c.ListResources(context.Background(), "monitor-status", map[string]interface{}{"name": "Operational"}, map[string]bool{
		"_id":             true,
		"name":            true,
		"isDegradedState": true,
	}, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0]["name"] != "Operational" {
		t.Errorf("unexpected results: %v", results)
	}
	if attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}

func TestReadResource_OtherBadRequestNotRetried(t *testing.T) {
	var attempts int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "malformed id"})
	}))
	defer server.Close()

	c := &Client{BaseURL: server.URL, APIKey: "k", HTTPClient: http.DefaultClient}
	_, err := c.ReadResource(context.Background(), "team", "abc", map[string]bool{"name": true})
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts != 1 {
		t.Errorf("expected 1 attempt (no retry on unrelated 400), got %d", attempts)
	}
}
