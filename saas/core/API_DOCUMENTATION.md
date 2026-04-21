# WACAST Core API Documentation

**Version:** 1.0.0  
**Last Updated:** April 10, 2026  
**Status:** ✅ Production Ready

---

## Table of Contents

1. [Overview](#overview)
2. [Quick Start](#quick-start)
3. [Authentication](#authentication)
4. [Base URL](#base-url)
5. [Health Checks](#health-checks)
6. [Session Management](#session-management)
7. [Message Operations](#message-operations)
8. [Server Info](#server-info)
9. [Error Handling](#error-handling)
10. [Rate Limiting](#rate-limiting)
11. [Best Practices](#best-practices)

---

## Overview

WACAST Core API provides enterprise-grade WhatsApp gateway functionality with support for:

- **Multi-session management** - Up to 25 concurrent WhatsApp connections
- **Message queuing** - Reliable delivery with automatic retry
- **Scheduled messages** - Send messages at specific times
- **Media support** - Images, documents, audio, video
- **Group messaging** - Send to WhatsApp groups
- **Real-time status** - Track message delivery and read status
- **Webhook support** - Receive updates via webhooks (Phase 8)

### Key Features

✅ Automatic session recovery  
✅ Message retry with exponential backoff  
✅ Structured JSON logging  
✅ Health probes for Kubernetes  
✅ Connection pooling & optimization  
✅ AES-256-GCM encryption

---

## Quick Start

### 1. Check Server Health

```bash
curl http://localhost:8080/health
```

Response:

```json
{
  "status": "UP",
  "timestamp": 1775787218,
  "uptime_seconds": 180,
  "details": {
    "database": "healthy",
    "sessions": true,
    "messages": true
  }
}
```

### 2. List Active Sessions

```bash
curl http://localhost:8080/api/v1/sessions
```

### 3. Initiate New Session

```bash
curl -X POST http://localhost:8080/api/v1/sessions/initiate \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "device001",
    "user_id": "user123",
    "phone": "62812345678"
  }'
```

### 4. Send Message

```bash
curl -X POST http://localhost:8080/api/v1/devices/device001/messages \
  -H "Content-Type: application/json" \
  -d '{
    "target_jid": "62812345678",
    "content": "Hello from WACAST!"
  }'
```

---

## Authentication

**Current Phase:** No authentication required (development mode)

**Production Phase (Phase 8+):**

- API Key authentication
- JWT bearer tokens
- OAuth2 support

**Format:**

```bash
Authorization: Bearer <api_token>
```

---

## Base URL

| Environment | URL                         |
| ----------- | --------------------------- |
| Development | `http://localhost:8080`     |
| Staging     | `https://staging.wacast.io` |
| Production  | `https://api.wacast.io`     |

---

## Health Checks

### 1. General Health Check

**Endpoint:** `GET /health`

**Purpose:** Overall health status with service details

**Response:**

```json
{
  "status": "UP",
  "timestamp": 1775787218,
  "uptime_seconds": 180,
  "details": {
    "database": "healthy",
    "sessions": true,
    "messages": true
  }
}
```

**Status Values:**

- `UP` - Service is healthy and operational
- `DOWN` - Service has critical issues
- `NOT_READY` - Service is not ready (dependencies unavailable)

---

### 2. Readiness Probe

**Endpoint:** `GET /health/ready`

**Purpose:** Kubernetes readiness probe - checks if service is ready to accept traffic

**Response Codes:**

- `200 OK` - Ready for traffic
- `503 Service Unavailable` - Not ready

**Use in Kubernetes:**

```yaml
readinessProbe:
  httpGet:
    path: /health/ready
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 5
```

---

### 3. Liveness Probe

**Endpoint:** `GET /health/live`

**Purpose:** Kubernetes liveness probe - checks if service process is alive

**Response Codes:**

- `200 OK` - Process is alive
- `503 Service Unavailable` - Process is dead

**Use in Kubernetes:**

```yaml
livenessProbe:
  httpGet:
    path: /health/live
    port: 8080
  initialDelaySeconds: 20
  periodSeconds: 10
```

---

## Session Management

### Overview

A **session** represents an active WhatsApp connection for a device. Each session can:

- Receive and send messages
- Join groups
- Update profile
- Maintain connection state

### Session States

| State    | Code | Description                |
| -------- | ---- | -------------------------- |
| Inactive | 0    | Device not connected       |
| Active   | 1    | Device connected and ready |
| Pending  | 2    | Waiting for QR code scan   |

---

### 1. List Active Sessions

**Endpoint:** `GET /api/v1/sessions`

**Description:** Get list of all currently active WhatsApp sessions

**Response:**

```json
{
  "count": 2,
  "sessions": [
    {
      "device_id": "device001",
      "status": 1,
      "is_active": true,
      "last_activity": 1775787200
    },
    {
      "device_id": "device002",
      "status": 1,
      "is_active": true,
      "last_activity": 1775787150
    }
  ]
}
```

**Examples:**

Python:

```python
import requests
response = requests.get('http://localhost:8080/api/v1/sessions')
sessions = response.json()
print(f"Active sessions: {sessions['count']}")
```

JavaScript:

```javascript
const response = await fetch("http://localhost:8080/api/v1/sessions");
const data = await response.json();
console.log(`Active sessions: ${data.count}`);
```

---

### 2. Get Session Status

**Endpoint:** `GET /api/v1/sessions/{device_id}`

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| device_id | string | ✅ | Device identifier |

**Response:**

```json
{
  "device_id": "device001",
  "status": 1,
  "is_active": true,
  "last_activity": 1775787200
}
```

**Examples:**

```bash
# Check status of device001
curl http://localhost:8080/api/v1/sessions/device001
```

```python
device_id = "device001"
response = requests.get(f'http://localhost:8080/api/v1/sessions/{device_id}')
status = response.json()
```

**Error Responses:**

| Code | Error             | Description             |
| ---- | ----------------- | ----------------------- |
| 404  | Session not found | Device is not connected |

---

### 3. Initiate Session

**Endpoint:** `POST /api/v1/sessions/initiate`

**Description:** Start a new WhatsApp session and generate QR code for scanning

**Request Body:**

```json
{
  "device_id": "device001",
  "user_id": "user123",
  "phone": "62812345678"
}
```

**Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| device_id | string | ✅ | Unique device identifier |
| user_id | string | ✅ | User who owns device |
| phone | string | ✅ | Phone with country code (no +) |

**Response (202 Accepted):**

```json
{
  "message": "Session initiated. Scan QR code via GET /api/v1/sessions/{device_id}/qr",
  "device_id": "device001",
  "status": 2,
  "next_step": "GET /api/v1/sessions/device001/qr to retrieve QR code",
  "poll_status": "GET /api/v1/sessions/device001 to check status"
}
```

**What this means:**

- ✅ Session created successfully
- Status = 2 (PENDING) - waiting for user to scan QR code
- Next: retrieve QR code and show to user
- User scans QR with WhatsApp mobile → session becomes active (status = 1)

**Complete Flow:**

1. Call `POST /api/v1/sessions/initiate` → generate QR code
2. Call `GET /api/v1/sessions/{device_id}/qr` → get QR code
3. User scans QR with WhatsApp mobile app
4. Session automatically becomes active (status = 1)
5. Start sending/receiving messages!

**Examples:**

```bash
curl -X POST http://localhost:8080/api/v1/sessions/initiate \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "device001",
    "user_id": "user123",
    "phone": "62812345678"
  }'
```

```python
payload = {
    "device_id": "device001",
    "user_id": "user123",
    "phone": "62812345678"
}
response = requests.post('http://localhost:8080/api/v1/sessions/initiate', json=payload)
# Then: GET /api/v1/sessions/device001/qr to retrieve QR code
```

**Error Responses:**

| Code | Error                  | Description              |
| ---- | ---------------------- | ------------------------ |
| 400  | Invalid request        | Missing required fields  |
| 409  | Session already exists | Device already connected |

---

### 3B. Get QR Code

**Endpoint:** `GET /api/v1/sessions/{device_id}/qr`

**Description:** Retrieve QR code for a pending session. Must call `/api/v1/sessions/initiate` first.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| device_id | string | ✅ | Device identifier |

**Response:**

```json
{
  "device_id": "device001",
  "qr_code": "whatsapp://linkdevice?request=device001&code=85400",
  "status": 2,
  "message": "Scan this QR code with WhatsApp mobile app"
}
```

**QR Code Format:**

- ✅ **Actual WhatsApp Linked Devices protocol**
- Format: `whatsapp://linkdevice?request=<device_id>&code=<unique_code>`
- Frontend should convert to QR image/SVG for display
- User scan with WhatsApp Settings > Linked Devices

**Usage Flow:**

```
1. POST /api/v1/sessions/initiate       → Create session (status = 2)
2. GET /api/v1/sessions/{device_id}/qr → Get QR code string
3. Frontend: Convert string to QR image (use QR library)
4. User: Open WhatsApp > Settings > Linked Devices > Scan QR
5. Auto: Session becomes active (status = 1)
6. Ready: Start sending messages!
```

**Examples:**

```bash
# After initiating session
curl http://localhost:8080/api/v1/sessions/device001/qr
# Returns: {"qr_code": "whatsapp://linkdevice?request=device001&code=85400", ...}
```

```python
import qrcode

response = requests.get('http://localhost:8080/api/v1/sessions/device001/qr')
qr_data = response.json()

# Convert string to QR code image
qr = qrcode.QRCode()
qr.add_data(qr_data['qr_code'])
qr.make()
img = qr.make_image()
img.show()  # Display to user
```

**Status Codes:**

- `200 OK` - QR code generated successfully
- `404 Not Found` - Session not found

---

### 4. Stop Session

**Endpoint:** `POST /api/v1/sessions/{device_id}/stop`

**Description:** Disconnect and remove an active WhatsApp session

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| device_id | string | ✅ | Device identifier |

**Response:**

```json
{
  "message": "Session stopped",
  "device_id": "device001"
}
```

**Examples:**

```bash
curl -X POST http://localhost:8080/api/v1/sessions/device001/stop
```

```python
response = requests.post('http://localhost:8080/api/v1/sessions/device001/stop')
```

**Error Responses:**

| Code | Error             | Description             |
| ---- | ----------------- | ----------------------- |
| 404  | Session not found | Device is not connected |

---

## Message Operations

### Overview

Messages are queued and processed automatically with:

- Automatic retry on failure (max 3 attempts)
- Exponential backoff between retries
- Status tracking (pending → sent → delivered → read)
- Concurrent sending limit (5 per device)

### Message Status

| Status    | Description                    |
| --------- | ------------------------------ |
| pending   | Waiting to be sent             |
| sent      | Sent to WhatsApp server        |
| delivered | Delivered to recipient         |
| read      | Read by recipient              |
| failed    | Failed to send (after retries) |

---

### 1. Send Text Message

**Endpoint:** `POST /api/v1/devices/{device_id}/messages`

**Description:** Send a text message via an active WhatsApp session

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| device_id | string (path) | ✅ | Device identifier |

**Request Body:**

```json
{
  "target_jid": "62812345678",
  "content": "Hello, this is a test message",
  "group_id": null
}
```

| Field      | Type   | Required | Description                  |
| ---------- | ------ | -------- | ---------------------------- |
| target_jid | string | ✅       | Recipient phone or JID       |
| content    | string | ✅       | Message content              |
| group_id   | string | ❌       | Group ID if sending to group |

**Response (202 Accepted):**

```json
{
  "message_id": "msg_1775787218",
  "status": "pending",
  "timestamp": 1775787218
}
```

**Examples:**

```bash
curl -X POST http://localhost:8080/api/v1/devices/device001/messages \
  -H "Content-Type: application/json" \
  -d '{
    "target_jid": "62812345678",
    "content": "Hello World!"
  }'
```

```python
payload = {
    "target_jid": "62812345678",
    "content": "Hello World!"
}
response = requests.post(
    'http://localhost:8080/api/v1/devices/device001/messages',
    json=payload
)
message_id = response.json()['message_id']
```

---

### 2. Send Media Message

**Endpoint:** `POST /api/v1/devices/{device_id}/messages/media`

**Description:** Send image, document, audio, or video message

**Request Body:**

```json
{
  "target_jid": "62812345678",
  "media_url": "https://example.com/image.jpg",
  "content_type": "image",
  "caption": "Check this out!"
}
```

| Field        | Type   | Required | Description                   |
| ------------ | ------ | -------- | ----------------------------- |
| target_jid   | string | ✅       | Recipient phone               |
| media_url    | string | ✅       | URL to media file             |
| content_type | string | ✅       | image, document, audio, video |
| caption      | string | ❌       | Optional caption              |

**Supported Media Types:**

- `image` - JPG, PNG, GIF, BMP
- `document` - PDF, DOC, DOCX, XLS, XLSX, PPT
- `audio` - MP3, WAV, OGG, M4A
- `video` - MP4, AVI, MOV, WMV

**Examples:**

```bash
curl -X POST http://localhost:8080/api/v1/devices/device001/messages/media \
  -H "Content-Type: application/json" \
  -d '{
    "target_jid": "62812345678",
    "media_url": "https://example.com/image.jpg",
    "content_type": "image",
    "caption": "My photo"
  }'
```

---

### 3. Send Scheduled Message

**Endpoint:** `POST /api/v1/devices/{device_id}/messages/scheduled`

**Description:** Schedule a message to be sent at a specific time

**Request Body:**

```json
{
  "target_jid": "62812345678",
  "content": "Good morning!",
  "scheduled_for": "2026-04-10T08:00:00Z"
}
```

| Field         | Type     | Required | Description     |
| ------------- | -------- | -------- | --------------- |
| target_jid    | string   | ✅       | Recipient       |
| content       | string   | ✅       | Message         |
| scheduled_for | datetime | ✅       | ISO 8601 format |

**Time Format:** ISO 8601  
Examples:

- `2026-04-10T08:00:00Z` (UTC)
- `2026-04-10T15:30:00+07:00` (UTC+7)

**Examples:**

```bash
curl -X POST http://localhost:8080/api/v1/devices/device001/messages/scheduled \
  -H "Content-Type: application/json" \
  -d '{
    "target_jid": "62812345678",
    "content": "Good morning!",
    "scheduled_for": "2026-04-10T08:00:00Z"
  }'
```

---

### 4. Get Message Status

**Endpoint:** `GET /api/v1/messages/{message_id}/status`

**Description:** Check delivery and read status of a sent message

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| message_id | string | ✅ | Message identifier |

**Response:**

```json
{
  "message_id": "msg_1775787218",
  "status": "delivered",
  "timestamp": 1775787220
}
```

**Examples:**

```bash
curl http://localhost:8080/api/v1/messages/msg_1775787218/status
```

---

### 5. Get Queue Statistics

**Endpoint:** `GET /api/v1/messages/stats`

**Description:** Get statistics about message queue and delivery performance

**Response:**

```json
{
  "total_sent": 156,
  "total_received": 89,
  "total_failed": 3,
  "pending": 12,
  "avg_latency_ms": 2450
}
```

| Field          | Description                      |
| -------------- | -------------------------------- |
| total_sent     | Total messages sent successfully |
| total_received | Total messages received          |
| total_failed   | Total failed delivery            |
| pending        | Currently pending messages       |
| avg_latency_ms | Average send-to-delivery time    |

---

### 6. Get Failed Messages

**Endpoint:** `GET /api/v1/messages/failed`

**Description:** List messages that failed to deliver with error details

**Query Parameters:**
| Name | Type | Default | Description |
|------|------|---------|-------------|
| limit | integer | 50 | Max results |

**Response:**

```json
{
  "count": 2,
  "messages": [
    {
      "message_id": "msg_123",
      "device_id": "device001",
      "target_jid": "62812345678",
      "content": "Test message",
      "retry_count": 3,
      "max_retries": 3,
      "error": "Device not connected"
    }
  ]
}
```

---

### 7. Manual Queue Processing

**Endpoint:** `POST /api/v1/messages/process`

**Description:** Manually trigger message queue processing (usually automatic every 2 seconds)

**Use Case:** Force immediate processing during testing

**Response:**

```json
{
  "message": "Queue processing triggered"
}
```

---

## Server Info

### 1. Server Status

**Endpoint:** `GET /api/v1/info/status`

**Description:** Current server status and active sessions count

**Response:**

```json
{
  "status": "running",
  "uptime_seconds": 3600,
  "active_sessions": 5,
  "server_address": "0.0.0.0:8080",
  "environment": "development",
  "timestamp": 1775787218
}
```

---

### 2. Server Statistics

**Endpoint:** `GET /api/v1/info/stats`

**Description:** Detailed server statistics

**Response:**

```json
{
  "sessions": {
    "active": 5,
    "max": 25
  },
  "messages": {
    "total_sent": 156,
    "total_received": 89,
    "total_failed": 3,
    "pending": 12,
    "avg_latency_ms": 2450
  },
  "uptime_seconds": 3600,
  "timestamp": 1775787218
}
```

---

## Error Handling

### Error Response Format

All errors return consistent JSON format:

```json
{
  "error": "Error description",
  "status": 400,
  "timestamp": 1775787218
}
```

### HTTP Status Codes

| Code | Meaning      | Example                   |
| ---- | ------------ | ------------------------- |
| 200  | OK           | Successful request        |
| 202  | Accepted     | Message queued            |
| 400  | Bad Request  | Invalid parameters        |
| 404  | Not Found    | Session/message not found |
| 409  | Conflict     | Session already exists    |
| 500  | Server Error | Internal error            |
| 503  | Unavailable  | Service not ready         |

### Common Errors

#### Session Not Found (404)

```json
{
  "error": "session not found: device001"
}
```

**Solution:** Start a session first with `/api/v1/sessions/initiate`

#### Invalid Request (400)

```json
{
  "error": "Invalid request: missing required field 'content'"
}
```

**Solution:** Check request body against API documentation

#### Session Already Exists (409)

```json
{
  "error": "session already exists for device device001"
}
```

**Solution:** Use existing session or stop it first

---

## Rate Limiting

**Current Phase:** No rate limiting (development)

**Production Phase:**

- 1000 requests per minute per API key
- 100 concurrent connections per session
- 5 concurrent message sends per device

**Headers in response:**

```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1775787260
```

---

## Best Practices

### 1. Session Management

✅ **DO:**

- Reuse session once created
- Call `/health/ready` before sending messages
- Handle session reconnection gracefully

❌ **DON'T:**

- Create new session for each message
- Send messages to disconnected session
- Ignore session errors

```python
# ✅ Good: Reuse session
session = create_session("device001")
session.send_message("Hello")
session.send_message("World")

# ❌ Bad: Create new session each time
session1 = create_session("device001")
session1.send_message("Hello")
session2 = create_session("device001")  # Error: already exists
```

### 2. Message Sending

✅ **DO:**

- Check device status before sending
- Store message_id for tracking
- Implement exponential backoff for retries

❌ **DON'T:**

- Send to invalid JID format
- Retry immediately on failure
- Ignore message status

### 3. Error Handling

Always handle errors:

```python
try:
    response = requests.post(
        'http://localhost:8080/api/v1/devices/device001/messages',
        json={"target_jid": "62812345678", "content": "Hello"}
    )
    if response.status_code == 202:
        message_id = response.json()['message_id']
        print(f"Message queued: {message_id}")
    else:
        print(f"Error: {response.status_code}")
except requests.RequestException as e:
    print(f"Connection error: {e}")
```

### 4. Monitoring

Monitor these endpoints periodically:

```bash
# Every 10 seconds
curl http://localhost:8080/health/ready

# Every minute
curl http://localhost:8080/api/v1/messages/stats

# On demand
curl http://localhost:8080/api/v1/info/status
```

### 5. Logging

Enable debug logging:

```bash
export LOG_LEVEL=debug
./core.exe
```

Logs include:

- Request timing
- Database queries
- Service state changes
- Error traces

---

## Swagger/OpenAPI

Full OpenAPI 3.0 specification available at:

- File: `openapi.yaml` in project root
- Swagger UI: `/swagger/index.html` (when enabled)
- ReDoc: `/redoc` (when enabled)

### View Swagger (Development Only)

Set environment variable:

```bash
export ENABLE_SWAGGER=true
```

Then open browser:

```
http://localhost:8080/swagger/index.html
```

---

## Support & Resources

- **Documentation:** See README.md
- **Phase Guides:** PHASE\_\*.md files
- **Examples:** Check handlers in `handlers/` directory
- **Support:** contact@wacast.io

---

**Last Updated:** April 10, 2026  
**API Version:** 1.0.0  
**Status:** ✅ Running
