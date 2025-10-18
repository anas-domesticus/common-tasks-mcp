# Configuration Examples

This directory contains example configurations demonstrating how to use the Generic DAG MCP Server for different domains.

## Overview

Each example includes:
- `mcp.yaml` - Server identity and terminology configuration
- `relationships.yaml` - Relationship type definitions for that domain
- `README.md` - Documentation and usage examples
- Sample node structures showing typical data

## Available Examples

### [Repository Tasks](./repo-tasks/)
**Domain**: Software development workflows and task management

**Tools Generated**: `add_task`, `get_task`, `list_tasks`, `update_task`, `delete_task`

**Relationships**:
- `prerequisites` - Tasks that must complete first
- `downstream_required` - Tasks that must follow
- `downstream_suggested` - Recommended next steps

**Use Cases**:
- CI/CD pipeline workflows
- Build and deployment procedures
- Testing workflows
- Release checklists
- Database migrations

### [Recipe Knowledge Base](./recipes/)
**Domain**: Culinary recipes and meal planning

**Tools Generated**: `add_recipe`, `get_recipe`, `list_recipes`, `update_recipe`, `delete_recipe`

**Relationships**:
- `requires` - Ingredients and base recipes needed
- `produces` - Dishes or components created
- `pairs_with` - Complementary recipes
- `variations` - Alternative preparations

**Use Cases**:
- Meal planning
- Shopping list generation
- Recipe pairing suggestions
- Dietary restriction management
- Ingredient substitution

## Using an Example

### Quick Start

```bash
# Navigate to an example directory
cd docs/examples/repo-tasks

# Start the server with this configuration
mcp serve --directory .
```

### Copy to Your Data Directory

```bash
# Create your data directory
mkdir -p ./my-data

# Copy configuration files
cp docs/examples/repo-tasks/mcp.yaml ./my-data/
cp docs/examples/repo-tasks/relationships.yaml ./my-data/

# Start the server
mcp serve --directory ./my-data
```

### Customize for Your Domain

1. Copy an example that's closest to your use case
2. Edit `mcp.yaml`:
   - Change `server.name` to your domain
   - Update `server.display_name`
   - Customize `server.instructions`
   - Modify `naming.node` terminology
3. Edit `relationships.yaml`:
   - Add, remove, or rename relationship types
   - Adjust directions (backward/forward/none)
   - Update descriptions
4. Start creating nodes with your domain data

## Creating Your Own Configuration

### Minimal Configuration

**mcp.yaml**:
```yaml
server:
  name: my-domain-mcp
  display_name: My Domain
  instructions: Description of your domain and how to use it.

naming:
  node:
    singular: item
    plural: items
    display_singular: Item
    display_plural: Items
```

**relationships.yaml**:
```yaml
relationships:
  - name: depends_on
    description: Items this item depends on
    direction: backward

  - name: enables
    description: Items this item enables
    direction: forward
```

### Relationship Direction Guidelines

Choose the appropriate direction for your relationship type:

**backward** - Points to things that come **before**:
- Dependencies (`depends_on`, `requires`, `prerequisites`)
- Inputs (`uses`, `consumes`, `built_from`)
- Blockers (`blocked_by`, `waits_for`)

**forward** - Points to things that come **after**:
- Outputs (`produces`, `generates`, `creates`)
- Consequences (`triggers`, `enables`, `unlocks`)
- Follow-ups (`leads_to`, `suggests`)

**none** - No temporal ordering:
- Associations (`related_to`, `similar_to`, `pairs_with`)
- References (`documents`, `tests`, `validates`)
- Categories (`type_of`, `instance_of`)

## More Example Ideas

### Service Dependency Graph
```yaml
naming:
  node: { singular: service, plural: services }
relationships:
  - name: depends_on
  - name: consumed_by
  - name: shares_data_with
```

### Learning Path System
```yaml
naming:
  node: { singular: lesson, plural: lessons }
relationships:
  - name: prerequisites
  - name: unlocks
  - name: practice_for
```

### Project Management
```yaml
naming:
  node: { singular: story, plural: stories }
relationships:
  - name: blocked_by
  - name: blocks
  - name: related_to
```

### Infrastructure as Code
```yaml
naming:
  node: { singular: component, plural: components }
relationships:
  - name: requires_infrastructure
  - name: provides_services
  - name: monitored_by
```

### Documentation Graph
```yaml
naming:
  node: { singular: document, plural: documents }
relationships:
  - name: references
  - name: supersedes
  - name: related_topics
```

## Contributing Examples

If you create an interesting configuration for a new domain, consider contributing it:

1. Create a new directory under `docs/examples/your-domain`
2. Include `mcp.yaml`, `relationships.yaml`, and `README.md`
3. Add sample node structures in the README
4. Update this index with your example
5. Submit a pull request

## Testing Your Configuration

```bash
# Validate your configuration files are valid YAML
yamllint mcp.yaml relationships.yaml

# Start the server in verbose mode to see detailed logs
mcp serve --directory . --verbose

# Test the MCP tools (requires an MCP client)
# The tools will be named based on your configuration:
# - list_{plural}
# - get_{singular}
# - add_{singular}
# - update_{singular}
# - delete_{singular}
```

## Need Help?

- Check the [main README](../../README.md) for architecture details
- Review the [graph_manager documentation](../../pkg/graph_manager/README.md)
- Look at existing examples for patterns
- Open an issue if you have questions
