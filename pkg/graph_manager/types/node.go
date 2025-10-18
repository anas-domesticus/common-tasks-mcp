package types

import (
	"fmt"
	"time"
)

// Node represents a node in the task graph with configurable DAG relationships
type Node struct {
	ID          string   `json:"id" yaml:"id"`
	Name        string   `json:"name" yaml:"name"`
	Summary     string   `json:"summary" yaml:"summary"`
	Description string   `json:"description" yaml:"description"`
	Tags        []string `json:"tags" yaml:"tags"`

	// EdgeIDs maps relationship names to lists of target node IDs (persisted to YAML)
	// Example: {"prerequisites": ["task-a", "task-b"], "downstream_required": ["task-c"]}
	EdgeIDs map[string][]string `json:"edges" yaml:"edges"`

	// Edges maps relationship names to lists of resolved edges (runtime only, not persisted)
	// The edges are populated by the Manager after loading from disk
	Edges map[string][]Edge `json:"-" yaml:"-"`

	CreatedAt time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `json:"updated_at" yaml:"updated_at"`
}

// Equals compares all persisted fields of two nodes for equality.
// Returns true if all fields match, false otherwise.
// Note: This only compares persisted fields (excludes the resolved Edges map).
func (n *Node) Equals(other *Node) bool {
	if n == nil && other == nil {
		return true
	}
	if n == nil || other == nil {
		return false
	}

	// Compare scalar fields
	if n.ID != other.ID ||
		n.Name != other.Name ||
		n.Summary != other.Summary ||
		n.Description != other.Description {
		return false
	}

	// Compare timestamps
	if !n.CreatedAt.Equal(other.CreatedAt) || !n.UpdatedAt.Equal(other.UpdatedAt) {
		return false
	}

	// Compare Tags slice
	if len(n.Tags) != len(other.Tags) {
		return false
	}
	for i := range n.Tags {
		if n.Tags[i] != other.Tags[i] {
			return false
		}
	}

	// Compare EdgeIDs map
	if len(n.EdgeIDs) != len(other.EdgeIDs) {
		return false
	}
	for relationshipName, ids := range n.EdgeIDs {
		otherIDs, exists := other.EdgeIDs[relationshipName]
		if !exists {
			return false
		}
		if len(ids) != len(otherIDs) {
			return false
		}
		for i := range ids {
			if ids[i] != otherIDs[i] {
				return false
			}
		}
	}

	return true
}

// Clone creates a deep copy of the node, copying all persisted fields.
// Note: The Edges map (resolved pointers) is NOT cloned - it should be resolved by the Manager after cloning.
func (n *Node) Clone() *Node {
	if n == nil {
		return nil
	}

	// Clone the node with all scalar fields
	clone := &Node{
		ID:          n.ID,
		Name:        n.Name,
		Summary:     n.Summary,
		Description: n.Description,
		CreatedAt:   n.CreatedAt,
		UpdatedAt:   n.UpdatedAt,
	}

	// Deep copy Tags slice
	if n.Tags != nil {
		clone.Tags = make([]string, len(n.Tags))
		copy(clone.Tags, n.Tags)
	}

	// Deep copy EdgeIDs map
	if n.EdgeIDs != nil {
		clone.EdgeIDs = make(map[string][]string, len(n.EdgeIDs))
		for relationshipName, ids := range n.EdgeIDs {
			idsCopy := make([]string, len(ids))
			copy(idsCopy, ids)
			clone.EdgeIDs[relationshipName] = idsCopy
		}
	}

	// Note: Edges map (resolved pointers) is intentionally left as nil
	// and will be resolved by Manager.ResolveEdges()

	return clone
}

// GetEdgeIDs returns the IDs of nodes connected by the specified relationship.
// Returns nil if the relationship doesn't exist or the receiver is nil.
func (n *Node) GetEdgeIDs(relationshipName string) []string {
	if n == nil || n.EdgeIDs == nil {
		return nil
	}
	return n.EdgeIDs[relationshipName]
}

// GetEdges returns the resolved edges for the specified relationship.
// Returns nil if the relationship doesn't exist or the receiver is nil.
func (n *Node) GetEdges(relationshipName string) []Edge {
	if n == nil || n.Edges == nil {
		return nil
	}
	return n.Edges[relationshipName]
}

// SetEdgeIDs sets the IDs for a specific relationship, replacing any existing IDs.
// Returns an error if the receiver is nil.
func (n *Node) SetEdgeIDs(relationshipName string, ids []string) error {
	if n == nil {
		return fmt.Errorf("cannot set edge IDs on nil node")
	}
	if n.EdgeIDs == nil {
		n.EdgeIDs = make(map[string][]string)
	}
	n.EdgeIDs[relationshipName] = ids
	return nil
}

// AddEdgeID adds a single ID to a specific relationship.
// Returns an error if the receiver is nil.
func (n *Node) AddEdgeID(relationshipName string, id string) error {
	if n == nil {
		return fmt.Errorf("cannot add edge ID to nil node")
	}
	if n.EdgeIDs == nil {
		n.EdgeIDs = make(map[string][]string)
	}
	n.EdgeIDs[relationshipName] = append(n.EdgeIDs[relationshipName], id)
	return nil
}

// SetEdges sets the resolved edges for a specific relationship, replacing any existing edges.
// Also updates the corresponding EdgeIDs map.
// Returns an error if the receiver is nil or if any edge has a nil To node.
func (n *Node) SetEdges(relationshipName string, edges []Edge) error {
	if n == nil {
		return fmt.Errorf("cannot set edges on nil node")
	}

	// Extract IDs from edges
	ids := make([]string, len(edges))
	for i, edge := range edges {
		if edge.To == nil {
			return fmt.Errorf("edge at index %d has nil To node", i)
		}
		ids[i] = edge.To.ID
	}

	// Update both maps
	if n.Edges == nil {
		n.Edges = make(map[string][]Edge)
	}
	if n.EdgeIDs == nil {
		n.EdgeIDs = make(map[string][]string)
	}

	n.Edges[relationshipName] = edges
	n.EdgeIDs[relationshipName] = ids
	return nil
}

// AddEdge adds a single edge to a specific relationship.
// Also updates the corresponding EdgeIDs map.
// Returns an error if the receiver is nil or if the edge has a nil To node.
func (n *Node) AddEdge(relationshipName string, edge Edge) error {
	if n == nil {
		return fmt.Errorf("cannot add edge to nil node")
	}
	if edge.To == nil {
		return fmt.Errorf("edge has nil To node")
	}

	if n.Edges == nil {
		n.Edges = make(map[string][]Edge)
	}
	if n.EdgeIDs == nil {
		n.EdgeIDs = make(map[string][]string)
	}

	n.Edges[relationshipName] = append(n.Edges[relationshipName], edge)
	n.EdgeIDs[relationshipName] = append(n.EdgeIDs[relationshipName], edge.To.ID)
	return nil
}
