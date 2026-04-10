# WACAST Core - WhatsApp Gateway

WhatsApp gateway dengan multi-session support menggunakan Go & PostgreSQL.

## Quick Start

```bash
# 1. Setup environment
cp .env.example .env
# Edit .env and set ENCRYPTION_KEY to a 64-character hex string

# 2. Start PostgreSQL (with Docker)
docker-compose up -d postgres

# 3. Install dependencies & download modules
go mod download

# 4. Run core service
go run main.go
```

See [QUICKSTART.md](./QUICKSTART.md) for detailed setup instructions.

## Phase Status

| Phase | Name             | Status      | Files | LOC  |
| ----- | ---------------- | ----------- | ----- | ---- |
| 1     | Foundation       | ✅ Complete | 4     | 200  |
| 2     | Models           | ✅ Complete | 8     | 600  |
| 3     | Database Queries | ✅ Complete | 9     | 1200 |
| 4     | Migrations       | ✅ Complete | 3     | 600  |
| 5     | Session Service  | ✅ Complete | 5     | 800  |
| 6     | Message Service  | ✅ Complete | 4     | 1000 |
| 7     | HTTP Server      | ✅ Complete | 2     | 280  |
| 8     | Webhooks         | ⏳ Pending  | -     | -    |

**Total (Phases 1-7):** 56 files, 6680+ lines of code

---

## Architecture Overview

```
WACAST Core - Layered Architecture:

┌─────────────────────────────────────────┐
│  HTTP API Layer (Phase 7)               │
│  ├─ Session Endpoints                   │
│  ├─ Message Endpoints                   │
│  └─ Device Endpoints                    │
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│  Handlers (Phase 5+)                    │
│  ├─ session_handler.go                  │
│  ├─ message_handler.go                  │
│  └─ device_handler.go                   │
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│  Services (Phase 5+)                    │
│  ├─ session/service.go (Multi-session) │
│  ├─ message/service.go (In Progress)   │
│  └─ device/service.go (Planned)        │
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│  Models & Queries (Phase 2-3)           │
│  ├─ models/ (17 types)                  │
│  ├─ database/queries (200+ functions)   │
│  └─ models validation/response helpers  │
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│  Database Layer (Phase 1, 4)            │
│  ├─ PostgreSQL 15+ with lib/pq          │
│  ├─ Connection pooling (25 max)         │
│  ├─ Migrations system (auto-run)        │
│  └─ Health check & transactions         │
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│  Infrastructure (Phase 1)               │
│  ├─ Config (environment variables)      │
│  ├─ Logger (Zap JSON structured)        │
│  ├─ Encryption (AES-256-GCM)            │
│  └─ Utils & helpers                     │
└─────────────────────────────────────────┘
```

## Services Breakdown

### Session Service (Phase 5) ✅

**File:** `services/session/`

Multi-session WhatsApp management:

- Up to 25 concurrent sessions per instance
- QR code generation and scanning
- Connection state tracking
- Session persistence (encrypted in DB)
- Auto-restore on application restart
- Event handlers for QR, connection, messages
- Endpoints: `/sessions` RESTful API

**Key Methods:**

```go
StartSession()           // Initiate new WhatsApp connection
StopSession()           // Gracefully disconnect
RestoreSession()        // Resume from encrypted data
SendMessage()           // Send text message
RestorePreviousSessions() // Auto-restore on startup
```

### Message Service (Phase 6) ✅

**File:** `services/message/`

Reliable message delivery with status tracking:

- Queue message outgoing/incoming
- Track delivery status (pending→sent→delivered→read)
- Retry failed messages (exponential backoff)
- Handle concurrent sends (max 5 per device)
- Scheduled message support
- Error logging and diagnostics
- Webhook integration ready
- Endpoints: `/devices/:device_id/messages` RESTful API

**Key Methods:**

```go
SendMessage()              // Queue text message
SendMessageWithMedia()     // Queue with attachment
SendScheduledMessage()     // Queue for future
ReceiveMessage()          // Process incoming
GetMessageStatus()        // Query delivery status
ProcessQueue()            // Manual queue processing
```

### Webhook Service (Phase 8) ⏳

**File:** `services/webhook/` (planned)

Push incoming events to API layer:

- Message received notifications
- Status updates
- Connection changes
- Device events

---

## Directory Structure

