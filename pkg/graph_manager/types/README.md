# Graph Types

This package defines the core data structures for task management, with a focus on expressing flexible, configurable relationships between nodes in a directed acyclic graph (DAG).

## Overview: Flexible Relationship System

The task system is built around **configurable directed acyclic graphs (DAGs)** where:
- **Nodes** represent tasks or work items
- **Edges** represent typed relationships between nodes
- **Relationships** define the categories/types of edges and their temporal direction

Unlike systems with hardcoded relationship types, this implementation allows you to define any number of relationship types dynamically. Each relationship type forms its own DAG, sharing the same nodes but representing different semantic connections.

## Core Types

### Node (`node.go`)

A Node represents a task or work item in the system.

**Metadata fields:**
- `ID` - Unique identifier
- `Name` - Human-readable name
- `Summary` - One-line description
- `Description` - Detailed explanation
- `Tags` - Categorization labels
- `CreatedAt`, `UpdatedAt` - Timestamps

**Relationship storage (two-level system):**
- `EdgeIDs map[string][]string` - **Persisted** to disk (YAML/JSON key: "edges")
  - Maps relationship names to lists of target node IDs
  - Example: `{"prerequisites": ["task-a", "task-b"], "validates": ["task-c"]}`
- `Edges map[string][]Edge` - **Runtime only** (not persisted, tagged `json:"-" yaml:"-"`)
  - Maps relationship names to lists of resolved Edge objects with pointers
  - Populated by the Manager after loading from disk

**Why two levels?**
- EdgeIDs are the source of truth on disk (simple, serializable strings)
- Edges provide convenient access to full node objects at runtime
- Separation allows safe serialization without circular references

### Edge (`edge.go`)

An Edge represents a directed connection from one node to another with a specific relationship type.

```go
type Edge struct {
    To   *Node         // Destination node
    Type *Relationship // Relationship category
}
```

Edges are runtime constructs - they're not persisted directly, but reconstructed from EdgeIDs when loading.

### Relationship (`relationship.go`)

A Relationship defines a type/category of relationship between nodes.

```go
type Relationship struct {
    Name        string                 // e.g., "prerequisites", "validates"
    Description string                 // Human-readable explanation
    Direction   RelationshipDirection  // Temporal flow
}
```

### RelationshipDirection (`relationship_direction.go`)

Defines how a relationship flows relative to execution order:

- `DirectionBackward` - Points to nodes that come **before** in execution order
  - Example: "prerequisites" (these must complete before this node)
- `DirectionForward` - Points to nodes that come **after** in execution order
  - Example: "downstream_required" (these must complete after this node)
- `DirectionNone` - No temporal ordering
  - Example: "related_to", "validates" (conceptual links without ordering)

## Example: Task Workflow Relationships

Let's model a microservice deployment workflow with three relationship types:

### Defining Relationships

```go
prerequisitesRel := Relationship{
    Name:        "prerequisites",
    Description: "Tasks that must be completed before this task",
    Direction:   DirectionBackward,
}

downstreamRequiredRel := Relationship{
    Name:        "downstream_required",
    Description: "Tasks that must be completed after this task",
    Direction:   DirectionForward,
}

downstreamSuggestedRel := Relationship{
    Name:        "downstream_suggested",
    Description: "Recommended follow-up tasks",
    Direction:   DirectionForward,
}
```

### Building the Graph

```
Task A: "Create new microservice"
   EdgeIDs: {
       "prerequisites": ["task-x", "task-y"],
       "downstream_required": ["task-b", "task-c"],
       "downstream_suggested": ["task-d", "task-e", "task-f"]
   }

Task X: "Create k8s manifests"
Task Y: "Create container build flow"
Task B: "Run unit tests"
Task C: "Run integration tests"
Task D: "Document API endpoints"
Task E: "Add monitoring dashboards"
Task F: "Create runbook for on-call"
```

### Workflow Semantics

**Prerequisites (DirectionBackward):**
- Task A cannot start until Tasks X and Y complete
- Hard blocker - these must be satisfied first

**Required Downstream Tasks (DirectionForward):**
- After Task A completes, Tasks B and C **must** run
- Hard requirement - workflow incomplete without them

**Suggested Downstream Tasks (DirectionForward):**
- Tasks D, E, and F are recommended but optional
- Soft suggestions - provide guidance without blocking

**Execution order:**
1. Tasks X and Y complete (prerequisites)
2. Task A executes
3. Tasks B and C must run (required downstream)
4. Tasks D, E, F are recommended (suggested downstream)
5. Workflow complete when A, B, C are done

## Flexibility: Custom Relationship Types

The system isn't limited to prerequisites and downstream tasks. You can define any relationships that make sense for your domain:

