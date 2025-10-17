package server

import "fmt"

// Config holds the configuration for the MCP server
type Config struct {
	// Transport mode: "stdio" or "http"
	Transport string `env:"MCP_TRANSPORT" yaml:"transport" default:"stdio"`

	// HTTP port (only used when Transport is "http")
	HTTPPort int `env:"MCP_HTTP_PORT" yaml:"httpPort" default:"8080"`

	// Directory where tasks are stored (git repository)
	Directory string `env:"MCP_DIRECTORY" yaml:"directory" default:"."`

	// Verbose logging
	Verbose bool `env:"MCP_VERBOSE" yaml:"verbose" default:"false"`
}

// Validate implements the config.Validator interface
func (c Config) Validate() error {
	if c.Transport != "stdio" && c.Transport != "http" {
		return fmt.Errorf("transport must be 'stdio' or 'http', got '%s'", c.Transport)
	}

	if c.Transport == "http" {
		if c.HTTPPort < 1 || c.HTTPPort > 65535 {
			return fmt.Errorf("httpPort must be between 1-65535, got %d", c.HTTPPort)
		}
	}

	return nil
}
