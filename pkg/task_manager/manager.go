package task_manager

import (
	"fmt"
	"os"
	"path/filepath"

	"common-tasks-mcp/pkg/task_manager/types"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// Manager handles task graph operations
type Manager struct {
	tasks    map[string]*types.Task
	tagCache map[string][]*types.Task
	logger   *zap.Logger
}

// NewManager creates a new task manager instance
func NewManager(logger *zap.Logger) *Manager {
	logger.Debug("Creating new task manager")
	return &Manager{
		tasks:    make(map[string]*types.Task),
		tagCache: make(map[string][]*types.Task),
		logger:   logger,
	}
}

// AddTask adds a task to the manager.
// It uses a clone-validate-commit pattern to ensure the addition doesn't introduce cycles.
func (m *Manager) AddTask(task *types.Task) error {
	m.logger.Debug("Adding task")

	if task == nil {
		m.logger.Error("Attempted to add nil task")
		return fmt.Errorf("task cannot be nil")
	}
	if task.ID == "" {
		m.logger.Error("Attempted to add task with empty ID")
		return fmt.Errorf("task ID cannot be empty")
	}
	if _, exists := m.tasks[task.ID]; exists {
		m.logger.Warn("Task already exists", zap.String("task_id", task.ID))
		return fmt.Errorf("task with ID %s already exists", task.ID)
	}

	m.logger.Debug("Validating task addition for cycles", zap.String("task_id", task.ID))

	// Clone the manager to test the addition
	testManager := m.Clone()

	// Perform the addition in the test manager
	testManager.tasks[task.ID] = task

	// Check for cycles in the test manager
	if err := testManager.DetectCycles(); err != nil {
		m.logger.Error("Task addition would introduce cycle",
			zap.String("task_id", task.ID),
			zap.Error(err),
		)
		return fmt.Errorf("addition would introduce cycle: %w", err)
	}

	// If no cycles detected, commit the addition to the original manager
	m.tasks[task.ID] = task
	m.logger.Debug("Task added to internal storage", zap.String("task_id", task.ID))

	// Update tag cache with the new task
	m.PopulateTagCache()
	m.logger.Info("Task added successfully",
		zap.String("task_id", task.ID),
		zap.String("task_name", task.Name),
		zap.Int("total_tasks", len(m.tasks)),
	)

	return nil
}

// UpdateTask updates an existing task in the manager.
// It uses a clone-validate-commit pattern to ensure the update doesn't introduce cycles,
// and automatically refreshes all task pointers to prevent stale references.
func (m *Manager) UpdateTask(task *types.Task) error {
	m.logger.Debug("Updating task")

	if task == nil {
		m.logger.Error("Attempted to update with nil task")
		return fmt.Errorf("task cannot be nil")
	}
	if task.ID == "" {
		m.logger.Error("Attempted to update task with empty ID")
		return fmt.Errorf("task ID cannot be empty")
	}
	if _, exists := m.tasks[task.ID]; !exists {
		m.logger.Warn("Task not found for update", zap.String("task_id", task.ID))
		return fmt.Errorf("task with ID %s not found", task.ID)
	}

	m.logger.Debug("Validating task update for cycles", zap.String("task_id", task.ID))

	// Clone the manager to test the update
	testManager := m.Clone()

	// Perform the update in the test manager
	testManager.tasks[task.ID] = task

	// Check for cycles in the test manager
	if err := testManager.DetectCycles(); err != nil {
		m.logger.Error("Task update would introduce cycle",
			zap.String("task_id", task.ID),
			zap.Error(err),
		)
		return fmt.Errorf("update would introduce cycle: %w", err)
	}

	// If no cycles detected, commit the update to the original manager
	m.tasks[task.ID] = task
	m.logger.Debug("Task updated in internal storage", zap.String("task_id", task.ID))

	// Resolve all task pointers to fix stale references
	// This ensures that any tasks pointing to the updated task get fresh pointers
	m.logger.Debug("Resolving task pointers after update")
	if err := m.ResolveTaskPointers(); err != nil {
		m.logger.Error("Failed to resolve task pointers", zap.Error(err))
		return err
	}

	// Update tag cache since tags may have changed
	m.PopulateTagCache()
	m.logger.Info("Task updated successfully",
		zap.String("task_id", task.ID),
		zap.String("task_name", task.Name),
	)

	return nil
}

// DeleteTask removes a task from the manager and cleans up all references to it
// from other tasks' prerequisite and downstream lists.
func (m *Manager) DeleteTask(id string) error {
	m.logger.Debug("Deleting task", zap.String("task_id", id))

	if id == "" {
		m.logger.Error("Attempted to delete task with empty ID")
		return fmt.Errorf("task ID cannot be empty")
	}
	if _, exists := m.tasks[id]; !exists {
		m.logger.Warn("Task not found for deletion", zap.String("task_id", id))
		return fmt.Errorf("task with ID %s not found", id)
	}

	m.logger.Debug("Removing task references from other tasks", zap.String("task_id", id))

	// Remove references to this task from all other tasks
	referencesRemoved := 0
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
			referencesRemoved++
		}
	}

	m.logger.Debug("Removed references from other tasks",
		zap.String("task_id", id),
		zap.Int("references_cleaned", referencesRemoved),
	)

	// Delete the task itself
	delete(m.tasks, id)

	// Update tag cache since a task was removed
	m.PopulateTagCache()
	m.logger.Info("Task deleted successfully",
		zap.String("task_id", id),
		zap.Int("remaining_tasks", len(m.tasks)),
	)

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

	m.logger.Debug("Cloning manager", zap.Int("task_count", len(m.tasks)))

	// Create new manager with same logger
	clone := &Manager{
		tasks:    make(map[string]*types.Task),
		tagCache: make(map[string][]*types.Task),
		logger:   m.logger,
	}

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

	m.logger.Debug("Manager cloned successfully")

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
	m.logger.Info("Loading tasks from directory", zap.String("path", dirPath))

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		m.logger.Error("Failed to create directory", zap.String("path", dirPath), zap.Error(err))
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Read all .yaml files in the directory
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		m.logger.Error("Failed to read directory", zap.String("path", dirPath), zap.Error(err))
		return fmt.Errorf("failed to read directory: %w", err)
	}

	m.logger.Debug("Found directory entries", zap.Int("count", len(entries)))

	tasksLoaded := 0
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" {
			m.logger.Debug("Skipping non-YAML file", zap.String("filename", entry.Name()))
			continue
		}

		filename := filepath.Join(dirPath, entry.Name())
		m.logger.Debug("Reading task file", zap.String("filename", filename))

		data, err := os.ReadFile(filename)
		if err != nil {
			m.logger.Error("Failed to read file", zap.String("filename", entry.Name()), zap.Error(err))
			return fmt.Errorf("failed to read file %s: %w", entry.Name(), err)
		}

		var task types.Task
		if err := yaml.Unmarshal(data, &task); err != nil {
			m.logger.Error("Failed to unmarshal task",
				zap.String("filename", entry.Name()),
				zap.Error(err),
			)
			return fmt.Errorf("failed to unmarshal task from %s: %w", entry.Name(), err)
		}

		m.tasks[task.ID] = &task
		tasksLoaded++
		m.logger.Debug("Loaded task from file",
			zap.String("task_id", task.ID),
			zap.String("task_name", task.Name),
			zap.String("filename", entry.Name()),
		)
	}

	m.logger.Info("Finished loading task files", zap.Int("tasks_loaded", tasksLoaded))

	// Detect cycles before resolving pointers
	m.logger.Debug("Detecting cycles in task graph")
	if err := m.DetectCycles(); err != nil {
		m.logger.Error("Cycle detected in task graph", zap.Error(err))
		return fmt.Errorf("cycle detected in task graph: %w", err)
	}
	m.logger.Debug("No cycles detected")

	// Resolve task pointers after loading all tasks and validating no cycles
	m.logger.Debug("Resolving task pointers")
	if err := m.ResolveTaskPointers(); err != nil {
		m.logger.Error("Failed to resolve task pointers", zap.Error(err))
		return err
	}
	m.logger.Debug("Task pointers resolved")

	// Populate tag cache for efficient tag-based lookups
	m.logger.Debug("Populating tag cache")
	m.PopulateTagCache()

	m.logger.Info("Successfully loaded tasks from directory",
		zap.String("path", dirPath),
		zap.Int("total_tasks", len(m.tasks)),
	)

	return nil
}

