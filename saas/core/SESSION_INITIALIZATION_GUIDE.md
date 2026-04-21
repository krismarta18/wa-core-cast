# Session Initialization Guide

**Konsep:** Session = Koneksi WhatsApp Web aktif yang di-link ke nomor tertentu via QR code

---

## QR Code dari WhatsApp? ✅

**YA!** QR code yang di-generate adalah dari **WhatsApp Linked Devices protocol**.

Flow-nya:

1. Sistem membuat **WhatsApp Web session** via whatsmeow library
2. WhatsApp generate **actual QR code** untuk linking
3. User scan QR code dengan **WhatsApp mobile app**
4. Mobile app link ke Web session → **Session ACTIVE**

**Format QR:** `whatsapp://linkdevice?request=<device_id>&code=<unique_code>`

---

## Apa itu Initialize?

`/api/v1/sessions/initiate` adalah endpoint untuk **membuat koneksi WhatsApp baru** pada device tertentu.

**Yang terjadi:**

1. Sistem membuat "session" baru
2. Generate QR code
3. User scan QR code dengan WhatsApp mobile app
4. Session menjadi active setelah di-scan

---

## Flow Lengkap

### Step 1: Initiate Session

```bash
curl -X POST http://localhost:8080/api/v1/sessions/initiate \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "device001",
    "user_id": "user123",
    "phone": "62812345678"
  }'
```

**Response:**

```json
{
  "message": "Session initiated. Scan QR code via GET /api/v1/sessions/{device_id}/qr",
  "device_id": "device001",
  "status": 2,
  "next_step": "GET /api/v1/sessions/device001/qr to retrieve QR code",
  "poll_status": "GET /api/v1/sessions/device001 to check status"
}
```

**Status Codes:**

- `202 Accepted` - Session created successfully, waiting untuk QR scan
- `400 Bad Request` - Parameter tidak valid
- `500 Server Error` - Failed to create session

**Parameters:**
| Parameter | Required | Tipe | Keterangan |
|-----------|----------|------|-----------|
| device_id | ✅ | string | Unique ID untuk device (bisa apa saja, contoh: "device001", "instance-1", dll) |
| user_id | ✅ | string | User yang punya device ini |
| phone | ✅ | string | Nomor WhatsApp dengan country code, **tanpa +** (contoh: "62812345678") |

---

### Step 2: Retrieve QR Code

```bash
curl http://localhost:8080/api/v1/sessions/device001/qr
```

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
- User scan dengan WhatsApp mobile → di-redirect ke link device flow
- After scan: session auto menjadi ACTIVE

**Cara menggunakan:**

1. Browser/aplikasi generate QR code dari string ini
2. User buka WhatsApp mobile app
3. Go to Settings > Linked Devices > Scan QR code
4. Scan hasil dari endpoint ini
5. Device will connect secara otomatis

---

### Step 3: Poll Session Status

```bash
# Check apakah session sudah active
curl http://localhost:8080/api/v1/sessions/device001
```

**Response sebelum scan:**

```json
{
  "device_id": "device001",
  "status": 2,
  "is_active": false
}
```

**Response setelah scan:**

```json
{
  "device_id": "device001",
  "status": 1,
  "is_active": true
}
```

**Status Codes:**
| Status | Value | Arti |
|--------|-------|------|
| Inactive | 0 | Device tidak connect |
| Active | 1 | ✅ Ready buat send/receive messages |
| Pending | 2 | Waiting untuk QR code scan |

**Polling Strategy:**

```javascript
// Poll setiap 2 detik sampai status jadi 1 (active)
const interval = setInterval(async () => {
  const res = await fetch("/api/v1/sessions/device001");
  const data = await res.json();

  if (data.status === 1) {
    console.log("Session active! Ready to use");
    clearInterval(interval);
  }
}, 2000);
```

---

## Session Lifecycle

```
┌─────────────────────────────────────────────────────────┐
│ 1. POST /api/v1/sessions/initiate                      │
│    (Create session, generate QR code)                   │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│ Status: 2 (PENDING)                                    │
│ 2. GET /api/v1/sessions/{device_id}/qr                │
│    (Show QR code to user)                              │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
        ┌─────────────────────────┐
        │ User scans QR with      │
        │ WhatsApp mobile app     │
        └────────────┬────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│ Status: 1 (ACTIVE) ✅                                   │
│ 3. POST /api/v1/devices/{device_id}/messages           │
│    (Now can send messages!)                             │
└─────────────────────────────────────────────────────────┘
```

