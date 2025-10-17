package server

import (
	"context"
	"fmt"
	"net/http"

	"common-tasks-mcp/pkg/task_manager"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server wraps the MCP server
type Server struct {
	mcp         *mcp.Server
	config      Config
	taskManager *task_manager.Manager
}

// New creates a new MCP server instance with a task manager
func New(cfg Config) (*Server, error) {
	// Create task manager
	taskMgr := task_manager.NewManager()

	// Load tasks from directory if any exist
	if err := taskMgr.LoadFromDir(cfg.Directory); err != nil {
		// Log the error but continue if directory doesn't exist or is empty
		if cfg.Verbose {
			fmt.Printf("Warning: Could not load tasks from %s: %v\n", cfg.Directory, err)
		}
	}

	// Create MCP server
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "common-tasks-mcp",
		Version: "0.1.0",
	}, nil)

	return &Server{
		mcp:         mcpServer,
		config:      cfg,
		taskManager: taskMgr,
	}, nil
}

// RunHTTP starts the MCP server with HTTP transport
func (s *Server) RunHTTP(ctx context.Context) error {
	// Create streamable HTTP handler
	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
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
		fmt.Printf("Starting MCP server on http://localhost%s\n", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for context cancellation or error
	select {
	case <-ctx.Done():
		return httpServer.Shutdown(context.Background())
	case err := <-errChan:
		return err
	}
}

// Run starts the MCP server with stdio transport
func (s *Server) Run(ctx context.Context) error {
	return s.mcp.Run(ctx, &mcp.StdioTransport{})
}
