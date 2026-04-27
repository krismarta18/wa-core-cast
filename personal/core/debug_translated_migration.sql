/*
 Navicat Premium Dump SQL

 Source Server         : LocalDev
 Source Server Type    : PostgreSQL
 Source Server Version : 180003 (180003)
 Source Host           : localhost:5432
 Source Catalog        : wacast
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 180003 (180003)
 File Encoding         : 65001

 Date: 27/04/2026 15:36:04
*/


-- ----------------------------
-- Type structure for enum_direction
-- ----------------------------

-- ----------------------------
-- Sequence structure for lookup_id_seq
-- ----------------------------

-- ----------------------------
-- Sequence structure for migrations_id_seq
-- ----------------------------

-- ----------------------------
-- Sequence structure for system_settings_id_seq
-- ----------------------------

-- ----------------------------
-- Table structure for api_keys
-- ----------------------------
DROP TABLE IF EXISTS "api_keys";
CREATE TABLE "api_keys" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "user_id" TEXT NOT NULL,
  "name" varchar(100)  NOT NULL,
  "prefix" varchar(20)  NOT NULL,
  "key_hash" varchar(255)  NOT NULL,
  "last_used_at" DATETIME,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" DATETIME
)
;

-- ----------------------------
-- Table structure for api_logs
-- ----------------------------
DROP TABLE IF EXISTS "api_logs";
CREATE TABLE "api_logs" (
  "id" TEXT NOT NULL,
  "user_id" TEXT,
  "endpoint" varchar(255) ,
  "req_body" TEXT,
  "response_body" TEXT,
  "created_at" DATETIME,
  "ip_address" varchar(255) ,
  "device_id" TEXT
)
;

-- ----------------------------
-- Table structure for audit_logs
-- ----------------------------
DROP TABLE IF EXISTS "audit_logs";
CREATE TABLE "audit_logs" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "user_id" TEXT,
  "action_type" varchar(50)  NOT NULL,
  "resource_type" varchar(50)  NOT NULL,
  "resource_id" TEXT,
  "metadata" TEXT,
  "ip_address" TEXT,
  "user_agent" text ,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for auto_response
-- ----------------------------
DROP TABLE IF EXISTS "auto_response";
CREATE TABLE "auto_response" (
  "id" TEXT NOT NULL,
  "device_id" TEXT,
  "keyword" varchar(255) ,
  "response_text" varchar(255) ,
  "is_active" BOOLEAN,
  "created_at" DATETIME
)
;

-- ----------------------------
-- Table structure for auto_response_keywords
-- ----------------------------
DROP TABLE IF EXISTS "auto_response_keywords";
CREATE TABLE "auto_response_keywords" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "user_id" TEXT NOT NULL,
  "device_id" TEXT,
  "keyword" varchar(255)  NOT NULL,
  "response_text" text  NOT NULL,
  "is_active" BOOLEAN NOT NULL DEFAULT 1,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "match_type" varchar(20)  DEFAULT 'contains'
)
;

-- ----------------------------
-- Table structure for auto_response_logs
-- ----------------------------
DROP TABLE IF EXISTS "auto_response_logs";
CREATE TABLE "auto_response_logs" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "keyword_id" TEXT NOT NULL,
  "message_id" TEXT,
  "triggered_by_phone" varchar(30)  NOT NULL,
  "matched_keyword" varchar(255)  NOT NULL,
  "response_sent" text  NOT NULL,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for billing_plans
-- ----------------------------
DROP TABLE IF EXISTS "billing_plans";
CREATE TABLE "billing_plans" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "name" varchar(100)  NOT NULL,
  "price" numeric(10,2) NOT NULL,
  "max_devices" INTEGER NOT NULL,
  "max_messages_per_day" INTEGER NOT NULL,
  "features" TEXT,
  "billing_cycle" varchar(20)  NOT NULL DEFAULT 'monthly',
  "is_active" BOOLEAN NOT NULL DEFAULT 1,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for blacklists
-- ----------------------------
DROP TABLE IF EXISTS "blacklists";
CREATE TABLE "blacklists" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "user_id" TEXT NOT NULL,
  "phone_number" varchar(30)  NOT NULL,
  "reason" varchar(255) ,
  "blocked_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "unblocked_at" DATETIME,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for broadcast_campaigns
