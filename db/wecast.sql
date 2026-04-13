/*
 WACAST Database Schema - Refactored
 
 Versi ini telah direfactor berdasarkan:
   - DB_DESIGN_REVIEW_AND_TARGET_SCHEMA.md
   - DB_TARGET_ERD.md
   - DASHBOARD_DATA_TABLES.md
   - Fitur dashboard di src/app/(dashboard)/
 
 Perubahan utama dari schema lama:
   - Semua naming dikonversi ke snake_case konsisten
   - Foreign key, unique constraint, dan index ditambahkan
   - OTP dipisah dari tabel users ke otp_verifications
   - Contact many-to-many via pivot contact_group_members
   - Session device dipisah ke device_sessions
   - Billing dilengkapi: invoices, invoice_items, usage_quotas
   - Webhook dilengkapi: event_subscriptions, deliveries
   - API keys ditambahkan
   - Messaging dilengkapi: templates, scheduled, broadcast rapi
   - Analytics ditambahkan: daily_message_stats, failure_records
   - Notifikasi, onboarding, audit, warming dipertahankan
   - broadcast_messages dihapus (payload digabung ke broadcast_campaigns)
 
 Target: PostgreSQL 14+
 Date: 13/04/2026
*/


-- ============================================================
-- CLEAN UP EXISTING SCHEMA
-- ============================================================

DROP TABLE IF EXISTS "public"."warming_sessions" CASCADE;
DROP TABLE IF EXISTS "public"."warming_pool" CASCADE;
DROP TABLE IF EXISTS "public"."audit_logs" CASCADE;
DROP TABLE IF EXISTS "public"."notifications" CASCADE;
DROP TABLE IF EXISTS "public"."notification_settings" CASCADE;
DROP TABLE IF EXISTS "public"."onboarding_progress" CASCADE;
DROP TABLE IF EXISTS "public"."resource_usage_metrics" CASCADE;
DROP TABLE IF EXISTS "public"."service_health_checks" CASCADE;
DROP TABLE IF EXISTS "public"."failure_records" CASCADE;
DROP TABLE IF EXISTS "public"."daily_message_stats" CASCADE;
DROP TABLE IF EXISTS "public"."api_logs" CASCADE;
DROP TABLE IF EXISTS "public"."webhook_deliveries" CASCADE;
DROP TABLE IF EXISTS "public"."webhook_event_subscriptions" CASCADE;
DROP TABLE IF EXISTS "public"."webhooks" CASCADE;
DROP TABLE IF EXISTS "public"."api_keys" CASCADE;
DROP TABLE IF EXISTS "public"."auto_response_logs" CASCADE;
DROP TABLE IF EXISTS "public"."auto_response_keywords" CASCADE;
DROP TABLE IF EXISTS "public"."auto_response" CASCADE;
DROP TABLE IF EXISTS "public"."broadcast_recipients" CASCADE;
DROP TABLE IF EXISTS "public"."broadcast_campaigns" CASCADE;
DROP TABLE IF EXISTS "public"."broadcast_messages" CASCADE;
DROP TABLE IF EXISTS "public"."scheduled_message_recipients" CASCADE;
DROP TABLE IF EXISTS "public"."scheduled_messages" CASCADE;
DROP TABLE IF EXISTS "public"."messages" CASCADE;
DROP TABLE IF EXISTS "public"."message_templates" CASCADE;
DROP TABLE IF EXISTS "public"."blacklists" CASCADE;
DROP TABLE IF EXISTS "public"."contact_group_members" CASCADE;
DROP TABLE IF EXISTS "public"."contact_groups" CASCADE;
DROP TABLE IF EXISTS "public"."contact" CASCADE;
DROP TABLE IF EXISTS "public"."contacts" CASCADE;
DROP TABLE IF EXISTS "public"."contact_labels" CASCADE;
DROP TABLE IF EXISTS "public"."groups" CASCADE;
DROP TABLE IF EXISTS "public"."device_metrics" CASCADE;
DROP TABLE IF EXISTS "public"."device_qr_codes" CASCADE;
DROP TABLE IF EXISTS "public"."device_sessions" CASCADE;
DROP TABLE IF EXISTS "public"."devices" CASCADE;
DROP TABLE IF EXISTS "public"."usage_quotas" CASCADE;
DROP TABLE IF EXISTS "public"."invoice_items" CASCADE;
DROP TABLE IF EXISTS "public"."invoices" CASCADE;
DROP TABLE IF EXISTS "public"."subscriptions" CASCADE;
DROP TABLE IF EXISTS "public"."billing_plans" CASCADE;
DROP TABLE IF EXISTS "public"."user_sessions" CASCADE;
DROP TABLE IF EXISTS "public"."otp_verifications" CASCADE;
DROP TABLE IF EXISTS "public"."users" CASCADE;
DROP TABLE IF EXISTS "public"."system_settings" CASCADE;
DROP TABLE IF EXISTS "public"."lookup" CASCADE;


-- ============================================================
-- TYPES / ENUMS
-- ============================================================

DROP TYPE IF EXISTS "public"."enum_direction";
CREATE TYPE "public"."enum_direction" AS ENUM (
  'IN',
  'OUT'
);


-- ============================================================
-- DOMAIN 1: IDENTITY AND AUTH
-- ============================================================

-- ----------------------------
-- Table: users
-- ----------------------------
CREATE TABLE "public"."users" (
  "id"              uuid            NOT NULL DEFAULT gen_random_uuid(),
  "phone_number"    varchar(30)     NOT NULL,
  "full_name"       varchar(100)    NOT NULL,
  "email"           varchar(255)    NULL,
  "company_name"    varchar(100)    NULL,
  "timezone"        varchar(50)     NOT NULL DEFAULT 'Asia/Jakarta',
  "is_verified"     boolean         NOT NULL DEFAULT false,
  "is_banned"       boolean         NOT NULL DEFAULT false,
  "is_api_enabled"  boolean         NOT NULL DEFAULT false,
  "created_at"      timestamptz     NOT NULL DEFAULT now(),
  "updated_at"      timestamptz     NOT NULL DEFAULT now(),
  "last_login_at"   timestamptz     NULL,
  CONSTRAINT "users_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "users_phone_number_key" UNIQUE ("phone_number"),
  CONSTRAINT "users_email_key" UNIQUE ("email")
);

