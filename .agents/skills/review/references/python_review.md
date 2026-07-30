---
title: Python Code Review Best Practices
description: Python review covering type hints, resource management, exception handling, and common pitfalls.
allowed-tools:
  # Python linting and type checking
  - Bash(ruff:*)
  - Bash(mypy:*)
  - Bash(pyright:*)
  - Bash(pytest:*)
  - Bash(bandit:*)
  - Bash(radon:*)
  - Bash(pip-audit:*)
  - Bash(safety:*)
  # Python tool installation
  - Bash(pip install:*)
  - Bash(pip:*)
  - Bash(uv:*)
---

# Python Code Review Best Practices

## Critical Issues (Must Fix)

### Bare Except Clauses

```python
# ❌ BAD: Catches everything including KeyboardInterrupt
try:
    do_something()
except:
    pass

# ❌ BAD: Still too broad
try:
    do_something()
except Exception:
    pass

# ✅ GOOD: Catch specific exceptions
try:
    do_something()
except ValueError as e:
    logger.error(f"Invalid value: {e}")
except IOError as e:
    logger.error(f"IO error: {e}")
    raise
```

### Mutable Default Arguments

```python
# ❌ BAD: Mutable default (shared across calls!)
def add_item(item, items=[]):
    items.append(item)
    return items

add_item(1)  # [1]
add_item(2)  # [1, 2] - NOT [2]!

# ✅ GOOD: Use None as default
def add_item(item, items=None):
    if items is None:
        items = []
    items.append(item)
    return items
```

### Security: Pickle/Eval

```python
# ❌ BAD: Arbitrary code execution
import pickle
data = pickle.loads(user_input)  # Remote code execution!

# ❌ BAD: Never eval user input
result = eval(user_expression)

# ✅ GOOD: Use safe alternatives
import json
data = json.loads(user_input)

# ✅ GOOD: Use ast.literal_eval for simple literals
import ast
result = ast.literal_eval(user_expression)
```

### SQL Injection

```python
# ❌ BAD: String formatting in SQL
cursor.execute(f"SELECT * FROM users WHERE id = {user_id}")

# ✅ GOOD: Parameterized queries
cursor.execute("SELECT * FROM users WHERE id = ?", (user_id,))

# ✅ GOOD: With SQLAlchemy
session.query(User).filter(User.id == user_id).first()
```

## High Priority Issues

### Missing Type Hints

```python
# ❌ BAD: No type information
def process(data):
    return data.get("value")

# ✅ GOOD: Full type hints
from typing import Any

def process(data: dict[str, Any]) -> Any | None:
    return data.get("value")

# ✅ BETTER: Specific types
def process(data: dict[str, str]) -> str | None:
    return data.get("value")
```

### Resource Management

```python
# ❌ BAD: Resource leak
file = open("data.txt")
content = file.read()
# file never closed!

# ✅ GOOD: Context manager
with open("data.txt") as file:
    content = file.read()

# ✅ GOOD: For async
async with aiofiles.open("data.txt") as file:
    content = await file.read()
```

### Boolean Comparisons

```python
# ❌ BAD: Explicit comparison to True/False
if is_valid == True:
    do_something()

if items != None:
    process(items)

# ✅ GOOD: Pythonic comparisons
if is_valid:
    do_something()

if items is not None:
    process(items)

# ✅ GOOD: Check for empty collections
if items:  # True if non-empty
    process(items)
```

### String Formatting

```python
# ❌ BAD: Old-style formatting
message = "Hello %s, you have %d messages" % (name, count)

# ❌ BAD: .format() is verbose
message = "Hello {}, you have {} messages".format(name, count)

# ✅ GOOD: f-strings (Python 3.6+)
message = f"Hello {name}, you have {count} messages"

# ✅ GOOD: For logging (lazy evaluation)
logger.debug("Processing %s items", len(items))
```

## Code Quality

### Style

When reviewing, ensure code meets organizational style guidance:

```
# Follow csgda-platform-docs python style guide including upstream references
ReferenceSearch(query="python style guide", repo="csgda-platform-docs", scope="org")
```

### List Comprehensions

```python
# ❌ BAD: Loop with append
result = []
for item in items:
    if item.active:
        result.append(item.value)

# ✅ GOOD: List comprehension
result = [item.value for item in items if item.active]

# ✅ GOOD: Generator for large datasets
result = (item.value for item in items if item.active)
```

### Dictionary Operations

```python
# ❌ BAD: Checking then accessing
if key in dictionary:
    value = dictionary[key]
else:
    value = default

# ✅ GOOD: Use .get()
value = dictionary.get(key, default)

# ✅ GOOD: setdefault for mutable defaults
cache.setdefault(key, []).append(item)

# ✅ GOOD: defaultdict
from collections import defaultdict
cache = defaultdict(list)
cache[key].append(item)
```

### Unpacking

```python
# ❌ BAD: Index access
first = items[0]
rest = items[1:]

# ✅ GOOD: Unpacking
first, *rest = items

# ✅ GOOD: Swap without temp
a, b = b, a

# ✅ GOOD: Enumerate for index
for i, item in enumerate(items):
    print(f"{i}: {item}")
```

