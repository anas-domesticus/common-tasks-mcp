package task_manager

import (
	"testing"
	"time"

	"common-tasks-mcp/pkg/types"
)

func TestPopulateTagCache(t *testing.T) {
	manager := NewManager()
	now := time.Now().UTC().Truncate(time.Second)

	// Create tasks with various tags
	taskA := &types.Task{
		ID:          "task-a",
		Name:        "Task A",
		Summary:     "First task",
		Description: "Task with backend and api tags",
		Tags:        []string{"backend", "api"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	taskB := &types.Task{
		ID:          "task-b",
		Name:        "Task B",
		Summary:     "Second task",
		Description: "Task with frontend and api tags",
		Tags:        []string{"frontend", "api"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	taskC := &types.Task{
		ID:          "task-c",
		Name:        "Task C",
		Summary:     "Third task",
		Description: "Task with testing tag",
		Tags:        []string{"testing"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	taskD := &types.Task{
		ID:          "task-d",
		Name:        "Task D",
		Summary:     "Fourth task",
		Description: "Task with backend tag",
		Tags:        []string{"backend"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Add tasks to manager
	if err := manager.AddTask(taskA); err != nil {
		t.Fatalf("Failed to add task A: %v", err)
	}
	if err := manager.AddTask(taskB); err != nil {
		t.Fatalf("Failed to add task B: %v", err)
	}
	if err := manager.AddTask(taskC); err != nil {
		t.Fatalf("Failed to add task C: %v", err)
	}
	if err := manager.AddTask(taskD); err != nil {
		t.Fatalf("Failed to add task D: %v", err)
	}

	// Populate the tag cache
	manager.PopulateTagCache()

	// Verify tag cache structure
	if len(manager.tagCache) != 4 {
		t.Errorf("Expected 4 unique tags in cache, got %d", len(manager.tagCache))
	}

	// Verify "api" tag has 2 tasks (A and B)
	apiTasks, exists := manager.tagCache["api"]
	if !exists {
		t.Error("Expected 'api' tag to exist in cache")
	} else if len(apiTasks) != 2 {
		t.Errorf("Expected 2 tasks with 'api' tag, got %d", len(apiTasks))
	} else {
		// Verify the correct tasks are present
		foundA, foundB := false, false
		for _, task := range apiTasks {
			if task.ID == "task-a" {
				foundA = true
			}
			if task.ID == "task-b" {
				foundB = true
			}
		}
		if !foundA || !foundB {
			t.Error("Expected tasks A and B to be tagged with 'api'")
		}
	}

	// Verify "backend" tag has 2 tasks (A and D)
	backendTasks, exists := manager.tagCache["backend"]
	if !exists {
		t.Error("Expected 'backend' tag to exist in cache")
	} else if len(backendTasks) != 2 {
		t.Errorf("Expected 2 tasks with 'backend' tag, got %d", len(backendTasks))
	} else {
		foundA, foundD := false, false
		for _, task := range backendTasks {
			if task.ID == "task-a" {
				foundA = true
			}
			if task.ID == "task-d" {
				foundD = true
			}
		}
		if !foundA || !foundD {
			t.Error("Expected tasks A and D to be tagged with 'backend'")
		}
	}

	// Verify "frontend" tag has 1 task (B)
	frontendTasks, exists := manager.tagCache["frontend"]
	if !exists {
		t.Error("Expected 'frontend' tag to exist in cache")
	} else if len(frontendTasks) != 1 {
		t.Errorf("Expected 1 task with 'frontend' tag, got %d", len(frontendTasks))
	} else if frontendTasks[0].ID != "task-b" {
		t.Errorf("Expected task B to be tagged with 'frontend', got %s", frontendTasks[0].ID)
	}

	// Verify "testing" tag has 1 task (C)
	testingTasks, exists := manager.tagCache["testing"]
	if !exists {
		t.Error("Expected 'testing' tag to exist in cache")
	} else if len(testingTasks) != 1 {
		t.Errorf("Expected 1 task with 'testing' tag, got %d", len(testingTasks))
	} else if testingTasks[0].ID != "task-c" {
		t.Errorf("Expected task C to be tagged with 'testing', got %s", testingTasks[0].ID)
	}

	// Test cache repopulation (should clear and rebuild)
	taskE := &types.Task{
		ID:          "task-e",
		Name:        "Task E",
		Summary:     "Fifth task",
		Description: "Additional task with api tag",
		Tags:        []string{"api"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := manager.AddTask(taskE); err != nil {
		t.Fatalf("Failed to add task E: %v", err)
	}

	// Repopulate cache
	manager.PopulateTagCache()

	// Verify "api" tag now has 3 tasks (A, B, and E)
	apiTasks, exists = manager.tagCache["api"]
	if !exists {
		t.Error("Expected 'api' tag to exist in cache after repopulation")
	} else if len(apiTasks) != 3 {
		t.Errorf("Expected 3 tasks with 'api' tag after repopulation, got %d", len(apiTasks))
	}
}

func TestGetTasksByTag(t *testing.T) {
	manager := NewManager()
	now := time.Now().UTC().Truncate(time.Second)

	// Create tasks with various tags
	taskA := &types.Task{
		ID:          "task-a",
		Name:        "Task A",
		Summary:     "First task",
		Description: "Task with backend and api tags",
		Tags:        []string{"backend", "api"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	taskB := &types.Task{
		ID:          "task-b",
		Name:        "Task B",
		Summary:     "Second task",
		Description: "Task with frontend and api tags",
		Tags:        []string{"frontend", "api"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	taskC := &types.Task{
		ID:          "task-c",
		Name:        "Task C",
		Summary:     "Third task",
		Description: "Task with testing tag",
		Tags:        []string{"testing"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Add tasks
	if err := manager.AddTask(taskA); err != nil {
		t.Fatalf("Failed to add task A: %v", err)
	}
	if err := manager.AddTask(taskB); err != nil {
		t.Fatalf("Failed to add task B: %v", err)
	}
	if err := manager.AddTask(taskC); err != nil {
		t.Fatalf("Failed to add task C: %v", err)
	}

	// Populate tag cache
	manager.PopulateTagCache()

	// Test retrieving tasks by existing tag with multiple tasks
	apiTasks, err := manager.GetTasksByTag("api")
	if err != nil {
		t.Fatalf("Expected no error when retrieving tasks by 'api' tag, got: %v", err)
	}
	if len(apiTasks) != 2 {
		t.Errorf("Expected 2 tasks with 'api' tag, got %d", len(apiTasks))
	}
	foundA, foundB := false, false
	for _, task := range apiTasks {
		if task.ID == "task-a" {
			foundA = true
		}
		if task.ID == "task-b" {
			foundB = true
		}
	}
	if !foundA || !foundB {
		t.Error("Expected tasks A and B to be returned for 'api' tag")
	}

	// Test retrieving tasks by existing tag with single task
	testingTasks, err := manager.GetTasksByTag("testing")
	if err != nil {
		t.Fatalf("Expected no error when retrieving tasks by 'testing' tag, got: %v", err)
	}
	if len(testingTasks) != 1 {
		t.Errorf("Expected 1 task with 'testing' tag, got %d", len(testingTasks))
	}
	if testingTasks[0].ID != "task-c" {
		t.Errorf("Expected task C for 'testing' tag, got %s", testingTasks[0].ID)
	}

	// Test retrieving tasks by non-existent tag
	nonExistentTasks, err := manager.GetTasksByTag("non-existent")
	if err != nil {
		t.Fatalf("Expected no error when retrieving tasks by non-existent tag, got: %v", err)
	}
	if len(nonExistentTasks) != 0 {
		t.Errorf("Expected 0 tasks for non-existent tag, got %d", len(nonExistentTasks))
	}

	// Test retrieving with empty tag
	_, err = manager.GetTasksByTag("")
	if err == nil {
		t.Error("Expected error when retrieving with empty tag, got nil")
	} else if err.Error() != "tag cannot be empty" {
		t.Errorf("Expected 'tag cannot be empty' error, got: %v", err)
	}
}