-- ----------------------------
-- Table: otp_verifications
-- ----------------------------
CREATE TABLE "public"."otp_verifications" (
  "id"              uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"         uuid            NULL,
  "phone_number"    varchar(30)     NOT NULL,
  "context"         varchar(50)     NOT NULL,  -- login, register, reset
  "otp_code_hash"   varchar(255)    NOT NULL,
  "attempt_count"   int             NOT NULL DEFAULT 0,
  "expires_at"      timestamptz     NOT NULL,
  "verified_at"     timestamptz     NULL,
  "created_at"      timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "otp_verifications_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "otp_verifications_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE SET NULL,
  CONSTRAINT "otp_verifications_attempt_count_check" CHECK ("attempt_count" >= 0)
);

-- ----------------------------
-- Table: user_sessions
-- ----------------------------
CREATE TABLE "public"."user_sessions" (
  "id"                  uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"             uuid            NOT NULL,
  "session_token_hash"  varchar(255)    NOT NULL,
  "refresh_token_hash"  varchar(255)    NULL,
  "ip_address"          inet            NULL,
  "user_agent"          text            NULL,
  "expires_at"          timestamptz     NOT NULL,
  "revoked_at"          timestamptz     NULL,
  "created_at"          timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "user_sessions_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "user_sessions_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE
);


-- ============================================================
-- DOMAIN 2: BILLING AND SUBSCRIPTION
-- ============================================================

-- ----------------------------
-- Table: billing_plans
-- ----------------------------
CREATE TABLE "public"."billing_plans" (
  "id"                    uuid            NOT NULL DEFAULT gen_random_uuid(),
  "name"                  varchar(100)    NOT NULL,
  "price"                 numeric(12,2)   NOT NULL DEFAULT 0,
  "max_devices"           int             NOT NULL DEFAULT 1,
  "max_messages_per_day"  int             NULL,
  "billing_cycle"         varchar(10)     NOT NULL DEFAULT 'monthly',  -- monthly, yearly
  "features"              jsonb           NOT NULL DEFAULT '{}',
  "is_active"             boolean         NOT NULL DEFAULT true,
  "created_at"            timestamptz     NOT NULL DEFAULT now(),
  "updated_at"            timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "billing_plans_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "billing_plans_name_key" UNIQUE ("name"),
  CONSTRAINT "billing_plans_price_check" CHECK ("price" >= 0),
  CONSTRAINT "billing_plans_max_devices_check" CHECK ("max_devices" > 0)
);

-- ----------------------------
-- Table: subscriptions
-- ----------------------------
CREATE TABLE "public"."subscriptions" (
  "id"            uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"       uuid            NOT NULL,
  "plan_id"       uuid            NOT NULL,
  "status"        varchar(20)     NOT NULL DEFAULT 'active',  -- active, expired, cancelled, trial
  "start_date"    timestamptz     NOT NULL DEFAULT now(),
  "end_date"      timestamptz     NULL,
  "renewal_date"  timestamptz     NULL,
  "auto_renew"    boolean         NOT NULL DEFAULT true,
  "created_at"    timestamptz     NOT NULL DEFAULT now(),
  "updated_at"    timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "subscriptions_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "subscriptions_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE RESTRICT,
  CONSTRAINT "subscriptions_plan_id_fkey"
    FOREIGN KEY ("plan_id") REFERENCES "public"."billing_plans" ("id") ON DELETE RESTRICT
);

-- ----------------------------
-- Table: usage_quotas
-- ----------------------------
CREATE TABLE "public"."usage_quotas" (
  "id"                uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"           uuid            NOT NULL,
  "subscription_id"   uuid            NOT NULL,
  "period_key"        varchar(20)     NOT NULL,  -- e.g. '2026-04'
  "messages_used"     int             NOT NULL DEFAULT 0,
  "messages_limit"    int             NULL,
  "devices_used"      int             NOT NULL DEFAULT 0,
  "devices_limit"     int             NULL,
  "created_at"        timestamptz     NOT NULL DEFAULT now(),
  "updated_at"        timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "usage_quotas_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "usage_quotas_user_id_period_key" UNIQUE ("user_id", "subscription_id", "period_key"),
  CONSTRAINT "usage_quotas_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE,
  CONSTRAINT "usage_quotas_subscription_id_fkey"
    FOREIGN KEY ("subscription_id") REFERENCES "public"."subscriptions" ("id") ON DELETE CASCADE,
  CONSTRAINT "usage_quotas_messages_used_check" CHECK ("messages_used" >= 0),
  CONSTRAINT "usage_quotas_devices_used_check" CHECK ("devices_used" >= 0)
);

-- ----------------------------
-- Table: invoices
-- ----------------------------
CREATE TABLE "public"."invoices" (
  "id"                uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"           uuid            NOT NULL,
  "subscription_id"   uuid            NULL,
  "invoice_number"    varchar(50)     NOT NULL,
  "issue_date"        date            NOT NULL DEFAULT current_date,
  "due_date"          date            NULL,
  "paid_at"           timestamptz     NULL,
  "amount"            numeric(12,2)   NOT NULL,
  "currency"          varchar(10)     NOT NULL DEFAULT 'IDR',
  "status"            varchar(20)     NOT NULL DEFAULT 'unpaid',  -- unpaid, paid, overdue, cancelled
  "payment_method"    varchar(50)     NULL,
  "created_at"        timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "invoices_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "invoices_invoice_number_key" UNIQUE ("invoice_number"),
  CONSTRAINT "invoices_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE RESTRICT,
  CONSTRAINT "invoices_subscription_id_fkey"
    FOREIGN KEY ("subscription_id") REFERENCES "public"."subscriptions" ("id") ON DELETE SET NULL,
  CONSTRAINT "invoices_amount_check" CHECK ("amount" >= 0)
);