### Class Design

```python
# ❌ BAD: No __slots__ for simple classes (memory overhead)
class Point:
    def __init__(self, x, y):
        self.x = x
        self.y = y

# ✅ GOOD: Use dataclasses
from dataclasses import dataclass

@dataclass
class Point:
    x: float
    y: float

# ✅ GOOD: For immutability
@dataclass(frozen=True)
class Point:
    x: float
    y: float
```

### Path Handling

```python
# ❌ BAD: String concatenation for paths
path = folder + "/" + filename

# ✅ GOOD: Use pathlib
from pathlib import Path
path = Path(folder) / filename
```

## Performance

### Avoid Global Lookups in Loops

```python
# ❌ BAD: Global lookup each iteration
for item in items:
    result.append(len(item))  # len() looked up each time

# ✅ GOOD: Local reference
local_len = len
for item in items:
    result.append(local_len(item))

# ✅ BETTER: List comprehension
result = [len(item) for item in items]
```

### Generator vs List

```python
# ❌ BAD: Materializes entire list
total = sum([x * 2 for x in range(1000000)])

# ✅ GOOD: Generator expression (lazy)
total = sum(x * 2 for x in range(1000000))
```

### String Joining

```python
# ❌ BAD: String concatenation in loop
result = ""
for s in strings:
    result += s  # O(n²)!

# ✅ GOOD: join()
result = "".join(strings)
```

### Dictionary Key Operations

```python
# ❌ BAD: List to check membership
items_to_check = [1, 2, 3, 4, 5]
if x in items_to_check:  # O(n)

# ✅ GOOD: Set for membership testing
items_to_check = {1, 2, 3, 4, 5}
if x in items_to_check:  # O(1)
```

## Testing

### Pytest Best Practices

```python
import pytest

class TestUserService:
    @pytest.fixture
    def user_service(self, db_session):
        return UserService(db_session)

    def test_create_user_success(self, user_service):
        # Arrange
        user_data = {"name": "Test", "email": "test@example.com"}

        # Act
        result = user_service.create_user(user_data)

        # Assert
        assert result.id is not None
        assert result.name == "Test"
        assert result.email == "test@example.com"

    def test_create_user_invalid_email(self, user_service):
        user_data = {"name": "Test", "email": "invalid"}

        with pytest.raises(ValueError, match="Invalid email"):
            user_service.create_user(user_data)

    @pytest.mark.parametrize("email,valid", [
        ("test@example.com", True),
        ("invalid", False),
        ("", False),
        ("test@test", False),
    ])
    def test_email_validation(self, email, valid):
        assert validate_email(email) == valid
```

### Mocking

```python
from unittest.mock import Mock, patch, AsyncMock

def test_external_api_call():
    with patch("module.requests.get") as mock_get:
        mock_get.return_value.json.return_value = {"data": "test"}

        result = fetch_data()

        assert result == {"data": "test"}
        mock_get.assert_called_once_with("https://api.example.com")

# For async
async def test_async_call():
    mock_client = AsyncMock()
    mock_client.fetch.return_value = {"data": "test"}

    result = await process_data(mock_client)

    assert result == {"data": "test"}
```

## Static Analysis Tools

| Tool | Purpose | Command |
|------|---------|---------|
| `ruff` | Fast linter + formatter | `ruff check . && ruff format --check .` |
| `mypy` | Type checking | `mypy .` |
| `pyright` | Type checking (faster) | `pyright` |
| `bandit` | Security linting | `bandit -r .` |
| `radon` | Complexity metrics | `radon cc -s .` |
| `vulture` | Dead code detection | `vulture .` |
| `pip-audit` | Dependency vulnerabilities | `pip-audit` |
| `safety` | Dependency vulnerabilities | `safety check` |
| `black` | Code formatting | `black --check .` |
| `isort` | Import sorting | `isort --check .` |

## Python Version Features

### Python 3.10+

```python
# Structural pattern matching
match command:
    case ["quit"]:
        return
    case ["load", filename]:
        load_file(filename)
    case _:
        print("Unknown command")

# Union type syntax
def process(value: int | str) -> str:
    return str(value)
```

### Python 3.11+

```python
# Exception groups
try:
    async with asyncio.TaskGroup() as tg:
        tg.create_task(task1())
        tg.create_task(task2())
except* ValueError as eg:
    for exc in eg.exceptions:
        handle_value_error(exc)
```

### Python 3.12+

```python
# Type parameter syntax
def first[T](items: list[T]) -> T:
    return items[0]

class Stack[T]:
    def push(self, item: T) -> None: ...
    def pop(self) -> T: ...
```

## Security Checklist

- [ ] No hardcoded secrets (use environment variables)
- [ ] Input validation on all external data
- [ ] Parameterized SQL queries
- [ ] No `eval()`, `exec()`, or `pickle.loads()` on user input
- [ ] HTTPS for all external requests
- [ ] Dependencies audited for vulnerabilities
- [ ] No sensitive data in logs
- [ ] Proper authentication/authorization
- [ ] Rate limiting on APIs
- [ ] CSRF protection on forms
