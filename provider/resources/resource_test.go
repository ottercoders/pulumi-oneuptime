package resources_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/blang/semver"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/integration"
	presource "github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/property"

	provider "github.com/ottercoders/pulumi-oneuptime/provider"
)

// setupTestServer creates an integration.Server backed by the real provider,
// configured to talk to the given httptest handler.
func setupTestServer(t *testing.T, handler http.Handler) integration.Server {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	prov := provider.Provider()
	s, err := integration.NewServer(context.Background(), "oneuptime",
		semver.MustParse("0.0.1"),
		integration.WithProvider(prov),
	)
	if err != nil {
		t.Fatalf("failed to create integration server: %v", err)
	}

	// Configure the provider with a test API key and the httptest URL
	err = s.Configure(p.ConfigureRequest{
		Args: property.NewMap(map[string]property.Value{
			"apiKey":    property.New("test-api-key"),
			"baseUrl":   property.New(ts.URL),
			"projectId": property.New("test-project-id"),
		}),
	})
	if err != nil {
		t.Fatalf("failed to configure provider: %v", err)
	}

	return s
}

// mockAPI is a simple HTTP handler that records requests and returns canned responses.
type mockAPI struct {
	mu       sync.Mutex
	requests []recordedRequest
}

type recordedRequest struct {
	Method string
	Path   string
	Body   map[string]interface{}
}

func (m *mockAPI) getRequests() []recordedRequest {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make([]recordedRequest, len(m.requests))
	copy(cp, m.requests)
	return cp
}

func newTeamMockAPI(t *testing.T) *mockAPI {
	t.Helper()
	return &mockAPI{}
}

func (m *mockAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if r.Body != nil {
		json.NewDecoder(r.Body).Decode(&body)
	}

	m.mu.Lock()
	m.requests = append(m.requests, recordedRequest{
		Method: r.Method,
		Path:   r.URL.Path,
		Body:   body,
	})
	m.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	switch {
	// Create team
	case r.Method == http.MethodPost && r.URL.Path == "/api/team":
		name, _ := body["name"].(string)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"_id":       "team-abc123",
			"name":      name,
			"projectId": body["projectId"],
			"slug":      strings.ToLower(name),
			"createdAt": "2024-01-01T00:00:00Z",
			"updatedAt": "2024-01-01T00:00:00Z",
		})

	// Read team
	case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/get-item"):
		json.NewEncoder(w).Encode(map[string]interface{}{
			"_id":       "team-abc123",
			"name":      "Engineering",
			"projectId": "test-project-id",
			"slug":      "engineering",
			"createdAt": "2024-01-01T00:00:00Z",
			"updatedAt": "2024-01-01T00:00:00Z",
		})

	// Update team
	case r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/api/team/"):
		w.WriteHeader(http.StatusOK)

	// Delete team
	case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/api/team/"):
		w.WriteHeader(http.StatusOK)

	// Create project
	case r.Method == http.MethodPost && r.URL.Path == "/api/project":
		name, _ := body["name"].(string)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"_id":       "proj-abc123",
			"name":      name,
			"slug":      strings.ToLower(name),
			"createdAt": "2024-01-01T00:00:00Z",
			"updatedAt": "2024-01-01T00:00:00Z",
		})

	// Read project
	case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/api/project/") && strings.HasSuffix(r.URL.Path, "/get-item"):
		json.NewEncoder(w).Encode(map[string]interface{}{
			"_id":       "proj-abc123",
			"name":      "My Project",
			"slug":      "my-project",
			"createdAt": "2024-01-01T00:00:00Z",
			"updatedAt": "2024-01-01T00:00:00Z",
		})

	// Update project
	case r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/api/project/"):
		w.WriteHeader(http.StatusOK)

	// Delete project
	case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/api/project/"):
		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "not found"})
	}
}

