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

	// Create nodes with various tags
	nodeA := &types.Node{
		ID:          "task-a",
		Name:        "Node A",
		Summary:     "First node",
		Description: "Node with backend and api tags",
		Tags:        []string{"backend", "api"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	nodeB := &types.Node{
		ID:          "task-b",
		Name:        "Node B",
		Summary:     "Second node",
		Description: "Node with frontend and api tags",
		Tags:        []string{"frontend", "api"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	nodeC := &types.Node{
		ID:          "task-c",
		Name:        "Node C",
		Summary:     "Third node",
		Description: "Node with testing tag",
		Tags:        []string{"testing"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	nodeD := &types.Node{
		ID:          "task-d",
		Name:        "Node D",
		Summary:     "Fourth node",
		Description: "Node with backend tag",
		Tags:        []string{"backend"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Add nodes to manager
	if err := manager.AddNode(nodeA); err != nil {
		t.Fatalf("Failed to add node A: %v", err)
	}
	if err := manager.AddNode(nodeB); err != nil {
		t.Fatalf("Failed to add node B: %v", err)
	}
	if err := manager.AddNode(nodeC); err != nil {
		t.Fatalf("Failed to add node C: %v", err)
	}
	if err := manager.AddNode(nodeD); err != nil {
		t.Fatalf("Failed to add node D: %v", err)
	}

	// Populate the tag cache
	manager.PopulateTagCache()

	// Verify tag cache structure
	if len(manager.tagCache) != 4 {
		t.Errorf("Expected 4 unique tags in cache, got %d", len(manager.tagCache))
	}

	// Verify "api" tag has 2 nodes (A and B)
	apiNodes, exists := manager.tagCache["api"]
	if !exists {
		t.Error("Expected 'api' tag to exist in cache")
	} else if len(apiNodes) != 2 {
		t.Errorf("Expected 2 nodes with 'api' tag, got %d", len(apiNodes))
	} else {
		// Verify the correct nodes are present
		foundA, foundB := false, false
		for _, node := range apiNodes {
			if node.ID == "task-a" {
				foundA = true
			}
			if node.ID == "task-b" {
				foundB = true
			}
		}
		if !foundA || !foundB {
			t.Error("Expected nodes A and B to be tagged with 'api'")
		}
	}

	// Verify "backend" tag has 2 nodes (A and D)
	backendNodes, exists := manager.tagCache["backend"]
	if !exists {
		t.Error("Expected 'backend' tag to exist in cache")
	} else if len(backendNodes) != 2 {
		t.Errorf("Expected 2 nodes with 'backend' tag, got %d", len(backendNodes))
	} else {
		foundA, foundD := false, false
		for _, node := range backendNodes {
			if node.ID == "task-a" {
				foundA = true
			}
			if node.ID == "task-d" {
				foundD = true
			}
		}
		if !foundA || !foundD {
			t.Error("Expected nodes A and D to be tagged with 'backend'")
		}
	}

	// Verify "frontend" tag has 1 node (B)
	frontendNodes, exists := manager.tagCache["frontend"]
	if !exists {
		t.Error("Expected 'frontend' tag to exist in cache")
	} else if len(frontendNodes) != 1 {
		t.Errorf("Expected 1 node with 'frontend' tag, got %d", len(frontendNodes))
	} else if frontendNodes[0].ID != "task-b" {
		t.Errorf("Expected node B to be tagged with 'frontend', got %s", frontendNodes[0].ID)
	}

	// Verify "testing" tag has 1 node (C)
	testingNodes, exists := manager.tagCache["testing"]
	if !exists {
		t.Error("Expected 'testing' tag to exist in cache")
	} else if len(testingNodes) != 1 {
		t.Errorf("Expected 1 node with 'testing' tag, got %d", len(testingNodes))
	} else if testingNodes[0].ID != "task-c" {
		t.Errorf("Expected node C to be tagged with 'testing', got %s", testingNodes[0].ID)
	}

	// Test cache repopulation (should clear and rebuild)
	nodeE := &types.Node{
		ID:          "task-e",
		Name:        "Node E",
		Summary:     "Fifth node",
		Description: "Additional node with api tag",
		Tags:        []string{"api"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := manager.AddNode(nodeE); err != nil {
		t.Fatalf("Failed to add node E: %v", err)
	}

	// Repopulate cache
	manager.PopulateTagCache()

	// Verify "api" tag now has 3 nodes (A, B, and E)
	apiNodes, exists = manager.tagCache["api"]
	if !exists {
		t.Error("Expected 'api' tag to exist in cache after repopulation")
	} else if len(apiNodes) != 3 {
		t.Errorf("Expected 3 nodes with 'api' tag after repopulation, got %d", len(apiNodes))
	}
}

func TestGetNodesByTag(t *testing.T) {
	log, _ := logger.New(false)
	manager := NewManager(log)
	now := time.Now().UTC().Truncate(time.Second)

	// Create nodes with various tags
	nodeA := &types.Node{
		ID:          "task-a",
		Name:        "Node A",
		Summary:     "First node",
		Description: "Node with backend and api tags",
		Tags:        []string{"backend", "api"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	nodeB := &types.Node{
		ID:          "task-b",
		Name:        "Node B",
		Summary:     "Second node",
		Description: "Node with frontend and api tags",
		Tags:        []string{"frontend", "api"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	nodeC := &types.Node{
		ID:          "task-c",
		Name:        "Node C",
		Summary:     "Third node",
		Description: "Node with testing tag",
		Tags:        []string{"testing"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Add nodes
	if err := manager.AddNode(nodeA); err != nil {
		t.Fatalf("Failed to add node A: %v", err)
	}
	if err := manager.AddNode(nodeB); err != nil {
		t.Fatalf("Failed to add node B: %v", err)
	}
	if err := manager.AddNode(nodeC); err != nil {
		t.Fatalf("Failed to add node C: %v", err)
	}

	// Populate tag cache
	manager.PopulateTagCache()

	// Test retrieving nodes by existing tag with multiple nodes
	apiNodes, err := manager.GetNodesByTag("api")
	if err != nil {
		t.Fatalf("Expected no error when retrieving nodes by 'api' tag, got: %v", err)
	}
	if len(apiNodes) != 2 {
		t.Errorf("Expected 2 nodes with 'api' tag, got %d", len(apiNodes))
	}
	foundA, foundB := false, false
	for _, node := range apiNodes {
		if node.ID == "task-a" {
			foundA = true
		}
		if node.ID == "task-b" {
			foundB = true
		}
	}
	if !foundA || !foundB {
		t.Error("Expected nodes A and B to be returned for 'api' tag")
	}

	// Test retrieving nodes by existing tag with single node
	testingNodes, err := manager.GetNodesByTag("testing")
	if err != nil {
		t.Fatalf("Expected no error when retrieving nodes by 'testing' tag, got: %v", err)
	}
	if len(testingNodes) != 1 {
		t.Errorf("Expected 1 node with 'testing' tag, got %d", len(testingNodes))
	}
	if testingNodes[0].ID != "task-c" {
		t.Errorf("Expected node C for 'testing' tag, got %s", testingNodes[0].ID)
	}

	// Test retrieving nodes by non-existent tag
	nonExistentNodes, err := manager.GetNodesByTag("non-existent")
	if err != nil {
		t.Fatalf("Expected no error when retrieving nodes by non-existent tag, got: %v", err)
	}
	if len(nonExistentNodes) != 0 {
		t.Errorf("Expected 0 nodes for non-existent tag, got %d", len(nonExistentNodes))
	}

	// Test retrieving with empty tag
	_, err = manager.GetNodesByTag("")
	if err == nil {
		t.Error("Expected error when retrieving with empty tag, got nil")
	} else if err.Error() != "tag cannot be empty" {
		t.Errorf("Expected 'tag cannot be empty' error, got: %v", err)
	}
}

func TestGetAllTags(t *testing.T) {
	log, _ := logger.New(false)
	manager := NewManager(log)
	now := time.Now().UTC().Truncate(time.Second)

	// Test with empty manager
	tags := manager.GetAllTags()
	if len(tags) != 0 {
		t.Errorf("Expected 0 tags from empty manager, got %d", len(tags))
	}

	// Create nodes with various tags
	nodeA := &types.Node{
		ID:          "task-a",
		Name:        "Node A",
		Summary:     "First node",
		Description: "Node with backend and api tags",
		Tags:        []string{"backend", "api"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	nodeB := &types.Node{
		ID:          "task-b",
		Name:        "Node B",
		Summary:     "Second node",
		Description: "Node with frontend and api tags",
		Tags:        []string{"frontend", "api"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	nodeC := &types.Node{
		ID:          "task-c",
		Name:        "Node C",
		Summary:     "Third node",
		Description: "Node with testing tag",
		Tags:        []string{"testing"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	nodeD := &types.Node{
		ID:          "task-d",
		Name:        "Node D",
		Summary:     "Fourth node",
		Description: "Node with backend tag",
		Tags:        []string{"backend"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Add nodes
	if err := manager.AddNode(nodeA); err != nil {
		t.Fatalf("Failed to add node A: %v", err)
	}
	if err := manager.AddNode(nodeB); err != nil {
		t.Fatalf("Failed to add node B: %v", err)
	}
	if err := manager.AddNode(nodeC); err != nil {
		t.Fatalf("Failed to add node C: %v", err)
	}
	if err := manager.AddNode(nodeD); err != nil {
		t.Fatalf("Failed to add node D: %v", err)
	}

	// Get all tags
	tags = manager.GetAllTags()

	// Verify we have 4 unique tags
	if len(tags) != 4 {
		t.Errorf("Expected 4 unique tags, got %d", len(tags))
	}

	// Verify tag counts
	expectedCounts := map[string]int{
		"api":      2, // nodes A and B
		"backend":  2, // nodes A and D
		"frontend": 1, // node B
		"testing":  1, // node C
	}

	for tag, expectedCount := range expectedCounts {
		count, exists := tags[tag]
		if !exists {
			t.Errorf("Expected tag '%s' to exist, but it doesn't", tag)
		} else if count != expectedCount {
			t.Errorf("Expected tag '%s' to have count %d, got %d", tag, expectedCount, count)
		}
	}

	// Verify no unexpected tags
	for tag := range tags {
		if _, expected := expectedCounts[tag]; !expected {
			t.Errorf("Unexpected tag '%s' found in results", tag)
		}
	}
}
