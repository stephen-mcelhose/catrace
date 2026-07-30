---
title: SQL Code Review Best Practices
description: SQL review for BigQuery/Spanner performance, security, query optimization, and best practices.
---

# SQL Code Review Best Practices

> **Primary Platforms:** This guide focuses on **BigQuery** and **Cloud Spanner**, our primary SQL platforms.

## BigQuery Best Practices

### Cost Optimization

```sql
-- ❌ BAD: SELECT * scans all columns (expensive!)
SELECT * FROM `project.dataset.large_table`;

-- ✅ GOOD: Select only needed columns
SELECT user_id, event_type, created_at
FROM `project.dataset.large_table`;

-- ❌ BAD: No partition filter (full table scan)
SELECT * FROM `project.dataset.events`
WHERE event_type = 'login';

-- ✅ GOOD: Filter on partition column
SELECT * FROM `project.dataset.events`
WHERE DATE(event_timestamp) = '2024-01-15'
  AND event_type = 'login';
```

### Partitioning and Clustering

```sql
-- ✅ GOOD: Create partitioned and clustered table
CREATE TABLE `project.dataset.events`
(
    event_id STRING NOT NULL,
    user_id STRING,
    event_type STRING,
    event_data JSON,
    event_timestamp TIMESTAMP NOT NULL
)
PARTITION BY DATE(event_timestamp)
CLUSTER BY user_id, event_type
OPTIONS(
    partition_expiration_days = 365,
    require_partition_filter = true
);
```

### Use QUALIFY for Window Functions

```sql
-- ❌ BAD: Subquery for row_number filter
SELECT * FROM (
    SELECT *, ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY created_at DESC) as rn
    FROM `project.dataset.orders`
)
WHERE rn = 1;

-- ✅ GOOD: QUALIFY clause (BigQuery)
SELECT *
FROM `project.dataset.orders`
QUALIFY ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY created_at DESC) = 1;
```

### Avoid Repeated Subqueries

```sql
-- ❌ BAD: Repeated expensive subquery
SELECT
    a.metric_1 / (SELECT SUM(value) FROM `project.dataset.totals`),
    a.metric_2 / (SELECT SUM(value) FROM `project.dataset.totals`)
FROM `project.dataset.metrics` a;

-- ✅ GOOD: Use CTE or CROSS JOIN
WITH totals AS (
    SELECT SUM(value) as total FROM `project.dataset.totals`
)
SELECT
    a.metric_1 / t.total,
    a.metric_2 / t.total
FROM `project.dataset.metrics` a
CROSS JOIN totals t;
```

### BigQuery-Specific Functions

```sql
-- ✅ GOOD: Use BigQuery array functions
SELECT
    user_id,
    ARRAY_AGG(STRUCT(product_id, quantity)) as order_items
FROM `project.dataset.order_items`
GROUP BY user_id;

-- ✅ GOOD: UNNEST for array expansion
SELECT user_id, item.product_id
FROM `project.dataset.orders`,
UNNEST(items) as item;

-- ✅ GOOD: Use SAFE functions to avoid errors
SELECT SAFE_DIVIDE(revenue, users) as arpu
FROM `project.dataset.metrics`;

-- ✅ GOOD: Parse JSON safely
SELECT
    JSON_VALUE(event_data, '$.user_id') as user_id,
    JSON_VALUE(event_data, '$.action') as action
FROM `project.dataset.events`;
```

### Dry Run Before Execution

```sql
-- Always check query cost before running large queries
-- Use BigQuery console's "Dry Run" or:
-- bq query --dry_run --use_legacy_sql=false "SELECT ..."
```

---

## Cloud Spanner Best Practices

### Primary Key Design

```sql
-- ❌ BAD: Sequential primary key (hot spots!)
CREATE TABLE Orders (
    OrderId INT64 NOT NULL,
    UserId STRING(36),
    CreatedAt TIMESTAMP
) PRIMARY KEY (OrderId);

-- ✅ GOOD: Use UUID or hash-prefixed keys
CREATE TABLE Orders (
    OrderId STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
    UserId STRING(36),
    CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (OrderId);

-- ✅ GOOD: Composite key to distribute load
CREATE TABLE UserOrders (
    UserId STRING(36) NOT NULL,
    OrderId STRING(36) NOT NULL,
    CreatedAt TIMESTAMP
) PRIMARY KEY (UserId, OrderId);
```

### Interleaved Tables

```sql
-- ✅ GOOD: Interleave child tables for locality
CREATE TABLE Users (
    UserId STRING(36) NOT NULL,
    Email STRING(255),
    Name STRING(100)
) PRIMARY KEY (UserId);

CREATE TABLE UserOrders (
    UserId STRING(36) NOT NULL,
    OrderId STRING(36) NOT NULL,
    Total NUMERIC,
    CreatedAt TIMESTAMP
) PRIMARY KEY (UserId, OrderId),
  INTERLEAVE IN PARENT Users ON DELETE CASCADE;

-- Reads of user + orders are fast (same split)
```

