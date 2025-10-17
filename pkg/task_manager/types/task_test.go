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

func TestCheckEdgeConsistency(t *testing.T) {
	tests := []struct {
		name      string
		tasks     []*Task
		ids       []string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "both empty slices",
			tasks:     []*Task{},
			ids:       []string{},
			wantError: false,
		},
		{
			name:      "both nil slices",
			tasks:     nil,
			ids:       nil,
			wantError: false,
		},
		{
			name: "single task and ID match",
			tasks: []*Task{
				{ID: "task-1"},
			},
			ids:       []string{"task-1"},
			wantError: false,
		},
		{
			name: "multiple tasks and IDs match (same order)",
			tasks: []*Task{
				{ID: "task-1"},
				{ID: "task-2"},
				{ID: "task-3"},
			},
			ids:       []string{"task-1", "task-2", "task-3"},
			wantError: false,
		},
		{
			name: "multiple tasks and IDs match (different order)",
			tasks: []*Task{
				{ID: "task-3"},
				{ID: "task-1"},
				{ID: "task-2"},
			},
			ids:       []string{"task-1", "task-2", "task-3"},
			wantError: false,
		},
		{
			name: "length mismatch - more tasks",
			tasks: []*Task{
				{ID: "task-1"},
				{ID: "task-2"},
			},
			ids:       []string{"task-1"},
			wantError: true,
			errorMsg:  "length mismatch: 2 tasks but 1 IDs",
		},
		{
			name: "length mismatch - more IDs",
			tasks: []*Task{
				{ID: "task-1"},
			},
			ids:       []string{"task-1", "task-2"},
			wantError: true,
			errorMsg:  "length mismatch: 1 tasks but 2 IDs",
		},
		{
			name: "duplicate ID in string slice",
			tasks: []*Task{
				{ID: "task-1"},
				{ID: "task-2"},
			},
			ids:       []string{"task-1", "task-1"},
			wantError: true,
			errorMsg:  "duplicate ID in string slice: task-1",
		},
		{
			name: "duplicate task ID in task slice",
			tasks: []*Task{
				{ID: "task-1"},
				{ID: "task-1"},
			},
			ids:       []string{"task-1", "task-2"},
			wantError: true,
			errorMsg:  "duplicate task ID in task slice: task-1",
		},
		{
			name: "task ID not found in string slice",
			tasks: []*Task{
				{ID: "task-1"},
				{ID: "task-2"},
			},
			ids:       []string{"task-1", "task-3"},
			wantError: true,
			errorMsg:  "task ID task-2 not found in string slice",
		},
		{
			name: "mismatched IDs - different sets",
			tasks: []*Task{
				{ID: "task-1"},
				{ID: "task-3"},
			},
			ids:       []string{"task-1", "task-2"},
			wantError: true,
			errorMsg:  "task ID task-3 not found in string slice",
		},
		{
			name: "nil task in slice",
			tasks: []*Task{
				{ID: "task-1"},
				nil,
			},
			ids:       []string{"task-1", "task-2"},
			wantError: true,
			errorMsg:  "nil task at index 1",
		},
		{
			name: "nil task at first position",
			tasks: []*Task{
				nil,
				{ID: "task-2"},
			},
			ids:       []string{"task-1", "task-2"},
			wantError: true,
			errorMsg:  "nil task at index 0",
		},
		{
			name: "complex valid case with many tasks",
			tasks: []*Task{
				{ID: "alpha"},
				{ID: "beta"},
				{ID: "gamma"},
				{ID: "delta"},
				{ID: "epsilon"},
			},
			ids:       []string{"epsilon", "gamma", "alpha", "delta", "beta"},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkEdgeConsistency(tt.tasks, tt.ids)

			if tt.wantError {
				if err == nil {
					t.Errorf("checkEdgeConsistency() expected error but got nil")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("checkEdgeConsistency() error = %q, want %q", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("checkEdgeConsistency() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestTask_CheckEdgeConsistency(t *testing.T) {
	tests := []struct {
		name             string
		task             *Task
		wantError        bool
		errorContains    string
		errorContainsAll []string
	}{
		{
			name:      "nil task returns nil",
			task:      nil,
			wantError: false,
		},
		{
			name: "all edges consistent",
			task: &Task{
				ID:                     "task-1",
				PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
				Prerequisites:          []*Task{{ID: "prereq-1"}, {ID: "prereq-2"}},
				DownstreamRequiredIDs:  []string{"req-1"},
				DownstreamRequired:     []*Task{{ID: "req-1"}},
				DownstreamSuggestedIDs: []string{"sug-1", "sug-2"},
				DownstreamSuggested:    []*Task{{ID: "sug-1"}, {ID: "sug-2"}},
			},
			wantError: false,
		},
		{
			name: "all edges consistent with different order",
			task: &Task{
				ID:                     "task-1",
				PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
				Prerequisites:          []*Task{{ID: "prereq-2"}, {ID: "prereq-1"}},
				DownstreamRequiredIDs:  []string{"req-1", "req-2"},
				DownstreamRequired:     []*Task{{ID: "req-2"}, {ID: "req-1"}},
				DownstreamSuggestedIDs: []string{"sug-1", "sug-2"},
				DownstreamSuggested:    []*Task{{ID: "sug-2"}, {ID: "sug-1"}},
			},
			wantError: false,
		},
		{
			name: "all edges empty",
			task: &Task{
				ID:                     "task-1",
				PrerequisiteIDs:        []string{},
				Prerequisites:          []*Task{},
				DownstreamRequiredIDs:  []string{},
				DownstreamRequired:     []*Task{},
				DownstreamSuggestedIDs: []string{},
				DownstreamSuggested:    []*Task{},
			},
			wantError: false,
		},
		{
			name: "all edges nil",
			task: &Task{
				ID:                     "task-1",
				PrerequisiteIDs:        nil,
				Prerequisites:          nil,
				DownstreamRequiredIDs:  nil,
				DownstreamRequired:     nil,
				DownstreamSuggestedIDs: nil,
				DownstreamSuggested:    nil,
			},
			wantError: false,
		},
		{
			name: "prerequisite edges inconsistent",
			task: &Task{
				ID:                     "task-1",
				PrerequisiteIDs:        []string{"prereq-1"},
				Prerequisites:          []*Task{{ID: "prereq-2"}},
				DownstreamRequiredIDs:  []string{"req-1"},
				DownstreamRequired:     []*Task{{ID: "req-1"}},
				DownstreamSuggestedIDs: []string{"sug-1"},
				DownstreamSuggested:    []*Task{{ID: "sug-1"}},
			},
			wantError:     true,
			errorContains: "prerequisite edges",
		},
		{
			name: "downstream required edges inconsistent",
			task: &Task{
				ID:                     "task-1",
				PrerequisiteIDs:        []string{"prereq-1"},
				Prerequisites:          []*Task{{ID: "prereq-1"}},
				DownstreamRequiredIDs:  []string{"req-1", "req-2"},
				DownstreamRequired:     []*Task{{ID: "req-1"}},
				DownstreamSuggestedIDs: []string{"sug-1"},
				DownstreamSuggested:    []*Task{{ID: "sug-1"}},
			},
			wantError:     true,
			errorContains: "downstream required edges",
		},
		{
			name: "downstream suggested edges inconsistent",
			task: &Task{
				ID:                     "task-1",
				PrerequisiteIDs:        []string{"prereq-1"},
				Prerequisites:          []*Task{{ID: "prereq-1"}},
				DownstreamRequiredIDs:  []string{"req-1"},
				DownstreamRequired:     []*Task{{ID: "req-1"}},
				DownstreamSuggestedIDs: []string{"sug-1"},
				DownstreamSuggested:    []*Task{{ID: "sug-2"}},
			},
			wantError:     true,
			errorContains: "downstream suggested edges",
		},
		{
			name: "prerequisite edges inconsistent - nil task",
			task: &Task{
				ID:                     "task-1",
				PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
				Prerequisites:          []*Task{{ID: "prereq-1"}, nil},
				DownstreamRequiredIDs:  []string{},
				DownstreamRequired:     []*Task{},
				DownstreamSuggestedIDs: []string{},
				DownstreamSuggested:    []*Task{},
			},
			wantError:     true,
			errorContains: "prerequisite edges",
		},
		{
			name: "multiple edges inconsistent - returns all errors",
			task: &Task{
				ID:                     "task-1",
				PrerequisiteIDs:        []string{"prereq-1"},
				Prerequisites:          []*Task{{ID: "wrong-id"}},
				DownstreamRequiredIDs:  []string{"req-1"},
				DownstreamRequired:     []*Task{{ID: "also-wrong"}},
				DownstreamSuggestedIDs: []string{"sug-1"},
				DownstreamSuggested:    []*Task{{ID: "wrong-again"}},
			},
			wantError:        true,
			errorContainsAll: []string{"prerequisite edges", "downstream required edges", "downstream suggested edges"},
		},
		{
			name: "all three edge types inconsistent with different errors",
			task: &Task{
				ID:                     "task-1",
				PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
				Prerequisites:          []*Task{{ID: "prereq-1"}}, // length mismatch
				DownstreamRequiredIDs:  []string{"req-1"},
				DownstreamRequired:     []*Task{nil}, // nil task
				DownstreamSuggestedIDs: []string{"sug-1"},
				DownstreamSuggested:    []*Task{{ID: "sug-2"}}, // ID mismatch
			},
			wantError:        true,
			errorContainsAll: []string{"prerequisite edges", "downstream required edges", "downstream suggested edges"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.CheckEdgeConsistency()

			if tt.wantError {
				if err == nil {
					t.Errorf("CheckEdgeConsistency() expected error but got nil")
				} else {
					errStr := err.Error()
					if tt.errorContains != "" && !strings.Contains(errStr, tt.errorContains) {
						t.Errorf("CheckEdgeConsistency() error = %q, want error containing %q", errStr, tt.errorContains)
					}
					for _, substr := range tt.errorContainsAll {
						if !strings.Contains(errStr, substr) {
							t.Errorf("CheckEdgeConsistency() error = %q, want error containing %q", errStr, substr)
						}
					}
				}
			} else {
				if err != nil {
					t.Errorf("CheckEdgeConsistency() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestTask_Clone(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	tests := []struct {
		name string
		task *Task
	}{
		{
			name: "nil task returns nil",
			task: nil,
		},
		{
			name: "task with all fields populated",
			task: &Task{
				ID:                     "test-1",
				Name:                   "Test Task",
				Summary:                "Summary",
				Description:            "Description",
				Tags:                   []string{"tag1", "tag2"},
				PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
				DownstreamRequiredIDs:  []string{"required-1"},
				DownstreamSuggestedIDs: []string{"suggested-1", "suggested-2"},
				// Populate pointer fields - these should NOT be cloned
				Prerequisites: []*Task{
					{ID: "prereq-1"},
					{ID: "prereq-2"},
				},
				DownstreamRequired: []*Task{
					{ID: "required-1"},
				},
				DownstreamSuggested: []*Task{
					{ID: "suggested-1"},
					{ID: "suggested-2"},
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name: "task with empty slices",
			task: &Task{
				ID:                     "test-2",
				Name:                   "Empty Task",
				Summary:                "",
				Description:            "",
				Tags:                   []string{},
				PrerequisiteIDs:        []string{},
				DownstreamRequiredIDs:  []string{},
				DownstreamSuggestedIDs: []string{},
				CreatedAt:              now,
				UpdatedAt:              now,
			},
		},
		{
			name: "task with nil slices",
			task: &Task{
				ID:                     "test-3",
				Name:                   "Nil Slices Task",
				Summary:                "Summary",
				Description:            "Description",
				Tags:                   nil,
				PrerequisiteIDs:        nil,
				DownstreamRequiredIDs:  nil,
				DownstreamSuggestedIDs: nil,
				CreatedAt:              now,
				UpdatedAt:              now,
			},
		},
		{
			name: "task with only ID and timestamps",
			task: &Task{
				ID:        "test-4",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clone := tt.task.Clone()

			// Test 1: nil task returns nil
			if tt.task == nil {
				if clone != nil {
					t.Errorf("Clone() of nil task should return nil, got %v", clone)
				}
				return
			}

			// Test 2: clone should not be nil for non-nil task
			if clone == nil {
				t.Fatal("Clone() returned nil for non-nil task")
			}

			// Test 3: clone should be a different instance (different pointer)
			if tt.task == clone {
				t.Error("Clone() returned the same pointer, expected a new instance")
			}

			// Test 4: clone should equal the original using Equals method
			if !tt.task.Equals(clone) {
				t.Error("Clone() should equal the original task")
			}

			// Test 5: verify all scalar fields are copied
			if clone.ID != tt.task.ID {
				t.Errorf("Clone().ID = %v, want %v", clone.ID, tt.task.ID)
			}
			if clone.Name != tt.task.Name {
				t.Errorf("Clone().Name = %v, want %v", clone.Name, tt.task.Name)
			}
			if clone.Summary != tt.task.Summary {
				t.Errorf("Clone().Summary = %v, want %v", clone.Summary, tt.task.Summary)
			}
			if clone.Description != tt.task.Description {
				t.Errorf("Clone().Description = %v, want %v", clone.Description, tt.task.Description)
			}
			if !clone.CreatedAt.Equal(tt.task.CreatedAt) {
				t.Errorf("Clone().CreatedAt = %v, want %v", clone.CreatedAt, tt.task.CreatedAt)
			}
			if !clone.UpdatedAt.Equal(tt.task.UpdatedAt) {
				t.Errorf("Clone().UpdatedAt = %v, want %v", clone.UpdatedAt, tt.task.UpdatedAt)
			}

			// Test 6: verify slice fields are deep copied (not the same slice reference)
			if tt.task.Tags != nil {
				if len(clone.Tags) != len(tt.task.Tags) {
					t.Errorf("Clone().Tags length = %d, want %d", len(clone.Tags), len(tt.task.Tags))
				}
				// Check that it's a different slice (different memory address)
				if len(tt.task.Tags) > 0 && &clone.Tags[0] == &tt.task.Tags[0] {
					t.Error("Clone().Tags should be a deep copy, not a reference to the same slice")
				}
				// Check values are the same
				for i := range tt.task.Tags {
					if clone.Tags[i] != tt.task.Tags[i] {
						t.Errorf("Clone().Tags[%d] = %v, want %v", i, clone.Tags[i], tt.task.Tags[i])
					}
				}
			}

			if tt.task.PrerequisiteIDs != nil {
				if len(clone.PrerequisiteIDs) != len(tt.task.PrerequisiteIDs) {
					t.Errorf("Clone().PrerequisiteIDs length = %d, want %d", len(clone.PrerequisiteIDs), len(tt.task.PrerequisiteIDs))
				}
				// Check that it's a different slice
				if len(tt.task.PrerequisiteIDs) > 0 && &clone.PrerequisiteIDs[0] == &tt.task.PrerequisiteIDs[0] {
					t.Error("Clone().PrerequisiteIDs should be a deep copy, not a reference to the same slice")
				}
				// Check values are the same
				for i := range tt.task.PrerequisiteIDs {
					if clone.PrerequisiteIDs[i] != tt.task.PrerequisiteIDs[i] {
						t.Errorf("Clone().PrerequisiteIDs[%d] = %v, want %v", i, clone.PrerequisiteIDs[i], tt.task.PrerequisiteIDs[i])
					}
				}
			}

			if tt.task.DownstreamRequiredIDs != nil {
				if len(clone.DownstreamRequiredIDs) != len(tt.task.DownstreamRequiredIDs) {
					t.Errorf("Clone().DownstreamRequiredIDs length = %d, want %d", len(clone.DownstreamRequiredIDs), len(tt.task.DownstreamRequiredIDs))
				}
				// Check that it's a different slice
				if len(tt.task.DownstreamRequiredIDs) > 0 && &clone.DownstreamRequiredIDs[0] == &tt.task.DownstreamRequiredIDs[0] {
					t.Error("Clone().DownstreamRequiredIDs should be a deep copy, not a reference to the same slice")
				}
				// Check values are the same
				for i := range tt.task.DownstreamRequiredIDs {
					if clone.DownstreamRequiredIDs[i] != tt.task.DownstreamRequiredIDs[i] {
						t.Errorf("Clone().DownstreamRequiredIDs[%d] = %v, want %v", i, clone.DownstreamRequiredIDs[i], tt.task.DownstreamRequiredIDs[i])
					}
				}
			}

			if tt.task.DownstreamSuggestedIDs != nil {
				if len(clone.DownstreamSuggestedIDs) != len(tt.task.DownstreamSuggestedIDs) {
					t.Errorf("Clone().DownstreamSuggestedIDs length = %d, want %d", len(clone.DownstreamSuggestedIDs), len(tt.task.DownstreamSuggestedIDs))
				}
				// Check that it's a different slice
				if len(tt.task.DownstreamSuggestedIDs) > 0 && &clone.DownstreamSuggestedIDs[0] == &tt.task.DownstreamSuggestedIDs[0] {
					t.Error("Clone().DownstreamSuggestedIDs should be a deep copy, not a reference to the same slice")
				}
				// Check values are the same
				for i := range tt.task.DownstreamSuggestedIDs {
					if clone.DownstreamSuggestedIDs[i] != tt.task.DownstreamSuggestedIDs[i] {
						t.Errorf("Clone().DownstreamSuggestedIDs[%d] = %v, want %v", i, clone.DownstreamSuggestedIDs[i], tt.task.DownstreamSuggestedIDs[i])
					}
				}
			}

			// Test 7: verify pointer fields are nil (not cloned)
			if clone.Prerequisites != nil {
				t.Error("Clone().Prerequisites should be nil, pointer fields should not be cloned")
			}
			if clone.DownstreamRequired != nil {
				t.Error("Clone().DownstreamRequired should be nil, pointer fields should not be cloned")
			}
			if clone.DownstreamSuggested != nil {
				t.Error("Clone().DownstreamSuggested should be nil, pointer fields should not be cloned")
			}

			// Test 8: modifying the clone should not affect the original
			if clone.Tags != nil && len(clone.Tags) > 0 {
				originalTags := make([]string, len(tt.task.Tags))
				copy(originalTags, tt.task.Tags)

				clone.Tags[0] = "modified-tag"

				if tt.task.Tags[0] == "modified-tag" {
					t.Error("Modifying clone.Tags should not affect the original task")
				}
				if tt.task.Tags[0] != originalTags[0] {
					t.Error("Original task.Tags should not be modified when clone is modified")
				}
			}

			if clone.PrerequisiteIDs != nil && len(clone.PrerequisiteIDs) > 0 {
				originalIDs := make([]string, len(tt.task.PrerequisiteIDs))
				copy(originalIDs, tt.task.PrerequisiteIDs)

				clone.PrerequisiteIDs[0] = "modified-id"

				if tt.task.PrerequisiteIDs[0] == "modified-id" {
					t.Error("Modifying clone.PrerequisiteIDs should not affect the original task")
				}
				if tt.task.PrerequisiteIDs[0] != originalIDs[0] {
					t.Error("Original task.PrerequisiteIDs should not be modified when clone is modified")
				}
			}

			// Modifying scalar fields should not affect original
			clone.Name = "Modified Name"
			if tt.task.Name == "Modified Name" {
				t.Error("Modifying clone.Name should not affect the original task")
			}

			clone.Description = "Modified Description"
			if tt.task.Description == "Modified Description" {
				t.Error("Modifying clone.Description should not affect the original task")
			}
		})
	}
}

func TestTask_SettersAndGetters(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	// Helper tasks for testing
	prereq1 := &Task{ID: "prereq-1", Name: "Prerequisite 1", CreatedAt: now, UpdatedAt: now}
	prereq2 := &Task{ID: "prereq-2", Name: "Prerequisite 2", CreatedAt: now, UpdatedAt: now}
	required1 := &Task{ID: "required-1", Name: "Required 1", CreatedAt: now, UpdatedAt: now}
	required2 := &Task{ID: "required-2", Name: "Required 2", CreatedAt: now, UpdatedAt: now}
	suggested1 := &Task{ID: "suggested-1", Name: "Suggested 1", CreatedAt: now, UpdatedAt: now}
	suggested2 := &Task{ID: "suggested-2", Name: "Suggested 2", CreatedAt: now, UpdatedAt: now}

	tests := []struct {
		name               string
		task               *Task
		operation          string // "set_prerequisites", "set_downstream_required", "set_downstream_suggested"
		inputTasks         []*Task
		wantError          bool
		errorContains      string
		expectedIDs        []string
		expectedTaskPtrs   []*Task
		verifyOtherGetters bool // verify other getters return correct values
	}{
		// Prerequisites tests
		{
			name:             "set prerequisites with valid tasks",
			task:             &Task{ID: "task-1"},
			operation:        "set_prerequisites",
			inputTasks:       []*Task{prereq1, prereq2},
			wantError:        false,
			expectedIDs:      []string{"prereq-1", "prereq-2"},
			expectedTaskPtrs: []*Task{prereq1, prereq2},
		},
		{
			name:             "set prerequisites with single task",
			task:             &Task{ID: "task-1"},
			operation:        "set_prerequisites",
			inputTasks:       []*Task{prereq1},
			wantError:        false,
			expectedIDs:      []string{"prereq-1"},
			expectedTaskPtrs: []*Task{prereq1},
		},
		{
			name:             "set prerequisites with empty slice",
			task:             &Task{ID: "task-1"},
			operation:        "set_prerequisites",
			inputTasks:       []*Task{},
			wantError:        false,
			expectedIDs:      []string{},
			expectedTaskPtrs: []*Task{},
		},
		{
			name:          "set prerequisites on nil task",
			task:          nil,
			operation:     "set_prerequisites",
			inputTasks:    []*Task{prereq1},
			wantError:     true,
			errorContains: "cannot set prerequisites on nil task",
		},
		{
			name:          "set prerequisites with nil task in slice",
			task:          &Task{ID: "task-1"},
			operation:     "set_prerequisites",
			inputTasks:    []*Task{prereq1, nil},
			wantError:     true,
			errorContains: "prerequisite task at index 1 is nil",
		},

		// Downstream Required tests
		{
			name:             "set downstream required with valid tasks",
			task:             &Task{ID: "task-1"},
			operation:        "set_downstream_required",
			inputTasks:       []*Task{required1, required2},
			wantError:        false,
			expectedIDs:      []string{"required-1", "required-2"},
			expectedTaskPtrs: []*Task{required1, required2},
		},
		{
			name:             "set downstream required with single task",
			task:             &Task{ID: "task-1"},
			operation:        "set_downstream_required",
			inputTasks:       []*Task{required1},
			wantError:        false,
			expectedIDs:      []string{"required-1"},
			expectedTaskPtrs: []*Task{required1},
		},
		{
			name:             "set downstream required with empty slice",
			task:             &Task{ID: "task-1"},
			operation:        "set_downstream_required",
			inputTasks:       []*Task{},
			wantError:        false,
			expectedIDs:      []string{},
			expectedTaskPtrs: []*Task{},
		},
		{
			name:          "set downstream required on nil task",
			task:          nil,
			operation:     "set_downstream_required",
			inputTasks:    []*Task{required1},
			wantError:     true,
			errorContains: "cannot set downstream required on nil task",
		},
		{
			name:          "set downstream required with nil task in slice",
			task:          &Task{ID: "task-1"},
			operation:     "set_downstream_required",
			inputTasks:    []*Task{required1, nil, required2},
			wantError:     true,
			errorContains: "downstream required task at index 1 is nil",
		},

		// Downstream Suggested tests
		{
			name:             "set downstream suggested with valid tasks",
			task:             &Task{ID: "task-1"},
			operation:        "set_downstream_suggested",
			inputTasks:       []*Task{suggested1, suggested2},
			wantError:        false,
			expectedIDs:      []string{"suggested-1", "suggested-2"},
			expectedTaskPtrs: []*Task{suggested1, suggested2},
		},
		{
			name:             "set downstream suggested with single task",
			task:             &Task{ID: "task-1"},
			operation:        "set_downstream_suggested",
			inputTasks:       []*Task{suggested1},
			wantError:        false,
			expectedIDs:      []string{"suggested-1"},
			expectedTaskPtrs: []*Task{suggested1},
		},
		{
			name:             "set downstream suggested with empty slice",
			task:             &Task{ID: "task-1"},
			operation:        "set_downstream_suggested",
			inputTasks:       []*Task{},
			wantError:        false,
			expectedIDs:      []string{},
			expectedTaskPtrs: []*Task{},
		},
		{
			name:          "set downstream suggested on nil task",
			task:          nil,
			operation:     "set_downstream_suggested",
			inputTasks:    []*Task{suggested1},
			wantError:     true,
			errorContains: "cannot set downstream suggested on nil task",
		},
		{
			name:          "set downstream suggested with nil task in slice",
			task:          &Task{ID: "task-1"},
			operation:     "set_downstream_suggested",
			inputTasks:    []*Task{nil, suggested1},
			wantError:     true,
			errorContains: "downstream suggested task at index 0 is nil",
		},

		// Test that all getters work correctly after setting
		{
			name:               "verify all getters after setting prerequisites",
			task:               &Task{ID: "task-1"},
			operation:          "set_prerequisites",
			inputTasks:         []*Task{prereq1, prereq2},
			wantError:          false,
			expectedIDs:        []string{"prereq-1", "prereq-2"},
			expectedTaskPtrs:   []*Task{prereq1, prereq2},
			verifyOtherGetters: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			// Execute the setter operation
			switch tt.operation {
			case "set_prerequisites":
				err = tt.task.SetPrerequisites(tt.inputTasks)
			case "set_downstream_required":
				err = tt.task.SetDownstreamRequired(tt.inputTasks)
			case "set_downstream_suggested":
				err = tt.task.SetDownstreamSuggested(tt.inputTasks)
			default:
				t.Fatalf("Unknown operation: %s", tt.operation)
			}

			// Check error expectations
			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Error = %q, want error containing %q", err.Error(), tt.errorContains)
				}
				return // Don't check further assertions if we expected an error
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify the IDs and task pointers were set correctly
			switch tt.operation {
			case "set_prerequisites":
				// Verify IDs
				gotIDs := tt.task.GetPrerequisiteIDs()
				if !equalStringSlices(gotIDs, tt.expectedIDs) {
					t.Errorf("GetPrerequisiteIDs() = %v, want %v", gotIDs, tt.expectedIDs)
				}

				// Verify task pointers
				gotTasks := tt.task.GetPrerequisites()
				if !equalTaskSlices(gotTasks, tt.expectedTaskPtrs) {
					t.Errorf("GetPrerequisites() returned different tasks than expected")
				}

			case "set_downstream_required":
				// Verify IDs
				gotIDs := tt.task.GetDownstreamRequiredIDs()
				if !equalStringSlices(gotIDs, tt.expectedIDs) {
					t.Errorf("GetDownstreamRequiredIDs() = %v, want %v", gotIDs, tt.expectedIDs)
				}

				// Verify task pointers
				gotTasks := tt.task.GetDownstreamRequired()
				if !equalTaskSlices(gotTasks, tt.expectedTaskPtrs) {
					t.Errorf("GetDownstreamRequired() returned different tasks than expected")
				}

			case "set_downstream_suggested":
				// Verify IDs
				gotIDs := tt.task.GetDownstreamSuggestedIDs()
				if !equalStringSlices(gotIDs, tt.expectedIDs) {
					t.Errorf("GetDownstreamSuggestedIDs() = %v, want %v", gotIDs, tt.expectedIDs)
				}

				// Verify task pointers
				gotTasks := tt.task.GetDownstreamSuggested()
				if !equalTaskSlices(gotTasks, tt.expectedTaskPtrs) {
					t.Errorf("GetDownstreamSuggested() returned different tasks than expected")
				}
			}

			// Verify other getters return nil/empty when not set
			if tt.verifyOtherGetters {
				switch tt.operation {
				case "set_prerequisites":
					if tt.task.GetDownstreamRequired() != nil {
						t.Errorf("GetDownstreamRequired() should be nil when not set")
					}
					if tt.task.GetDownstreamSuggested() != nil {
						t.Errorf("GetDownstreamSuggested() should be nil when not set")
					}
				case "set_downstream_required":
					if tt.task.GetPrerequisites() != nil {
						t.Errorf("GetPrerequisites() should be nil when not set")
					}
					if tt.task.GetDownstreamSuggested() != nil {
						t.Errorf("GetDownstreamSuggested() should be nil when not set")
					}
				case "set_downstream_suggested":
					if tt.task.GetPrerequisites() != nil {
						t.Errorf("GetPrerequisites() should be nil when not set")
					}
					if tt.task.GetDownstreamRequired() != nil {
						t.Errorf("GetDownstreamRequired() should be nil when not set")
					}
				}
			}
		})
	}

	// Additional test: verify nil task getters
	t.Run("getters on nil task return nil", func(t *testing.T) {
		var nilTask *Task

		if nilTask.GetPrerequisiteIDs() != nil {
			t.Error("GetPrerequisiteIDs() on nil task should return nil")
		}
		if nilTask.GetDownstreamRequiredIDs() != nil {
			t.Error("GetDownstreamRequiredIDs() on nil task should return nil")
		}
		if nilTask.GetDownstreamSuggestedIDs() != nil {
			t.Error("GetDownstreamSuggestedIDs() on nil task should return nil")
		}
		if nilTask.GetPrerequisites() != nil {
			t.Error("GetPrerequisites() on nil task should return nil")
		}
		if nilTask.GetDownstreamRequired() != nil {
			t.Error("GetDownstreamRequired() on nil task should return nil")
		}
		if nilTask.GetDownstreamSuggested() != nil {
			t.Error("GetDownstreamSuggested() on nil task should return nil")
		}
	})
}

// Helper function to compare string slices
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Helper function to compare task slices (compares pointers)
func equalTaskSlices(a, b []*Task) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestTask_GetAllPrerequisites(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	tests := []struct {
		name          string
		setupGraph    func() *Task
		expectedIDs   []string
		expectedCount int
	}{
		{
			name: "nil task returns empty slice",
			setupGraph: func() *Task {
				return nil
			},
			expectedIDs:   []string{},
			expectedCount: 0,
		},
		{
			name: "task with no prerequisites returns empty slice",
			setupGraph: func() *Task {
				task := &Task{ID: "task-a", Name: "Task A", CreatedAt: now, UpdatedAt: now}
				return task
			},
			expectedIDs:   []string{},
			expectedCount: 0,
		},
		{
			name: "task with single level prerequisites",
			setupGraph: func() *Task {
				prereq1 := &Task{ID: "prereq-1", Name: "Prerequisite 1", CreatedAt: now, UpdatedAt: now}
				prereq2 := &Task{ID: "prereq-2", Name: "Prerequisite 2", CreatedAt: now, UpdatedAt: now}
				task := &Task{ID: "task-a", Name: "Task A", CreatedAt: now, UpdatedAt: now}
				task.Prerequisites = []*Task{prereq1, prereq2}
				return task
			},
			expectedIDs:   []string{"prereq-1", "prereq-2"},
			expectedCount: 2,
		},
		{
			name: "task with two-level chain: A -> B -> C",
			setupGraph: func() *Task {
				taskC := &Task{ID: "task-c", Name: "Task C", CreatedAt: now, UpdatedAt: now}
				taskB := &Task{ID: "task-b", Name: "Task B", Prerequisites: []*Task{taskC}, CreatedAt: now, UpdatedAt: now}
				taskA := &Task{ID: "task-a", Name: "Task A", Prerequisites: []*Task{taskB}, CreatedAt: now, UpdatedAt: now}
				return taskA
			},
			expectedIDs:   []string{"task-b", "task-c"},
			expectedCount: 2,
		},
		{
			name: "task with three-level chain: A -> B -> C -> D",
			setupGraph: func() *Task {
				taskD := &Task{ID: "task-d", Name: "Task D", CreatedAt: now, UpdatedAt: now}
				taskC := &Task{ID: "task-c", Name: "Task C", Prerequisites: []*Task{taskD}, CreatedAt: now, UpdatedAt: now}
				taskB := &Task{ID: "task-b", Name: "Task B", Prerequisites: []*Task{taskC}, CreatedAt: now, UpdatedAt: now}
				taskA := &Task{ID: "task-a", Name: "Task A", Prerequisites: []*Task{taskB}, CreatedAt: now, UpdatedAt: now}
				return taskA
			},
			expectedIDs:   []string{"task-b", "task-c", "task-d"},
			expectedCount: 3,
		},
		{
			name: "task with diamond pattern (shared prerequisite)",
			setupGraph: func() *Task {
				// A depends on B and C, both B and C depend on D
				taskD := &Task{ID: "task-d", Name: "Task D", CreatedAt: now, UpdatedAt: now}
				taskB := &Task{ID: "task-b", Name: "Task B", Prerequisites: []*Task{taskD}, CreatedAt: now, UpdatedAt: now}
				taskC := &Task{ID: "task-c", Name: "Task C", Prerequisites: []*Task{taskD}, CreatedAt: now, UpdatedAt: now}
				taskA := &Task{ID: "task-a", Name: "Task A", Prerequisites: []*Task{taskB, taskC}, CreatedAt: now, UpdatedAt: now}
				return taskA
			},
			expectedIDs:   []string{"task-b", "task-c", "task-d"},
			expectedCount: 3, // task-d should only appear once
		},
		{
			name: "complex graph with multiple paths",
			setupGraph: func() *Task {
				// E -> {D, C}, D -> {B, A}, C -> B, B -> A
				taskA := &Task{ID: "task-a", Name: "Task A", CreatedAt: now, UpdatedAt: now}
				taskB := &Task{ID: "task-b", Name: "Task B", Prerequisites: []*Task{taskA}, CreatedAt: now, UpdatedAt: now}
				taskC := &Task{ID: "task-c", Name: "Task C", Prerequisites: []*Task{taskB}, CreatedAt: now, UpdatedAt: now}
				taskD := &Task{ID: "task-d", Name: "Task D", Prerequisites: []*Task{taskB, taskA}, CreatedAt: now, UpdatedAt: now}
				taskE := &Task{ID: "task-e", Name: "Task E", Prerequisites: []*Task{taskD, taskC}, CreatedAt: now, UpdatedAt: now}
				return taskE
			},
			expectedIDs:   []string{"task-d", "task-c", "task-b", "task-a"},
			expectedCount: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := tt.setupGraph()
			result := task.GetAllPrerequisites()

			if len(result) != tt.expectedCount {
				t.Errorf("GetAllPrerequisites() count = %d, want %d", len(result), tt.expectedCount)
			}

			// Build a set of returned IDs
			returnedIDs := make(map[string]bool)
			for _, task := range result {
				returnedIDs[task.ID] = true
			}

			// Verify all expected IDs are present
			for _, expectedID := range tt.expectedIDs {
				if !returnedIDs[expectedID] {
					t.Errorf("GetAllPrerequisites() missing expected ID: %s", expectedID)
				}
			}

			// Verify no unexpected IDs are present
			if len(returnedIDs) != len(tt.expectedIDs) {
				t.Errorf("GetAllPrerequisites() returned %d unique IDs, want %d", len(returnedIDs), len(tt.expectedIDs))
			}
		})
	}
}

func TestTask_GetAllDownstreamRequired(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	tests := []struct {
		name          string
		setupGraph    func() *Task
		expectedIDs   []string
		expectedCount int
	}{
		{
			name: "nil task returns empty slice",
			setupGraph: func() *Task {
				return nil
			},
			expectedIDs:   []string{},
			expectedCount: 0,
		},
		{
			name: "task with no downstream required returns empty slice",
			setupGraph: func() *Task {
				task := &Task{ID: "task-a", Name: "Task A", CreatedAt: now, UpdatedAt: now}
				return task
			},
			expectedIDs:   []string{},
			expectedCount: 0,
		},
		{
			name: "task with single level downstream required",
			setupGraph: func() *Task {
				downstream1 := &Task{ID: "downstream-1", Name: "Downstream 1", CreatedAt: now, UpdatedAt: now}
				downstream2 := &Task{ID: "downstream-2", Name: "Downstream 2", CreatedAt: now, UpdatedAt: now}
				task := &Task{ID: "task-a", Name: "Task A", DownstreamRequired: []*Task{downstream1, downstream2}, CreatedAt: now, UpdatedAt: now}
				return task
			},
			expectedIDs:   []string{"downstream-1", "downstream-2"},
			expectedCount: 2,
		},
		{
			name: "task with two-level chain: A -> B -> C",
			setupGraph: func() *Task {
				taskC := &Task{ID: "task-c", Name: "Task C", CreatedAt: now, UpdatedAt: now}
				taskB := &Task{ID: "task-b", Name: "Task B", DownstreamRequired: []*Task{taskC}, CreatedAt: now, UpdatedAt: now}
				taskA := &Task{ID: "task-a", Name: "Task A", DownstreamRequired: []*Task{taskB}, CreatedAt: now, UpdatedAt: now}
				return taskA
			},
			expectedIDs:   []string{"task-b", "task-c"},
			expectedCount: 2,
		},
		{
			name: "task with three-level chain: A -> B -> C -> D",
			setupGraph: func() *Task {
				taskD := &Task{ID: "task-d", Name: "Task D", CreatedAt: now, UpdatedAt: now}
				taskC := &Task{ID: "task-c", Name: "Task C", DownstreamRequired: []*Task{taskD}, CreatedAt: now, UpdatedAt: now}
				taskB := &Task{ID: "task-b", Name: "Task B", DownstreamRequired: []*Task{taskC}, CreatedAt: now, UpdatedAt: now}
				taskA := &Task{ID: "task-a", Name: "Task A", DownstreamRequired: []*Task{taskB}, CreatedAt: now, UpdatedAt: now}
				return taskA
			},
			expectedIDs:   []string{"task-b", "task-c", "task-d"},
			expectedCount: 3,
		},
		{
			name: "task with diamond pattern (converging downstream)",
			setupGraph: func() *Task {
				// A -> {B, C}, both B and C -> D
				taskD := &Task{ID: "task-d", Name: "Task D", CreatedAt: now, UpdatedAt: now}
				taskB := &Task{ID: "task-b", Name: "Task B", DownstreamRequired: []*Task{taskD}, CreatedAt: now, UpdatedAt: now}
				taskC := &Task{ID: "task-c", Name: "Task C", DownstreamRequired: []*Task{taskD}, CreatedAt: now, UpdatedAt: now}
				taskA := &Task{ID: "task-a", Name: "Task A", DownstreamRequired: []*Task{taskB, taskC}, CreatedAt: now, UpdatedAt: now}
				return taskA
			},
			expectedIDs:   []string{"task-b", "task-c", "task-d"},
			expectedCount: 3, // task-d should only appear once
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := tt.setupGraph()
			result := task.GetAllDownstreamRequired()

			if len(result) != tt.expectedCount {
				t.Errorf("GetAllDownstreamRequired() count = %d, want %d", len(result), tt.expectedCount)
			}

			// Build a set of returned IDs
			returnedIDs := make(map[string]bool)
			for _, task := range result {
				returnedIDs[task.ID] = true
			}

			// Verify all expected IDs are present
			for _, expectedID := range tt.expectedIDs {
				if !returnedIDs[expectedID] {
					t.Errorf("GetAllDownstreamRequired() missing expected ID: %s", expectedID)
				}
			}

			// Verify no unexpected IDs are present
			if len(returnedIDs) != len(tt.expectedIDs) {
				t.Errorf("GetAllDownstreamRequired() returned %d unique IDs, want %d", len(returnedIDs), len(tt.expectedIDs))
			}
		})
	}
}

func TestTask_GetAllDownstreamSuggested(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	tests := []struct {
		name          string
		setupGraph    func() *Task
		expectedIDs   []string
		expectedCount int
	}{
		{
			name: "nil task returns empty slice",
			setupGraph: func() *Task {
				return nil
			},
			expectedIDs:   []string{},
			expectedCount: 0,
		},
		{
			name: "task with no downstream suggested returns empty slice",
			setupGraph: func() *Task {
				task := &Task{ID: "task-a", Name: "Task A", CreatedAt: now, UpdatedAt: now}
				return task
			},
			expectedIDs:   []string{},
			expectedCount: 0,
		},
		{
			name: "task with single level downstream suggested",
			setupGraph: func() *Task {
				suggested1 := &Task{ID: "suggested-1", Name: "Suggested 1", CreatedAt: now, UpdatedAt: now}
				suggested2 := &Task{ID: "suggested-2", Name: "Suggested 2", CreatedAt: now, UpdatedAt: now}
				task := &Task{ID: "task-a", Name: "Task A", DownstreamSuggested: []*Task{suggested1, suggested2}, CreatedAt: now, UpdatedAt: now}
				return task
			},
			expectedIDs:   []string{"suggested-1", "suggested-2"},
			expectedCount: 2,
		},
		{
			name: "task with two-level chain: A -> B -> C",
			setupGraph: func() *Task {
				taskC := &Task{ID: "task-c", Name: "Task C", CreatedAt: now, UpdatedAt: now}
				taskB := &Task{ID: "task-b", Name: "Task B", DownstreamSuggested: []*Task{taskC}, CreatedAt: now, UpdatedAt: now}
				taskA := &Task{ID: "task-a", Name: "Task A", DownstreamSuggested: []*Task{taskB}, CreatedAt: now, UpdatedAt: now}
				return taskA
			},
			expectedIDs:   []string{"task-b", "task-c"},
			expectedCount: 2,
		},
		{
			name: "task with three-level chain: A -> B -> C -> D",
			setupGraph: func() *Task {
				taskD := &Task{ID: "task-d", Name: "Task D", CreatedAt: now, UpdatedAt: now}
				taskC := &Task{ID: "task-c", Name: "Task C", DownstreamSuggested: []*Task{taskD}, CreatedAt: now, UpdatedAt: now}
				taskB := &Task{ID: "task-b", Name: "Task B", DownstreamSuggested: []*Task{taskC}, CreatedAt: now, UpdatedAt: now}
				taskA := &Task{ID: "task-a", Name: "Task A", DownstreamSuggested: []*Task{taskB}, CreatedAt: now, UpdatedAt: now}
				return taskA
			},
			expectedIDs:   []string{"task-b", "task-c", "task-d"},
			expectedCount: 3,
		},
		{
			name: "task with diamond pattern (converging downstream)",
			setupGraph: func() *Task {
				// A -> {B, C}, both B and C -> D
				taskD := &Task{ID: "task-d", Name: "Task D", CreatedAt: now, UpdatedAt: now}
				taskB := &Task{ID: "task-b", Name: "Task B", DownstreamSuggested: []*Task{taskD}, CreatedAt: now, UpdatedAt: now}
				taskC := &Task{ID: "task-c", Name: "Task C", DownstreamSuggested: []*Task{taskD}, CreatedAt: now, UpdatedAt: now}
				taskA := &Task{ID: "task-a", Name: "Task A", DownstreamSuggested: []*Task{taskB, taskC}, CreatedAt: now, UpdatedAt: now}
				return taskA
			},
			expectedIDs:   []string{"task-b", "task-c", "task-d"},
			expectedCount: 3, // task-d should only appear once
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := tt.setupGraph()
			result := task.GetAllDownstreamSuggested()

			if len(result) != tt.expectedCount {
				t.Errorf("GetAllDownstreamSuggested() count = %d, want %d", len(result), tt.expectedCount)
			}

			// Build a set of returned IDs
			returnedIDs := make(map[string]bool)
			for _, task := range result {
				returnedIDs[task.ID] = true
			}

			// Verify all expected IDs are present
			for _, expectedID := range tt.expectedIDs {
				if !returnedIDs[expectedID] {
					t.Errorf("GetAllDownstreamSuggested() missing expected ID: %s", expectedID)
				}
			}

			// Verify no unexpected IDs are present
			if len(returnedIDs) != len(tt.expectedIDs) {
				t.Errorf("GetAllDownstreamSuggested() returned %d unique IDs, want %d", len(returnedIDs), len(tt.expectedIDs))
			}
		})
	}
}
