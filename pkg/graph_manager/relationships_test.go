package graph_manager

import (
	"os"
	"path/filepath"
	"testing"

	"common-tasks-mcp/pkg/graph_manager/types"
	"common-tasks-mcp/pkg/logger"
)

func TestRegisterRelationship(t *testing.T) {
	log, _ := logger.New(false)
	manager := NewManager(log)

	tests := []struct {
		name          string
		relationship  types.Relationship
		wantError     bool
		expectedError string
	}{
		{
			name: "valid relationship",
			relationship: types.Relationship{
				Name:        "prerequisites",
				Description: "Tasks that must complete before this one",
				Direction:   types.DirectionBackward,
			},
			wantError: false,
		},
		{
			name: "duplicate relationship",
			relationship: types.Relationship{
				Name:        "prerequisites",
				Description: "Duplicate",
				Direction:   types.DirectionBackward,
			},
			wantError:     true,
			expectedError: "already registered",
		},
		{
			name: "invalid relationship - empty name",
			relationship: types.Relationship{
				Name:        "",
				Description: "Invalid",
				Direction:   types.DirectionBackward,
			},
			wantError:     true,
			expectedError: "invalid relationship",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.RegisterRelationship(tt.relationship)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Verify it was registered
				if !manager.IsRelationshipRegistered(tt.relationship.Name) {
					t.Errorf("Relationship %s was not registered", tt.relationship.Name)
				}
			}
		})
	}
}

func TestGetRelationship(t *testing.T) {
	log, _ := logger.New(false)
	manager := NewManager(log)

	// Register a relationship
	rel := types.Relationship{
		Name:        "test-rel",
		Description: "Test relationship",
		Direction:   types.DirectionForward,
	}

	if err := manager.RegisterRelationship(rel); err != nil {
		t.Fatalf("Failed to register relationship: %v", err)
	}

	// Test getting registered relationship
	retrieved := manager.GetRelationship("test-rel")
	if retrieved == nil {
		t.Fatal("Expected to retrieve relationship but got nil")
	}

	if retrieved.Name != rel.Name {
		t.Errorf("Expected name %s, got %s", rel.Name, retrieved.Name)
	}

	// Test getting non-existent relationship
	missing := manager.GetRelationship("non-existent")
	if missing != nil {
		t.Errorf("Expected nil for non-existent relationship, got %v", missing)
	}
}

func TestGetAllRelationships(t *testing.T) {
	log, _ := logger.New(false)
	manager := NewManager(log)

	// Register multiple relationships
	relationships := []types.Relationship{
		{Name: "rel1", Description: "First", Direction: types.DirectionBackward},
		{Name: "rel2", Description: "Second", Direction: types.DirectionForward},
		{Name: "rel3", Description: "Third", Direction: types.DirectionNone},
	}

	for _, rel := range relationships {
		if err := manager.RegisterRelationship(rel); err != nil {
			t.Fatalf("Failed to register relationship %s: %v", rel.Name, err)
		}
	}

	// Get all relationships
	all := manager.GetAllRelationships()

	if len(all) != len(relationships) {
		t.Errorf("Expected %d relationships, got %d", len(relationships), len(all))
	}

	for _, rel := range relationships {
		if _, exists := all[rel.Name]; !exists {
			t.Errorf("Relationship %s not found in GetAllRelationships", rel.Name)
		}
	}

	// Verify it's a copy (modification doesn't affect original)
	all["rel1"] = types.Relationship{Name: "modified", Description: "Modified", Direction: types.DirectionBackward}

	original := manager.GetRelationship("rel1")
	if original.Description == "Modified" {
		t.Error("Modifying returned map should not affect original relationships")
	}
}

