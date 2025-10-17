package task_manager

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"common-tasks-mcp/pkg/types"
)

func TestPersistAndLoad(t *testing.T) {
	// Create a temporary directory for test files
	testDir := filepath.Join(t.TempDir(), "tasks")

	// Create a manager with 3 tasks
	manager1 := NewManager()

	now := time.Now().UTC().Truncate(time.Second)

	// Task A - no dependencies
	taskA := &types.Task{
		ID:           "task-a",
		Name:         "Task A",
		Summary:      "First task",
		Description:  "This is the first task with no dependencies",
		Tags:         []string{"api", "backend"},
		Dependencies: []string{},
		Dependents:   []string{"task-b"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Task B - depends on A, has C as dependent
	taskB := &types.Task{
		ID:           "task-b",
		Name:         "Task B",
		Summary:      "Second task",
		Description:  "This task depends on A and has C as dependent",
		Tags:         []string{"frontend", "api"},
		Dependencies: []string{"task-a"},
		Dependents:   []string{"task-c"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Task C - depends on B
	taskC := &types.Task{
		ID:           "task-c",
		Name:         "Task C",
		Summary:      "Third task",
		Description:  "This task depends on B",
		Tags:         []string{"testing"},
		Dependencies: []string{"task-b"},
		Dependents:   []string{},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	manager1.tasks["task-a"] = taskA
	manager1.tasks["task-b"] = taskB
	manager1.tasks["task-c"] = taskC

	// Persist tasks to directory
	if err := manager1.Persist(testDir); err != nil {
		t.Fatalf("Failed to persist tasks: %v", err)
	}

	// Verify files were created
	for _, id := range []string{"task-a", "task-b", "task-c"} {
		filename := filepath.Join(testDir, id+".yaml")
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not created", filename)
		}
	}

	// Create a new manager and load tasks
	manager2 := NewManager()
	if err := manager2.Load(testDir); err != nil {
		t.Fatalf("Failed to load tasks: %v", err)
	}

	// Compare the managers
	if len(manager1.tasks) != len(manager2.tasks) {
		t.Errorf("Task count mismatch: expected %d, got %d", len(manager1.tasks), len(manager2.tasks))
	}

	// Compare each task
	for id, originalTask := range manager1.tasks {
		loadedTask, exists := manager2.tasks[id]
		if !exists {
			t.Errorf("Task %s not found in loaded manager", id)
			continue
		}

		if originalTask.ID != loadedTask.ID {
			t.Errorf("Task %s ID mismatch: expected %s, got %s", id, originalTask.ID, loadedTask.ID)
		}
		if originalTask.Name != loadedTask.Name {
			t.Errorf("Task %s Name mismatch: expected %s, got %s", id, originalTask.Name, loadedTask.Name)
		}
		if originalTask.Summary != loadedTask.Summary {
			t.Errorf("Task %s Summary mismatch: expected %s, got %s", id, originalTask.Summary, loadedTask.Summary)
		}
		if originalTask.Description != loadedTask.Description {
			t.Errorf("Task %s Description mismatch: expected %s, got %s", id, originalTask.Description, loadedTask.Description)
		}

		// Compare tags
		if len(originalTask.Tags) != len(loadedTask.Tags) {
			t.Errorf("Task %s Tags count mismatch: expected %d, got %d", id, len(originalTask.Tags), len(loadedTask.Tags))
		} else {
			for i, tag := range originalTask.Tags {
				if tag != loadedTask.Tags[i] {
					t.Errorf("Task %s Tag %d mismatch: expected %s, got %s", id, i, tag, loadedTask.Tags[i])
				}
			}
		}

		// Compare dependencies
		if len(originalTask.Dependencies) != len(loadedTask.Dependencies) {
			t.Errorf("Task %s Dependencies count mismatch: expected %d, got %d", id, len(originalTask.Dependencies), len(loadedTask.Dependencies))
		} else {
			for i, dep := range originalTask.Dependencies {
				if dep != loadedTask.Dependencies[i] {
					t.Errorf("Task %s Dependency %d mismatch: expected %s, got %s", id, i, dep, loadedTask.Dependencies[i])
				}
			}
		}

		// Compare dependents
		if len(originalTask.Dependents) != len(loadedTask.Dependents) {
			t.Errorf("Task %s Dependents count mismatch: expected %d, got %d", id, len(originalTask.Dependents), len(loadedTask.Dependents))
		} else {
			for i, dep := range originalTask.Dependents {
				if dep != loadedTask.Dependents[i] {
					t.Errorf("Task %s Dependent %d mismatch: expected %s, got %s", id, i, dep, loadedTask.Dependents[i])
				}
			}
		}

		// Compare timestamps (truncated to second precision due to YAML serialization)
		if !originalTask.CreatedAt.Equal(loadedTask.CreatedAt) {
			t.Errorf("Task %s CreatedAt mismatch: expected %v, got %v", id, originalTask.CreatedAt, loadedTask.CreatedAt)
		}
		if !originalTask.UpdatedAt.Equal(loadedTask.UpdatedAt) {
			t.Errorf("Task %s UpdatedAt mismatch: expected %v, got %v", id, originalTask.UpdatedAt, loadedTask.UpdatedAt)
		}
	}
}
