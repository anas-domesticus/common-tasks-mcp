package types

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestNodeYAMLMarshalling(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	// Create a node with EdgeIDs populated
	node := &Node{
		ID:          "test-1",
		Name:        "Test Node",
		Summary:     "Summary",
		Description: "Description",
		Tags:        []string{"tag1", "tag2"},
		EdgeIDs: map[string][]string{
			"prerequisites":        {"prereq-1", "prereq-2"},
			"downstream_required":  {"required-1"},
			"downstream_suggested": {"suggested-1", "suggested-2"},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Marshal to YAML
	data, err := yaml.Marshal(node)
	if err != nil {
		t.Fatalf("Failed to marshal node to YAML: %v", err)
	}

	yamlStr := string(data)

	// Verify that edges field is present
	if !strings.Contains(yamlStr, "edges:") {
		t.Error("Expected 'edges:' field in YAML output")
	}

	// Unmarshal back and verify
	var unmarshalled Node
	if err := yaml.Unmarshal(data, &unmarshalled); err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// Verify EdgeIDs were preserved
	if len(unmarshalled.EdgeIDs) != 3 {
		t.Errorf("Expected 3 edge types, got %d", len(unmarshalled.EdgeIDs))
	}
	if len(unmarshalled.EdgeIDs["prerequisites"]) != 2 {
		t.Errorf("Expected 2 prerequisites, got %d", len(unmarshalled.EdgeIDs["prerequisites"]))
	}
	if len(unmarshalled.EdgeIDs["downstream_required"]) != 1 {
		t.Errorf("Expected 1 downstream required, got %d", len(unmarshalled.EdgeIDs["downstream_required"]))
	}

	// Verify Edges map (runtime) is nil after unmarshalling
	if unmarshalled.Edges != nil {
		t.Errorf("Edges map should be nil after unmarshal, got %v", unmarshalled.Edges)
	}
}

func TestNodeJSONMarshalling(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	node := &Node{
		ID:          "test-1",
		Name:        "Test Node",
		Summary:     "Summary",
		Description: "Description",
		Tags:        []string{"tag1", "tag2"},
		EdgeIDs: map[string][]string{
			"prerequisites": {"prereq-1", "prereq-2"},
			"validates":     {"test-1"},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Marshal to JSON
	data, err := json.Marshal(node)
	if err != nil {
		t.Fatalf("Failed to marshal node to JSON: %v", err)
	}

	jsonStr := string(data)

	// Verify that edges field is present
	if !strings.Contains(jsonStr, `"edges"`) {
		t.Error("Expected 'edges' field in JSON output")
	}
	if !strings.Contains(jsonStr, `"prerequisites"`) {
		t.Error("Expected 'prerequisites' in JSON output")
	}

	// Unmarshal back and verify
	var unmarshalled Node
	if err := json.Unmarshal(data, &unmarshalled); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify EdgeIDs were preserved
	if len(unmarshalled.EdgeIDs["prerequisites"]) != 2 {
		t.Errorf("Expected 2 prerequisites, got %d", len(unmarshalled.EdgeIDs["prerequisites"]))
	}
}

func TestNodeEquals(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	baseNode := &Node{
		ID:          "test-1",
		Name:        "Test Node",
		Summary:     "Summary",
		Description: "Description",
		Tags:        []string{"tag1", "tag2"},
		EdgeIDs: map[string][]string{
			"prerequisites": {"prereq-1", "prereq-2"},
			"validates":     {"test-1"},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		name        string
		node1       *Node
		node2       *Node
		shouldEqual bool
	}{
		{
			name:        "same node equals itself",
			node1:       baseNode,
			node2:       baseNode,
			shouldEqual: true,
		},
		{
			name:  "identical node should be equal",
			node1: baseNode,
			node2: &Node{
				ID:          "test-1",
				Name:        "Test Node",
				Summary:     "Summary",
				Description: "Description",
				Tags:        []string{"tag1", "tag2"},
				EdgeIDs: map[string][]string{
					"prerequisites": {"prereq-1", "prereq-2"},
					"validates":     {"test-1"},
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
			shouldEqual: true,
		},
		{
			name:        "both nil should be equal",
			node1:       nil,
			node2:       nil,
			shouldEqual: true,
		},
		{
			name:        "node and nil should not be equal",
			node1:       baseNode,
			node2:       nil,
			shouldEqual: false,
		},
		{
			name:  "different EdgeIDs should not be equal",
			node1: baseNode,
			node2: &Node{
				ID:          "test-1",
				Name:        "Test Node",
				Summary:     "Summary",
				Description: "Description",
				Tags:        []string{"tag1", "tag2"},
				EdgeIDs: map[string][]string{
					"prerequisites": {"prereq-3"},
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
			shouldEqual: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.node1.Equals(tt.node2)
			if got != tt.shouldEqual {
				t.Errorf("Equals() = %v, want %v", got, tt.shouldEqual)
			}
		})
	}
}

func TestNodeClone(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	node := &Node{
		ID:          "test-1",
		Name:        "Test Node",
		Summary:     "Summary",
		Description: "Description",
		Tags:        []string{"tag1", "tag2"},
		EdgeIDs: map[string][]string{
			"prerequisites": {"prereq-1", "prereq-2"},
			"validates":     {"test-1"},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	clone := node.Clone()

	// Verify clone is not nil
	if clone == nil {
		t.Fatal("Clone() returned nil for non-nil node")
	}

	// Verify clone is a different instance
	if node == clone {
		t.Error("Clone() returned the same pointer")
	}

	// Verify clone equals original
	if !node.Equals(clone) {
		t.Error("Clone() should equal the original node")
	}

	// Verify modifying clone doesn't affect original
	clone.EdgeIDs["prerequisites"][0] = "modified"
	if node.EdgeIDs["prerequisites"][0] == "modified" {
		t.Error("Modifying clone.EdgeIDs should not affect original")
	}
}

func TestNodeEdgeMethods(t *testing.T) {
	t.Run("GetEdgeIDs", func(t *testing.T) {
		node := &Node{
			EdgeIDs: map[string][]string{
				"prerequisites": {"prereq-1", "prereq-2"},
			},
		}

		ids := node.GetEdgeIDs("prerequisites")
		if len(ids) != 2 {
			t.Errorf("Expected 2 IDs, got %d", len(ids))
		}

		// Non-existent relationship
		ids = node.GetEdgeIDs("nonexistent")
		if ids != nil {
			t.Errorf("Expected nil for non-existent relationship, got %v", ids)
		}
	})

	t.Run("SetEdgeIDs", func(t *testing.T) {
		node := &Node{}
		err := node.SetEdgeIDs("prerequisites", []string{"prereq-1", "prereq-2"})
		if err != nil {
			t.Fatalf("SetEdgeIDs failed: %v", err)
		}

		ids := node.GetEdgeIDs("prerequisites")
		if len(ids) != 2 {
			t.Errorf("Expected 2 IDs, got %d", len(ids))
		}
	})

	t.Run("AddEdgeID", func(t *testing.T) {
		node := &Node{}
		err := node.AddEdgeID("prerequisites", "prereq-1")
		if err != nil {
			t.Fatalf("AddEdgeID failed: %v", err)
		}

		err = node.AddEdgeID("prerequisites", "prereq-2")
		if err != nil {
			t.Fatalf("AddEdgeID failed: %v", err)
		}

		ids := node.GetEdgeIDs("prerequisites")
		if len(ids) != 2 {
			t.Errorf("Expected 2 IDs, got %d", len(ids))
		}
	})

	t.Run("SetEdges", func(t *testing.T) {
		node1 := &Node{ID: "node-1"}
		node2 := &Node{ID: "node-2"}
		relationship := &Relationship{Name: "prerequisites"}

		node := &Node{}
		edges := []Edge{
			{To: node1, Type: relationship},
			{To: node2, Type: relationship},
		}

		err := node.SetEdges("prerequisites", edges)
		if err != nil {
			t.Fatalf("SetEdges failed: %v", err)
		}

		// Verify both Edges and EdgeIDs are set
		gotEdges := node.GetEdges("prerequisites")
		if len(gotEdges) != 2 {
			t.Errorf("Expected 2 edges, got %d", len(gotEdges))
		}

		gotIDs := node.GetEdgeIDs("prerequisites")
		if len(gotIDs) != 2 {
			t.Errorf("Expected 2 IDs, got %d", len(gotIDs))
		}
		if gotIDs[0] != "node-1" || gotIDs[1] != "node-2" {
			t.Errorf("Unexpected IDs: %v", gotIDs)
		}
	})

	t.Run("AddEdge", func(t *testing.T) {
		node1 := &Node{ID: "node-1"}
		relationship := &Relationship{Name: "prerequisites"}

		node := &Node{}
		err := node.AddEdge("prerequisites", Edge{To: node1, Type: relationship})
		if err != nil {
			t.Fatalf("AddEdge failed: %v", err)
		}

		edges := node.GetEdges("prerequisites")
		if len(edges) != 1 {
			t.Errorf("Expected 1 edge, got %d", len(edges))
		}

		ids := node.GetEdgeIDs("prerequisites")
		if len(ids) != 1 || ids[0] != "node-1" {
			t.Errorf("Unexpected IDs: %v", ids)
		}
	})
}
