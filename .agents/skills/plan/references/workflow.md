---
title: Interactive Planning Workflow
description: 5-phase process for creating implementation plans interactively
---

# Interactive Planning Workflow

This workflow creates high-quality implementation plans through iterative collaboration rather than isolated planning.

**Remember**: You are in READ-ONLY mode for the codebase. Only write plan documents to the output directory.

## Overview

```
┌─────────────────────────────────────────────────────────────┐
│  Phase 1: Initial Understanding  →  Understand the task    │
│  Phase 2: Research & Discovery   →  Explore codebase       │
│  Phase 3: Plan Structure         →  Outline, get feedback  │
│  Phase 4: Detailed Plan          →  Write complete plan    │
│  Phase 5: Review & Completion    →  Refine until approved  │
└─────────────────────────────────────────────────────────────┘
```

---

## Phase 1: Initial Understanding

**Goal**: Deeply understand the task before planning.

### Check for Prior Research

Before exploring from scratch, check if research artifacts already exist:

```typescript
// Check for existing research documents
ListArtifacts({ type: "research" })

// Retrieve relevant findings
RetrieveArtifact({ id: "artifact-id-from-list" })
```

Prior research from `explore` or `research` skills may already contain the codebase patterns, file inventories, and architectural context you need. Build on existing work rather than duplicating it.

### Explore the Codebase

Fill gaps not covered by prior research:

```typescript
// Read mentioned files completely — don't skim
Read({ file_path: "path/to/relevant/file" })

// Search for related code
Grep({ pattern: "ServiceName", output_mode: "files_with_matches" })
Glob({ pattern: "src/**/*feature*" })

// Understand structure
Lsp({ operation: "documentSymbol", filePath: "path/to/module" })
```

### Present Understanding

Summarize what you've learned:

```markdown
## My Understanding

**Task**: {what the user wants}

**Current State**:
- {relevant existing code/infrastructure}
- {patterns already in use}

**Proposed Approach**:
- {high-level strategy}

**Questions**:
1. {anything ambiguous}
```

### USER CHECKPOINT (Conditional)

**Use when**: Requirements are vague, ambiguous, or you made assumptions.

```typescript
AskParentQuestion({
  question: "Does my understanding look correct? Any clarifications?",
  type: "approval",
  routing: "user",
  context: "Confirming I understand your requirements before proceeding."
})
```

**Skip if**: Requirements are crystal clear and you made no assumptions.

---

## Phase 2: Research & Discovery

**Goal**: Gather detailed information using parallel exploration.

**Skip if**: Task is trivial, or Phase 1 gave you complete understanding.

### Spawn Parallel Research Agents

For complex tasks, spawn multiple agents to explore different aspects simultaneously:

```typescript
// Agent 1: Find existing patterns
SpawnAgent({
  agent_type: "explore",
  task: "Find how similar features are implemented. Look at stores, services, and routes for patterns.",
  context_summary: "Researching implementation patterns for the planned feature"
})

// Agent 2: Check testing conventions
SpawnAgent({
  agent_type: "explore",
  task: "Find test patterns — what framework, what directories, what coverage expectations.",
  context_summary: "Researching testing conventions"
})
```

Scale agent count to task complexity:

| Scenario | Agents |
|----------|--------|
| Isolated to known files | 1 |
| Multiple areas involved | 2-3 |
| Complex architectural task | 3+ |

### When to Use Web Search

Use when looking for best practices, library APIs, or security considerations:

```typescript
WebSearch({ query: "database migration best practices rollback strategy" })
```

### Synthesize Findings

After all research completes, present options with trade-offs:

```markdown
## Research Findings

**Existing Patterns**:
- {what you found in the codebase}

**Options Identified**:

**Option A**: {approach}
- Pros: ...
- Cons: ...

**Option B**: {approach}
- Pros: ...
- Cons: ...

**Recommendation**: {which option and why}
```

---

## Phase 3: Plan Structure Development

**Goal**: Get user buy-in on structure before detailed writing.

