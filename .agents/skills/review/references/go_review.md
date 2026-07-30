---
title: Go Code Review Best Practices
description: Go code review focusing on error handling, concurrency, idiomatic patterns, and csgda-kit module reuse.
allowed-tools:
  - Bash(go test:*)
  - Bash(go mod:*)
---

# Go Code Review Best Practices

## Reuse csgda-kit Modules (Check First!)

Before implementing new functionality, **always search csgda-kit** for existing modules and patterns:

```
ReferenceSearch(query="HTTP client with retry", repo="csgda-kit")
ReferenceSearch(query="logging middleware", repo="csgda-kit")
ReferenceSearch(query="config loading", repo="csgda-kit")
ReferenceSearch(query="database connection pool", repo="csgda-kit")
```

### Common Modules to Check

| Need              | Search Query                        | Likely Module    |
|-------------------|-------------------------------------|------------------|
| HTTP clients      | `"http client"`                     | `pkg/httpclient` |
| Logging           | `"structured logging" OR "slog"`    | `pkg/log`        |
| Configuration     | `"config" OR "viper"`               | `pkg/config`     |
| Authentication    | `"auth" OR "jwt" OR "oauth"`        | `pkg/auth`       |
| Database          | `"database" OR "postgres" OR "sql"` | `pkg/db`         |
| Metrics           | `"metrics" OR "prometheus"`         | `pkg/metrics`    |
| Tracing           | `"tracing" OR "opentelemetry"`      | `pkg/tracing`    |
| Validation        | `"validation" OR "validate"`        | `pkg/validation` |
| Error handling    | `"errors" OR "error types"`         | `pkg/errors`     |
| Testing utilities | `"test helper" OR "mock"`           | `pkg/testutil`   |

### Review Questions for Reuse

- ❓ Is there an existing module in csgda-kit for this functionality?
- ❓ Does this code duplicate patterns already established in csgda-kit?
- ❓ Should this new utility be contributed back to csgda-kit?
- ❓ Are we following the same patterns used in csgda-kit (error handling, logging, config)?

### Style

When reviewing, ensure code meets organizational style guidance:

```
# Follow csgda-platform-docs go style guide including upstream references
ReferenceSearch(query="go style guide", repo="csgda-platform-docs", scope="org")
```

### Pattern Alignment

When reviewing, ensure code follows csgda-kit patterns:

```
# Check how csgda-kit handles similar concerns
ReferenceSearch(query="error wrapping pattern", repo="csgda-kit")
ReferenceSearch(query="context propagation", repo="csgda-kit")
ReferenceSearch(query="graceful shutdown", repo="csgda-kit")
```

---

## Critical Issues (Must Fix)

### Error Handling

```go
// ❌ BAD: Ignoring errors
result, _ := doSomething()

// ✅ GOOD: Handle all errors
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doSomething failed: %w", err)
}
```

### Goroutine Leaks

```go
// ❌ BAD: Goroutine with no exit condition
go func() {
    for {
        processData()
    }
}()

// ✅ GOOD: Context-aware goroutine
go func(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            processData()
        }
    }
}(ctx)
```

### Defer in Loops

```go
// ❌ BAD: Defer accumulates until function returns
for _, file := range files {
    f, _ := os.Open(file)
    defer f.Close() // Resources held until loop completes!
}

// ✅ GOOD: Use a closure or explicit close
for _, file := range files {
    func() {
        f, _ := os.Open(file)
        defer f.Close()
        // process file
    }()
}
```

### Nil Pointer Dereference

```go
// ❌ BAD: No nil check
func process(user *User) {
    fmt.Println(user.Name) // Panic if nil!
}

// ✅ GOOD: Check for nil
func process(user *User) error {
    if user == nil {
        return errors.New("user cannot be nil")
    }
    fmt.Println(user.Name)
    return nil
}
```

## High Priority Issues

### Context Propagation

```go
// ❌ BAD: Creating new background context
func handler(ctx context.Context) {
    go doWork(context.Background()) // Loses cancellation!
}

// ✅ GOOD: Propagate parent context
func handler(ctx context.Context) {
    go doWork(ctx)
}
```

### Race Conditions

```go
// ❌ BAD: Shared state without synchronization
var counter int
go func() { counter++ }()
go func() { counter++ }()

// ✅ GOOD: Use sync primitives or channels
var (
    counter int
    mu      sync.Mutex
)
go func() { mu.Lock(); counter++; mu.Unlock() }()
```

