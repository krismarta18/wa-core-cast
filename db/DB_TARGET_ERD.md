# WACAST Target ERD

ERD ini diturunkan dari review di `DB_DESIGN_REVIEW_AND_TARGET_SCHEMA.md` dan disusun agar tetap terbaca. Diagram dibagi per domain supaya tidak terlalu padat.

## 1. Identity, Billing, and Devices

```mermaid
erDiagram
    USERS {
        uuid id PK
        varchar phone_number UK
        varchar full_name
        varchar email UK
        varchar company_name
        varchar timezone
        boolean is_verified
        boolean is_banned
        boolean is_api_enabled
        timestamptz created_at
        timestamptz updated_at
        timestamptz last_login_at
    }

    OTP_VERIFICATIONS {
        uuid id PK
        uuid user_id FK
        varchar phone_number
        varchar context
        varchar otp_code_hash
        int attempt_count
        timestamptz expires_at
        timestamptz verified_at
        timestamptz created_at
    }

    USER_SESSIONS {
        uuid id PK
        uuid user_id FK
        varchar session_token_hash
        inet ip_address
        text user_agent
        timestamptz expires_at
        timestamptz revoked_at
        timestamptz created_at
    }

    BILLING_PLANS {
        uuid id PK
        varchar name UK
        numeric price
        int max_devices
        int max_messages_per_day
        jsonb features
        boolean is_active
        timestamptz created_at
        timestamptz updated_at
    }

    SUBSCRIPTIONS {
        uuid id PK
        uuid user_id FK
        uuid plan_id FK
        varchar status
        timestamptz start_date
        timestamptz end_date
        timestamptz renewal_date
        boolean auto_renew
        timestamptz created_at
        timestamptz updated_at
    }

    USAGE_QUOTAS {
        uuid id PK
        uuid user_id FK
        uuid subscription_id FK
        varchar period_key
        int messages_used
        int messages_limit
        int devices_used
        int devices_limit
        timestamptz created_at
        timestamptz updated_at
    }

    INVOICES {
        uuid id PK
        uuid user_id FK
        uuid subscription_id FK
        varchar invoice_number UK
        date issue_date
        date due_date
        timestamptz paid_at
        numeric amount
        varchar currency
        varchar status
        varchar payment_method
        timestamptz created_at
    }

    INVOICE_ITEMS {
        uuid id PK
        uuid invoice_id FK
        varchar description
        int qty
        numeric unit_price
        numeric total_price
    }

    DEVICES {
        uuid id PK
        uuid user_id FK
        varchar unique_name
        varchar display_name
        varchar phone_number
        varchar status
        varchar platform
        varchar wa_version
        int battery_level
        timestamptz last_seen_at
        timestamptz connected_since
        timestamptz created_at
        timestamptz updated_at
    }

    DEVICE_SESSIONS {
        uuid id PK
        uuid device_id FK
        bytea session_blob
        varchar session_status
        timestamptz started_at
        timestamptz ended_at
        timestamptz last_restart_at
        timestamptz created_at
    }

    DEVICE_QR_CODES {
        uuid id PK
        uuid user_id FK
        uuid device_id FK
        text qr_string
        text qr_image_url
        varchar status
        timestamptz generated_at
        timestamptz expired_at
    }

    DEVICE_METRICS {
        uuid id PK
        uuid device_id FK
        bigint uptime_seconds
        int messages_sent_count
        int messages_received_count
        numeric success_rate
        timestamptz recorded_at
    }

    USERS ||--o{ OTP_VERIFICATIONS : requests
    USERS ||--o{ USER_SESSIONS : has
    USERS ||--o{ SUBSCRIPTIONS : owns
    BILLING_PLANS ||--o{ SUBSCRIPTIONS : assigned_to
    USERS ||--o{ USAGE_QUOTAS : tracks
    SUBSCRIPTIONS ||--o{ USAGE_QUOTAS : defines
    USERS ||--o{ INVOICES : billed
    SUBSCRIPTIONS ||--o{ INVOICES : generates
    INVOICES ||--o{ INVOICE_ITEMS : contains
    USERS ||--o{ DEVICES : owns
    DEVICES ||--o{ DEVICE_SESSIONS : has
    USERS ||--o{ DEVICE_QR_CODES : generates
    DEVICES ||--o{ DEVICE_QR_CODES : pairs
    DEVICES ||--o{ DEVICE_METRICS : records
```

## 2. Contacts and Messaging

