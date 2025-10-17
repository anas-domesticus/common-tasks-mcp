package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"common-tasks-mcp/mcp/server"
	"common-tasks-mcp/pkg/config"

	"github.com/spf13/cobra"
)

var (
	// Flags
	configPath string

	// Serve command flags
	transport string
	httpPort  int
	directory string
	verbose   bool
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Common Tasks MCP Server",
	Long: `A Go-based Model Context Protocol (MCP) server for storing and managing
commonly performed tasks in a git repository.

This MCP server enables you to store, retrieve, and manage frequently used tasks
directly from your git repository. It captures institutional knowledge about what
needs to happen for different types of changes and can be integrated with any
client that supports the Model Context Protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Common Tasks MCP Server")
		fmt.Println("Use 'mcp serve' to start the server or 'mcp --help' for more information.")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Common Tasks MCP Server v0.1.0")
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the MCP server",
	Long: `Start the MCP server to handle task management requests via the Model Context Protocol.

The server enables MCP-compatible clients (Claude Desktop, Claude Code, etc.) to
store, retrieve, and manage commonly performed tasks from a git-backed repository.

Transport modes:
  stdio - Standard input/output (default, for MCP clients)
  http  - HTTP server (for REST API access)`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		var cfg ServerConfig
		if err := config.GetConfig(&cfg, configPath, true); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		// Override config with command-line flags if provided
		if cmd.Flags().Changed("transport") {
			cfg.Transport = transport
		}
		if cmd.Flags().Changed("port") {
			cfg.HTTPPort = httpPort
		}
		if cmd.Flags().Changed("directory") {
			cfg.Directory = directory
		}
		if cmd.Flags().Changed("verbose") {
			cfg.Verbose = verbose
		}

		// Validate configuration after flag overrides
		if err := cfg.Validate(); err != nil {
			fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
			os.Exit(1)
		}

		// Print server info
		fmt.Println("Common Tasks MCP Server v0.1.0")
		fmt.Printf("Transport mode: %s\n", cfg.Transport)
		fmt.Printf("Task directory: %s\n", cfg.Directory)
		if cfg.Transport == "http" {
			fmt.Printf("HTTP port: %d\n", cfg.HTTPPort)
		}
		if cfg.Verbose {
			fmt.Println("Verbose logging enabled")
		}
		fmt.Println()

		// Create server
		srv, err := server.New(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating server: %v\n", err)
			os.Exit(1)
		}

		// Setup context with signal handling
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Handle shutdown signals
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		// Start server in goroutine
		errChan := make(chan error, 1)
		go func() {
			var err error
			if cfg.Transport == "http" {
				err = srv.RunHTTP(ctx)
			} else {
				err = srv.Run(ctx)
			}
			errChan <- err
		}()

		// Wait for shutdown signal or error
		select {
		case <-sigChan:
			fmt.Println("\nShutting down server...")
			cancel()
		case err := <-errChan:
			if err != nil {
				fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config file path (YAML)")

	// Serve command flags
	serveCmd.Flags().StringVarP(&transport, "transport", "t", "stdio", "transport mode: stdio or http")
	serveCmd.Flags().IntVarP(&httpPort, "port", "p", 8080, "HTTP port (only used with --transport=http)")
	serveCmd.Flags().StringVarP(&directory, "directory", "d", ".", "directory where tasks are stored (git repository)")
	serveCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")

	// Add subcommands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serveCmd)
}
