# Common Tasks MCP Server

A Go-based Model Context Protocol (MCP) server for storing and managing commonly performed tasks in a git repository.

## Overview

This MCP server enables you to store, retrieve, and manage frequently used tasks directly from your git repository. It can be integrated with any client that supports the Model Context Protocol.

The goal is to break down knowledge silos by providing AI coding assistants with institutional knowledge about what needs to happen for different types of changes. When an AI tool like Claude Code queries this MCP server, it receives structured instructions on all the tasks, prerequisites, and follow-up actions required for a given change—capturing the tribal knowledge that typically lives only in developers' heads or scattered documentation.

This repository will also include a Claude Code agent that monitors coding sessions and automatically surfaces relevant task workflows based on what developers are working on, proactively suggesting the full context of what needs to happen for a given change.

## Features

- Store commonly performed tasks in a git-backed repository
- Integration with MCP-compatible clients (Claude Desktop, Claude Code, etc.)
- Version control for your tasks using git
- Three-DAG relationship model supporting prerequisites, required downstream tasks, and suggested downstream tasks
- Tag-based task organization and retrieval

## Installation

```bash
go install
```

## Usage

Coming soon.

## Requirements

- Go 1.21 or higher
- Git

## License

MIT
