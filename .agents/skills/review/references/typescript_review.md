---
title: TypeScript Code Review Best Practices
description: Critical issues, type safety, async/await patterns, and common anti-patterns for TypeScript/JavaScript code review.
allowed-tools:
  # TypeScript/JavaScript package managers
  - Bash(bun run:*)
  - Bash(bun test:*)
  - Bash(bun tsc:*)
  - Bash(bunx:*)
  - Bash(npm run:*)
  - Bash(npm test:*)
  - Bash(npm audit:*)
  - Bash(npx:*)
  - Bash(pnpm run:*)
  - Bash(pnpm test:*)
  - Bash(pnpm audit:*)
  - Bash(pnpm exec:*)
  - Bash(yarn run:*)
  - Bash(yarn test:*)
  - Bash(yarn audit:*)
  - Bash(yarn exec:*)
  # TypeScript linting/type checking
  - Bash(eslint:*)
  - Bash(tsc:*)
  - Bash(prettier:*)
  - Bash(madge:*)
  - Bash(depcheck:*)
  # Tool installation
  - Bash(npm install -g:*)
---

# TypeScript Code Review Best Practices

## Critical Issues (Must Fix)

### Avoid `any` Type

```typescript
// ❌ BAD: any defeats type safety
function process(data: any) {
    return data.foo.bar; // No type checking!
}

// ✅ GOOD: Use proper types or unknown
function process(data: unknown) {
    if (isValidData(data)) {
        return data.foo.bar;
    }
    throw new Error('Invalid data');
}

// ✅ GOOD: Use generics
function process<T extends { foo: { bar: string } }>(data: T) {
    return data.foo.bar;
}
```

### Null/Undefined Checks

```typescript
// ❌ BAD: Assuming value exists
function getName(user: User | null) {
    return user.name; // TypeError if null!
}

// ✅ GOOD: Optional chaining and nullish coalescing
function getName(user: User | null) {
    return user?.name ?? 'Anonymous';
}

// ✅ GOOD: Type guard
function getName(user: User | null) {
    if (!user) {
        throw new Error('User required');
    }
    return user.name; // TypeScript knows user is not null
}
```

### Async/Await Error Handling

```typescript
// ❌ BAD: Unhandled promise rejection
async function fetchData() {
    const response = await fetch(url); // No error handling!
    return response.json();
}

// ✅ GOOD: Proper error handling
async function fetchData() {
    try {
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}`);
        }
        return await response.json();
    } catch (error) {
        console.error('Fetch failed:', error);
        throw error;
    }
}
```

### Promise Handling

```typescript
// ❌ BAD: Floating promise (no await)
function handleClick() {
    saveData(); // Promise not awaited!
}

// ❌ BAD: Promise in non-async function
function handleClick() {
    return saveData(); // Returns Promise but caller may ignore
}

// ✅ GOOD: Await in async function
async function handleClick() {
    await saveData();
}

// ✅ GOOD: Explicit void for fire-and-forget
function handleClick() {
    void saveData().catch(console.error);
}
```

## High Priority Issues

### Type Assertions

```typescript
// ❌ BAD: Unsafe type assertion
const user = response as User; // No runtime check!

// ❌ WORSE: Double assertion
const user = response as unknown as User;

// ✅ GOOD: Type guard with runtime check
function isUser(value: unknown): value is User {
    return (
        typeof value === 'object' &&
        value !== null &&
        'id' in value &&
        'name' in value
    );
}

if (isUser(response)) {
    console.log(response.name);
}

// ✅ GOOD: Use zod or similar for validation
const UserSchema = z.object({
    id: z.string(),
    name: z.string(),
});
const user = UserSchema.parse(response);
```

### Non-null Assertion

```typescript
// ❌ BAD: Non-null assertion without guarantee
function process(item: Item | undefined) {
    return item!.name; // Crashes if undefined!
}

// ✅ GOOD: Proper null check
function process(item: Item | undefined) {
    if (!item) {
        throw new Error('Item required');
    }
    return item.name;
}
```

### Mutable Default Parameters

```typescript
// ❌ BAD: Mutable default (shared reference)
function addItem(item: string, list: string[] = []) {
    list.push(item);
    return list;
}
// Multiple calls share same array!

// ✅ GOOD: Create new instance
function addItem(item: string, list?: string[]) {
    const items = list ?? [];
    return [...items, item];
}
```

## Code Quality

### Prefer `const` over `let`

```typescript
// ❌ BAD: Using let when value doesn't change
let config = { timeout: 5000 };

// ✅ GOOD: Use const
const config = { timeout: 5000 };

// ✅ GOOD: Use const even for reassignment (use array methods)
// ❌ BAD
let result = [];
for (const item of items) {
    result.push(transform(item));
}

