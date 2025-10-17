# Claude Code Integration Guide

This document provides guidance for Claude Code (and other AI coding assistants) on how to use the Common Tasks MCP server to capture and preserve institutional knowledge.

## Overview

The Common Tasks MCP server helps break down knowledge silos by capturing repeatable development workflows as structured task definitions. When you (Claude) complete work, you should proactively identify whether that work represents a repeatable task and capture it for future reference.

## When to Capture Tasks

At the end of a **logical unit of work**, evaluate whether you've completed a task that:
- Is likely to be repeated in this codebase
- Represents a common development operation
- Would benefit other developers (human or AI) working on this project
- Has clear prerequisites or downstream requirements

### Examples of Logical Units of Work

✅ **Should be captured:**
- Adding a new API endpoint
- Setting up a new database migration
- Configuring a new CI/CD workflow
- Adding a new dependency or package
- Creating a new component or module
- Running tests in a specific way
- Building or deploying the application
- Updating documentation after a feature

❌ **Should NOT be captured:**
- One-off bug fixes for unique issues
- Exploratory code changes
- Simple typo corrections
- Single-file formatting changes
- Answering questions without making changes

## How to Capture Tasks

### Step 1: Identify the Task

After completing a logical unit of work, ask yourself:
1. What did I just accomplish?
2. Would someone else need to do this again?
3. What are the actual commands or steps I performed?
4. What needed to happen first (prerequisites)?
5. What should happen next (downstream tasks)?

### Step 2: Query for Existing Tasks

Before creating a new task, check if a similar task already exists:

```
Use mcp__common-tasks__list_tasks with relevant tags to see if this workflow is already captured.
```

### Step 3: Create or Update the Task

Use the `mcp__common-tasks__add_task` tool to capture the workflow:

```javascript
mcp__common-tasks__add_task({
  id: "descriptive-kebab-case-id",
  name: "Human Readable Task Name",
  summary: "One-line description of what this task does",
  description: "Detailed explanation including:\n- Actual commands used\n- File locations\n- Context about when to use this\n- Any important notes or gotchas",
  tags: ["relevant", "tags", "for", "categorization"],
  prerequisiteIDs: ["tasks-that-must-run-first"],
  downstreamRequiredIDs: ["tasks-that-must-follow"],
  downstreamSuggestedIDs: ["recommended-next-tasks"]
})
```

If updating an existing task, use `mcp__common-tasks__update_task` instead.

## Best Practices for Task Descriptions

### Be Specific and Actionable

❌ **Vague:** "Update the configuration"
✅ **Specific:** "Add new environment variable to docker-compose.yml for feature flags"

### Include Actual Commands

❌ **Abstract:** "Run the tests"
✅ **Concrete:** "Run integration tests with `npm run test:integration`"

### Mention File Locations

❌ **Generic:** "Update the API documentation"
✅ **Located:** "Update API documentation in `docs/api/` directory"

### Provide Context

