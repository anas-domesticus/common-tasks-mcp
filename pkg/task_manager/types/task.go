package types

import (
	"errors"
	"fmt"
	"time"
)

// Task represents a node in the task graph with three distinct DAG relationships
type Task struct {
	ID          string   `json:"id" yaml:"id"`
	Name        string   `json:"name" yaml:"name"`
	Summary     string   `json:"summary" yaml:"summary"`
	Description string   `json:"description" yaml:"description"`
	Tags        []string `json:"tags" yaml:"tags"`
	// DAG #1: Prerequisites - tasks that must be completed before this task
	PrerequisiteIDs []string `json:"prerequisites" yaml:"prerequisites"` // IDs of prerequisite tasks
	Prerequisites   []*Task  `json:"-" yaml:"-"`                         // Resolved prerequisite task pointers (not persisted)
	// DAG #2: Downstream Required - tasks that must be completed after this task
	DownstreamRequiredIDs []string `json:"downstream_required" yaml:"downstream_required"` // IDs of required downstream tasks
	DownstreamRequired    []*Task  `json:"-" yaml:"-"`                                     // Resolved required downstream task pointers (not persisted)
	// DAG #3: Downstream Suggested - tasks that are recommended but optional after this task
	DownstreamSuggestedIDs []string  `json:"downstream_suggested" yaml:"downstream_suggested"` // IDs of suggested downstream tasks
	DownstreamSuggested    []*Task   `json:"-" yaml:"-"`                                       // Resolved suggested downstream task pointers (not persisted)
	CreatedAt              time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" yaml:"updated_at"`
}

// Equals compares all persisted fields of two tasks for equality.
// Returns true if all fields match, false otherwise.
// Note: This only compares persisted fields (excludes Prerequisites, DownstreamRequired, and DownstreamSuggested pointer fields).
func (t *Task) Equals(other *Task) bool {
	if t == nil && other == nil {
		return true
	}
	if t == nil || other == nil {
		return false
	}

	// Compare scalar fields
	if t.ID != other.ID ||
		t.Name != other.Name ||
		t.Summary != other.Summary ||
		t.Description != other.Description {
		return false
	}

	// Compare timestamps
	if !t.CreatedAt.Equal(other.CreatedAt) || !t.UpdatedAt.Equal(other.UpdatedAt) {
		return false
	}

	// Compare Tags slice
	if len(t.Tags) != len(other.Tags) {
		return false
	}
	for i := range t.Tags {
		if t.Tags[i] != other.Tags[i] {
			return false
		}
	}

	// Compare PrerequisiteIDs slice
	if len(t.PrerequisiteIDs) != len(other.PrerequisiteIDs) {
		return false
	}
	for i := range t.PrerequisiteIDs {
		if t.PrerequisiteIDs[i] != other.PrerequisiteIDs[i] {
			return false
		}
	}

	// Compare DownstreamRequiredIDs slice
	if len(t.DownstreamRequiredIDs) != len(other.DownstreamRequiredIDs) {
		return false
	}
	for i := range t.DownstreamRequiredIDs {
		if t.DownstreamRequiredIDs[i] != other.DownstreamRequiredIDs[i] {
			return false
		}
	}

	// Compare DownstreamSuggestedIDs slice
	if len(t.DownstreamSuggestedIDs) != len(other.DownstreamSuggestedIDs) {
		return false
	}
	for i := range t.DownstreamSuggestedIDs {
		if t.DownstreamSuggestedIDs[i] != other.DownstreamSuggestedIDs[i] {
			return false
		}
	}

	return true
}