// ✅ GOOD
const result = items.map(transform);
```

### Object Shorthand

```typescript
// ❌ BAD: Redundant property names
const user = { name: name, age: age };

// ✅ GOOD: Object shorthand
const user = { name, age };
```

### Template Literals

```typescript
// ❌ BAD: String concatenation
const message = 'Hello, ' + name + '! You have ' + count + ' messages.';

// ✅ GOOD: Template literals
const message = `Hello, ${name}! You have ${count} messages.`;
```

### Destructuring

```typescript
// ❌ BAD: Repetitive property access
function process(user: User) {
    console.log(user.name);
    console.log(user.email);
    console.log(user.age);
}

// ✅ GOOD: Destructuring
function process({ name, email, age }: User) {
    console.log(name);
    console.log(email);
    console.log(age);
}
```

### Enum Alternatives

```typescript
// ⚠️ AVOID: Enums have runtime overhead
enum Status {
    Active = 'ACTIVE',
    Inactive = 'INACTIVE',
}

// ✅ PREFER: const objects or union types
const Status = {
    Active: 'ACTIVE',
    Inactive: 'INACTIVE',
} as const;
type Status = typeof Status[keyof typeof Status];

// ✅ SIMPLE: Union type
type Status = 'ACTIVE' | 'INACTIVE';
```

## Performance

### Array Method Chaining

```typescript
// ❌ BAD: Multiple iterations
const result = items
    .filter(item => item.active)
    .map(item => item.value)
    .filter(value => value > 0);

// ✅ GOOD: Reduce iterations when performance matters
const result = items.reduce<number[]>((acc, item) => {
    if (item.active && item.value > 0) {
        acc.push(item.value);
    }
    return acc;
}, []);
```

### Avoid Unnecessary Spread

```typescript
// ❌ BAD: Spreading for no reason
const newItems = [...items].sort();

// ✅ GOOD: Use toSorted() (ES2023+)
const newItems = items.toSorted();

// ❌ BAD: Spreading just to add properties
const updated = { ...user, ...{ name: 'new' } };

// ✅ GOOD: Direct object spread
const updated = { ...user, name: 'new' };
```

### Lazy Evaluation

```typescript
// ❌ BAD: Compute even if not needed
function getValue(condition: boolean) {
    const expensiveValue = computeExpensiveValue();
    return condition ? expensiveValue : defaultValue;
}

// ✅ GOOD: Lazy evaluation
function getValue(condition: boolean) {
    return condition ? computeExpensiveValue() : defaultValue;
}
```

## Testing

### Proper Mocking

```typescript
// ❌ BAD: Mocking implementation details
jest.mock('./utils', () => ({
    formatDate: jest.fn().mockReturnValue('2024-01-01'),
}));

// ✅ GOOD: Mock at boundary (API, DB)
const mockFetch = jest.fn().mockResolvedValue({
    ok: true,
    json: () => Promise.resolve({ data: 'test' }),
});
global.fetch = mockFetch;
```

### Test Structure

```typescript
describe('UserService', () => {
    describe('createUser', () => {
        it('should create user with valid data', async () => {
            // Arrange
            const userData = { name: 'Test', email: 'test@example.com' };

            // Act
            const result = await userService.createUser(userData);

            // Assert
            expect(result).toMatchObject({
                id: expect.any(String),
                ...userData,
            });
        });

        it('should throw on invalid email', async () => {
            // Arrange
            const userData = { name: 'Test', email: 'invalid' };

            // Act & Assert
            await expect(userService.createUser(userData))
                .rejects.toThrow('Invalid email');
        });
    });
});
```

## Static Analysis Tools

| Tool | Purpose | Command |
|------|---------|---------|
| `tsc` | Type checking | `tsc --noEmit` |
| `eslint` | Linting | `eslint . --ext .ts,.tsx` |
| `prettier` | Formatting | `prettier --check .` |
| `typescript-eslint` | TS-specific rules | (configure in eslint) |
| `knip` | Dead code detection | `knip` |
| `depcheck` | Unused dependencies | `depcheck` |
| `madge` | Circular dependencies | `madge --circular src/` |

## Strict Mode Checklist

Ensure `tsconfig.json` has:

```json
{
    "compilerOptions": {
        "strict": true,
        "noUncheckedIndexedAccess": true,
        "noImplicitReturns": true,
        "noFallthroughCasesInSwitch": true,
        "exactOptionalPropertyTypes": true
    }
}
```

## Security Checklist

- [ ] No `eval()` or `Function()` constructors
- [ ] No `dangerouslySetInnerHTML` without sanitization
- [ ] User input validated (prefer zod/yup)
- [ ] No secrets in client-side code
- [ ] CORS properly configured
- [ ] Content Security Policy headers set
- [ ] Dependencies audited (`npm audit`)
- [ ] No prototype pollution vulnerabilities
