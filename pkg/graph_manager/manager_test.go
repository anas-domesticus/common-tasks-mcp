package graph_manager

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"common-tasks-mcp/pkg/graph_manager/types"
	"common-tasks-mcp/pkg/logger"
)

func TestPersistAndLoad(t *testing.T) {
	// Create a temporary directory for test files
	testDir := filepath.Join(t.TempDir(), "tasks")

	// Create a manager with 3 tasks
	log, _ := logger.New(false)
	manager1 := NewManager(log)

	now := time.Now().UTC().Truncate(time.Second)

	// Node A - no dependencies
	nodeA := &types.Node{
		ID:          "task-a",
		Name:        "Task A",
		Summary:     "First task",
		Description: "This is the first task with no dependencies",
		Tags:        []string{"api", "backend"},
		EdgeIDs: map[string][]string{
			"downstream_required": {"task-b"},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Node B - depends on A, has C as dependent
	nodeB := &types.Node{
		ID:          "task-b",
		Name:        "Task B",
		Summary:     "Second task",
		Description: "This task depends on A and has C as dependent",
		Tags:        []string{"frontend", "api"},
		EdgeIDs: map[string][]string{
			"prerequisites":       {"task-a"},
			"downstream_required": {"task-c"},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Node C - depends on B
	nodeC := &types.Node{
		ID:          "task-c",
		Name:        "Task C",
		Summary:     "Third task",
		Description: "This task depends on B",
		Tags:        []string{"testing"},
		EdgeIDs: map[string][]string{
			"prerequisites": {"task-b"},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	manager1.nodes["task-a"] = nodeA
	manager1.nodes["task-b"] = nodeB
	manager1.nodes["task-c"] = nodeC

	// Persist tasks to directory
	if err := manager1.PersistToDir(testDir); err != nil {
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
	log2, _ := logger.New(false)
	manager2 := NewManager(log2)
	if err := manager2.LoadFromDir(testDir); err != nil {
		t.Fatalf("Failed to load tasks: %v", err)
	}

	// Compare the managers
	if len(manager1.nodes) != len(manager2.nodes) {
		t.Errorf("Node count mismatch: expected %d, got %d", len(manager1.nodes), len(manager2.nodes))
	}

	// Compare each node
	for id, originalNode := range manager1.nodes {
		loadedNode, exists := manager2.nodes[id]
		if !exists {
			t.Errorf("Node %s not found in loaded manager", id)
			continue
		}

		if !originalNode.Equals(loadedNode) {
			t.Errorf("Node %s does not match after persist/load cycle", id)
		}
	}
}

func TestAddNode(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	tests := []struct {
		name           string
		setupNodes     []*types.Node // Nodes to add before the test node
		nodeToAdd      *types.Node
		wantError      bool
		expectedError  string
		expectedCount  int // Expected number of nodes after operation
		validateResult func(t *testing.T, manager *Manager)
	}{
		{
			name:       "add valid node to empty manager",
			setupNodes: []*types.Node{},
			nodeToAdd: &types.Node{
				ID:          "task-1",
				Name:        "Test Task 1",
				Summary:     "First test task",
				Description: "A valid task to add",
				Tags:        []string{"test"},
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantError:     false,
			expectedCount: 1,
			validateResult: func(t *testing.T, manager *Manager) {
				addedNode, exists := manager.nodes["task-1"]
				if !exists {
					t.Error("Node was not added to manager")
				} else if addedNode.ID != "task-1" {
					t.Errorf("Added node ID mismatch: expected task-1, got %s", addedNode.ID)
				}
			},
		},
		{
			name:          "add nil node",
			setupNodes:    []*types.Node{},
			nodeToAdd:     nil,
			wantError:     true,
			expectedError: "node cannot be nil",
			expectedCount: 0,
		},
		{
			name:       "add node with empty ID",
			setupNodes: []*types.Node{},
			nodeToAdd: &types.Node{
				ID:        "",
				Name:      "No ID Node",
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantError:     true,
			expectedError: "node ID cannot be empty",
			expectedCount: 0,
		},
		{
			name: "add duplicate node",
			setupNodes: []*types.Node{
				{
					ID:        "task-1",
					Name:      "Existing Node",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			nodeToAdd: &types.Node{
				ID:        "task-1",
				Name:      "Duplicate Node",
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantError:     true,
			expectedError: "node with ID task-1 already exists",
			expectedCount: 1,
		},
		{
			name: "add node with valid prerequisite reference",
			setupNodes: []*types.Node{
				{
					ID:        "prereq-1",
					Name:      "Prerequisite",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			nodeToAdd: &types.Node{
				ID:        "task-2",
				Name:      "Task with Prereq",
				EdgeIDs:   map[string][]string{"prerequisites": {"prereq-1"}},
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantError:     false,
			expectedCount: 2,
		},
		{
			name: "add node that would create a cycle",
			setupNodes: []*types.Node{
				{
					ID:        "task-a",
					Name:      "Task A",
					EdgeIDs:   map[string][]string{"prerequisites": {"task-b"}},
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			nodeToAdd: &types.Node{
				ID:        "task-b",
				Name:      "Task B",
				EdgeIDs:   map[string][]string{"prerequisites": {"task-a"}},
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantError:     true,
			expectedError: "addition would introduce cycle",
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log, _ := logger.New(false)
			manager := NewManager(log)

			// Setup nodes
			for _, node := range tt.setupNodes {
				manager.nodes[node.ID] = node
			}

			// Attempt to add the test node
			err := manager.AddNode(tt.nodeToAdd)

			// Check error expectation
			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			// Check node count
			if len(manager.nodes) != tt.expectedCount {
				t.Errorf("Expected %d nodes, got %d", tt.expectedCount, len(manager.nodes))
			}

			// Run custom validation if provided
			if tt.validateResult != nil && !tt.wantError {
				tt.validateResult(t, manager)
			}
		})
	}
}

func TestGetNode(t *testing.T) {
	log, _ := logger.New(false)
	manager := NewManager(log)
	now := time.Now().UTC().Truncate(time.Second)

	// Add some test nodes
	node1 := &types.Node{
		ID:        "task-1",
		Name:      "Task 1",
		Summary:   "First task",
		Tags:      []string{"tag1"},
		CreatedAt: now,
		UpdatedAt: now,
	}
	node2 := &types.Node{
		ID:        "task-2",
		Name:      "Task 2",
		Summary:   "Second task",
		Tags:      []string{"tag2"},
		CreatedAt: now,
		UpdatedAt: now,
	}

	manager.nodes["task-1"] = node1
	manager.nodes["task-2"] = node2

	tests := []struct {
		name      string
		id        string
		wantNode  *types.Node
		wantError bool
	}{
		{
			name:      "get existing node",
			id:        "task-1",
			wantNode:  node1,
			wantError: false,
		},
		{
			name:      "get another existing node",
			id:        "task-2",
			wantNode:  node2,
			wantError: false,
		},
		{
			name:      "get non-existent node",
			id:        "task-99",
			wantNode:  nil,
			wantError: true,
		},
		{
			name:      "get with empty ID",
			id:        "",
			wantNode:  nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := manager.GetNode(tt.id)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if got != tt.wantNode {
					t.Errorf("Got wrong node: expected %v, got %v", tt.wantNode, got)
				}
			}
		})
	}
}

func TestDeleteNode(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	tests := []struct {
		name          string
		setupNodes    []*types.Node
		deleteID      string
		wantError     bool
		expectedCount int
		validateGraph func(t *testing.T, manager *Manager)
	}{
		{
			name: "delete node with no references",
			setupNodes: []*types.Node{
				{
					ID:        "task-1",
					Name:      "Task 1",
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "task-2",
					Name:      "Task 2",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			deleteID:      "task-1",
			wantError:     false,
			expectedCount: 1,
			validateGraph: func(t *testing.T, manager *Manager) {
				if _, exists := manager.nodes["task-1"]; exists {
					t.Error("task-1 should be deleted")
				}
				if _, exists := manager.nodes["task-2"]; !exists {
					t.Error("task-2 should still exist")
				}
			},
		},
		{
			name: "delete node and clean up references",
			setupNodes: []*types.Node{
				{
					ID:        "task-a",
					Name:      "Task A",
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:   "task-b",
					Name: "Task B",
					EdgeIDs: map[string][]string{
						"prerequisites": {"task-a"},
					},
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:   "task-c",
					Name: "Task C",
					EdgeIDs: map[string][]string{
						"prerequisites": {"task-a"},
					},
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			deleteID:      "task-a",
			wantError:     false,
			expectedCount: 2,
			validateGraph: func(t *testing.T, manager *Manager) {
				// task-a should be deleted
				if _, exists := manager.nodes["task-a"]; exists {
					t.Error("task-a should be deleted")
				}

				// task-b should no longer reference task-a
				taskB, exists := manager.nodes["task-b"]
				if !exists {
					t.Fatal("task-b should exist")
				}
				if prereqs, ok := taskB.EdgeIDs["prerequisites"]; ok && len(prereqs) > 0 {
					for _, id := range prereqs {
						if id == "task-a" {
							t.Error("task-b should not reference deleted task-a")
						}
					}
				}

				// task-c should no longer reference task-a
				taskC, exists := manager.nodes["task-c"]
				if !exists {
					t.Fatal("task-c should exist")
				}
				if prereqs, ok := taskC.EdgeIDs["prerequisites"]; ok && len(prereqs) > 0 {
					for _, id := range prereqs {
						if id == "task-a" {
							t.Error("task-c should not reference deleted task-a")
						}
					}
				}
			},
		},
		{
			name:          "delete non-existent node",
			setupNodes:    []*types.Node{},
			deleteID:      "non-existent",
			wantError:     true,
			expectedCount: 0,
		},
		{
			name:          "delete with empty ID",
			setupNodes:    []*types.Node{},
			deleteID:      "",
			wantError:     true,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log, _ := logger.New(false)
			manager := NewManager(log)

			// Setup nodes
			for _, node := range tt.setupNodes {
				manager.nodes[node.ID] = node
			}

			// Attempt to delete
			err := manager.DeleteNode(tt.deleteID)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			// Check node count
			if len(manager.nodes) != tt.expectedCount {
				t.Errorf("Expected %d nodes, got %d", tt.expectedCount, len(manager.nodes))
			}

			// Run custom validation
			if tt.validateGraph != nil && !tt.wantError {
				tt.validateGraph(t, manager)
			}
		})
	}
}

func TestDetectCycles(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	tests := []struct {
		name       string
		setupNodes []*types.Node
		wantCycle  bool
	}{
		{
			name: "no cycles - simple chain",
			setupNodes: []*types.Node{
				{
					ID:        "task-a",
					Name:      "Task A",
					EdgeIDs:   map[string][]string{"prerequisites": {"task-b"}},
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "task-b",
					Name:      "Task B",
					EdgeIDs:   map[string][]string{"prerequisites": {"task-c"}},
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "task-c",
					Name:      "Task C",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			wantCycle: false,
		},
		{
			name: "cycle in prerequisites",
			setupNodes: []*types.Node{
				{
					ID:        "task-a",
					Name:      "Task A",
					EdgeIDs:   map[string][]string{"prerequisites": {"task-b"}},
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "task-b",
					Name:      "Task B",
					EdgeIDs:   map[string][]string{"prerequisites": {"task-a"}},
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			wantCycle: true,
		},
		{
			name: "cycle in custom relationship",
			setupNodes: []*types.Node{
				{
					ID:        "task-x",
					Name:      "Task X",
					EdgeIDs:   map[string][]string{"validates": {"task-y"}},
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "task-y",
					Name:      "Task Y",
					EdgeIDs:   map[string][]string{"validates": {"task-x"}},
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			wantCycle: true,
		},
		{
			name: "three-node cycle",
			setupNodes: []*types.Node{
				{
					ID:        "task-1",
					Name:      "Task 1",
					EdgeIDs:   map[string][]string{"downstream_required": {"task-2"}},
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "task-2",
					Name:      "Task 2",
					EdgeIDs:   map[string][]string{"downstream_required": {"task-3"}},
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "task-3",
					Name:      "Task 3",
					EdgeIDs:   map[string][]string{"downstream_required": {"task-1"}},
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			wantCycle: true,
		},
		{
			name: "no cycles - diamond pattern",
			setupNodes: []*types.Node{
				{
					ID:        "task-a",
					Name:      "Task A",
					EdgeIDs:   map[string][]string{"prerequisites": {"task-b", "task-c"}},
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "task-b",
					Name:      "Task B",
					EdgeIDs:   map[string][]string{"prerequisites": {"task-d"}},
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "task-c",
					Name:      "Task C",
					EdgeIDs:   map[string][]string{"prerequisites": {"task-d"}},
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "task-d",
					Name:      "Task D",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			wantCycle: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log, _ := logger.New(false)
			manager := NewManager(log)

			// Setup nodes
			for _, node := range tt.setupNodes {
				manager.nodes[node.ID] = node
			}

			err := manager.DetectCycles()

			if tt.wantCycle {
				if err == nil {
					t.Error("Expected cycle detection error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected cycle detection error: %v", err)
				}
			}
		})
	}
}

func TestListAllNodes(t *testing.T) {
	log, _ := logger.New(false)
	manager := NewManager(log)
	now := time.Now().UTC().Truncate(time.Second)

	// Empty manager
	nodes := manager.ListAllNodes()
	if len(nodes) != 0 {
		t.Errorf("Empty manager should return 0 nodes, got %d", len(nodes))
	}

	// Add some nodes
	for i := 1; i <= 3; i++ {
		node := &types.Node{
			ID:        string(rune('a' + i - 1)),
			Name:      "Task " + string(rune('A'+i-1)),
			CreatedAt: now,
			UpdatedAt: now,
		}
		manager.nodes[node.ID] = node
	}

	nodes = manager.ListAllNodes()
	if len(nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(nodes))
	}
}

func TestResolveNodePointers(t *testing.T) {
	log, _ := logger.New(false)
	manager := NewManager(log)
	now := time.Now().UTC().Truncate(time.Second)

	// Create nodes with EdgeIDs
	nodeA := &types.Node{
		ID:        "task-a",
		Name:      "Task A",
		CreatedAt: now,
		UpdatedAt: now,
	}

	nodeB := &types.Node{
		ID:        "task-b",
		Name:      "Task B",
		EdgeIDs:   map[string][]string{"prerequisites": {"task-a"}},
		CreatedAt: now,
		UpdatedAt: now,
	}

	nodeC := &types.Node{
		ID:        "task-c",
		Name:      "Task C",
		EdgeIDs:   map[string][]string{"prerequisites": {"task-a", "task-b"}},
		CreatedAt: now,
		UpdatedAt: now,
	}

	manager.nodes["task-a"] = nodeA
	manager.nodes["task-b"] = nodeB
	manager.nodes["task-c"] = nodeC

	// Resolve pointers
	err := manager.ResolveNodePointers()
	if err != nil {
		t.Fatalf("Failed to resolve node pointers: %v", err)
	}

	// Verify nodeB has resolved edges
	if nodeB.Edges == nil {
		t.Fatal("nodeB.Edges should be initialized")
	}
	prereqs, ok := nodeB.Edges["prerequisites"]
	if !ok {
		t.Fatal("nodeB should have prerequisites edges")
	}
	if len(prereqs) != 1 {
		t.Fatalf("nodeB should have 1 prerequisite, got %d", len(prereqs))
	}
	if prereqs[0].To != nodeA {
		t.Error("nodeB's prerequisite should point to nodeA")
	}

	// Verify nodeC has resolved edges
	if nodeC.Edges == nil {
		t.Fatal("nodeC.Edges should be initialized")
	}
	prereqsC, ok := nodeC.Edges["prerequisites"]
	if !ok {
		t.Fatal("nodeC should have prerequisites edges")
	}
	if len(prereqsC) != 2 {
		t.Fatalf("nodeC should have 2 prerequisites, got %d", len(prereqsC))
	}

	// Test error case - missing node
	nodeD := &types.Node{
		ID:        "task-d",
		Name:      "Task D",
		EdgeIDs:   map[string][]string{"prerequisites": {"non-existent"}},
		CreatedAt: now,
		UpdatedAt: now,
	}
	manager.nodes["task-d"] = nodeD

	err = manager.ResolveNodePointers()
	if err == nil {
		t.Error("Expected error when resolving pointer to non-existent node")
	}
}