### Slice Gotchas

```go
// ❌ BAD: Appending to shared backing array
original := []int{1, 2, 3, 4, 5}
slice1 := original[:3]
slice2 := append(slice1, 6) // Modifies original[3]!

// ✅ GOOD: Use full slice expression or copy
slice1 := original[:3:3] // Capacity limited
// or
slice1 := make([]int, 3)
copy(slice1, original[:3])
```

## Code Quality

### Interface Design

```go
// ❌ BAD: Large interface (hard to mock, implement)
type UserService interface {
    Create(user User) error
    Update(user User) error
    Delete(id string) error
    Get(id string) (User, error)
    List() ([]User, error)
    Search(query string) ([]User, error)
    // ... 20 more methods
}

// ✅ GOOD: Small, focused interfaces
type UserReader interface {
    Get(id string) (User, error)
}

type UserWriter interface {
    Create(user User) error
    Update(user User) error
}
```

### Error Wrapping

```go
// ❌ BAD: Loses context
if err != nil {
    return err
}

// ❌ BAD: Hides original error
if err != nil {
    return errors.New("operation failed")
}

// ✅ GOOD: Wrap with context
if err != nil {
    return fmt.Errorf("failed to fetch user %s: %w", userID, err)
}
```

### Struct Initialization

```go
// ❌ BAD: Positional initialization (fragile)
user := User{"John", 30, "john@example.com"}

// ✅ GOOD: Named fields
user := User{
    Name:  "John",
    Age:   30,
    Email: "john@example.com",
}
```

## Performance

### String Building

```go
// ❌ BAD: String concatenation in loop
var result string
for _, s := range items {
    result += s // O(n²) allocations!
}

// ✅ GOOD: Use strings.Builder
var builder strings.Builder
for _, s := range items {
    builder.WriteString(s)
}
result := builder.String()
```

### Preallocate Slices

```go
// ❌ BAD: Growing slice repeatedly
var results []Item
for _, id := range ids {
    results = append(results, fetchItem(id))
}

// ✅ GOOD: Preallocate when size known
results := make([]Item, 0, len(ids))
for _, id := range ids {
    results = append(results, fetchItem(id))
}
```

### Map Preallocation

```go
// ❌ BAD: Map grows dynamically
m := make(map[string]int)

// ✅ GOOD: Preallocate when size known
m := make(map[string]int, expectedSize)
```

## Testing

### Table-Driven Tests

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 1, 2, 3},
        {"negative", -1, -2, -3},
        {"zero", 0, 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d",
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

### Test Helpers

```go
// Use t.Helper() for test utilities
func assertEqual(t *testing.T, got, want interface{}) {
    t.Helper() // Marks this as helper - errors report caller's line
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}
```

## Static Analysis

Run `/lint` to handle all linting (tool setup, `go vet`, `golangci-lint` with CSGDAA org linters). It auto-installs missing tools.

For review-specific checks beyond linting (mod health, tests, TODOs), use the review script:

```bash
# Full Go review (mod checks + tests — run after /lint)
bash {baseDir}/scripts/analyze-go.sh

# With auto-fix
bash {baseDir}/scripts/analyze-go.sh --fix

# Skip tests
SKIP_TESTS=true bash {baseDir}/scripts/analyze-go.sh
```

### CSGDAA Org Linters

Five custom Go analyzers at `github.com/bayer-int/csgda-go-linters`, run via golangci-lint:

- **ctxname** — `context.Context` variables must be named `ctx`
- **emptyinitialism** — Use `var x T` not `x := T{}`; use `new(T)` not `&T{}`
- **nohiddenexports** — Exported functions must not accept/return unexported types
- **optionnames** — `Set` for single-value options, `With` for additive/variadic

For the full style guide and examples, see the `/lint` skill reference.

## Go 1.21+ Features

### `slices` and `maps` packages

```go
// Use standard library generics
import "slices"

slices.Sort(items)
slices.Contains(items, target)
idx := slices.Index(items, target)
```

### Structured Logging (slog)

```go
import "log/slog"

slog.Info("user created",
    "user_id", user.ID,
    "email", user.Email)
```

## Security Checklist

- [ ] SQL queries use parameterized statements
- [ ] User input is validated and sanitized
- [ ] Secrets not hardcoded (use environment variables)
- [ ] TLS/HTTPS used for external connections
- [ ] Rate limiting on public endpoints
- [ ] No sensitive data in logs
- [ ] Proper authentication/authorization checks
