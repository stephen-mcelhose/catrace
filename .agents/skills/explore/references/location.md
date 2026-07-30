---
title: Code Location
description: How to find WHERE code lives in a codebase
---

# Code Location

When you need to find WHERE code lives, not analyze what it does.

## Core Principle

**Find and categorize files by purpose. Don't read their contents.**

Your job is to locate relevant files and organize them by purpose. This is faster than analysis because you're mapping territory, not studying it.

## Search Strategy

### Initial Broad Search

Think deeply about effective search patterns:

- Common naming conventions in this codebase
- Language-specific directory structures
- Related terms and synonyms

**Find files containing a pattern:**

```
Grep
  pattern: "UserService"
  output_mode: "files_with_matches"
```

**Grep Parameters:**

- `pattern`: Regex pattern to search for (required)
- `path`: Directory to search (optional, defaults to workspace)
- `output_mode`: `"files_with_matches"` (default), `"content"`, or `"count"`
- `glob`: Filter files by glob pattern, e.g., `"*.ts"`
- `type`: Filter by file type, e.g., `"ts"`, `"py"`, `"js"`
- `case_insensitive`: Boolean for case-insensitive search
- `context_lines`: Lines of context around matches (content mode)
- `head_limit`: Limit to first N results

**Find files by name pattern:**

```
Glob
  pattern: "src/**/*auth*.ts"
```

**Glob Parameters:**

- `pattern`: Glob pattern (required), e.g., `"**/*.ts"`, `"src/**/test_*.py"`
- `path`: Directory to search (optional, defaults to workspace)
- `limit`: Maximum results (default: 100)

Results are sorted by modification time (newest first).

### Refine by Language/Framework

**JavaScript/TypeScript**: `src/`, `lib/`, `components/`, `pages/`, `api/`
**Python**: `src/`, `lib/`, `pkg/`, module names
**Go**: `pkg/`, `internal/`, `cmd/`

### Common Patterns to Find

- `*service*`, `*handler*`, `*controller*` - Business logic
- `*test*`, `*spec*` - Test files
- `*.config.*`, `*rc*` - Configuration
- `*.d.ts`, `*.types.*` - Type definitions
- `README*`, `*.md` in feature dirs - Documentation

## Output Format

```markdown
## File Locations for [Feature/Topic]

### Implementation Files
- `src/services/feature.ts` - Main service logic
- `src/handlers/feature-handler.ts` - Request handling
- `src/models/feature.ts` - Data models

### Test Files
- `src/services/__tests__/feature.test.ts` - Service tests
- `e2e/feature.spec.ts` - End-to-end tests

### Configuration
- `config/feature.json` - Feature-specific config

### Type Definitions
- `types/feature.d.ts` - TypeScript definitions

### Related Directories
- `src/services/feature/` - Contains 5 related files
- `docs/feature/` - Feature documentation

### Entry Points
- `src/index.ts` - Imports feature module at line 23
- `api/routes.ts` - Registers feature routes
```

## Guidelines

- **Don't read file contents** - Just report locations
- **Be thorough** - Check multiple naming patterns
- **Group logically** - Make it easy to understand code organization
- **Include counts** - "Contains X files" for directories
- **Note naming patterns** - Help understand conventions
- **Check multiple extensions** - `.js/.ts`, `.py`, `.go`, etc.

## What NOT to Do

- Don't analyze what the code does
- Don't read files to understand implementation
- Don't make assumptions about functionality
- Don't skip test or config files
- Don't critique file organization
- Don't suggest better structures

## Remember

You're a file finder and organizer. Create a map of the existing territory, don't redesign the landscape.
