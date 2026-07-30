---
name: explore
description: Fast codebase exploration - find files, patterns, and code structures using LSP-first semantic navigation.
license: MIT
allowed-tools: Glob Grep Read Lsp LspDiagnostics ReferenceSearch TodoWrite Write
metadata:
  agent:
    model: gemini-3-flash-preview
  version: "1.1.0"
---

# Explore

Fast codebase exploration with LSP-first semantic navigation.

## When to Use

Use this skill for:

- Finding where something is defined: "find where X is defined"
- Understanding code flow: "what calls Y", "what does Z do"
- Type hierarchy exploration: "understand the type hierarchy"
- Pattern discovery: "find all implementations of interface X"

## References

For detailed guidance on different exploration approaches, read these references:

| Reference | Content |
|-----------|---------|
| `references/location.md` | Finding WHERE code lives - file location strategies |
| `references/analysis.md` | Understanding HOW code works - implementation analysis |
| `references/patterns.md` | Finding patterns to model new implementations after |

## Output Location

Write exploration findings to the **OUTPUT DIRECTORY specified in context** using the Write tool.

The output directory path is injected when you are spawned as a subagent (check the "## Output Directory" section in your context).

Use the filename pattern: `{adjective}-{adjective}-{noun}.md` (e.g., `bold-keen-beam.md`)

## Document Format

The explore skill produces `analysis` type artifacts (not `explore` - there is no explore artifact type).

```markdown
---
title: "Exploration Title"
type: analysis
workspace: {current workspace path}
created_at: {ISO timestamp}
skill: explore
---

# Exploration: {Title}

## Summary
Brief overview of what was explored and key findings.

## Files Found

### High Relevance
- `path/to/file.ts`: Purpose and key exports

### Medium Relevance
- `path/to/other.ts`: Why it's relevant

## Patterns Found

### Pattern: `pattern-name`
**Files**: `file1.ts`, `file2.ts`

Description of the pattern and how it's used.

```typescript
// Example code snippet
```

## Code Snippets

### `path/to/file.ts:10-25`

Description of what this code does.

```typescript
// Relevant code
```

## Recommendations

- Actionable insights from exploration

```

## Exploration Workflow

### Step 1: Locate Files
Use `Glob` to find files by pattern. Avoid broad searches.

```

Glob(pattern: "src/**/*.ts")

```

### Step 2: Understand Structure (LSP First!)
**IMPORTANT**: Use LSP before reading files to understand structure without consuming context.

```

Lsp(operation: "documentSymbol", filePath: "src/file.ts")

```

This returns all functions, classes, types, and exports without reading the entire file.

### Step 3: Navigate Semantically
Use LSP to trace code flow:

```

Lsp(operation: "goToDefinition", filePath: "file.ts", line: 10, character: 15)
Lsp(operation: "findReferences", filePath: "file.ts", line: 10, character: 15)
Lsp(operation: "hover", filePath: "file.ts", line: 10, character: 15)

```

### Step 4: Read When Needed
Only read files when you need implementation details:

```

Read(file_path: "src/file.ts")

```

### Step 5: Write Findings
Write your exploration document to the OUTPUT DIRECTORY from context:

```

Write(file_path: "{OUTPUT_DIR}/bold-keen-beam.md", content: "...")

```

## Tool Usage Guidelines

### Prefer LSP Over Read
- Use `documentSymbol` to understand file structure
- Use `goToDefinition` to find where things are defined
- Use `findReferences` to find all usages
- Only `Read` when you need the actual implementation

### Avoid Broad Grep
Grep without specific patterns causes massive output. Use targeted patterns:

```

# Good - specific pattern

Grep(pattern: "export function handleAuth", path: "src/")

# Bad - too broad

Grep(pattern: "function", path: "src/")

```

### Track Progress
Use TodoWrite to track exploration progress for complex searches.

## Example Exploration

**Task**: "Find where user authentication is handled"

1. `Glob(pattern: "src/**/*auth*.ts")` - Find auth-related files
2. `Lsp(documentSymbol)` on each file - Understand structure
3. `Lsp(goToDefinition)` on `handleAuth` - Find definition
4. `Lsp(findReferences)` on `handleAuth` - Find all usages
5. `Read` the key files identified
6. `Write` findings to `{OUTPUT_DIR}/bold-keen-beam.md`
