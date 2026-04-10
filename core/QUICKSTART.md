# QUICKSTART - WACAST Core Setup

## Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Docker (optional, untuk PostgreSQL)

## Step 1: Setup Environment

```bash
cd core

# Copy .env template
cp .env.example .env

# Edit .env with your credentials (or use defaults for local dev)
```

## Step 2: Start PostgreSQL

### Option A: Using Docker

```bash
docker-compose up -d postgres
```

This will:

- Start PostgreSQL container
- Load initial schema dari `db/wecast.sql`
- Port: 5432

### Option B: Using Local PostgreSQL

```bash
# Create database
createdb wacast

# Import schema
psql -U postgres -d wacast -f ../db/wecast.sql
```

## Step 3: Install Dependencies

```bash
go mod download
```

## Step 4: Run Core Service

```bash
go run main.go
```

Expected output:

```
{"level":"info","timestamp":"2026-04-09T...","message":"Starting WACAST Core","service":"wacast-core",...}
{"level":"info","timestamp":"2026-04-09T...","message":"Database connected successfully",...}
{"level":"info","timestamp":"2026-04-09T...","message":"WACAST Core started successfully"}
```

## Step 5: Verify Setup

Check database stats:

```bash
psql -U postgres -d wacast -c "\dt"
```

Should show all tables:

- api_logs
- auto_response
- billing_plans
- broadcast_campaigns
- broadcast_messages
- broadcast_recipients
- contact
- devices
- groups
- lookup
- messages
- subscriptions
- system_settings
- users
- warming_pool
- warming_sessions
- webhooks

## Step 6: View Logs

The application uses structured logging (JSON format).

Filter by level:

```bash
# Watch info logs only
go run main.go | grep "info"

# Watch error logs only
go run main.go | grep "error"
```

## Configuration Details

### Environment Variables

| Variable                 | Default     | Description                |
| ------------------------ | ----------- | -------------------------- |
| ENVIRONMENT              | development | Environment mode           |
| SERVER_PORT              | 8080        | Server port                |
| DB_HOST                  | localhost   | Database host              |
| DB_PORT                  | 5432        | Database port              |
| DB_USER                  | postgres    | Database user              |
| DB_PASSWORD              | 123456      | Database password          |
| DB_NAME                  | wacast      | Database name              |
| LOG_LEVEL                | debug       | Logging level              |
| WHATSAPP_SESSION_TIMEOUT | 300         | Session timeout in seconds |

### Database Connection Pool

- Max Open Connections: 25
- Max Idle Connections: 5
- Connection Max Lifetime: 5 minutes
- Connection Max Idle Time: 2 minutes
- Connection Timeout: 5 seconds

Adjust in `.env` using:

- `DB_MAX_OPEN_CONNS`
- `DB_MAX_IDLE_CONNS`
- `DB_CONN_MAX_LIFETIME`
- `DB_CONN_MAX_IDLE_TIME`
- `DB_CONNECTION_TIMEOUT`

## Troubleshooting

### PostgreSQL Connection Failed

```
error: "failed to ping database"
```

**Solution:**

1. Check PostgreSQL is running: `psql -U postgres -c "SELECT 1"`
2. Verify credentials in `.env`
3. Check DB_HOST and DB_PORT

### Database Schema Not Found

```
schema "public" does not exist
```

**Solution:**

```bash
# Reimport schema
psql -U postgres -d wacast -f ../db/wecast.sql
```

### Port Already in Use

```
listen tcp :8080: bind: address already in use
```

**Solution:**

```bash
# Change SERVER_PORT in .env
SERVER_PORT=8081

# Or kill process on port 8080
lsof -ti:8080 | xargs kill -9
```

## Next Steps

After successful setup:

1. ✅ Config & Environment Setup
2. ✅ Database Connection
3. ✅ Create Models (17 entities with full CRUD)
4. ✅ Create Database Queries (200+ functions)
5. ⬜ Setup Migrations Runner
6. ⬜ Build Session Service (WhatsApp)
7. ⬜ Build Message Service
8. ⬜ Create HTTP Server & Handlers

## Models & Database Queries

All models and database query functions are ready!

**Models** (`/models/`):

- User management
- Device/Session management
- Messages handling
- Subscriptions & billing
- Broadcast campaigns
- Contacts & groups
- Account warming
- Webhooks, auto-response, settings

**Database Queries** (`/database/`):

- 200+ CRUD functions
- Pagination support
- Status filtering
- Multi-entity operations

See [MODELS_GUIDE.md](./MODELS_GUIDE.md) for detailed documentation.
