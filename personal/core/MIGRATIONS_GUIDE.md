# Database Migrations Guide

## Overview

Sistem migrations untuk WACAST memungkinkan version control database schema dengan rollout yang aman.

**Migration Runner** (`database/migration.go`) menangani:

- Load migration files dari disk
- Track applied migrations di database
- Execute pending migrations
- Rollback support (future)

## How It Works

### Flow

```
Application Start
    ↓
Load Config → Connect to DB
    ↓
Initialize Migration Table (migrations)
    ↓
Load Migration Files from ./migrations/
    ↓
Check Applied Migrations
    ↓
Execute Pending Migrations → Record in migrations table
    ↓
Continue Application Startup
```

### Migration Files

**Location:** `core/migrations/`

**Naming Convention:** `XXX_migration_description.sql`

```
001_initial_schema.sql      ← First migration (creates all tables)
002_add_indexes.sql          ← Second migration (adds performance indexes)
003_add_constraints.sql      ← Future migration
...
```

**Version Sorting:** Migrations run in version order (001, 002, 003, etc.)

## Current Migrations

### 001_initial_schema.sql

Creates all base tables with relationships:

- `users` - User accounts
- `devices` - WhatsApp sessions
- `messages` - Message logs
- `groups` & `contact` - Contact management
- `billing_plans` & `subscriptions` - Billing
- `broadcast_campaigns`, `broadcast_messages`, `broadcast_recipients` - Broadcasting
- `warming_pool` & `warming_sessions` - Account warming
- `auto_response` - Auto responder
- `webhooks` - Webhook configurations
- `api_logs` - API request logging
- `system_settings` & `lookup` - System configuration

**Foreign Keys Added:**

```
subscriptions → users, billing_plans
devices → users
messages → devices
groups → users
contact → groups
broadcast_campaigns → users, devices
broadcast_messages → broadcast_campaigns
broadcast_recipients → broadcast_campaigns, groups, contact
auto_response → devices
warming_pool → devices
warming_sessions → devices
webhooks → devices
api_logs → users, devices
```

### 002_add_indexes.sql

Creates performance indexes untuk common queries:

**User Queries:**

- idx_users_phone
- idx_users_is_ban
- idx_users_is_verify

**Device Queries:**

- idx_devices_user_id
- idx_devices_phone
- idx_devices_status
- idx_devices_last_seen

**Message Queries:**

- idx_messages_device_id
- idx_messages_created_at
- idx_messages_status
- idx_messages_direction
- idx_messages_receipt

**Broadcast Queries:**

- idx_broadcast_campaigns_user_id
- idx_broadcast_campaigns_device_id
- idx_broadcast_campaigns_status
- idx_broadcast_recipients_campaign_id
- idx_broadcast_recipients_status

**Other Indexes:**

- groups, contacts, auto_response, warming_pool, warming_sessions, webhooks, api_logs, subscriptions

**Total: 40+ indexes** untuk database performance

## Migration Tracking

### migrations Table

```sql
CREATE TABLE migrations (
    id SERIAL PRIMARY KEY,
    version VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    applied_at TIMESTAMPTZ NOT NULL
)
```

**Tracks:**

- Migration version (e.g., "001")
- Migration name (e.g., "initial_schema")
- When applied (timestamp)

**Example:**

```
| id | version | name              | applied_at              |
|----|---------|-------------------|------------------------|
| 1  | 001     | initial_schema    | 2026-04-09 10:30:00+00 |
| 2  | 002     | add_indexes       | 2026-04-09 10:30:15+00 |
```

## Running Migrations

### Automatic (On Startup)

Migrations run automatically saat aplikasi start:

```bash
go run main.go
```

**Output:**

```
{"level":"info",...,"message":"Running database migrations"}
{"level":"info",...,"message":"Migrations table initialized"}
{"level":"debug",...,"message":"Loaded migration","version":"001"}
{"level":"debug",...,"message":"Loaded migration","version":"002"}
{"level":"info",...,"message":"Running migration","version":"001"}
{"level":"info",...,"message":"Migration applied successfully","version":"001"}
{"level":"info",...,"message":"Running migration","version":"002"}
{"level":"info",...,"message":"Migration applied successfully","version":"002"}
{"level":"info",...,"message":"Applied migrations:","count":2}
{"level":"info",...,"message":"  - 001","applied_at":"2026-04-09T10:30:00Z"}
{"level":"info",...,"message":"  - 002","applied_at":"2026-04-09T10:30:15Z"}
{"level":"info",...,"message":"Database migrations completed successfully"}
```

### Manual (Programmatic)

```go
import "wacast/core/database"

// Create migration runner
runner := database.NewMigrationRunner(db)

// Load migrations from directory
err := runner.LoadMigrationsFromDirectory("./migrations")

// Run migrations
err = runner.RunMigrations()

// Print status
err = runner.PrintMigrationStatus()
```

## Creating New Migrations

### Step 1: Create SQL File

Create file di `migrations/` dengan naming convention:

```bash
# 003_add_user_settings.sql
```

**Template:**

```sql
-- Migration: 003_add_user_settings
-- Description: Add user settings table and relationships

-- Create table
CREATE TABLE IF NOT EXISTS user_settings (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    theme VARCHAR(50),
    language VARCHAR(10),
    created_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Add indexes
CREATE INDEX idx_user_settings_user_id ON user_settings(user_id);

-- Migration complete
```

### Step 2: Naming Convention

**Important:** Version number must be unique and sortable

