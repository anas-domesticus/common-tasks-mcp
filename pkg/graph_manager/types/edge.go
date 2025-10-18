package types

// Edge represents a directed edge in the task graph with a specific relationship type/category.
// For example, an edge from task A to task B with category "prerequisites" means
// "task A has task B as a prerequisite" or "B must be completed before A".
type Edge struct {
	// To is the destination task ID (the task being referenced)
	To *Node

	// Category is the relationship type (e.g., "prerequisites", "downstream_required", "validates")
	Type *Relationship `json:"type" yaml:"type"`
}

// NewEdge creates a new edge with the given parameters
func NewEdge(to *Node, category *Relationship) Edge {
	return Edge{
		To:   to,
		Type: category,
	}
}