```mermaid
erDiagram
    USERS {
        uuid id PK
    }

    CONTACT_LABELS {
        uuid id PK
        uuid user_id FK
        varchar name
        varchar color
        timestamptz created_at
    }

    CONTACTS {
        uuid id PK
        uuid user_id FK
        uuid label_id FK
        varchar name
        varchar phone_number
        jsonb additional_data
        text note
        timestamptz created_at
        timestamptz updated_at
        timestamptz deleted_at
    }

    CONTACT_GROUPS {
        uuid id PK
        uuid user_id FK
        varchar name
        text description
        timestamptz created_at
        timestamptz updated_at
        timestamptz deleted_at
    }

    CONTACT_GROUP_MEMBERS {
        uuid id PK
        uuid group_id FK
        uuid contact_id FK
        timestamptz created_at
    }

    BLACKLISTS {
        uuid id PK
        uuid user_id FK
        varchar phone_number
        varchar reason
        timestamptz blocked_at
        timestamptz unblocked_at
        timestamptz created_at
    }

    MESSAGE_TEMPLATES {
        uuid id PK
        uuid user_id FK
        varchar name
        varchar category
        text content
        int used_count
        timestamptz created_at
        timestamptz updated_at
    }

    SCHEDULED_MESSAGES {
        uuid id PK
        uuid user_id FK
        uuid device_id FK
        uuid template_id FK
        uuid group_id FK
        varchar recipient_mode
        jsonb recipient_payload
        text message_content
        timestamptz scheduled_at
        varchar status
        timestamptz created_at
        timestamptz updated_at
        timestamptz executed_at
    }

    SCHEDULED_MESSAGE_RECIPIENTS {
        uuid id PK
        uuid scheduled_message_id FK
        uuid contact_id FK
        uuid group_id FK
        varchar phone_number
        varchar status
        timestamptz created_at
    }

    BROADCAST_CAMPAIGNS {
        uuid id PK
        uuid user_id FK
        uuid device_id FK
        uuid template_id FK
        varchar name
        text message_content
        int delay_seconds
        int total_recipients
        int processed_count
        int success_count
        int failed_count
        timestamptz scheduled_at
        varchar status
        timestamptz created_at
        timestamptz updated_at
        timestamptz deleted_at
    }

    BROADCAST_RECIPIENTS {
        uuid id PK
        uuid campaign_id FK
        uuid contact_id FK
        uuid group_id FK
        varchar phone_number
        varchar status
        timestamptz sent_at
        text error_message
        int retry_count
    }

    MESSAGES {
        uuid id PK
        uuid user_id FK
        uuid device_id FK
        uuid template_id FK
        uuid broadcast_id FK
        uuid scheduled_message_id FK
        varchar direction
        varchar recipient_phone
        varchar sender_phone
        varchar message_type
        text content
        varchar status
        text error_log
        timestamptz created_at
        timestamptz sent_at
        timestamptz delivered_at
        timestamptz failed_at
        timestamptz updated_at
    }

    AUTO_RESPONSE_KEYWORDS {
        uuid id PK
        uuid user_id FK
        uuid device_id FK
        varchar keyword
        text response_text
        boolean is_active
        timestamptz created_at
        timestamptz updated_at
    }

    USERS ||--o{ CONTACT_LABELS : owns
    CONTACT_LABELS ||--o{ CONTACTS : classifies
    USERS ||--o{ CONTACTS : owns
    USERS ||--o{ CONTACT_GROUPS : owns
    CONTACT_GROUPS ||--o{ CONTACT_GROUP_MEMBERS : contains
    CONTACTS ||--o{ CONTACT_GROUP_MEMBERS : joins
    USERS ||--o{ BLACKLISTS : blocks
    USERS ||--o{ MESSAGE_TEMPLATES : creates
    MESSAGE_TEMPLATES ||--o{ SCHEDULED_MESSAGES : reused_by
    CONTACT_GROUPS ||--o{ SCHEDULED_MESSAGES : targets
    SCHEDULED_MESSAGES ||--o{ SCHEDULED_MESSAGE_RECIPIENTS : expands_to
    CONTACTS ||--o{ SCHEDULED_MESSAGE_RECIPIENTS : receives
    CONTACT_GROUPS ||--o{ SCHEDULED_MESSAGE_RECIPIENTS : source_group
    MESSAGE_TEMPLATES ||--o{ BROADCAST_CAMPAIGNS : reused_by
    BROADCAST_CAMPAIGNS ||--o{ BROADCAST_RECIPIENTS : delivers_to
    CONTACTS ||--o{ BROADCAST_RECIPIENTS : targets
    CONTACT_GROUPS ||--o{ BROADCAST_RECIPIENTS : source_group
    USERS ||--o{ MESSAGES : owns
    MESSAGE_TEMPLATES ||--o{ MESSAGES : formats
    BROADCAST_CAMPAIGNS ||--o{ MESSAGES : produces
    SCHEDULED_MESSAGES ||--o{ MESSAGES : executes
    USERS ||--o{ AUTO_RESPONSE_KEYWORDS : defines
```

## 3. Integrations, Monitoring, and App Support