-- ----------------------------
DROP TABLE IF EXISTS "broadcast_campaigns";
CREATE TABLE "broadcast_campaigns" (
  "id" TEXT NOT NULL,
  "user_id" TEXT,
  "device_id" TEXT,
  "name_broadcast" varchar(255) ,
  "total_recipients" INTEGER,
  "processed_count" INTEGER,
  "scheduled_at" DATETIME,
  "status" varchar(20)  DEFAULT 'draft',
  "created_at" DATETIME,
  "deleted_at" DATETIME,
  "template_id" TEXT,
  "name" varchar(255) ,
  "message_content" text ,
  "delay_seconds" INTEGER NOT NULL DEFAULT 5,
  "success_count" INTEGER NOT NULL DEFAULT 0,
  "failed_count" INTEGER NOT NULL DEFAULT 0,
  "started_at" DATETIME,
  "completed_at" DATETIME,
  "updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for broadcast_messages
-- ----------------------------
DROP TABLE IF EXISTS "broadcast_messages";
CREATE TABLE "broadcast_messages" (
  "id" TEXT NOT NULL,
  "campaign_id" TEXT,
  "message_type" INTEGER,
  "message_text" varchar(255) ,
  "media_url" varchar(255) ,
  "button_data" TEXT
)
;

-- ----------------------------
-- Table structure for broadcast_recipients
-- ----------------------------
DROP TABLE IF EXISTS "broadcast_recipients";
CREATE TABLE "broadcast_recipients" (
  "id" TEXT NOT NULL,
  "campaign_id" TEXT,
  "groups_id" TEXT,
  "contact_id" TEXT,
  "status" varchar(20)  DEFAULT 'pending',
  "sent_at" DATETIME,
  "error_messages" varchar(255) ,
  "retry_count" INTEGER,
  "created_at" DATETIME,
  "group_id" TEXT,
  "phone_number" varchar(30) ,
  "failed_at" DATETIME,
  "error_message" text 
)
;

-- ----------------------------
-- Table structure for contact
-- ----------------------------
DROP TABLE IF EXISTS "contact";
CREATE TABLE "contact" (
  "id" TEXT NOT NULL,
  "group_id" TEXT,
  "name" varchar(100) ,
  "phone" varchar(30) ,
  "additional_data" TEXT,
  "created_at" DATETIME,
  "updated_at" DATETIME,
  "deleted_at" DATETIME
)
;

-- ----------------------------
-- Table structure for contact_group_members
-- ----------------------------
DROP TABLE IF EXISTS "contact_group_members";
CREATE TABLE "contact_group_members" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "group_id" TEXT NOT NULL,
  "contact_id" TEXT NOT NULL,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for contact_groups
-- ----------------------------
DROP TABLE IF EXISTS "contact_groups";
CREATE TABLE "contact_groups" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "user_id" TEXT NOT NULL,
  "name" varchar(100)  NOT NULL,
  "description" text ,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" DATETIME
)
;

-- ----------------------------
-- Table structure for contact_labels
-- ----------------------------
DROP TABLE IF EXISTS "contact_labels";
CREATE TABLE "contact_labels" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "user_id" TEXT NOT NULL,
  "name" varchar(50)  NOT NULL,
  "color" varchar(20) ,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for contacts
-- ----------------------------
DROP TABLE IF EXISTS "contacts";
CREATE TABLE "contacts" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "user_id" TEXT NOT NULL,
  "label_id" TEXT,
  "name" varchar(100)  NOT NULL,
  "phone_number" varchar(30)  NOT NULL,
  "additional_data" TEXT,
  "note" text ,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" DATETIME
)
;

-- ----------------------------
-- Table structure for daily_message_stats
-- ----------------------------
DROP TABLE IF EXISTS "daily_message_stats";
CREATE TABLE "daily_message_stats" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "user_id" TEXT NOT NULL,
  "device_id" TEXT,
  "stat_date" date NOT NULL,
  "sent_count" INTEGER NOT NULL DEFAULT 0,
  "failed_count" INTEGER NOT NULL DEFAULT 0,
  "delivered_count" INTEGER NOT NULL DEFAULT 0,
  "received_count" INTEGER NOT NULL DEFAULT 0,
  "success_rate" numeric(5,2),
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for device_metrics
-- ----------------------------
DROP TABLE IF EXISTS "device_metrics";
CREATE TABLE "device_metrics" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "device_id" TEXT NOT NULL,
  "uptime_seconds" INTEGER NOT NULL DEFAULT 0,
  "messages_sent_count" INTEGER NOT NULL DEFAULT 0,
  "messages_received_count" INTEGER NOT NULL DEFAULT 0,
  "success_rate" numeric(5,2),
  "recorded_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for device_qr_codes
-- ----------------------------
DROP TABLE IF EXISTS "device_qr_codes";
CREATE TABLE "device_qr_codes" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "user_id" TEXT NOT NULL,
  "device_id" TEXT,
  "qr_string" text ,
  "qr_image_url" text ,
  "status" varchar(20)  NOT NULL DEFAULT 'pending',
  "generated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "expired_at" DATETIME
)
;

-- ----------------------------
-- Table structure for device_sessions
-- ----------------------------
DROP TABLE IF EXISTS "device_sessions";
CREATE TABLE "device_sessions" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "device_id" TEXT NOT NULL,
  "session_blob" BLOB,
  "session_status" varchar(20)  NOT NULL DEFAULT 'inactive',
  "started_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "ended_at" DATETIME,
  "restart_count" INTEGER NOT NULL DEFAULT 0,
  "last_restart_at" DATETIME,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for devices
-- ----------------------------
DROP TABLE IF EXISTS "devices";
CREATE TABLE "devices" (
  "id" TEXT NOT NULL,
  "user_id" TEXT,
  "unique_name" varchar(100) ,
  "display_name" varchar(100) ,
  "phone_number" varchar(30) ,
  "status" varchar(20) ,
  "last_seen_at" DATETIME,
  "session_data" BLOB,
  "created_at" DATETIME,
  "updated_at" DATETIME,
  "connected_since" DATETIME,
  "platform" varchar(255) ,
  "wa_version" varchar(50) ,
  "battery_level" INTEGER
)
;

-- ----------------------------
-- Table structure for failure_records
-- ----------------------------
DROP TABLE IF EXISTS "failure_records";
CREATE TABLE "failure_records" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "user_id" TEXT NOT NULL,
  "device_id" TEXT,
  "message_id" TEXT,
  "recipient_phone" varchar(30)  NOT NULL,
  "failure_type" varchar(50)  NOT NULL,
  "failure_reason" text ,
  "occurred_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for groups
-- ----------------------------
DROP TABLE IF EXISTS "groups";
CREATE TABLE "groups" (
  "id" TEXT NOT NULL,
  "user_id" TEXT,
  "group_name" varchar(100) ,
  "created_at" DATETIME,
  "deleted_at" DATETIME
)
;

-- ----------------------------
-- Table structure for invoice_items
-- ----------------------------
DROP TABLE IF EXISTS "invoice_items";
CREATE TABLE "invoice_items" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "invoice_id" TEXT NOT NULL,
  "description" varchar(255)  NOT NULL,
  "qty" INTEGER NOT NULL DEFAULT 1,
  "unit_price" numeric(12,2) NOT NULL,
  "total_price" numeric(12,2) NOT NULL
)
;

