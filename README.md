# Common Tasks MCP Server

A Go-based Model Context Protocol (MCP) server for storing and managing commonly performed tasks in a local directory.

## Overview

This MCP server enables you to store, retrieve, and manage frequently used tasks as YAML files in a local directory. It can be integrated with any client that supports the Model Context Protocol.

The goal is to break down knowledge silos by providing AI coding assistants with institutional knowledge about what needs to happen for different types of changes. When an AI tool like Claude Code queries this MCP server, it receives structured instructions on all the tasks, prerequisites, and follow-up actions required for a given change—capturing the tribal knowledge that typically lives only in developers' heads or scattered documentation.

This repository will also include a Claude Code agent that monitors coding sessions and automatically surfaces relevant task workflows based on what developers are working on, proactively suggesting the full context of what needs to happen for a given change.

## Features

- Store commonly performed tasks as YAML files in a local directory
- Integration with MCP-compatible clients (Claude Desktop, Claude Code, etc.)
- Three-DAG relationship model supporting prerequisites, required downstream tasks, and suggested downstream tasks
- Tag-based task organization and retrieval
- Cycle detection to ensure valid DAG relationships
- Both stdio and HTTP transport modes

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
- Store tasks in `./.tasks` directory (mounted as a volume)
- Run in HTTP mode on port 8080
- Enable verbose logging

## Usage

### Starting the Server

**Stdio mode (default for MCP clients):**
```bash
mcp serve --directory ./tasks
```

**HTTP mode:**
```bash
mcp serve --transport http --port 8080 --directory ./tasks
```

### Command-line Options

- `--directory, -d`: Directory where tasks are stored (default: ".")
- `--transport, -t`: Transport mode: stdio or http (default: "stdio")
- `--port, -p`: HTTP port when using http transport (default: 8080)
- `--verbose, -v`: Enable verbose logging
- `--config, -c`: Path to YAML config file

### Configuration

Configuration can be provided via YAML file or environment variables:

**YAML:**
```yaml
transport: stdio
http_port: 8080
directory: ./tasks
verbose: false
```

**Environment Variables:**
- `MCP_TRANSPORT`: Transport mode (stdio or http)
- `MCP_HTTP_PORT`: HTTP port number
- `MCP_DIRECTORY`: Tasks directory path
- `MCP_VERBOSE`: Enable verbose logging (true/false)

### MCP Tools

The server provides the following MCP tools:

- **list_tasks**: List all tasks or filter by tags
- **get_task**: Get a specific task by ID with full relationship details
- **add_task**: Create a new task with relationships
- **update_task**: Update an existing task
- **delete_task**: Delete a task and clean up all references

### Task Structure

Tasks are stored as YAML files with the following structure:

```yaml
id: example-task
name: Example Task
summary: Brief description of the task
description: |
  Detailed description explaining what needs to be done.
tags:
  - backend
  - database
prerequisites:
  - setup-database
downstream_required:
  - run-migrations
downstream_suggested:
  - update-docs
```

## Requirements

- Go 1.21 or higher
- Docker (optional, for containerized deployment)

## License

MIT
