---
name: research
description: Deep research combining web search and codebase exploration. Use for investigating patterns, finding documentation, and understanding how to implement features.
license: MIT
allowed-tools: Glob Grep Read Lsp LspDiagnostics WebSearch WebFetch ReferenceSearch TodoWrite Write
metadata:
  agent:
    model: gemini-3-flash-preview
  version: "1.1.0"
---

# Research

Deep research combining web sources and codebase exploration to produce comprehensive research documents.

## When to Use

Use this skill for:

- Researching implementation approaches: "how to implement X"
- Investigating patterns: "what patterns exist for Y"
- Understanding architecture: "understand how Z works"
- Finding documentation: "find docs for API X"
- Comparing approaches: "pros/cons of A vs B"

## References

For detailed guidance on different research approaches, read these references:

| Reference | Content |
|-----------|---------|
| `references/web-research.md` | Conducting effective web research with proper attribution |
| `references/artifact-discovery.md` | Finding existing documents and artifacts |
| `references/artifact-analysis.md` | Extracting high-value insights from documents |

## Output Location

Write research findings to the **OUTPUT DIRECTORY specified in context** using the Write tool.

The output directory path is injected when you are spawned as a subagent (check the "## Output Directory" section in your context).

Use the filename pattern: `<adjective>-<adjective>-<noun>.md` (e.g., `calm-deft-dawn.md`)

## Document Format

```markdown
---
title: "Research Title"
type: research
workspace: {current workspace path}
created_at: {ISO timestamp}
skill: research
---

# Research: {Title}

## Summary
2-5 paragraph executive summary of findings.

## Files Found

### High Relevance
- `path/to/file.ts`: Purpose and key exports
  ```typescript
  // Key code snippet
  ```

### Medium Relevance

- `path/to/other.ts`: Why it's relevant

## Workspace Analysis

### Technology Stack

- TypeScript
- Bun runtime
- SolidJS for UI

### Architecture Patterns

- Registry pattern for skills/tools
- Factory pattern for tool creation
- Event-driven subagent coordination

### File Structure

Overview of how the codebase is organized.

## Organization Patterns

### Relevant Repos

- repo-name: What it contains

### Conventions

- Convention 1: Description
- Convention 2: Description

### Reusable Components

- Component: Where and how to use

## Web Resources

### [Resource Title](url)

**Relevance**: Why this is relevant

**Key Insights**:

- Insight 1
- Insight 2

## Recommendations

- Actionable recommendations based on research

## Constraints

- Known limitations or requirements

## Open Questions

- Questions that need further investigation

```

## Research Workflow

### Phase 1: Codebase Research

1. **Find relevant files**
   ```

   Glob(pattern: "src/**/*auth*.ts")

   ```

2. **Understand structure with LSP**
   ```

   Lsp(operation: "documentSymbol", filePath: "src/auth/index.ts")

   ```

3. **Trace implementations**
   ```

   Lsp(operation: "goToDefinition", ...)
   Lsp(operation: "findReferences", ...)

   ```

4. **Read critical files**
   ```

   Read(file_path: "src/auth/handler.ts")

   ```

### Phase 2: Web Research

1. **Search for documentation**
   ```

   WebSearch(query: "TypeScript authentication best practices 2025")

   ```

2. **Fetch relevant pages**
   ```

   WebFetch(url: "<https://docs.example.com/auth>", prompt: "Extract authentication patterns")

   ```

### Phase 3: Reference Search

Search indexed organization repositories for patterns:

```

ReferenceSearch(query: "authentication middleware patterns", language: "typescript")
ReferenceSearch(query: "JWT token handling", repo: "shared-lib")

```

### Phase 4: Synthesize and Write

Combine findings into a coherent research document (use the OUTPUT_DIR from context):

```

Write(file_path: "{OUTPUT_DIR}/calm-deft-dawn.md", content: "...")

```

## Tool Usage Guidelines

### Codebase Tools
- **Glob**: Find files by pattern
- **Grep**: Search for specific patterns (use targeted queries)
- **Read**: Read file contents
- **Lsp**: Semantic navigation (documentSymbol, goToDefinition, findReferences)
- **LspDiagnostics**: Check for errors/warnings

### Web Tools
- **WebSearch**: Search the web for documentation and examples
- **WebFetch**: Fetch and extract content from URLs

### Organization Tools
- **ReferenceSearch**: Search indexed org repos semantically

### Planning Tools
- **TodoWrite**: Track research progress

## Research Quality Standards

### Good Research
- Comprehensive coverage of the topic
- Multiple sources (codebase + web + org repos)
- Concrete code examples
- Actionable recommendations
- Identified constraints and open questions

### Poor Research
- Surface-level findings
- Single source only
- No code examples
- Vague recommendations
- Ignores constraints

## Example Research Task

**Task**: "Research how to add OAuth2 authentication"

1. **Codebase**: Find existing auth code with Glob/Lsp
2. **Web**: Search for OAuth2 best practices with WebSearch
3. **Org**: Find OAuth2 examples in org repos with ReferenceSearch
4. **Synthesize**: Combine into research doc
5. **Write**: Save to `{OUTPUT_DIR}/eager-fair-flux.md`
