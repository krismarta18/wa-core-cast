# Database

PostgreSQL connection, query builders, dan data access layer.

## Files:

- `db.go` - Database initialization & connection pool
- `user_queries.go` - User-related queries
- `device_queries.go` - Device & session queries
- `message_queries.go` - Message logging & tracking
- `subscription_queries.go` - Subscription & billing queries

## Purpose:

Abstrak semua database operations. Services hanya call functions di sini, tidak langsung SQL query.
