---
title: Implementation Patterns
description: Common implementation patterns with step-by-step breakdowns
---

# Implementation Patterns

Common patterns for structuring implementation plans. Each pattern provides a proven sequence of steps.

> **Note**: Examples below use TypeScript for illustration, but these patterns apply to any language. Adapt file extensions, commands, and idioms to your project.

## Pattern Categories

- [Database Changes](#database-changes)
- [New Features](#new-features)
- [Refactoring](#refactoring)
- [Bug Fixes](#bug-fixes)
- [API Design](#api-design)

---

## Database Changes

**Sequence**: Schema → Store Methods → Business Logic → API → Clients

This pattern minimizes risk by making changes in dependency order.

### Phase 1: Schema Changes (Non-breaking)

```markdown
### Step 1: Add Migration [simple]
**Files**: `migrations/YYYYMMDD_description.sql`

Add new columns/tables without removing existing ones:

```sql
-- Always additive first
ALTER TABLE users ADD COLUMN preferences JSONB DEFAULT '{}';
```

**Verification**: `{migrate command} && {db status command}`

### Step 2: Update Types [simple]

**Files**: `src/db/schema.ts`

Add TypeScript types for new columns.

**Verification**: `{typecheck command}`

```

### Phase 2: Data Access Layer

```markdown
### Step 3: Add Store Methods [moderate]
**Files**: `src/stores/user-store.ts`

Add methods to read/write new data:

```typescript
async getPreferences(userId: string): Promise<UserPreferences> {
  // ...
}

async setPreferences(userId: string, prefs: UserPreferences): Promise<void> {
  // ...
}
```

**Verification**: `{test command} src/stores/user-store.test`

```

### Phase 3: Business Logic & API

```markdown
### Step 4: Add Service Methods [moderate]
**Files**: `src/services/user-service.ts`

Add business logic using store methods.

### Step 5: Add API Endpoints [simple]
**Files**: `src/routes/users.ts`

Expose new functionality via API.

**Verification**: API endpoint responds correctly
```

### Phase 4: Cleanup (If Replacing)

```markdown
### Step 6: Migrate Data [complex]
**Files**: `migrations/YYYYMMDD_migrate_data.sql`

Move data from old structure to new.

### Step 7: Remove Old Schema [simple]
**Files**: `migrations/YYYYMMDD_cleanup.sql`

Remove deprecated columns after confirming migration.
```

---

## New Features

**Sequence**: Research Patterns → Data Model → Backend Logic → API Endpoints → UI

### Phase 1: Research & Design

```markdown
### Step 1: Research Existing Patterns [trivial]
Search codebase for similar features:

- How do other features handle this?
- What patterns are already established?
- Are there utilities to reuse?

**Verification**: Document findings in plan

### Step 2: Design Data Model [simple]
**Files**: `src/types/feature.ts`

Define types and interfaces:

```typescript
interface FeatureConfig {
  enabled: boolean;
  settings: FeatureSettings;
}
```

**Verification**: Types compile

```

### Phase 2: Backend Implementation

```markdown
### Step 3: Implement Core Logic [moderate]
**Files**: `src/services/feature-service.ts`

Core business logic without external dependencies first.

### Step 4: Add Persistence [moderate]
**Files**: `src/stores/feature-store.ts`

Storage layer for feature data.

### Step 5: Add API Routes [simple]
**Files**: `src/routes/feature.ts`

REST/GraphQL endpoints.

**Verification**: API responds correctly
```

### Phase 3: UI Integration

```markdown
### Step 6: Add UI Components [moderate]
**Files**: `src/components/Feature/*.tsx`

User interface components.

### Step 7: Wire Up to API [simple]
**Files**: `src/hooks/useFeature.ts`

Connect UI to backend.

**Verification**: Feature works end-to-end
```

---

## Refactoring

**Sequence**: Document Current Behavior → Plan Incremental Changes → Maintain Compatibility

### Step 1: Document Current Behavior

```markdown
### Step 1: Capture Current Behavior [trivial]
**Files**: `docs/refactor-notes.md` (temporary)

Document:
- All public APIs and their contracts
- All call sites
- All tests and what they verify

**Verification**: Notes reviewed
```

### Step 2: Add Tests If Missing

```markdown
### Step 2: Add Characterization Tests [moderate]
**Files**: `src/module/__tests__/module.test.ts`

Capture current behavior in tests before changing:

```typescript
describe('Module (before refactor)', () => {
  it('handles edge case X', () => {
    // Capture current behavior
  });
});
```

**Verification**: Tests pass, coverage adequate

```

### Step 3: Incremental Changes

```markdown
### Step 3: Extract/Rename/Move [complexity varies]

Make one structural change at a time:

1. Extract class/function
2. Run tests
3. Commit
4. Move to new location
5. Run tests
6. Commit
7. Update imports
8. Run tests
9. Commit

**Verification**: Tests pass after each change
```

### Step 4: Cleanup

```markdown
### Step 4: Remove Deprecated Code [simple]
**Files**: Various

Remove old code, update exports, clean imports.

**Verification**: No dead code, all tests pass
```

---

## Bug Fixes

**Sequence**: Root Cause Analysis → Targeted Fix → Prevention Measures

### Phase 1: Understand the Bug

```markdown
### Step 1: Reproduce the Bug [trivial]
Document exact reproduction steps:

```bash
# Steps to reproduce
1. Start the application
2. Navigate to /dashboard
3. Click "Export" button
4. Expected: Download starts
5. Actual: Error shown
```

### Step 2: Find Root Cause [moderate]

**Files**: `src/features/export.ts:45`

Trace through code to find where behavior diverges:

```typescript
// Bug is here: response.body is null when content-type missing
const data = await response.json(); // Throws on null
```

**Verification**: Root cause identified and confirmed

```

### Phase 2: Fix

```markdown
### Step 3: Add Failing Test [simple]
**Files**: `src/features/__tests__/export.test.ts`

Write test that fails with current code:

```typescript
it('handles missing content-type header', () => {
  // This test should fail before fix
});
```

### Step 4: Implement Fix [simple]

**Files**: `src/features/export.ts`

Make minimal change to fix the issue:

```typescript
// Before
const data = await response.json();

// After
const contentType = response.headers.get('content-type');
if (!contentType?.includes('application/json')) {
  throw new ExportError('Invalid response format');
}
const data = await response.json();
```

**Verification**: New test passes, existing tests still pass

```

### Phase 3: Prevent Recurrence

```markdown
### Step 5: Add Defensive Measures [simple]
Consider if this bug class can be prevented:

- Add type validation at boundaries
- Add runtime checks
- Improve error messages
- Add monitoring/logging

**Verification**: Similar bugs would be caught
```

---

## API Design

**Sequence**: Contract First → Handlers → Validation → Documentation

### Step 1: Define Contract

```markdown
### Step 1: Design API Contract [moderate]
**Files**: `src/api/contracts/feature.ts`

Define request/response types:

```typescript
// Request
interface CreateFeatureRequest {
  name: string;
  config: FeatureConfig;
}

// Response
interface CreateFeatureResponse {
  id: string;
  createdAt: string;
}

// Errors
type CreateFeatureError =
  | { code: 'INVALID_NAME'; message: string }
  | { code: 'DUPLICATE'; message: string };
```

**Verification**: Types reviewed and approved

```

### Step 2: Implement Handlers

```markdown
### Step 2: Add Route Handlers [moderate]
**Files**: `src/routes/feature.ts`

Implement endpoints matching contract:

```typescript
app.post('/api/features', async (req, res) => {
  const result = await featureService.create(req.body);
  res.json(result);
});
```

**Verification**: Routes respond correctly

```

### Step 3: Add Validation

```markdown
### Step 3: Add Input Validation [simple]
**Files**: `src/api/validation/feature.ts`

Validate requests match contract:

```typescript
const createFeatureSchema = z.object({
  name: z.string().min(1).max(100),
  config: featureConfigSchema,
});
```

**Verification**: Invalid inputs rejected with clear errors

```

### Step 4: Document

```markdown
### Step 4: Add API Documentation [simple]
**Files**: `docs/api/feature.md` or OpenAPI spec

Document endpoints, parameters, responses, errors.

**Verification**: Documentation matches implementation
```

---

## Extending This Document

To add a new pattern:

1. Identify the common sequence of steps
2. Document each phase with:
   - Step title and complexity
   - Files affected
   - Code examples
   - Verification command
3. Include both happy path and error handling
4. Add to the pattern categories list above
