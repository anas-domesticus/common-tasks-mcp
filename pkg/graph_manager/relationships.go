package graph_manager

import (
	"fmt"
	"os"
	"path/filepath"

	"common-tasks-mcp/pkg/graph_manager/types"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// RelationshipsConfig represents the YAML structure for relationship definitions
type RelationshipsConfig struct {
	Relationships []types.Relationship `yaml:"relationships"`
}

// RegisterRelationship registers a new relationship type with the manager.
// Returns an error if the relationship is invalid or already exists.
func (m *Manager) RegisterRelationship(rel types.Relationship) error {
	m.logger.Debug("Registering relationship", zap.String("name", rel.Name))

	// Validate the relationship
	if err := rel.Validate(); err != nil {
		m.logger.Error("Invalid relationship",
			zap.String("name", rel.Name),
			zap.Error(err),
		)
		return fmt.Errorf("invalid relationship %s: %w", rel.Name, err)
	}

	// Check if already exists
	if _, exists := m.relationshipTypes[rel.Name]; exists {
		m.logger.Warn("Relationship already registered", zap.String("name", rel.Name))
		return fmt.Errorf("relationship %s already registered", rel.Name)
	}

	m.relationshipTypes[rel.Name] = &rel
	m.logger.Info("Registered relationship",
		zap.String("name", rel.Name),
		zap.String("direction", string(rel.Direction)),
	)

	return nil
}

// LoadRelationshipsFromFile loads relationship definitions from a YAML file.
// The file should contain a "relationships" array with relationship definitions.
func (m *Manager) LoadRelationshipsFromFile(filePath string) error {
	m.logger.Info("Loading relationships from file", zap.String("path", filePath))

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		m.logger.Error("Failed to read relationships file",
			zap.String("path", filePath),
			zap.Error(err),
		)
		return fmt.Errorf("failed to read relationships file: %w", err)
	}

	// Parse YAML
	var config RelationshipsConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		m.logger.Error("Failed to parse relationships file",
			zap.String("path", filePath),
			zap.Error(err),
		)
		return fmt.Errorf("failed to parse relationships file: %w", err)
	}

	// Register each relationship
	registered := 0
	for _, rel := range config.Relationships {
		if err := m.RegisterRelationship(rel); err != nil {
			m.logger.Warn("Skipping invalid relationship",
				zap.String("name", rel.Name),
				zap.Error(err),
			)
			continue
		}
		registered++
	}

	m.logger.Info("Loaded relationships from file",
		zap.String("path", filePath),
		zap.Int("registered", registered),
		zap.Int("total", len(config.Relationships)),
	)

	return nil
}

// LoadRelationshipsFromDir loads relationships from a directory.
// It looks for a "relationships.yaml" file in the specified directory.
func (m *Manager) LoadRelationshipsFromDir(dirPath string) error {
	relPath := filepath.Join(dirPath, "relationships.yaml")

	// Check if the file exists
	if _, err := os.Stat(relPath); os.IsNotExist(err) {
		m.logger.Debug("No relationships file found, skipping",
			zap.String("path", relPath),
		)
		return nil // Not an error if the file doesn't exist
	}

	return m.LoadRelationshipsFromFile(relPath)
}

// GetRelationship retrieves a relationship by name.
// Returns nil if the relationship is not registered.
func (m *Manager) GetRelationship(name string) *types.Relationship {
	if rel, exists := m.relationshipTypes[name]; exists {
		return rel
	}
	return nil
}

// GetAllRelationships returns all registered relationships.
// Returns a copy to prevent external modification.
func (m *Manager) GetAllRelationships() map[string]types.Relationship {
	copy := make(map[string]types.Relationship, len(m.relationshipTypes))
	for k, v := range m.relationshipTypes {
		copy[k] = *v
	}
	return copy
}

// IsRelationshipRegistered checks if a relationship type is registered.
func (m *Manager) IsRelationshipRegistered(name string) bool {
	_, exists := m.relationshipTypes[name]
	return exists
}

// GetRegisteredRelationshipNames returns a list of all registered relationship names.
func (m *Manager) GetRegisteredRelationshipNames() []string {
	names := make([]string, 0, len(m.relationshipTypes))
	for name := range m.relationshipTypes {
		names = append(names, name)
	}
	return names
}

// ValidateRelationships checks that all relationship types used in nodes are registered.
// Returns an error listing any unregistered relationships found.
func (m *Manager) ValidateRelationships() error {
	m.logger.Debug("Validating all relationships are registered")

	unregistered := make(map[string]bool)

	// Check all nodes for unregistered relationship types
	for _, node := range m.nodes {
		if node.EdgeIDs != nil {
			for relationshipName := range node.EdgeIDs {
				if !m.IsRelationshipRegistered(relationshipName) {
					unregistered[relationshipName] = true
				}
			}
		}
	}

	if len(unregistered) > 0 {
		unregisteredList := make([]string, 0, len(unregistered))
		for name := range unregistered {
			unregisteredList = append(unregisteredList, name)
		}

		m.logger.Error("Found unregistered relationships in use",
			zap.Strings("unregistered", unregisteredList),
		)

		return fmt.Errorf("unregistered relationships in use: %v", unregisteredList)
	}

	m.logger.Debug("All relationships are registered")
	return nil
}
