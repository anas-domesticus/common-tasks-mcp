# Graph Manager

A flexible, type-safe Go library for managing directed acyclic graphs (DAGs) with configurable relationship types. Designed for workflow management, task dependencies, and any system requiring complex multi-DAG relationships between nodes.

## Features

- **Multiple Independent DAGs**: Define unlimited relationship types, each forming its own DAG
- **Cycle Detection**: Automatic validation prevents invalid graph structures
- **Safe Mutations**: Clone-validate-commit pattern for transactional changes
- **Persistence**: YAML-based storage with automatic pointer resolution
- **Tag-Based Indexing**: Fast lookups by tags
- **Auto-Cleanup**: Deleting nodes automatically removes all references
- **Type-Safe**: Strongly typed nodes, edges, and relationships

## Core Concepts

### Nodes

Nodes are the vertices in your graph, representing tasks, work items, or any entity you need to model.

```go
node := &types.Node{
    ID:          "deploy-api",
    Name:        "Deploy API Service",
    Summary:     "Deploy REST API to production",
    Description: "Full deployment including tests and rollout",
    Tags:        []string{"deployment", "production", "api"},
    EdgeIDs: map[string][]string{
        "prerequisites":        {"build-binary", "run-tests"},
        "downstream_required":  {"smoke-test", "update-docs"},
        "downstream_suggested": {"add-monitoring"},
    },
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
}
```

### Edges

Edges are directed connections between nodes with a specific relationship type:

```go
edge := types.Edge{
    To:   targetNode,  // Pointer to destination node
    Type: relationship, // Relationship category
}
```

### Relationships

Relationships define the semantic meaning and temporal direction of edges:

```go
prerequisitesRel := types.Relationship{
    Name:        "prerequisites",
    Description: "Tasks that must be completed before this task",
    Direction:   types.DirectionBackward, // Points to nodes that come before
}

downstreamRel := types.Relationship{
    Name:        "downstream_required",
    Description: "Tasks that must be completed after this task",
    Direction:   types.DirectionForward, // Points to nodes that come after
}

relatesRel := types.Relationship{
    Name:        "related_to",
    Description: "Conceptually related tasks",
    Direction:   types.DirectionNone, // No temporal ordering
}
```

### Relationship Directions

- **DirectionBackward**: Points to nodes that come **before** in execution order
  - Example: "prerequisites" (must complete before this node)
  - Example: "depends_on" (this node depends on these)

- **DirectionForward**: Points to nodes that come **after** in execution order
  - Example: "downstream_required" (must complete after this node)
  - Example: "triggers" (this node triggers these)

- **DirectionNone**: No temporal ordering
  - Example: "related_to" (conceptual link)
  - Example: "validates" (quality assurance link)

## Quick Start

### Installation

```bash
go get common-tasks-mcp/pkg/graph_manager
```

### Basic Usage

```go
package main

import (
    "common-tasks-mcp/pkg/graph_manager"
    "common-tasks-mcp/pkg/graph_manager/types"
    "common-tasks-mcp/pkg/logger"
    "time"
)

func main() {
    // Create a manager
    log, _ := logger.New(false)
    manager := graph_manager.NewManager(log)

    // Create nodes
    buildNode := &types.Node{
        ID:        "build-binary",
        Name:      "Build Binary",
        Summary:   "Compile the application",
        Tags:      []string{"build"},
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    testNode := &types.Node{
        ID:      "run-tests",
        Name:    "Run Tests",
        Summary: "Execute test suite",
        Tags:    []string{"testing"},
        EdgeIDs: map[string][]string{
            "prerequisites": {"build-binary"},
        },
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    deployNode := &types.Node{
        ID:      "deploy",
        Name:    "Deploy to Production",
        Summary: "Deploy the application",
        Tags:    []string{"deployment"},
        EdgeIDs: map[string][]string{
            "prerequisites": {"build-binary", "run-tests"},
        },
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    // Add nodes to manager (validates for cycles)
    manager.AddNode(buildNode)
    manager.AddNode(testNode)
    manager.AddNode(deployNode)

    // Save to disk
    manager.PersistToDir("./tasks")

    // Load from disk (resolves pointers)
    manager.LoadFromDir("./tasks")

    // Query nodes
    node, _ := manager.GetNode("deploy")
    nodes, _ := manager.GetNodesByTag("testing")
    allNodes := manager.ListAllNodes()
}
```

