package task_manager

import (
	"fmt"
	"os"
	"path/filepath"

	"common-tasks-mcp/pkg/task_manager/types"

	"gopkg.in/yaml.v3"
)

// Manager handles task graph operations
type Manager struct {
	tasks    map[string]*types.Task
	tagCache map[string][]*types.Task
}

// NewManager creates a new task manager instance
func NewManager() *Manager {
	return &Manager{
		tasks:    make(map[string]*types.Task),
		tagCache: make(map[string][]*types.Task),
	}
}

// AddTask adds a task to the manager.
// It uses a clone-validate-commit pattern to ensure the addition doesn't introduce cycles.
func (m *Manager) AddTask(task *types.Task) error {
	if task == nil {
		return fmt.Errorf("task cannot be nil")
	}
	if task.ID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	if _, exists := m.tasks[task.ID]; exists {
		return fmt.Errorf("task with ID %s already exists", task.ID)
	}

	// Clone the manager to test the addition
	testManager := m.Clone()

	// Perform the addition in the test manager
	testManager.tasks[task.ID] = task

	// Check for cycles in the test manager
	if err := testManager.DetectCycles(); err != nil {
		return fmt.Errorf("addition would introduce cycle: %w", err)
	}

	// If no cycles detected, commit the addition to the original manager
	m.tasks[task.ID] = task

	// Update tag cache with the new task
	m.PopulateTagCache()

	return nil
}

// UpdateTask updates an existing task in the manager.
// It uses a clone-validate-commit pattern to ensure the update doesn't introduce cycles,
// and automatically refreshes all task pointers to prevent stale references.
func (m *Manager) UpdateTask(task *types.Task) error {
	if task == nil {
		return fmt.Errorf("task cannot be nil")
	}
	if task.ID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	if _, exists := m.tasks[task.ID]; !exists {
		return fmt.Errorf("task with ID %s not found", task.ID)
	}

	// Clone the manager to test the update
	testManager := m.Clone()

	// Perform the update in the test manager
	testManager.tasks[task.ID] = task

	// Check for cycles in the test manager
	if err := testManager.DetectCycles(); err != nil {
		return fmt.Errorf("update would introduce cycle: %w", err)
	}

	// If no cycles detected, commit the update to the original manager
	m.tasks[task.ID] = task

	// Resolve all task pointers to fix stale references
	// This ensures that any tasks pointing to the updated task get fresh pointers
	if err := m.ResolveTaskPointers(); err != nil {
		return err
	}

	// Update tag cache since tags may have changed
	m.PopulateTagCache()

	return nil
}

// DeleteTask removes a task from the manager and cleans up all references to it
// from other tasks' prerequisite and downstream lists.
func (m *Manager) DeleteTask(id string) error {
	if id == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	if _, exists := m.tasks[id]; !exists {
		return fmt.Errorf("task with ID %s not found", id)
	}

	// Remove references to this task from all other tasks
	for _, task := range m.tasks {
		if task.ID == id {
			continue // Skip the task being deleted
		}

		// Remove from PrerequisiteIDs
		task.PrerequisiteIDs = removeStringFromSlice(task.PrerequisiteIDs, id)
		// Remove from Prerequisites pointers if populated
		if task.Prerequisites != nil {
			task.Prerequisites = removeTaskFromSlice(task.Prerequisites, id)
		}

		// Remove from DownstreamRequiredIDs
		task.DownstreamRequiredIDs = removeStringFromSlice(task.DownstreamRequiredIDs, id)
		// Remove from DownstreamRequired pointers if populated
		if task.DownstreamRequired != nil {
			task.DownstreamRequired = removeTaskFromSlice(task.DownstreamRequired, id)
		}

		// Remove from DownstreamSuggestedIDs
		task.DownstreamSuggestedIDs = removeStringFromSlice(task.DownstreamSuggestedIDs, id)
		// Remove from DownstreamSuggested pointers if populated
		if task.DownstreamSuggested != nil {
			task.DownstreamSuggested = removeTaskFromSlice(task.DownstreamSuggested, id)
		}
	}

	// Delete the task itself
	delete(m.tasks, id)

	// Update tag cache since a task was removed
	m.PopulateTagCache()

	return nil
}

// removeStringFromSlice removes all occurrences of a string from a slice
func removeStringFromSlice(slice []string, value string) []string {
	if slice == nil {
		return nil
	}

	result := make([]string, 0, len(slice))
	for _, v := range slice {
		if v != value {
			result = append(result, v)
		}
	}

	// Return nil if the result is empty to maintain nil vs empty slice distinction
	if len(result) == 0 && slice != nil {
		return []string{}
	}
	return result
}

