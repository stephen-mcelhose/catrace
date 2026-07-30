---
name: review
description: Review code changes for quality, bugs, and improvements. Use when the user asks to review changes, check a PR, or review code before committing.
license: MIT
allowed-tools:
  # Version control (read-only)
  - Bash(git:*)
  - Bash(gh pr:*)
  # File tools (always needed)
  - Read
  - Glob
  - Grep
  - Lsp
  - LspDiagnostics
  - ReferenceSearch
  # Skill chaining (delegates linting to /lint)
  - Skill
  # Tool detection
  - Bash(command -v:*)
  - Bash(which:*)
  # Tool installation (common)
  - Bash(brew install:*)
  # Note: Language-specific tools are added when you read the corresponding
  # reference file (e.g., go_review.md adds Go tools, typescript_review.md adds TS tools)
metadata:
  version: "2.5.0"
---

# Code Review

Review code changes systematically for quality, correctness, security, and maintainability.

## Workflow

### 1. Get Repository Context

```bash
# Check current branch and status
git status
git branch -v

# View recent commits
git log --oneline -10
```

### 2. Get the Diff

```bash
# Unstaged changes
git diff

# Staged changes
git diff --staged

# Branch comparison (common patterns)
git diff main...HEAD
git diff origin/main...HEAD
git diff master...HEAD

# PR diff (if reviewing a PR)
gh pr diff <number>
```

### 3. Detect Languages & Load Review Guidelines

**This step is critical - load guidelines BEFORE analysis to inform your review.**

1. Scan the diff for file extensions to detect languages:
   - `.ts`, `.tsx`, `.js`, `.jsx` → TypeScript/JavaScript
   - `.go` → Go
   - `.py` → Python
   - `.sh`, `.bash` → Shell
   - `.sql` → SQL
   - `.proto` → API/Protobuf
   - `.yaml`, `.yml` in `.github/workflows/` → GitHub Actions

2. Read the relevant reference files for detected languages:

| Language             | Reference File                                                           |
|----------------------|--------------------------------------------------------------------------|
| Go                   | [references/go_review.md](./references/go_review.md)                     |
| Python               | [references/python_review.md](./references/python_review.md)             |
| TypeScript           | [references/typescript_review.md](./references/typescript_review.md)     |
| Shell                | [references/shell_review.md](./references/shell_review.md)               |
| SQL (BigQuery/Spanner) | [references/sql_review.md](./references/sql_review.md)                 |
| GitHub Actions       | [references/github_actions_review.md](./references/github_actions_review.md) |
| API/Protobuf         | [references/api_review.md](./references/api_review.md)                   |

3. Always read security guidelines regardless of language:
   - [references/security_review.md](./references/security_review.md)

> **Note:** These reference files contain critical review criteria. Load them NOW before proceeding to static analysis.

### 4. Run Static Analysis

Delegate linting to the `/lint` skill — it handles tool installation and
runs the appropriate linters for each detected language:

```
/lint
```

For languages not yet covered by `/lint`, or project-specific linting:

| Language       | Commands to Run                                                |
|----------------|----------------------------------------------------------------|
| TypeScript/JS  | `bun run typecheck`, `bun run lint`                            |
| Python         | `ruff check .`, `mypy .`                                       |
| Shell          | `shellcheck *.sh`, `shfmt -d .`                                |
| Rust           | `cargo check`, `cargo clippy`                                  |

**Note:** Use the project's package manager (`bun run`, `npm run`, `pnpm run`, `yarn run`) for project-configured linting when available.

### 5. Check for Reusable Components

**For Go code**, search csgda-kit for existing modules:

```
ReferenceSearch(query="HTTP client", repo="csgda-kit")
ReferenceSearch(query="logging middleware", repo="csgda-kit")
```

**For GitHub Actions**, search csgda-platform-actions:

```
ReferenceSearch(query="deploy workflow", repo="csgda-platform-actions")
ReferenceSearch(query="docker build", repo="csgda-platform-actions")
```

### 6. Understand Context with LSP

Use LSP for deeper code understanding:

```
Lsp(operation: "goToDefinition", filePath: "...", line: N, character: M)
Lsp(operation: "findReferences", filePath: "...", line: N, character: M)
Lsp(operation: "hover", filePath: "...", line: N, character: M)
LspDiagnostics(path: "src/file.ts")
```

### 7. Apply Review Checklist

#### Correctness

- [ ] Logic handles edge cases
- [ ] Error conditions handled properly
- [ ] No silent failures

#### Security

- [ ] Input validation present
- [ ] No hardcoded secrets
- [ ] Injection risks addressed (SQL, command, XSS)
- [ ] Proper authentication/authorization

#### Performance

- [ ] No N+1 queries or unnecessary loops
- [ ] Reasonable allocations
- [ ] Appropriate data structures

#### Maintainability

- [ ] Clear naming and structure
- [ ] Appropriate comments for complex logic
- [ ] No dead code
- [ ] Follows codebase patterns

#### Testing

- [ ] New code paths tested
- [ ] Edge cases covered
- [ ] Tests are readable

### 8. Provide Structured Feedback

## Output Format

> **Important:** When referencing files in your response, always use clickable file annotations:
>
> - Format: `[filename:line](path/to/file:line)`
> - Example: `[user.ts:42](src/models/user.ts:42)`
> - This allows users to click directly to the referenced location

```markdown
## Summary
Brief overview of the changes and overall assessment.

## Analysis Results
Output from static analysis tools (if run).

## Issues
- **[filename:line](path/to/file:line) SEVERITY** Description and suggested fix

Example:
- **[auth.ts:127](src/auth/auth.ts:127) HIGH** Missing null check before accessing user.email

## Suggestions
- **[filename:line](path/to/file:line)** Optional improvement ideas

## Pattern Compliance
- [ ] Follows csgda-kit patterns (for Go)
- [ ] Uses csgda-platform-actions (for CI/CD)
- [ ] Follows language best practices

## Verdict: Approved / Changes Requested
Final decision with reasoning.
```

## Severity Levels

| Level        | Description                                          |
|--------------|------------------------------------------------------|
| **CRITICAL** | Security vulnerability, data loss risk, crash        |
| **HIGH**     | Bug, significant performance issue, breaking change  |
| **MEDIUM**   | Code smell, minor performance issue, maintainability |
| **LOW**      | Style, naming, optional improvements                 |

## Git Commands (Read-Only)

This skill has read-only git access:

- `git status` / `git branch` - Repository state
- `git diff` / `git diff --staged` - View local changes
- `git diff main...HEAD` - Compare branch to main
- `git log` - View commit history
- `git show <commit>` - Inspect specific commits
- `gh pr diff` / `gh pr view` - Read-only PR access

## Quick Reference: Common Issues by Language

### Go

- Ignoring errors (`_`)
- Goroutine leaks (no exit condition)
- Defer in loops
- Race conditions

### Python

- Bare `except:` clauses
- Mutable default arguments
- Missing type hints
- Resource leaks (no `with`)

### TypeScript

- Using `any` type
- Missing null checks
- Unhandled promise rejections
- Type assertions without validation

### Shell

- Unquoted variables
- Missing `set -euo pipefail`
- Using `eval` with user input
- Missing error handling

### GitHub Actions

- Unpinned action versions
- Secrets in logs
- Command injection via `${{ }}`
- Missing permission restrictions

### SQL (BigQuery/Spanner)

- SELECT * (cost/performance)
- Missing partition filters (BigQuery)
- Sequential primary keys (Spanner hotspots)
- Non-parameterized queries
