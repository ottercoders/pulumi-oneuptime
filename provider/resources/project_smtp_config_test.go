package resources_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	p "github.com/pulumi/pulumi-go-provider"
	presource "github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/property"
)

func TestProjectSmtpConfig_LifeCycle(t *testing.T) {
	t.Parallel()

	var createBody map[string]interface{}
	var updateBody map[string]interface{}
	var lastUpdateID string

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		if r.Body != nil {
			json.NewDecoder(r.Body).Decode(&body)
		}
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/smtp-config":
			data, _ := body["data"].(map[string]interface{})
			createBody = data
			json.NewEncoder(w).Encode(map[string]interface{}{
				"_id":       "smtp-1",
				"name":      data["name"],
				"hostname":  data["hostname"],
				"port":      data["port"],
				"fromEmail": data["fromEmail"],
				"fromName":  data["fromName"],
				"projectId": data["projectId"],
				"slug":      "primary",
				"createdAt": "2024-01-01T00:00:00Z",
				"updatedAt": "2024-01-01T00:00:00Z",
			})

		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/get-item"):
			// Server's read ACL omits password and clientSecret.
			json.NewEncoder(w).Encode(map[string]interface{}{
				"_id":       "smtp-1",
				"name":      "primary",
				"hostname":  "smtp.example.com",
				"port":      587,
				"username":  "mailer",
				"fromEmail": "alerts@example.com",
				"fromName":  "Example Alerts",
				"secure":    true,
				"authType":  "Username and Password",
				"projectId": "test-project-id",
				"slug":      "primary",
				"createdAt": "2024-01-01T00:00:00Z",
				"updatedAt": "2024-01-01T00:00:00Z",
			})

		case r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/api/smtp-config/"):
			data, _ := body["data"].(map[string]interface{})
			updateBody = data
			lastUpdateID = strings.TrimPrefix(r.URL.Path, "/api/smtp-config/")
			w.WriteHeader(http.StatusOK)

		case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/api/smtp-config/"):
			w.WriteHeader(http.StatusOK)

		default:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "not found"})
		}
	})

	s := setupTestServer(t, handler)
	urn := presource.NewURN("test", "provider", "", "oneuptime:resources:ProjectSmtpConfig", "primary")

	// Create
	createResp, err := s.Create(p.CreateRequest{
		Urn: urn,
		Properties: property.NewMap(map[string]property.Value{
			"name":         property.New("primary"),
			"hostname":     property.New("smtp.example.com"),
			"port":         property.New(587.0),
			"username":     property.New("mailer"),
			"password":     property.New("hunter2").WithSecret(true),
			"fromEmail":    property.New("alerts@example.com"),
			"fromName":     property.New("Example Alerts"),
			"secure":       property.New(true),
			"authType":     property.New("Username and Password"),
			"clientSecret": property.New("oauth-secret").WithSecret(true),
		}),
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if createResp.ID != "smtp-1" {
		t.Errorf("expected ID smtp-1, got %q", createResp.ID)
	}

	// Wire-shape checks on the create body.
	if createBody == nil {
		t.Fatal("expected create body to be captured")
	}
	if createBody["projectId"] != "test-project-id" {
		t.Errorf("projectId not injected from config; got %v", createBody["projectId"])
	}
	if createBody["password"] != "hunter2" {
		t.Errorf("password missing on wire; got %v", createBody["password"])
	}
	if createBody["clientSecret"] != "oauth-secret" {
		t.Errorf("clientSecret missing on wire; got %v", createBody["clientSecret"])
	}
	if createBody["hostname"] != "smtp.example.com" {
		t.Errorf("hostname wrong on wire; got %v", createBody["hostname"])
	}

	// Read — server omits secrets, provider must preserve from prior state.
	readResp, err := s.Read(p.ReadRequest{
		ID:         createResp.ID,
		Urn:        urn,
		Properties: createResp.Properties,
		Inputs: property.NewMap(map[string]property.Value{
			"name":         property.New("primary"),
			"hostname":     property.New("smtp.example.com"),
			"port":         property.New(587.0),
			"fromEmail":    property.New("alerts@example.com"),
			"fromName":     property.New("Example Alerts"),
			"password":     property.New("hunter2").WithSecret(true),
			"clientSecret": property.New("oauth-secret").WithSecret(true),
		}),
	})
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if readResp.ID != "smtp-1" {
		t.Errorf("expected Read ID smtp-1, got %q", readResp.ID)
	}
	if got, ok := readResp.Properties.GetOk("password"); !ok || got.AsString() != "hunter2" {
		t.Errorf("Read should preserve password from prior state; got %v (ok=%v)", got, ok)
	}
	if got, ok := readResp.Properties.GetOk("clientSecret"); !ok || got.AsString() != "oauth-secret" {
		t.Errorf("Read should preserve clientSecret from prior state; got %v (ok=%v)", got, ok)
	}

	// Update — change fromName, rotate password.
	_, err = s.Update(p.UpdateRequest{
		ID:    createResp.ID,
		Urn:   urn,
		State: createResp.Properties,
		Inputs: property.NewMap(map[string]property.Value{
			"name":      property.New("primary"),
			"hostname":  property.New("smtp.example.com"),
			"port":      property.New(587.0),
			"username":  property.New("mailer"),
			"password":  property.New("hunter3").WithSecret(true),
			"fromEmail": property.New("alerts@example.com"),
			"fromName":  property.New("Example Ops"),
			"secure":    property.New(true),
			"authType":  property.New("Username and Password"),
		}),
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if lastUpdateID != "smtp-1" {
		t.Errorf("expected PUT to /api/smtp-config/smtp-1, got %q", lastUpdateID)
	}
	if updateBody == nil || updateBody["password"] != "hunter3" {
		t.Errorf("rotated password missing on update wire; got %v", updateBody["password"])
	}
	if updateBody["fromName"] != "Example Ops" {
		t.Errorf("fromName not updated on wire; got %v", updateBody["fromName"])
	}

	// Delete
	if err := s.Delete(p.DeleteRequest{
		ID:         createResp.ID,
		Urn:        urn,
		Properties: createResp.Properties,
	}); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestProjectSmtpConfig_ReadEmpty(t *testing.T) {
	t.Parallel()

	// Server returns 200 {} (deleted-resource convention).
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/get-item") {
			w.Write([]byte(`{}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})
	s := setupTestServer(t, handler)

	urn := presource.NewURN("test", "provider", "", "oneuptime:resources:ProjectSmtpConfig", "missing")
	resp, err := s.Read(p.ReadRequest{
		ID:         "deleted-id",
		Urn:        urn,
		Properties: property.Map{},
		Inputs:     property.Map{},
	})
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if resp.ID != "" {
		t.Errorf("expected empty ID for deleted resource, got %q", resp.ID)
	}
}