-- ----------------------------
-- Table structure for invoices
-- ----------------------------
DROP TABLE IF EXISTS "invoices";
CREATE TABLE "invoices" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "user_id" TEXT NOT NULL,
  "subscription_id" TEXT,
  "invoice_number" varchar(50)  NOT NULL,
  "issue_date" date NOT NULL DEFAULT CURRENT_DATE,
  "due_date" date,
  "paid_at" DATETIME,
  "amount" numeric(12,2) NOT NULL,
  "currency" varchar(10)  NOT NULL DEFAULT 'IDR',
  "status" varchar(20)  NOT NULL DEFAULT 'unpaid',
  "payment_method" varchar(50) ,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for lookup
-- ----------------------------
DROP TABLE IF EXISTS "lookup";
CREATE TABLE "lookup" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT  ,
  "keys" varchar(255) ,
  "values" varchar(100) 
)
;

-- ----------------------------
-- Table structure for message_templates
-- ----------------------------
DROP TABLE IF EXISTS "message_templates";
CREATE TABLE "message_templates" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "user_id" TEXT NOT NULL,
  "name" varchar(100)  NOT NULL,
  "category" varchar(50)  NOT NULL DEFAULT 'general',
  "content" text  NOT NULL,
  "used_count" INTEGER NOT NULL DEFAULT 0,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for messages
-- ----------------------------
DROP TABLE IF EXISTS "messages";
CREATE TABLE "messages" (
  "id" TEXT NOT NULL,
  "device_id" TEXT,
  "direction" varchar(10) ,
  "receipt_number" text ,
  "message_type" INTEGER,
  "content" text ,
  "status_message" INTEGER,
  "error_log" text ,
  "created_at" DATETIME,
  "target_jid" varchar(100) ,
  "updated_at" DATETIME DEFAULT CURRENT_TIMESTAMP,
  "priority" INTEGER DEFAULT 3,
  "retry_count" INTEGER DEFAULT 0,
  "max_retries" INTEGER DEFAULT 3,
  "scheduled_for" DATETIME,
  "whatsapp_message_id" varchar(255) ,
  "media_url" text ,
  "caption" text ,
  "broadcast_id" TEXT,
  "scheduled_message_id" TEXT
)
;

-- ----------------------------
-- Table structure for migrations
-- ----------------------------
DROP TABLE IF EXISTS "migrations";
CREATE TABLE "migrations" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT  ,
  "version" varchar(255)  NOT NULL,
  "name" varchar(255)  NOT NULL,
  "applied_at" DATETIME NOT NULL
)
;

-- ----------------------------
-- Table structure for notification_settings
-- ----------------------------
DROP TABLE IF EXISTS "notification_settings";
CREATE TABLE "notification_settings" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "user_id" TEXT NOT NULL,
  "event_key" varchar(100)  NOT NULL,
  "email_enabled" BOOLEAN NOT NULL DEFAULT 0,
  "in_app_enabled" BOOLEAN NOT NULL DEFAULT 1,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for notifications
-- ----------------------------
DROP TABLE IF EXISTS "notifications";
CREATE TABLE "notifications" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "user_id" TEXT NOT NULL,
  "type" varchar(20)  NOT NULL DEFAULT 'info',
  "title" varchar(255)  NOT NULL,
  "body" text  NOT NULL,
  "is_read" BOOLEAN NOT NULL DEFAULT 0,
  "read_at" DATETIME,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for onboarding_progress
-- ----------------------------
DROP TABLE IF EXISTS "onboarding_progress";
CREATE TABLE "onboarding_progress" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "user_id" TEXT NOT NULL,
  "step_key" varchar(50)  NOT NULL,
  "is_completed" BOOLEAN NOT NULL DEFAULT 0,
  "completed_at" DATETIME,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for otp_verifications
-- ----------------------------
DROP TABLE IF EXISTS "otp_verifications";
CREATE TABLE "otp_verifications" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "user_id" TEXT,
  "phone_number" varchar(30)  NOT NULL,
  "context" varchar(50)  NOT NULL,
  "otp_code" varchar(10)  NOT NULL,
  "attempt_count" INTEGER NOT NULL DEFAULT 0,
  "expires_at" DATETIME NOT NULL,
  "verified_at" DATETIME,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for resource_usage_metrics
-- ----------------------------
DROP TABLE IF EXISTS "resource_usage_metrics";
CREATE TABLE "resource_usage_metrics" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "metric_name" varchar(100)  NOT NULL,
  "metric_value" numeric(12,4) NOT NULL,
  "metric_unit" varchar(20)  NOT NULL,
  "recorded_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for scheduled_message_recipients
-- ----------------------------
DROP TABLE IF EXISTS "scheduled_message_recipients";
CREATE TABLE "scheduled_message_recipients" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "scheduled_message_id" TEXT NOT NULL,
  "contact_id" TEXT,
  "group_id" TEXT,
  "phone_number" varchar(30)  NOT NULL,
  "status" varchar(20)  DEFAULT 'pending',
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for scheduled_messages
-- ----------------------------
DROP TABLE IF EXISTS "scheduled_messages";
CREATE TABLE "scheduled_messages" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "user_id" TEXT NOT NULL,
  "device_id" TEXT NOT NULL,
  "template_id" TEXT,
  "group_id" TEXT,
  "recipient_mode" varchar(20)  NOT NULL DEFAULT 'single',
  "recipient_payload" TEXT,
  "message_content" text  NOT NULL,
  "scheduled_at" DATETIME NOT NULL,
  "status" varchar(20)  NOT NULL DEFAULT 'pending',
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "executed_at" DATETIME
)
;

