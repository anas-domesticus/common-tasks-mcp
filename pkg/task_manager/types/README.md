# Task Types

This package defines the core data structures for task management, with a focus on expressing relationships between tasks.

## Task Relationships: Three DAGs

The task system is built around **three distinct Directed Acyclic Graphs (DAGs)** that share the same nodes (tasks) but represent different types of relationships:

1. **Prerequisites** - Tasks that must be completed before/as part of this task
2. **Required Downstream Tasks** - Tasks that must be completed after this task
3. **Suggested Downstream Tasks** - Tasks that are recommended or related, but optional

All three graphs must be acyclic to maintain a valid task system.

### 1. Prerequisites (DAG #1)

Prerequisites represent tasks that **must be completed before** or **are components of** the current task. The current task cannot be considered to be completed unless all its prerequisites are completed.

- **Field**: `PrerequisiteIDs` (persisted) / `Prerequisites` (resolved pointers)
- **Meaning**: "This task cannot proceed until these tasks are complete"
- **Constraint**: Hard blocker - must be satisfied
- **Use case**: Enforcing execution order, blocking tasks until prerequisites are met

**Example**: Task A is "Create new microservice"
- Prerequisites might include:
  - Task X: "Create k8s manifests"
  - Task Y: "Create container build flow"

Task A is **blocked** and cannot start until both Task X and Task Y are completed.

### 2. Required Downstream Tasks (DAG #2)

Required downstream tasks are tasks that **must be completed after** the current task finishes. These represent mandatory follow-up actions that are triggered by or necessitated by completing the current task.

- **Field**: `DownstreamRequiredIDs` (persisted) / `DownstreamRequired` (resolved pointers)
- **Meaning**: "These tasks must be done after this task completes"
- **Constraint**: Hard requirement - must be completed as part of the workflow
- **Use case**: Enforcing post-completion actions, validation steps, mandatory follow-ups

**Example**: Continuing with Task A "Create new microservice"
- Required downstream tasks might include:
  - Task B: "Run unit tests for new microservice"
  - Task C: "Run integration tests"

After Task A completes, Tasks B and C **must be run** before the workflow can be considered complete.

### 3. Suggested Downstream Tasks (DAG #3)

Suggested downstream tasks are tasks that are **conceptually related** or **recommended next steps**, but are not mandatory. These represent optional work that logically follows from the current task.

- **Field**: `DownstreamSuggestedIDs` (persisted) / `DownstreamSuggested` (resolved pointers)
- **Meaning**: "These tasks are recommended or related, but optional"
- **Constraint**: Soft suggestion - no blocking
- **Use case**: Suggesting next actions to users, tracking related work, maintaining project context

**Example**: Continuing with Task A "Create new microservice"
- Downstream suggested tasks might include:
  - Task D: "Document API endpoints"
  - Task E: "Add monitoring dashboards"
  - Task F: "Create runbook for on-call"

These tasks are good ideas and logically follow from creating a new service, but they could be done independently, deferred, or skipped based on priorities.

### Key Distinctions

**Three independent DAGs:**
- Each DAG represents a different type of relationship between tasks
- A task can appear in multiple DAGs with different edges
- The edges in one DAG are independent of edges in another DAG

**Constraint levels:**
1. **Prerequisites**: Hard blocker - cannot consider this task to be complete until these are also completed
2. **Required Downstream Tasks**: Hard requirement - must complete these tasks after this task finishes
3. **Suggested Downstream Tasks**: Soft suggestion - optional related work

**Workflow implications:**
- Prerequisites block task start
- Required Downstream Tasks block workflow completion
- Suggested Downstream Tasks provide guidance but don't block anything

### Visual Example

```
Task A: "Create new microservice"
   Prerequisites (MUST complete first):
      Task X: "Create k8s manifests"
      Task Y: "Create container build flow"

   Required Downstream Tasks (MUST complete after):
      Task B: "Run unit tests for new microservice"
      Task C: "Run integration tests"

   Suggested Downstream Tasks (OPTIONAL):
      Task D: "Document API endpoints"
      Task E: "Add monitoring dashboards"
      Task F: "Create runbook for on-call"
```

**Workflow execution:**
1. Tasks X and Y must complete **before** Task A can start (Prerequisites)
2. Task A executes
3. After Task A completes, Tasks B and C **must be run** (Required Downstream Tasks)
4. Tasks D, E, and F are **recommended** but optional (Suggested Downstream Tasks)
5. The workflow is only complete when A, B, and C are all done (X and Y are prerequisites)

## Implementation Notes

- **Persisted fields**: `PrerequisiteIDs`, `DownstreamRequiredIDs`, and `DownstreamSuggestedIDs` store task IDs as strings
- **Runtime fields**: `Prerequisites`, `DownstreamRequired`, and `DownstreamSuggested` store resolved pointers to actual Task objects (not persisted)
- The resolved pointers are populated by the task manager when loading and building the task graph

## Use Cases

### Prerequisites
- Workflow enforcement (can't deploy before building)
- Resource availability (can't use infrastructure before it's created)
- Logical prerequisites (can't test code before it's written)

### Required Downstream Tasks
- Enforcing post-completion actions
- Validation steps that must run after task completion
- Mandatory follow-up actions

### Suggested Downstream Tasks
- Suggesting next actions to users
- Tracking related work items
- Maintaining project context
- Visualizing project flow without strict blocking
