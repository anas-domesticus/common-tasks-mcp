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

// AddTask adds a task to the manager
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

	m.tasks[task.ID] = task
	return nil
}

// UpdateTask updates an existing task in the manager
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

	m.tasks[task.ID] = task
	return nil
}

// DeleteTask removes a task from the manager
func (m *Manager) DeleteTask(id string) error {
	if id == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	if _, exists := m.tasks[id]; !exists {
		return fmt.Errorf("task with ID %s not found", id)
	}

	delete(m.tasks, id)
	return nil
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

// GetPrerequisites retrieves all prerequisite tasks for the given task
func (m *Manager) GetPrerequisites(task *types.Task) ([]*types.Task, error) {
	if task == nil {
		return nil, fmt.Errorf("task cannot be nil")
	}

	return m.getTasks(task.PrerequisiteIDs)
}

// GetDownstreamRequired retrieves all required downstream tasks for the given task
func (m *Manager) GetDownstreamRequired(task *types.Task) ([]*types.Task, error) {
	if task == nil {
		return nil, fmt.Errorf("task cannot be nil")
	}

	return m.getTasks(task.DownstreamRequiredIDs)
}

// GetDownstreamSuggested retrieves all suggested downstream tasks for the given task
func (m *Manager) GetDownstreamSuggested(task *types.Task) ([]*types.Task, error) {
	if task == nil {
		return nil, fmt.Errorf("task cannot be nil")
	}

	return m.getTasks(task.DownstreamSuggestedIDs)
}

// DetectCycles checks all three DAGs (Prerequisites, Downstream Required, and Downstream Suggested)
// for cycles. Returns an error if any cycles are detected.
func (m *Manager) DetectCycles() error {
	var errors []error

	// Check Prerequisites DAG for cycles
	if err := m.detectCyclesInDAG("prerequisites", func(task *types.Task) []string {
		return task.PrerequisiteIDs
	}); err != nil {
		errors = append(errors, fmt.Errorf("cycle detected in prerequisites DAG: %w", err))
	}

	// Check Downstream Required DAG for cycles
	if err := m.detectCyclesInDAG("downstream required", func(task *types.Task) []string {
		return task.DownstreamRequiredIDs
	}); err != nil {
		errors = append(errors, fmt.Errorf("cycle detected in downstream required DAG: %w", err))
	}

	// Check Downstream Suggested DAG for cycles
	if err := m.detectCyclesInDAG("downstream suggested", func(task *types.Task) []string {
		return task.DownstreamSuggestedIDs
	}); err != nil {
		errors = append(errors, fmt.Errorf("cycle detected in downstream suggested DAG: %w", err))
	}

	if len(errors) > 0 {
		// Combine all errors into one
		msg := ""
		for i, err := range errors {
			if i > 0 {
				msg += "; "
			}
			msg += err.Error()
		}
		return fmt.Errorf("%s", msg)
	}

	return nil
}

// detectCyclesInDAG performs cycle detection on a specific DAG using DFS
func (m *Manager) detectCyclesInDAG(dagName string, getEdges func(*types.Task) []string) error {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	// Check each task as a potential starting point
	for taskID := range m.tasks {
		if !visited[taskID] {
			if m.hasCycleDFS(taskID, visited, recStack, getEdges) {
				return fmt.Errorf("cycle found starting from task %s", taskID)
			}
		}
	}

	return nil
}

// hasCycleDFS performs depth-first search to detect cycles
func (m *Manager) hasCycleDFS(taskID string, visited, recStack map[string]bool, getEdges func(*types.Task) []string) bool {
	// Mark current node as visited and add to recursion stack
	visited[taskID] = true
	recStack[taskID] = true

	// Get the task
	task, exists := m.tasks[taskID]
	if !exists {
		// If task doesn't exist, we can't traverse it, so no cycle from this path
		recStack[taskID] = false
		return false
	}

	// Get edges for this task based on the DAG we're checking
	edges := getEdges(task)

	// Recursively check all adjacent nodes
	for _, adjacentID := range edges {
		// If adjacent node is not visited, recurse on it
		if !visited[adjacentID] {
			if m.hasCycleDFS(adjacentID, visited, recStack, getEdges) {
				return true
			}
		} else if recStack[adjacentID] {
			// If adjacent node is in recursion stack, we found a cycle
			return true
		}
	}

	// Remove from recursion stack before returning
	recStack[taskID] = false
	return false
}

// Load reads all YAML files from the specified directory and loads tasks
func (m *Manager) Load(dirPath string) error {
	// Check if directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", dirPath)
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

	return nil
}

// Persist writes all tasks to the specified directory as YAML files
func (m *Manager) Persist(dirPath string) error {
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
