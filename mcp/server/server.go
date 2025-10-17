package server

import (
	"context"
	"fmt"
	"net/http"

	"common-tasks-mcp/pkg/task_manager"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

// Server wraps the MCP server
type Server struct {
	mcp         *mcp.Server
	config      Config
	taskManager *task_manager.Manager
	logger      *zap.Logger
}

// New creates a new MCP server instance with a task manager
func New(cfg Config, logger *zap.Logger) (*Server, error) {
	logger.Info("Creating MCP server", zap.String("directory", cfg.Directory))

	// Create task manager
	taskMgr := task_manager.NewManager(logger)

	// Load tasks from directory if any exist
	logger.Info("Loading tasks from directory", zap.String("path", cfg.Directory))
	if err := taskMgr.LoadFromDir(cfg.Directory); err != nil {
		logger.Warn("Could not load tasks from directory",
			zap.String("directory", cfg.Directory),
			zap.Error(err),
		)
		// Log the error but continue if directory doesn't exist or is empty
		if cfg.Verbose {
			fmt.Printf("Warning: Could not load tasks from %s: %v\n", cfg.Directory, err)
		}
	} else {
		taskCount := len(taskMgr.ListAllTasks())
		logger.Info("Tasks loaded successfully", zap.Int("count", taskCount))
	}

	// Create MCP server
	logger.Debug("Initializing MCP server instance")
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "common-tasks-mcp",
		Version: "0.1.0",
	}, &mcp.ServerOptions{
		Instructions: "This server provides access to commonly performed development tasks and workflows. Each task includes: what needs to be done first (prerequisites), what must be done after (required follow-ups), and what's recommended to do after (suggested follow-ups). Use this to understand the complete workflow for any development task, not just the immediate action. Start by listing tasks with relevant tags to find what you need, then get the full task details to see the complete workflow.",
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
