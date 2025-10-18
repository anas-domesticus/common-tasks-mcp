---
description: "Prompt for capturing workflows as tasks during active development. Guides recognition of repeatable operations and helps maintain tasks as you work, ensuring institutional knowledge is captured in real-time."
---

# Capture Workflow as Task

As you work on development tasks, recognize when your current work represents a repeatable workflow that should be captured in the task system for future reference.

## When to Capture a Task

Capture tasks when you:
- **Run commands repeatedly**: If you just ran the same command for the second or third time
- **Follow a multi-step process**: When steps naturally follow each other (format → test → build)
- **Discover how to do something**: When you figure out the right command or sequence
- **Fix or update a workflow**: When you realize a task's command or process has changed
- **Complete a common operation**: Testing, building, deploying, formatting, etc.

## What to Capture

### The Basics
- **The actual command you ran**: Copy it exactly, with flags and paths
- **Where it should be run from**: Working directory context if relevant
- **What files/directories it affects**: Just the paths, not the contents
- **What needs to happen first**: Prerequisites you discovered along the way

### Workflow Context
- If you ran multiple commands in sequence, consider if each should be its own task
- Note which tasks naturally follow each other (format before test, test before build)
- Identify which relationships are required vs. suggested

## How to Capture

1. **Check if task exists**: List tasks to see if this workflow is already captured
2. **Create or update**:
   - Create new task if this is a novel workflow
   - Update existing task if you found the command has changed - **but you MUST provide evidence** (see below)
3. **Connect the workflow**: Add prerequisite/downstream relationships based on what you actually did
4. **Keep it simple**: Command + relevant paths. No need to explain why or document edge cases

### Evidence Required for Updates

When updating an existing task, you MUST provide evidence from the codebase that justifies the change:
- **Command changes**: Show recent commits, updated build files (Makefile, Taskfile.yml, package.json), or output from running both old and new commands
- **Description changes**: Reference documentation updates, code changes, or commit messages
- **Relationship changes**: Demonstrate the new workflow order through actual usage, CI/CD changes, or script modifications

Do NOT update tasks based on:
- ❌ Speculation about what might work better
- ❌ Assumptions about how things should be done
- ❌ Personal preferences without codebase evidence
- ❌ Changes you haven't actually verified

## Guidelines

- **Capture what you actually did**, not what the docs say or what you think should happen
- **Verify before capturing**: If creating a task proactively (not based on a command you just ran), verify non-destructive commands work before adding them to the task system
- **One task = one command** (or one clear atomic operation)
- **Use the exact command** with the exact flags and paths you used
- **Don't overthink it**: If you ran it twice, it's probably worth capturing
- **Update immediately**: Don't wait - capture while it's fresh
- **Skip one-off operations**: If this is truly unique to one situation, don't create a task

## Examples

**You just ran**: `go test ./pkg/...` three times while debugging
→ Create task: "run-pkg-tests" with that exact command

**You ran**: `gofmt -w . && go test ./...` as a sequence
→ Two tasks: "format-go-code" → "run-unit-tests" with relationship

**You discovered**: The build command changed from `go build ./cli/app` to `go build -o bin/app ./cli/app`
→ Update task: "build-app" with the new command

**You completed**: Run tests, then build Docker image, then start compose
→ Three tasks with downstream relationships: test → build → run

## Anti-patterns

❌ Don't create tasks for operations you haven't actually performed
❌ Don't create umbrella tasks like "complete-development-workflow"
❌ Don't add detailed explanations or documentation to descriptions
❌ Don't capture commands you only ran once for a unique situation
❌ Don't create tasks based on what you think developers might need
❌ Don't update existing tasks without evidence from the codebase justifying the change

## Remember

The goal is to **capture institutional knowledge** as you work. When you figure out how to do something, or when you follow a workflow, capture it so the next person (or AI assistant) can benefit from your discovery.

Tasks should reflect **reality** - what actually gets run - not ideals or documentation.
