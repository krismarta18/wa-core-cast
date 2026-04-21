# Models & Database Queries Guide

## Overview

Semua 17 tables dari schema SQL sudah di-convert menjadi Go models dan database query functions.

## Models Breakdown

### 1. User Management (`models/user.go`)

```
User struct fields:
- ID, Phone, NamaLengkap, IsVerify, OTPCode, OTPExpired
- IDSubscribed (subscription relation), MaxDevice
- IsBan, IsAPI

Requests:
- CreateUserRequest
- UpdateUserRequest

Database Functions (user_queries.go):
- CreateUser()
- GetUserByID(), GetUserByPhone()
- UpdateUser(), DeleteUser()
- VerifyUser(), BanUser()
- UpdateUserMaxDevice()
- GetAllUsers(), GetUserCount()
```

### 2. Device/Session Management (`models/device.go`)

```
Device struct fields:
- ID, UserID, UniqueName, NameDevice, Phone
- Status (0: inactive, 1: active, 2: disconnect)
- LastSeen, SessionData (encrypted []byte)

Database Functions (device_queries.go):
- CreateDevice()
- GetDeviceByID(), GetDevicesByUserID()
- GetActiveDevices(), GetDeviceByPhone()
- UpdateDeviceStatus(), UpdateDeviceSessionData()
- UpdateDeviceLastSeen()
- CountUserDevices()
- DeleteDevice()

Multi-Session Support:
- Each user can have multiple devices (limited by billing plan)
- Session data stored encrypted in DB
- Status tracking untuk auto-reconnect on startup
```

### 3. Message Handling (`models/message.go`)

```
Message struct fields:
- ID, DeviceID, Direction (IN/OUT)
- ReceiptNumber, MessageType, Content
- StatusMessage (0:pending, 1:sent, 2:delivered, 4:failed)
- ErrorLog, CreatedAt

Database Functions (message_queries.go):
- CreateMessage()
- GetMessageByID(), GetMessageByReceiptNumber()
- GetMessagesByDeviceID(), GetMessagesByStatusAndDevice()
- GetIncomingMessages(), GetOutgoingMessages()
- GetPendingMessages()
- UpdateMessageStatus(), UpdateMessageStatusWithError()
- CountMessagesByStatus()
```

### 4. Subscription & Billing (`models/subscription.go`)

```
BillingPlan struct:
- ID, Name, Price, MaxDevice, MaxMessagesDay
- Features (JSONB for flexible features)

Subscription struct:
- ID, UserID, PlanID, Status (0/1)
- CreatedAt, UpdatedAt

Database Functions (subscription_queries.go):
- CreateSubscription(), GetSubscriptionByID()
- GetSubscriptionByUserID()
- UpdateSubscription(), DeactivateSubscription()
- ActivateSubscription()
- CreateBillingPlan(), GetBillingPlanByID()
- GetAllBillingPlans()
- GetSubscriptionCount()
```

### 5. Broadcast Campaign (`models/broadcast.go`)

```
BroadcastCampaign struct:
- ID, UserID, DeviceID
- NameBroadcast, TotalRecipients, ProcessedCount
- ScheduledAt, Status (0:draft, 1:scheduled, 2:running, 3:completed)

BroadcastMessage struct:
- ID, CampaignID, MessageType
- MessageText, MediaUrl, ButtonData (JSONB)

BroadcastRecipient struct:
- ID, CampaignID, GroupsID, ContactID
- Status (0:pending, 1:sent, 2:failed)
- SentAt, ErrorMessages, RetryCount

Database Functions (broadcast_queries.go):
- Campaign: Create, GetByID, GetByUserID, UpdateStatus, UpdateProgress, Delete
- Message: Create, GetByCampaignID
- Recipient: Create, GetByCampaignID, GetPending, UpdateStatus
- Counting: CountRecipientsByStatus()
```

### 6. Contact & Group (`models/contact.go`)

```
Contact struct:
- ID, GroupID, Name, Phone
- AdditionalData (JSONB)
- CreatedAt, UpdatedAt, DeletedAt (soft delete)

Group struct:
- ID, UserID, GroupName
- CreatedAt, DeletedAt (soft delete)

Database Functions (contact_queries.go):
- Contact: Create, GetByID, GetByPhone, GetByGroupID, Update, Delete
- Group: Create, GetByID, GetByUserID, Update, Delete
- Utility: CountGroupContacts()
```

### 7. Account Warming (`models/warming.go`)

```
WarmingPool struct:
- ID, DeviceID, Intensity, DailyLimit
- MessageSendToday, IsActive
- NextActionAt (scheduler)

WarmingSession struct:
- ID, DeviceID, TargetPhone
- MessageSent, ResponseReceived
- Status (0:pending, 1:sent, 2:response_received, 3:failed)

Database Functions (warming_queries.go):
- Pool: Create, GetByDeviceID, Update, ResetDaily, Activate/Deactivate
- Session: Create, GetByDeviceID, UpdateStatus
- Utility: GetActiveWarmingPools(), UpdateMessageCount()
```

### 8. Other Features (`models/other.go`)

