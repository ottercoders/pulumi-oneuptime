package resources

import (
	"strings"
	"testing"
)

// --- ToMap ---

func TestToMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    interface{}
		wantKeys []string
		wantVals map[string]interface{}
	}{
		{
			name: "all fields populated",
			input: struct {
				Name  string `json:"name"`
				Color string `json:"color"`
			}{Name: "prod", Color: "#ff0000"},
			wantVals: map[string]interface{}{
				"name":  "prod",
				"color": "#ff0000",
			},
		},
		{
			name: "nil pointer omitted with omitempty",
			input: struct {
				Name string  `json:"name"`
				Desc *string `json:"description,omitempty"`
			}{Name: "prod"},
			wantKeys: []string{"name"},
		},
		{
			name: "zero string with omitempty omitted",
			input: struct {
				Name string `json:"name,omitempty"`
			}{Name: ""},
			wantKeys: []string{},
		},
		{
			name:     "empty struct",
			input:    struct{}{},
			wantKeys: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, err := ToMap(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.wantVals != nil {
				for k, want := range tc.wantVals {
					got, ok := result[k]
					if !ok {
						t.Errorf("expected key %q in result", k)
						continue
					}
					if got != want {
						t.Errorf("key %q: got %v, want %v", k, got, want)
					}
				}
			}

			if tc.wantKeys != nil {
				if len(result) != len(tc.wantKeys) {
					t.Errorf("expected %d keys, got %d: %v", len(tc.wantKeys), len(result), result)
				}
				for _, k := range tc.wantKeys {
					if _, ok := result[k]; !ok {
						t.Errorf("expected key %q in result", k)
					}
				}
			}
		})
	}
}

func TestToMap_WithRealArgs(t *testing.T) {
	t.Parallel()

	desc := "Engineering team"
	args := TeamArgs{
		Name:        "Engineering",
		Description: &desc,
	}

	result, err := ToMap(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["name"] != "Engineering" {
		t.Errorf("expected name 'Engineering', got %v", result["name"])
	}
	if result["description"] != "Engineering team" {
		t.Errorf("expected description 'Engineering team', got %v", result["description"])
	}
	// projectId should be omitted (nil pointer with omitempty)
	if _, ok := result["projectId"]; ok {
		t.Error("expected projectId to be omitted when nil")
	}
}

// --- FromMap ---

func TestFromMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   map[string]interface{}
		wantErr bool
	}{
		{
			name: "full map",
			input: map[string]interface{}{
				"_id":       "abc123",
				"name":      "Engineering",
				"slug":      "engineering",
				"createdAt": "2024-01-01T00:00:00Z",
			},
		},
		{
			name: "partial map",
			input: map[string]interface{}{
				"_id":  "abc123",
				"name": "Engineering",
			},
		},
		{
			name: "extra keys ignored",
			input: map[string]interface{}{
				"_id":     "abc123",
				"name":    "Engineering",
				"unknown": "should be ignored",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var state TeamState
			err := FromMap(tc.input, &state)
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tc.wantErr {
				if state.ResourceID != tc.input["_id"] {
					t.Errorf("ResourceID: got %q, want %q", state.ResourceID, tc.input["_id"])
				}
				if state.Name != tc.input["name"] {
					t.Errorf("Name: got %q, want %q", state.Name, tc.input["name"])
				}
			}
		})
	}
}

// --- SelectFields ---

