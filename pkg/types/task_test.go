package types

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestTaskYAMLMarshalling(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	// Create a task with both ID fields and pointer fields populated
	task := &Task{
		ID:            "test-1",
		Name:          "Test Task",
		Summary:       "Summary",
		Description:   "Description",
		Tags:          []string{"tag1", "tag2"},
		DependencyIDs: []string{"dep-1", "dep-2"},
		DependentIDs:  []string{"dependent-1"},
		// Populate pointer fields - these should NOT be marshalled
		Dependencies: []*Task{
			{ID: "should-not-appear-1"},
			{ID: "should-not-appear-2"},
		},
		Dependents: []*Task{
			{ID: "should-not-appear-3"},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Marshal to YAML
	data, err := yaml.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal task to YAML: %v", err)
	}

	yamlStr := string(data)

	// Verify that DependencyIDs and DependentIDs are present
	if !strings.Contains(yamlStr, "dependencies:") {
		t.Error("Expected 'dependencies:' field in YAML output")
	}
	if !strings.Contains(yamlStr, "dep-1") {
		t.Error("Expected 'dep-1' in YAML output")
	}
	if !strings.Contains(yamlStr, "dep-2") {
		t.Error("Expected 'dep-2' in YAML output")
	}
	if !strings.Contains(yamlStr, "dependents:") {
		t.Error("Expected 'dependents:' field in YAML output")
	}
	if !strings.Contains(yamlStr, "dependent-1") {
		t.Error("Expected 'dependent-1' in YAML output")
	}

	// Verify that pointer field values are NOT present
	if strings.Contains(yamlStr, "should-not-appear") {
		t.Errorf("Pointer field values should not appear in YAML output. Got:\n%s", yamlStr)
	}

	// Unmarshal back and verify
	var unmarshalled Task
	if err := yaml.Unmarshal(data, &unmarshalled); err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// Verify ID fields were preserved
	if len(unmarshalled.DependencyIDs) != 2 {
		t.Errorf("Expected 2 dependency IDs, got %d", len(unmarshalled.DependencyIDs))
	}
	if len(unmarshalled.DependentIDs) != 1 {
		t.Errorf("Expected 1 dependent ID, got %d", len(unmarshalled.DependentIDs))
	}

	// Verify pointer fields are nil/empty after unmarshalling
	if unmarshalled.Dependencies != nil {
		t.Errorf("Dependencies pointer field should be nil after unmarshal, got %v", unmarshalled.Dependencies)
	}
	if unmarshalled.Dependents != nil {
		t.Errorf("Dependents pointer field should be nil after unmarshal, got %v", unmarshalled.Dependents)
	}
}

func TestTaskJSONMarshalling(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	// Create a task with both ID fields and pointer fields populated
	task := &Task{
		ID:            "test-1",
		Name:          "Test Task",
		Summary:       "Summary",
		Description:   "Description",
		Tags:          []string{"tag1", "tag2"},
		DependencyIDs: []string{"dep-1", "dep-2"},
		DependentIDs:  []string{"dependent-1"},
		// Populate pointer fields - these should NOT be marshalled
		Dependencies: []*Task{
			{ID: "should-not-appear-1"},
			{ID: "should-not-appear-2"},
		},
		Dependents: []*Task{
			{ID: "should-not-appear-3"},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Marshal to JSON
	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal task to JSON: %v", err)
	}

	jsonStr := string(data)

	// Verify that DependencyIDs and DependentIDs are present with JSON field names
	if !strings.Contains(jsonStr, `"dependencies"`) {
		t.Error("Expected 'dependencies' field in JSON output")
	}
	if !strings.Contains(jsonStr, `"dep-1"`) {
		t.Error("Expected 'dep-1' in JSON output")
	}
	if !strings.Contains(jsonStr, `"dependents"`) {
		t.Error("Expected 'dependents' field in JSON output")
	}

	// Verify that pointer field values are NOT present
	if strings.Contains(jsonStr, "should-not-appear") {
		t.Errorf("Pointer field values should not appear in JSON output. Got:\n%s", jsonStr)
	}

	// Unmarshal back and verify
	var unmarshalled Task
	if err := json.Unmarshal(data, &unmarshalled); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify ID fields were preserved
	if len(unmarshalled.DependencyIDs) != 2 {
		t.Errorf("Expected 2 dependency IDs, got %d", len(unmarshalled.DependencyIDs))
	}
	if len(unmarshalled.DependentIDs) != 1 {
		t.Errorf("Expected 1 dependent ID, got %d", len(unmarshalled.DependentIDs))
	}

	// Verify pointer fields are nil/empty after unmarshalling
	if unmarshalled.Dependencies != nil {
		t.Errorf("Dependencies pointer field should be nil after unmarshal, got %v", unmarshalled.Dependencies)
	}
	if unmarshalled.Dependents != nil {
		t.Errorf("Dependents pointer field should be nil after unmarshal, got %v", unmarshalled.Dependents)
	}
}
