# MangaHub Database Documentation

---

## 1. Overview

MangaHub uses **SQLite3** as its database management system. SQLite was chosen for its simplicity, portability, and suitability for the academic project scope (50-100 concurrent users).

**Database Technology**: SQLite 3
**Driver**: `github.com/mattn/go-sqlite3`
**Location**: `./data/mangahub.db`
**ORM**: None (raw SQL using `database/sql`)

### Key Features

- **ACID Compliance**: Full transaction support
- **Foreign Key Constraints**: Referential integrity enforcement
- **Lightweight**: Single-file database, no separate server process
- **Zero Configuration**: Embedded database, works out of the box
- **Cross-Platform**: Works on Windows, Linux, macOS

### Why SQLite?

**Advantages for MangaHub:**

- ✅ Academic project requirement from specification
- ✅ Simple deployment (single file)
- ✅ No external database server needed
- ✅ Sufficient for 50-100 concurrent users
- ✅ ACID-compliant transactions
- ✅ Full SQL support with foreign keys
- ✅ Easy backup (copy single file)

**Limitations:**

- ❌ Single writer at a time (write serialization)
- ❌ Not suitable for high-concurrency writes (>100 concurrent users)
- ❌ No network access (must be on same machine)
- ❌ Limited scalability compared to PostgreSQL/MySQL

---

## 2. Database Schema

### Entity-Relationship Diagram

```
┌─────────────────────────────────────┐
│            users                    │
├─────────────────────────────────────┤
│ id (PK)              TEXT           │
│ username (UNIQUE)    TEXT           │
│ email (UNIQUE)       TEXT           │
│ password_hash        TEXT           │
│ created_at           TIMESTAMP      │
│ updated_at           TIMESTAMP      │
└────────────┬────────────────────────┘
             │
             │ 1:N
             │
      ┌──────▼──────────────────────────────┐
      │       user_progress                 │
      ├─────────────────────────────────────┤
      │ user_id (PK, FK) → users.id         │
      │ manga_id (PK, FK) → manga.id        │
      │ current_chapter      INTEGER        │
      │ status               TEXT           │
      │ rating               INTEGER        │
      │ started_at           TIMESTAMP      │
      │ completed_at         TIMESTAMP      │
      │ updated_at           TIMESTAMP      │
      └──────────┬──────────────────────────┘
                 │
                 │ N:1
                 │
  ┌──────────────▼───────────────────┐
  │         manga                    │
  ├──────────────────────────────────┤
  │ id (PK)           TEXT           │
  │ title             TEXT           │
  │ author            TEXT           │
  │ genres            TEXT (JSON)    │
  │ status            TEXT           │
  │ total_chapters    INTEGER        │
  │ description       TEXT           │
  │ cover_image_url   TEXT           │
  │ created_at        TIMESTAMP      │
  │ updated_at        TIMESTAMP      │
  └──────────────────────────────────┘
```

### Schema Summary

| Table           | Purpose                   | Rows (Expected) |
| --------------- | ------------------------- | --------------- |
| `users`         | User authentication       | 10-50           |
| `manga`         | Manga catalog/library     | 200+            |
| `user_progress` | Reading progress tracking | 500+            |

---

## 3. Table Definitions

### 3.1 `users` Table

**Purpose**: Store user authentication and profile information

```sql
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
```

**Columns:**

| Column          | Type      | Constraints      | Description                              |
| --------------- | --------- | ---------------- | ---------------------------------------- |
| `id`            | TEXT      | PRIMARY KEY      | Unique user identifier (UUID format)     |
| `username`      | TEXT      | UNIQUE, NOT NULL | Username for login (unique)              |
| `email`         | TEXT      | UNIQUE, NOT NULL | Email address (unique)                   |
| `password_hash` | TEXT      | NOT NULL         | Bcrypt hashed password (never plaintext) |
| `created_at`    | TIMESTAMP | DEFAULT NOW      | Account creation timestamp               |
| `updated_at`    | TIMESTAMP | DEFAULT NOW      | Last profile update timestamp            |

**Sample Data:**

```sql
-- Password is "password123" hashed with bcrypt
INSERT INTO users (id, username, email, password_hash)
VALUES
  ('user-testuser', 'testuser', 'testuser@example.com', '$2a$10$...'),
  ('user-alice', 'alice', 'alice@example.com', '$2a$10$...'),
  ('user-bob', 'bob', 'bob@example.com', '$2a$10$...');
```

