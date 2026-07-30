---
title: Phase Template
description: Template for multi-phase implementation plans
---

# Phase Template

Use this template for large implementations that span multiple phases. Each phase should be independently deployable and testable.

## When to Use Phases

- Implementation spans multiple days or sessions
- Changes can be incrementally released
- Rollback points are needed between major changes
- Multiple team members may work on different phases

## Phase Structure

```markdown
## Phase N: {Phase Title}

### Objective
{1-2 sentence description of what this phase accomplishes}

### Prerequisites
{What must be completed before this phase can begin}

- [ ] Phase N-1 completed and verified
- [ ] {Specific prerequisite}
- [ ] {Another prerequisite}

### Steps

#### Step N.1: {Step Title} [complexity]
**Files**: `path/to/file.ts`

{Description}

**Verification**: {Command to verify}

#### Step N.2: {Step Title} [complexity]
...

### Phase Success Criteria

**Automated**:
- [ ] All tests pass
- [ ] No type errors
- [ ] {Phase-specific check}

**Manual**:
- [ ] {Phase-specific manual verification}

### Rollback Strategy
{How to undo this phase if needed}

1. Revert commits: `git revert {range}`
2. {Additional rollback steps}
3. Verify rollback: `{command}`

### Phase Dependencies
{What other phases or systems this phase affects}

- **Depends on**: Phase N-1
- **Blocks**: Phase N+1
- **Affects**: {External system or component}
```

## Example: Database Migration in Phases

```markdown
## Phase 1: Schema Addition (Non-breaking)

### Objective
Add new columns without affecting existing functionality.

### Prerequisites
- [ ] Database backup completed
- [ ] Migration script reviewed

### Steps

#### Step 1.1: Add Migration [simple]
**Files**: `migrations/20250117_add_user_preferences.sql`

```sql
ALTER TABLE users ADD COLUMN preferences JSONB DEFAULT '{}';
```

**Verification**: `bun run db:migrate && bun run db:status`

### Phase Success Criteria

**Automated**: `bun test:db` passes
**Manual**: Existing features work unchanged

### Rollback Strategy

```sql
ALTER TABLE users DROP COLUMN preferences;
```

---

## Phase 2: Code Integration (Feature Flag)

### Objective

Add code to use new schema, behind feature flag.

### Prerequisites

- [ ] Phase 1 completed and deployed
- [ ] Feature flag `USER_PREFERENCES_V2` configured

...

```

## Phase Sizing Guidelines

| Phase Size | Duration | Steps | Risk Level |
|------------|----------|-------|------------|
| Small | 1-2 hours | 2-4 | Low |
| Medium | 2-4 hours | 4-8 | Medium |
| Large | 4-8 hours | 8-15 | High |

**Rule of thumb**: If a phase has more than 15 steps, split it.

## Phase Boundaries

Good phase boundaries are:

- **Deployable**: Phase can be released independently
- **Testable**: Phase has clear success criteria
- **Rollbackable**: Phase can be undone without data loss
- **Reviewable**: Phase can be code-reviewed as a unit

Poor phase boundaries:

- Partial feature that breaks without next phase
- Changes that can't be tested in isolation
- Irreversible changes mixed with reversible ones