-- ----------------------------
-- Table structure for service_health_checks
-- ----------------------------
DROP TABLE IF EXISTS "service_health_checks";
CREATE TABLE "service_health_checks" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "service_name" varchar(100)  NOT NULL,
  "status" varchar(20)  NOT NULL DEFAULT 'unknown',
  "latency_ms" INTEGER,
  "uptime_percent" numeric(5,2),
  "started_at" DATETIME,
  "last_incident_at" DATETIME,
  "checked_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for settings
-- ----------------------------
DROP TABLE IF EXISTS "settings";
CREATE TABLE "settings" (
  "key" varchar(100)  NOT NULL,
  "value" text  NOT NULL,
  "description" text ,
  "updated_at" DATETIME DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for subscriptions
-- ----------------------------
DROP TABLE IF EXISTS "subscriptions";
CREATE TABLE "subscriptions" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "user_id" TEXT,
  "plan_id" TEXT,
  "status" varchar(20)  NOT NULL DEFAULT 'inactive',
  "created_at" DATETIME,
  "start_date" DATETIME,
  "end_date" DATETIME,
  "renewal_date" DATETIME,
  "auto_renew" BOOLEAN NOT NULL DEFAULT 1,
  "updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "max_devices" INTEGER DEFAULT 0,
  "max_messages_per_day" INTEGER DEFAULT 0
)
;

-- ----------------------------
-- Table structure for system_settings
-- ----------------------------
DROP TABLE IF EXISTS "system_settings";
CREATE TABLE "system_settings" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT  ,
  "keys" varchar(100) ,
  "value" varchar(255) ,
  "description" varchar(255) ,
  "created_at" DATETIME
)
;

-- ----------------------------
-- Table structure for usage_quotas
-- ----------------------------
DROP TABLE IF EXISTS "usage_quotas";
CREATE TABLE "usage_quotas" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "user_id" TEXT NOT NULL,
  "subscription_id" TEXT NOT NULL,
  "period_key" varchar(20)  NOT NULL,
  "messages_used" INTEGER NOT NULL DEFAULT 0,
  "messages_limit" INTEGER,
  "devices_used" INTEGER NOT NULL DEFAULT 0,
  "devices_limit" INTEGER,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for user_sessions
-- ----------------------------
DROP TABLE IF EXISTS "user_sessions";
CREATE TABLE "user_sessions" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "user_id" TEXT NOT NULL,
  "session_token_hash" varchar(255)  NOT NULL,
  "refresh_token_hash" varchar(255) ,
  "ip_address" TEXT,
  "user_agent" text ,
  "expires_at" DATETIME NOT NULL,
  "revoked_at" DATETIME,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "refresh_expires_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "last_active_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS "users";
CREATE TABLE "users" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "phone_number" varchar(30)  NOT NULL,
  "full_name" varchar(50)  NOT NULL,
  "is_verified" BOOLEAN NOT NULL DEFAULT 0,
  "otp_code" varchar(10) ,
  "otp_expired" DATETIME,
  "id_subscribed" TEXT,
  "max_device" INTEGER,
  "is_banned" BOOLEAN NOT NULL DEFAULT 0,
  "is_api_enabled" BOOLEAN NOT NULL DEFAULT 0,
  "email" varchar(255) ,
  "company_name" varchar(255) ,
  "timezone" varchar(100)  NOT NULL DEFAULT 'Asia/Jakarta',
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "last_login_at" DATETIME
)
;

-- ----------------------------
-- Table structure for warming_pool
-- ----------------------------
DROP TABLE IF EXISTS "warming_pool";
CREATE TABLE "warming_pool" (
  "id" TEXT NOT NULL,
  "device_id" TEXT,
  "intensity" INTEGER,
  "daily_limit" INTEGER,
  "message_send_today" INTEGER,
  "is_active" BOOLEAN,
  "next_action_at" DATETIME,
  "created_at" DATETIME
)
;

-- ----------------------------
-- Table structure for warming_sessions
-- ----------------------------
DROP TABLE IF EXISTS "warming_sessions";
CREATE TABLE "warming_sessions" (
  "id" TEXT NOT NULL,
  "device_id" TEXT,
  "target_phone" varchar(30) ,
  "message_sent" varchar(255) ,
  "response_received" varchar(255) ,
  "status" INTEGER,
  "created_at" DATETIME
)
;

-- ----------------------------
-- Table structure for webhook_deliveries
-- ----------------------------
DROP TABLE IF EXISTS "webhook_deliveries";
CREATE TABLE "webhook_deliveries" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "webhook_id" TEXT NOT NULL,
  "event_key" varchar(100)  NOT NULL,
  "payload" TEXT NOT NULL DEFAULT '{}',
  "attempt" INTEGER NOT NULL DEFAULT 1,
  "http_status" INTEGER,
  "response_body" text ,
  "status" varchar(20)  NOT NULL DEFAULT 'pending',
  "sent_at" DATETIME,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for webhook_event_subscriptions
-- ----------------------------
DROP TABLE IF EXISTS "webhook_event_subscriptions";
CREATE TABLE "webhook_event_subscriptions" (
  "id" TEXT NOT NULL DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
  "webhook_id" TEXT NOT NULL,
  "event_key" varchar(100)  NOT NULL,
  "is_enabled" BOOLEAN NOT NULL DEFAULT 1
)
;

