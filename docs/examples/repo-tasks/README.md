# Repository Tasks Example

This example configures the MCP server for managing development workflows and repository tasks.

## Use Case

Track institutional knowledge about development workflows:
- Build and deployment procedures
- Testing workflows
- CI/CD pipeline tasks
- Code review processes
- Release procedures
- Database migration workflows

## Configuration

### MCP Tools Generated

With this configuration, the server exposes:
- `add_task` - Create a new development task
- `get_task` - Retrieve a task with its full workflow
- `list_tasks` - Browse tasks by tags
- `update_task` - Modify an existing task
- `delete_task` - Remove a task

### Relationship Types

**prerequisites** (backward direction)
- Points to tasks that must complete before this task can start
- Example: "run-tests" requires "build-binary" as a prerequisite
- Forms a dependency DAG

**downstream_required** (forward direction)
- Points to tasks that MUST follow this task
- Example: After "create-migration", you MUST run "apply-migration"
- Enforces mandatory workflow steps

**downstream_suggested** (forward direction)
- Points to tasks that SHOULD follow this task
- Example: After "deploy-api", you SHOULD "update-api-docs"
- Provides guidance without strict enforcement

## Example Task

```yaml
id: deploy-api-service
name: Deploy API Service to Production
summary: Deploy the REST API service to production environment
description: |
  Deploys the API service to production using the deployment pipeline.

  Steps:
  1. Ensure all tests pass
  2. Build the production binary
  3. Push to container registry
  4. Update kubernetes manifests
  5. Apply with kubectl
  6. Verify health checks

  See deployment runbook: docs/deployment.md
tags:
  - deployment
  - production
  - api
edges:
  prerequisites:
    - run-integration-tests
    - run-e2e-tests
    - build-production-binary
  downstream_required:
    - verify-health-checks
    - smoke-test-production
  downstream_suggested:
    - update-api-documentation
    - notify-stakeholders
    - update-changelog
created_at: 2024-01-15T10:30:00Z
updated_at: 2024-01-15T10:30:00Z
```

## Common Tags

Suggested tags for organizing development tasks:

- **Lifecycle**: `build`, `test`, `deploy`, `release`
- **Environment**: `development`, `staging`, `production`
- **Component**: `frontend`, `backend`, `database`, `api`, `infrastructure`
- **Type**: `cicd`, `documentation`, `migration`, `configuration`
- **Technology**: `go`, `python`, `docker`, `kubernetes`, `postgres`

## Workflow Example

A typical deployment workflow might look like:

```
run-unit-tests ────┐
                   ├──> build-production-binary ──> deploy-api-service ──> verify-health-checks
run-integration-tests ─┘                                    │
                                                            └──> update-api-documentation
```

This structure is captured automatically through the task relationships, allowing AI assistants to understand and guide developers through the complete workflow.

## Running This Example

```bash
# Start the server with this configuration
mcp serve --directory ./docs/examples/repo-tasks

# Or copy to your data directory
cp docs/examples/repo-tasks/mcp.yaml ./my-tasks/
cp docs/examples/repo-tasks/relationships.yaml ./my-tasks/
mcp serve --directory ./my-tasks
```
