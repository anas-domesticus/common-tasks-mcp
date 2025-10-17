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

func TestTaskEquals(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	// Create a base task
	baseTask := &Task{
		ID:            "test-1",
		Name:          "Test Task",
		Summary:       "Summary",
		Description:   "Description",
		Tags:          []string{"tag1", "tag2"},
		DependencyIDs: []string{"dep-1", "dep-2"},
		DependentIDs:  []string{"dependent-1"},
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Test 1: Same task should equal itself
	if !baseTask.Equals(baseTask) {
		t.Error("Task should equal itself")
	}

	// Test 2: Identical task should be equal
	identicalTask := &Task{
		ID:            "test-1",
		Name:          "Test Task",
		Summary:       "Summary",
		Description:   "Description",
		Tags:          []string{"tag1", "tag2"},
		DependencyIDs: []string{"dep-1", "dep-2"},
		DependentIDs:  []string{"dependent-1"},
		CreatedAt:     now,
		UpdatedAt:     now,
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
		ID:            "test-2",
		Name:          "Test Task",
		Summary:       "Summary",
		Description:   "Description",
		Tags:          []string{"tag1", "tag2"},
		DependencyIDs: []string{"dep-1", "dep-2"},
		DependentIDs:  []string{"dependent-1"},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if baseTask.Equals(differentID) {
		t.Error("Tasks with different IDs should not be equal")
	}

	// Test 6: Different Name
	differentName := &Task{
		ID:            "test-1",
		Name:          "Different Name",
		Summary:       "Summary",
		Description:   "Description",
		Tags:          []string{"tag1", "tag2"},
		DependencyIDs: []string{"dep-1", "dep-2"},
		DependentIDs:  []string{"dependent-1"},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if baseTask.Equals(differentName) {
		t.Error("Tasks with different Names should not be equal")
	}

	// Test 7: Different Summary
	differentSummary := &Task{
		ID:            "test-1",
		Name:          "Test Task",
		Summary:       "Different Summary",
		Description:   "Description",
		Tags:          []string{"tag1", "tag2"},
		DependencyIDs: []string{"dep-1", "dep-2"},
		DependentIDs:  []string{"dependent-1"},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if baseTask.Equals(differentSummary) {
		t.Error("Tasks with different Summaries should not be equal")
	}

	// Test 8: Different Description
	differentDescription := &Task{
		ID:            "test-1",
		Name:          "Test Task",
		Summary:       "Summary",
		Description:   "Different Description",
		Tags:          []string{"tag1", "tag2"},
		DependencyIDs: []string{"dep-1", "dep-2"},
		DependentIDs:  []string{"dependent-1"},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if baseTask.Equals(differentDescription) {
		t.Error("Tasks with different Descriptions should not be equal")
	}

	// Test 9: Different CreatedAt
	differentCreatedAt := &Task{
		ID:            "test-1",
		Name:          "Test Task",
		Summary:       "Summary",
		Description:   "Description",
		Tags:          []string{"tag1", "tag2"},
		DependencyIDs: []string{"dep-1", "dep-2"},
		DependentIDs:  []string{"dependent-1"},
		CreatedAt:     now.Add(time.Hour),
		UpdatedAt:     now,
	}
	if baseTask.Equals(differentCreatedAt) {
		t.Error("Tasks with different CreatedAt should not be equal")
	}

	// Test 10: Different UpdatedAt
	differentUpdatedAt := &Task{
		ID:            "test-1",
		Name:          "Test Task",
		Summary:       "Summary",
		Description:   "Description",
		Tags:          []string{"tag1", "tag2"},
		DependencyIDs: []string{"dep-1", "dep-2"},
		DependentIDs:  []string{"dependent-1"},
		CreatedAt:     now,
		UpdatedAt:     now.Add(time.Hour),
	}
	if baseTask.Equals(differentUpdatedAt) {
		t.Error("Tasks with different UpdatedAt should not be equal")
	}

	// Test 11: Different Tags length
	differentTagsLength := &Task{
		ID:            "test-1",
		Name:          "Test Task",
		Summary:       "Summary",
		Description:   "Description",
		Tags:          []string{"tag1"},
		DependencyIDs: []string{"dep-1", "dep-2"},
		DependentIDs:  []string{"dependent-1"},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if baseTask.Equals(differentTagsLength) {
		t.Error("Tasks with different Tags length should not be equal")
	}

	// Test 12: Different Tags values
	differentTagsValues := &Task{
		ID:            "test-1",
		Name:          "Test Task",
		Summary:       "Summary",
		Description:   "Description",
		Tags:          []string{"tag1", "tag3"},
		DependencyIDs: []string{"dep-1", "dep-2"},
		DependentIDs:  []string{"dependent-1"},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if baseTask.Equals(differentTagsValues) {
		t.Error("Tasks with different Tags values should not be equal")
	}

	// Test 13: Different DependencyIDs length
	differentDepsLength := &Task{
		ID:            "test-1",
		Name:          "Test Task",
		Summary:       "Summary",
		Description:   "Description",
		Tags:          []string{"tag1", "tag2"},
		DependencyIDs: []string{"dep-1"},
		DependentIDs:  []string{"dependent-1"},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if baseTask.Equals(differentDepsLength) {
		t.Error("Tasks with different DependencyIDs length should not be equal")
	}

	// Test 14: Different DependencyIDs values
	differentDepsValues := &Task{
		ID:            "test-1",
		Name:          "Test Task",
		Summary:       "Summary",
		Description:   "Description",
		Tags:          []string{"tag1", "tag2"},
		DependencyIDs: []string{"dep-1", "dep-3"},
		DependentIDs:  []string{"dependent-1"},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if baseTask.Equals(differentDepsValues) {
		t.Error("Tasks with different DependencyIDs values should not be equal")
	}

	// Test 15: Different DependentIDs length
	differentDependentsLength := &Task{
		ID:            "test-1",
		Name:          "Test Task",
		Summary:       "Summary",
		Description:   "Description",
		Tags:          []string{"tag1", "tag2"},
		DependencyIDs: []string{"dep-1", "dep-2"},
		DependentIDs:  []string{},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if baseTask.Equals(differentDependentsLength) {
		t.Error("Tasks with different DependentIDs length should not be equal")
	}

	// Test 16: Different DependentIDs values
	differentDependentsValues := &Task{
		ID:            "test-1",
		Name:          "Test Task",
		Summary:       "Summary",
		Description:   "Description",
		Tags:          []string{"tag1", "tag2"},
		DependencyIDs: []string{"dep-1", "dep-2"},
		DependentIDs:  []string{"dependent-2"},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if baseTask.Equals(differentDependentsValues) {
		t.Error("Tasks with different DependentIDs values should not be equal")
	}

	// Test 17: Empty slices vs nil slices (should be equal based on length comparison)
	emptySlices := &Task{
		ID:            "test-empty",
		Name:          "Empty Slices",
		Summary:       "Summary",
		Description:   "Description",
		Tags:          []string{},
		DependencyIDs: []string{},
		DependentIDs:  []string{},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	nilSlices := &Task{
		ID:            "test-empty",
		Name:          "Empty Slices",
		Summary:       "Summary",
		Description:   "Description",
		Tags:          nil,
		DependencyIDs: nil,
		DependentIDs:  nil,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if !emptySlices.Equals(nilSlices) {
		t.Error("Tasks with empty slices and nil slices should be equal")
	}

	// Test 18: Pointer fields should not affect equality
	withPointers := &Task{
		ID:            "test-1",
		Name:          "Test Task",
		Summary:       "Summary",
		Description:   "Description",
		Tags:          []string{"tag1", "tag2"},
		DependencyIDs: []string{"dep-1", "dep-2"},
		DependentIDs:  []string{"dependent-1"},
		Dependencies:  []*Task{{ID: "different"}},
		Dependents:    []*Task{{ID: "different"}},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if !baseTask.Equals(withPointers) {
		t.Error("Pointer fields should not affect equality comparison")
	}
}
