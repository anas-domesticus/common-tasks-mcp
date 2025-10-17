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
		ID:                     "test-1",
		Name:                   "Test Task",
		Summary:                "Summary",
		Description:            "Description",
		Tags:                   []string{"tag1", "tag2"},
		PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
		DownstreamRequiredIDs:  []string{"required-1"},
		DownstreamSuggestedIDs: []string{"suggested-1", "suggested-2"},
		// Populate pointer fields - these should NOT be marshalled
		Prerequisites: []*Task{
			{ID: "should-not-appear-1"},
			{ID: "should-not-appear-2"},
		},
		DownstreamRequired: []*Task{
			{ID: "should-not-appear-3"},
		},
		DownstreamSuggested: []*Task{
			{ID: "should-not-appear-4"},
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

	// Verify that PrerequisiteIDs, DownstreamRequiredIDs, and DownstreamSuggestedIDs are present
	if !strings.Contains(yamlStr, "prerequisites:") {
		t.Error("Expected 'prerequisites:' field in YAML output")
	}
	if !strings.Contains(yamlStr, "prereq-1") {
		t.Error("Expected 'prereq-1' in YAML output")
	}
	if !strings.Contains(yamlStr, "prereq-2") {
		t.Error("Expected 'prereq-2' in YAML output")
	}
	if !strings.Contains(yamlStr, "downstream_required:") {
		t.Error("Expected 'downstream_required:' field in YAML output")
	}
	if !strings.Contains(yamlStr, "required-1") {
		t.Error("Expected 'required-1' in YAML output")
	}
	if !strings.Contains(yamlStr, "downstream_suggested:") {
		t.Error("Expected 'downstream_suggested:' field in YAML output")
	}
	if !strings.Contains(yamlStr, "suggested-1") {
		t.Error("Expected 'suggested-1' in YAML output")
	}
	if !strings.Contains(yamlStr, "suggested-2") {
		t.Error("Expected 'suggested-2' in YAML output")
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
	if len(unmarshalled.PrerequisiteIDs) != 2 {
		t.Errorf("Expected 2 prerequisite IDs, got %d", len(unmarshalled.PrerequisiteIDs))
	}
	if len(unmarshalled.DownstreamRequiredIDs) != 1 {
		t.Errorf("Expected 1 downstream required ID, got %d", len(unmarshalled.DownstreamRequiredIDs))
	}
	if len(unmarshalled.DownstreamSuggestedIDs) != 2 {
		t.Errorf("Expected 2 downstream suggested IDs, got %d", len(unmarshalled.DownstreamSuggestedIDs))
	}

	// Verify pointer fields are nil/empty after unmarshalling
	if unmarshalled.Prerequisites != nil {
		t.Errorf("Prerequisites pointer field should be nil after unmarshal, got %v", unmarshalled.Prerequisites)
	}
	if unmarshalled.DownstreamRequired != nil {
		t.Errorf("DownstreamRequired pointer field should be nil after unmarshal, got %v", unmarshalled.DownstreamRequired)
	}
	if unmarshalled.DownstreamSuggested != nil {
		t.Errorf("DownstreamSuggested pointer field should be nil after unmarshal, got %v", unmarshalled.DownstreamSuggested)
	}
}

func TestTaskJSONMarshalling(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	// Create a task with both ID fields and pointer fields populated
	task := &Task{
		ID:                     "test-1",
		Name:                   "Test Task",
		Summary:                "Summary",
		Description:            "Description",
		Tags:                   []string{"tag1", "tag2"},
		PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
		DownstreamRequiredIDs:  []string{"required-1"},
		DownstreamSuggestedIDs: []string{"suggested-1", "suggested-2"},
		// Populate pointer fields - these should NOT be marshalled
		Prerequisites: []*Task{
			{ID: "should-not-appear-1"},
			{ID: "should-not-appear-2"},
		},
		DownstreamRequired: []*Task{
			{ID: "should-not-appear-3"},
		},
		DownstreamSuggested: []*Task{
			{ID: "should-not-appear-4"},
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

	// Verify that PrerequisiteIDs, DownstreamRequiredIDs, and DownstreamSuggestedIDs are present with JSON field names
	if !strings.Contains(jsonStr, `"prerequisites"`) {
		t.Error("Expected 'prerequisites' field in JSON output")
	}
	if !strings.Contains(jsonStr, `"prereq-1"`) {
		t.Error("Expected 'prereq-1' in JSON output")
	}
	if !strings.Contains(jsonStr, `"prereq-2"`) {
		t.Error("Expected 'prereq-2' in JSON output")
	}
	if !strings.Contains(jsonStr, `"downstream_required"`) {
		t.Error("Expected 'downstream_required' field in JSON output")
	}
	if !strings.Contains(jsonStr, `"required-1"`) {
		t.Error("Expected 'required-1' in JSON output")
	}
	if !strings.Contains(jsonStr, `"downstream_suggested"`) {
		t.Error("Expected 'downstream_suggested' field in JSON output")
	}
	if !strings.Contains(jsonStr, `"suggested-1"`) {
		t.Error("Expected 'suggested-1' in JSON output")
	}
	if !strings.Contains(jsonStr, `"suggested-2"`) {
		t.Error("Expected 'suggested-2' in JSON output")
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
	if len(unmarshalled.PrerequisiteIDs) != 2 {
		t.Errorf("Expected 2 prerequisite IDs, got %d", len(unmarshalled.PrerequisiteIDs))
	}
	if len(unmarshalled.DownstreamRequiredIDs) != 1 {
		t.Errorf("Expected 1 downstream required ID, got %d", len(unmarshalled.DownstreamRequiredIDs))
	}
	if len(unmarshalled.DownstreamSuggestedIDs) != 2 {
		t.Errorf("Expected 2 downstream suggested IDs, got %d", len(unmarshalled.DownstreamSuggestedIDs))
	}

	// Verify pointer fields are nil/empty after unmarshalling
	if unmarshalled.Prerequisites != nil {
		t.Errorf("Prerequisites pointer field should be nil after unmarshal, got %v", unmarshalled.Prerequisites)
	}
	if unmarshalled.DownstreamRequired != nil {
		t.Errorf("DownstreamRequired pointer field should be nil after unmarshal, got %v", unmarshalled.DownstreamRequired)
	}
	if unmarshalled.DownstreamSuggested != nil {
		t.Errorf("DownstreamSuggested pointer field should be nil after unmarshal, got %v", unmarshalled.DownstreamSuggested)
	}
}

func TestTaskEquals(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	// Create a base task for comparisons
	baseTask := &Task{
		ID:                     "test-1",
		Name:                   "Test Task",
		Summary:                "Summary",
		Description:            "Description",
		Tags:                   []string{"tag1", "tag2"},
		PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
		DownstreamRequiredIDs:  []string{"required-1"},
		DownstreamSuggestedIDs: []string{"suggested-1"},
		CreatedAt:              now,
		UpdatedAt:              now,
	}

	tests := []struct {
		name        string
		task1       *Task
		task2       *Task
		shouldEqual bool
	}{
		{
			name:        "same task should equal itself",
			task1:       baseTask,
			task2:       baseTask,
			shouldEqual: true,
		},
		{
			name:  "identical task should be equal",
			task1: baseTask,
			task2: &Task{
				ID:                     "test-1",
				Name:                   "Test Task",
				Summary:                "Summary",
				Description:            "Description",
				Tags:                   []string{"tag1", "tag2"},
				PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
				DownstreamRequiredIDs:  []string{"required-1"},
				DownstreamSuggestedIDs: []string{"suggested-1"},
				CreatedAt:              now,
				UpdatedAt:              now,
			},
			shouldEqual: true,
		},
		{
			name:        "both nil should be equal",
			task1:       nil,
			task2:       nil,
			shouldEqual: true,
		},
		{
			name:        "task and nil should not be equal",
			task1:       baseTask,
			task2:       nil,
			shouldEqual: false,
		},
		{
			name:        "nil and task should not be equal",
			task1:       nil,
			task2:       baseTask,
			shouldEqual: false,
		},
		{
			name:  "different ID",
			task1: baseTask,
			task2: &Task{
				ID:                     "test-2",
				Name:                   "Test Task",
				Summary:                "Summary",
				Description:            "Description",
				Tags:                   []string{"tag1", "tag2"},
				PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
				DownstreamRequiredIDs:  []string{"required-1"},
				DownstreamSuggestedIDs: []string{"suggested-1"},
				CreatedAt:              now,
				UpdatedAt:              now,
			},
			shouldEqual: false,
		},
		{
			name:  "different Name",
			task1: baseTask,
			task2: &Task{
				ID:                     "test-1",
				Name:                   "Different Name",
				Summary:                "Summary",
				Description:            "Description",
				Tags:                   []string{"tag1", "tag2"},
				PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
				DownstreamRequiredIDs:  []string{"required-1"},
				DownstreamSuggestedIDs: []string{"suggested-1"},
				CreatedAt:              now,
				UpdatedAt:              now,
			},
			shouldEqual: false,
		},
		{
			name:  "different Summary",
			task1: baseTask,
			task2: &Task{
				ID:                     "test-1",
				Name:                   "Test Task",
				Summary:                "Different Summary",
				Description:            "Description",
				Tags:                   []string{"tag1", "tag2"},
				PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
				DownstreamRequiredIDs:  []string{"required-1"},
				DownstreamSuggestedIDs: []string{"suggested-1"},
				CreatedAt:              now,
				UpdatedAt:              now,
			},
			shouldEqual: false,
		},
		{
			name:  "different Description",
			task1: baseTask,
			task2: &Task{
				ID:                     "test-1",
				Name:                   "Test Task",
				Summary:                "Summary",
				Description:            "Different Description",
				Tags:                   []string{"tag1", "tag2"},
				PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
				DownstreamRequiredIDs:  []string{"required-1"},
				DownstreamSuggestedIDs: []string{"suggested-1"},
				CreatedAt:              now,
				UpdatedAt:              now,
			},
			shouldEqual: false,
		},
		{
			name:  "different CreatedAt",
			task1: baseTask,
			task2: &Task{
				ID:                     "test-1",
				Name:                   "Test Task",
				Summary:                "Summary",
				Description:            "Description",
				Tags:                   []string{"tag1", "tag2"},
				PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
				DownstreamRequiredIDs:  []string{"required-1"},
				DownstreamSuggestedIDs: []string{"suggested-1"},
				CreatedAt:              now.Add(time.Hour),
				UpdatedAt:              now,
			},
			shouldEqual: false,
		},
		{
			name:  "different UpdatedAt",
			task1: baseTask,
			task2: &Task{
				ID:                     "test-1",
				Name:                   "Test Task",
				Summary:                "Summary",
				Description:            "Description",
				Tags:                   []string{"tag1", "tag2"},
				PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
				DownstreamRequiredIDs:  []string{"required-1"},
				DownstreamSuggestedIDs: []string{"suggested-1"},
				CreatedAt:              now,
				UpdatedAt:              now.Add(time.Hour),
			},
			shouldEqual: false,
		},
		{
			name:  "different Tags length",
			task1: baseTask,
			task2: &Task{
				ID:                     "test-1",
				Name:                   "Test Task",
				Summary:                "Summary",
				Description:            "Description",
				Tags:                   []string{"tag1"},
				PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
				DownstreamRequiredIDs:  []string{"required-1"},
				DownstreamSuggestedIDs: []string{"suggested-1"},
				CreatedAt:              now,
				UpdatedAt:              now,
			},
			shouldEqual: false,
		},
		{
			name:  "different Tags values",
			task1: baseTask,
			task2: &Task{
				ID:                     "test-1",
				Name:                   "Test Task",
				Summary:                "Summary",
				Description:            "Description",
				Tags:                   []string{"tag1", "tag3"},
				PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
				DownstreamRequiredIDs:  []string{"required-1"},
				DownstreamSuggestedIDs: []string{"suggested-1"},
				CreatedAt:              now,
				UpdatedAt:              now,
			},
			shouldEqual: false,
		},
		{
			name:  "different PrerequisiteIDs length",
			task1: baseTask,
			task2: &Task{
				ID:                     "test-1",
				Name:                   "Test Task",
				Summary:                "Summary",
				Description:            "Description",
				Tags:                   []string{"tag1", "tag2"},
				PrerequisiteIDs:        []string{"prereq-1"},
				DownstreamRequiredIDs:  []string{"required-1"},
				DownstreamSuggestedIDs: []string{"suggested-1"},
				CreatedAt:              now,
				UpdatedAt:              now,
			},
			shouldEqual: false,
		},
		{
			name:  "different PrerequisiteIDs values",
			task1: baseTask,
			task2: &Task{
				ID:                     "test-1",
				Name:                   "Test Task",
				Summary:                "Summary",
				Description:            "Description",
				Tags:                   []string{"tag1", "tag2"},
				PrerequisiteIDs:        []string{"prereq-1", "prereq-3"},
				DownstreamRequiredIDs:  []string{"required-1"},
				DownstreamSuggestedIDs: []string{"suggested-1"},
				CreatedAt:              now,
				UpdatedAt:              now,
			},
			shouldEqual: false,
		},
		{
			name:  "different DownstreamRequiredIDs length",
			task1: baseTask,
			task2: &Task{
				ID:                     "test-1",
				Name:                   "Test Task",
				Summary:                "Summary",
				Description:            "Description",
				Tags:                   []string{"tag1", "tag2"},
				PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
				DownstreamRequiredIDs:  []string{},
				DownstreamSuggestedIDs: []string{"suggested-1"},
				CreatedAt:              now,
				UpdatedAt:              now,
			},
			shouldEqual: false,
		},
		{
			name:  "different DownstreamRequiredIDs values",
			task1: baseTask,
			task2: &Task{
				ID:                     "test-1",
				Name:                   "Test Task",
				Summary:                "Summary",
				Description:            "Description",
				Tags:                   []string{"tag1", "tag2"},
				PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
				DownstreamRequiredIDs:  []string{"required-2"},
				DownstreamSuggestedIDs: []string{"suggested-1"},
				CreatedAt:              now,
				UpdatedAt:              now,
			},
			shouldEqual: false,
		},
		{
			name:  "different DownstreamSuggestedIDs length",
			task1: baseTask,
			task2: &Task{
				ID:                     "test-1",
				Name:                   "Test Task",
				Summary:                "Summary",
				Description:            "Description",
				Tags:                   []string{"tag1", "tag2"},
				PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
				DownstreamRequiredIDs:  []string{"required-1"},
				DownstreamSuggestedIDs: []string{},
				CreatedAt:              now,
				UpdatedAt:              now,
			},
			shouldEqual: false,
		},
		{
			name:  "different DownstreamSuggestedIDs values",
			task1: baseTask,
			task2: &Task{
				ID:                     "test-1",
				Name:                   "Test Task",
				Summary:                "Summary",
				Description:            "Description",
				Tags:                   []string{"tag1", "tag2"},
				PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
				DownstreamRequiredIDs:  []string{"required-1"},
				DownstreamSuggestedIDs: []string{"suggested-2"},
				CreatedAt:              now,
				UpdatedAt:              now,
			},
			shouldEqual: false,
		},
		{
			name: "empty slices vs nil slices should be equal",
			task1: &Task{
				ID:                     "test-empty",
				Name:                   "Empty Slices",
				Summary:                "Summary",
				Description:            "Description",
				Tags:                   []string{},
				PrerequisiteIDs:        []string{},
				DownstreamRequiredIDs:  []string{},
				DownstreamSuggestedIDs: []string{},
				CreatedAt:              now,
				UpdatedAt:              now,
			},
			task2: &Task{
				ID:                     "test-empty",
				Name:                   "Empty Slices",
				Summary:                "Summary",
				Description:            "Description",
				Tags:                   nil,
				PrerequisiteIDs:        nil,
				DownstreamRequiredIDs:  nil,
				DownstreamSuggestedIDs: nil,
				CreatedAt:              now,
				UpdatedAt:              now,
			},
			shouldEqual: true,
		},
		{
			name:  "pointer fields should not affect equality",
			task1: baseTask,
			task2: &Task{
				ID:                     "test-1",
				Name:                   "Test Task",
				Summary:                "Summary",
				Description:            "Description",
				Tags:                   []string{"tag1", "tag2"},
				PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
				DownstreamRequiredIDs:  []string{"required-1"},
				DownstreamSuggestedIDs: []string{"suggested-1"},
				Prerequisites:          []*Task{{ID: "different"}},
				DownstreamRequired:     []*Task{{ID: "different"}},
				DownstreamSuggested:    []*Task{{ID: "different"}},
				CreatedAt:              now,
				UpdatedAt:              now,
			},
			shouldEqual: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.task1.Equals(tt.task2)
			if got != tt.shouldEqual {
				t.Errorf("Equals() = %v, want %v", got, tt.shouldEqual)
			}
		})
	}
}
