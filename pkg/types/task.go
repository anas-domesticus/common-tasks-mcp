package types

import "time"

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