-- ----------------------------
-- Table: invoice_items
-- ----------------------------
CREATE TABLE "public"."invoice_items" (
  "id"            uuid            NOT NULL DEFAULT gen_random_uuid(),
  "invoice_id"    uuid            NOT NULL,
  "description"   varchar(255)    NOT NULL,
  "qty"           int             NOT NULL DEFAULT 1,
  "unit_price"    numeric(12,2)   NOT NULL,
  "total_price"   numeric(12,2)   NOT NULL,
  CONSTRAINT "invoice_items_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "invoice_items_invoice_id_fkey"
    FOREIGN KEY ("invoice_id") REFERENCES "public"."invoices" ("id") ON DELETE CASCADE,
  CONSTRAINT "invoice_items_qty_check" CHECK ("qty" > 0),
  CONSTRAINT "invoice_items_unit_price_check" CHECK ("unit_price" >= 0)
);


-- ============================================================
-- DOMAIN 3: DEVICES AND WA SESSIONS
-- ============================================================

-- ----------------------------
-- Table: devices
-- ----------------------------
CREATE TABLE "public"."devices" (
  "id"                uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"           uuid            NOT NULL,
  "unique_name"       varchar(100)    NOT NULL,
  "display_name"      varchar(100)    NOT NULL,
  "phone_number"      varchar(30)     NULL,
  "status"            varchar(30)     NOT NULL DEFAULT 'disconnected',  -- connected, disconnected, pending_qr, banned
  "platform"          varchar(50)     NULL,
  "wa_version"        varchar(30)     NULL,
  "battery_level"     int             NULL,
  "last_seen_at"      timestamptz     NULL,
  "connected_since"   timestamptz     NULL,
  "created_at"        timestamptz     NOT NULL DEFAULT now(),
  "updated_at"        timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "devices_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "devices_user_id_unique_name_key" UNIQUE ("user_id", "unique_name"),
  CONSTRAINT "devices_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE,
  CONSTRAINT "devices_battery_level_check"
    CHECK ("battery_level" IS NULL OR ("battery_level" >= 0 AND "battery_level" <= 100))
);

-- ----------------------------
-- Table: device_sessions
-- ----------------------------
CREATE TABLE "public"."device_sessions" (
  "id"                uuid            NOT NULL DEFAULT gen_random_uuid(),
  "device_id"         uuid            NOT NULL,
  "session_blob"      bytea           NULL,
  "session_status"    varchar(20)     NOT NULL DEFAULT 'inactive',  -- active, inactive, expired
  "started_at"        timestamptz     NOT NULL DEFAULT now(),
  "ended_at"          timestamptz     NULL,
  "restart_count"     int             NOT NULL DEFAULT 0,
  "last_restart_at"   timestamptz     NULL,
  "created_at"        timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "device_sessions_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "device_sessions_device_id_fkey"
    FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id") ON DELETE CASCADE,
  CONSTRAINT "device_sessions_restart_count_check" CHECK ("restart_count" >= 0)
);

-- ----------------------------
-- Table: device_qr_codes
-- ----------------------------
CREATE TABLE "public"."device_qr_codes" (
  "id"            uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"       uuid            NOT NULL,
  "device_id"     uuid            NULL,
  "qr_string"     text            NULL,
  "qr_image_url"  text            NULL,
  "status"        varchar(20)     NOT NULL DEFAULT 'pending',  -- pending, scanned, expired
  "generated_at"  timestamptz     NOT NULL DEFAULT now(),
  "expired_at"    timestamptz     NULL,
  CONSTRAINT "device_qr_codes_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "device_qr_codes_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE,
  CONSTRAINT "device_qr_codes_device_id_fkey"
    FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id") ON DELETE SET NULL
);

-- ----------------------------
-- Table: device_metrics
-- ----------------------------
CREATE TABLE "public"."device_metrics" (
  "id"                        uuid            NOT NULL DEFAULT gen_random_uuid(),
  "device_id"                 uuid            NOT NULL,
  "uptime_seconds"            bigint          NOT NULL DEFAULT 0,
  "messages_sent_count"       int             NOT NULL DEFAULT 0,
  "messages_received_count"   int             NOT NULL DEFAULT 0,
  "success_rate"              numeric(5,2)    NULL,
  "recorded_at"               timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "device_metrics_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "device_metrics_device_id_fkey"
    FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id") ON DELETE CASCADE,
  CONSTRAINT "device_metrics_success_rate_check"
    CHECK ("success_rate" IS NULL OR ("success_rate" >= 0 AND "success_rate" <= 100))
);


-- ============================================================
-- DOMAIN 4: CONTACTS
-- ============================================================

-- ----------------------------
-- Table: contact_labels
-- ----------------------------
CREATE TABLE "public"."contact_labels" (
  "id"          uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"     uuid            NOT NULL,
  "name"        varchar(50)     NOT NULL,
  "color"       varchar(20)     NULL,
  "created_at"  timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "contact_labels_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "contact_labels_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE
);

-- ----------------------------
-- Table: contacts
-- ----------------------------
CREATE TABLE "public"."contacts" (
  "id"                uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"           uuid            NOT NULL,
  "label_id"          uuid            NULL,
  "name"              varchar(100)    NOT NULL,
  "phone_number"      varchar(30)     NOT NULL,
  "additional_data"   jsonb           NULL,
  "note"              text            NULL,
  "created_at"        timestamptz     NOT NULL DEFAULT now(),
  "updated_at"        timestamptz     NOT NULL DEFAULT now(),
  "deleted_at"        timestamptz     NULL,
  CONSTRAINT "contacts_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "contacts_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE,
  CONSTRAINT "contacts_label_id_fkey"
    FOREIGN KEY ("label_id") REFERENCES "public"."contact_labels" ("id") ON DELETE SET NULL
);

-- ----------------------------
-- Table: contact_groups
-- ----------------------------
CREATE TABLE "public"."contact_groups" (
  "id"          uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"     uuid            NOT NULL,
  "name"        varchar(100)    NOT NULL,
  "description" text            NULL,
  "created_at"  timestamptz     NOT NULL DEFAULT now(),
  "updated_at"  timestamptz     NOT NULL DEFAULT now(),
  "deleted_at"  timestamptz     NULL,
  CONSTRAINT "contact_groups_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "contact_groups_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE
);

