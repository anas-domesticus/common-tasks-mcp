package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"common-tasks-mcp/mcp/server"
	"common-tasks-mcp/pkg/config"
	"common-tasks-mcp/pkg/logger"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
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

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Output the prompt for generating initial tasks",
	Long: `Output the prompt that can be used with Claude Code or other AI assistants
to generate an initial set of tasks for a codebase.

This prompt guides the AI to:
- Explore the codebase structure
- Find build systems, CI/CD configs, and documentation
- Create tasks with proper relationships and workflows
- Use actual commands and paths from the project`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(server.GetGenerateInitialTasksPrompt())
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
		// Create a basic logger first for config loading
		// We'll recreate it with proper verbosity later
		basicLog, err := logger.New(false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}
		defer basicLog.Sync()

		// Set logger for config package
		config.SetLogger(basicLog)

		// Load configuration
		var cfg ServerConfig
		if err := config.GetConfig(&cfg, configPath, true); err != nil {
			basicLog.Error("Failed to load configuration", zap.Error(err))
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
			basicLog.Error("Configuration validation failed", zap.Error(err))
			fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
			os.Exit(1)
		}

		// Initialize logger with proper verbosity
		log, err := logger.New(cfg.Verbose)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}
		defer log.Sync()

		// Update logger for config package
		config.SetLogger(log)

		log.Info("Starting Common Tasks MCP Server",
			zap.String("version", "0.1.0"),
			zap.String("transport", cfg.Transport),
			zap.String("task_directory", cfg.Directory),
			zap.Bool("verbose", cfg.Verbose),
		)

		if cfg.Transport == "http" {
			log.Info("HTTP transport configured", zap.Int("port", cfg.HTTPPort))
		}

		// Print server info (keep for user visibility)
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
		srv, err := server.New(cfg, log)
		if err != nil {
			log.Error("Failed to create server", zap.Error(err))
			fmt.Fprintf(os.Stderr, "Error creating server: %v\n", err)
			os.Exit(1)
		}

		log.Info("Server created successfully")

		// Setup context with signal handling
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Handle shutdown signals
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		log.Debug("Signal handlers registered")

		// Start server in goroutine
		errChan := make(chan error, 1)
		go func() {
			var err error
			if cfg.Transport == "http" {
				log.Info("Starting HTTP server")
				err = srv.RunHTTP(ctx)
			} else {
				log.Info("Starting stdio server")
				err = srv.Run(ctx)
			}
			errChan <- err
		}()

		// Wait for shutdown signal or error
		select {
		case sig := <-sigChan:
			log.Info("Received shutdown signal", zap.String("signal", sig.String()))
			fmt.Println("\nShutting down server...")
			cancel()
		case err := <-errChan:
			if err != nil {
				log.Error("Server exited with error", zap.Error(err))
				fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
				os.Exit(1)
			}
			log.Info("Server exited normally")
		}

		log.Info("Shutdown complete")
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
	rootCmd.AddCommand(promptCmd)
}