**Quality assurance:**
```go
{
    Name:        "validates",
    Description: "Tasks that validate this task's output",
    Direction:   DirectionNone,  // Conceptual link, not temporal
}
```

**Component relationships:**
```go
{
    Name:        "depends_on",
    Description: "Components this task depends on",
    Direction:   DirectionBackward,
}
```

**Documentation links:**
```go
{
    Name:        "documents",
    Description: "Documentation tasks for this feature",
    Direction:   DirectionForward,
}
```

## Implementation Patterns

### Creating a Node with Relationships

```go
node := &Node{
    ID:          "deploy-api",
    Name:        "Deploy API Service",
    Summary:     "Deploy the REST API to production",
    Description: "Full deployment including tests and rollout",
    Tags:        []string{"deployment", "production", "api"},
    EdgeIDs: map[string][]string{
        "prerequisites":        {"build-binary", "run-tests"},
        "downstream_required":  {"smoke-test", "update-docs"},
    },
    CreatedAt:   time.Now(),
    UpdatedAt:   time.Now(),
}
```

### Working with Edge Methods

The Node type provides methods to manage both EdgeIDs and Edges:

**ID-level operations (work with strings):**
```go
node.GetEdgeIDs("prerequisites")                    // Returns []string
node.SetEdgeIDs("prerequisites", []string{"a", "b"}) // Replaces all IDs
node.AddEdgeID("prerequisites", "c")                 // Appends one ID
```

**Edge-level operations (work with resolved pointers):**
```go
node.GetEdges("prerequisites")         // Returns []Edge
node.SetEdges("prerequisites", edges)  // Updates both Edges and EdgeIDs
node.AddEdge("prerequisites", edge)    // Updates both Edges and EdgeIDs
```

**Important:** Edge-level operations automatically keep EdgeIDs in sync.

### Serialization Behavior

**Marshalling to YAML/JSON:**
```yaml
id: deploy-api
name: Deploy API Service
edges:
  prerequisites:
    - build-binary
    - run-tests
  downstream_required:
    - smoke-test
    - update-docs
```

Only EdgeIDs are persisted (under the key "edges"). The Edges map is ignored.

**Unmarshalling from YAML/JSON:**
- EdgeIDs map is populated from the "edges" key
- Edges map remains nil (to be resolved by the Manager)

### Cloning Nodes

```go
clone := node.Clone()
```

Clone creates a deep copy of all persisted fields:
- ✅ Copies: ID, Name, Summary, Description, Tags, EdgeIDs, timestamps
- ❌ Skips: Edges map (left as nil, to be re-resolved)

### Comparing Nodes

```go
equal := node1.Equals(node2)
```

Equals compares all persisted fields, including:
- Scalar fields (ID, Name, etc.)
- Tags slice (order matters)
- EdgeIDs map (order within slices matters)
- Timestamps (using time.Equal for proper comparison)

## Multiple DAGs, Shared Nodes

Each relationship type forms an independent DAG:
- Different relationship types can have completely different edge structures
- A node can have different edges in different relationship graphs
- All DAGs share the same set of nodes
- Each DAG must remain acyclic independently

**Example:**
```
Prerequisites DAG: X → A, Y → A
Downstream DAG:    A → B, A → C
Validates DAG:     C → A, B → A
```

Node A participates in all three DAGs with different edge configurations.

## Design Principles

1. **Separation of persistence and runtime** - EdgeIDs are serialized, Edges are computed
2. **Flexibility over hardcoding** - Define relationship types as needed, not baked into the struct
3. **Type safety** - Edges carry their relationship type, preventing category confusion
4. **Dual access patterns** - Work with IDs (simple) or Edges (convenient) as needed
5. **Deep copy support** - Clone nodes safely without worrying about pointer aliasing
6. **Consistency** - Edge methods maintain EdgeIDs and Edges in sync

## Use Cases

**Prerequisites (DirectionBackward):**
- Workflow enforcement (can't deploy before building)
- Resource availability (can't use infrastructure before it's created)
- Logical dependencies (can't test code before it's written)

**Required Downstream Tasks (DirectionForward):**
- Enforcing post-completion actions (must run tests after build)
- Validation steps (must verify after deployment)
- Mandatory follow-up actions (must update docs after feature)

**Suggested Downstream Tasks (DirectionForward):**
- Recommending next actions to users
- Tracking related work items
- Maintaining project context
- Visualizing project flow without strict blocking

**Custom Relationships:**
- Code review workflows ("reviewed_by")
- Component dependencies ("depends_on")
- Test coverage ("tested_by")
- Documentation links ("documented_in")
- Ownership tracking ("owned_by")

The flexibility of the relationship system allows modeling any domain-specific relationships your task management needs require.