```mermaid
erDiagram
    USERS {
        uuid id PK
    }

    DEVICES {
        uuid id PK
        uuid user_id FK
    }

    WEBHOOKS {
        uuid id PK
        uuid user_id FK
        uuid device_id FK
        varchar webhook_url
        varchar secret_key_hash
        boolean is_active
        timestamptz created_at
        timestamptz updated_at
    }

    WEBHOOK_EVENT_SUBSCRIPTIONS {
        uuid id PK
        uuid webhook_id FK
        varchar event_key
        boolean is_enabled
    }

    WEBHOOK_DELIVERIES {
        uuid id PK
        uuid webhook_id FK
        varchar event_key
        jsonb payload
        int attempt
        int http_status
        text response_body
        varchar status
        timestamptz sent_at
        timestamptz created_at
    }

    API_KEYS {
        uuid id PK
        uuid user_id FK
        varchar name
        varchar key_prefix
        varchar key_hash
        timestamptz last_used_at
        timestamptz expires_at
        boolean is_active
        timestamptz created_at
        timestamptz revoked_at
    }

    API_LOGS {
        uuid id PK
        uuid user_id FK
        uuid device_id FK
        varchar endpoint
        jsonb req_body
        jsonb response_body
        inet ip_address
        timestamptz created_at
    }

    DAILY_MESSAGE_STATS {
        uuid id PK
        uuid user_id FK
        uuid device_id FK
        date stat_date
        int sent_count
        int failed_count
        int delivered_count
        numeric success_rate
        timestamptz created_at
    }

    FAILURE_RECORDS {
        uuid id PK
        uuid user_id FK
        uuid device_id FK
        uuid message_id FK
        varchar recipient_phone
        varchar failure_type
        text failure_reason
        timestamptz occurred_at
        timestamptz created_at
    }

    SERVICE_HEALTH_CHECKS {
        uuid id PK
        varchar service_name
        varchar status
        int latency_ms
        numeric uptime_percent
        timestamptz last_incident_at
        timestamptz checked_at
    }

    RESOURCE_USAGE_METRICS {
        uuid id PK
        varchar metric_name
        numeric metric_value
        varchar metric_unit
        timestamptz recorded_at
    }

    ONBOARDING_PROGRESS {
        uuid id PK
        uuid user_id FK
        varchar step_key
        boolean is_completed
        timestamptz completed_at
        timestamptz created_at
        timestamptz updated_at
    }

    NOTIFICATION_SETTINGS {
        uuid id PK
        uuid user_id FK
        varchar event_key
        boolean email_enabled
        boolean in_app_enabled
        timestamptz created_at
        timestamptz updated_at
    }

    NOTIFICATIONS {
        uuid id PK
        uuid user_id FK
        varchar type
        varchar title
        text body
        boolean is_read
        timestamptz read_at
        timestamptz created_at
    }

    AUDIT_LOGS {
        uuid id PK
        uuid user_id FK
        varchar action_type
        varchar resource_type
        uuid resource_id
        jsonb metadata
        inet ip_address
        text user_agent
        timestamptz created_at
    }

    WARMING_POOL {
        uuid id PK
        uuid device_id FK
        int intensity
        int daily_limit
        int message_send_today
        boolean is_active
        timestamptz next_action_at
    }

    WARMING_SESSIONS {
        uuid id PK
        uuid device_id FK
        varchar target_phone
        text message_sent
        text response_received
        varchar status
    }

    USERS ||--o{ WEBHOOKS : owns
    DEVICES ||--o{ WEBHOOKS : scoped_to
    WEBHOOKS ||--o{ WEBHOOK_EVENT_SUBSCRIPTIONS : subscribes
    WEBHOOKS ||--o{ WEBHOOK_DELIVERIES : logs
    USERS ||--o{ API_KEYS : owns
    USERS ||--o{ API_LOGS : produces
    DEVICES ||--o{ API_LOGS : scoped_to
    USERS ||--o{ DAILY_MESSAGE_STATS : aggregates
    DEVICES ||--o{ DAILY_MESSAGE_STATS : measures
    USERS ||--o{ FAILURE_RECORDS : experiences
    DEVICES ||--o{ FAILURE_RECORDS : causes
    USERS ||--o{ ONBOARDING_PROGRESS : tracks
    USERS ||--o{ NOTIFICATION_SETTINGS : configures
    USERS ||--o{ NOTIFICATIONS : receives
    USERS ||--o{ AUDIT_LOGS : triggers
    DEVICES ||--o{ WARMING_POOL : has
    DEVICES ||--o{ WARMING_SESSIONS : runs
```

## 4. Notes

- `SERVICE_HEALTH_CHECKS` dan `RESOURCE_USAGE_METRICS` bersifat optional karena fitur System Health saat ini sudah tidak tampil di navigasi dashboard.
- `WARMING_POOL` dan `WARMING_SESSIONS` dipertahankan di ERD karena sudah ada di schema existing backend, walau belum terefleksi penuh di UI dashboard.
- `SCHEDULED_MESSAGE_RECIPIENTS` sengaja dipisah agar scheduled message tidak bergantung penuh pada JSON dan tetap queryable saat data membesar.
- `BROADCAST_MESSAGES` tidak dimasukkan ke ERD inti. Kalau satu campaign memang hanya satu payload, lebih baik field payload digabung ke `broadcast_campaigns`. Kalau nanti butuh multi-step payload, tabel itu bisa dihidupkan lagi sebagai `broadcast_contents`.
