package graph_manager

import (
	"fmt"

	"common-tasks-mcp/pkg/graph_manager/types"
)

// PopulateTagCache builds the tag cache by iterating through all nodes
// and indexing them by their tags for efficient tag-based lookups
func (m *Manager) PopulateTagCache() {
	// Clear existing cache
	m.tagCache = make(map[string][]*types.Node)

	// Iterate through all nodes and populate cache
	for _, node := range m.nodes {
		for _, tag := range node.Tags {
			m.tagCache[tag] = append(m.tagCache[tag], node)
		}
	}
}

// GetNodesByTag retrieves all nodes with the specified tag
func (m *Manager) GetNodesByTag(tag string) ([]*types.Node, error) {
	if tag == "" {
		return nil, fmt.Errorf("tag cannot be empty")
	}

	nodes, exists := m.tagCache[tag]
	if !exists {
		return []*types.Node{}, nil
	}

	return nodes, nil
}

// GetAllTags retrieves all unique tags from the tag cache
// Returns a map where keys are tag names and values are the count of nodes with that tag
func (m *Manager) GetAllTags() map[string]int {
	tags := make(map[string]int)

	for tag, nodes := range m.tagCache {
		tags[tag] = len(nodes)
	}

	return tags
}
