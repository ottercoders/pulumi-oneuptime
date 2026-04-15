package resources

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/ottercoders/pulumi-oneuptime/provider/client"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// IsNotFound re-exports client.IsNotFound for use in resources.
var IsNotFound = client.IsNotFound

// ToMap converts a struct to map[string]interface{} using json tags.
// Nil pointer fields and zero-value fields with omitempty are excluded.
func ToMap(v interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshalling to map: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling to map: %w", err)
	}
	return result, nil
}

// FromMap populates a struct from a map using json tags.
func FromMap(m map[string]interface{}, target interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshalling map: %w", err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("unmarshalling to struct: %w", err)
	}
	return nil
}

// SelectFields returns a map of {jsonFieldName: true} for all json-tagged
// fields in the given struct type, suitable for OneUptime's select body.
func SelectFields(v interface{}) map[string]bool {
	result := make(map[string]bool)
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	collectFields(t, result)
	return result
}

func collectFields(t reflect.Type, result map[string]bool) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// Handle embedded structs
		if field.Anonymous {
			ft := field.Type
			if ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
			}
			if ft.Kind() == reflect.Struct {
				collectFields(ft, result)
				continue
			}
		}
		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		// Parse the tag name (before any comma)
		name := tag
		for i, c := range tag {
			if c == ',' {
				name = tag[:i]
				break
			}
		}
		if name != "" {
			result[name] = true
		}
	}
}

// ResolveProjectID resolves the project ID from the resource arg or provider config.
func ResolveProjectID(resourceProjectID *string, configProjectID *string) (string, error) {
	if resourceProjectID != nil && *resourceProjectID != "" {
		return *resourceProjectID, nil
	}
	if configProjectID != nil && *configProjectID != "" {
		return *configProjectID, nil
	}
	return "", infer.ProviderErrorf("projectId is required: set it on the resource or in provider config or ONEUPTIME_PROJECT_ID env var")
}
