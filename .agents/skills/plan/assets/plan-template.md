---
title: Plan Template
description: Complete implementation plan document structure
---

# Plan Template

Use this template when writing implementation plans. Copy and adapt to your specific task.

## YAML Frontmatter

```yaml
---
title: "{Descriptive Plan Title}"
type: plan
workspace: {absolute workspace path}
created_at: {ISO 8601 timestamp}
skill: plan
sources:
  - "{research artifact ID or path}"
  - "{additional sources}"
---
```

## Document Structure

```markdown
# Implementation Plan: {Title}

## Summary

{2-3 paragraphs describing:
- What this plan implements
- The high-level approach
- Key decisions made and why}

## Files to Modify

### High Relevance (Must Change)
{Files that will definitely be modified}

- `path/to/file`: {What changes and why}
- `path/to/another`: {What changes and why}

### Medium Relevance (May Change)
{Files that might need changes depending on implementation}

- `path/to/maybe`: {Potential change}

### New Files
{Any new files to be created}

- `path/to/new/file`: {Purpose of new file}

## Implementation Steps

### Step 1: {Step Title} [complexity]
**Files**: `path/to/file`

{Detailed description of what to do}

```
// Code example showing the change
```

**Verification**: {Command or check to verify this step worked}

### Step 2: {Step Title} [complexity]

...

## Success Criteria

### Automated Verification

{Commands that can be run by agents to verify success — use the project's actual test/build commands}

- [ ] Type checking passes: `{typecheck command}`
- [ ] Tests pass: `{test command} path/to/tests`
- [ ] Build succeeds: `{build command}`
- [ ] Linting passes: `{lint command}`

### Manual Verification

{Checks requiring human testing}

- [ ] Feature works as expected
- [ ] Performance is acceptable under load
- [ ] Edge cases handled correctly
- [ ] No regressions in related functionality

## Testing Strategy

### Unit Tests

{What unit tests to add or update}

- Test case 1: {description}
- Test case 2: {description}

### Integration Tests

{Integration testing approach}

- Scenario 1: {description}
- Scenario 2: {description}

### Manual Testing

{Steps for manual verification}

```bash
# Start the application
{dev command}

# Test the feature
{manual test steps}
```

## Risks

### {Risk Title} [severity: high|medium|low]

**Issue**: {What could go wrong}
**Mitigation**: {How to prevent or handle it}

### {Another Risk} [severity]

**Issue**: ...
**Mitigation**: ...

## Out of Scope

{Explicitly list what this plan does NOT cover to prevent scope creep}

- {Feature or concern explicitly deferred}
- {Optimization or enhancement for a future task}

## Assumptions

{List assumptions made during planning}

- Assumption 1: {description}
- Assumption 2: {description}

## Open Questions

{Questions that need answers before or during implementation}

- [ ] Question 1: {description}
- [ ] Question 2: {description}

## Constraints

{Technical or business constraints affecting implementation}

- Constraint 1: {description}
- Constraint 2: {description}

```

## Complexity Markers

Use these markers for each step:

| Marker | Description | Typical Scope |
|--------|-------------|---------------|
| `[trivial]` | Single-line change, obvious fix | Config update, typo fix |
| `[simple]` | Small, well-understood change | Single-file, clear pattern |
| `[moderate]` | Multi-file, some complexity | New feature, pattern change |
| `[complex]` | Architectural, careful needed | Migration, refactor, new system |

## Tips for Good Plans

1. **Every file touched should be listed** - No surprises during implementation
2. **Verification for each step** - Catch issues early, not at the end
3. **Code examples for complex changes** - Show what the result looks like
4. **Separate automated vs manual criteria** - Know what can be CI'd
5. **Risks with mitigations** - Show you've thought through failure modes