## API Reference

### Manager

#### Creating a Manager

```go
func NewManager(logger *zap.Logger) *Manager
```

Creates a new graph manager instance.

#### Adding Nodes

```go
func (m *Manager) AddNode(node *types.Node) error
```

Adds a node to the graph. Uses clone-validate-commit pattern to prevent cycles.

**Returns error if:**
- Node is nil or has empty ID
- Node with same ID already exists
- Addition would create a cycle

#### Updating Nodes

```go
func (m *Manager) UpdateNode(node *types.Node) error
```

Updates an existing node and refreshes all pointers.

**Returns error if:**
- Node is nil or has empty ID
- Node doesn't exist
- Update would create a cycle

#### Deleting Nodes

```go
func (m *Manager) DeleteNode(id string) error
```

Removes a node and cleans up all references to it from other nodes.

**Returns error if:**
- ID is empty
- Node doesn't exist

#### Querying Nodes

```go
func (m *Manager) GetNode(id string) (*types.Node, error)
func (m *Manager) GetNodesByTag(tag string) ([]*types.Node, error)
func (m *Manager) ListAllNodes() []*types.Node
```

Retrieve nodes by ID, tag, or all nodes.

#### Persistence

```go
func (m *Manager) PersistToDir(dirPath string) error
func (m *Manager) LoadFromDir(dirPath string) error
```

**PersistToDir**: Writes all nodes as YAML files to the specified directory.

**LoadFromDir**: Reads all YAML files from the directory, validates for cycles, and resolves node pointers.

#### Graph Operations

```go
func (m *Manager) DetectCycles() error
func (m *Manager) ResolveNodePointers() error
func (m *Manager) Clone() *Manager
```

**DetectCycles**: Checks all relationship types for cycles. Returns error with detailed cycle information if found.

**ResolveNodePointers**: Populates the Edges map for all nodes by looking up EdgeIDs. Called automatically by LoadFromDir.

**Clone**: Creates a deep copy of the manager for transactional testing.

### Node Methods

#### Edge ID Operations (String-based)

```go
func (n *Node) GetEdgeIDs(relationshipName string) []string
func (n *Node) SetEdgeIDs(relationshipName string, ids []string) error
func (n *Node) AddEdgeID(relationshipName string, id string) error
```

Work with edge IDs as strings (persisted format).

#### Edge Operations (Pointer-based)

```go
func (n *Node) GetEdges(relationshipName string) []Edge
func (n *Node) SetEdges(relationshipName string, edges []Edge) error
func (n *Node) AddEdge(relationshipName string, edge Edge) error
```

Work with resolved edge objects (runtime format). Automatically keeps EdgeIDs in sync.

#### Node Operations

```go
func (n *Node) Clone() *Node
func (n *Node) Equals(other *Node) bool
```

**Clone**: Creates a deep copy of the node (excludes Edges map).

**Equals**: Compares all persisted fields for equality.

## Advanced Usage

### Working with Multiple DAGs

Each relationship type forms an independent DAG sharing the same nodes:

```go
// Create a node with multiple relationship types
node := &types.Node{
    ID:   "task-a",
    Name: "Task A",
    EdgeIDs: map[string][]string{
        "prerequisites":        {"setup-env"},
        "downstream_required":  {"run-tests"},
        "validates":            {"code-review"},
        "documents":            {"update-readme"},
    },
}
```

This creates 4 independent DAGs:
- **prerequisites DAG**: Workflow dependencies
- **downstream_required DAG**: Required follow-up tasks
- **validates DAG**: Quality assurance links
- **documents DAG**: Documentation links

Each DAG is validated independently for cycles.

### Transactional Updates

Use the clone-validate-commit pattern for complex operations:

