# Generic DAG MCP Server

A configuration-driven Model Context Protocol (MCP) server for managing directed acyclic graphs (DAGs) with flexible relationship types. Built in Go for performance and type safety.

## Overview

This is a **generic, reusable DAG framework** that exposes graph operations through MCP. The server's domain, terminology, and relationship types are entirely configuration-driven—allowing you to model tasks, recipes, services, infrastructure, learning paths, or any other domain with directed relationships.

**Key architectural principles:**
- **Zero-code domain adaptation**: Change from "tasks" to "recipes" to "services" by editing YAML files
- **Multiple independent DAGs**: Each relationship type forms its own DAG sharing the same nodes
- **Configuration-driven tools**: MCP tools and descriptions auto-generate based on your domain config
- **Type-safe at runtime**: Strongly typed despite configuration flexibility
- **Automatic cycle detection**: Ensures all relationship DAGs remain valid

The system breaks down knowledge silos by providing AI assistants with structured, queryable knowledge graphs. Whether it's institutional knowledge about development workflows, service dependencies, recipe relationships, or learning paths—the same framework adapts to your domain.

> **⚠️ Active Development Notice**
>
> This project is under active development. The API, configuration schema, and implementation details are subject to change as the design evolves. While the core concepts (DAG management, configuration-driven domain modeling) are stable, specific interfaces and file formats may change in future releases.

## Motivation

AI assistants and agentic workflows increasingly need domain-specific knowledge to function effectively. However, traditional approaches—dumping large documentation files into prompts or expecting AI to parse scattered tribal knowledge—are inefficient and hit context window limits quickly. Modern AI systems work best with **small, targeted pieces of information** delivered precisely when needed.

This is where structured knowledge graphs excel. Instead of providing an AI with a 50-page deployment runbook, you can let it query: "What are the prerequisites for deploying to production?" The system returns only the relevant nodes and their relationships—exactly the context needed, nothing more. This **"less is more"** principle is fundamental to effective context engineering.

The DAG structure provides semantic relationships that AI systems can navigate intelligently. When an AI assistant asks about a task, it doesn't just get the task description—it learns what must happen first (prerequisites), what must follow (required downstream), and what's recommended (suggested downstream). This allows agentic workflows to plan multi-step operations, understand dependencies, and make informed decisions without overloading their context windows.

**Common applications:**
- **AI code assistants** understanding repository workflows and build processes
- **Agentic systems** planning multi-step operations with dependency awareness
- **Automation tools** that need to know "what comes next" in a workflow
- **Knowledge management** systems that surface institutional knowledge on-demand
- **Decision support** systems navigating complex operational procedures

By encoding domain knowledge as queryable graphs rather than static documents, this server enables AI systems to be more precise, efficient, and context-aware.

## Dogfooding

This project uses itself to manage development workflows. The `.tasks/` directory contains our own task repository tracking common development operations like running tests, updating documentation, and reviewing examples. This serves both as a practical development tool and a real-world demonstration of the system.

See [.tasks/README.md](.tasks/README.md) for setup instructions and usage examples.

## Features

- **Configuration-driven identity**: Server name, terminology, and tool names adapt to your domain
- **Flexible relationship system**: Define unlimited relationship types with temporal directionality
- **Multiple independent DAGs**: Each relationship type validated independently for cycles
- **Tag-based indexing**: Fast lookups and filtering by arbitrary tags
- **Dual storage model**: Simple string IDs on disk, resolved pointers at runtime
- **YAML persistence**: Human-readable and git-friendly storage format
- **MCP integration**: Works with Claude Desktop, Claude Code, and any MCP client
- **Both transports**: Stdio (for desktop clients) and HTTP (for web services)
- **Safe mutations**: Clone-validate-commit pattern prevents invalid graph states

## Example Use Cases

The server adapts to any domain through simple YAML configuration. Here's a quick example:

### Development Tasks & Workflows
```yaml
# mcp.yaml
server:
  name: common-tasks-mcp
  display_name: Common Tasks
naming:
  node:
    singular: task
    plural: tasks

# relationships.yaml
relationships:
  - name: prerequisites
    description: Tasks that must be completed before this task
    direction: backward
  - name: downstream_required
    description: Tasks that must be completed after this task
    direction: forward
  - name: downstream_suggested
    description: Recommended follow-up tasks
    direction: forward
```

**Result**: Tools like `add_task`, `list_tasks`, `get_task` for managing development workflows.

### More Examples

For complete, working examples including:
- **Recipe Knowledge Base** - Culinary recipes with ingredient dependencies
- **Microservices Dependency Tracking** - Service mesh relationships
- **Learning Path System** - Educational content with prerequisites
- **Infrastructure Components** - Infrastructure as code dependencies
- **And more...**