-- ----------------------------
-- Table structure for webhook_settings
-- ----------------------------
DROP TABLE IF EXISTS "webhook_settings";
CREATE TABLE "webhook_settings" (
  "user_id" TEXT NOT NULL,
  "url" text  NOT NULL,
  "secret" varchar(255)  NOT NULL,
  "is_active" BOOLEAN DEFAULT 1,
  "enabled_events" TEXT,
  "updated_at" DATETIME DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for webhooks
-- ----------------------------
DROP TABLE IF EXISTS "webhooks";
CREATE TABLE "webhooks" (
  "id" TEXT NOT NULL,
  "device_id" TEXT,
  "webhook_url" varchar(255) ,
  "secret_key" varchar(255) ,
  "created_at" DATETIME
)
;

-- ----------------------------
-- ----------------------------

-- ----------------------------
-- ----------------------------

-- ----------------------------
-- ----------------------------

-- ----------------------------
-- ----------------------------

-- ----------------------------
-- ----------------------------

-- ----------------------------
-- ----------------------------

-- ----------------------------
-- ----------------------------

-- ----------------------------
-- ----------------------------

-- ----------------------------
-- ----------------------------

-- ----------------------------
-- ----------------------------

-- ----------------------------
-- ----------------------------

-- ----------------------------
-- ----------------------------

-- ----------------------------
-- ----------------------------

-- ----------------------------
-- ----------------------------

-- ----------------------------
-- ----------------------------

-- ----------------------------
-- ----------------------------

-- ----------------------------
-- Function structure for armor
-- ----------------------------

-- ----------------------------
-- Function structure for armor
-- ----------------------------

-- ----------------------------
-- Function structure for crypt
-- ----------------------------

-- ----------------------------
-- Function structure for dearmor
-- ----------------------------

-- ----------------------------
-- Function structure for decrypt
-- ----------------------------

-- ----------------------------
-- Function structure for decrypt_iv
-- ----------------------------

-- ----------------------------
-- Function structure for digest
-- ----------------------------

-- ----------------------------
-- Function structure for digest
-- ----------------------------

-- ----------------------------
-- Function structure for encrypt
-- ----------------------------

-- ----------------------------
-- Function structure for encrypt_iv
-- ----------------------------

-- ----------------------------
-- Function structure for fips_mode
-- ----------------------------

-- ----------------------------
-- Function structure for gen_random_bytes
-- ----------------------------

-- ----------------------------
-- Function structure for gen_random_uuid
-- ----------------------------

-- ----------------------------
-- Function structure for gen_salt
-- ----------------------------

-- ----------------------------
-- Function structure for gen_salt
-- ----------------------------

-- ----------------------------
-- Function structure for hmac
-- ----------------------------

-- ----------------------------
-- Function structure for hmac
-- ----------------------------

-- ----------------------------
-- Function structure for pgp_armor_headers
-- ----------------------------

-- ----------------------------
-- Function structure for pgp_key_id
-- ----------------------------

-- ----------------------------
-- Function structure for pgp_pub_decrypt
-- ----------------------------

-- ----------------------------
-- Function structure for pgp_pub_decrypt
-- ----------------------------

-- ----------------------------
-- Function structure for pgp_pub_decrypt
-- ----------------------------

-- ----------------------------
-- Function structure for pgp_pub_decrypt_bytea
-- ----------------------------

-- ----------------------------
-- Function structure for pgp_pub_decrypt_bytea
-- ----------------------------

-- ----------------------------
-- Function structure for pgp_pub_decrypt_bytea
-- ----------------------------

-- ----------------------------
-- Function structure for pgp_pub_encrypt
-- ----------------------------

-- ----------------------------
-- Function structure for pgp_pub_encrypt
-- ----------------------------

-- ----------------------------
-- Function structure for pgp_pub_encrypt_bytea
-- ----------------------------

-- ----------------------------
-- Function structure for pgp_pub_encrypt_bytea
-- ----------------------------

-- ----------------------------
-- Function structure for pgp_sym_decrypt
-- ----------------------------

-- ----------------------------
-- Function structure for pgp_sym_decrypt
-- ----------------------------

-- ----------------------------
-- Function structure for pgp_sym_decrypt_bytea
-- ----------------------------

-- ----------------------------
-- Function structure for pgp_sym_decrypt_bytea
-- ----------------------------

-- ----------------------------
-- Function structure for pgp_sym_encrypt
-- ----------------------------

-- ----------------------------
-- Function structure for pgp_sym_encrypt
-- ----------------------------

-- ----------------------------
-- Function structure for pgp_sym_encrypt_bytea
-- ----------------------------

-- ----------------------------
-- Function structure for pgp_sym_encrypt_bytea
-- ----------------------------

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------

-- ----------------------------
-- Indexes structure for table api_keys
-- ----------------------------
CREATE INDEX "idx_api_keys_hash" ON "api_keys"  (
  "key_hash"   
);

-- ----------------------------
-- Primary Key structure for table api_keys
-- ----------------------------

-- ----------------------------
-- Indexes structure for table api_logs
-- ----------------------------
CREATE INDEX "idx_api_logs_created_at" ON "api_logs"  (
  "created_at"  
);
CREATE INDEX "idx_api_logs_endpoint" ON "api_logs"  (
  "endpoint"   
);
CREATE INDEX "idx_api_logs_user_id" ON "api_logs"  (
  "user_id"  
);

-- ----------------------------
-- Primary Key structure for table api_logs
-- ----------------------------

-- ----------------------------
-- Indexes structure for table audit_logs
-- ----------------------------
CREATE INDEX "idx_audit_logs_resource" ON "audit_logs"  (
  "resource_type"   ,
  "resource_id"  
) WHERE resource_id IS NOT NULL;
CREATE INDEX "idx_audit_logs_user_created" ON "audit_logs"  (
  "user_id"  ,
  "created_at"  
) WHERE user_id IS NOT NULL;

-- ----------------------------
-- Primary Key structure for table audit_logs
-- ----------------------------

-- ----------------------------
-- Indexes structure for table auto_response
-- ----------------------------
CREATE INDEX "idx_auto_response_is_active" ON "auto_response"  (
  "is_active"  
);

-- ----------------------------
-- Primary Key structure for table auto_response
-- ----------------------------

-- ----------------------------
-- Indexes structure for table auto_response_keywords
-- ----------------------------
CREATE INDEX "idx_auto_response_device_id" ON "auto_response_keywords"  (
  "device_id"  
) WHERE device_id IS NOT NULL;
CREATE INDEX "idx_auto_response_user_id" ON "auto_response_keywords"  (
  "user_id"  
);

-- ----------------------------
-- Primary Key structure for table auto_response_keywords
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table auto_response_logs
-- ----------------------------

-- ----------------------------
-- Indexes structure for table billing_plans
-- ----------------------------
CREATE INDEX "idx_billing_plans_active_price" ON "billing_plans"  (
  "is_active"  ,
  "price"  
);

-- ----------------------------
-- Primary Key structure for table billing_plans
-- ----------------------------

-- ----------------------------
-- Indexes structure for table blacklists
-- ----------------------------
CREATE INDEX "idx_blacklists_user_id" ON "blacklists"  (
  "user_id"  
);

-- ----------------------------
-- Uniques structure for table blacklists
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table blacklists
-- ----------------------------

-- ----------------------------
-- Indexes structure for table broadcast_campaigns
-- ----------------------------
CREATE INDEX "idx_broadcast_campaigns_deleted_at" ON "broadcast_campaigns"  (
  "deleted_at"  
);
CREATE INDEX "idx_broadcast_campaigns_device_id" ON "broadcast_campaigns"  (
  "device_id"  
);
CREATE INDEX "idx_broadcast_campaigns_status" ON "broadcast_campaigns"  (
  "status"   
);
CREATE INDEX "idx_broadcast_campaigns_user_id" ON "broadcast_campaigns"  (
  "user_id"  
);

-- ----------------------------
-- Primary Key structure for table broadcast_campaigns
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table broadcast_messages
-- ----------------------------

-- ----------------------------
-- Indexes structure for table broadcast_recipients
-- ----------------------------
CREATE INDEX "idx_broadcast_recipients_campaign_id" ON "broadcast_recipients"  (
  "campaign_id"  
);
CREATE INDEX "idx_broadcast_recipients_status" ON "broadcast_recipients"  (
  "status"   
);

-- ----------------------------
-- Primary Key structure for table broadcast_recipients
-- ----------------------------

-- ----------------------------
-- Indexes structure for table contact
-- ----------------------------
CREATE INDEX "idx_contact_deleted_at" ON "contact"  (
  "deleted_at"  
);
CREATE INDEX "idx_contact_group_id" ON "contact"  (
  "group_id"  
);
CREATE INDEX "idx_contact_phone" ON "contact"  (
  "phone"   
);

-- ----------------------------
-- Primary Key structure for table contact
-- ----------------------------

-- ----------------------------
-- Uniques structure for table contact_group_members
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table contact_group_members
-- ----------------------------

-- ----------------------------
-- Indexes structure for table contact_groups
-- ----------------------------
CREATE INDEX "idx_contact_groups_user_id" ON "contact_groups"  (
  "user_id"  
);

-- ----------------------------
-- Primary Key structure for table contact_groups
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table contact_labels
-- ----------------------------

-- ----------------------------
-- Indexes structure for table contacts
-- ----------------------------
CREATE INDEX "idx_contacts_label_id" ON "contacts"  (
  "label_id"  
) WHERE label_id IS NOT NULL;
CREATE INDEX "idx_contacts_user_name" ON "contacts"  (
  "user_id"  ,
  "name"   
);
CREATE UNIQUE INDEX "idx_contacts_user_phone_active" ON "contacts"  (
  "user_id"  ,
  "phone_number"   
) WHERE deleted_at IS NULL;

-- ----------------------------
-- Primary Key structure for table contacts
-- ----------------------------

-- ----------------------------
-- Indexes structure for table daily_message_stats
-- ----------------------------
CREATE INDEX "idx_daily_stats_user_date" ON "daily_message_stats"  (
  "user_id"  ,
  "stat_date"  
);
CREATE UNIQUE INDEX "idx_daily_stats_user_device" ON "daily_message_stats"  (
  "user_id"  ,
  "device_id"  ,
  "stat_date"  
) WHERE device_id IS NOT NULL;
CREATE UNIQUE INDEX "idx_daily_stats_user_no_device" ON "daily_message_stats"  (
  "user_id"  ,
  "stat_date"  
) WHERE device_id IS NULL;

-- ----------------------------
-- Uniques structure for table daily_message_stats
-- ----------------------------

-- ----------------------------
-- Checks structure for table daily_message_stats
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table daily_message_stats
-- ----------------------------

-- ----------------------------
-- Indexes structure for table device_metrics
-- ----------------------------
CREATE INDEX "idx_device_metrics_device_id" ON "device_metrics"  (
  "device_id"  
);

-- ----------------------------
-- Checks structure for table device_metrics
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table device_metrics
-- ----------------------------

-- ----------------------------
-- Indexes structure for table device_qr_codes
-- ----------------------------
CREATE INDEX "idx_device_qr_device_id" ON "device_qr_codes"  (
  "device_id"  
) WHERE device_id IS NOT NULL;
CREATE INDEX "idx_device_qr_user_id" ON "device_qr_codes"  (
  "user_id"  
);

-- ----------------------------
-- Primary Key structure for table device_qr_codes
-- ----------------------------

-- ----------------------------
-- Indexes structure for table device_sessions
-- ----------------------------
CREATE INDEX "idx_device_sessions_device_id" ON "device_sessions"  (
  "device_id"  
);

-- ----------------------------
-- Checks structure for table device_sessions
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table device_sessions
-- ----------------------------

-- ----------------------------
-- Indexes structure for table devices
-- ----------------------------
CREATE INDEX "idx_devices_last_seen" ON "devices"  (
  "last_seen_at"  
);
CREATE INDEX "idx_devices_phone" ON "devices"  (
  "phone_number"   
);
CREATE INDEX "idx_devices_status" ON "devices"  (
  "status"   
);
CREATE INDEX "idx_devices_user_id" ON "devices"  (
  "user_id"  
);

-- ----------------------------
-- Primary Key structure for table devices
-- ----------------------------

-- ----------------------------
-- Indexes structure for table failure_records
-- ----------------------------
CREATE INDEX "idx_failure_records_device" ON "failure_records"  (
  "device_id"  
) WHERE device_id IS NOT NULL;
CREATE INDEX "idx_failure_records_user" ON "failure_records"  (
  "user_id"  ,
  "occurred_at"  
);

-- ----------------------------
-- Primary Key structure for table failure_records
-- ----------------------------

-- ----------------------------
-- Indexes structure for table groups
-- ----------------------------
CREATE INDEX "idx_groups_deleted_at" ON "groups"  (
  "deleted_at"  
);
CREATE INDEX "idx_groups_user_id" ON "groups"  (
  "user_id"  
);

-- ----------------------------
-- Primary Key structure for table groups
-- ----------------------------

-- ----------------------------
-- Checks structure for table invoice_items
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table invoice_items
-- ----------------------------

-- ----------------------------
-- Indexes structure for table invoices
-- ----------------------------
CREATE INDEX "idx_invoices_status" ON "invoices"  (
  "status"   
);
CREATE INDEX "idx_invoices_user_id" ON "invoices"  (
  "user_id"  
);

-- ----------------------------
-- Uniques structure for table invoices
-- ----------------------------

-- ----------------------------
-- Checks structure for table invoices
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table invoices
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table lookup
-- ----------------------------

-- ----------------------------
-- Indexes structure for table message_templates
-- ----------------------------
CREATE INDEX "idx_message_templates_user_id" ON "message_templates"  (
  "user_id"  
);

-- ----------------------------
-- Checks structure for table message_templates
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table message_templates
-- ----------------------------

-- ----------------------------
-- Indexes structure for table messages
-- ----------------------------
CREATE INDEX "idx_messages_created_at" ON "messages"  (
  "created_at"  
);
CREATE INDEX "idx_messages_device_id" ON "messages"  (
  "device_id"  
);
CREATE INDEX "idx_messages_device_status_dir" ON "messages"  (
  "device_id"  ,
  "status_message"  ,
  "direction"   
);
CREATE INDEX "idx_messages_direction" ON "messages"  (
  "direction"   
);
CREATE INDEX "idx_messages_receipt" ON "messages"  (
  "receipt_number"   
);
CREATE INDEX "idx_messages_scheduled_msg_id" ON "messages"  (
  "scheduled_message_id"  
);
CREATE INDEX "idx_messages_status" ON "messages"  (
  "status_message"  
);
CREATE INDEX "idx_messages_whatsapp_message_id" ON "messages"  (
  "whatsapp_message_id"   
) WHERE whatsapp_message_id IS NOT NULL;

-- ----------------------------
-- Primary Key structure for table messages
-- ----------------------------

-- ----------------------------
-- Uniques structure for table migrations
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table migrations
-- ----------------------------

-- ----------------------------
-- Uniques structure for table notification_settings
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table notification_settings
-- ----------------------------

-- ----------------------------
-- Indexes structure for table notifications
-- ----------------------------
CREATE INDEX "idx_notifications_user_created" ON "notifications"  (
  "user_id"  ,
  "created_at"  
);
CREATE INDEX "idx_notifications_user_unread" ON "notifications"  (
  "user_id"  ,
  "is_read"  
);

-- ----------------------------
-- Primary Key structure for table notifications
-- ----------------------------

-- ----------------------------
-- Uniques structure for table onboarding_progress
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table onboarding_progress
-- ----------------------------

-- ----------------------------
-- Indexes structure for table otp_verifications
-- ----------------------------
CREATE INDEX "idx_otp_expires_at" ON "otp_verifications"  (
  "expires_at"  
);
CREATE INDEX "idx_otp_phone_number" ON "otp_verifications"  (
  "phone_number"   
);
CREATE INDEX "idx_otp_verifications_expires_at" ON "otp_verifications"  (
  "expires_at"  
);
CREATE INDEX "idx_otp_verifications_phone_number" ON "otp_verifications"  (
  "phone_number"   
);
CREATE INDEX "idx_otp_verifications_user_id" ON "otp_verifications"  (
  "user_id"  
);

-- ----------------------------
-- Checks structure for table otp_verifications
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table otp_verifications
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table resource_usage_metrics
-- ----------------------------

-- ----------------------------
-- Uniques structure for table scheduled_message_recipients
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table scheduled_message_recipients
-- ----------------------------

-- ----------------------------
-- Indexes structure for table scheduled_messages
-- ----------------------------
CREATE INDEX "idx_scheduled_messages_status" ON "scheduled_messages"  (
  "status"   
);
CREATE INDEX "idx_scheduled_messages_user" ON "scheduled_messages"  (
  "user_id"  ,
  "scheduled_at"  
);

-- ----------------------------
-- Primary Key structure for table scheduled_messages
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table service_health_checks
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table settings
-- ----------------------------

-- ----------------------------
-- Indexes structure for table subscriptions
-- ----------------------------
CREATE INDEX "idx_subscriptions_plan_id" ON "subscriptions"  (
  "plan_id"  
);
CREATE INDEX "idx_subscriptions_status" ON "subscriptions"  (
  "status"   
);
CREATE INDEX "idx_subscriptions_user_id" ON "subscriptions"  (
  "user_id"  
);
CREATE INDEX "idx_subscriptions_user_status" ON "subscriptions"  (
  "user_id"  ,
  "status"   
);

-- ----------------------------
-- Primary Key structure for table subscriptions
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table system_settings
-- ----------------------------

-- ----------------------------
-- Indexes structure for table usage_quotas
-- ----------------------------
CREATE INDEX "idx_usage_quotas_user_period" ON "usage_quotas"  (
  "user_id"  ,
  "period_key"   
);

-- ----------------------------
-- Uniques structure for table usage_quotas
-- ----------------------------

-- ----------------------------
-- Checks structure for table usage_quotas
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table usage_quotas
-- ----------------------------

-- ----------------------------
-- Indexes structure for table user_sessions
-- ----------------------------
CREATE INDEX "idx_user_sessions_expires_at" ON "user_sessions"  (
  "expires_at"  
);
CREATE INDEX "idx_user_sessions_refresh_token_hash" ON "user_sessions"  (
  "refresh_token_hash"   
);
CREATE UNIQUE INDEX "idx_user_sessions_session_token_hash" ON "user_sessions"  (
  "session_token_hash"   
);
CREATE INDEX "idx_user_sessions_user_id" ON "user_sessions"  (
  "user_id"  
);

-- ----------------------------
-- Primary Key structure for table user_sessions
-- ----------------------------

-- ----------------------------
-- Indexes structure for table users
-- ----------------------------
CREATE INDEX "idx_users_is_ban" ON "users"  (
  "is_banned"  
);
CREATE INDEX "idx_users_is_verify" ON "users"  (
  "is_verified"  
);
CREATE INDEX "idx_users_phone" ON "users"  (
  "phone_number"   
);
CREATE UNIQUE INDEX "idx_users_phone_number" ON "users"  (
  "phone_number"   
);

-- ----------------------------
-- Primary Key structure for table users
-- ----------------------------

-- ----------------------------
-- Indexes structure for table warming_pool
-- ----------------------------
CREATE INDEX "idx_warming_pool_device_id" ON "warming_pool"  (
  "device_id"  
);
CREATE INDEX "idx_warming_pool_is_active" ON "warming_pool"  (
  "is_active"  
);
CREATE INDEX "idx_warming_pool_next_action" ON "warming_pool"  (
  "next_action_at"  
);

-- ----------------------------
-- Primary Key structure for table warming_pool
-- ----------------------------

-- ----------------------------
-- Indexes structure for table warming_sessions
-- ----------------------------
CREATE INDEX "idx_warming_sessions_device_id" ON "warming_sessions"  (
  "device_id"  
);
CREATE INDEX "idx_warming_sessions_status" ON "warming_sessions"  (
  "status"  
);

-- ----------------------------
-- Primary Key structure for table warming_sessions
-- ----------------------------

-- ----------------------------
-- Indexes structure for table webhook_deliveries
-- ----------------------------
CREATE INDEX "idx_webhook_deliveries_webhook" ON "webhook_deliveries"  (
  "webhook_id"  ,
  "created_at"  
);

-- ----------------------------
-- Checks structure for table webhook_deliveries
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table webhook_deliveries
-- ----------------------------

-- ----------------------------
-- Uniques structure for table webhook_event_subscriptions
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table webhook_event_subscriptions
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table webhook_settings
-- ----------------------------

-- ----------------------------
-- Indexes structure for table webhooks
-- ----------------------------
CREATE INDEX "idx_webhooks_device_id" ON "webhooks"  (
  "device_id"  
);

-- ----------------------------
-- Primary Key structure for table webhooks
-- ----------------------------

-- ----------------------------
-- ----------------------------

-- ----------------------------
-- ----------------------------

-- ----------------------------
-- ----------------------------

-- ----------------------------
-- Foreign Keys structure for table auto_response
-- ----------------------------

-- ----------------------------
-- Foreign Keys structure for table auto_response_logs
-- ----------------------------

-- ----------------------------
-- Foreign Keys structure for table broadcast_campaigns
-- ----------------------------

-- ----------------------------
-- Foreign Keys structure for table broadcast_messages
-- ----------------------------

-- ----------------------------
-- Foreign Keys structure for table broadcast_recipients
-- ----------------------------

-- ----------------------------
-- Foreign Keys structure for table contact
-- ----------------------------

-- ----------------------------
-- Foreign Keys structure for table contact_group_members
-- ----------------------------

-- ----------------------------
-- Foreign Keys structure for table contacts
-- ----------------------------

-- ----------------------------
-- Foreign Keys structure for table devices
-- ----------------------------

-- ----------------------------
-- Foreign Keys structure for table groups
-- ----------------------------

-- ----------------------------
-- Foreign Keys structure for table invoice_items
-- ----------------------------

-- ----------------------------
-- Foreign Keys structure for table messages
-- ----------------------------

-- ----------------------------
-- Foreign Keys structure for table scheduled_message_recipients
-- ----------------------------

-- ----------------------------
-- Foreign Keys structure for table scheduled_messages
-- ----------------------------

-- ----------------------------
-- Foreign Keys structure for table subscriptions
-- ----------------------------

-- ----------------------------
-- Foreign Keys structure for table warming_pool
-- ----------------------------

-- ----------------------------
-- Foreign Keys structure for table warming_sessions
-- ----------------------------

-- ----------------------------
-- Foreign Keys structure for table webhooks
-- ----------------------------

-- ----------------------------
-- ----------------------------