-- ----------------------------
-- Table: contact_group_members
-- ----------------------------
CREATE TABLE "public"."contact_group_members" (
  "id"          uuid            NOT NULL DEFAULT gen_random_uuid(),
  "group_id"    uuid            NOT NULL,
  "contact_id"  uuid            NOT NULL,
  "created_at"  timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "contact_group_members_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "contact_group_members_group_id_contact_id_key" UNIQUE ("group_id", "contact_id"),
  CONSTRAINT "contact_group_members_group_id_fkey"
    FOREIGN KEY ("group_id") REFERENCES "public"."contact_groups" ("id") ON DELETE CASCADE,
  CONSTRAINT "contact_group_members_contact_id_fkey"
    FOREIGN KEY ("contact_id") REFERENCES "public"."contacts" ("id") ON DELETE CASCADE
);

-- ----------------------------
-- Table: blacklists
-- ----------------------------
CREATE TABLE "public"."blacklists" (
  "id"            uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"       uuid            NOT NULL,
  "phone_number"  varchar(30)     NOT NULL,
  "reason"        varchar(255)    NULL,
  "blocked_at"    timestamptz     NOT NULL DEFAULT now(),
  "unblocked_at"  timestamptz     NULL,
  "created_at"    timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "blacklists_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "blacklists_user_id_phone_number_key" UNIQUE ("user_id", "phone_number"),
  CONSTRAINT "blacklists_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE
);


-- ============================================================
-- DOMAIN 5: MESSAGING
-- ============================================================

-- ----------------------------
-- Table: message_templates
-- (dibuat sebelum messages, broadcast_campaigns, scheduled_messages)
-- ----------------------------
CREATE TABLE "public"."message_templates" (
  "id"          uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"     uuid            NOT NULL,
  "name"        varchar(100)    NOT NULL,
  "category"    varchar(50)     NOT NULL DEFAULT 'general',
  "content"     text            NOT NULL,
  "used_count"  int             NOT NULL DEFAULT 0,
  "created_at"  timestamptz     NOT NULL DEFAULT now(),
  "updated_at"  timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "message_templates_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "message_templates_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE,
  CONSTRAINT "message_templates_used_count_check" CHECK ("used_count" >= 0)
);

-- ----------------------------
-- Table: broadcast_campaigns
-- (dibuat sebelum messages karena messages FK ke broadcast_id)
-- ----------------------------
CREATE TABLE "public"."broadcast_campaigns" (
  "id"                uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"           uuid            NOT NULL,
  "device_id"         uuid            NOT NULL,
  "template_id"       uuid            NULL,
  "name"              varchar(255)    NOT NULL,
  "message_content"   text            NULL,
  "delay_seconds"     int             NOT NULL DEFAULT 0,
  "total_recipients"  int             NOT NULL DEFAULT 0,
  "processed_count"   int             NOT NULL DEFAULT 0,
  "success_count"     int             NOT NULL DEFAULT 0,
  "failed_count"      int             NOT NULL DEFAULT 0,
  "scheduled_at"      timestamptz     NULL,
  "started_at"        timestamptz     NULL,
  "completed_at"      timestamptz     NULL,
  "status"            varchar(20)     NOT NULL DEFAULT 'draft',  -- draft, queued, sending, completed, failed, cancelled
  "created_at"        timestamptz     NOT NULL DEFAULT now(),
  "updated_at"        timestamptz     NOT NULL DEFAULT now(),
  "deleted_at"        timestamptz     NULL,
  CONSTRAINT "broadcast_campaigns_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "broadcast_campaigns_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE,
  CONSTRAINT "broadcast_campaigns_device_id_fkey"
    FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id") ON DELETE RESTRICT,
  CONSTRAINT "broadcast_campaigns_template_id_fkey"
    FOREIGN KEY ("template_id") REFERENCES "public"."message_templates" ("id") ON DELETE SET NULL,
  CONSTRAINT "broadcast_campaigns_delay_check" CHECK ("delay_seconds" >= 0),
  CONSTRAINT "broadcast_campaigns_counts_check"
    CHECK ("total_recipients" >= 0 AND "processed_count" >= 0
      AND "success_count" >= 0 AND "failed_count" >= 0)
);

-- ----------------------------
-- Table: scheduled_messages
-- (dibuat sebelum messages karena messages FK ke scheduled_message_id)
-- ----------------------------
CREATE TABLE "public"."scheduled_messages" (
  "id"                uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"           uuid            NOT NULL,
  "device_id"         uuid            NOT NULL,
  "template_id"       uuid            NULL,
  "group_id"          uuid            NULL,
  "recipient_mode"    varchar(20)     NOT NULL DEFAULT 'single',  -- single, group, bulk
  "recipient_payload" jsonb           NULL,
  "message_content"   text            NOT NULL,
  "scheduled_at"      timestamptz     NOT NULL,
  "status"            varchar(20)     NOT NULL DEFAULT 'pending',  -- pending, sent, failed, cancelled
  "created_at"        timestamptz     NOT NULL DEFAULT now(),
  "updated_at"        timestamptz     NOT NULL DEFAULT now(),
  "executed_at"       timestamptz     NULL,
  CONSTRAINT "scheduled_messages_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "scheduled_messages_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE,
  CONSTRAINT "scheduled_messages_device_id_fkey"
    FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id") ON DELETE RESTRICT,
  CONSTRAINT "scheduled_messages_template_id_fkey"
    FOREIGN KEY ("template_id") REFERENCES "public"."message_templates" ("id") ON DELETE SET NULL,
  CONSTRAINT "scheduled_messages_group_id_fkey"
    FOREIGN KEY ("group_id") REFERENCES "public"."contact_groups" ("id") ON DELETE SET NULL
);