### Use Query Parameters

```go
// ❌ BAD: String interpolation (Spanner doesn't allow anyway)
// query := fmt.Sprintf("SELECT * FROM Users WHERE UserId = '%s'", userId)

// ✅ GOOD: Parameters
stmt := spanner.Statement{
    SQL: "SELECT * FROM Users WHERE UserId = @userId",
    Params: map[string]interface{}{
        "userId": userId,
    },
}
```

### Avoid Hot Spots

```sql
-- ❌ BAD: Monotonically increasing key causes hot spots
-- All inserts go to the same split!
CREATE TABLE Events (
    EventId INT64 NOT NULL,  -- Auto-increment pattern
    ...
) PRIMARY KEY (EventId);

-- ✅ GOOD: UUID or reverse timestamp
CREATE TABLE Events (
    EventId STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
    -- or ShardId + Timestamp pattern
    ShardId INT64 NOT NULL,
    EventTimestamp TIMESTAMP NOT NULL,
    ...
) PRIMARY KEY (EventId);
-- or PRIMARY KEY (ShardId, EventTimestamp DESC)
```

### Commit Timestamps

```sql
-- ✅ GOOD: Use commit timestamps for ordering
CREATE TABLE AuditLog (
    LogId STRING(36) NOT NULL,
    Action STRING(50),
    CommitTime TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (LogId);

-- Insert with commit timestamp
INSERT INTO AuditLog (LogId, Action, CommitTime)
VALUES ('uuid-here', 'user_created', PENDING_COMMIT_TIMESTAMP());
```

### Read-Write vs Read-Only Transactions

```go
// ✅ GOOD: Use read-only for queries (better performance)
iter := client.Single().Query(ctx, stmt)  // Strong read
// or
iter := client.Single().WithTimestampBound(
    spanner.ExactStaleness(15*time.Second),
).Query(ctx, stmt)  // Stale read (even faster)

// ✅ GOOD: Use read-write only when needed
_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
    // mutations here
    return nil
})
```

---

## Common SQL Issues (Both Platforms)

### Avoid SELECT *

```sql
-- ❌ BAD: SELECT * (unknown columns, extra data)
SELECT * FROM users;

-- ✅ GOOD: Explicit column list
SELECT id, name, email, created_at
FROM users;
```

### NULL Handling

```sql
-- ❌ BAD: Comparing with NULL
SELECT * FROM users WHERE phone = NULL;
-- Returns no rows! NULL = NULL is NULL, not TRUE

-- ✅ GOOD: Use IS NULL
SELECT * FROM users WHERE phone IS NULL;

-- ✅ GOOD: Handle NULLs with COALESCE/IFNULL
SELECT
    name,
    COALESCE(nickname, name) as display_name,
    IFNULL(phone, 'No phone') as contact
FROM users;
```

### Use CTEs for Readability

```sql
-- ✅ GOOD: CTEs make complex queries readable
WITH active_users AS (
    SELECT user_id, email
    FROM users
    WHERE status = 'active'
),
recent_orders AS (
    SELECT user_id, COUNT(*) as order_count
    FROM orders
    WHERE created_at > TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 30 DAY)
    GROUP BY user_id
)
SELECT
    u.email,
    COALESCE(o.order_count, 0) as orders_last_30_days
FROM active_users u
LEFT JOIN recent_orders o ON u.user_id = o.user_id;
```

### Batch Operations

```sql
-- ❌ BAD: Individual inserts
INSERT INTO logs (message) VALUES ('msg1');
INSERT INTO logs (message) VALUES ('msg2');
INSERT INTO logs (message) VALUES ('msg3');

-- ✅ GOOD: Batch insert
INSERT INTO logs (message) VALUES
    ('msg1'),
    ('msg2'),
    ('msg3');
```

## Static Analysis Tools

| Tool | Purpose | Command |
|------|---------|---------|
| `bq query --dry_run` | Estimate BigQuery cost | `bq query --dry_run --use_legacy_sql=false "SQL"` |
| `sqlfluff` | SQL linting | `sqlfluff lint query.sql --dialect bigquery` |
| `sqlfmt` | SQL formatting | `sqlfmt query.sql` |

## Security Checklist

- [ ] All queries use parameterized statements
- [ ] No user input in SQL strings
- [ ] IAM permissions follow least privilege
- [ ] Audit logging enabled
- [ ] Data Access Logs configured (BigQuery)
- [ ] Fine-grained access control where needed (Spanner)
- [ ] Column-level security for sensitive data
- [ ] VPC Service Controls configured
- [ ] CMEK encryption if required