// PersistToDir writes all tasks to the specified directory as YAML files
func (m *Manager) PersistToDir(dirPath string) error {
	m.logger.Info("Persisting tasks to directory",
		zap.String("path", dirPath),
		zap.Int("task_count", len(m.tasks)),
	)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		m.logger.Error("Failed to create directory", zap.String("path", dirPath), zap.Error(err))
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write each task as a separate YAML file
	tasksPersisted := 0
	for id, task := range m.tasks {
		filename := filepath.Join(dirPath, fmt.Sprintf("%s.yaml", id))

		m.logger.Debug("Marshaling task", zap.String("task_id", id))
		data, err := yaml.Marshal(task)
		if err != nil {
			m.logger.Error("Failed to marshal task", zap.String("task_id", id), zap.Error(err))
			return fmt.Errorf("failed to marshal task %s: %w", id, err)
		}

		m.logger.Debug("Writing task file", zap.String("filename", filename))
		if err := os.WriteFile(filename, data, 0644); err != nil {
			m.logger.Error("Failed to write task file",
				zap.String("task_id", id),
				zap.String("filename", filename),
				zap.Error(err),
			)
			return fmt.Errorf("failed to write task %s: %w", id, err)
		}

		tasksPersisted++
	}

	m.logger.Info("Successfully persisted tasks to directory",
		zap.String("path", dirPath),
		zap.Int("tasks_persisted", tasksPersisted),
	)

	return nil
}
