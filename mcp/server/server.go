package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"common-tasks-mcp/pkg/graph_manager"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// PromptInfo holds prompt content and metadata
type PromptInfo struct {
	Content     string
	Description string
}

// PromptFrontmatter represents the YAML frontmatter in prompt files
type PromptFrontmatter struct {
	Description string `yaml:"description"`
}

// Server wraps the MCP server
type Server struct {
	mcp         *mcp.Server
	config      Config
	taskManager *graph_manager.Manager
	logger      *zap.Logger
	prompts     map[string]*PromptInfo // Map of prompt name to prompt info
}

// New creates a new MCP server instance with a node manager
func New(cfg Config, logger *zap.Logger) (*Server, error) {
	logger.Info("Creating MCP server", zap.String("directory", cfg.Directory))

	// Load MCP configuration from mcp.yaml
	mcpConfig, err := LoadMCPConfig(cfg.Directory)
	if err != nil {
		logger.Warn("Could not load MCP configuration, using defaults",
			zap.Error(err),
		)
		mcpConfig = DefaultMCPConfig()
	} else {
		logger.Info("MCP configuration loaded",
			zap.String("server_name", mcpConfig.Server.Name),
			zap.String("display_name", mcpConfig.Server.DisplayName),
		)
	}

	// Store MCP config in server config
	cfg.MCP = mcpConfig

	// Create node manager
	taskMgr := graph_manager.NewManager(logger)

	// Load relationships configuration if it exists
	relationshipsPath := filepath.Join(cfg.Directory, "relationships.yaml")
	logger.Info("Loading relationship definitions", zap.String("path", relationshipsPath))
	if err := taskMgr.LoadRelationshipsFromFile(relationshipsPath); err != nil {
		logger.Warn("Could not load relationships configuration",
			zap.String("path", relationshipsPath),
			zap.Error(err),
		)
		// Log the error but continue - relationships file is optional
		if cfg.Verbose {
			fmt.Printf("Warning: Could not load relationships from %s: %v\n", relationshipsPath, err)
		}
	} else {
		// TODO: Add method to get relationship count from manager
		logger.Info("Relationships loaded successfully", zap.String("path", relationshipsPath))
	}

	// Load nodes from directory if any exist
	nodesPath := filepath.Join(cfg.Directory, "nodes")
	logger.Info("Loading nodes from directory", zap.String("path", nodesPath))
	if err := taskMgr.LoadNodesFromDir(nodesPath); err != nil {
		logger.Warn("Could not load nodes from directory",
			zap.String("directory", nodesPath),
			zap.Error(err),
		)
		// Log the error but continue if directory doesn't exist or is empty
		if cfg.Verbose {
			fmt.Printf("Warning: Could not load nodes from %s: %v\n", nodesPath, err)
		}
	} else {
		nodeCount := len(taskMgr.ListAllNodes())
		logger.Info("Nodes loaded successfully", zap.Int("count", nodeCount))
	}

	// Create MCP server using configuration from mcp.yaml
	logger.Debug("Initializing MCP server instance")
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    mcpConfig.Server.Name,
		Version: "0.1.0",
	}, &mcp.ServerOptions{
		Instructions: mcpConfig.Server.Instructions,
	})

	srv := &Server{
		mcp:         mcpServer,
		config:      cfg,
		taskManager: taskMgr,
		logger:      logger,
		prompts:     make(map[string]*PromptInfo),
	}

	// Load prompts from disk
	promptsPath := filepath.Join(cfg.Directory, "prompts")
	logger.Info("Loading prompts from directory", zap.String("path", promptsPath))
	if err := srv.loadPrompts(promptsPath); err != nil {
		logger.Warn("Could not load prompts from directory",
			zap.String("directory", promptsPath),
			zap.Error(err),
		)
		// Continue without prompts - they're optional
		if cfg.Verbose {
			fmt.Printf("Warning: Could not load prompts from %s: %v\n", promptsPath, err)
		}
	} else {
		logger.Info("Prompts loaded successfully", zap.Int("count", len(srv.prompts)))
	}

	// Register all MCP tools
	logger.Debug("Registering MCP tools")
	srv.registerTools()
	if cfg.ReadOnly {
		logger.Info("MCP tools registered successfully (read-only mode: write tools suppressed)")
	} else {
		logger.Info("MCP tools registered successfully")
	}

	// Register all MCP prompts
	logger.Debug("Registering MCP prompts")
	srv.registerPrompts()
	logger.Info("MCP prompts registered successfully")

	return srv, nil
}

