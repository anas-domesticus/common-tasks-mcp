package types

import "time"

// Task represents a node in the task graph
type Task struct {
	ID            string    `json:"id" yaml:"id"`
	Name          string    `json:"name" yaml:"name"`
	Summary       string    `json:"summary" yaml:"summary"`
	Description   string    `json:"description" yaml:"description"`
	Tags          []string  `json:"tags" yaml:"tags"`
	DependencyIDs []string  `json:"dependencies" yaml:"dependencies"` // IDs of upstream tasks (tasks this depends on)
	DependentIDs  []string  `json:"dependents" yaml:"dependents"`     // IDs of downstream tasks (tasks that depend on this)
	Dependencies  []*Task   `json:"-" yaml:"-"`                       // Resolved upstream task pointers (not persisted)
	Dependents    []*Task   `json:"-" yaml:"-"`                       // Resolved downstream task pointers (not persisted)
	CreatedAt     time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" yaml:"updated_at"`
}

// Equals compares all persisted fields of two tasks for equality.
// Returns true if all fields match, false otherwise.
// Note: This only compares persisted fields (excludes Dependencies and Dependents pointer fields).
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

	// Compare DependencyIDs slice
	if len(t.DependencyIDs) != len(other.DependencyIDs) {
		return false
	}
	for i := range t.DependencyIDs {
		if t.DependencyIDs[i] != other.DependencyIDs[i] {
			return false
		}
	}

	// Compare DependentIDs slice
	if len(t.DependentIDs) != len(other.DependentIDs) {
		return false
	}
	for i := range t.DependentIDs {
		if t.DependentIDs[i] != other.DependentIDs[i] {
			return false
		}
	}

	return true
}
