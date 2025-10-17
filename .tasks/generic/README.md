# Generic Task Templates

This directory contains language-agnostic and framework-agnostic task templates that can be adapted to work with most codebases. These tasks represent common software development workflows that apply across different technology stacks.

## Purpose

Generic tasks provide:
- **Starting point** for creating project-specific tasks
- **Reference documentation** for common development workflows
- **Cross-language examples** showing equivalent commands across different ecosystems
- **Best practices** for common operations

## How to Use

1. **Reference**: Use these as documentation when you need to perform a common task
2. **Copy & Customize**: Copy a generic task and adapt it to your specific project needs
3. **Learn**: See how different languages/frameworks approach the same problem

## Task Categories

### Code Quality & Formatting
- `format-code.yaml` - Format source code using language-specific formatters
- `lint-code.yaml` - Run static analysis and linting tools
- `check-code-style.yaml` - Verify code follows style guidelines (dry-run)
- `run-type-check.yaml` - Run static type checking

### Testing
- `run-unit-tests.yaml` - Run all unit tests
- `run-tests-with-coverage.yaml` - Run tests with coverage reporting
- `run-integration-tests.yaml` - Run integration tests
- `run-e2e-tests.yaml` - Run end-to-end tests

### Build & Dependencies
- `install-dependencies.yaml` - Install or download project dependencies
- `update-dependencies.yaml` - Update dependencies to latest compatible versions
- `clean-build-artifacts.yaml` - Remove compiled binaries and build outputs
- `build-production.yaml` - Create optimized production build

### Docker & Containers
- `build-docker-image.yaml` - Build Docker container image
- `start-docker-compose.yaml` - Start services with Docker Compose
- `stop-docker-compose.yaml` - Stop Docker Compose services
- `view-docker-logs.yaml` - View container logs

### Database Operations
- `run-database-migrations.yaml` - Apply pending database migrations
- `backup-database.yaml` - Create database backup
- `seed-database.yaml` - Populate database with test data

### Documentation
- `update-readme.yaml` - Update README.md file
- `update-contributing.yaml` - Update CONTRIBUTING.md
- `update-changelog.yaml` - Update CHANGELOG.md
- `generate-api-docs.yaml` - Generate API documentation from code

### CI/CD & Git
- `update-github-workflow.yaml` - Modify GitHub Actions workflows
- `update-gitlab-ci.yaml` - Modify GitLab CI pipeline
- `create-git-tag.yaml` - Create git tag for release
- `setup-pre-commit-hooks.yaml` - Configure git pre-commit hooks

### Security
- `run-security-scan.yaml` - Scan dependencies for vulnerabilities

## Language/Framework Coverage

These generic tasks include commands for:
- **Go** - golang projects
- **Python** - Python projects (with common frameworks like Django)
- **JavaScript/TypeScript** - Node.js and frontend projects
- **Java** - Maven and Gradle projects
- **Rust** - Cargo projects
- **Ruby** - Ruby and Rails projects
- **PHP** - Composer and Laravel projects
- **C/C++** - Compiled C/C++ projects

## Adapting Generic Tasks

When creating project-specific tasks from these templates:

1. **Remove the `generic` tag** - Replace with project-specific tags
2. **Update the ID** - Remove `generic-` prefix: `generic-run-unit-tests` → `run-unit-tests`
3. **Specify exact commands** - Choose the specific command for your stack
4. **Add context** - Include project-specific paths, configurations, or notes
5. **Add relationships** - Define prerequisites and downstream tasks for your workflow
6. **Remove irrelevant content** - Delete examples for languages you don't use

## Example: Adapting a Generic Task

**Generic task** (`.tasks/generic/run-unit-tests.yaml`):
```yaml
id: generic-run-unit-tests
name: Run Unit Tests
description: |
  Common commands by language:
  - Go: `go test ./...`
  - Python: `pytest`
  - Node.js: `npm test`
  ...
tags:
  - testing
  - generic
```

**Project-specific task** (`.tasks/run-unit-tests.yaml`):
```yaml
id: run-unit-tests
name: Run Unit Tests
description: |
  Runs all unit tests using pytest.

  Command: `pytest tests/unit/ -v`

  Test files are located in `tests/unit/` directory.
tags:
  - testing
  - python
prerequisites: []
downstream_suggested:
  - run-integration-tests
```

## Contributing

When adding new generic tasks:
- Ensure they apply to **multiple languages/frameworks** (not specific to one)
- Include examples for **at least 3-4 different ecosystems**
- Focus on **common development workflows** most projects need
- Use clear, concise descriptions
- Always include the `generic` tag
