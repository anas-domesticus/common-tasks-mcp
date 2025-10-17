package task_manager

import (
	"fmt"

	"common-tasks-mcp/pkg/types"
)

// PopulateTagCache builds the tag cache by iterating through all tasks
// and indexing them by their tags for efficient tag-based lookups
func (m *Manager) PopulateTagCache() {
	// Clear existing cache
	m.tagCache = make(map[string][]*types.Task)

	// Iterate through all tasks and populate cache
	for _, task := range m.tasks {
		for _, tag := range task.Tags {
			m.tagCache[tag] = append(m.tagCache[tag], task)
		}
	}
}

// GetTasksByTag retrieves all tasks with the specified tag
func (m *Manager) GetTasksByTag(tag string) ([]*types.Task, error) {
	if tag == "" {
		return nil, fmt.Errorf("tag cannot be empty")
	}

	tasks, exists := m.tagCache[tag]
	if !exists {
		return []*types.Task{}, nil
	}

	return tasks, nil
}