**Indexes:**

- `idx_users_username`: Speeds up login queries (username lookups)
- `idx_users_email`: Speeds up email verification and password reset

---

### 3.2 `manga` Table

**Purpose**: Store manga catalog information

```sql
CREATE TABLE IF NOT EXISTS manga (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT,
    genres TEXT,  -- JSON array stored as TEXT: ["Action", "Adventure"]
    status TEXT CHECK(status IN ('ongoing', 'completed', 'hiatus', 'cancelled')),
    total_chapters INTEGER DEFAULT 0,
    description TEXT,
    cover_image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for search and filtering
CREATE INDEX IF NOT EXISTS idx_manga_title ON manga(title);
CREATE INDEX IF NOT EXISTS idx_manga_status ON manga(status);
```

**Columns:**

| Column            | Type      | Constraints      | Description                             |
| ----------------- | --------- | ---------------- | --------------------------------------- |
| `id`              | TEXT      | PRIMARY KEY      | Unique manga identifier (slug format)   |
| `title`           | TEXT      | NOT NULL         | Manga title                             |
| `author`          | TEXT      | -                | Author name                             |
| `genres`          | TEXT      | -                | JSON array of genres (stored as string) |
| `status`          | TEXT      | CHECK constraint | Publication status (ongoing/completed)  |
| `total_chapters`  | INTEGER   | DEFAULT 0        | Total number of chapters                |
| `description`     | TEXT      | -                | Manga synopsis/description              |
| `cover_image_url` | TEXT      | -                | URL to cover image                      |
| `created_at`      | TIMESTAMP | DEFAULT NOW      | Record creation timestamp               |
| `updated_at`      | TIMESTAMP | DEFAULT NOW      | Last update timestamp                   |

**Status Values:**

- `ongoing` - Currently being published
- `completed` - Finished series
- `hiatus` - Temporarily paused
- `cancelled` - Series cancelled

**Genres Storage:**
Genres are stored as JSON text for simplicity:

```json
["Action", "Adventure", "Shounen"]
```

**Sample Data:**

```sql
INSERT INTO manga (id, title, author, genres, status, total_chapters, description)
VALUES (
  'one-piece',
  'One Piece',
  'Eiichiro Oda',
  '["Action", "Adventure", "Comedy", "Fantasy"]',
  'ongoing',
  1150,
  'Monkey D. Luffy wants to be King of the Pirates...'
);
```

**Indexes:**

- `idx_manga_title`: Enables fast title-based searches
- `idx_manga_status`: Filters manga by publication status

---

### 3.3 `user_progress` Table

**Purpose**: Track user reading progress for each manga

```sql
CREATE TABLE IF NOT EXISTS user_progress (
    user_id TEXT NOT NULL,
    manga_id TEXT NOT NULL,
    current_chapter INTEGER DEFAULT 0,
    status TEXT CHECK(status IN ('reading', 'completed', 'plan_to_read', 'on_hold', 'dropped')),
    rating INTEGER CHECK(rating >= 1 AND rating <= 10),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, manga_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (manga_id) REFERENCES manga(id) ON DELETE CASCADE
);

-- Indexes for efficient queries
CREATE INDEX IF NOT EXISTS idx_user_progress_user_id ON user_progress(user_id);
CREATE INDEX IF NOT EXISTS idx_user_progress_manga_id ON user_progress(manga_id);
CREATE INDEX IF NOT EXISTS idx_user_progress_status ON user_progress(user_id, status);
```

**Columns:**

| Column            | Type      | Constraints                  | Description                         |
| ----------------- | --------- | ---------------------------- | ----------------------------------- |
| `user_id`         | TEXT      | PK, FK → users(id), NOT NULL | Reference to user                   |
| `manga_id`        | TEXT      | PK, FK → manga(id), NOT NULL | Reference to manga                  |
| `current_chapter` | INTEGER   | DEFAULT 0                    | Last chapter read (0 = not started) |
| `status`          | TEXT      | CHECK constraint             | Reading status                      |
| `rating`          | INTEGER   | CHECK (1-10)                 | User rating (1-10 scale, nullable)  |
| `started_at`      | TIMESTAMP | -                            | When user started reading           |
| `completed_at`    | TIMESTAMP | -                            | When user marked as completed       |
| `updated_at`      | TIMESTAMP | DEFAULT NOW                  | Last progress update                |

