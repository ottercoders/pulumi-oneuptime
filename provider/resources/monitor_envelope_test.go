package resources

import (
	"testing"
)

func strPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }
func boolPtr(b bool) *bool    { return &b }

func TestAttachMonitorSteps_Wrap(t *testing.T) {
	t.Parallel()
	data := map[string]interface{}{}

	steps := []MonitorStep{{
		ID:                     "step-1",
		MonitorDestination:     strPtr("https://example.com"),
		MonitorDestinationPort: intPtr(443),
		RequestType:            strPtr("GET"),
		MonitorCriteria: MonitorCriteria{
			CriteriaInstances: []MonitorCriteriaInstance{{
				ID:              "ci-1",
				Name:            "Down if 5xx",
				Description:     "Mark offline on server error",
				FilterCondition: "Any",
				Filters: []CriteriaFilter{{
					CheckOn:    "Response Status Code",
					FilterType: strPtr("Greater Than"),
					Value:      strPtr("499"),
				}},
				MonitorStatusID:     strPtr("status-offline"),
				ChangeMonitorStatus: boolPtr(true),
			}},
		},
	}}

	if err := attachMonitorSteps(data, steps, strPtr("status-operational")); err != nil {
		t.Fatalf("attachMonitorSteps: %v", err)
	}

	outer, ok := data["monitorSteps"].(map[string]interface{})
	if !ok {
		t.Fatalf("monitorSteps not enveloped: %v", data["monitorSteps"])
	}
	if outer["_type"] != "MonitorSteps" {
		t.Errorf("outer _type = %v, want MonitorSteps", outer["_type"])
	}
	outerValue := outer["value"].(map[string]interface{})
	if outerValue["defaultMonitorStatusId"] != "status-operational" {
		t.Errorf("defaultMonitorStatusId = %v", outerValue["defaultMonitorStatusId"])
	}

	arr := outerValue["monitorStepsInstanceArray"].([]map[string]interface{})
	if len(arr) != 1 {
		t.Fatalf("expected 1 step, got %d", len(arr))
	}

	stepEnv := arr[0]
	if stepEnv["_type"] != "MonitorStep" {
		t.Errorf("step _type = %v, want MonitorStep", stepEnv["_type"])
	}
	stepVal := stepEnv["value"].(map[string]interface{})
	if stepVal["id"] != "step-1" {
		t.Errorf("step id = %v", stepVal["id"])
	}

	dest := stepVal["monitorDestination"].(map[string]interface{})
	if dest["_type"] != "URL" || dest["value"] != "https://example.com" {
		t.Errorf("monitorDestination envelope wrong: %v", dest)
	}

	port := stepVal["monitorDestinationPort"].(map[string]interface{})
	if port["_type"] != "Port" || port["value"] != "443" {
		t.Errorf("monitorDestinationPort envelope wrong: %v", port)
	}

	crit := stepVal["monitorCriteria"].(map[string]interface{})
	if crit["_type"] != "MonitorCriteria" {
		t.Errorf("criteria _type = %v", crit["_type"])
	}
	critValue := crit["value"].(map[string]interface{})
	instArr := critValue["monitorCriteriaInstanceArray"].([]map[string]interface{})
	if len(instArr) != 1 {
		t.Fatalf("expected 1 criteria instance, got %d", len(instArr))
	}
	inst := instArr[0]
	if inst["_type"] != "MonitorCriteriaInstance" {
		t.Errorf("instance _type = %v", inst["_type"])
	}
	instVal := inst["value"].(map[string]interface{})
	if instVal["name"] != "Down if 5xx" || instVal["filterCondition"] != "Any" {
		t.Errorf("criteria instance value wrong: %v", instVal)
	}
	if instVal["id"] != "ci-1" {
		t.Errorf("criteria instance id = %v", instVal["id"])
	}

	filters := instVal["filters"].([]interface{})
	if len(filters) != 1 {
		t.Fatalf("expected 1 filter, got %d", len(filters))
	}
	f0 := filters[0].(map[string]interface{})
	if f0["checkOn"] != "Response Status Code" || f0["filterType"] != "Greater Than" || f0["value"] != "499" {
		t.Errorf("filter wrong: %v", f0)
	}
}