-- ----------------------------
-- Table: messages
-- ----------------------------
CREATE TABLE "public"."messages" (
  "id"                    uuid                      NOT NULL DEFAULT gen_random_uuid(),
  "user_id"               uuid                      NOT NULL,
  "device_id"             uuid                      NOT NULL,
  "template_id"           uuid                      NULL,
  "broadcast_id"          uuid                      NULL,
  "scheduled_message_id"  uuid                      NULL,
  "direction"             "public"."enum_direction" NOT NULL,
  "recipient_phone"       varchar(30)               NULL,
  "sender_phone"          varchar(30)               NULL,
  "message_type"          varchar(30)               NOT NULL DEFAULT 'text',  -- text, image, video, document, audio
  "content"               text                      NOT NULL,
  "status"                varchar(20)               NOT NULL DEFAULT 'pending',  -- pending, sent, delivered, read, failed, cancelled
  "error_log"             text                      NULL,
  "created_at"            timestamptz               NOT NULL DEFAULT now(),
  "sent_at"               timestamptz               NULL,
  "delivered_at"          timestamptz               NULL,
  "read_at"               timestamptz               NULL,
  "failed_at"             timestamptz               NULL,
  "updated_at"            timestamptz               NOT NULL DEFAULT now(),
  CONSTRAINT "messages_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "messages_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE,
  CONSTRAINT "messages_device_id_fkey"
    FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id") ON DELETE RESTRICT,
  CONSTRAINT "messages_template_id_fkey"
    FOREIGN KEY ("template_id") REFERENCES "public"."message_templates" ("id") ON DELETE SET NULL,
  CONSTRAINT "messages_broadcast_id_fkey"
    FOREIGN KEY ("broadcast_id") REFERENCES "public"."broadcast_campaigns" ("id") ON DELETE SET NULL,
  CONSTRAINT "messages_scheduled_message_id_fkey"
    FOREIGN KEY ("scheduled_message_id") REFERENCES "public"."scheduled_messages" ("id") ON DELETE SET NULL
);

-- ----------------------------
-- Table: scheduled_message_recipients
-- ----------------------------
CREATE TABLE "public"."scheduled_message_recipients" (
  "id"                    uuid            NOT NULL DEFAULT gen_random_uuid(),
  "scheduled_message_id"  uuid            NOT NULL,
  "contact_id"            uuid            NULL,
  "group_id"              uuid            NULL,
  "phone_number"          varchar(30)     NOT NULL,
  "status"                varchar(20)     NULL DEFAULT 'pending',
  "created_at"            timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "scheduled_message_recipients_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "scheduled_message_recipients_unique" UNIQUE ("scheduled_message_id", "phone_number"),
  CONSTRAINT "scheduled_message_recipients_scheduled_message_id_fkey"
    FOREIGN KEY ("scheduled_message_id") REFERENCES "public"."scheduled_messages" ("id") ON DELETE CASCADE,
  CONSTRAINT "scheduled_message_recipients_contact_id_fkey"
    FOREIGN KEY ("contact_id") REFERENCES "public"."contacts" ("id") ON DELETE SET NULL,
  CONSTRAINT "scheduled_message_recipients_group_id_fkey"
    FOREIGN KEY ("group_id") REFERENCES "public"."contact_groups" ("id") ON DELETE SET NULL
);

-- ----------------------------
-- Table: broadcast_recipients
-- ----------------------------
CREATE TABLE "public"."broadcast_recipients" (
  "id"            uuid            NOT NULL DEFAULT gen_random_uuid(),
  "campaign_id"   uuid            NOT NULL,
  "contact_id"    uuid            NULL,
  "group_id"      uuid            NULL,
  "phone_number"  varchar(30)     NOT NULL,
  "status"        varchar(20)     NOT NULL DEFAULT 'pending',  -- pending, sent, failed, skipped
  "sent_at"       timestamptz     NULL,
  "failed_at"     timestamptz     NULL,
  "error_message" text            NULL,
  "retry_count"   int             NOT NULL DEFAULT 0,
  CONSTRAINT "broadcast_recipients_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "broadcast_recipients_campaign_id_phone_key" UNIQUE ("campaign_id", "phone_number"),
  CONSTRAINT "broadcast_recipients_campaign_id_fkey"
    FOREIGN KEY ("campaign_id") REFERENCES "public"."broadcast_campaigns" ("id") ON DELETE CASCADE,
  CONSTRAINT "broadcast_recipients_contact_id_fkey"
    FOREIGN KEY ("contact_id") REFERENCES "public"."contacts" ("id") ON DELETE SET NULL,
  CONSTRAINT "broadcast_recipients_group_id_fkey"
    FOREIGN KEY ("group_id") REFERENCES "public"."contact_groups" ("id") ON DELETE SET NULL,
  CONSTRAINT "broadcast_recipients_retry_count_check" CHECK ("retry_count" >= 0)
);

-- ----------------------------
-- Table: auto_response_keywords
-- ----------------------------
CREATE TABLE "public"."auto_response_keywords" (
  "id"            uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"       uuid            NOT NULL,
  "device_id"     uuid            NULL,
  "keyword"       varchar(255)    NOT NULL,
  "response_text" text            NOT NULL,
  "is_active"     boolean         NOT NULL DEFAULT true,
  "created_at"    timestamptz     NOT NULL DEFAULT now(),
  "updated_at"    timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "auto_response_keywords_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "auto_response_keywords_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE,
  CONSTRAINT "auto_response_keywords_device_id_fkey"
    FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id") ON DELETE SET NULL
);

-- ----------------------------
-- Table: auto_response_logs
-- ----------------------------
CREATE TABLE "public"."auto_response_logs" (
  "id"                    uuid            NOT NULL DEFAULT gen_random_uuid(),
  "keyword_id"            uuid            NOT NULL,
  "message_id"            uuid            NULL,
  "triggered_by_phone"    varchar(30)     NOT NULL,
  "matched_keyword"       varchar(255)    NOT NULL,
  "response_sent"         text            NOT NULL,
  "created_at"            timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "auto_response_logs_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "auto_response_logs_keyword_id_fkey"
    FOREIGN KEY ("keyword_id") REFERENCES "public"."auto_response_keywords" ("id") ON DELETE CASCADE,
  CONSTRAINT "auto_response_logs_message_id_fkey"
    FOREIGN KEY ("message_id") REFERENCES "public"."messages" ("id") ON DELETE SET NULL
);


-- ============================================================
-- DOMAIN 6: API AND WEBHOOKS
-- ============================================================

-- ----------------------------
-- Table: api_keys
-- ----------------------------
CREATE TABLE "public"."api_keys" (
  "id"            uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"       uuid            NOT NULL,
  "name"          varchar(100)    NOT NULL,
  "key_prefix"    varchar(20)     NOT NULL,
  "key_hash"      varchar(255)    NOT NULL,
  "last_used_at"  timestamptz     NULL,
  "expires_at"    timestamptz     NULL,
  "is_active"     boolean         NOT NULL DEFAULT true,
  "created_at"    timestamptz     NOT NULL DEFAULT now(),
  "revoked_at"    timestamptz     NULL,
  CONSTRAINT "api_keys_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "api_keys_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE
);

