package task_manager

import (
	"fmt"
	"os"
	"path/filepath"

	"common-tasks-mcp/pkg/types"

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