See the **[docs/examples/](docs/examples/)** directory for full configuration files, sample data, and usage instructions.

## Installation

### From Source

```bash
go install ./cli/mcp
```

### Using Docker

```bash
docker compose up -d
```

The Docker setup will:
- Store data in `./.tasks` directory (mounted as a volume)
- Run in HTTP mode on port 8080
- Enable verbose logging

## Configuration

### Domain Configuration (`mcp.yaml`)

Place this file in your data directory to configure the server's identity and terminology:

```yaml
# Server metadata
server:
  # Name of the MCP server (used in protocol identification)
  name: my-graph-mcp

  # Human-friendly display name
  display_name: My Knowledge Graph

  # Instructions shown to clients about how to use this server
  instructions: |-
    This server provides access to [your domain description here].
    Use list_[plural] to browse, get_[singular] to retrieve details,
    and add_[singular] to create new entries.

# Friendly names for graph entities
naming:
  node:
    singular: item        # API uses "add_item", "get_item"
    plural: items         # API uses "list_items"
    display_singular: Item
    display_plural: Items
```

### Relationship Configuration (`relationships.yaml`)

Define your relationship types in the same directory:

```yaml
relationships:
  # Example: Things that come before
  - name: prerequisites
    description: Items that must be completed before this item
    direction: backward

  # Example: Things that come after
  - name: next_steps
    description: Items that should follow this item
    direction: forward

  # Example: Related without temporal ordering
  - name: related_to
    description: Conceptually related items
    direction: none
```

**Relationship directions:**
- `backward`: Points to nodes that come **before** in execution/dependency order
- `forward`: Points to nodes that come **after** in execution/dependency order
- `none`: No temporal ordering implied (conceptual links)

### Runtime Configuration

Configuration can be provided via YAML file or environment variables:

**Config file (config.yaml):**
```yaml
transport: stdio
httpPort: 8080
directory: ./data
verbose: false
readOnly: false
```

**Environment Variables:**
- `MCP_TRANSPORT`: Transport mode (stdio or http)
- `MCP_HTTP_PORT`: HTTP port number
- `MCP_DIRECTORY`: Data directory path
- `MCP_VERBOSE`: Enable verbose logging (true/false)
- `MCP_READ_ONLY`: Enable read-only mode (true/false)

## Usage

### Starting the Server

**Stdio mode (default for MCP clients):**
```bash
mcp serve --directory ./data
```

**HTTP mode:**
```bash
mcp serve --transport http --port 8080 --directory ./data
```

### Command-line Options

- `--directory, -d`: Directory where data is stored (default: ".")
- `--transport, -t`: Transport mode: stdio or http (default: "stdio")
- `--port, -p`: HTTP port when using http transport (default: 8080)
- `--verbose, -v`: Enable verbose logging
- `--read-only, -r`: Enable read-only mode (suppresses write tools)
- `--config, -c`: Path to YAML config file

### MCP Tools (Auto-Generated)

The server dynamically generates tools based on your `mcp.yaml` configuration:

- **list_[plural]**: List all nodes or filter by tags
- **get_[singular]**: Get a specific node by ID with full relationship details
- **list_tags**: Get all unique tags with usage counts
- **add_[singular]**: Create a new node with relationships
- **update_[singular]**: Update an existing node
- **delete_[singular]**: Delete a node and clean up all references

**Example**: If you configure `singular: recipe`, the tools become `add_recipe`, `get_recipe`, `list_recipes`, etc.

### MCP Prompts

The server may include domain-specific prompts (check the `prompts/` directory in your data directory):

- **generate-initial-[plural]**: Prompt for generating initial graph content
- **capture-workflow**: Prompt for capturing workflows during active use

When prompts are available, the server also registers:
- **list_prompts**: Get all available prompts with descriptions
- **get_prompt**: Retrieve the full content of a specific prompt

### Node Structure

Nodes are stored as YAML files with this structure (relationship names adapt to your config):

```yaml
id: example-node
name: Example Node
summary: Brief description
description: |
  Detailed description explaining this node.
tags:
  - category-a
  - category-b
edges:
  prerequisites:
    - prerequisite-node-1
    - prerequisite-node-2
  next_steps:
    - next-node-1
  related_to:
    - related-node-1
created_at: 2024-01-15T10:30:00Z
updated_at: 2024-01-15T10:30:00Z
```

**Note**: The `edges` key contains all relationship types. Only the IDs are persisted—pointers are resolved at runtime.

## Contributing

Contributions welcome! This is a generic framework—if you build an interesting domain configuration, consider sharing it as an example.

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, architecture details, and guidelines.

## License

MIT