-- ----------------------------
-- Table: webhooks
-- ----------------------------
CREATE TABLE "public"."webhooks" (
  "id"                uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"           uuid            NOT NULL,
  "device_id"         uuid            NULL,
  "webhook_url"       varchar(500)    NOT NULL,
  "secret_key_hash"   varchar(255)    NOT NULL,
  "is_active"         boolean         NOT NULL DEFAULT true,
  "created_at"        timestamptz     NOT NULL DEFAULT now(),
  "updated_at"        timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "webhooks_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "webhooks_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE,
  CONSTRAINT "webhooks_device_id_fkey"
    FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id") ON DELETE SET NULL
);

-- ----------------------------
-- Table: webhook_event_subscriptions
-- ----------------------------
CREATE TABLE "public"."webhook_event_subscriptions" (
  "id"          uuid            NOT NULL DEFAULT gen_random_uuid(),
  "webhook_id"  uuid            NOT NULL,
  "event_key"   varchar(100)    NOT NULL,  -- message.sent, message.received, device.connected
  "is_enabled"  boolean         NOT NULL DEFAULT true,
  CONSTRAINT "webhook_event_subscriptions_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "webhook_event_subscriptions_webhook_id_event_key_key" UNIQUE ("webhook_id", "event_key"),
  CONSTRAINT "webhook_event_subscriptions_webhook_id_fkey"
    FOREIGN KEY ("webhook_id") REFERENCES "public"."webhooks" ("id") ON DELETE CASCADE
);

-- ----------------------------
-- Table: webhook_deliveries
-- ----------------------------
CREATE TABLE "public"."webhook_deliveries" (
  "id"            uuid            NOT NULL DEFAULT gen_random_uuid(),
  "webhook_id"    uuid            NOT NULL,
  "event_key"     varchar(100)    NOT NULL,
  "payload"       jsonb           NOT NULL DEFAULT '{}',
  "attempt"       int             NOT NULL DEFAULT 1,
  "http_status"   int             NULL,
  "response_body" text            NULL,
  "status"        varchar(20)     NOT NULL DEFAULT 'pending',  -- pending, success, failed, retrying
  "sent_at"       timestamptz     NULL,
  "created_at"    timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "webhook_deliveries_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "webhook_deliveries_webhook_id_fkey"
    FOREIGN KEY ("webhook_id") REFERENCES "public"."webhooks" ("id") ON DELETE CASCADE,
  CONSTRAINT "webhook_deliveries_attempt_check" CHECK ("attempt" >= 1)
);

-- ----------------------------
-- Table: api_logs
-- ----------------------------
CREATE TABLE "public"."api_logs" (
  "id"            uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"       uuid            NULL,
  "device_id"     uuid            NULL,
  "endpoint"      varchar(255)    NULL,
  "req_body"      jsonb           NULL,
  "response_body" jsonb           NULL,
  "ip_address"    inet            NULL,
  "created_at"    timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "api_logs_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "api_logs_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE SET NULL,
  CONSTRAINT "api_logs_device_id_fkey"
    FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id") ON DELETE SET NULL
);


-- ============================================================
-- DOMAIN 7: MONITORING AND ANALYTICS
-- ============================================================

-- ----------------------------
-- Table: daily_message_stats
-- ----------------------------
CREATE TABLE "public"."daily_message_stats" (
  "id"                uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"           uuid            NOT NULL,
  "device_id"         uuid            NULL,
  "stat_date"         date            NOT NULL,
  "sent_count"        int             NOT NULL DEFAULT 0,
  "failed_count"      int             NOT NULL DEFAULT 0,
  "delivered_count"   int             NOT NULL DEFAULT 0,
  "received_count"    int             NOT NULL DEFAULT 0,
  "success_rate"      numeric(5,2)    NULL,
  "created_at"        timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "daily_message_stats_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "daily_message_stats_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE,
  CONSTRAINT "daily_message_stats_device_id_fkey"
    FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id") ON DELETE SET NULL,
  CONSTRAINT "daily_message_stats_counts_check"
    CHECK ("sent_count" >= 0 AND "failed_count" >= 0
      AND "delivered_count" >= 0 AND "received_count" >= 0)
);

-- ----------------------------
-- Table: failure_records
-- ----------------------------
CREATE TABLE "public"."failure_records" (
  "id"                uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"           uuid            NOT NULL,
  "device_id"         uuid            NULL,
  "message_id"        uuid            NULL,
  "recipient_phone"   varchar(30)     NOT NULL,
  "failure_type"      varchar(50)     NOT NULL,  -- send_failed, timeout, invalid_number, banned
  "failure_reason"    text            NULL,
  "occurred_at"       timestamptz     NOT NULL DEFAULT now(),
  "created_at"        timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "failure_records_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "failure_records_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE,
  CONSTRAINT "failure_records_device_id_fkey"
    FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id") ON DELETE SET NULL,
  CONSTRAINT "failure_records_message_id_fkey"
    FOREIGN KEY ("message_id") REFERENCES "public"."messages" ("id") ON DELETE SET NULL
);

-- ----------------------------
-- Table: service_health_checks
-- (opsional - halaman /monitoring/health masih ada di source code)
-- ----------------------------
CREATE TABLE "public"."service_health_checks" (
  "id"                uuid            NOT NULL DEFAULT gen_random_uuid(),
  "service_name"      varchar(100)    NOT NULL,
  "status"            varchar(20)     NOT NULL DEFAULT 'unknown',  -- healthy, degraded, down, unknown
  "latency_ms"        int             NULL,
  "uptime_percent"    numeric(5,2)    NULL,
  "started_at"        timestamptz     NULL,
  "last_incident_at"  timestamptz     NULL,
  "checked_at"        timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "service_health_checks_pkey" PRIMARY KEY ("id")
);