```go
// Clone the manager
testManager := manager.Clone()

// Make changes to the clone
testManager.AddNode(newNode)
testManager.UpdateNode(modifiedNode)

// Validate
if err := testManager.DetectCycles(); err != nil {
    // Changes would create cycle, abort
    return err
}

// If valid, apply changes to original
manager.AddNode(newNode)
manager.UpdateNode(modifiedNode)
```

### Custom Relationship Types

Define domain-specific relationships:

```go
// Code review workflow
reviewRel := types.Relationship{
    Name:        "reviewed_by",
    Description: "Code reviews that validate this change",
    Direction:   types.DirectionNone,
}

// Component dependencies
dependsRel := types.Relationship{
    Name:        "depends_on",
    Description: "Service dependencies",
    Direction:   types.DirectionBackward,
}

// Feature documentation
docsRel := types.Relationship{
    Name:        "documented_in",
    Description: "Documentation for this feature",
    Direction:   types.DirectionForward,
}
```

### Persistence Format

Nodes are stored as YAML files with the structure:

```yaml
id: deploy-api
name: Deploy API Service
summary: Deploy REST API to production
description: |
  Full deployment including:
  - Binary upload
  - Health checks
  - Rollout monitoring
tags:
  - deployment
  - production
edges:
  prerequisites:
    - build-binary
    - run-tests
  downstream_required:
    - smoke-test
created_at: 2024-01-15T10:30:00Z
updated_at: 2024-01-15T10:30:00Z
```

**Note**: Only EdgeIDs are persisted (as `edges` key). The Edges map is runtime-only.

## Architecture

### Two-Level Edge Storage

The package uses a dual-storage approach for edges:

1. **EdgeIDs** (Persisted)
   - Stored as `map[string][]string`
   - Maps relationship names to target node IDs
   - Serialized to YAML/JSON
   - Source of truth on disk

2. **Edges** (Runtime)
   - Stored as `map[string][]Edge`
   - Maps relationship names to resolved edge objects
   - Tagged `json:"-" yaml:"-"` (not serialized)
   - Populated by Manager.ResolveNodePointers()

This separation allows:
- Safe serialization without circular references
- Convenient runtime access to full node objects
- Clear distinction between persisted and derived data

### Cycle Detection

The manager validates all relationship types for cycles using depth-first search:

1. Collects all unique relationship types from nodes
2. For each relationship type, constructs a directed graph
3. Performs DFS to detect back edges (cycles)
4. Reports all cycles found with detailed paths

Example cycle error:
```
detected 2 cycle(s):
  1. prerequisites: task-a -> task-b -> task-c -> task-a
  2. downstream_required: task-x -> task-y -> task-x
```

### Safe Updates

All mutating operations (Add, Update, Delete) follow these patterns:

**Add/Update:**
1. Validate input (non-nil, non-empty ID, exists/doesn't exist)
2. Clone the manager
3. Apply change to clone
4. Detect cycles in clone
5. If valid, apply to original
6. Refresh tag cache and pointers

**Delete:**
1. Validate input
2. Purge node from graph (removes all edges pointing to it)
3. Refresh tag cache

## Testing

The package includes comprehensive tests:

```bash
# Run all tests
go test ./pkg/graph_manager/...

# Run with coverage
go test -cover ./pkg/graph_manager/...

# Run specific test
go test -run TestPersistAndLoad ./pkg/graph_manager
```

## Design Principles

1. **Flexibility over hardcoding**: Define relationship types dynamically
2. **Safety over performance**: Validate before mutating
3. **Separation of concerns**: Persistence vs. runtime representation
4. **Type safety**: Strongly typed throughout
5. **Developer ergonomics**: Dual access patterns (IDs vs. Edges)
6. **Fail fast**: Immediate validation on operations

## Use Cases

### Workflow Management
- Task dependencies and ordering
- Build pipeline stages
- Deployment workflows
- Release checklists

### Dependency Tracking
- Service dependencies
- Package relationships
- Infrastructure dependencies
- Component relationships

### Knowledge Graphs
- Institutional knowledge capture
- Runbook workflows
- Documentation linking
- Quality assurance processes

### Project Management
- Feature relationships
- Epic breakdowns
- Sprint planning
- Milestone tracking

## Contributing

See the main [repository README](../../README.md) for contribution guidelines.

## License

MIT