// Clone creates a deep copy of the task, copying all persisted fields.
// Note: The pointer fields (Prerequisites, DownstreamRequired, DownstreamSuggested)
// are NOT cloned - they should be resolved by the Manager after cloning.
func (t *Task) Clone() *Task {
	if t == nil {
		return nil
	}

	// Clone the task with all scalar fields
	clone := &Task{
		ID:          t.ID,
		Name:        t.Name,
		Summary:     t.Summary,
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}

	// Deep copy Tags slice
	if t.Tags != nil {
		clone.Tags = make([]string, len(t.Tags))
		copy(clone.Tags, t.Tags)
	}

	// Deep copy PrerequisiteIDs slice
	if t.PrerequisiteIDs != nil {
		clone.PrerequisiteIDs = make([]string, len(t.PrerequisiteIDs))
		copy(clone.PrerequisiteIDs, t.PrerequisiteIDs)
	}

	// Deep copy DownstreamRequiredIDs slice
	if t.DownstreamRequiredIDs != nil {
		clone.DownstreamRequiredIDs = make([]string, len(t.DownstreamRequiredIDs))
		copy(clone.DownstreamRequiredIDs, t.DownstreamRequiredIDs)
	}

	// Deep copy DownstreamSuggestedIDs slice
	if t.DownstreamSuggestedIDs != nil {
		clone.DownstreamSuggestedIDs = make([]string, len(t.DownstreamSuggestedIDs))
		copy(clone.DownstreamSuggestedIDs, t.DownstreamSuggestedIDs)
	}

	// Note: pointer fields (Prerequisites, DownstreamRequired, DownstreamSuggested)
	// are intentionally left as nil and will be resolved by Manager.ResolveTaskPointers()

	return clone
}

// GetPrerequisiteIDs returns the IDs of prerequisite tasks.
// Returns nil if the receiver is nil.
func (t *Task) GetPrerequisiteIDs() []string {
	if t == nil {
		return nil
	}
	return t.PrerequisiteIDs
}

// GetDownstreamRequiredIDs returns the IDs of required downstream tasks.
// Returns nil if the receiver is nil.
func (t *Task) GetDownstreamRequiredIDs() []string {
	if t == nil {
		return nil
	}
	return t.DownstreamRequiredIDs
}

// GetDownstreamSuggestedIDs returns the IDs of suggested downstream tasks.
// Returns nil if the receiver is nil.
func (t *Task) GetDownstreamSuggestedIDs() []string {
	if t == nil {
		return nil
	}
	return t.DownstreamSuggestedIDs
}

// GetPrerequisites returns the resolved prerequisite task pointers.
// Returns nil if the receiver is nil.
func (t *Task) GetPrerequisites() []*Task {
	if t == nil {
		return nil
	}
	return t.Prerequisites
}

// GetDownstreamRequired returns the resolved required downstream task pointers.
// Returns nil if the receiver is nil.
func (t *Task) GetDownstreamRequired() []*Task {
	if t == nil {
		return nil
	}
	return t.DownstreamRequired
}

// GetDownstreamSuggested returns the resolved suggested downstream task pointers.
// Returns nil if the receiver is nil.
func (t *Task) GetDownstreamSuggested() []*Task {
	if t == nil {
		return nil
	}
	return t.DownstreamSuggested
}

// SetPrerequisites sets both the prerequisite task pointers and their IDs.
// Extracts IDs from the provided tasks and updates both Prerequisites and PrerequisiteIDs.
// Returns an error if the receiver is nil or if any task in the slice is nil.
func (t *Task) SetPrerequisites(tasks []*Task) error {
	if t == nil {
		return fmt.Errorf("cannot set prerequisites on nil task")
	}

	// Extract IDs from tasks
	ids := make([]string, len(tasks))
	for i, task := range tasks {
		if task == nil {
			return fmt.Errorf("prerequisite task at index %d is nil", i)
		}
		ids[i] = task.ID
	}

	t.Prerequisites = tasks
	t.PrerequisiteIDs = ids
	return nil
}