### Create Outline

Before writing detailed steps, share the structure:

```markdown
## Proposed Plan Structure

### Phase 1: {Layer/Area}
1. {step}
2. {step}
3. {step}

### Phase 2: {Layer/Area}
4. {step}
5. {step}

Does this breakdown make sense? Should any phases be reordered or split?
```

### Get Feedback

Ask targeted questions for scope decisions:

```typescript
AskParentQuestion({
  question: "Should we implement all phases, or start with Phase 1 only?",
  type: "choice",
  routing: "user",
  options: [
    { label: "All phases", description: "Complete implementation" },
    { label: "Phase 1 only", description: "Start small, iterate" }
  ]
})
```

### USER CHECKPOINT (Conditional)

**Use when**: Multiple valid approaches exist, or task is moderate/complex.

**Skip if**: There's only one obvious approach and the task is straightforward.

---

## Phase 4: Detailed Plan Writing

**Goal**: Write a complete, actionable plan.

### Use Templates

Load the plan template:

```typescript
Read({ file_path: "{baseDir}/assets/plan-template.md" })
```

For large, multi-phase implementations (spans multiple days, incremental releases, rollback points needed), load the phase template instead:

```typescript
Read({ file_path: "{baseDir}/assets/phase-template.md" })
```

### Required Sections

Every plan must include:

1. **Summary** — What, why, and key decisions
2. **Files to Modify** — Every file listed with what changes
3. **Implementation Steps** — With `[complexity]` markers and verification per step
4. **Success Criteria** — Automated (runnable commands) AND manual checks
5. **Testing Strategy** — Unit, integration, manual
6. **Risks** — With mitigations
7. **Out of Scope** — Prevent scope creep
8. **Assumptions** — Document decisions made

### Write to Output Directory

```typescript
Write({
  file_path: "{OUTPUT_DIR}/swift-bright-arch.md",
  content: planContent
})
```

### Quality Checklist

Before presenting:

- [ ] Every step has a verification command
- [ ] All affected files are listed
- [ ] Complexity markers assigned to each step
- [ ] Success criteria are measurable and runnable
- [ ] No unresolved questions remain
- [ ] Out-of-scope items documented

---

## Phase 5: Review & Completion

**Goal**: Refine plan until user approves.

### Present Draft

```markdown
## Plan Complete

I've written the implementation plan to:
`{OUTPUT_DIR}/swift-bright-arch.md`

**Summary**:
- {N} implementation steps across {N} phases
- Estimated complexity: {overall}

Please review and let me know if you'd like any changes.
```

### Iterate Based on Feedback

| Feedback | Action |
|----------|--------|
| "Too detailed" | Combine steps, reduce examples |
| "Missing X" | Add section, research if needed |
| "Wrong approach" | Discuss alternative, rewrite |
| "Step is unclear" | Add specifics, include code examples |

### Completing Your Turn

**Your turn MUST end with one of:**

1. **AskParentQuestion** with `routing: "user"` — for clarifying scope/architecture/requirements
2. **Presenting the completed plan** — write document, present summary, await approval

**Do NOT** end without one of these actions, ask "Is this plan okay?" in plain text, or leave unresolved questions.

---

## Question Routing

```
Is this a technical question the codebase can answer?
├── YES → routing: "parent"
│         (e.g., "What testing framework is used?")
└── NO → Does this affect scope, architecture, or requirements?
         ├── YES → routing: "user"
         │         (e.g., "Should we prioritize performance or simplicity?")
         └── NO → routing: "parent-or-user"
                   (e.g., "What's the naming convention?")
```

## Completion Checklist

- [ ] Prior research consumed (ListArtifacts/RetrieveArtifact)
- [ ] Codebase patterns researched
- [ ] User's requirements confirmed
- [ ] Plan structure reviewed by user
- [ ] All steps have verification
- [ ] Success criteria are measurable
- [ ] Risks identified with mitigations
- [ ] Out-of-scope items documented
- [ ] No open questions remain
- [ ] Plan written to output directory
