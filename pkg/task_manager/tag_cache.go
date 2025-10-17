package task_manager

import "common-tasks-mcp/pkg/types"

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