**Composite Primary Key**: (`user_id`, `manga_id`)

- One progress record per user per manga
- Ensures users can't have duplicate entries for same manga

**Foreign Keys:**

- `user_id → users(id)` with `ON DELETE CASCADE`
  - Deleting a user automatically deletes their progress
- `manga_id → manga(id)` with `ON DELETE CASCADE`
  - Deleting a manga removes all user progress for it

**Reading Status Values:**

- `reading` - Currently reading
- `completed` - Finished reading
- `plan_to_read` - In user's to-read list
- `on_hold` - Paused reading
- `dropped` - User stopped reading

**Sample Data:**

```sql
INSERT INTO user_progress (user_id, manga_id, current_chapter, status, rating)
VALUES
  ('user-testuser', 'one-piece', 1050, 'reading', 10),
  ('user-alice', 'naruto', 700, 'completed', 9),
  ('user-bob', 'attack-on-titan', 0, 'plan_to_read', NULL);
```

**Indexes:**

- `idx_user_progress_user_id`: Get all manga for a user (library view)
- `idx_user_progress_manga_id`: Find all users reading a manga
- `idx_user_progress_status`: Filter user's library by status

---

## 4. Database Migrations

### Migration System

MangaHub uses a custom migration system implemented in `pkg/database/migrate.go`.

**Migration Files Location**: `migrations/`

**Naming Convention**: `{version}_{description}.{direction}.sql`

- Example: `001_create_users_table.up.sql`
- Example: `001_create_users_table.down.sql`

**Migration Tracking Table**:

```sql
CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Available Migrations

| Version | Name                       | Description                 |
| ------- | -------------------------- | --------------------------- |
| 001     | create_users_table         | Creates users table         |
| 002     | create_manga_table         | Creates manga catalog table |
| 003     | create_user_progress_table | Creates progress tracking   |

### Running Migrations

**Apply all pending migrations:**

```bash
make migrate-up
```

This will:

1. Create `schema_migrations` table if not exists
2. Load all migration files from `migrations/`
3. Check which migrations have been applied
4. Apply pending migrations in order
5. Record each migration in `schema_migrations`

**Rollback last migration:**

```bash
make migrate-down
```

**Reset database (rollback all, reapply all, reseed):**

```bash
make db-reset
```

### Migration Implementation

**Applying a Migration** (`pkg/database/migrate.go:applyMigration`):

```go
func (m *Migrator) applyMigration(migration Migration) error {
    tx, err := m.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Execute migration SQL
    if _, err := tx.Exec(migration.UpSQL); err != nil {
        return err  // Transaction auto-rolls back
    }

    // Record migration
    _, err = tx.Exec(`
        INSERT INTO schema_migrations (version, name) VALUES (?, ?)
    `, migration.Version, migration.Name)
    if err != nil {
        return err
    }

    return tx.Commit()  // Commit if everything succeeded
}
```

**Transaction Safety**: Each migration runs in a transaction - either fully applies or fully rolls back.

---

## 5. Seed Data

### Seeding Process

**Seed Script**: `scripts/seed/main.go`

**Run Seeding:**

```bash
make seed
```

**What Gets Seeded:**

1. **Sample Users** (3 users):
   - `testuser` / `testuser@example.com` / password: `password123`
   - `alice` / `alice@example.com` / password: `alice123`
   - `bob` / `bob@example.com` / password: `bob123`

2. **Manga Catalog** (from `data/manga.json`):
   - 200+ manga titles
   - Loaded from JSON file
   - Includes title, author, genres, status, chapters, description, cover URL

3. **Sample User Progress** (5 entries):
   - Test data showing various reading statuses
   - Links users to manga they're reading

### Seed Data Structure

**Users Seeding** (`scripts/seed/main.go:seedUsers`):

```go
users := []struct {
    username string
    email    string
    password string
}{
    {"testuser", "testuser@example.com", "password123"},
    {"alice", "alice@example.com", "alice123"},
    {"bob", "bob@example.com", "bob123"},
}

for _, u := range users {
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
    userID := fmt.Sprintf("user-%s", u.username)

    db.Exec(`INSERT OR IGNORE INTO users (id, username, email, password_hash)
             VALUES (?, ?, ?, ?)`, userID, u.username, u.email, string(hashedPassword))
}
```

**Manga Seeding** (`scripts/seed/main.go:seedMangaFromJSON`):

```go
content, _ := os.ReadFile("data/manga.json")
var mangaList []models.Manga
json.Unmarshal(content, &mangaList)