-- ----------------------------
-- Table: resource_usage_metrics
-- (opsional)
-- ----------------------------
CREATE TABLE "public"."resource_usage_metrics" (
  "id"            uuid            NOT NULL DEFAULT gen_random_uuid(),
  "metric_name"   varchar(100)    NOT NULL,
  "metric_value"  numeric(12,4)   NOT NULL,
  "metric_unit"   varchar(20)     NOT NULL,
  "recorded_at"   timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "resource_usage_metrics_pkey" PRIMARY KEY ("id")
);


-- ============================================================
-- DOMAIN 8: APP SUPPORT
-- ============================================================

-- ----------------------------
-- Table: onboarding_progress
-- ----------------------------
CREATE TABLE "public"."onboarding_progress" (
  "id"            uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"       uuid            NOT NULL,
  "step_key"      varchar(50)     NOT NULL,
  "is_completed"  boolean         NOT NULL DEFAULT false,
  "completed_at"  timestamptz     NULL,
  "created_at"    timestamptz     NOT NULL DEFAULT now(),
  "updated_at"    timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "onboarding_progress_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "onboarding_progress_user_id_step_key_key" UNIQUE ("user_id", "step_key"),
  CONSTRAINT "onboarding_progress_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE
);

-- ----------------------------
-- Table: notification_settings
-- ----------------------------
CREATE TABLE "public"."notification_settings" (
  "id"              uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"         uuid            NOT NULL,
  "event_key"       varchar(100)    NOT NULL,
  "email_enabled"   boolean         NOT NULL DEFAULT false,
  "in_app_enabled"  boolean         NOT NULL DEFAULT true,
  "created_at"      timestamptz     NOT NULL DEFAULT now(),
  "updated_at"      timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "notification_settings_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "notification_settings_user_id_event_key_key" UNIQUE ("user_id", "event_key"),
  CONSTRAINT "notification_settings_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE
);

-- ----------------------------
-- Table: notifications
-- ----------------------------
CREATE TABLE "public"."notifications" (
  "id"          uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"     uuid            NOT NULL,
  "type"        varchar(20)     NOT NULL DEFAULT 'info',  -- info, success, warning, alert
  "title"       varchar(255)    NOT NULL,
  "body"        text            NOT NULL,
  "is_read"     boolean         NOT NULL DEFAULT false,
  "read_at"     timestamptz     NULL,
  "created_at"  timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "notifications_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "notifications_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE
);

-- ----------------------------
-- Table: audit_logs
-- ----------------------------
CREATE TABLE "public"."audit_logs" (
  "id"            uuid            NOT NULL DEFAULT gen_random_uuid(),
  "user_id"       uuid            NULL,
  "action_type"   varchar(50)     NOT NULL,  -- create, update, delete, login, logout
  "resource_type" varchar(50)     NOT NULL,  -- user, device, contact, broadcast
  "resource_id"   uuid            NULL,
  "metadata"      jsonb           NULL,
  "ip_address"    inet            NULL,
  "user_agent"    text            NULL,
  "created_at"    timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "audit_logs_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "audit_logs_user_id_fkey"
    FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE SET NULL
);


-- ============================================================
-- DOMAIN 9: WARMING (PRESERVED FROM EXISTING BACKEND)
-- ============================================================

-- ----------------------------
-- Table: warming_pool
-- ----------------------------
CREATE TABLE "public"."warming_pool" (
  "id"                    uuid            NOT NULL DEFAULT gen_random_uuid(),
  "device_id"             uuid            NOT NULL,
  "intensity"             int             NOT NULL DEFAULT 1,
  "daily_limit"           int             NOT NULL DEFAULT 10,
  "message_send_today"    int             NOT NULL DEFAULT 0,
  "is_active"             boolean         NOT NULL DEFAULT true,
  "next_action_at"        timestamptz     NULL,
  CONSTRAINT "warming_pool_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "warming_pool_device_id_fkey"
    FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id") ON DELETE CASCADE,
  CONSTRAINT "warming_pool_daily_limit_check" CHECK ("daily_limit" > 0),
  CONSTRAINT "warming_pool_message_send_today_check" CHECK ("message_send_today" >= 0)
);

-- ----------------------------
-- Table: warming_sessions
-- ----------------------------
CREATE TABLE "public"."warming_sessions" (
  "id"                uuid            NOT NULL DEFAULT gen_random_uuid(),
  "device_id"         uuid            NOT NULL,
  "target_phone"      varchar(30)     NOT NULL,
  "message_sent"      varchar(500)    NULL,
  "response_received" varchar(500)    NULL,
  "status"            varchar(20)     NOT NULL DEFAULT 'pending',  -- pending, sent, replied, failed
  "created_at"        timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "warming_sessions_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "warming_sessions_device_id_fkey"
    FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id") ON DELETE CASCADE
);


-- ============================================================
-- DOMAIN 10: SYSTEM
-- ============================================================

-- ----------------------------
-- Table: system_settings
-- ----------------------------
CREATE TABLE "public"."system_settings" (
  "id"          int             NOT NULL GENERATED ALWAYS AS IDENTITY (
                                  INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1
                                ),
  "key"         varchar(100)    NOT NULL,
  "value"       text            NULL,
  "description" varchar(255)    NULL,
  "created_at"  timestamptz     NOT NULL DEFAULT now(),
  "updated_at"  timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "system_settings_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "system_settings_key_key" UNIQUE ("key")
);

-- ----------------------------
-- Table: lookup
-- ----------------------------
CREATE TABLE "public"."lookup" (
  "id"      int             NOT NULL GENERATED ALWAYS AS IDENTITY (
                              INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1
                            ),
  "key"     varchar(100)    NOT NULL,
  "value"   varchar(255)    NOT NULL,
  CONSTRAINT "lookup_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "lookup_key_value_key" UNIQUE ("key", "value")
);


-- ============================================================
-- INDEXES
-- ============================================================

-- Domain 1: Identity
CREATE INDEX "idx_user_sessions_user_id"         ON "public"."user_sessions" ("user_id");
CREATE INDEX "idx_user_sessions_expires_at"       ON "public"."user_sessions" ("expires_at");
CREATE INDEX "idx_otp_phone_number"               ON "public"."otp_verifications" ("phone_number");
CREATE INDEX "idx_otp_expires_at"                 ON "public"."otp_verifications" ("expires_at");