func TestLoadRelationshipsFromFile(t *testing.T) {
	log, _ := logger.New(false)

	// Create a temporary directory
	tempDir := t.TempDir()

	tests := []struct {
		name           string
		yamlContent    string
		wantError      bool
		expectedCount  int
		validateResult func(*Manager) error
	}{
		{
			name: "valid relationships file",
			yamlContent: `relationships:
  - name: prerequisites
    description: Tasks that must complete before this one
    direction: backward
  - name: downstream_required
    description: Tasks that must complete after this one
    direction: forward
  - name: downstream_suggested
    description: Suggested tasks after this one
    direction: forward
`,
			wantError:     false,
			expectedCount: 3,
		},
		{
			name: "empty relationships",
			yamlContent: `relationships: []
`,
			wantError:     false,
			expectedCount: 0,
		},
		{
			name: "invalid yaml",
			yamlContent: `relationships:
  - name: bad
    description: Invalid
    direction: invalid_direction
`,
			wantError:     false, // File loads but relationship is skipped
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewManager(log)

			// Create test file
			filePath := filepath.Join(tempDir, tt.name+".yaml")
			if err := os.WriteFile(filePath, []byte(tt.yamlContent), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Load relationships
			err := manager.LoadRelationshipsFromFile(filePath)

			if tt.wantError && err == nil {
				t.Errorf("Expected error but got nil")
			}

			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check count
			all := manager.GetAllRelationships()
			if len(all) != tt.expectedCount {
				t.Errorf("Expected %d relationships, got %d", tt.expectedCount, len(all))
			}
		})
	}
}

func TestLoadRelationshipsFromDir(t *testing.T) {
	log, _ := logger.New(false)

	t.Run("directory with relationships.yaml", func(t *testing.T) {
		tempDir := t.TempDir()
		manager := NewManager(log)

		yamlContent := `relationships:
  - name: test-rel
    description: Test relationship
    direction: backward
`

		relFile := filepath.Join(tempDir, "relationships.yaml")
		if err := os.WriteFile(relFile, []byte(yamlContent), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		err := manager.LoadRelationshipsFromDir(tempDir)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !manager.IsRelationshipRegistered("test-rel") {
			t.Error("Expected test-rel to be registered")
		}
	})

	t.Run("directory without relationships.yaml", func(t *testing.T) {
		tempDir := t.TempDir()
		manager := NewManager(log)

		// Should not error if file doesn't exist
		err := manager.LoadRelationshipsFromDir(tempDir)
		if err != nil {
			t.Errorf("Unexpected error when file doesn't exist: %v", err)
		}

		all := manager.GetAllRelationships()
		if len(all) != 0 {
			t.Errorf("Expected 0 relationships, got %d", len(all))
		}
	})
}

func TestValidateRelationships(t *testing.T) {
	log, _ := logger.New(false)

	t.Run("all relationships registered", func(t *testing.T) {
		manager := NewManager(log)

		// Register relationships
		manager.RegisterRelationship(types.Relationship{
			Name:        "prerequisites",
			Description: "Test",
			Direction:   types.DirectionBackward,
		})

		// Add a node using the registered relationship
		node := &types.Node{
			ID:   "test-node",
			Name: "Test Node",
			EdgeIDs: map[string][]string{
				"prerequisites": {"other-node"},
			},
		}
		manager.AddNode(node)

		// Validate should pass
		err := manager.ValidateRelationships()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("unregistered relationship in use", func(t *testing.T) {
		manager := NewManager(log)

		// Add a node using an unregistered relationship
		node := &types.Node{
			ID:   "test-node",
			Name: "Test Node",
			EdgeIDs: map[string][]string{
				"unregistered": {"other-node"},
			},
		}
		manager.AddNode(node)

		// Validate should fail
		err := manager.ValidateRelationships()
		if err == nil {
			t.Error("Expected error for unregistered relationship")
		}
	})

	t.Run("no relationships used", func(t *testing.T) {
		manager := NewManager(log)

		// Add a node with no relationships
		node := &types.Node{
			ID:   "test-node",
			Name: "Test Node",
		}
		manager.AddNode(node)

		// Validate should pass
		err := manager.ValidateRelationships()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}