---

## Use Cases & Examples

### Use Case 1: Multi-Device Manager

Setup 5 devices untuk send bulk messages:

```bash
# Create 5 sessions
for i in {1..5}; do
  curl -X POST http://localhost:8080/api/v1/sessions/initiate \
    -H "Content-Type: application/json" \
    -d '{
      "device_id": "device_'$i'",
      "user_id": "bulk_sender",
      "phone": "6281234567'$i'"
    }'
done

# Get QR codes untuk semua device
# Display di web interface
# User scan QR codes dari 5 WhatsApp accounts

# Check status penuh scan
curl http://localhost:8080/api/v1/sessions | jq '.sessions | map(select(.is_active == true)) | length'

# Setelah semua active, broadcast message ke 5 device sekaligus
```

### Use Case 2: Disaster Recovery

Restore session yang previous:

```bash
# Check session yang tersimpan di database
curl http://localhost:8080/api/v1/sessions

# Jika inactive, re-initiate:
curl -X POST http://localhost:8080/api/v1/sessions/initiate \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "device001",
    "user_id": "user123",
    "phone": "62812345678"
  }'
```

### Use Case 3: Session Timeout

Reconnect ke inactive session:

```bash
# Check current status
curl http://localhost:8080/api/v1/sessions/device001
# Response: {"device_id": "device001", "status": 0, "is_active": false}

# Re-initiate (QR code sama seperti pertama kali)
curl -X POST http://localhost:8080/api/v1/sessions/initiate \
  -H "Content-Type: application/json" \
  -d '{"device_id": "device001", "user_id": "user123", "phone": "62812345678"}'

# Scan QR code lagi
```

---

## Error Handling

### Session Already Exists

```json
{
  "error": "session already exists for device device001"
}
```

**Solusi:**

- Stop session lama dulu: `POST /api/v1/sessions/{device_id}/stop`
- Atau gunakan device_id berbeda

### QR Code Not Available

```json
{
  "error": "QR code not available for this device"
}
```

**Solusi:**

- Device belum di-initialize
- Atau session sudah expired
- Re-initiate session baru

### Session Not Active

Tidak bisa send message ketika status bukan 1:

```json
{
  "error": "session not active: device001"
}
```

**Solusi:**

- Pastikan user sudah scan QR code
- Check status dengan: `GET /api/v1/sessions/{device_id}`
- Re-scan jika perlu

---

## API Endpoints Summary

| Endpoint                                     | Method | Tujuan                     |
| -------------------------------------------- | ------ | -------------------------- |
| `/api/v1/sessions`                           | GET    | List semua active sessions |
| `/api/v1/sessions/{device_id}`               | GET    | Check status session       |
| `/api/v1/sessions/{device_id}/qr`            | GET    | Retrieve QR code           |
| `/api/v1/sessions/initiate`                  | POST   | Create new session         |
| `/api/v1/sessions/{device_id}/stop`          | POST   | Disconnect session         |
| `/api/v1/devices/{device_id}/messages`       | POST   | Send text message          |
| `/api/v1/devices/{device_id}/messages/media` | POST   | Send media                 |

---

## Best Practices

✅ **DO:**

- Reuse session yang sudah active
- Poll status dengan interval 2-5 detik
- Store device_id dan user_id di database
- Handle session timeout dengan auto-reconnect
- Log QR scan events untuk audit

❌ **DON'T:**

- Create multiple session dengan device_id same
- Ignore session status checks
- Send message sebelum status 1 (active)
- Assume QR code universal (generate per session)
- Leave session inactive untuk lama (akan disconnect)

---

## Production Considerations

**Phase 8+ TODO:**

- [ ] WebSocket untuk real-time QR code delivery
- [ ] Webhook notifications untuk status change
- [ ] Automated QR code expiry (10 minutes)
- [ ] Session persistence di database
- [ ] Auto-reconnect dengan credential save
- [ ] Rate limiting per device
- [ ] Session analytics & monitoring

---

**Last Updated:** April 10, 2026  
**API Version:** 1.0.0
