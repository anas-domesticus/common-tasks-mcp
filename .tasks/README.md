# Common Tasks Repository

This directory contains the task repository for the **common-tasks-mcp** project itself. We're dogfooding our own software by using it to track common development workflows for this codebase.

## What is Dogfooding?

This project is a generic DAG-based task management system. Rather than maintaining development workflows in scattered documentation, we use our own MCP server to manage tasks like:
- Running tests
- Updating documentation
- Reviewing examples
- Adding new prompts
- And other common repository operations

This serves as both a practical tool for development and a real-world example of the system in action.

## Running the Server

### Using Docker Compose

From the project root directory:

```bash
docker compose up -d
```

This will:
- Start the MCP server in HTTP mode on port 8080
- Mount this `.tasks/` directory as the data directory
- Enable verbose logging for debugging

### Using the CLI (from source)

```bash
# Install the CLI
go install ./cli/mcp

# Run the server
mcp serve --directory .tasks --transport stdio
```

## Adding to Claude Code

To use this task repository in Claude Code, add the server to your MCP configuration:

```bash
# For Docker setup (HTTP transport)
claude mcp add http://localhost:8080 common-tasks

# For CLI setup (stdio transport)
claude mcp add "mcp serve --directory /path/to/common-tasks-mcp/.tasks" common-tasks
```

Replace `/path/to/common-tasks-mcp/.tasks` with the absolute path to this directory.

## Usage

Once configured, you can use these tools in Claude Code:

- `list_tasks` - Browse available development tasks
- `get_task` - Get detailed workflow for a specific task
- `list_tags` - See all available tags
- `add_task` - Create new development tasks
- `update_task` - Modify existing tasks
- `delete_task` - Remove tasks

## Configuration Files

- `mcp.yaml` - Server identity and terminology configuration
- `relationships.yaml` - Task relationship definitions (prerequisites, downstream_required, downstream_suggested)
- `nodes/` - Individual task definitions in YAML format
- `prompts/` - AI assistant prompts for task generation and workflow capture

## More Information

See the [main README](../README.md) for complete documentation on the Generic DAG MCP Server architecture, configuration options, and examples.