```yaml
description: |-
  Adds a new REST API endpoint to the server.

  Steps:
  1. Create handler function in `pkg/handlers/`
  2. Register route in `pkg/server/routes.go`
  3. Add request/response types in `pkg/types/`
  4. Write unit tests in `pkg/handlers/*_test.go`
  5. Update OpenAPI spec in `api/openapi.yaml`

  See existing handlers for patterns and conventions.
```

## Workflow Relationships

### Prerequisites

Tasks that **must** be completed before this one:

```yaml
prerequisites:
  - install-dependencies
  - setup-database
```

### Required Downstream

Tasks that **must** follow this one:

```yaml
downstream_required:
  - run-database-migrations
```

### Suggested Downstream

Tasks that are **recommended** after this one:

```yaml
downstream_suggested:
  - run-unit-tests
  - update-readme
```

## Task Granularity

### Atomic Tasks (Preferred)

Create discrete, single-purpose tasks:

✅ `format-code` → `run-linter` → `run-unit-tests` → `build-binary`

Each task is atomic and has clear relationships.

### Avoid Meta-Tasks

❌ Don't create umbrella tasks that just aggregate others:

```yaml
# ❌ BAD - Meta-task
id: complete-development-workflow
name: Complete Development Workflow
description: "Run formatting, linting, testing, and building"
```

The relationships between tasks already express the workflow. Meta-tasks add unnecessary complexity.

## Proactive Workflow Capture

### Scenario: Adding a New Feature

You've just added a new feature that required:
1. Creating new database migrations
2. Adding API endpoints
3. Writing tests
4. Updating documentation

**Action:** Create/update tasks for each atomic operation:
- `create-database-migration` (if not exists)
- `add-api-endpoint` (if pattern is repeatable)
- `update-api-documentation` (if not exists)

Link them with appropriate relationships.

### Scenario: Setting Up CI/CD

You've configured a new GitHub Actions workflow.

**Action:** Create task:
```yaml
id: add-github-workflow
name: Add GitHub Actions Workflow
summary: Create a new CI/CD workflow in .github/workflows/
description: |-
  Creates a new GitHub Actions workflow file.

  1. Create YAML file in `.github/workflows/`
  2. Define triggers (push, pull_request, schedule)
  3. Define jobs and steps
  4. Configure secrets in repository settings if needed

  See existing workflows for examples and conventions.
tags:
  - cicd
  - github
downstream_suggested:
  - update-contributing
  - update-readme
```

### Scenario: Running Tests in a Specific Way

You discovered the project has integration tests that need special flags.

**Action:** Capture the command:
```yaml
id: run-integration-tests
name: Run Integration Tests
summary: Run integration tests with database dependencies
description: |-
  Runs integration tests that require a running PostgreSQL database.

  Command: `DATABASE_URL=postgres://localhost/test_db npm run test:integration`

  Prerequisites:
  - PostgreSQL must be running (via Docker Compose or local install)
  - Test database must be created and migrated

  Integration tests are located in `tests/integration/`
tags:
  - testing
  - integration
prerequisites:
  - start-database
  - run-migrations
```

## Integration with Development Workflow

### Pattern 1: Post-Implementation Capture

After completing work:

1. **Reflect** on what you just did
2. **Check** if similar task exists using `list_tasks`
3. **Create/Update** task definition with `add_task` or `update_task`
4. **Inform** the user that you've captured the workflow

Example message:
> I've captured this workflow as a task (`add-api-endpoint`) so it can be easily referenced in the future. The task includes the steps we just followed and links to the testing workflow.

### Pattern 2: Real-time Workflow Recognition

While working, recognize when you're following a repeatable pattern:

1. **Recognize** the pattern early
2. **Query** existing tasks to see if it's already documented
3. **Follow** the documented workflow if it exists
4. **Update** the task if you discover improvements

### Pattern 3: Workflow Suggestion

When a user asks you to do something, check if there's already a documented workflow:

1. **Query** tasks using `list_tasks` with relevant tags
2. **Retrieve** full task details with `get_task`
3. **Follow** the documented steps
4. **Inform** user you're following the established workflow

## Tags for Organization

Use clear, simple tags for categorization:

- **Language/Tech**: `go`, `python`, `javascript`, `docker`, `kubernetes`
- **Operation Type**: `testing`, `build`, `deployment`, `documentation`
- **Workflow Stage**: `cicd`, `development`, `production`, `setup`
- **Component**: `backend`, `frontend`, `database`, `api`

Aim for 2-4 tags per task. Don't over-tag.

## Example Workflow Session

```
User: "Add a new feature to handle user authentication"

Claude: [Implements authentication feature]
Claude: [Runs tests to verify]
Claude: [Updates documentation]

Claude: I've completed the authentication feature. Let me capture this workflow...

[Calls mcp__common-tasks__add_task to create "add-authentication-endpoint" task]

Claude to User:
"I've added the authentication feature and captured this workflow as a task
for future reference. The task documents:
- Creating the auth handlers in pkg/handlers/auth.go
- Adding JWT middleware
- Writing tests in pkg/handlers/auth_test.go
- Updating the API documentation

This will help maintain consistency when adding more authenticated endpoints."
```

## Remember

🎯 **Capture workflows, not one-offs**
📝 **Be specific with commands and locations**
🔗 **Define relationships between tasks**
🏷️ **Use clear, simple tags**
🚫 **Avoid meta-tasks - use relationships instead**
💡 **Update existing tasks when you learn better approaches**

By capturing tasks as you work, you're building institutional knowledge that helps everyone (humans and AIs) working on the codebase.