// RunHTTP starts the MCP server with HTTP transport
func (s *Server) RunHTTP(ctx context.Context) error {
	s.logger.Info("Initializing HTTP transport", zap.Int("port", s.config.HTTPPort))

	// Create streamable HTTP handler
	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		s.logger.Debug("Incoming HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
		)
		return s.mcp
	}, nil)

	// Create HTTP server
	addr := fmt.Sprintf(":%d", s.config.HTTPPort)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Run server with graceful shutdown
	errChan := make(chan error, 1)
	go func() {
		s.logger.Info("HTTP server listening", zap.String("address", addr))
		fmt.Printf("Starting MCP server on http://localhost%s\n", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("HTTP server error", zap.Error(err))
			errChan <- err
		}
	}()

	// Wait for context cancellation or error
	select {
	case <-ctx.Done():
		s.logger.Info("Context cancelled, shutting down HTTP server")
		shutdownCtx := context.Background()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("Error during HTTP server shutdown", zap.Error(err))
			return err
		}
		s.logger.Info("HTTP server shutdown complete")
		return nil
	case err := <-errChan:
		return err
	}
}

// Run starts the MCP server with stdio transport
func (s *Server) Run(ctx context.Context) error {
	s.logger.Info("Starting MCP server with stdio transport")
	err := s.mcp.Run(ctx, &mcp.StdioTransport{})
	if err != nil {
		s.logger.Error("Stdio server error", zap.Error(err))
	} else {
		s.logger.Info("Stdio server exited normally")
	}
	return err
}

// loadPrompts loads prompt files from the specified directory
func (s *Server) loadPrompts(promptsDir string) error {
	// Check if prompts directory exists
	if _, err := os.Stat(promptsDir); os.IsNotExist(err) {
		return fmt.Errorf("prompts directory does not exist: %s", promptsDir)
	}

	// Read all .md files from the prompts directory
	entries, err := os.ReadDir(promptsDir)
	if err != nil {
		return fmt.Errorf("failed to read prompts directory: %w", err)
	}

	loadedCount := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Only process .md files
		if filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		// Read prompt content
		promptPath := filepath.Join(promptsDir, entry.Name())
		content, err := os.ReadFile(promptPath)
		if err != nil {
			s.logger.Warn("Failed to read prompt file",
				zap.String("file", entry.Name()),
				zap.Error(err),
			)
			continue
		}

		// Use filename without extension as prompt name
		promptName := entry.Name()[:len(entry.Name())-3]

		// Parse frontmatter and content
		promptInfo := s.parsePromptFile(string(content), promptName)

		s.prompts[promptName] = promptInfo
		loadedCount++

		s.logger.Debug("Loaded prompt",
			zap.String("name", promptName),
			zap.String("file", entry.Name()),
			zap.String("description", promptInfo.Description),
			zap.Int("content_size", len(promptInfo.Content)),
		)
	}

	if loadedCount == 0 {
		return fmt.Errorf("no prompt files found in directory")
	}

	return nil
}

// parsePromptFile parses a prompt file with optional YAML frontmatter
func (s *Server) parsePromptFile(content string, promptName string) *PromptInfo {
	info := &PromptInfo{
		Content:     content,
		Description: fmt.Sprintf("Prompt: %s", promptName), // Default fallback
	}

	// Check for YAML frontmatter (starts with ---)
	if !strings.HasPrefix(content, "---\n") {
		// No frontmatter, return as-is
		return info
	}

	// Find the closing ---
	parts := strings.SplitN(content[4:], "\n---\n", 2)
	if len(parts) != 2 {
		// Invalid frontmatter, return as-is
		s.logger.Warn("Invalid frontmatter format in prompt",
			zap.String("name", promptName),
		)
		return info
	}

	// Parse YAML frontmatter
	var frontmatter PromptFrontmatter
	if err := yaml.Unmarshal([]byte(parts[0]), &frontmatter); err != nil {
		s.logger.Warn("Failed to parse frontmatter in prompt",
			zap.String("name", promptName),
			zap.Error(err),
		)
		return info
	}

	// Update info with parsed values
	if frontmatter.Description != "" {
		info.Description = frontmatter.Description
	}
	info.Content = strings.TrimSpace(parts[1])

	return info
}

// registerPrompts registers all MCP prompts with the server
func (s *Server) registerPrompts() {
	// Register all prompts that were successfully loaded
	for promptName, promptInfo := range s.prompts {
		prompt := &mcp.Prompt{
			Name:        promptName,
			Description: promptInfo.Description,
		}

		s.mcp.AddPrompt(prompt, s.handlePrompt)
		s.logger.Debug("Registered prompt",
			zap.String("name", promptName),
			zap.String("description", promptInfo.Description),
		)
	}

	if len(s.prompts) == 0 {
		s.logger.Info("No prompts registered (none found in prompts directory)")
	}
}

// handlePrompt is a generic handler for all prompts
func (s *Server) handlePrompt(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	promptName := req.Params.Name
	s.logger.Debug("Handling prompt request", zap.String("name", promptName))

	// Get prompt info from loaded prompts
	promptInfo, exists := s.prompts[promptName]
	if !exists {
		s.logger.Error("Prompt not found", zap.String("name", promptName))
		return nil, fmt.Errorf("prompt %s not found", promptName)
	}

	s.logger.Info("Successfully retrieved prompt",
		zap.String("name", promptName),
		zap.Int("content_length", len(promptInfo.Content)),
	)

	return &mcp.GetPromptResult{
		Description: promptInfo.Description,
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: promptInfo.Content,
				},
			},
		},
	}, nil
}