// removeTaskFromSlice removes all tasks with the given ID from a task slice
func removeTaskFromSlice(slice []*types.Task, id string) []*types.Task {
	if slice == nil {
		return nil
	}

	result := make([]*types.Task, 0, len(slice))
	for _, task := range slice {
		if task != nil && task.ID != id {
			result = append(result, task)
		}
	}

	// Return nil if the result is empty to maintain nil vs empty slice distinction
	if len(result) == 0 && slice != nil {
		return []*types.Task{}
	}
	return result
}

// ListAllTasks returns all tasks in the manager
func (m *Manager) ListAllTasks() []*types.Task {
	tasks := make([]*types.Task, 0, len(m.tasks))
	for _, task := range m.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// GetTask retrieves a task by ID
func (m *Manager) GetTask(id string) (*types.Task, error) {
	if id == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}

	task, exists := m.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task with ID %s not found", id)
	}

	return task, nil
}

// getTasks retrieves multiple tasks by their IDs
func (m *Manager) getTasks(ids []string) ([]*types.Task, error) {
	if len(ids) == 0 {
		return []*types.Task{}, nil
	}

	tasks := make([]*types.Task, 0, len(ids))
	var notFound []string

	for _, id := range ids {
		if id == "" {
			return nil, fmt.Errorf("task ID cannot be empty")
		}

		task, exists := m.tasks[id]
		if !exists {
			notFound = append(notFound, id)
			continue
		}

		tasks = append(tasks, task)
	}

	if len(notFound) > 0 {
		return tasks, fmt.Errorf("tasks not found: %v", notFound)
	}

	return tasks, nil
}

// ResolveTaskPointers populates the task pointer fields (Prerequisites, DownstreamRequired,
// DownstreamSuggested) by looking up the corresponding IDs for all tasks in the manager.
// Should be called after loading tasks from disk to restore the pointer relationships.
// Returns an error if any referenced task IDs cannot be found.
func (m *Manager) ResolveTaskPointers() error {
	for _, task := range m.tasks {
		// Resolve prerequisites
		if len(task.PrerequisiteIDs) > 0 {
			prereqs, err := m.getTasks(task.PrerequisiteIDs)
			if err != nil {
				return fmt.Errorf("failed to resolve prerequisites for task %s: %w", task.ID, err)
			}
			task.Prerequisites = prereqs
		}

		// Resolve downstream required
		if len(task.DownstreamRequiredIDs) > 0 {
			downstream, err := m.getTasks(task.DownstreamRequiredIDs)
			if err != nil {
				return fmt.Errorf("failed to resolve downstream required for task %s: %w", task.ID, err)
			}
			task.DownstreamRequired = downstream
		}

		// Resolve downstream suggested
		if len(task.DownstreamSuggestedIDs) > 0 {
			suggested, err := m.getTasks(task.DownstreamSuggestedIDs)
			if err != nil {
				return fmt.Errorf("failed to resolve downstream suggested for task %s: %w", task.ID, err)
			}
			task.DownstreamSuggested = suggested
		}
	}

	return nil
}

// Clone creates a deep copy of the manager and all its tasks.
// The cloned manager has independent tasks with resolved pointers.
// This is useful for making transactional changes that can be validated before committing.
func (m *Manager) Clone() *Manager {
	if m == nil {
		return nil
	}

	// Create new manager
	clone := NewManager()

	// Clone all tasks
	for id, task := range m.tasks {
		clonedTask := task.Clone()
		clone.tasks[id] = clonedTask
	}

	// Resolve task pointers in the cloned manager
	// Note: We ignore errors here because if the original manager was valid,
	// the clone should also be valid. If there are resolution errors, they
	// would have existed in the original manager too.
	_ = clone.ResolveTaskPointers()

	// Clone tag cache (we'll just rebuild it)
	clone.PopulateTagCache()

	return clone
}