func TestTeamResource_LifeCycle(t *testing.T) {
	t.Parallel()

	mock := newTeamMockAPI(t)
	s := setupTestServer(t, mock)

	urn := presource.NewURN("test", "provider", "", "oneuptime:resources:Team", "test-team")

	// Create
	createResp, err := s.Create(p.CreateRequest{
		Urn: urn,
		Properties: property.NewMap(map[string]property.Value{
			"name": property.New("Engineering"),
		}),
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if createResp.ID == "" {
		t.Error("expected non-empty ID from Create")
	}

	// Verify the API was called with projectId from config
	reqs := mock.getRequests()
	var createReq *recordedRequest
	for i := range reqs {
		if reqs[i].Method == http.MethodPost && reqs[i].Path == "/api/team" {
			createReq = &reqs[i]
			break
		}
	}
	if createReq == nil {
		t.Fatal("expected POST /api/team request")
	}
	// Body is wrapped in {"data": {...}}
	dataMap, ok := createReq.Body["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected create body to be wrapped in 'data' key")
	}
	if dataMap["projectId"] != "test-project-id" {
		t.Errorf("expected projectId 'test-project-id' in create body data, got %v", dataMap["projectId"])
	}

	// Read
	readResp, err := s.Read(p.ReadRequest{
		ID:  createResp.ID,
		Urn: urn,
		Properties: createResp.Properties,
		Inputs: property.NewMap(map[string]property.Value{
			"name": property.New("Engineering"),
		}),
	})
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if readResp.ID == "" {
		t.Error("expected non-empty ID from Read")
	}

	// Update
	_, err = s.Update(p.UpdateRequest{
		ID:    createResp.ID,
		Urn:   urn,
		State: createResp.Properties,
		Inputs: property.NewMap(map[string]property.Value{
			"name": property.New("Platform"),
		}),
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Delete
	err = s.Delete(p.DeleteRequest{
		ID:         createResp.ID,
		Urn:        urn,
		Properties: createResp.Properties,
	})
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestTeamResource_CreateDryRun(t *testing.T) {
	t.Parallel()

	mock := newTeamMockAPI(t)
	s := setupTestServer(t, mock)

	urn := presource.NewURN("test", "provider", "", "oneuptime:resources:Team", "test-team")

	resp, err := s.Create(p.CreateRequest{
		Urn: urn,
		Properties: property.NewMap(map[string]property.Value{
			"name": property.New("Engineering"),
		}),
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("DryRun Create failed: %v", err)
	}
	if resp.ID != "preview-id" {
		t.Errorf("expected preview-id, got %q", resp.ID)
	}

	// Verify no actual API calls were made (only the configure call may have been recorded)
	reqs := mock.getRequests()
	for _, req := range reqs {
		if req.Path == "/api/team" {
			t.Error("expected no API call during DryRun, but POST /api/team was called")
		}
	}
}

func TestTeamResource_ReadNotFound(t *testing.T) {
	t.Parallel()

	// Handler that always returns 404
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "not found"})
	})
	s := setupTestServer(t, handler)

	urn := presource.NewURN("test", "provider", "", "oneuptime:resources:Team", "test-team")

	resp, err := s.Read(p.ReadRequest{
		ID:         "nonexistent-id",
		Urn:        urn,
		Properties: property.Map{},
		Inputs:     property.Map{},
	})
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	// A 404 should return an empty ID (resource was deleted externally)
	if resp.ID != "" {
		t.Errorf("expected empty ID for not-found resource, got %q", resp.ID)
	}
}

func TestProjectResource_LifeCycle(t *testing.T) {
	t.Parallel()

	mock := newTeamMockAPI(t)
	s := setupTestServer(t, mock)

	urn := presource.NewURN("test", "provider", "", "oneuptime:resources:Project", "test-project")

	// Create
	createResp, err := s.Create(p.CreateRequest{
		Urn: urn,
		Properties: property.NewMap(map[string]property.Value{
			"name": property.New("My Project"),
		}),
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if createResp.ID == "" {
		t.Error("expected non-empty ID from Create")
	}

	// Verify no projectId was sent in the create body
	reqs := mock.getRequests()
	var createReq *recordedRequest
	for i := range reqs {
		if reqs[i].Method == http.MethodPost && reqs[i].Path == "/api/project" {
			createReq = &reqs[i]
			break
		}
	}
	if createReq == nil {
		t.Fatal("expected POST /api/project request")
	}
	if _, hasProjectID := createReq.Body["projectId"]; hasProjectID {
		t.Error("Project resource should not send projectId in request body")
	}

	// Delete
	err = s.Delete(p.DeleteRequest{
		ID:         createResp.ID,
		Urn:        urn,
		Properties: createResp.Properties,
	})
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestTeamResource_ApiKeyHeader(t *testing.T) {
	t.Parallel()

	var capturedAPIKey string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAPIKey = r.Header.Get("APIKey")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"_id":       "team-abc123",
			"name":      "Test",
			"slug":      "test",
			"createdAt": "2024-01-01T00:00:00Z",
			"updatedAt": "2024-01-01T00:00:00Z",
		})
	})
	s := setupTestServer(t, handler)

	urn := presource.NewURN("test", "provider", "", "oneuptime:resources:Team", "test-team")
	_, err := s.Create(p.CreateRequest{
		Urn: urn,
		Properties: property.NewMap(map[string]property.Value{
			"name": property.New("Test"),
		}),
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if capturedAPIKey != "test-api-key" {
		t.Errorf("expected ApiKey header 'test-api-key', got %q", capturedAPIKey)
	}
}