func TestSelectFields(t *testing.T) {
	t.Parallel()

	t.Run("simple struct", func(t *testing.T) {
		t.Parallel()
		type simple struct {
			Name  string `json:"name"`
			Color string `json:"color"`
		}
		fields := SelectFields(simple{})
		expectFields(t, fields, []string{"name", "color"})
	})

	t.Run("TeamState with embedded TeamArgs", func(t *testing.T) {
		t.Parallel()
		fields := SelectFields(TeamState{})
		// Should contain fields from both TeamArgs and TeamState
		expectFields(t, fields, []string{
			"projectId", "name", "description", // from TeamArgs
			"_id", "slug", "createdAt", "updatedAt", // from TeamState
		})
	})

	t.Run("omitempty tag still extracted", func(t *testing.T) {
		t.Parallel()
		type tagged struct {
			Name string `json:"name,omitempty"`
		}
		fields := SelectFields(tagged{})
		if !fields["name"] {
			t.Error("expected 'name' in fields despite omitempty")
		}
	})

	t.Run("dash tag excluded", func(t *testing.T) {
		t.Parallel()
		type tagged struct {
			Internal string `json:"-"`
			Name     string `json:"name"`
		}
		fields := SelectFields(tagged{})
		if fields["-"] {
			t.Error("expected '-' tag to be excluded")
		}
		if !fields["name"] {
			t.Error("expected 'name' in fields")
		}
		if len(fields) != 1 {
			t.Errorf("expected 1 field, got %d: %v", len(fields), fields)
		}
	})

	t.Run("no json tag excluded", func(t *testing.T) {
		t.Parallel()
		type tagged struct {
			Internal string
			Name     string `json:"name"`
		}
		fields := SelectFields(tagged{})
		if len(fields) != 1 {
			t.Errorf("expected 1 field, got %d: %v", len(fields), fields)
		}
	})

	t.Run("pointer to struct", func(t *testing.T) {
		t.Parallel()
		fields := SelectFields(&TeamState{})
		if !fields["name"] {
			t.Error("expected 'name' in fields when passing pointer")
		}
	})

	t.Run("IncidentStateResourceState with embedded args", func(t *testing.T) {
		t.Parallel()
		fields := SelectFields(IncidentStateResourceState{})
		expectFields(t, fields, []string{
			"projectId", "name", "color", "description",
			"isCreatedState", "isAcknowledgedState", "isResolvedState", "order",
			"_id", "slug", "createdAt", "updatedAt",
		})
	})
}

func expectFields(t *testing.T, fields map[string]bool, expected []string) {
	t.Helper()
	for _, name := range expected {
		if !fields[name] {
			t.Errorf("expected field %q in result, got: %v", name, fields)
		}
	}
}

// --- ResolveProjectID ---

func TestResolveProjectID(t *testing.T) {
	t.Parallel()

	strPtr := func(s string) *string { return &s }

	tests := []struct {
		name      string
		resource  *string
		config    *string
		want      string
		wantErr   bool
		errSubstr string
	}{
		{
			name:     "resource set, config set -> resource wins",
			resource: strPtr("res-proj"),
			config:   strPtr("cfg-proj"),
			want:     "res-proj",
		},
		{
			name:     "resource nil, config set -> config wins",
			resource: nil,
			config:   strPtr("cfg-proj"),
			want:     "cfg-proj",
		},
		{
			name:     "resource set, config nil -> resource wins",
			resource: strPtr("res-proj"),
			config:   nil,
			want:     "res-proj",
		},
		{
			name:      "both nil -> error",
			resource:  nil,
			config:    nil,
			wantErr:   true,
			errSubstr: "projectId is required",
		},
		{
			name:     "resource empty string, config set -> config wins",
			resource: strPtr(""),
			config:   strPtr("cfg-proj"),
			want:     "cfg-proj",
		},
		{
			name:      "both empty string -> error",
			resource:  strPtr(""),
			config:    strPtr(""),
			wantErr:   true,
			errSubstr: "projectId is required",
		},
		{
			name:      "resource nil, config empty string -> error",
			resource:  nil,
			config:    strPtr(""),
			wantErr:   true,
			errSubstr: "projectId is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := ResolveProjectID(tc.resource, tc.config)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tc.errSubstr != "" && !strings.Contains(err.Error(), tc.errSubstr) {
					t.Errorf("error %q does not contain %q", err.Error(), tc.errSubstr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
