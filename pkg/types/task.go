package types

import "time"

// Task represents a node in the task graph
type Task struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Summary      string    `json:"summary"`
	Description  string    `json:"description"`
	Tags         []string  `json:"tags"`
	Dependencies []string  `json:"dependencies"` // IDs of upstream tasks (tasks this depends on)
	Dependents   []string  `json:"dependents"`   // IDs of downstream tasks (tasks that depend on this)
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
