---
title: Pattern Finding
description: How to find similar implementations to model after
---

# Pattern Finding

When you need to find existing patterns and implementations to model new code after.

## Core Principle

**Find concrete examples, not abstract recommendations.**

Search for comparable features, usage examples, and established patterns. Extract reusable code structures with snippets and `file:line` references.

## Pattern Categories

### API Patterns

- Route structures and naming
- Middleware usage
- Error handling approaches
- Authentication/authorization
- Validation patterns
- Pagination implementations

### Data Patterns

- Database query patterns
- Caching strategies
- Data transformation
- Migration patterns

### Component Patterns

- File organization
- State management
- Event handling
- Lifecycle methods

### Testing Patterns

- Unit test structure
- Integration test setup
- Mocking strategies
- Assertion patterns

## Search Strategy

### Step 1: Identify What You're Looking For

Before searching, clarify:

- What type of pattern? (API, data, component, test)
- What specific behavior? (CRUD, auth, validation)
- What's the closest existing feature?

### Step 2: Search for Similar Implementations

**Find files with similar exports:**

```
Grep
  pattern: "export.*Service"
  output_mode: "files_with_matches"
```

**Find specific class patterns in a directory:**

```
Grep
  pattern: "class.*Controller"
  path: "src/"
```

**Find test files by naming convention:**

```
Glob
  pattern: "**/*.test.ts"
```

**Use LSP to understand structure without reading full files:**

```
Lsp
  operation: "documentSymbol"
  filePath: "src/services/user.service.ts"
```

### Step 3: Extract and Document

Read the found files and extract:

- The pattern structure
- Key code snippets
- Configuration requirements
- Related utilities

## Output Format

```markdown
## Patterns Found for [Feature Type]

### Pattern: Service Layer Structure

**Files**: `src/services/user.service.ts`, `src/services/order.service.ts`

**Usage Context**: All business logic services follow this pattern.

**Example** (`src/services/user.service.ts:10-45`):
```typescript
export class UserService {
  constructor(private db: Database, private cache: Cache) {}

  async findById(id: string): Promise<User | null> {
    // Check cache first
    const cached = await this.cache.get(`user:${id}`);
    if (cached) return cached;

    // Query database
    const user = await this.db.users.findUnique({ where: { id } });

    // Cache result
    if (user) await this.cache.set(`user:${id}`, user, { ttl: 300 });

    return user;
  }
}
```

**Key Aspects**:

- Constructor injection for dependencies
- Cache-first read pattern
- TTL-based cache expiration

**Frequency**: Found in 8 service files

**Related Utilities**:

- `src/utils/cache.ts` - Cache wrapper
- `src/utils/db.ts` - Database client

---

### Pattern: Error Handling

**Files**: `src/middleware/error.ts`, `src/handlers/*.ts`

**Example** (`src/middleware/error.ts:20-35`):

```typescript
// Error handling pattern used across handlers
```

```

## Guidelines

- **Show concrete code** - Not abstract descriptions
- **Include file:line references** - Be precise
- **Note frequency** - How common is this pattern?
- **List related utilities** - What helpers exist?
- **Don't evaluate** - Document patterns without judging them

## What NOT to Do

- Don't suggest which pattern is "better"
- Don't critique existing patterns
- Don't recommend refactoring
- Don't evaluate efficiency
- Don't comment on whether patterns are optimal

## Remember

You are a pattern documentarian. Catalog existing implementations exactly as they appear so they can be modeled for new code.