-- Domain 2: Billing
CREATE INDEX "idx_subscriptions_user_id"          ON "public"."subscriptions" ("user_id");
CREATE INDEX "idx_subscriptions_user_status"      ON "public"."subscriptions" ("user_id", "status");
CREATE INDEX "idx_usage_quotas_user_period"       ON "public"."usage_quotas" ("user_id", "period_key");
CREATE INDEX "idx_invoices_user_id"               ON "public"."invoices" ("user_id");
CREATE INDEX "idx_invoices_status"                ON "public"."invoices" ("status");

-- Domain 3: Devices
CREATE INDEX "idx_devices_user_status"            ON "public"."devices" ("user_id", "status");
CREATE INDEX "idx_devices_user_last_seen"         ON "public"."devices" ("user_id", "last_seen_at" DESC);
CREATE INDEX "idx_device_sessions_device_id"      ON "public"."device_sessions" ("device_id");
CREATE INDEX "idx_device_qr_user_id"              ON "public"."device_qr_codes" ("user_id");
CREATE INDEX "idx_device_qr_device_id"            ON "public"."device_qr_codes" ("device_id") WHERE "device_id" IS NOT NULL;
CREATE INDEX "idx_device_metrics_device_id"       ON "public"."device_metrics" ("device_id");

-- Domain 4: Contacts
CREATE INDEX "idx_contacts_user_name"             ON "public"."contacts" ("user_id", "name");
CREATE INDEX "idx_contacts_label_id"              ON "public"."contacts" ("label_id") WHERE "label_id" IS NOT NULL;
-- Partial unique: active contacts saja (soft delete)
CREATE UNIQUE INDEX "idx_contacts_user_phone_active" ON "public"."contacts" ("user_id", "phone_number") WHERE "deleted_at" IS NULL;
CREATE INDEX "idx_contact_groups_user_id"         ON "public"."contact_groups" ("user_id");
CREATE INDEX "idx_blacklists_user_id"             ON "public"."blacklists" ("user_id");

-- Domain 5: Messaging
CREATE INDEX "idx_messages_user_created"          ON "public"."messages" ("user_id", "created_at" DESC);
CREATE INDEX "idx_messages_device_created"        ON "public"."messages" ("device_id", "created_at" DESC);
CREATE INDEX "idx_messages_status_created"        ON "public"."messages" ("status", "created_at" DESC);
CREATE INDEX "idx_messages_recipient_phone"       ON "public"."messages" ("recipient_phone");
CREATE INDEX "idx_messages_broadcast_id"          ON "public"."messages" ("broadcast_id") WHERE "broadcast_id" IS NOT NULL;
CREATE INDEX "idx_messages_scheduled_id"          ON "public"."messages" ("scheduled_message_id") WHERE "scheduled_message_id" IS NOT NULL;
CREATE INDEX "idx_broadcast_campaigns_user"       ON "public"."broadcast_campaigns" ("user_id", "created_at" DESC);
CREATE INDEX "idx_broadcast_campaigns_status"     ON "public"."broadcast_campaigns" ("status");
CREATE INDEX "idx_broadcast_recip_campaign"       ON "public"."broadcast_recipients" ("campaign_id", "status");
CREATE INDEX "idx_scheduled_messages_user"        ON "public"."scheduled_messages" ("user_id", "scheduled_at");
CREATE INDEX "idx_scheduled_messages_status"      ON "public"."scheduled_messages" ("status");
CREATE INDEX "idx_message_templates_user_id"      ON "public"."message_templates" ("user_id");
CREATE INDEX "idx_auto_response_user_id"          ON "public"."auto_response_keywords" ("user_id");
CREATE INDEX "idx_auto_response_device_id"        ON "public"."auto_response_keywords" ("device_id") WHERE "device_id" IS NOT NULL;

-- Domain 6: API and Webhooks
CREATE INDEX "idx_api_keys_user_active"           ON "public"."api_keys" ("user_id", "is_active");
CREATE INDEX "idx_webhooks_user_active"           ON "public"."webhooks" ("user_id", "is_active");
CREATE INDEX "idx_webhook_deliveries_webhook"     ON "public"."webhook_deliveries" ("webhook_id", "created_at" DESC);
CREATE INDEX "idx_api_logs_user_created"          ON "public"."api_logs" ("user_id", "created_at" DESC);
CREATE INDEX "idx_api_logs_created_at"            ON "public"."api_logs" ("created_at" DESC);

-- Domain 7: Monitoring
-- Partial unique agar NULL-safe pada device_id nullable
CREATE UNIQUE INDEX "idx_daily_stats_user_no_device"    ON "public"."daily_message_stats" ("user_id", "stat_date") WHERE "device_id" IS NULL;
CREATE UNIQUE INDEX "idx_daily_stats_user_device"       ON "public"."daily_message_stats" ("user_id", "device_id", "stat_date") WHERE "device_id" IS NOT NULL;
CREATE INDEX "idx_daily_stats_user_date"          ON "public"."daily_message_stats" ("user_id", "stat_date" DESC);
CREATE INDEX "idx_failure_records_user"           ON "public"."failure_records" ("user_id", "occurred_at" DESC);
CREATE INDEX "idx_failure_records_device"         ON "public"."failure_records" ("device_id") WHERE "device_id" IS NOT NULL;

-- Domain 8: App Support
CREATE INDEX "idx_notifications_user_unread"      ON "public"."notifications" ("user_id", "is_read");
CREATE INDEX "idx_notifications_user_created"     ON "public"."notifications" ("user_id", "created_at" DESC);
CREATE INDEX "idx_audit_logs_user_created"        ON "public"."audit_logs" ("user_id", "created_at" DESC) WHERE "user_id" IS NOT NULL;
CREATE INDEX "idx_audit_logs_resource"            ON "public"."audit_logs" ("resource_type", "resource_id") WHERE "resource_id" IS NOT NULL;

-- Domain 9: Warming
CREATE INDEX "idx_warming_pool_device_id"         ON "public"."warming_pool" ("device_id");
CREATE INDEX "idx_warming_sessions_device_id"     ON "public"."warming_sessions" ("device_id");
