---
title: Complexity Estimation
description: Guidelines for estimating and marking step complexity
---

# Complexity Estimation

Mark each implementation step with a complexity level to set expectations and help prioritize effort.

> **Note**: Examples below use TypeScript for illustration, but complexity markers apply to any language. Adapt to your project's conventions.

## Complexity Levels

| Marker | Typical Duration | Scope | Risk |
|--------|------------------|-------|------|
| `[trivial]` | Minutes | Single line/file | Very Low |
| `[simple]` | 15-30 min | Single file, clear pattern | Low |
| `[moderate]` | 30-120 min | Multi-file, some unknowns | Medium |
| `[complex]` | Hours+ | Architectural, many files | High |

---

## [trivial] - Obvious Changes

**Characteristics:**

- Single line or few lines
- No logic changes
- Well-defined location
- No testing needed beyond typecheck

**Examples:**

```markdown
### Step 1: Update Version [trivial]
**Files**: `package.json`

Change version from "1.0.0" to "1.1.0".

**Verification**: `cat package.json | grep version`
```

```markdown
### Step 2: Fix Typo [trivial]
**Files**: `src/components/Header:15`

Change "Recieve" to "Receive".

**Verification**: `{typecheck command}`
```

```markdown
### Step 3: Add Environment Variable [trivial]
**Files**: `.env.example`

Add `FEATURE_FLAG=true` to example env.

**Verification**: File updated
```

**When NOT trivial:**

- If you need to find where to make the change → `[simple]`
- If change affects behavior → `[simple]` or higher
- If multiple files need same change → `[simple]`

---

## [simple] - Clear Changes

**Characteristics:**

- Single file, well-understood pattern
- Follows existing conventions
- Has clear examples in codebase
- Limited testing needed

**Examples:**

```markdown
### Step 1: Add New Config Option [simple]
**Files**: `src/config/schema`

Add validation schema for new option.

**Verification**: `{typecheck command}`

```

```markdown
### Step 2: Add Utility Function [simple]
**Files**: `src/utils/string`

Add string formatting helper.

**Verification**: `{test command} src/utils/string.test`

```

```markdown
### Step 3: Add Route [simple]
**Files**: `src/routes/index`

Register new endpoint following existing pattern.

**Verification**: API endpoint responds correctly
```

**When NOT simple:**

- Multiple files affected → `[moderate]`
- New pattern introduced → `[moderate]`
- Significant logic required → `[moderate]`

---

## [moderate] - Multi-file or New Patterns

**Characteristics:**

- Changes span multiple files
- Introduces new pattern or component
- Some complexity in logic
- Requires thoughtful testing
- May have edge cases to consider

**Examples:**

```markdown
### Step 1: Add New Service [moderate]
**Files**:
- `src/services/notification-service` (new)
- `src/services/index`
- `src/types/notification` (new)

Create notification service with:
- Send notification method
- Queue management
- Retry logic

**Verification**: `{test command} src/services/notification-service.test`
```

```markdown
### Step 2: Refactor Component [moderate]
**Files**:
- `src/components/UserCard`
- `src/components/UserCard.test`
- `src/hooks/useUser`

Extract user data fetching, update component.

**Verification**: `{test command} src/components/UserCard.test`
```

```markdown
### Step 3: Add Database Table [moderate]
**Files**:
- `migrations/YYYYMMDD_add_notifications.sql`
- `src/db/schema`
- `src/stores/notification-store` (new)

Create notifications table and data access layer.

**Verification**: `{migrate command} && {test command} src/stores/`
```

**When NOT moderate:**

- Architectural decision required → `[complex]`
- Breaking changes to public API → `[complex]`
- Data migration involved → `[complex]`

---

## [complex] - Architectural Changes

**Characteristics:**

- Affects system architecture
- Many files and components
- Requires careful planning
- Has migration/rollback concerns
- High risk of unintended consequences
- May need phased implementation

**Examples:**

```markdown
### Step 1: Replace Auth System [complex]
**Files**: (10+ files across auth, middleware, routes, tests)

Migrate from session-based to JWT authentication:
1. Add JWT utilities
2. Update middleware chain
3. Modify all protected routes
4. Update client token handling
5. Migrate existing sessions
6. Remove session code

**Verification**: Full integration test suite passes
```

```markdown
### Step 2: Database Schema Migration [complex]
**Files**:
- Multiple migration files
- Multiple store files
- Service layer updates
- API changes

Normalize user preferences from JSON column to separate tables:
1. Create new tables
2. Add dual-write to new tables
3. Backfill existing data
4. Switch reads to new tables
5. Remove old column

**Verification**: All data migrated, no regressions
```

```markdown
### Step 3: Extract Microservice [complex]
**Files**: Multiple services, new repository, deployment changes

Extract notification system into separate service:
1. Define service boundary
2. Create API contract
3. Implement service
4. Add client in monolith
5. Migrate traffic
6. Remove old code

**Verification**: Feature parity with monitoring
```

---

## Decision Tree

Use this flowchart to estimate complexity:

```
Is it a single-line change with obvious location?
├── Yes → [trivial]
└── No ↓

Is it contained to a single file with clear pattern?
├── Yes → [simple]
└── No ↓

Does it span 2-5 files with moderate logic?
├── Yes → [moderate]
└── No ↓

Does it involve architecture, migration, or 6+ files?
├── Yes → [complex]
└── No → Reassess: break into smaller steps
```

## Adjusting Estimates

**Increase complexity when:**

- Codebase is unfamiliar
- No existing patterns to follow
- External dependencies involved
- Data migration required
- Breaking changes possible

**Decrease complexity when:**

- Similar change was made recently
- Strong test coverage exists
- Clear documentation available
- Experienced with this area

## When Steps Are Too Complex

If a step is `[complex]`, consider:

1. **Can it be broken into phases?**
   - Database migration → Schema, then code, then cleanup

2. **Can it be feature-flagged?**
   - Roll out incrementally, reduce blast radius

3. **Are there intermediate states?**
   - Dual-write, shadow traffic, gradual migration

A plan with many `[complex]` steps suggests the task should be broken down further or phased over multiple plans.
