package resources

import (
	"fmt"
	"net"
	"strings"
)

// OneUptime stores and round-trips structured monitor config in a
// `{_type, value}` envelope. The inner `value` is itself a map whose nested
// fields may be further enveloped (monitorDestination/Port at the leaf,
// MonitorCriteria/MonitorCriteriaInstance in the middle). The default Pulumi
// ToMap can't express that shape via json tags alone, so the fields involved
// carry `json:"-"` on the structs and the helpers below assemble the
// enveloped wire form after ToMap has handled the flat fields.

// attachMonitorSteps writes the MonitorSteps envelope to data["monitorSteps"]
// and data["defaultMonitorStatusId"]. The envelope matches what OneUptime's
// MonitorSteps.toJSON emits: the outer MonitorSteps wrapper, each MonitorStep
// wrapped, each MonitorCriteria/Instance wrapped, and the destination/port
// leaves wrapped as URL/Hostname/IP/Port values.
//
// Passing a nil steps slice leaves both keys untouched so the server treats
// it as "no change" (useful for Manual monitors and for Update when steps are
// unchanged).
func attachMonitorSteps(data map[string]interface{}, steps []MonitorStep, defaultStatusID *string) error {
	if steps == nil && defaultStatusID == nil {
		return nil
	}

	instances := make([]map[string]interface{}, 0, len(steps))
	for i, step := range steps {
		stepMap, err := monitorStepToEnvelope(step)
		if err != nil {
			return fmt.Errorf("step %d: %w", i, err)
		}
		instances = append(instances, stepMap)
	}

	envValue := map[string]interface{}{
		"monitorStepsInstanceArray": instances,
	}
	if defaultStatusID != nil && *defaultStatusID != "" {
		envValue["defaultMonitorStatusId"] = *defaultStatusID
	}

	data["monitorSteps"] = map[string]interface{}{
		"_type": "MonitorSteps",
		"value": envValue,
	}
	return nil
}

// attachLabels rewrites a []string of label IDs into the ManyToMany write
// shape OneUptime expects (`[{_id: "..."}]`) at data["labels"]. Nil input is
// a no-op; an empty slice explicitly clears labels on the server.
func attachLabels(data map[string]interface{}, labels []string) {
	if labels == nil {
		return
	}
	refs := make([]map[string]interface{}, 0, len(labels))
	for _, id := range labels {
		refs = append(refs, map[string]interface{}{"_id": id})
	}
	data["labels"] = refs
}

// attachIDRefs is the generic counterpart to attachLabels for any ManyToMany
// column written as an array of `{_id: ...}` references (e.g.
// MonitorSecret.monitors).
func attachIDRefs(data map[string]interface{}, key string, ids []string) {
	if ids == nil {
		return
	}
	refs := make([]map[string]interface{}, 0, len(ids))
	for _, id := range ids {
		refs = append(refs, map[string]interface{}{"_id": id})
	}
	data[key] = refs
}

// monitorStepToEnvelope builds the `{_type: "MonitorStep", value: {...}}` wire
// form for a single step, including the nested criteria and destination
// envelopes.
func monitorStepToEnvelope(step MonitorStep) (map[string]interface{}, error) {
	// ToMap serializes the flat json-tagged fields (requestType, headers,
	// sub-monitor configs, etc.). The json:"-" fields — monitorDestination,
	// monitorDestinationPort, monitorCriteria, id — are re-attached below.
	value, err := ToMap(step)
	if err != nil {
		return nil, err
	}

	if step.ID != "" {
		value["id"] = step.ID
	}

	if step.MonitorDestination != nil && *step.MonitorDestination != "" {
		value["monitorDestination"] = map[string]interface{}{
			"_type": inferDestinationType(*step.MonitorDestination),
			"value": *step.MonitorDestination,
		}
	}

	if step.MonitorDestinationPort != nil {
		value["monitorDestinationPort"] = map[string]interface{}{
			"_type": "Port",
			"value": fmt.Sprintf("%d", *step.MonitorDestinationPort),
		}
	}

	value["monitorCriteria"] = monitorCriteriaToEnvelope(step.MonitorCriteria)

	return map[string]interface{}{
		"_type": "MonitorStep",
		"value": value,
	}, nil
}

func monitorCriteriaToEnvelope(c MonitorCriteria) map[string]interface{} {
	instances := make([]map[string]interface{}, 0, len(c.CriteriaInstances))
	for _, inst := range c.CriteriaInstances {
		instances = append(instances, monitorCriteriaInstanceToEnvelope(inst))
	}
	return map[string]interface{}{
		"_type": "MonitorCriteria",
		"value": map[string]interface{}{
			"monitorCriteriaInstanceArray": instances,
		},
	}
}

func monitorCriteriaInstanceToEnvelope(inst MonitorCriteriaInstance) map[string]interface{} {
	value, err := ToMap(inst)
	if err != nil {
		// ToMap only fails on json.Marshal errors; MonitorCriteriaInstance
		// contains nothing that can fail (no channels, no functions).
		// Fall back to an empty value map rather than propagating, since
		// this helper is called from a hot path that doesn't return error.
		value = map[string]interface{}{}
	}
	// Always include id so the server can match criteria across updates.
	if inst.ID != "" {
		value["id"] = inst.ID
	}
	// OneUptime's TS interface declares incidents/alerts as non-optional
	// arrays — the keys must be present on the wire even when empty. The
	// Go struct uses omitempty so nil/empty slices drop out of ToMap;
	// re-attach here to match the server contract.
	if _, ok := value["incidents"]; !ok {
		value["incidents"] = []interface{}{}
	}
	if _, ok := value["alerts"]; !ok {
		value["alerts"] = []interface{}{}
	}
	return map[string]interface{}{
		"_type": "MonitorCriteriaInstance",
		"value": value,
	}
}

// inferDestinationType picks OneUptime's _type discriminator for a
// monitorDestination leaf based on the shape of the string. HTTPS/HTTP →
// "URL"; numeric dotted → "IP"; anything else → "Hostname".
func inferDestinationType(dest string) string {
	lower := strings.ToLower(dest)
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
		return "URL"
	}
	if net.ParseIP(dest) != nil {
		return "IP"
	}
	return "Hostname"
}
