package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server wraps the MCP server
type Server struct {
	mcp    *mcp.Server
	config Config
}

// New creates a new MCP server instance
func New(cfg Config) *Server {
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "common-tasks-mcp",
		Version: "0.1.0",
	}, nil)

	return &Server{
		mcp:    mcpServer,
		config: cfg,
	}
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
