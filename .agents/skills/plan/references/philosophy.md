---
title: Planning Philosophy
description: Core principles for creating effective implementation plans
---

# Planning Philosophy

These principles guide the creation of high-quality implementation plans that lead to successful outcomes.

## Core Principles

### 1. Be Skeptical

**Question vague requirements.** Don't accept "make it faster" or "improve the UX" without clarification. Ask:

- What specific metric defines "faster"?
- What user actions should be easier?
- How will we measure success?

**Identify potential issues early.** Before planning:

- What could go wrong?
- What assumptions am I making?
- What don't I understand yet?

**Verify with code, not assumptions.** Read the actual implementation:

```
❌ BAD: "The auth system probably uses JWT"
✅ GOOD: "Read src/auth/token.ts - confirmed JWT with RS256"
```

### 2. Be Interactive

**Adapt involvement to task complexity.** Don't plan in isolation, but don't over-ask either:

1. **Gather context** → Share understanding → Confirm if ambiguous
2. **Research** → Present options → Get feedback on trade-offs
3. **Outline structure** → Validate approach if multiple valid paths
4. **Write detailed plan** → Iterate → Final review for complex tasks

**Use judgment for user involvement:**

| Task Clarity | User Involvement |
|--------------|------------------|
| Crystal clear requirements | Proceed efficiently, checkpoint at end |
| Some ambiguity | Confirm understanding, validate structure |
| Very ambiguous | Checkpoints at every step |
| User requested collaboration | All checkpoints |

**Allow course corrections.** Plans should evolve:

```
❌ BAD: Writing 50 steps without checking in
✅ GOOD: "Here's the outline. Does this approach make sense before I detail each step?"
```

**Route critical decisions to the user.** Architecture, scope, and requirement decisions should involve the user:

```typescript
AskParentQuestion({
  question: "Should we use Redis or in-memory caching?",
  type: "choice",
  routing: "user",  // Architecture decision - user decides
  options: [
    { label: "Redis", description: "Persistent, shared across instances" },
    { label: "In-memory", description: "Faster, simpler, single instance only" }
  ]
})
```

### 3. Be Thorough

**Read complete files.** Don't skim:

```
❌ BAD: Reading first 50 lines of a 500-line file
✅ GOOD: Reading the entire file to understand all edge cases
```

**Research patterns with parallel agents.** For complex tasks:

```typescript
// Spawn multiple research agents
SpawnAgent({
  agent_type: "explore",
  task: "Find all usages of UserService",
  context_summary: "Looking for UserService patterns"
})

SpawnAgent({
  agent_type: "explore",
  task: "Find similar features to model after",
  context_summary: "Looking for similar implementations"
})
```

**Include specific file paths.** Be precise:

```
❌ BAD: "Update the config file"
✅ GOOD: "Update src/config/server-config.ts lines 45-60"
```

**Write measurable success criteria.** Make verification concrete:

```
❌ BAD: "The feature should work"
✅ GOOD: "{test command} passes with >80% coverage"
```

### 4. Be Practical

**Focus on incremental, testable changes.** Each step should:

- Be independently verifiable
- Not break existing functionality
- Be small enough to review

**Consider migration and rollback.** For each change:

- How do we migrate existing data?
- How do we roll back if something goes wrong?
- What's the blast radius of failure?

**Think through edge cases.** Before finalizing:

- What if the input is empty?
- What if the service is unavailable?
- What if the user lacks permissions?

**Explicitly list out-of-scope items.** Prevent scope creep:

```markdown
## Out of Scope
- Performance optimization (separate task)
- Mobile responsive design (v2)
- Internationalization (not required for MVP)
```

### 5. No Unresolved Questions

**Stop if open questions emerge.** Don't continue planning with unknowns:

```
❌ BAD: Planning continues with "[TBD: decide auth approach]"
✅ GOOD: Pausing to research or ask user for clarification
```

**Research or clarify immediately.** When blocked:

1. Try to find the answer (code search, web search)
2. If still unclear, ask the user
3. Document the decision and rationale

**Final plans must be complete and actionable.** No placeholders:

```
❌ BAD: "Step 3: Implement the main logic [details TBD]"
✅ GOOD: "Step 3: Add UserPreferences class with get/set/reset methods..."
```

## Anti-Patterns to Avoid

### The Assumption Trap

Planning based on what you think the code does instead of reading it.

**Fix**: Always read relevant files before planning changes.

### The Kitchen Sink

Including every possible feature and edge case in v1.

**Fix**: Define MVP scope clearly. Use "Out of Scope" section.

### The Vague Step

Steps like "set up the infrastructure" or "implement the feature."

**Fix**: Break into specific, actionable sub-steps with file paths.

### The Missing Verification

Planning changes without defining how to verify them.

**Fix**: Every step needs a verification command or check.

### The Silent Assumption

Making decisions without documenting why.

**Fix**: Use "Assumptions" section. Ask when uncertain.

### The Linear Path

Assuming everything will work perfectly on first try.

**Fix**: Include error handling, rollback strategies, and risk mitigations.

## When to Research More

Before finalizing a plan, ask yourself:

- [ ] Do I understand the current implementation?
- [ ] Do I know what patterns this codebase uses?
- [ ] Have I identified all affected files?
- [ ] Do I know the testing approach?
- [ ] Are there similar implementations to model after?

If any answer is "no," research more before planning.
