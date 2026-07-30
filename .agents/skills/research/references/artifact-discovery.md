---
title: Artifact Discovery
description: How to find relevant documents and artifacts for a research task
---

# Artifact Discovery

When you need to find existing research, plans, and documents relevant to your task.

## Core Principle

**Locate and categorize relevant artifacts without deep analysis.**

Use the artifact system tools to discover what documents exist and retrieve their content.

## Artifact Types

The system supports these artifact types:

| Type | Description | Created By |
|------|-------------|------------|
| `research` | Deep investigations with web and codebase findings | Research skill |
| `plan` | Implementation plans with steps and file changes | Plan skill |
| `analysis` | Code analysis and exploration results | Explore skill |
| `summary` | Condensed summaries of larger documents | Various skills |
| `code` | Generated code snippets or modules | Implementation work |
| `handoff` | Session transfer documents for continuity | Handoff skill |

**Note**: The `explore` skill produces `analysis` type artifacts, not `explore` type.

## Artifact Storage

Artifacts are stored in:

```
~/.config/csgdaa-code/projects/{projectId}/artifacts/
├── index.json           # Artifact index with metadata
├── {id}.json            # Full artifact content files
└── manifest.json        # Lightweight index for discovery
```

You don't need to access these files directly - use the artifact tools.

## Search Strategy

### Step 1: List Available Artifacts

```
ListArtifacts
```

**Parameters** (all optional):

- `types`: Filter by artifact types, e.g., `["research", "plan"]`
- `tags`: Filter by tags (AND logic - must have all)
- `unconsumedOnly`: Only return artifacts not yet consumed (default: false)
- `limit`: Maximum results (default: 20)
- `format`: Output format - `"references"` (default), `"summary"`, or `"json"`

**Return Format** (references mode):

```
[REF:research:abc123] Auth Middleware - Middleware pattern for request authentication
  Tags: auth, middleware
  Path: ~/.config/csgdaa-code/projects/.../artifacts/abc123.json
```

### Step 2: Filter by Relevance

From the ListArtifacts output:

- Check titles and summaries for keyword matches
- Note artifact types (research, plan, analysis)
- Look for relevant tags
- Check if marked as consumed (may already be incorporated elsewhere)

### Step 3: Retrieve Artifact Content

```
RetrieveArtifact
  artifact_id: "abc123"
```

**Parameters**:

- `artifact_id`: Required - the artifact ID from ListArtifacts
- `markConsumed`: Whether to mark as consumed (default: true)
- `consumerType`: Who's consuming - `"planner"`, `"consolidator"`, or `"main-agent"`

**Returns**: The full artifact content directly (not just a path). For research documents, also returns parsed structured data with sections like:

- `criticalFiles`: Array of `{ path, purpose, exports }`
- `codeExamples`: Array of code snippets
- `implementationHints`: Actionable insights
- `constraints`: Known limitations
- `openQuestions`: Unresolved questions

### Step 4: Categorize Findings

Group documents by purpose:

- **Research** - Background information, investigations
- **Plans** - Implementation approaches, specifications
- **Analysis** - Codebase findings from exploration
- **Handoffs** - Context from previous sessions

## Output Format

```markdown
## Relevant Artifacts for [Topic]

### Research Documents
- `[abc123]` **Auth System Research** - Authentication patterns and best practices
  - Tags: auth, security, jwt
  - Consumed: No
  - Relevance: Covers JWT implementation we need

### Implementation Plans
- `[def456]` **API Refactor Plan** - Steps to modernize REST endpoints
  - Tags: api, refactor
  - Consumed: Yes (by previous planning)
  - Relevance: May have relevant endpoint patterns

### Analysis (from Explore)
- `[ghi789]` **Service Layer Analysis** - How services are structured
  - Tags: services, architecture
  - Relevance: Shows existing patterns to follow

### Handoffs
- `[jkl012]` **Session Handoff 2024-01-15** - Context from previous session
  - Contains: Auth implementation progress, blockers encountered

### Recommended Reading Order
1. Start with: `abc123` - Foundational research on auth
2. Then: `ghi789` - Understand existing service patterns
3. If needed: `def456` - Reference for API patterns
```

## Guidelines

- **Scan metadata first** - Use ListArtifacts before retrieving full content
- **Check consumption status** - Consumed artifacts may be incorporated elsewhere
- **Note artifact types** - Different types serve different purposes
- **Suggest reading order** - Help prioritize what to read first
- **Use tags** - Filter by relevant tags to narrow results

## What NOT to Do

- Don't retrieve all artifacts just to scan them - use ListArtifacts first
- Don't analyze document content deeply during discovery
- Don't evaluate document quality
- Don't mark artifacts consumed unless you're actually using them

## Remember

You're creating a map of existing knowledge. Use the artifact system tools efficiently to identify what's available so the right documents can be retrieved and analyzed.
