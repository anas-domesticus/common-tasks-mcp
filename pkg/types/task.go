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
