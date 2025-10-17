# Generate Initial Task Set

You have access to the common-tasks MCP server. Your goal is to create discrete, atomic tasks that capture individual operations developers perform repeatedly. Task relationships (prerequisites/downstream) will naturally express complete workflows.

## What to Create

Create tasks for workflows that developers perform repeatedly, such as:
- Running tests (unit, integration, with coverage, with race detection)
- Code quality checks (formatting, linting, type checking)
- Building artifacts (binaries, containers, packages)
- Deployment workflows
- Development setup steps
- Git/commit workflows
- Documentation updates (updating README.md, API docs, CONTRIBUTING.md)
- CI/CD workflow updates (modifying GitHub Actions, GitLab CI, CircleCI configs)

**Good starting points**: Tasks for updating documentation (README.md, CONTRIBUTING.md), updating CI flows (.github/workflows/, .gitlab-ci.yml), and other change-related workflows are excellent foundational tasks to create first, as these represent common operations that accompany most development work.

## What NOT to Create

Avoid creating meta-workflow or umbrella tasks that simply aggregate other tasks:
- ❌ "Complete Development Workflow" - the relationships already define this
- ❌ "Release Workflow" - relationships between tasks express the release flow
- ❌ "Setup Development Environment" - too broad, break into discrete steps
- ❌ "Pre-commit Checks" - the format-code → lint-code → run-tests relationships already capture this

**Why?** The task relationship DAG IS the workflow. Meta-tasks duplicate what relationships already express and create unnecessary complexity.

**Create discrete, atomic tasks instead:**
- ✅ "Run Unit Tests" with `go test ./...`
- ✅ "Format Code" with `gofmt -w .`
- ✅ "Build Docker Image" with `docker build`
- ✅ "Run Tests with Coverage" - a specific variant with distinct flags

Each task should represent a single, actionable operation that a developer would run directly.

## Where to Find Information

**IMPORTANT**: Base tasks on actual workflows found in the codebase, not speculation. Only create tasks for operations that developers have actually performed or that are explicitly defined in build systems and CI/CD configs.

1. **Build systems**: Check for Makefile, Taskfile.yml, package.json scripts, build.gradle, etc. - these define tasks developers actually run
2. **CI/CD**: Look at .github/workflows/, .gitlab-ci.yml, .circleci/, etc. - these show automated workflows that are actually used
3. **Documentation**: Read README.md, CONTRIBUTING.md, docs/ - these document actual development practices
4. **Test patterns**: Examine how tests are organized and run - look for actual test commands and patterns
5. **Scripts**: Check scripts/ or similar directories - these are real automation that developers use
6. **Project files**: go.mod, Cargo.toml, requirements.txt, etc. - these show actual dependencies and tooling
7. **Git history**: Use `git log`, `git log --oneline`, and commit messages to understand:
   - Common workflows (what commands appear in commit messages)
   - Frequent change patterns (what files/areas change together)
   - Recent development focus (what's being actively worked on)
   - Release processes (tags, version bumps, changelog patterns)
   - Team practices (commit message conventions, branch patterns)

Do NOT create tasks based on assumptions about what might be useful. If you don't see evidence of a workflow being used, don't create a task for it.

## Task Structure

Each task should include:
- **ID**: kebab-case identifier
- **Name**: Human-readable name
- **Summary**: One-line description
- **Description**: Detailed explanation with actual commands from this codebase
- **Tags**: For categorization - use 2-4 simple, clear tags (e.g., cicd, k8s, go, java, build, testing, deployment)
- **Prerequisites**: Tasks that MUST complete before this one
- **Required downstream**: Tasks that MUST follow this one
- **Suggested downstream**: Tasks that are recommended after this one

## Guidelines

1. **Use actual commands** from the codebase (exact paths, tool versions if specified)
2. **Verify commands before adding tasks** - for non-destructive tasks (tests, builds, formatters, linters), run the command first to verify it works before creating the task. This ensures you're documenting real, working workflows.
3. **Create workflows** - connect related tasks with prerequisites/downstream
4. **Be specific** - include flags, arguments, paths from the project
5. **Think about order** - what must happen before/after each task?
6. **Cover common scenarios** - what do developers do daily? weekly? on release?
7. **Keep tags simple** - prefer clear, single-word tags (cicd, k8s, go, java, build). Don't overuse tags; only add new ones when necessary for clarity and organization. Aim for 2-4 tags per task.
8. **Include directory paths** - mention specific directories in task descriptions to help LLMs navigate directly to relevant code locations (e.g., "Tests are located in pkg/task_manager/" or "Configuration files in config/"). DO NOT list the contents of files or enumerate what's inside directories - just provide the paths for navigation.
9. **Keep tasks atomic** - each task should be a single, discrete operation. Don't create meta-tasks that just aggregate other tasks. The DAG relationships express workflows naturally without needing umbrella tasks.
10. **Keep descriptions concise** - focus on what the task does and where relevant files are located. Avoid listing file contents, dependency versions, environment variables, or other detailed configuration unless absolutely necessary to run the command.
11. **Require evidence for changes** - when modifying existing tasks (updating commands, changing descriptions, or altering relationships), you MUST provide evidence from the codebase that justifies the change. Evidence includes: recent commits showing the change, updated documentation, modified build files, or actual command execution showing the new behavior. Do not modify tasks based on speculation or assumptions.

## Example Task Relationships

```
format-code → run-tests → build-binary
                ↓
         run-tests-with-coverage (suggested)
```

## Start

Begin by exploring the codebase to understand its structure, build system, and common workflows. Then create tasks systematically, starting with fundamental operations (testing, formatting) and building up to complex workflows (deployment, release).
