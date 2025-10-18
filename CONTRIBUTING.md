# Contributing to Generic DAG MCP Server

Thank you for your interest in contributing! This document provides guidelines and information for developers working on this project.

## Development Setup

### Requirements

- Go 1.25 or higher
- Docker (optional, for containerized deployment)

### Building from Source

```bash
# Build the CLI
go build -o mcp ./cli/mcp

# Install globally
go install ./cli/mcp
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Test specific package
go test ./pkg/graph_manager/...
```

## Project Structure

```
.
├── cli/mcp/                 # CLI binary
├── pkg/
│   ├── graph_manager/       # Generic DAG operations
│   │   └── types/           # Core data structures
│   ├── config/              # Configuration loading
│   └── logger/              # Logging setup
├── mcp/server/              # MCP server implementation
│   └── prompts/             # Embedded prompt templates
└── .tasks/                  # Example data directory
    ├── mcp.yaml             # Domain configuration
    ├── relationships.yaml   # Relationship definitions
    └── nodes/               # Node storage (YAML files)
```

## Architecture

### Core Components

- **`pkg/graph_manager/`**: Generic DAG operations (no domain knowledge)
  - Cycle detection across multiple relationship types
  - Node persistence and pointer resolution
  - Tag-based indexing
  - Safe mutation with clone-validate-commit pattern

- **`pkg/graph_manager/types/`**: Core data structures
  - `Node`: Generic graph vertex with arbitrary edge types
  - `Edge`: Directed connection with relationship type
  - `Relationship`: Metadata about edge categories
  - `RelationshipDirection`: Temporal flow enumeration

- **`mcp/server/`**: Configuration-driven MCP layer
  - Dynamic tool generation from `mcp.yaml`
  - Relationship loading from `relationships.yaml`
  - Transport abstraction (stdio/HTTP)

### Two-Level Edge Storage

The system uses a dual storage model to handle relationships efficiently:

**EdgeIDs** (Persisted):
```yaml
edges:
  prerequisites: ["node-a", "node-b"]
```
- Simple string IDs
- Stored in YAML/JSON
- Source of truth on disk

**Edges** (Runtime):
```go
Edges: map[string][]Edge{
  "prerequisites": {
    {To: *Node, Type: *Relationship},
    {To: *Node, Type: *Relationship},
  },
}
```
- Resolved pointers to full objects
- Computed after loading from disk
- Not persisted (avoids circular references)

### Multiple Independent DAGs

Each relationship type forms its own DAG:
```
Prerequisites DAG: A → B → C
Next Steps DAG:    C → D → E
Related DAG:       A ↔ D (no ordering)
```

All three share the same nodes but have independent edge structures. Each is validated separately for cycles.

### Design Principles

1. **Configuration-driven**: The framework adapts to any domain through YAML configuration
2. **Type safety**: Strong typing at runtime despite configuration flexibility
3. **Immutability**: Clone-validate-commit pattern prevents invalid states
4. **Separation of concerns**: Graph operations are domain-agnostic
5. **Extensibility**: New relationship types and domains without code changes

## Contributing Guidelines

### Code Style

- Follow standard Go conventions and formatting (`gofmt`)
- Write tests for new functionality
- Keep the graph_manager package domain-agnostic
- Document exported functions and types

### Adding Features

When adding new features:

1. **Keep it generic**: The framework should remain domain-agnostic
2. **Configuration-first**: Prefer configuration options over hardcoded behavior
3. **Test thoroughly**: Include unit tests for new functionality
4. **Update docs**: Keep README and CONTRIBUTING.md up to date

### Testing Approach

- Unit tests for graph operations and validation
- Integration tests for MCP server functionality
- Configuration validation tests
- Cycle detection edge cases

### Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Ensure all tests pass
5. Update documentation as needed
6. Submit a pull request

## Domain Configuration Examples

We welcome contributions of interesting domain configurations! If you've built a configuration for a specific use case (recipe management, infrastructure tracking, learning paths, etc.), consider sharing it as an example in the `docs/examples/` directory.

### Example Structure

```
docs/examples/your-domain/
├── README.md              # Description of the use case
├── mcp.yaml              # Domain configuration
├── relationships.yaml    # Relationship definitions
└── nodes/                # Sample nodes (optional)
```

## Questions or Issues?

- Open an issue for bugs or feature requests
- Check existing issues before creating new ones
- Provide clear reproduction steps for bugs

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
