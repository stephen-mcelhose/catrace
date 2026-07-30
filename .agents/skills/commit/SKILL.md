---
name: commit
description: Create git commits with well-crafted messages. Use when the user says "commit", "save changes", "create a commit", or after completing a task that modified files.
license: MIT
allowed-tools: Bash(git status:*) Bash(git diff:*) Bash(git log:*) Bash(git add:*) Bash(git commit:*) Read Glob Grep
metadata:
  version: "1.3.0"
  hooks:
    - event: PreToolUse
      matcher: "Bash"
      type: command
      command: "jq -r '.tool_input.command // \"\"' | grep -qE '\\bgit\\b.*\\bcommit\\b' && { [ ! -f \"$HOOK_WORKSPACE/.gitignore\" ] && echo 'WARNING: No .gitignore found at repo root. Build artifacts, secrets, or environment files could be committed by accident. Consider adding a .gitignore before committing.' || true; } || true"
      timeout: 5000
---

# Git Commit

Create a commit with a clear, descriptive message based on staged changes.

## Workflow

1. Run `git status` to see staged and unstaged changes
2. **Pre-commit safety check**: If no `.gitignore` exists at the repo root, warn the user before proceeding — build artifacts, secrets, and environment files can easily be committed by accident. Also check if any staged files match common artifact patterns (`.build/`, `node_modules/`, `dist/`, `DerivedData/`, `__pycache__/`, `.venv/`, `.env`) and warn if so. These are warnings, not blockers — some repos legitimately have no `.gitignore`.
3. Run `git diff --staged` to review what will be committed
4. Run `git log --oneline -5` to match the repository's commit style
5. Draft a message explaining the **why**, not just the what, be brief, but descriptive
6. Create the commit using a heredoc for proper formatting

## Commit Message Format

```text
<type>(<scope>): <subject>

<body>
```

### Types

**Version-bumping (end-user-facing changes only):**

- `feat` — New feature visible to end users of the deployed application
- `fix` — Bug fix affecting end users of the deployed application

Write these from the user's perspective: what benefit do they see?

**Non-bumping (developer/internal changes):**

- `chore` — Maintenance tasks, dependency updates, configs
- `ci` — CI/CD pipeline changes
- `refactor` — Code restructuring without behavior change
- `docs` — Documentation only
- `style` — Formatting, whitespace (no code change)
- `test` — Adding or updating tests

Use these for skills, tooling, internal refactors, and infrastructure.

### Examples

| ✅ Do | ❌ Don't |
| ----- | -------- |
| `feat(checkout): complete purchases without leaving the page` | `feat(checkout): add AJAX form submission` |
| `fix(upload): files over 10MB no longer fail silently` | `fix(upload): handle chunked transfer edge case` |
| `chore: update eslint to v9` | `feat: update eslint to v9` |
| `ci: run tests in parallel to speed up builds` | `feat(ci): add parallel test jobs` |
| `refactor(auth): simplify token validation logic` | `feat(auth): refactor token validation` |
| `docs: add deployment guide for new contributors` | `feat(docs): add deployment guide` |
| `chore(infra): remove legacy terraform variables` | `fix(infra): remove old vars from terraform` |

**Key distinction:** If the change doesn't affect what end users experience in the deployed app, it's not a `feat` or `fix`.

### Rules

- Always sign commits (git commit -S)
- Never use emoji in commit subject or body
- Subject: 50 chars max, imperative mood ("Add" not "Added")
- Body: Explain what and why, reference issues if applicable, and line wrap at 72 characters
- Be brief, but descriptive
- Use heredoc for multi-line messages:

```bash
git commit -m "$(cat <<'EOF'
fix(auth): resolve token refresh race condition

The refresh token was being invalidated before the new token
was stored, causing intermittent 401 errors.

Fixes #123
EOF
)"
```

## Safety Rules

### NEVER use these operations (blocked by permissions)

| Command | Risk |
| --------- | ------ |
| `git push --force` / `-f` | Destroys remote history, breaks collaborators |
| `git reset --hard` | Permanently loses uncommitted work |
| `git clean -f` | Permanently deletes untracked files |
| `git branch -D` | Force-deletes branch ignoring merge status |
| `git rebase` | Rewrites history, dangerous on shared branches |
| `git stash drop/clear` | Permanently loses stashed changes |
| `git config --global/--system` | Could break git installation |

### Commit safety

- NEVER skip hooks (`--no-verify`) unless explicitly requested
- NEVER amend pushed commits without explicit permission
- NEVER create empty commits
- Always verify `git status` before committing
- Always review `git diff --staged` before committing

### Amend rules

Only use `git commit --amend` when ALL conditions are met:

1. User explicitly requested amend, OR pre-commit hook auto-modified files
2. The commit has NOT been pushed to remote (`git status` shows "ahead")
3. You created the commit in this session

If unsure, create a NEW commit instead of amending.
