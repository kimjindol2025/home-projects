# Phase 11: Complete SQLite API Reference

**Status**: 🟢 **COMPLETE** - Comprehensive API documentation for all phases

---

## 📚 Table of Contents

1. [Phase 11.1: Native SQLite Interface](#phase-111-native-sqlite-interface)
2. [Phase 11.2: Connection Pool](#phase-112-connection-pool)
3. [Phase 11.3: Query Builder](#phase-113-query-builder)
4. [Phase 11.4: Performance Cache](#phase-114-performance-cache)
5. [Integration Patterns](#integration-patterns)
6. [Performance Tuning](#performance-tuning)
7. [Deployment](#deployment)

---

## Phase 11.1: Native SQLite Interface

### Overview

Low-level bindings to SQLite3 C API. Use when you need maximum control over database operations.

### Classes & Interfaces

```typescript
interface SQLiteConnection {
  handle: number;           // Opaque pointer
  filename: string;         // Database file path
  isOpen: boolean;          // Connection state
  lastError: string | null; // Last error message
}

interface SQLiteStatement {
  handle: number;    // Statement pointer
  sql: string;       // SQL text
  columnCount: number;
  paramCount: number;
}

type SQLiteRow = {
  [columnName: string]: any;
}
```

### API Methods

#### Connection Management

**`SQLiteNative.openDatabase(filename: string): SQLiteConnection`**

Opens a connection to SQLite database.

```typescript
const db = SQLiteNative.openDatabase(':memory:');          // In-memory
const db = SQLiteNative.openDatabase('./data.db');         // File-based
const db = SQLiteNative.openDatabase('./data.db?mode=ro'); // Read-only
```

**`SQLiteNative.closeDatabase(conn: SQLiteConnection): boolean`**

Closes connection and releases resources.

```typescript
const success = SQLiteNative.closeDatabase(db);
if (!success) console.error('Already closed');
```

#### Statement Operations

**`SQLiteNative.prepareStatement(conn: SQLiteConnection, sql: string): SQLiteStatement`**

Compiles SQL statement for execution.

```typescript
const stmt = SQLiteNative.prepareStatement(
  db,
  'SELECT * FROM users WHERE id = ? AND status = ?'
);
console.log(`Parameters: ${stmt.paramCount}`);  // 2
console.log(`Result columns: ${stmt.columnCount}`);
```

**`SQLiteNative.finalizeStatement(stmt: SQLiteStatement): boolean`**

Deallocates statement and releases resources.

```typescript
const success = SQLiteNative.finalizeStatement(stmt);
```

#### Parameter Binding

**`SQLiteNative.bindParameter(stmt: SQLiteStatement, index: number, value: any): boolean`**

Binds value to parameter. Index is 1-based (SQL standard).

```typescript
// String parameter
SQLiteNative.bindParameter(stmt, 1, 'active');

// Integer parameter
SQLiteNative.bindParameter(stmt, 2, 42);

// Null parameter
SQLiteNative.bindParameter(stmt, 3, null);

// Binary/Blob parameter
const buffer = Buffer.from('binary data');
SQLiteNative.bindParameter(stmt, 4, buffer);
```

**Supported Types**:
- `string` → `sqlite3_bind_text()`
- `number` → `sqlite3_bind_int64()` / `sqlite3_bind_double()`
- `null` → `sqlite3_bind_null()`
- `Buffer` / `Uint8Array` → `sqlite3_bind_blob()`

#### Result Fetching

**`SQLiteNative.step(stmt: SQLiteStatement): SQLiteRow | null`**

Fetches next row from result set.

```typescript
let row;
while ((row = SQLiteNative.step(stmt)) !== null) {
  console.log(row);  // { id: 1, name: 'Alice', status: 'active' }
}
```

**`SQLiteNative.reset(stmt: SQLiteStatement): boolean`**

Resets statement for re-execution.

```typescript
SQLiteNative.reset(stmt);  // Now can call step() again from beginning
```

#### Error Handling

**`SQLiteNative.getLastError(conn: SQLiteConnection): string`**

Gets error message from last failed operation.

```typescript
const error = SQLiteNative.getLastError(db);
console.error(`Database error: ${error}`);
```

#### Convenience Methods

**`SQLiteNative.queryAll(conn, sql, params?): SQLiteRow[]`**

Execute SELECT and get all results.

```typescript
const users = SQLiteNative.queryAll(
  db,
  'SELECT * FROM users WHERE age > ?',
  [18]
);
// [{ id: 1, name: 'Alice', age: 25 }, { id: 2, name: 'Bob', age: 30 }]
```

**`SQLiteNative.execute(conn, sql, params?): number`**

Execute INSERT/UPDATE/DELETE and get affected row count.

```typescript
const affected = SQLiteNative.execute(
  db,
  'UPDATE users SET status = ? WHERE id = ?',
  ['inactive', 42]
);
console.log(`Updated ${affected} rows`);
```

---

## Phase 11.2: Connection Pool

### Overview

Manages reusable database connections with intelligent queuing, timeout handling, and statistics.

### Configuration

```typescript
interface PoolOptions {
  maxConnections?: number;        // Default: 5
  idleTimeout?: number;           // Default: 30000ms
  acquireTimeout?: number;        // Default: 10000ms
  connectionCheckInterval?: number; // Default: 5000ms
}
```

### API Methods

#### Lifecycle

**`new SQLiteConnectionPool(filename: string, options?: PoolOptions)`**

```typescript
const pool = new SQLiteConnectionPool('./app.db', {
  maxConnections: 10,           // Up to 10 concurrent connections
  idleTimeout: 60000,           // Close idle after 1 minute
  acquireTimeout: 5000,         // Wait max 5 seconds
  connectionCheckInterval: 10000 // Check idle connections every 10s
});
```

**`async pool.acquire(): Promise<SQLiteConnection>`**

Acquire connection from pool.

```typescript
const conn = await pool.acquire();  // Reuses if available, else creates new
try {
  const rows = SQLiteNative.queryAll(conn, 'SELECT * FROM users');
} finally {
  pool.release(conn);  // Important: always release
}
```

**`pool.release(conn: SQLiteConnection): void`**

Return connection to pool for reuse.

```typescript
pool.release(conn);  // Now available for other requests
```

**`async pool.close(): Promise<void>`**

Shutdown pool and close all connections.

```typescript
await pool.close();  // Wait for all operations to complete
```

#### Monitoring

**`pool.getStats(): PoolStats`**

```typescript
interface PoolStats {
  totalCreated: number;        // Total connections ever created
  activeConnections: number;   // Currently in use
  idleConnections: number;     // Available in pool
  waitingRequests: number;     // Queued requests
  totalAcquired: number;       // Total acquire() calls
  totalReleased: number;       // Total release() calls
  idleTimeouts: number;        // Closed due to timeout
  averageAcquireTime: number;  // Average ms to acquire
}
```

```typescript
const stats = pool.getStats();
console.log(`Active: ${stats.activeConnections}/${maxConnections}`);
console.log(`Waiting: ${stats.waitingRequests}`);
console.log(`Avg acquire time: ${stats.averageAcquireTime}ms`);
```

**`pool.getHealthStatus(): HealthStatus`**

```typescript
interface HealthStatus {
  isHealthy: boolean;
  activeConnections: number;
  availableConnections: number;
  waitingRequests: number;
}
```

```typescript
const health = pool.getHealthStatus();
if (!health.isHealthy) console.warn('Pool health issue detected');
```

**`pool.resetStats(): void`**

Reset all counters without affecting current connections.

```typescript
pool.resetStats();  // Clear metrics for new measurement period
```

---

## Phase 11.3: Query Builder

### Overview

Type-safe SQL query construction with automatic parameter binding and SQL injection prevention.

### Enums

```freelang
enum QueryType { SELECT, INSERT, UPDATE, DELETE }
enum ComparisonOp { EQ, NEQ, GT, GTE, LT, LTE, LIKE, IN, BETWEEN }
enum LogicalOp { AND, OR }
```

### API Functions

#### Initialization

**`selectQuery(columns: [string]): QueryBuilder`**

```freelang
let builder = selectQuery(["id", "name", "email"])
```

#### Building Clauses

**`from(builder: QueryBuilder, table: string): QueryBuilder`**

```freelang
let b = from(builder, "users")
```

**`where(builder: QueryBuilder, column: string, value: any): QueryBuilder`**

Equals comparison (=)

```freelang
let b = where(builder, "status", "active")  // status = ?
```

**`whereOp(builder: QueryBuilder, column: string, op: ComparisonOp, value: any): QueryBuilder`**

Custom comparison operator.

```freelang
let b = whereOp(builder, "age", ComparisonOp.GT, 18)     // age > ?
let b = whereOp(builder, "name", ComparisonOp.LIKE, "%alice%")
let b = whereOp(builder, "id", ComparisonOp.IN, [1, 2, 3])
```

**`orWhere(builder: QueryBuilder, column: string, value: any): QueryBuilder`**

Add OR condition.

```freelang
let b = orWhere(builder, "premium", true)  // ... OR premium = ?
```

**`orderBy(builder: QueryBuilder, column: string, ascending: boolean): QueryBuilder`**

```freelang
let b = orderBy(builder, "created_at", false)  // ORDER BY created_at DESC
```

**`limit(builder: QueryBuilder, limitValue: i32): QueryBuilder`**

```freelang
let b = limit(builder, 10)  // LIMIT 10
```

**`offset(builder: QueryBuilder, offsetValue: i32): QueryBuilder`**

```freelang
let b = offset(builder, 20)  // OFFSET 20
```

**`distinct(builder: QueryBuilder): QueryBuilder`**

```freelang
let b = distinct(builder)  // SELECT DISTINCT
```

#### SQL Generation

**`buildSQL(builder: QueryBuilder): string`**

```freelang
let sql = buildSQL(builder)
// "SELECT id, name, email FROM users WHERE status = ? ORDER BY created_at DESC LIMIT 10"
```

**`getParameters(builder: QueryBuilder): [any]`**

Extract parameters in order.

```freelang
let params = getParameters(builder)
// ["active"]
```

### Examples

**Simple SELECT**

```freelang
let builder = selectQuery(["id", "name"])
let b = from(builder, "users")
let b2 = where(b, "age", 18)
let b3 = orderBy(b2, "name", true)
let b4 = limit(b3, 10)

let sql = buildSQL(b4)        // "SELECT id, name FROM users WHERE age = ? ORDER BY name LIMIT 10"
let params = getParameters(b4) // [18]
```

**Complex Query**

```freelang
let builder = selectQuery(["id", "name", "email"])
let b = from(builder, "users")
let b2 = whereOp(b, "age", ComparisonOp.GT, 21)
let b3 = orWhere(b2, "premium", true)
let b4 = orderBy(b3, "created_at", false)
let b5 = limit(b4, 20)
let b6 = offset(b5, 40)

// "SELECT id, name, email FROM users WHERE age > ? OR premium = ? ORDER BY created_at DESC LIMIT 20 OFFSET 40"
```

---

## Phase 11.4: Performance Cache

### Overview

LRU query result caching with TTL expiration and pattern-based invalidation.

### Data Types

```freelang
type CacheEntry = {
  key: string
  value: any
  createdAt: i64
  ttl: i64
  accessCount: i32
}

type CacheConfig = {
  maxSize: i32
  ttlMillis: i64
  cleanupInterval: i64
}

type CacheStats = {
  totalEntries: i32
  maxSize: i32
  hits: i32
  misses: i32
  evictions: i32
  expirations: i32
  averageAccessCount: f64
  memorySizeMB: f64
}
```

### API Functions

#### Lifecycle

**`cacheNew(config: CacheConfig): QueryCache`**

```freelang
let cache = cacheNew({
  maxSize: 1000,
  ttlMillis: 60000,      // 1 minute TTL
  cleanupInterval: 5000  // Check every 5 seconds
})
```

**`cacheClear(cache: QueryCache): void`**

Remove all entries.

```freelang
cacheClear(cache)
```

#### Operations

**`cacheGet(cache: QueryCache, key: string): any`**

Get value (tracks hits/misses).

```freelang
let value = cacheGet(cache, "SELECT * FROM users")
if value == null {
  // Cache miss - execute query and cache result
}
```

**`cacheSet(cache: QueryCache, key: string, value: any, ttl: i64): boolean`**

Set or update value.

```freelang
cacheSet(cache, "SELECT * FROM users", rows, 60000)
```

**`cacheInvalidate(cache: QueryCache, pattern: string): i32`**

Remove entries matching pattern (substring match).

```freelang
let removed = cacheInvalidate(cache, "SELECT")  // Remove all SELECT queries
let removed = cacheInvalidate(cache, "user:")   // Remove user-related queries
```

#### Statistics

**`cacheGetStats(cache: QueryCache): CacheStats`**

```freelang
let stats = cacheGetStats(cache)
let hitRate = cacheGetHitRate(cache)     // Percentage
let evictionRate = cacheGetEvictionRate(cache)
```

**`cacheResetStats(cache: QueryCache): void`**

Reset counters without clearing data.

```freelang
cacheResetStats(cache)
```

---

## Integration Patterns

### Pattern 1: Simple Query with Cache

```typescript
import SQLiteNative from './sqlite-native';
import { cacheNew, cacheGet, cacheSet } from './query-cache';

const cache = cacheNew({ maxSize: 100, ttlMillis: 60000, cleanupInterval: 5000 });
const db = SQLiteNative.openDatabase(':memory:');

function getCachedUsers() {
  const cacheKey = 'SELECT * FROM users';
  let users = cacheGet(cache, cacheKey);

  if (users === null) {
    // Cache miss - query database
    users = SQLiteNative.queryAll(db, 'SELECT * FROM users');
    cacheSet(cache, cacheKey, users, 60000);
  }

  return users;
}

getCachedUsers();  // Hits database
getCachedUsers();  // Hits cache (50-100x faster)
```

### Pattern 2: Query Builder with Cache and Pool

```typescript
import SQLiteConnectionPool from './sqlite-connection-pool';
import { selectQuery, from, where, buildSQL, getParameters } from './query-builder';
import { cacheNew, cacheGet, cacheSet } from './query-cache';

const pool = new SQLiteConnectionPool('./app.db', { maxConnections: 5 });
const cache = cacheNew({ maxSize: 500, ttlMillis: 60000, cleanupInterval: 5000 });

async function getUsersOver18() {
  // Build query
  let builder = selectQuery(['id', 'name', 'email']);
  let b = from(builder, 'users');
  let b2 = where(b, 'age', 18);

  const sql = buildSQL(b2);
  const params = getParameters(b2);
  const cacheKey = sql + JSON.stringify(params);

  // Check cache
  let result = cacheGet(cache, cacheKey);
  if (result !== null) {
    return result;
  }

  // Get connection from pool
  const conn = await pool.acquire();
  try {
    result = SQLiteNative.queryAll(conn, sql, params);
    cacheSet(cache, cacheKey, result, 60000);
    return result;
  } finally {
    pool.release(conn);
  }
}

// Usage
const users = await getUsersOver18();  // Hits database first time
const users2 = await getUsersOver18(); // Hits cache, 50-100x faster
```

### Pattern 3: Cache Invalidation on Write

```typescript
async function updateUserStatus(userId: number, status: string) {
  const conn = await pool.acquire();
  try {
    const affected = SQLiteNative.execute(
      conn,
      'UPDATE users SET status = ? WHERE id = ?',
      [status, userId]
    );

    // Invalidate related cache entries
    if (affected > 0) {
      cacheInvalidate(cache, 'SELECT');  // Invalidate all SELECT queries
      cacheInvalidate(cache, `user:${userId}`);  // Invalidate user-specific
    }

    return affected;
  } finally {
    pool.release(conn);
  }
}
```

---

## Performance Tuning

### Cache Configuration

| Use Case | Config |
|----------|--------|
| Small app (dev) | `maxSize: 100, ttl: 60s, interval: 5s` |
| Medium app | `maxSize: 1000, ttl: 300s, interval: 10s` |
| Large app | `maxSize: 10000, ttl: 600s, interval: 30s` |
| High concurrency | `maxSize: 5000, ttl: 120s, interval: 5s` |

### Pool Configuration

| Use Case | Config |
|----------|--------|
| Single-threaded | `maxConnections: 1-2` |
| Web app (10 concurrent) | `maxConnections: 5` |
| Web app (100 concurrent) | `maxConnections: 20-30` |
| Data pipeline | `maxConnections: 50+` |

### Optimization Tips

1. **Cache Hot Queries**: Cache frequently-repeated queries
2. **Pattern Invalidation**: Invalidate groups instead of individual entries
3. **TTL Tuning**: Balance freshness vs performance
4. **Pool Sizing**: Monitor `waitingRequests` in stats
5. **Memory Monitoring**: Check `memorySizeMB` in cache stats

---

## Deployment

### Production Checklist

- [ ] Set reasonable `maxSize` for cache (avoid unbounded growth)
- [ ] Configure `idleTimeout` for connections (30-60 seconds typical)
- [ ] Set `acquireTimeout` (5-10 seconds typical)
- [ ] Enable monitoring of pool statistics
- [ ] Set up cache invalidation strategy for writes
- [ ] Test with expected concurrency level
- [ ] Monitor hit rate and adjust TTL if needed
- [ ] Set backup/recovery strategy for database file
- [ ] Configure logging for errors and timeouts

### Monitoring

```typescript
setInterval(() => {
  const stats = pool.getStats();
  const health = pool.getHealthStatus();
  const cacheStats = cacheGetStats(cache);

  console.log({
    poolHealth: health.isHealthy,
    activeConnections: stats.activeConnections,
    waitingRequests: stats.waitingRequests,
    cacheHitRate: cacheGetHitRate(cache),
    cacheEntries: cacheStats.totalEntries
  });
}, 60000);  // Every 60 seconds
```

### Scaling

- **Vertical**: Increase `maxConnections`, `maxSize`
- **Horizontal**: Run multiple processes with separate pools
- **Distributed**: Use shared cache (Redis) for multi-process

---

**Version**: Phase 11.5
**Last Updated**: 2026-03-14