for _, manga := range mangaList {
    genresJSON, _ := manga.MarshalGenres()  // Convert []string to JSON

    db.Exec(`INSERT OR IGNORE INTO manga (id, title, author, genres, status, total_chapters, description, cover_image_url)
             VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
             manga.ID, manga.Title, manga.Author, genresJSON, manga.Status,
             manga.TotalChapters, manga.Description, manga.CoverImageURL)
}
```

**Idempotent Seeding**: Uses `INSERT OR IGNORE` - running seed multiple times won't create duplicates.

---

## 6. Connection Configuration

### Configuration File

**Location**: `configs/dev.yaml`

```yaml
database:
  path: "./data/mangahub.db"
```

**Environment Override**:

```bash
export DB_PATH="/custom/path/mangahub.db"
```

### Connection Pool Settings

**Code**: `pkg/database/sqlite.go`

```go
type Config struct {
    Path            string
    MaxOpenConns    int           // Maximum open connections
    MaxIdleConns    int           // Maximum idle connections
    ConnMaxLifetime time.Duration // Connection lifetime
}

func DefaultConfig() Config {
    return Config{
        Path:            "./data/mangahub.db",
        MaxOpenConns:    25,   // SQLite recommendation: 1-25 for concurrency
        MaxIdleConns:    5,    // Keep 5 idle connections ready
        ConnMaxLifetime: 5 * time.Minute,
    }
}
```

**Connection Pooling**:

```go
db.SetMaxOpenConns(25)  // Limit to 25 concurrent connections
db.SetMaxIdleConns(5)   // Keep 5 idle connections in pool
db.SetConnMaxLifetime(5 * time.Minute)  // Recycle connections after 5 minutes
```

**Foreign Keys Enforcement**:

```go
// CRITICAL: SQLite doesn't enforce foreign keys by default
db.Exec("PRAGMA foreign_keys = ON")
```

Without this pragma, foreign key constraints are ignored!

### Health Check

```go
func HealthCheck(db *sql.DB) error {
    if db == nil {
        return fmt.Errorf("database connection is nil")
    }
    return db.Ping()  // Verifies connection is alive
}
```

**API Endpoint**: `GET /health` uses this to verify database health.

---

## 7. Common Queries

### User Queries

**Register New User**:

```sql
INSERT INTO users (id, username, email, password_hash)
VALUES (?, ?, ?, ?);
```

**Login (Verify Credentials)**:

```sql
SELECT id, username, password_hash
FROM users
WHERE username = ?;
-- Then verify password_hash with bcrypt.CompareHashAndPassword()
```

**Get User Profile**:

```sql
SELECT id, username, email, created_at, updated_at
FROM users
WHERE id = ?;
```

---

### Manga Queries

**Search Manga by Title**:

```sql
SELECT id, title, author, genres, status, total_chapters, cover_image_url
FROM manga
WHERE title LIKE '%' || ? || '%'
ORDER BY title
LIMIT 20;
```

**Filter by Status**:

```sql
SELECT id, title, author, status, total_chapters
FROM manga
WHERE status = 'ongoing'
ORDER BY title;
```

**Get Manga by ID**:

```sql
SELECT * FROM manga WHERE id = ?;
```

---

### User Progress Queries

**Get User's Library**:

```sql
SELECT m.id, m.title, m.author, m.cover_image_url,
       up.current_chapter, up.status, up.rating, up.updated_at
FROM user_progress up
JOIN manga m ON up.manga_id = m.id
WHERE up.user_id = ?
ORDER BY up.updated_at DESC;
```

**Get Currently Reading Manga**:

```sql
SELECT m.id, m.title, up.current_chapter, m.total_chapters
FROM user_progress up
JOIN manga m ON up.manga_id = m.id
WHERE up.user_id = ? AND up.status = 'reading'
ORDER BY up.updated_at DESC;
```

**Update Reading Progress**:

```sql
INSERT INTO user_progress (user_id, manga_id, current_chapter, status, updated_at)
VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(user_id, manga_id) DO UPDATE SET
    current_chapter = excluded.current_chapter,
    status = excluded.status,
    updated_at = CURRENT_TIMESTAMP;
```

**Get Progress for Specific Manga**:

```sql
SELECT current_chapter, status, rating, started_at, completed_at
FROM user_progress
WHERE user_id = ? AND manga_id = ?;
```

---

## 8. Transactions

### Why Use Transactions?

- **Atomicity**: Multiple operations succeed or fail together
- **Consistency**: Database always in valid state
- **Isolation**: Concurrent transactions don't interfere
- **Durability**: Committed changes persist

### Transaction Example

**Update Progress and Mark Completed**:

```go
func (r *Repository) CompleteReading(userID, mangaID string, finalChapter int, rating int) error {
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()  // Rollback if not committed

    // Update progress
    _, err = tx.Exec(`
        UPDATE user_progress
        SET current_chapter = ?, status = 'completed', rating = ?, completed_at = CURRENT_TIMESTAMP
        WHERE user_id = ? AND manga_id = ?
    `, finalChapter, rating, userID, mangaID)
    if err != nil {
        return err  // Rolls back automatically
    }

    // Additional operations...

    return tx.Commit()  // All or nothing
}
```

---

## 9. Performance Optimization

### Indexes

All critical indexes are created in migration files:

```sql
-- Users
CREATE INDEX idx_users_username ON users(username);  -- Login queries
CREATE INDEX idx_users_email ON users(email);        -- Email lookups

-- Manga
CREATE INDEX idx_manga_title ON manga(title);        -- Search queries
CREATE INDEX idx_manga_status ON manga(status);      -- Filter by status

-- User Progress
CREATE INDEX idx_user_progress_user_id ON user_progress(user_id);           -- Get user's library
CREATE INDEX idx_user_progress_manga_id ON user_progress(manga_id);         -- Manga popularity
CREATE INDEX idx_user_progress_status ON user_progress(user_id, status);    -- Filter by status
```

### Query Optimization Tips

**Use Prepared Statements**:

```go
// ✅ GOOD - Uses prepared statement (prevents SQL injection, faster)
stmt, _ := db.Prepare("SELECT * FROM manga WHERE id = ?")
defer stmt.Close()
stmt.QueryRow(mangaID)

// ❌ BAD - String concatenation (SQL injection risk, no statement caching)
db.Query("SELECT * FROM manga WHERE id = '" + mangaID + "'")
```

**Limit Result Sets**:

```sql
-- Always use LIMIT for large tables
SELECT * FROM manga ORDER BY title LIMIT 20;
```

**Use Covering Indexes** (when possible):

```sql
-- This query can be answered entirely from the index
SELECT user_id, status FROM user_progress WHERE user_id = ?;
```

### SQLite-Specific Optimizations

**Analyze Database** (update statistics for query planner):

```sql
ANALYZE;
```

**Vacuum** (reclaim space after deletions):

```sql
VACUUM;
```

**Pragmas for Performance**:

```sql
PRAGMA journal_mode = WAL;        -- Write-Ahead Logging (better concurrency)
PRAGMA synchronous = NORMAL;      -- Balance between speed and safety
PRAGMA cache_size = -64000;       -- 64MB cache
PRAGMA temp_store = MEMORY;       -- Store temp tables in memory
```

---

## 10. Backup and Restore

### Backup Database

**Method 1: File Copy** (Safest - requires downtime):

```bash
# Stop all servers first
cp data/mangahub.db data/backups/mangahub_$(date +%Y%m%d_%H%M%S).db
```

**Method 2: SQLite Backup API** (Online backup):

```bash
sqlite3 data/mangahub.db ".backup data/backups/mangahub_backup.db"
```

**Method 3: SQL Dump**:

```bash
sqlite3 data/mangahub.db .dump > mangahub_backup.sql
```

### Restore Database

**From File Backup**:

```bash
# Stop servers
cp data/backups/mangahub_20251226_120000.db data/mangahub.db
# Restart servers
```

**From SQL Dump**:

```bash
sqlite3 data/mangahub.db < mangahub_backup.sql
```

### Automated Backups

**Cron Job** (Linux/Mac):

```bash
# Add to crontab: backup every day at 2 AM
0 2 * * * cd /path/to/mangahub && sqlite3 data/mangahub.db ".backup data/backups/mangahub_$(date +\%Y\%m\%d).db"
```

---

## 11. Security Considerations

### Password Security

**Hashing Algorithm**: bcrypt (cost factor 10)

```go
import "golang.org/x/crypto/bcrypt"

// Hash password before storing
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
db.Exec("INSERT INTO users (..., password_hash) VALUES (..., ?)", string(hashedPassword))

// Verify password during login
err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(providedPassword))
if err == nil {
    // Password correct
}
```

**Never store plaintext passwords!**

### SQL Injection Prevention

**Always use parameterized queries**:

```go
// ✅ SAFE - Uses prepared statement
db.Query("SELECT * FROM users WHERE username = ?", username)

// ❌ VULNERABLE - String concatenation
db.Query("SELECT * FROM users WHERE username = '" + username + "'")
// Attack: username = "'; DROP TABLE users; --"
```

### Access Control

- Database file permissions: `chmod 600 data/mangahub.db` (owner read/write only)
- No remote database access (SQLite is local-only)
- Application-level authorization via JWT middleware

### Data Privacy

**Consider encrypting sensitive data**:

- Email addresses (PII - Personally Identifiable Information)
- Reading history (user privacy)

**GDPR Compliance** (if deployed in EU):

- Implement user data export
- Implement account deletion (CASCADE handles progress deletion)

---

## 12. Troubleshooting

### Database Locked Error

```
Error: database is locked
```

**Cause**: Another process has an exclusive lock (SQLite allows only one writer at a time)

**Solutions**:

1. Ensure `MaxOpenConns` is set appropriately (recommended: 1-25)
2. Use WAL mode: `PRAGMA journal_mode = WAL;`
3. Keep transactions short
4. Consider upgrading to PostgreSQL for high concurrency

---

### Foreign Key Constraint Violation

```
Error: FOREIGN KEY constraint failed
```

**Cause**: Trying to insert/update with invalid foreign key reference

**Debug**:

```sql
-- Check if referenced manga exists
SELECT id FROM manga WHERE id = ?;

-- Check if referenced user exists
SELECT id FROM users WHERE id = ?;
```

**Remember**: Enable foreign keys!

```sql
PRAGMA foreign_keys = ON;
```

---

### Migration Already Applied

```
Error: migration already applied
```

**Check migration status**:

```sql
SELECT * FROM schema_migrations ORDER BY version;
```

**Manual rollback**:

```bash
make migrate-down  # Rollback last migration
```

---

## 13. Upgrading to PostgreSQL (Future)

### Why Upgrade?

- More than 100 concurrent users
- High write concurrency requirements
- Network database access needed
- Advanced features (JSON columns, full-text search)

### Migration Steps

1. **Export Data from SQLite**:

   ```bash
   sqlite3 data/mangahub.db .dump > sqlite_dump.sql
   ```

2. **Convert Schema** (SQLite → PostgreSQL):
   - Change `TEXT PRIMARY KEY` → `UUID PRIMARY KEY` or `SERIAL`
   - Change `INTEGER` → `INT` or `BIGINT`
   - Replace `AUTOINCREMENT` → `SERIAL`
   - Update `TIMESTAMP DEFAULT CURRENT_TIMESTAMP` → `TIMESTAMP DEFAULT NOW()`

3. **Import to PostgreSQL**:

   ```bash
   psql -U postgres -d mangahub < postgres_schema.sql
   ```

4. **Update Driver**:

   ```go
   import _ "github.com/lib/pq"  // PostgreSQL driver
   db, err := sql.Open("postgres", "postgres://user:pass@localhost/mangahub?sslmode=disable")
   ```

5. **Test Thoroughly** before production deployment

---

## 14. References

### Internal Documentation

- [Architecture Documentation](./architecture.md) - System design
- [API Documentation](./api-documentation.md) - HTTP endpoints
- [CLAUDE.md](../CLAUDE.md) - Development guidelines

### External Resources

- [SQLite Documentation](https://www.sqlite.org/docs.html)
- [SQLite Foreign Keys](https://www.sqlite.org/foreignkeys.html)
- [SQLite WAL Mode](https://www.sqlite.org/wal.html)
- [Go database/sql Tutorial](https://go.dev/doc/database/querying)
- [Bcrypt in Go](https://pkg.go.dev/golang.org/x/crypto/bcrypt)

---

**Last Updated**: 2025-12-26
**Version**: 1.0.0
**Database Version**: 3 (latest migration: 003_create_user_progress_table)
**Status**: ✅ Production Ready
