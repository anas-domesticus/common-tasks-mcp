# Claude Code Integration Guide

This MCP server helps capture and preserve institutional knowledge as structured task definitions.

## Major Components

### Available Tools
- `list_tasks`: Browse tasks, optionally filtered by tags
- `get_task`: Get full task details including relationships
- `list_tags`: See all available tags with usage counts
- `add_task`: Create new tasks (requires all relationship fields)
- `update_task`: Modify existing tasks
- `delete_task`: Remove tasks (automatically cleans up references)

### Task Structure
Each task has:
- **Identity**: `id`, `name`, `summary`, `description`
- **Organization**: `tags` for categorization
- **Relationships**: `prerequisiteIDs`, `downstreamRequiredIDs`, `downstreamSuggestedIDs`

### Relationships
- **Prerequisites**: Tasks that must be completed first (backward direction)
- **Downstream Required**: Tasks that must follow (forward direction)
- **Downstream Suggested**: Recommended follow-up tasks (forward direction)

## When to Add Tasks

**After completing commonly performed work**, consider capturing it as a task if:
- It's likely to be repeated in this codebase
- It represents a standard development operation
- It would benefit future developers (human or AI)
- It has clear prerequisites or downstream requirements

**Examples**: Adding API endpoints, running tests with specific flags, deployment procedures, setting up CI/CD workflows, database migrations, dependency updates.

**Don't capture**: One-off bug fixes, simple typo corrections, exploratory changes.

## Best Practices

1. **Check first**: Use `list_tasks` to see if a similar task exists
2. **Be specific**: Include actual commands, file locations, and concrete steps
3. **Stay atomic**: One clear purpose per task, use relationships to connect them
4. **Tag appropriately**: 2-4 tags (tech stack, operation type, component, etc.)
5. **No meta-tasks**: Don't create umbrella tasks—use relationships instead

By capturing commonly performed tasks, you build institutional knowledge that helps everyone working on the codebase.