// SetDownstreamRequired sets both the required downstream task pointers and their IDs.
// Extracts IDs from the provided tasks and updates both DownstreamRequired and DownstreamRequiredIDs.
// Returns an error if the receiver is nil or if any task in the slice is nil.
func (t *Task) SetDownstreamRequired(tasks []*Task) error {
	if t == nil {
		return fmt.Errorf("cannot set downstream required on nil task")
	}

	// Extract IDs from tasks
	ids := make([]string, len(tasks))
	for i, task := range tasks {
		if task == nil {
			return fmt.Errorf("downstream required task at index %d is nil", i)
		}
		ids[i] = task.ID
	}

	t.DownstreamRequired = tasks
	t.DownstreamRequiredIDs = ids
	return nil
}

// SetDownstreamSuggested sets both the suggested downstream task pointers and their IDs.
// Extracts IDs from the provided tasks and updates both DownstreamSuggested and DownstreamSuggestedIDs.
// Returns an error if the receiver is nil or if any task in the slice is nil.
func (t *Task) SetDownstreamSuggested(tasks []*Task) error {
	if t == nil {
		return fmt.Errorf("cannot set downstream suggested on nil task")
	}

	// Extract IDs from tasks
	ids := make([]string, len(tasks))
	for i, task := range tasks {
		if task == nil {
			return fmt.Errorf("downstream suggested task at index %d is nil", i)
		}
		ids[i] = task.ID
	}

	t.DownstreamSuggested = tasks
	t.DownstreamSuggestedIDs = ids
	return nil
}

// checkEdgeConsistency validates that a task slice and string slice conform.
// The IDs in the string slice should match the IDs in the task slice (same set, any order).
// Returns nil if consistent, or a descriptive error if inconsistent.
func checkEdgeConsistency(tasks []*Task, ids []string) error {
	// Handle nil/empty cases
	if len(tasks) == 0 && len(ids) == 0 {
		return nil
	}

	// Check length mismatch
	if len(tasks) != len(ids) {
		return fmt.Errorf("length mismatch: %d tasks but %d IDs", len(tasks), len(ids))
	}

	// Build a map of IDs from the string slice
	idMap := make(map[string]bool, len(ids))
	for _, id := range ids {
		if idMap[id] {
			return fmt.Errorf("duplicate ID in string slice: %s", id)
		}
		idMap[id] = true
	}

	// Verify each task's ID is in the string slice
	taskIDMap := make(map[string]bool, len(tasks))
	for i, task := range tasks {
		if task == nil {
			return fmt.Errorf("nil task at index %d", i)
		}
		if !idMap[task.ID] {
			return fmt.Errorf("task ID %s not found in string slice", task.ID)
		}
		if taskIDMap[task.ID] {
			return fmt.Errorf("duplicate task ID in task slice: %s", task.ID)
		}
		taskIDMap[task.ID] = true
	}

	// Verify each string ID has a corresponding task (should be covered by length check + above, but explicit)
	for _, id := range ids {
		if !taskIDMap[id] {
			return fmt.Errorf("ID %s not found in task slice", id)
		}
	}

	return nil
}

// checkPrerequisiteConsistency validates that Prerequisites and PrerequisiteIDs are consistent.
// Returns nil if consistent, or a descriptive error if inconsistent.
func (t *Task) checkPrerequisiteConsistency() error {
	if t == nil {
		return nil
	}
	return checkEdgeConsistency(t.Prerequisites, t.PrerequisiteIDs)
}

// checkDownstreamRequiredConsistency validates that DownstreamRequired and DownstreamRequiredIDs are consistent.
// Returns nil if consistent, or a descriptive error if inconsistent.
func (t *Task) checkDownstreamRequiredConsistency() error {
	if t == nil {
		return nil
	}
	return checkEdgeConsistency(t.DownstreamRequired, t.DownstreamRequiredIDs)
}

// checkDownstreamSuggestedConsistency validates that DownstreamSuggested and DownstreamSuggestedIDs are consistent.
// Returns nil if consistent, or a descriptive error if inconsistent.
func (t *Task) checkDownstreamSuggestedConsistency() error {
	if t == nil {
		return nil
	}
	return checkEdgeConsistency(t.DownstreamSuggested, t.DownstreamSuggestedIDs)
}