```
core/
├── config/                  # Configuration & constants
│   ├── config.go           # Main config loader
│   ├── database.go         # Database configuration
│   └── constants.go        # 30+ global constants
│
├── database/               # Data access layer
│   ├── db.go              # Connection pool, transactions
│   ├── migration.go       # Migration runner
│   ├── user_queries.go
│   ├── device_queries.go
│   ├── message_queries.go
│   ├── broadcast_queries.go
│   ├── contact_queries.go
│   ├── subscription_queries.go
│   ├── warming_queries.go
│   └── other_queries.go
│
├── models/                 # Type-safe data structures (17 types)
│   ├── user.go
│   ├── device.go
│   ├── message.go
│   ├── subscription.go
│   ├── contact.go
│   ├── broadcast.go
│   ├── warming.go
│   └── other.go
│
├── services/               # Business logic & integrations
│   ├── session/            # PHASE 5 - WhatsApp multi-sessions
│   │   ├── types.go       # SessionStatus, Events, SessionManager
│   │   ├── encryption.go  # AES-256-GCM encryption
│   │   ├── service.go     # Main SessionService (400 lines)
│   │   └── manager.go     # Background operations (cleanup, reconnect)
│   │
│   └── message/           # PHASE 6 - Message queue & delivery
│       ├── types.go       # MessageStatus, QueuedMessage, ReceivedMessage
│       ├── store.go       # Database message store (300 lines)
│       ├── service.go     # Main MessageService (500 lines)
│       └── (handlers in handlers/)
│
├── handlers/               # HTTP request handlers
│   ├── session_handler.go # REST endpoints for sessions (Phase 5)
│   └── message_handler.go # REST endpoints for messages (Phase 6)
│
├── migrations/             # Database migrations (idempotent)
│   ├── 001_initial_schema.sql
│   ├── 002_add_indexes.sql
│   └── README.md
│
├── utils/                  # Shared utilities
│   ├── logger.go          # Zap wrapper
│   ├── encryption.go      # AES-GCM helpers
│   └── validators.go      # Input validation
│
├── cmd/                    # CLI tools
│   └── migrate/           # Migration status checker
│       └── main.go
│
├── main.go                # Application entry point
├── go.mod                 # Module definition & dependencies
├── .env.example           # Configuration template
├── .gitignore
├── Dockerfile             # Multi-stage build
├── docker-compose.yml     # PostgreSQL + core service
│
├── Documentation/
│   ├── README.md                    # This file
│   ├── QUICKSTART.md                # Setup guide
│   ├── SETUP_COMPLETE.md            # Phase 1-4 summary
│   ├── MODELS_GUIDE.md              # Model documentation
│   ├── MIGRATIONS_GUIDE.md          # Migration system
│   ├── SESSION_SERVICE_GUIDE.md     # Phase 5 documentation
│   ├── MESSAGE_SERVICE_GUIDE.md     # Phase 6 documentation
│   ├── PHASE_5_COMPLETE.md          # Phase 5 completion report
│   ├── PHASE_6_COMPLETE.md          # Phase 6 completion report
│   ├── DEVELOPMENT_STATUS.md        # Full status report
│   └── sample-requests.http         # Example Requests
```

## Features by Phase

### Phase 1-4: Foundation ✅

- Configuration management
- PostgreSQL connection pooling
- Structured logging (Zap)
- 17 database tables with relationships
- Idempotent migrations
- 200+ CRUD query functions
- AES-256-GCM encryption ready

### Phase 5: Session Service ✅

- ✅ Multi-session management (up to 25 concurrent)
- ✅ WhatsApp integration via `go.mau.fi/whatsmeow`
- ✅ QR code generation & scanning
- ✅ Session data encryption (AES-256-GCM)
- ✅ Auto-restore on app restart
- ✅ Event-driven callback system
- ✅ REST API for session control
- ✅ Background cleanup & reconnect manager

### Phase 6: Message Service ✅

- ✅ Message queuing with batch processing
- ✅ Automatic retry with exponential backoff
- ✅ Status tracking (pending→sent→delivered→read)
- ✅ Concurrent send limiting (max 5 per device)
- ✅ Incoming message handling
- ✅ Failed message recovery
- ✅ Scheduled message support
- ✅ Media attachment framework
- ✅ REST API for message operations
- ✅ Queue statistics & monitoring

