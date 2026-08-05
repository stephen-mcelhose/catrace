---
name: review
description: >
  Review code or docs changes. Default is normal review (correctness, clarity,
  conventions). Use GAN/adversarial mode when the user asks for GAN review or
  when reviewing merge-critical code (package-root Go, experiments/*/main.go,
  merge-ready runnable PRs). Also triggers on "review changes", "check a PR",
  "review before committing".
license: MIT
allowed-tools:
  # Version control (read-only)
  - Bash(git:*)
  - Bash(gh pr:*)
  - Bash(gh:* --jq:*)
  - Bash(jq:*)
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
  version: "3.1.0"
---

# Code Review

When inspecting `gh`/API JSON in the shell, **prefer `jq` / `gh --jq`**.

## Choose mode (do this first)

| Mode | When |
|------|------|
| **Normal** (default) | Docs, plans, wiki, example polish, general “review this” |
| **GAN** (adversarial) | User asks for GAN / adversarial review; **or** diff is merge-critical: package-root `*.go`, `experiments/*/main.go`, or a PR about to merge runnable behavior |

State the chosen mode in one line at the start of the review.

- **Normal:** correctness, clarity, conventions, maintainability. Not “try to break it.”
- **GAN:** assume the implementation is wrong until proven otherwise; prefer concrete failures; use attacks table. Walkthroughs are not GAN (peer-code-review only if asked).

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

3. Always read security guidelines:
   - [references/security_review.md](./references/security_review.md)

4. If mode is **GAN**, also read:
   - [references/gan_review.md](./references/gan_review.md)

> **Note:** Load language + security refs before analysis. Load GAN ref only in GAN mode.

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

### 8a. Normal mode — structured feedback

> File refs: `[filename:line](path/to/file:line)`.

```markdown
## Mode: normal

## Summary
Brief overview and assessment.

## Issues
- **[file:line](path) SEVERITY** Description and suggested fix

## Suggestions
- Optional improvements

## Verdict: Approved / Changes Requested
```

Severity for normal mode: **CRITICAL** / **HIGH** / **MEDIUM** / **LOW**.

### 8b. GAN mode — attack then report

Before writing, try to break the diff:

1. List MUST/SHOULD behaviors (tests, walkthrough, PR body).
2. Invent inputs/sequences that should fail if the code is weak.
3. Prefer a failing test or exact repro.
4. Classify: `blocker` / `major` / `minor` / `note`.

Domain angles: non-stochastic rows, joint index mismatches, Kronecker where
coupling was claimed, Stationary on non-ergodic chains, Trace with recurrent
hidden states, MFPT convention vs docs, silent ignored errors.

```markdown
# GAN review

## Scope
- Base / head:
- Files:

## Attacks
| # | Attack | Result | Severity |
|---|--------|--------|----------|
| 1 | … | reproduced / not reproducible | blocker |

## Required fixes before ship
- …

## Residual risk
- …

## Verdict: ship / changes required
```

Do **not** soft-pass blockers. Style-only nits are `note`.

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
