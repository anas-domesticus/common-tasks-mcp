package types

import "fmt"

// Relationship represents a type/category of relationship between nodes in the graph.
// All relationships form directed acyclic graphs (DAGs).
// For example, "prerequisites" (must be preceded by) or "downstream_required" (must be run after).
type Relationship struct {
	// Name is the identifier for this relationship type (e.g., "prerequisites", "validates")
	Name string `json:"name" yaml:"name"`

	// Description explains what this relationship means in human-readable form
	Description string `json:"description" yaml:"description"`

	// Direction indicates how this relationship flows relative to execution order.
	// "backward" means it points to things that come before (e.g., prerequisites)
	// "forward" means it points to things that come after (e.g., downstream tasks)
	// "none" means the relationship has no temporal ordering
	Direction RelationshipDirection `json:"direction" yaml:"direction"`
}

// NewRelationship creates a new relationship type with the given parameters
func NewRelationship(name, description string, direction RelationshipDirection) (Relationship, error) {
	return Relationship{
		Name:        name,
		Description: description,
		Direction:   direction,
	}, nil
}

// Validate checks if the relationship configuration is valid
func (r Relationship) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("relationship name is required")
	}

	// Validate direction
	switch r.Direction {
	case DirectionBackward, DirectionForward, DirectionNone:
		// Valid
	default:
		return fmt.Errorf("invalid direction: %s", r.Direction)
	}

	return nil
}