### Phase 7-8: In Progress

- HTTP server with Gin framework
- Webhook notification system
- Auto-response & warming
- Broadcast campaigns

---

## Key Configuration

```env
# WhatsApp Session
WHATSAPP_SESSION_TIMEOUT=300    # Seconds before idle disconnection
ENCRYPTION_KEY=<64-hex-chars>   # AES-256 key for session data

# Database
DB_MAX_OPEN_CONNS=25            # Connection pool size
DB_CONN_MAX_LIFETIME=5          # Minutes

# Logging
LOG_LEVEL=debug                 # debug|info|warn|error
```

---

## API Examples

### Session Management

```bash
# Get all active sessions
curl http://localhost:8080/sessions

# Get session status
curl http://localhost:8080/sessions/device-001

# Initiate new session (returns QR code event)
curl -X POST http://localhost:8080/sessions/initiate \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "device-001",
    "user_id": "user-123",
    "phone": "62812345678"
  }'

# Send message (after session is authenticated)
curl -X POST http://localhost:8080/sessions/device-001/messages \
  -H "Content-Type: application/json" \
  -d '{
    "target_jid": "62812345678@s.whatsapp.net",
    "message": "Hello World",
    "is_group": false
  }'

# Stop session
curl -X POST http://localhost:8080/sessions/device-001/stop
```

### Message Management

```bash
# Send text message
curl -X POST http://localhost:8080/devices/device-001/messages \
  -H "Content-Type: application/json" \
  -d '{
    "target_jid": "62812345678@s.whatsapp.net",
    "content": "Hello World"
  }'

# Check message status
curl http://localhost:8080/messages/msg-123/status

# Send media message
curl -X POST http://localhost:8080/devices/device-001/messages/media \
  -H "Content-Type: application/json" \
  -d '{
    "target_jid": "62812345678@s.whatsapp.net",
    "media_url": "https://example.com/image.jpg",
    "content_type": "image"
  }'

# Schedule message
curl -X POST http://localhost:8080/devices/device-001/messages/scheduled \
  -H "Content-Type: application/json" \
  -d '{
    "target_jid": "62812345678@s.whatsapp.net",
    "content": "Tomorrow message",
    "scheduled_for": "2026-04-10T10:00:00Z"
  }'

# Get queue statistics
curl http://localhost:8080/messages/stats

# Get failed messages
curl http://localhost:8080/messages/failed?limit=50
```

---

## Documentation Files

1. **[QUICKSTART.md](./QUICKSTART.md)** - Setup, Docker, database initialization
2. **[SETUP_COMPLETE.md](./SETUP_COMPLETE.md)** - Foundation summary (Phases 1-4)
3. **[MODELS_GUIDE.md](./MODELS_GUIDE.md)** - All 17 data models documented
4. **[MIGRATIONS_GUIDE.md](./MIGRATIONS_GUIDE.md)** - Migration system & usage
5. **[SESSION_SERVICE_GUIDE.md](./SESSION_SERVICE_GUIDE.md)** - Phase 5 deep dive (400+ lines)
6. **[MESSAGE_SERVICE_GUIDE.md](./MESSAGE_SERVICE_GUIDE.md)** - Phase 6 deep dive (600+ lines)
7. **[PHASE_5_COMPLETE.md](./PHASE_5_COMPLETE.md)** - Phase 5 completion report
8. **[PHASE_6_COMPLETE.md](./PHASE_6_COMPLETE.md)** - Phase 6 completion report
9. **[DEVELOPMENT_STATUS.md](./DEVELOPMENT_STATUS.md)** - Complete project status

---

## Development Workflow

### 1. Local Development

```bash
# Watch mode for auto-rebuild
# (Optional - install air for watch mode)
go install github.com/cosmtrek/air@latest
air

# Or manual rebuild
go run main.go
```

### 2. Testing Session Service

```bash
# Start service
go run main.go

# In another terminal, test endpoints
curl http://localhost:8080/sessions
```

### 3. Database Inspection

```bash
# Connect to PostgreSQL
docker exec -it wacast-postgres psql -U postgres -d wacast

# Check migrations
SELECT * FROM migrations;

# Check devices
SELECT id, phone, status FROM devices;
```

### 4. Logs

```bash
# All logs are JSON structured from Zap
# Easily parsed by log aggregation systems
# Example: tail -f app.log | jq '.msg'
```

