# Services

Business logic & WhatsApp integration.

## Files:

- `session_service.go` - Multi-session management, connect/disconnect
- `message_service.go` - Send/receive messages, status tracking
- `device_service.go` - Device activation & deactivation
- `webhook_service.go` - Outgoing webhook notifications

## Purpose:

Core logic untuk:

- Initialize sessions dari session_data
- Route messages in/out
- Track message status
- Trigger webhooks ke API layer