func TestAttachMonitorSteps_IPAndHostname(t *testing.T) {
	t.Parallel()
	data := map[string]interface{}{}
	steps := []MonitorStep{
		{MonitorDestination: strPtr("192.168.1.1"), MonitorCriteria: MonitorCriteria{}},
		{MonitorDestination: strPtr("db.internal"), MonitorCriteria: MonitorCriteria{}},
	}
	if err := attachMonitorSteps(data, steps, nil); err != nil {
		t.Fatal(err)
	}
	arr := data["monitorSteps"].(map[string]interface{})["value"].(map[string]interface{})["monitorStepsInstanceArray"].([]map[string]interface{})
	if got := arr[0]["value"].(map[string]interface{})["monitorDestination"].(map[string]interface{})["_type"]; got != "IP" {
		t.Errorf("expected IP envelope, got %v", got)
	}
	if got := arr[1]["value"].(map[string]interface{})["monitorDestination"].(map[string]interface{})["_type"]; got != "Hostname" {
		t.Errorf("expected Hostname envelope, got %v", got)
	}
}

func TestAttachMonitorSteps_Nil(t *testing.T) {
	t.Parallel()
	data := map[string]interface{}{"name": "keep"}
	if err := attachMonitorSteps(data, nil, nil); err != nil {
		t.Fatal(err)
	}
	if _, present := data["monitorSteps"]; present {
		t.Error("nil steps + nil default should leave monitorSteps unset")
	}
}

func TestAttachLabels_ManyToMany(t *testing.T) {
	t.Parallel()
	data := map[string]interface{}{}
	attachLabels(data, []string{"lbl-a", "lbl-b"})
	labels := data["labels"].([]map[string]interface{})
	if len(labels) != 2 {
		t.Fatalf("expected 2 refs, got %d", len(labels))
	}
	if labels[0]["_id"] != "lbl-a" || labels[1]["_id"] != "lbl-b" {
		t.Errorf("labels shape wrong: %v", labels)
	}
}

func TestAttachLabels_NilSkipsEmptyExplicitlyClears(t *testing.T) {
	t.Parallel()
	// nil = don't touch
	d1 := map[string]interface{}{"name": "x"}
	attachLabels(d1, nil)
	if _, p := d1["labels"]; p {
		t.Error("nil labels should leave key absent")
	}
	// empty = explicit clear
	d2 := map[string]interface{}{}
	attachLabels(d2, []string{})
	got, ok := d2["labels"].([]map[string]interface{})
	if !ok || len(got) != 0 {
		t.Errorf("empty labels should produce empty array, got %v", d2["labels"])
	}
}

func TestAttachIDRefs_Monitors(t *testing.T) {
	t.Parallel()
	data := map[string]interface{}{}
	attachIDRefs(data, "monitors", []string{"m1", "m2", "m3"})
	monitors := data["monitors"].([]map[string]interface{})
	if len(monitors) != 3 {
		t.Fatalf("expected 3 refs, got %d", len(monitors))
	}
	for i, id := range []string{"m1", "m2", "m3"} {
		if monitors[i]["_id"] != id {
			t.Errorf("refs[%d]._id = %v, want %s", i, monitors[i]["_id"], id)
		}
	}
}

func TestInferDestinationType(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   string
		want string
	}{
		{"https://example.com", "URL"},
		{"http://example.com/path", "URL"},
		{"HTTPS://UPPER.COM", "URL"},
		{"192.168.1.1", "IP"},
		{"::1", "IP"},
		{"fe80::1", "IP"},
		{"db.internal", "Hostname"},
		{"localhost", "Hostname"},
	}
	for _, tc := range cases {
		if got := inferDestinationType(tc.in); got != tc.want {
			t.Errorf("%q: got %s, want %s", tc.in, got, tc.want)
		}
	}
}