---

## Next Steps

1. **Phase 7:** HTTP Server - Gin framework with full REST API
2. **Phase 8:** Webhooks - Push events to API layer
3. **Phase 9:** Advanced - warming, auto-response, broadcasts

---

## Technology Stack

- **Runtime:** Go 1.21+
- **Database:** PostgreSQL 15+
- **WhatsApp:** go.mau.fi/whatsmeow
- **Logging:** Uber Zap (structured JSON)
- **HTTP:** Gin framework (Phase 7)
- **Encryption:** golang.org/x/crypto (AES-GCM)
- **QR Codes:** github.com/skip2/go-qrcode

---

## Contributing

Each Phase builds on the previous:

- Phase 5 → Session Service ✅
- Phase 6 → Message Service ✅
- Phase 7 → HTTP API
- Phase 8+ → Advanced features

All code follows Go conventions and includes comprehensive error handling.

---

## Support

See documentation files for detailed information:

- **Setup issues:** QUICKSTART.md
- **Database questions:** MIGRATIONS_GUIDE.md
- **Session code:** SESSION_SERVICE_GUIDE.md
- **Message code:** MESSAGE_SERVICE_GUIDE.md
- **Status:** DEVELOPMENT_STATUS.md

---

_Last Updated: Phase 6 - Message Service Complete_

- HTTP server & handlers
- Webhook system
- Auto-response system
- Account warming system

## Configuration

All configuration via environment variables in `.env`:

```bash
# Core
ENVIRONMENT=development
SERVER_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=123456
DB_NAME=wacast

# Connection Pool
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5
DB_CONN_MAX_IDLE_TIME=2

# Logging
LOG_LEVEL=debug

# WhatsApp
WHATSAPP_SESSION_TIMEOUT=300

# Encryption
ENCRYPTION_KEY=your-secret-key-here
```

See `.env.example` for all available options.

## Development

### Logging

Simple logging via `utils` package:

```go
import "wacast/core/utils"

utils.Info("User logged in", zap.String("user_id", userId))
utils.Error("Database connection failed", zap.Error(err))
utils.Debug("Processing message", zap.Any("message", msg))
```

### Database Operations

Use methods from `database.Database`:

```go
// Query single row
err := db.QueryRow("SELECT * FROM users WHERE id = $1", userID).
    Scan(&user.ID, &user.Phone, ...)

// Query multiple rows
rows, err := db.Query("SELECT * FROM messages WHERE device_id = $1", deviceID)

// Execute command
result, err := db.Exec("UPDATE devices SET status = $1 WHERE id = $2", status, deviceID)

// Transaction
tx, err := db.BeginTx()
tx.Exec(...)
tx.Commit()
```

## Docker Support

### Development with Docker

```bash
# Start all services
docker-compose up

# Start only database
docker-compose up -d postgres

# View logs
docker-compose logs -f wacast-core

# Stop services
docker-compose down
```

### Production Build

```bash
# Build image
docker build -t wacast-core:latest .

# Run container
docker run -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_USER=postgres \
  -e DB_PASSWORD=123456 \
  wacast-core:latest
```

## Directory Details

### `/config`

Manages all configuration loading, environment variables, and constants. Database configuration, server settings, and feature flags are defined here.

### `/database`

Data access layer providing connection pool management and query builders. Abstracts all direct SQL operations.

### `/models`

Go structs that mirror database schema. Used throughout the application for type safety and consistency.

### `/services`

Core business logic including WhatsApp session management, message routing, and event notifications. This is where the main gateway logic lives.

### `/handlers`

Request handlers that receive from API layer, call services, and return responses. Currently minimal as requests come directly from services are testing.

### `/utils`

Reusable functions: logging, encryption, validation. Helpers that don't fit in other packages.

### `/migrations`

SQL migration files for schema versioning and database updates. Run on application startup.

## Next Steps

1. Create Models (`/models/*`) - Done ✅
2. Create Database Queries (`/database/*_queries.go`) - Done ✅
3. Setup Migrations Runner - Next ⬜
4. Build WhatsApp Session Service - Core functionality ⬜
5. Build Message Service - Message routing ⬜
6. Create HTTP Server & Handlers - API layer ⬜

## Support

For setup issues, see [QUICKSTART.md](./QUICKSTART.md) troubleshooting section.