// DetectCycles checks all three DAGs (Prerequisites, Downstream Required, and Downstream Suggested)
// for cycles. Returns an error if any cycles are detected, with detailed information about all cycles found.
func (m *Manager) DetectCycles() error {
	var allCycles []string

	// Check Prerequisites DAG for cycles
	cycles := m.detectCyclesInDAG("prerequisites", func(task *types.Task) []string {
		return task.PrerequisiteIDs
	})
	if len(cycles) > 0 {
		for _, cycle := range cycles {
			allCycles = append(allCycles, fmt.Sprintf("Prerequisites DAG: %s", cycle))
		}
	}

	// Check Downstream Required DAG for cycles
	cycles = m.detectCyclesInDAG("downstream required", func(task *types.Task) []string {
		return task.DownstreamRequiredIDs
	})
	if len(cycles) > 0 {
		for _, cycle := range cycles {
			allCycles = append(allCycles, fmt.Sprintf("Downstream Required DAG: %s", cycle))
		}
	}

	// Check Downstream Suggested DAG for cycles
	cycles = m.detectCyclesInDAG("downstream suggested", func(task *types.Task) []string {
		return task.DownstreamSuggestedIDs
	})
	if len(cycles) > 0 {
		for _, cycle := range cycles {
			allCycles = append(allCycles, fmt.Sprintf("Downstream Suggested DAG: %s", cycle))
		}
	}

	if len(allCycles) > 0 {
		msg := fmt.Sprintf("detected %d cycle(s):\n", len(allCycles))
		for i, cycle := range allCycles {
			msg += fmt.Sprintf("  %d. %s\n", i+1, cycle)
		}
		return fmt.Errorf("%s", msg)
	}

	return nil
}

// detectCyclesInDAG performs cycle detection on a specific DAG using DFS
// Returns a slice of cycle descriptions (e.g., "task-a -> task-b -> task-c -> task-a")
func (m *Manager) detectCyclesInDAG(dagName string, getEdges func(*types.Task) []string) []string {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	var cycles []string
	path := []string{}

	// Check each task as a potential starting point
	for taskID := range m.tasks {
		if !visited[taskID] {
			m.findCyclesDFS(taskID, visited, recStack, &path, &cycles, getEdges)
		}
	}

	return cycles
}

// findCyclesDFS performs depth-first search to find all cycles
func (m *Manager) findCyclesDFS(taskID string, visited, recStack map[string]bool, path *[]string, cycles *[]string, getEdges func(*types.Task) []string) {
	// Mark current node as visited and add to recursion stack
	visited[taskID] = true
	recStack[taskID] = true
	*path = append(*path, taskID)

	// Get the task
	task, exists := m.tasks[taskID]
	if exists {
		// Get edges for this task based on the DAG we're checking
		edges := getEdges(task)

		// Recursively check all adjacent nodes
		for _, adjacentID := range edges {
			// If adjacent node is not visited, recurse on it
			if !visited[adjacentID] {
				m.findCyclesDFS(adjacentID, visited, recStack, path, cycles, getEdges)
			} else if recStack[adjacentID] {
				// If adjacent node is in recursion stack, we found a cycle
				// Find where the cycle starts in the path
				cycleStart := -1
				for i, id := range *path {
					if id == adjacentID {
						cycleStart = i
						break
					}
				}

				// Build the cycle description
				if cycleStart >= 0 {
					cyclePath := append((*path)[cycleStart:], adjacentID)
					cycleDesc := ""
					for i, id := range cyclePath {
						if i > 0 {
							cycleDesc += " -> "
						}
						cycleDesc += id
					}
					*cycles = append(*cycles, cycleDesc)
				}
			}
		}
	}

	// Remove from recursion stack and path before returning
	recStack[taskID] = false
	*path = (*path)[:len(*path)-1]
}

// LoadFromDir reads all YAML files from the specified directory and loads tasks
func (m *Manager) LoadFromDir(dirPath string) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Read all .yaml files in the directory
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}

		filename := filepath.Join(dirPath, entry.Name())
		data, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", entry.Name(), err)
		}

		var task types.Task
		if err := yaml.Unmarshal(data, &task); err != nil {
			return fmt.Errorf("failed to unmarshal task from %s: %w", entry.Name(), err)
		}

		m.tasks[task.ID] = &task
	}

	// Detect cycles before resolving pointers
	if err := m.DetectCycles(); err != nil {
		return fmt.Errorf("cycle detected in task graph: %w", err)
	}

	// Resolve task pointers after loading all tasks and validating no cycles
	if err := m.ResolveTaskPointers(); err != nil {
		return err
	}

	// Populate tag cache for efficient tag-based lookups
	m.PopulateTagCache()

	return nil
}

// PersistToDir writes all tasks to the specified directory as YAML files
func (m *Manager) PersistToDir(dirPath string) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write each task as a separate YAML file
	for id, task := range m.tasks {
		filename := filepath.Join(dirPath, fmt.Sprintf("%s.yaml", id))

		data, err := yaml.Marshal(task)
		if err != nil {
			return fmt.Errorf("failed to marshal task %s: %w", id, err)
		}

		if err := os.WriteFile(filename, data, 0644); err != nil {
			return fmt.Errorf("failed to write task %s: %w", id, err)
		}
	}

	return nil
}
