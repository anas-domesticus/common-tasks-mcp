package graph_manager

import (
	"testing"
	"time"

	"common-tasks-mcp/pkg/graph_manager/types"
	"common-tasks-mcp/pkg/logger"
)

func TestPopulateTagCache(t *testing.T) {
	log, _ := logger.New(false)
	manager := NewManager(log)
	now := time.Now().UTC().Truncate(time.Second)

	// Create tasks with various tags
	taskA := &types.Node{
		ID:          "task-a",
		Name:        "Node A",
		Summary:     "First task",
		Description: "Node with backend and api tags",
		Tags:        []string{"backend", "api"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	taskB := &types.Node{
		ID:          "task-b",
		Name:        "Node B",
		Summary:     "Second task",
		Description: "Node with frontend and api tags",
		Tags:        []string{"frontend", "api"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	taskC := &types.Node{
		ID:          "task-c",
		Name:        "Node C",
		Summary:     "Third task",
		Description: "Node with testing tag",
		Tags:        []string{"testing"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	taskD := &types.Node{
		ID:          "task-d",
		Name:        "Node D",
		Summary:     "Fourth task",
		Description: "Node with backend tag",
		Tags:        []string{"backend"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Add tasks to manager
	if err := manager.AddNode(taskA); err != nil {
		t.Fatalf("Failed to add task A: %v", err)
	}
	if err := manager.AddNode(taskB); err != nil {
		t.Fatalf("Failed to add task B: %v", err)
	}
	if err := manager.AddNode(taskC); err != nil {
		t.Fatalf("Failed to add task C: %v", err)
	}
	if err := manager.AddNode(taskD); err != nil {
		t.Fatalf("Failed to add task D: %v", err)
	}

	// Populate the tag cache
	manager.PopulateTagCache()

	// Verify tag cache structure
	if len(manager.tagCache) != 4 {
		t.Errorf("Expected 4 unique tags in cache, got %d", len(manager.tagCache))
	}

	// Verify "api" tag has 2 tasks (A and B)
	apiNodes, exists := manager.tagCache["api"]
	if !exists {
		t.Error("Expected 'api' tag to exist in cache")
	} else if len(apiNodes) != 2 {
		t.Errorf("Expected 2 tasks with 'api' tag, got %d", len(apiNodes))
	} else {
		// Verify the correct tasks are present
		foundA, foundB := false, false
		for _, task := range apiNodes {
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
	backendNodes, exists := manager.tagCache["backend"]
	if !exists {
		t.Error("Expected 'backend' tag to exist in cache")
	} else if len(backendNodes) != 2 {
		t.Errorf("Expected 2 tasks with 'backend' tag, got %d", len(backendNodes))
	} else {
		foundA, foundD := false, false
		for _, task := range backendNodes {
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
	frontendNodes, exists := manager.tagCache["frontend"]
	if !exists {
		t.Error("Expected 'frontend' tag to exist in cache")
	} else if len(frontendNodes) != 1 {
		t.Errorf("Expected 1 task with 'frontend' tag, got %d", len(frontendNodes))
	} else if frontendNodes[0].ID != "task-b" {
		t.Errorf("Expected task B to be tagged with 'frontend', got %s", frontendNodes[0].ID)
	}

	// Verify "testing" tag has 1 task (C)
	testingNodes, exists := manager.tagCache["testing"]
	if !exists {
		t.Error("Expected 'testing' tag to exist in cache")
	} else if len(testingNodes) != 1 {
		t.Errorf("Expected 1 task with 'testing' tag, got %d", len(testingNodes))
	} else if testingNodes[0].ID != "task-c" {
		t.Errorf("Expected task C to be tagged with 'testing', got %s", testingNodes[0].ID)
	}

	// Test cache repopulation (should clear and rebuild)
	taskE := &types.Node{
		ID:          "task-e",
		Name:        "Node E",
		Summary:     "Fifth task",
		Description: "Additional task with api tag",
		Tags:        []string{"api"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := manager.AddNode(taskE); err != nil {
		t.Fatalf("Failed to add task E: %v", err)
	}

	// Repopulate cache
	manager.PopulateTagCache()

	// Verify "api" tag now has 3 tasks (A, B, and E)
	apiNodes, exists = manager.tagCache["api"]
	if !exists {
		t.Error("Expected 'api' tag to exist in cache after repopulation")
	} else if len(apiNodes) != 3 {
		t.Errorf("Expected 3 tasks with 'api' tag after repopulation, got %d", len(apiNodes))
	}
}

func TestGetNodesByTag(t *testing.T) {
	log, _ := logger.New(false)
	manager := NewManager(log)
	now := time.Now().UTC().Truncate(time.Second)

	// Create tasks with various tags
	taskA := &types.Node{
		ID:          "task-a",
		Name:        "Node A",
		Summary:     "First task",
		Description: "Node with backend and api tags",
		Tags:        []string{"backend", "api"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	taskB := &types.Node{
		ID:          "task-b",
		Name:        "Node B",
		Summary:     "Second task",
		Description: "Node with frontend and api tags",
		Tags:        []string{"frontend", "api"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	taskC := &types.Node{
		ID:          "task-c",
		Name:        "Node C",
		Summary:     "Third task",
		Description: "Node with testing tag",
		Tags:        []string{"testing"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Add tasks
	if err := manager.AddNode(taskA); err != nil {
		t.Fatalf("Failed to add task A: %v", err)
	}
	if err := manager.AddNode(taskB); err != nil {
		t.Fatalf("Failed to add task B: %v", err)
	}
	if err := manager.AddNode(taskC); err != nil {
		t.Fatalf("Failed to add task C: %v", err)
	}

	// Populate tag cache
	manager.PopulateTagCache()

	// Test retrieving tasks by existing tag with multiple tasks
	apiNodes, err := manager.GetNodesByTag("api")
	if err != nil {
		t.Fatalf("Expected no error when retrieving tasks by 'api' tag, got: %v", err)
	}
	if len(apiNodes) != 2 {
		t.Errorf("Expected 2 tasks with 'api' tag, got %d", len(apiNodes))
	}
	foundA, foundB := false, false
	for _, task := range apiNodes {
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
	testingNodes, err := manager.GetNodesByTag("testing")
	if err != nil {
		t.Fatalf("Expected no error when retrieving tasks by 'testing' tag, got: %v", err)
	}
	if len(testingNodes) != 1 {
		t.Errorf("Expected 1 task with 'testing' tag, got %d", len(testingNodes))
	}
	if testingNodes[0].ID != "task-c" {
		t.Errorf("Expected task C for 'testing' tag, got %s", testingNodes[0].ID)
	}

	// Test retrieving tasks by non-existent tag
	nonExistentNodes, err := manager.GetNodesByTag("non-existent")
	if err != nil {
		t.Fatalf("Expected no error when retrieving tasks by non-existent tag, got: %v", err)
	}
	if len(nonExistentNodes) != 0 {
		t.Errorf("Expected 0 tasks for non-existent tag, got %d", len(nonExistentNodes))
	}

	// Test retrieving with empty tag
	_, err = manager.GetNodesByTag("")
	if err == nil {
		t.Error("Expected error when retrieving with empty tag, got nil")
	} else if err.Error() != "tag cannot be empty" {
		t.Errorf("Expected 'tag cannot be empty' error, got: %v", err)
	}
}