```
✓ Correct:
  - 001_initial_schema.sql
  - 002_add_indexes.sql
  - 003_add_constraints.sql
  - 010_add_new_feature.sql    (still sorts after 003)

✗ Wrong:
  - 01_initial_schema.sql      (doesn't sort properly)
  - 2_add_indexes.sql          (doesn't sort properly)
  - add_feature.sql            (no version number)
```

### Step 3: Write Safe SQL

**Best Practices:**

```sql
-- Use IF NOT EXISTS to prevent errors on re-runs
CREATE TABLE IF NOT EXISTS new_table (...);

-- Drop old columns safely
ALTER TABLE existing_table DROP COLUMN IF EXISTS old_column;

-- Use CASCADE only when necessary
DROP INDEX IF EXISTS old_index CASCADE;

-- Add constraints safely
ALTER TABLE table_name ADD CONSTRAINT constraint_name ...;

-- Rollback-friendly: separate concerns
-- ✓ Good: One feature per migration
-- ✗ Bad: Multiple unrelated changes in one migration
```

**Common Patterns:**

```sql
-- Add column
ALTER TABLE users ADD COLUMN IF NOT EXISTS extra_field VARCHAR(255);

-- Rename column
ALTER TABLE users RENAME COLUMN old_name TO new_name;

-- Add index
CREATE INDEX IF NOT EXISTS idx_name ON table_name(column_name);

-- Add constraint
ALTER TABLE table_name ADD CONSTRAINT fk_name
    FOREIGN KEY (column) REFERENCES other_table(id);
```

### Step 4: Test

```bash
# Create test database
createdb wacast_test

# Run migration
go run migrate.go
# or start the app
go run main.go

# Verify schema
psql -d wacast_test -c "\dt"
```

## Verification Checklist

Before committing migration:

- [ ] File named with unique version (e.g., `003_feature.sql`)
- [ ] SQL is idempotent (can run multiple times safely)
- [ ] Uses `IF NOT EXISTS` / `IF <condition>`
- [ ] Foreign keys are correct
- [ ] Indexes are named properly
- [ ] Comments explain what migration does
- [ ] Migration tested locally
- [ ] No hardcoded values (use params if needed)
- [ ] Handles existing data gracefully

## Troubleshooting

### Migration Not Running

**Problem:** Migration file not found

**Solution:**

```bash
# Check migrations directory exists
ls -la core/migrations/

# Check file permissions
chmod 644 core/migrations/*.sql

# Verify filename format
# Should be: XXX_name.sql
```

**Problem:** Migration fails with duplicate key error

**Solution:**

```sql
-- Check if table already exists
SELECT * FROM migrations WHERE version = '003';

-- If exists, migration already applied
-- If not, check original SQL error
```

### Database State Issues

**Problem:** Schema in inconsistent state

**Solution:**

```bash
# Connect to DB
psql -d wacast

# Check migrations table
SELECT * FROM migrations ORDER BY version;

# Manual cleanup if needed
DELETE FROM migrations WHERE version = '003';
```

**Problem:** Foreign key constraint violations

**Solution:**

```sql
-- Verify data before adding constraint
SELECT * FROM table_with_fk WHERE foreign_key_id NOT IN (
    SELECT id FROM referenced_table
);

-- Fix NULL values if needed
UPDATE table_with_fk SET foreign_key_id = valid_id WHERE foreign_key_id IS NULL;
```

## Rollback Support

Currently migrations are one-way (ups only).

**Future Implementation:** Add DOWN SQL support

```go
// In migration.go
type Migration struct {
    Version   string
    Name      string
    UpSQL     string     // Current
    DownSQL   string     // Future: for rollback
    AppliedAt *time.Time
}

// Future method
func (mr *MigrationRunner) RollbackMigration(version string) error {
    // Execute DownSQL
    // Remove from migrations table
}
```

## Performance Optimization

### Index Strategy

- **Primary Keys:** Covered in 001_initial_schema.sql via UUID PKs
- **Foreign Keys:** Covered in 001_initial_schema.sql via FOREIGN KEY constraints
- **Search Fields:** Covered in 002_add_indexes.sql
- **Timestamp Fields:** Indexed untuk historical queries (created_at, last_seen)
- **Status Fields:** Indexed untuk filtering (status, is_active, is_ban)

### Selective Indexing

Only frequently-queried columns are indexed:

- User lookup (phone, is_ban, is_verify)
- Device tracking (user_id, status, last_seen)
- Message queries (device_id, status, direction, created_at)
- Campaign progress (status, timestamps)

## Archive Migrations

Once migrations move to production, versioning becomes critical:

**Strategy:**

1. Keep all migrations in source control
2. Version control is the single source of truth
3. Never modify historical migrations
4. Always create new migrations for changes
5. Document breaking changes

Example:

```
migrations/
├── 001_initial_schema.sql       (Production v1.0)
├── 002_add_indexes.sql          (Production v1.0)
├── 003_add_constraints.sql      (Production v1.1)
├── 004_backfill_data.sql        (Production v1.1)
└── 005_add_audit_logs.sql       (Development)
```

## Integration with Services

Migrations run BEFORE services initialize:

```go
// main.go flow:
1. Load Config
2. Initialize Logger
3. Connect Database
4. Run Migrations ← Ensures schema is ready
5. Initialize Services ← Can now use tables/indexes
6. Start HTTP Server ← Ready to accept requests
```

This ensures no service tries to access non-existent tables.

## Summary

**Migrations are:**

- ✅ Versioned & tracked
- ✅ Automatic on startup
- ✅ Idempotent (safe to run multiple times)
- ✅ Sorted by version
- ✅ Logged for audit trail
- ✅ Performance-optimized with indexes

**Next:** Create more migrations as features are added!
