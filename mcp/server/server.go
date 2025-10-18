package server

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"path/filepath"

	"common-tasks-mcp/pkg/graph_manager"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

//go:embed prompts/generate-initial-tasks.md
var generateInitialTasksPrompt string

//go:embed prompts/capture-workflow.md
var captureWorkflowPrompt string

// Server wraps the MCP server
type Server struct {
	mcp         *mcp.Server
	config      Config
	taskManager *graph_manager.Manager
	logger      *zap.Logger
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
	}

	// Register all MCP tools
	logger.Debug("Registering MCP tools")
	srv.registerTools()
	logger.Info("MCP tools registered successfully")

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

// GetGenerateInitialTasksPrompt returns the embedded prompt for generating initial tasks
func GetGenerateInitialTasksPrompt() string {
	return generateInitialTasksPrompt
}

// registerPrompts registers all MCP prompts with the server
func (s *Server) registerPrompts() {
	// Generate initial tasks prompt
	generateTasksPrompt := &mcp.Prompt{
		Name:        "generate-initial-tasks",
		Description: "Prompt for generating an initial set of tasks for a codebase. Guides exploration of project structure, build systems, CI/CD configs, and documentation to create tasks with proper relationships and workflows.",
	}

	s.mcp.AddPrompt(generateTasksPrompt, s.handleGenerateTasksPrompt)

	// Capture workflow prompt
	captureWorkflowPrompt := &mcp.Prompt{
		Name:        "capture-workflow",
		Description: "Prompt for capturing workflows as tasks during active development. Guides recognition of repeatable operations and helps maintain tasks as you work, ensuring institutional knowledge is captured in real-time.",
	}

	s.mcp.AddPrompt(captureWorkflowPrompt, s.handleCaptureWorkflowPrompt)
}

// handleGenerateTasksPrompt handles the generate-initial-tasks prompt
func (s *Server) handleGenerateTasksPrompt(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	s.logger.Debug("Handling generate-initial-tasks prompt request")

	prompt := generateInitialTasksPrompt

	s.logger.Info("Successfully retrieved generate-initial-tasks prompt", zap.Int("length", len(prompt)))

	return &mcp.GetPromptResult{
		Description: "Prompt for generating an initial set of tasks for a codebase",
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: prompt,
				},
			},
		},
	}, nil
}

// handleCaptureWorkflowPrompt handles the capture-workflow prompt
func (s *Server) handleCaptureWorkflowPrompt(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	s.logger.Debug("Handling capture-workflow prompt request")

	prompt := captureWorkflowPrompt

	s.logger.Info("Successfully retrieved capture-workflow prompt", zap.Int("length", len(prompt)))

	return &mcp.GetPromptResult{
		Description: "Prompt for capturing workflows as tasks during active development",
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: prompt,
				},
			},
		},
	}, nil
}
