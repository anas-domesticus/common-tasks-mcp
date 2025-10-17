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

	// Create a base task
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

	// Test 1: Same task should equal itself
	if !baseTask.Equals(baseTask) {
		t.Error("Task should equal itself")
	}

	// Test 2: Identical task should be equal
	identicalTask := &Task{
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
	if !baseTask.Equals(identicalTask) {
		t.Error("Identical tasks should be equal")
	}

	// Test 3: Both nil should be equal
	var nilTask1 *Task
	var nilTask2 *Task
	if !nilTask1.Equals(nilTask2) {
		t.Error("Two nil tasks should be equal")
	}

	// Test 4: One nil, one non-nil should not be equal
	if baseTask.Equals(nilTask1) {
		t.Error("Task and nil should not be equal")
	}
	if nilTask1.Equals(baseTask) {
		t.Error("Nil and task should not be equal")
	}

	// Test 5: Different ID
	differentID := &Task{
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
	}
	if baseTask.Equals(differentID) {
		t.Error("Tasks with different IDs should not be equal")
	}

	// Test 6: Different Name
	differentName := &Task{
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
	}
	if baseTask.Equals(differentName) {
		t.Error("Tasks with different Names should not be equal")
	}

	// Test 7: Different Summary
	differentSummary := &Task{
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
	}
	if baseTask.Equals(differentSummary) {
		t.Error("Tasks with different Summaries should not be equal")
	}

	// Test 8: Different Description
	differentDescription := &Task{
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
	}
	if baseTask.Equals(differentDescription) {
		t.Error("Tasks with different Descriptions should not be equal")
	}

	// Test 9: Different CreatedAt
	differentCreatedAt := &Task{
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
	}
	if baseTask.Equals(differentCreatedAt) {
		t.Error("Tasks with different CreatedAt should not be equal")
	}

	// Test 10: Different UpdatedAt
	differentUpdatedAt := &Task{
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
	}
	if baseTask.Equals(differentUpdatedAt) {
		t.Error("Tasks with different UpdatedAt should not be equal")
	}

	// Test 11: Different Tags length
	differentTagsLength := &Task{
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
	}
	if baseTask.Equals(differentTagsLength) {
		t.Error("Tasks with different Tags length should not be equal")
	}

	// Test 12: Different Tags values
	differentTagsValues := &Task{
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
	}
	if baseTask.Equals(differentTagsValues) {
		t.Error("Tasks with different Tags values should not be equal")
	}

	// Test 13: Different PrerequisiteIDs length
	differentDepsLength := &Task{
		ID:                     "test-1",
		Name:                   "Test Task",
		Summary:                "Summary",
		Description:            "Description",
		Tags:                   []string{"tag1", "tag2"},
		PrerequisiteIDs:        []string{"prereq-1"}, // Only 1 instead of 2
		DownstreamRequiredIDs:  []string{"required-1"},
		DownstreamSuggestedIDs: []string{"suggested-1"},
		CreatedAt:              now,
		UpdatedAt:              now,
	}
	if baseTask.Equals(differentDepsLength) {
		t.Error("Tasks with different PrerequisiteIDs length should not be equal")
	}

	// Test 14: Different PrerequisiteIDs values
	differentDepsValues := &Task{
		ID:                     "test-1",
		Name:                   "Test Task",
		Summary:                "Summary",
		Description:            "Description",
		Tags:                   []string{"tag1", "tag2"},
		PrerequisiteIDs:        []string{"prereq-1", "prereq-3"}, // Different value
		DownstreamRequiredIDs:  []string{"required-1"},
		DownstreamSuggestedIDs: []string{"suggested-1"},
		CreatedAt:              now,
		UpdatedAt:              now,
	}
	if baseTask.Equals(differentDepsValues) {
		t.Error("Tasks with different PrerequisiteIDs values should not be equal")
	}

	// Test 15: Different DownstreamRequiredIDs length
	differentDependentsLength := &Task{
		ID:                     "test-1",
		Name:                   "Test Task",
		Summary:                "Summary",
		Description:            "Description",
		Tags:                   []string{"tag1", "tag2"},
		PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
		DownstreamRequiredIDs:  []string{}, // Empty instead of 1
		DownstreamSuggestedIDs: []string{"suggested-1"},
		CreatedAt:              now,
		UpdatedAt:              now,
	}
	if baseTask.Equals(differentDependentsLength) {
		t.Error("Tasks with different DownstreamRequiredIDs length should not be equal")
	}

	// Test 16: Different DownstreamRequiredIDs values
	differentDependentsValues := &Task{
		ID:                     "test-1",
		Name:                   "Test Task",
		Summary:                "Summary",
		Description:            "Description",
		Tags:                   []string{"tag1", "tag2"},
		PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
		DownstreamRequiredIDs:  []string{"required-2"}, // Different value
		DownstreamSuggestedIDs: []string{"suggested-1"},
		CreatedAt:              now,
		UpdatedAt:              now,
	}
	if baseTask.Equals(differentDependentsValues) {
		t.Error("Tasks with different DownstreamRequiredIDs values should not be equal")
	}

	// Test 17: Different DownstreamSuggestedIDs length
	differentSuggestedLength := &Task{
		ID:                     "test-1",
		Name:                   "Test Task",
		Summary:                "Summary",
		Description:            "Description",
		Tags:                   []string{"tag1", "tag2"},
		PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
		DownstreamRequiredIDs:  []string{"required-1"},
		DownstreamSuggestedIDs: []string{}, // Empty instead of 1
		CreatedAt:              now,
		UpdatedAt:              now,
	}
	if baseTask.Equals(differentSuggestedLength) {
		t.Error("Tasks with different DownstreamSuggestedIDs length should not be equal")
	}

	// Test 18: Different DownstreamSuggestedIDs values
	differentSuggestedValues := &Task{
		ID:                     "test-1",
		Name:                   "Test Task",
		Summary:                "Summary",
		Description:            "Description",
		Tags:                   []string{"tag1", "tag2"},
		PrerequisiteIDs:        []string{"prereq-1", "prereq-2"},
		DownstreamRequiredIDs:  []string{"required-1"},
		DownstreamSuggestedIDs: []string{"suggested-2"}, // Different value
		CreatedAt:              now,
		UpdatedAt:              now,
	}
	if baseTask.Equals(differentSuggestedValues) {
		t.Error("Tasks with different DownstreamSuggestedIDs values should not be equal")
	}

	// Test 19: Empty slices vs nil slices (should be equal based on length comparison)
	emptySlices := &Task{
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
	}
	nilSlices := &Task{
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
	}
	if !emptySlices.Equals(nilSlices) {
		t.Error("Tasks with empty slices and nil slices should be equal")
	}

	// Test 20: Pointer fields should not affect equality
	withPointers := &Task{
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
	}
	if !baseTask.Equals(withPointers) {
		t.Error("Pointer fields should not affect equality comparison")
	}
}