// CheckEdgeConsistency validates that all DAG edge relationships are consistent.
// It checks that Prerequisites, DownstreamRequired, and DownstreamSuggested slices
// match their corresponding ID slices (PrerequisiteIDs, DownstreamRequiredIDs, DownstreamSuggestedIDs).
// Returns nil if all edges are consistent, or a joined error containing all inconsistencies found.
func (t *Task) CheckEdgeConsistency() error {
	if t == nil {
		return nil
	}

	var errs []error

	if err := t.checkPrerequisiteConsistency(); err != nil {
		errs = append(errs, fmt.Errorf("prerequisite edges: %w", err))
	}

	if err := t.checkDownstreamRequiredConsistency(); err != nil {
		errs = append(errs, fmt.Errorf("downstream required edges: %w", err))
	}

	if err := t.checkDownstreamSuggestedConsistency(); err != nil {
		errs = append(errs, fmt.Errorf("downstream suggested edges: %w", err))
	}

	return errors.Join(errs...)
}

// collectTasksRecursive recursively collects all tasks reachable through the given getter function.
// It uses a visited map to avoid infinite loops in case of cycles (though DAGs should not have cycles).
// The getter function should return the task slice to traverse (e.g., Prerequisites, DownstreamRequired, etc.).
func collectTasksRecursive(tasks []*Task, getter func(*Task) []*Task, visited map[string]bool, result *[]*Task) {
	for _, task := range tasks {
		if task == nil {
			continue
		}

		// Skip if already visited to avoid infinite loops
		if visited[task.ID] {
			continue
		}

		// Mark as visited and add to result
		visited[task.ID] = true
		*result = append(*result, task)

		// Recursively traverse the next level
		nextTasks := getter(task)
		if len(nextTasks) > 0 {
			collectTasksRecursive(nextTasks, getter, visited, result)
		}
	}
}

// GetAllPrerequisites recursively collects all prerequisite tasks in the entire chain.
// For example, if task A has prerequisite B, and B has prerequisite C, this returns [B, C].
// The order is breadth-first traversal through the prerequisite chain.
// Returns an empty slice if there are no prerequisites or if the receiver is nil.
func (t *Task) GetAllPrerequisites() []*Task {
	if t == nil || len(t.Prerequisites) == 0 {
		return []*Task{}
	}

	visited := make(map[string]bool)
	result := make([]*Task, 0)

	collectTasksRecursive(t.Prerequisites, func(task *Task) []*Task {
		return task.Prerequisites
	}, visited, &result)

	return result
}

// GetAllDownstreamRequired recursively collects all required downstream tasks in the entire chain.
// For example, if task A has required downstream B, and B has required downstream C, this returns [B, C].
// The order is breadth-first traversal through the downstream required chain.
// Returns an empty slice if there are no required downstream tasks or if the receiver is nil.
func (t *Task) GetAllDownstreamRequired() []*Task {
	if t == nil || len(t.DownstreamRequired) == 0 {
		return []*Task{}
	}

	visited := make(map[string]bool)
	result := make([]*Task, 0)

	collectTasksRecursive(t.DownstreamRequired, func(task *Task) []*Task {
		return task.DownstreamRequired
	}, visited, &result)

	return result
}

// GetAllDownstreamSuggested recursively collects all suggested downstream tasks in the entire chain.
// For example, if task A has suggested downstream B, and B has suggested downstream C, this returns [B, C].
// The order is breadth-first traversal through the downstream suggested chain.
// Returns an empty slice if there are no suggested downstream tasks or if the receiver is nil.
func (t *Task) GetAllDownstreamSuggested() []*Task {
	if t == nil || len(t.DownstreamSuggested) == 0 {
		return []*Task{}
	}

	visited := make(map[string]bool)
	result := make([]*Task, 0)

	collectTasksRecursive(t.DownstreamSuggested, func(task *Task) []*Task {
		return task.DownstreamSuggested
	}, visited, &result)

	return result
}
