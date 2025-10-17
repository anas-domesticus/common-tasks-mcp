# Generate Initial Task Set

You have access to the common-tasks MCP server. Your goal is to create a comprehensive set of tasks that capture the common, repeatable workflows for this codebase.

## What to Create

Create tasks for workflows that developers perform repeatedly, such as:
- Running tests (unit, integration, with coverage, with race detection)
- Code quality checks (formatting, linting, type checking)
- Building artifacts (binaries, containers, packages)
- Deployment workflows
- Development setup steps
- Git/commit workflows

## Where to Find Information

1. **Build systems**: Check for Makefile, Taskfile.yml, package.json scripts, build.gradle, etc.
2. **CI/CD**: Look at .github/workflows/, .gitlab-ci.yml, .circleci/, etc.
3. **Documentation**: Read README.md, CONTRIBUTING.md, docs/
4. **Test patterns**: Examine how tests are organized and run
5. **Scripts**: Check scripts/ or similar directories
6. **Project files**: go.mod, Cargo.toml, requirements.txt, etc.
7. **Git history**: Use `git log`, `git log --oneline`, and commit messages to understand:
   - Common workflows (what commands appear in commit messages)
   - Frequent change patterns (what files/areas change together)
   - Recent development focus (what's being actively worked on)
   - Release processes (tags, version bumps, changelog patterns)
   - Team practices (commit message conventions, branch patterns)

## Task Structure

Each task should include:
- **ID**: kebab-case identifier
- **Name**: Human-readable name
- **Summary**: One-line description
- **Description**: Detailed explanation with actual commands from this codebase
- **Tags**: For categorization (testing, build, deployment, etc.)
- **Prerequisites**: Tasks that MUST complete before this one
- **Required downstream**: Tasks that MUST follow this one
- **Suggested downstream**: Tasks that are recommended after this one

## Guidelines

1. **Use actual commands** from the codebase (exact paths, tool versions if specified)
2. **Create workflows** - connect related tasks with prerequisites/downstream
3. **Be specific** - include flags, arguments, paths from the project
4. **Think about order** - what must happen before/after each task?
5. **Cover common scenarios** - what do developers do daily? weekly? on release?

## Example Task Relationships

```
format-code → run-tests → build-binary
                ↓
         run-tests-with-coverage (suggested)
```

## Start

Begin by exploring the codebase to understand its structure, build system, and common workflows. Then create tasks systematically, starting with fundamental operations (testing, formatting) and building up to complex workflows (deployment, release).