```
AutoResponse struct:
- ID, DeviceID, Keyword, ResponseText, IsActive

Webhook struct:
- ID, DeviceID, WebhookUrl, SecretKey

APILog struct:
- ID, UserID, Endpoint
- ReqBody, ResponseBody (JSONB)
- CreatedAt, IPAddress, DeviceID

SystemSetting struct:
- ID, Keys, Value, Description

Database Functions (other_queries.go - 30+ functions):
- AutoResponse: Create, GetByID, GetByDeviceID, Update, Delete
- Webhook: Create, GetByID, GetByDeviceID, Update, Delete
- APILog: Create, GetByUserID
- SystemSetting: GetSystemSetting(), UpdateSystemSetting()
```

## Database Query Patterns

### Query Example 1: Simple Get

```go
// GetUserByID retrieves a user by ID
func (d *Database) GetUserByID(userID uuid.UUID) (*models.User, error) {
    query := `SELECT id, phone, nama_lengkap, ... FROM users WHERE id = $1`

    user := &models.User{}
    err := d.QueryRow(query, userID).Scan(&user.ID, &user.Phone, ...)
    if err != nil {
        return nil, err
    }
    return user, nil
}
```

### Query Example 2: List with Pagination

```go
// GetMessagesByDeviceID retrieves messages with pagination
func (d *Database) GetMessagesByDeviceID(deviceID uuid.UUID, limit, offset int) {
    query := `SELECT ... FROM messages WHERE device_id = $1
              ORDER BY created_at DESC LIMIT $2 OFFSET $3`

    rows, err := d.Query(query, deviceID, limit, offset)
    // scan rows and return []models.Message
}
```

### Query Example 3: Update with Conditional Fields

```go
// UpdateUser updates only non-nil fields
func (d *Database) UpdateUser(userID uuid.UUID, update *models.UpdateUserRequest) error {
    query := `UPDATE users SET `
    args := []interface{}{}
    argCount := 1

    if update.NamaLengkap != nil {
        query += fmt.Sprintf("nama_lengkap = $%d", argCount)
        args = append(args, *update.NamaLengkap)
        argCount++
    }
    // ... repeat for other fields

    query += fmt.Sprintf(" WHERE id = $%d", argCount)
    args = append(args, userID)

    _, err := d.Exec(query, args...)
    return err
}
```

## Usage Examples

### Creating New Records

```go
import "wacast/core/database"
import "wacast/core/models"

// Create user
user := &models.User{
    ID:          uuid.New(),
    Phone:       "62812345678",
    NamaLengkap: "John Doe",
    IsVerify:    false,
}

err := database.DB.CreateUser(user)
if err != nil {
    log.Error("Failed to create user", err)
}
```

### Reading Records

```go
// Get user by ID
user, err := database.DB.GetUserByID(userID)
if err != nil {
    log.Error("User not found")
}

// Get devices for user
devices, err := database.DB.GetDevicesByUserID(userID)
for _, device := range devices {
    if device.IsActive() {
        // Handle active device
    }
}
```

### Updating Records

```go
// Update device status
err := database.DB.UpdateDeviceStatus(deviceID, 1) // 1 = active

// Update with partial fields
update := &models.UpdateUserRequest{
    NamaLengkap: ptr("New Name"),
    IsVerify:    ptr(true),
}
err := database.DB.UpdateUser(userID, update)
```

### Querying with Status Filters

```go
// Get pending messages
messages, err := database.DB.GetMessagesByStatusAndDevice(
    deviceID,
    0, // 0 = pending status
    limit,
    offset,
)

// Get active warming pools
pools, err := database.DB.GetActiveWarmingPools()
```

## Type Safety & Validation

### Models Include:

- Struct tags untuk JSON marshaling
- Helper methods (e.g., `IsActive()`, `GetStatusText()`)
- Response transformation methods (e.g., `ToResponse()`)

### Database Layer Ensures:

- Type-safe scanning dengan automatic type conversion
- Error handling untuk setiap operation
- Structured logging untuk debugging

## Integration with Services

Database queries dirancang untuk be called oleh service layer:

```
HTTP Request
    ↓
Handlers (handlers/*)
    ↓
Services (services/*) ← Uses Database Queries
    ↓
Database Layer (database/*_queries.go)
    ↓
PostgreSQL
```

Contoh flow:

```go
// In message_service.go
func (s *MessageService) SendMessage(req *models.SendMessageRequest) error {
    // 1. Get device
    device, err := database.DB.GetDeviceByID(req.DeviceID)

    // 2. Create message
    message := &models.Message{
        ID: uuid.New(),
        DeviceID: req.DeviceID,
        // ... fill fields
    }
    err = database.DB.CreateMessage(message)

    // 3. Send via WhatsApp library
    // ... whatsmeow integration

    // 4. Update status
    err = database.DB.UpdateMessageStatus(message.ID, 1) // sent

    return nil
}
```

## Performance Considerations

### Connection Pooling

- Max open connections: 25
- Max idle connections: 5
- Configurable via `.env`

### Queries Optimized For:

- Common access patterns (GetByID, GetByUserID, GetByDeviceID)
- Pagination support untuk large result sets
- Status-based filtering
- Soft deletes (deleted_at field)

### Future Optimizations:

- Add database indexes untuk frequently queried columns
- Implement caching layer for read-heavy operations
- Batch operations untuk bulk inserts/updates

## Next Phase: Migrations

Database migrations runner akan:

1. Run initial schema berdasarkan wecast.sql
2. Add indexes untuk performance
3. Support future schema updates

See next phase documentation untuk migration setup.
