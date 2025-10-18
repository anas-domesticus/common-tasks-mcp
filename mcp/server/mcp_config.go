package server

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// MCPConfig represents the configuration loaded from mcp.yaml
type MCPConfig struct {
	Server ServerMetadata `yaml:"server"`
	Naming NamingConfig   `yaml:"naming"`
}

// ServerMetadata contains the MCP server identification and description
type ServerMetadata struct {
	Name         string `yaml:"name"`
	DisplayName  string `yaml:"display_name"`
	Instructions string `yaml:"instructions"`
}

// NamingConfig contains friendly names for graph entities
type NamingConfig struct {
	Node NodeNaming `yaml:"node"`
}

// NodeNaming contains singular and plural forms for nodes
type NodeNaming struct {
	Singular        string `yaml:"singular"`
	Plural          string `yaml:"plural"`
	DisplaySingular string `yaml:"display_singular"`
	DisplayPlural   string `yaml:"display_plural"`
}

// DefaultMCPConfig returns the default configuration
func DefaultMCPConfig() MCPConfig {
	return MCPConfig{
		Server: ServerMetadata{
			Name:         "common-tasks-mcp",
			DisplayName:  "Common Tasks",
			Instructions: "This server provides access to commonly performed development tasks and workflows.",
		},
		Naming: NamingConfig{
			Node: NodeNaming{
				Singular:        "task",
				Plural:          "tasks",
				DisplaySingular: "Task",
				DisplayPlural:   "Tasks",
			},
		},
	}
}

// LoadMCPConfig loads the MCP configuration from mcp.yaml in the specified directory
// If the file doesn't exist, it returns the default configuration
func LoadMCPConfig(directory string) (MCPConfig, error) {
	configPath := filepath.Join(directory, "mcp.yaml")

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return DefaultMCPConfig(), nil
	}

	// Read the file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return MCPConfig{}, fmt.Errorf("failed to read mcp.yaml: %w", err)
	}

	// Parse YAML
	var config MCPConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return MCPConfig{}, fmt.Errorf("failed to parse mcp.yaml: %w", err)
	}

	// Validate required fields and set defaults if needed
	if config.Server.Name == "" {
		config.Server.Name = "common-tasks-mcp"
	}
	if config.Server.DisplayName == "" {
		config.Server.DisplayName = "Common Tasks"
	}
	if config.Server.Instructions == "" {
		config.Server.Instructions = "This server provides access to commonly performed development tasks and workflows."
	}

	// Set naming defaults if not provided
	if config.Naming.Node.Singular == "" {
		config.Naming.Node.Singular = "task"
	}
	if config.Naming.Node.Plural == "" {
		config.Naming.Node.Plural = "tasks"
	}
	if config.Naming.Node.DisplaySingular == "" {
		config.Naming.Node.DisplaySingular = "Task"
	}
	if config.Naming.Node.DisplayPlural == "" {
		config.Naming.Node.DisplayPlural = "Tasks"
	}

	return config, nil
}
