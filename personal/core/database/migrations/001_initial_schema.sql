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
DROP TYPE IF EXISTS "public"."enum_direction";
CREATE TYPE "public"."enum_direction" AS ENUM (
  'IN',
  'OUT'
);

-- ----------------------------
-- Sequence structure for lookup_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."lookup_id_seq";
CREATE SEQUENCE "public"."lookup_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for migrations_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."migrations_id_seq";
CREATE SEQUENCE "public"."migrations_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for system_settings_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."system_settings_id_seq";
CREATE SEQUENCE "public"."system_settings_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;

-- ----------------------------
-- Table structure for api_keys
-- ----------------------------
DROP TABLE IF EXISTS "public"."api_keys";
CREATE TABLE "public"."api_keys" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "prefix" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "key_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "last_used_at" timestamptz(6),
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "deleted_at" timestamptz(6)
)
;

-- ----------------------------
-- Table structure for api_logs
-- ----------------------------
DROP TABLE IF EXISTS "public"."api_logs";
CREATE TABLE "public"."api_logs" (
  "id" uuid NOT NULL,
  "user_id" uuid,
  "endpoint" varchar(255) COLLATE "pg_catalog"."default",
  "req_body" jsonb,
  "response_body" jsonb,
  "created_at" timestamptz(6),
  "ip_address" varchar(255) COLLATE "pg_catalog"."default",
  "device_id" uuid
)
;

-- ----------------------------
-- Table structure for audit_logs
-- ----------------------------
DROP TABLE IF EXISTS "public"."audit_logs";
CREATE TABLE "public"."audit_logs" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid,
  "action_type" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "resource_type" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "resource_id" uuid,
  "metadata" jsonb,
  "ip_address" inet,
  "user_agent" text COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for auto_response
-- ----------------------------
DROP TABLE IF EXISTS "public"."auto_response";
CREATE TABLE "public"."auto_response" (
  "id" uuid NOT NULL,
  "device_id" uuid,
  "keyword" varchar(255) COLLATE "pg_catalog"."default",
  "response_text" varchar(255) COLLATE "pg_catalog"."default",
  "is_active" bool,
  "created_at" timestamptz(6)
)
;

-- ----------------------------
-- Table structure for auto_response_keywords
-- ----------------------------
DROP TABLE IF EXISTS "public"."auto_response_keywords";
CREATE TABLE "public"."auto_response_keywords" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "device_id" uuid,
  "keyword" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "response_text" text COLLATE "pg_catalog"."default" NOT NULL,
  "is_active" bool NOT NULL DEFAULT true,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
  "match_type" varchar(20) COLLATE "pg_catalog"."default" DEFAULT 'contains'::character varying
)
;

-- ----------------------------
-- Table structure for auto_response_logs
-- ----------------------------
DROP TABLE IF EXISTS "public"."auto_response_logs";
CREATE TABLE "public"."auto_response_logs" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "keyword_id" uuid NOT NULL,
  "message_id" uuid,
  "triggered_by_phone" varchar(30) COLLATE "pg_catalog"."default" NOT NULL,
  "matched_keyword" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "response_sent" text COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for billing_plans
-- ----------------------------
DROP TABLE IF EXISTS "public"."billing_plans";
CREATE TABLE "public"."billing_plans" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "price" numeric(10,2) NOT NULL,
  "max_devices" int4 NOT NULL,
  "max_messages_per_day" int4 NOT NULL,
  "features" jsonb,
  "billing_cycle" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'monthly'::character varying,
  "is_active" bool NOT NULL DEFAULT true,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for blacklists
-- ----------------------------
DROP TABLE IF EXISTS "public"."blacklists";
CREATE TABLE "public"."blacklists" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "phone_number" varchar(30) COLLATE "pg_catalog"."default" NOT NULL,
  "reason" varchar(255) COLLATE "pg_catalog"."default",
  "blocked_at" timestamptz(6) NOT NULL DEFAULT now(),
  "unblocked_at" timestamptz(6),
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for broadcast_campaigns
-- ----------------------------
DROP TABLE IF EXISTS "public"."broadcast_campaigns";
CREATE TABLE "public"."broadcast_campaigns" (
  "id" uuid NOT NULL,
  "user_id" uuid,
  "device_id" uuid,
  "name_broadcast" varchar(255) COLLATE "pg_catalog"."default",
  "total_recipients" int4,
  "processed_count" int4,
  "scheduled_at" timestamptz(6),
  "status" varchar(20) COLLATE "pg_catalog"."default" DEFAULT 'draft'::character varying,
  "created_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "template_id" uuid,
  "name" varchar(255) COLLATE "pg_catalog"."default",
  "message_content" text COLLATE "pg_catalog"."default",
  "delay_seconds" int4 NOT NULL DEFAULT 5,
  "success_count" int4 NOT NULL DEFAULT 0,
  "failed_count" int4 NOT NULL DEFAULT 0,
  "started_at" timestamptz(6),
  "completed_at" timestamptz(6),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for broadcast_messages
-- ----------------------------
DROP TABLE IF EXISTS "public"."broadcast_messages";
CREATE TABLE "public"."broadcast_messages" (
  "id" uuid NOT NULL,
  "campaign_id" uuid,
  "message_type" int4,
  "message_text" varchar(255) COLLATE "pg_catalog"."default",
  "media_url" varchar(255) COLLATE "pg_catalog"."default",
  "button_data" jsonb
)
;

-- ----------------------------
-- Table structure for broadcast_recipients
-- ----------------------------
DROP TABLE IF EXISTS "public"."broadcast_recipients";
CREATE TABLE "public"."broadcast_recipients" (
  "id" uuid NOT NULL,
  "campaign_id" uuid,
  "groups_id" uuid,
  "contact_id" uuid,
  "status" varchar(20) COLLATE "pg_catalog"."default" DEFAULT 'pending'::character varying,
  "sent_at" timestamptz(6),
  "error_messages" varchar(255) COLLATE "pg_catalog"."default",
  "retry_count" int4,
  "created_at" timestamptz(6),
  "group_id" uuid,
  "phone_number" varchar(30) COLLATE "pg_catalog"."default",
  "failed_at" timestamptz(6),
  "error_message" text COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for contact
-- ----------------------------
DROP TABLE IF EXISTS "public"."contact";
CREATE TABLE "public"."contact" (
  "id" uuid NOT NULL,
  "group_id" uuid,
  "name" varchar(100) COLLATE "pg_catalog"."default",
  "phone" varchar(30) COLLATE "pg_catalog"."default",
  "additional_data" jsonb,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6)
)
;

-- ----------------------------
-- Table structure for contact_group_members
-- ----------------------------
DROP TABLE IF EXISTS "public"."contact_group_members";
CREATE TABLE "public"."contact_group_members" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "group_id" uuid NOT NULL,
  "contact_id" uuid NOT NULL,
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for contact_groups
-- ----------------------------
DROP TABLE IF EXISTS "public"."contact_groups";
CREATE TABLE "public"."contact_groups" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
  "deleted_at" timestamptz(6)
)
;

-- ----------------------------
-- Table structure for contact_labels
-- ----------------------------
DROP TABLE IF EXISTS "public"."contact_labels";
CREATE TABLE "public"."contact_labels" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "name" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "color" varchar(20) COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for contacts
-- ----------------------------
DROP TABLE IF EXISTS "public"."contacts";
CREATE TABLE "public"."contacts" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "label_id" uuid,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "phone_number" varchar(30) COLLATE "pg_catalog"."default" NOT NULL,
  "additional_data" jsonb,
  "note" text COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
  "deleted_at" timestamptz(6)
)
;

-- ----------------------------
-- Table structure for daily_message_stats
-- ----------------------------
DROP TABLE IF EXISTS "public"."daily_message_stats";
CREATE TABLE "public"."daily_message_stats" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "device_id" uuid,
  "stat_date" date NOT NULL,
  "sent_count" int4 NOT NULL DEFAULT 0,
  "failed_count" int4 NOT NULL DEFAULT 0,
  "delivered_count" int4 NOT NULL DEFAULT 0,
  "received_count" int4 NOT NULL DEFAULT 0,
  "success_rate" numeric(5,2),
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for device_metrics
-- ----------------------------
DROP TABLE IF EXISTS "public"."device_metrics";
CREATE TABLE "public"."device_metrics" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "device_id" uuid NOT NULL,
  "uptime_seconds" int8 NOT NULL DEFAULT 0,
  "messages_sent_count" int4 NOT NULL DEFAULT 0,
  "messages_received_count" int4 NOT NULL DEFAULT 0,
  "success_rate" numeric(5,2),
  "recorded_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for device_qr_codes
-- ----------------------------
DROP TABLE IF EXISTS "public"."device_qr_codes";
CREATE TABLE "public"."device_qr_codes" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "device_id" uuid,
  "qr_string" text COLLATE "pg_catalog"."default",
  "qr_image_url" text COLLATE "pg_catalog"."default",
  "status" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'pending'::character varying,
  "generated_at" timestamptz(6) NOT NULL DEFAULT now(),
  "expired_at" timestamptz(6)
)
;

-- ----------------------------
-- Table structure for device_sessions
-- ----------------------------
DROP TABLE IF EXISTS "public"."device_sessions";
CREATE TABLE "public"."device_sessions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "device_id" uuid NOT NULL,
  "session_blob" bytea,
  "session_status" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'inactive'::character varying,
  "started_at" timestamptz(6) NOT NULL DEFAULT now(),
  "ended_at" timestamptz(6),
  "restart_count" int4 NOT NULL DEFAULT 0,
  "last_restart_at" timestamptz(6),
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for devices
-- ----------------------------
DROP TABLE IF EXISTS "public"."devices";
CREATE TABLE "public"."devices" (
  "id" uuid NOT NULL,
  "user_id" uuid,
  "unique_name" varchar(100) COLLATE "pg_catalog"."default",
  "display_name" varchar(100) COLLATE "pg_catalog"."default",
  "phone_number" varchar(30) COLLATE "pg_catalog"."default",
  "status" varchar(20) COLLATE "pg_catalog"."default",
  "last_seen_at" timestamptz(6),
  "session_data" bytea,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "connected_since" timestamptz(6),
  "platform" varchar(255) COLLATE "pg_catalog"."default",
  "wa_version" varchar(50) COLLATE "pg_catalog"."default",
  "battery_level" int4
)
;

-- ----------------------------
-- Table structure for failure_records
-- ----------------------------
DROP TABLE IF EXISTS "public"."failure_records";
CREATE TABLE "public"."failure_records" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "device_id" uuid,
  "message_id" uuid,
  "recipient_phone" varchar(30) COLLATE "pg_catalog"."default" NOT NULL,
  "failure_type" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "failure_reason" text COLLATE "pg_catalog"."default",
  "occurred_at" timestamptz(6) NOT NULL DEFAULT now(),
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for groups
-- ----------------------------
DROP TABLE IF EXISTS "public"."groups";
CREATE TABLE "public"."groups" (
  "id" uuid NOT NULL,
  "user_id" uuid,
  "group_name" varchar(100) COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6),
  "deleted_at" timestamptz(6)
)
;

-- ----------------------------
-- Table structure for invoice_items
-- ----------------------------
DROP TABLE IF EXISTS "public"."invoice_items";
CREATE TABLE "public"."invoice_items" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "invoice_id" uuid NOT NULL,
  "description" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "qty" int4 NOT NULL DEFAULT 1,
  "unit_price" numeric(12,2) NOT NULL,
  "total_price" numeric(12,2) NOT NULL
)
;

-- ----------------------------
-- Table structure for invoices
-- ----------------------------
DROP TABLE IF EXISTS "public"."invoices";
CREATE TABLE "public"."invoices" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "subscription_id" uuid,
  "invoice_number" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "issue_date" date NOT NULL DEFAULT CURRENT_DATE,
  "due_date" date,
  "paid_at" timestamptz(6),
  "amount" numeric(12,2) NOT NULL,
  "currency" varchar(10) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'IDR'::character varying,
  "status" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'unpaid'::character varying,
  "payment_method" varchar(50) COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for lookup
-- ----------------------------
DROP TABLE IF EXISTS "public"."lookup";
CREATE TABLE "public"."lookup" (
  "id" int4 NOT NULL DEFAULT nextval('lookup_id_seq'::regclass),
  "keys" varchar(255) COLLATE "pg_catalog"."default",
  "values" varchar(100) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for message_templates
-- ----------------------------
DROP TABLE IF EXISTS "public"."message_templates";
CREATE TABLE "public"."message_templates" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "category" varchar(50) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'general'::character varying,
  "content" text COLLATE "pg_catalog"."default" NOT NULL,
  "used_count" int4 NOT NULL DEFAULT 0,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for messages
-- ----------------------------
DROP TABLE IF EXISTS "public"."messages";
CREATE TABLE "public"."messages" (
  "id" uuid NOT NULL,
  "device_id" uuid,
  "direction" varchar(10) COLLATE "pg_catalog"."default",
  "receipt_number" text COLLATE "pg_catalog"."default",
  "message_type" int4,
  "content" text COLLATE "pg_catalog"."default",
  "status_message" int4,
  "error_log" text COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6),
  "target_jid" varchar(100) COLLATE "pg_catalog"."default",
  "updated_at" timestamptz(6) DEFAULT now(),
  "priority" int4 DEFAULT 3,
  "retry_count" int4 DEFAULT 0,
  "max_retries" int4 DEFAULT 3,
  "scheduled_for" timestamptz(6),
  "whatsapp_message_id" varchar(255) COLLATE "pg_catalog"."default",
  "media_url" text COLLATE "pg_catalog"."default",
  "caption" text COLLATE "pg_catalog"."default",
  "broadcast_id" uuid,
  "scheduled_message_id" uuid
)
;

-- ----------------------------
-- Table structure for migrations
-- ----------------------------
DROP TABLE IF EXISTS "public"."migrations";
CREATE TABLE "public"."migrations" (
  "id" int4 NOT NULL DEFAULT nextval('migrations_id_seq'::regclass),
  "version" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "applied_at" timestamptz(6) NOT NULL
)
;

-- ----------------------------
-- Table structure for notification_settings
-- ----------------------------
DROP TABLE IF EXISTS "public"."notification_settings";
CREATE TABLE "public"."notification_settings" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "event_key" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "email_enabled" bool NOT NULL DEFAULT false,
  "in_app_enabled" bool NOT NULL DEFAULT true,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for notifications
-- ----------------------------
DROP TABLE IF EXISTS "public"."notifications";
CREATE TABLE "public"."notifications" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "type" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'info'::character varying,
  "title" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "body" text COLLATE "pg_catalog"."default" NOT NULL,
  "is_read" bool NOT NULL DEFAULT false,
  "read_at" timestamptz(6),
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for onboarding_progress
-- ----------------------------
DROP TABLE IF EXISTS "public"."onboarding_progress";
CREATE TABLE "public"."onboarding_progress" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "step_key" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "is_completed" bool NOT NULL DEFAULT false,
  "completed_at" timestamptz(6),
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for otp_verifications
-- ----------------------------
DROP TABLE IF EXISTS "public"."otp_verifications";
CREATE TABLE "public"."otp_verifications" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid,
  "phone_number" varchar(30) COLLATE "pg_catalog"."default" NOT NULL,
  "context" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "otp_code" varchar(10) COLLATE "pg_catalog"."default" NOT NULL,
  "attempt_count" int4 NOT NULL DEFAULT 0,
  "expires_at" timestamptz(6) NOT NULL,
  "verified_at" timestamptz(6),
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for resource_usage_metrics
-- ----------------------------
DROP TABLE IF EXISTS "public"."resource_usage_metrics";
CREATE TABLE "public"."resource_usage_metrics" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "metric_name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "metric_value" numeric(12,4) NOT NULL,
  "metric_unit" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "recorded_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for scheduled_message_recipients
-- ----------------------------
DROP TABLE IF EXISTS "public"."scheduled_message_recipients";
CREATE TABLE "public"."scheduled_message_recipients" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "scheduled_message_id" uuid NOT NULL,
  "contact_id" uuid,
  "group_id" uuid,
  "phone_number" varchar(30) COLLATE "pg_catalog"."default" NOT NULL,
  "status" varchar(20) COLLATE "pg_catalog"."default" DEFAULT 'pending'::character varying,
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for scheduled_messages
-- ----------------------------
DROP TABLE IF EXISTS "public"."scheduled_messages";
CREATE TABLE "public"."scheduled_messages" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "device_id" uuid NOT NULL,
  "template_id" uuid,
  "group_id" uuid,
  "recipient_mode" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'single'::character varying,
  "recipient_payload" jsonb,
  "message_content" text COLLATE "pg_catalog"."default" NOT NULL,
  "scheduled_at" timestamptz(6) NOT NULL,
  "status" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'pending'::character varying,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
  "executed_at" timestamptz(6)
)
;

-- ----------------------------
-- Table structure for service_health_checks
-- ----------------------------
DROP TABLE IF EXISTS "public"."service_health_checks";
CREATE TABLE "public"."service_health_checks" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "service_name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "status" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'unknown'::character varying,
  "latency_ms" int4,
  "uptime_percent" numeric(5,2),
  "started_at" timestamptz(6),
  "last_incident_at" timestamptz(6),
  "checked_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for settings
-- ----------------------------
DROP TABLE IF EXISTS "public"."settings";
CREATE TABLE "public"."settings" (
  "key" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "value" text COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default",
  "updated_at" timestamptz(6) DEFAULT now(),
  PRIMARY KEY ("key")
)
;

-- ----------------------------
-- Table structure for subscriptions
-- ----------------------------
DROP TABLE IF EXISTS "public"."subscriptions";
CREATE TABLE "public"."subscriptions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid,
  "plan_id" uuid,
  "status" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'inactive'::character varying,
  "created_at" timestamptz(6),
  "start_date" timestamptz(6),
  "end_date" timestamptz(6),
  "renewal_date" timestamptz(6),
  "auto_renew" bool NOT NULL DEFAULT true,
  "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
  "max_devices" int4 DEFAULT 0,
  "max_messages_per_day" int4 DEFAULT 0
)
;

-- ----------------------------
-- Table structure for system_settings
-- ----------------------------
DROP TABLE IF EXISTS "public"."system_settings";
CREATE TABLE "public"."system_settings" (
  "id" int4 NOT NULL DEFAULT nextval('system_settings_id_seq'::regclass),
  "keys" varchar(100) COLLATE "pg_catalog"."default",
  "value" varchar(255) COLLATE "pg_catalog"."default",
  "description" varchar(255) COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6)
)
;

-- ----------------------------
-- Table structure for usage_quotas
-- ----------------------------
DROP TABLE IF EXISTS "public"."usage_quotas";
CREATE TABLE "public"."usage_quotas" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "subscription_id" uuid NOT NULL,
  "period_key" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "messages_used" int4 NOT NULL DEFAULT 0,
  "messages_limit" int4,
  "devices_used" int4 NOT NULL DEFAULT 0,
  "devices_limit" int4,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for user_sessions
-- ----------------------------
DROP TABLE IF EXISTS "public"."user_sessions";
CREATE TABLE "public"."user_sessions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "session_token_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "refresh_token_hash" varchar(255) COLLATE "pg_catalog"."default",
  "ip_address" inet,
  "user_agent" text COLLATE "pg_catalog"."default",
  "expires_at" timestamptz(6) NOT NULL,
  "revoked_at" timestamptz(6),
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "refresh_expires_at" timestamptz(6) NOT NULL DEFAULT now(),
  "last_active_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS "public"."users";
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "phone_number" varchar(30) COLLATE "pg_catalog"."default" NOT NULL,
  "full_name" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "is_verified" bool NOT NULL DEFAULT false,
  "otp_code" varchar(10) COLLATE "pg_catalog"."default",
  "otp_expired" timestamptz(6),
  "id_subscribed" uuid,
  "max_device" int4,
  "is_banned" bool NOT NULL DEFAULT false,
  "is_api_enabled" bool NOT NULL DEFAULT false,
  "email" varchar(255) COLLATE "pg_catalog"."default",
  "company_name" varchar(255) COLLATE "pg_catalog"."default",
  "timezone" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'Asia/Jakarta'::character varying,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
  "last_login_at" timestamptz(6)
)
;

-- ----------------------------
-- Table structure for warming_pool
-- ----------------------------
DROP TABLE IF EXISTS "public"."warming_pool";
CREATE TABLE "public"."warming_pool" (
  "id" uuid NOT NULL,
  "device_id" uuid,
  "intensity" int4,
  "daily_limit" int4,
  "message_send_today" int4,
  "is_active" bool,
  "next_action_at" timestamptz(6),
  "created_at" timestamptz(6)
)
;

-- ----------------------------
-- Table structure for warming_sessions
-- ----------------------------
DROP TABLE IF EXISTS "public"."warming_sessions";
CREATE TABLE "public"."warming_sessions" (
  "id" uuid NOT NULL,
  "device_id" uuid,
  "target_phone" varchar(30) COLLATE "pg_catalog"."default",
  "message_sent" varchar(255) COLLATE "pg_catalog"."default",
  "response_received" varchar(255) COLLATE "pg_catalog"."default",
  "status" int4,
  "created_at" timestamptz(6)
)
;

-- ----------------------------
-- Table structure for webhook_deliveries
-- ----------------------------
DROP TABLE IF EXISTS "public"."webhook_deliveries";
CREATE TABLE "public"."webhook_deliveries" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "webhook_id" uuid NOT NULL,
  "event_key" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "payload" jsonb NOT NULL DEFAULT '{}'::jsonb,
  "attempt" int4 NOT NULL DEFAULT 1,
  "http_status" int4,
  "response_body" text COLLATE "pg_catalog"."default",
  "status" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'pending'::character varying,
  "sent_at" timestamptz(6),
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Table structure for webhook_event_subscriptions
-- ----------------------------
DROP TABLE IF EXISTS "public"."webhook_event_subscriptions";
CREATE TABLE "public"."webhook_event_subscriptions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "webhook_id" uuid NOT NULL,
  "event_key" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "is_enabled" bool NOT NULL DEFAULT true
)
;

-- ----------------------------
-- Table structure for webhook_settings
-- ----------------------------
DROP TABLE IF EXISTS "public"."webhook_settings";
CREATE TABLE "public"."webhook_settings" (
  "user_id" uuid NOT NULL,
  "url" text COLLATE "pg_catalog"."default" NOT NULL,
  "secret" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "is_active" bool DEFAULT true,
  "enabled_events" jsonb,
  "updated_at" timestamptz(6) DEFAULT CURRENT_TIMESTAMP
)
;

-- ----------------------------
-- Table structure for webhooks
-- ----------------------------
DROP TABLE IF EXISTS "public"."webhooks";
CREATE TABLE "public"."webhooks" (
  "id" uuid NOT NULL,
  "device_id" uuid,
  "webhook_url" varchar(255) COLLATE "pg_catalog"."default",
  "secret_key" varchar(255) COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6)
)
;

-- ----------------------------
-- Table structure for whatsmeow_app_state_mutation_macs
-- ----------------------------
DROP TABLE IF EXISTS "public"."whatsmeow_app_state_mutation_macs";
CREATE TABLE "public"."whatsmeow_app_state_mutation_macs" (
  "jid" text COLLATE "pg_catalog"."default" NOT NULL,
  "name" text COLLATE "pg_catalog"."default" NOT NULL,
  "version" int8 NOT NULL,
  "index_mac" bytea NOT NULL,
  "value_mac" bytea NOT NULL
)
;

-- ----------------------------
-- Table structure for whatsmeow_app_state_sync_keys
-- ----------------------------
DROP TABLE IF EXISTS "public"."whatsmeow_app_state_sync_keys";
CREATE TABLE "public"."whatsmeow_app_state_sync_keys" (
  "jid" text COLLATE "pg_catalog"."default" NOT NULL,
  "key_id" bytea NOT NULL,
  "key_data" bytea NOT NULL,
  "timestamp" int8 NOT NULL,
  "fingerprint" bytea NOT NULL
)
;

-- ----------------------------
-- Table structure for whatsmeow_app_state_version
-- ----------------------------
DROP TABLE IF EXISTS "public"."whatsmeow_app_state_version";
CREATE TABLE "public"."whatsmeow_app_state_version" (
  "jid" text COLLATE "pg_catalog"."default" NOT NULL,
  "name" text COLLATE "pg_catalog"."default" NOT NULL,
  "version" int8 NOT NULL,
  "hash" bytea NOT NULL
)
;

-- ----------------------------
-- Table structure for whatsmeow_chat_settings
-- ----------------------------
DROP TABLE IF EXISTS "public"."whatsmeow_chat_settings";
CREATE TABLE "public"."whatsmeow_chat_settings" (
  "our_jid" text COLLATE "pg_catalog"."default" NOT NULL,
  "chat_jid" text COLLATE "pg_catalog"."default" NOT NULL,
  "muted_until" int8 NOT NULL DEFAULT 0,
  "pinned" bool NOT NULL DEFAULT false,
  "archived" bool NOT NULL DEFAULT false
)
;

-- ----------------------------
-- Table structure for whatsmeow_contacts
-- ----------------------------
DROP TABLE IF EXISTS "public"."whatsmeow_contacts";
CREATE TABLE "public"."whatsmeow_contacts" (
  "our_jid" text COLLATE "pg_catalog"."default" NOT NULL,
  "their_jid" text COLLATE "pg_catalog"."default" NOT NULL,
  "first_name" text COLLATE "pg_catalog"."default",
  "full_name" text COLLATE "pg_catalog"."default",
  "push_name" text COLLATE "pg_catalog"."default",
  "business_name" text COLLATE "pg_catalog"."default",
  "redacted_phone" text COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for whatsmeow_device
-- ----------------------------
DROP TABLE IF EXISTS "public"."whatsmeow_device";
CREATE TABLE "public"."whatsmeow_device" (
  "jid" text COLLATE "pg_catalog"."default" NOT NULL,
  "lid" text COLLATE "pg_catalog"."default",
  "facebook_uuid" uuid,
  "registration_id" int8 NOT NULL,
  "noise_key" bytea NOT NULL,
  "identity_key" bytea NOT NULL,
  "signed_pre_key" bytea NOT NULL,
  "signed_pre_key_id" int4 NOT NULL,
  "signed_pre_key_sig" bytea NOT NULL,
  "adv_key" bytea NOT NULL,
  "adv_details" bytea NOT NULL,
  "adv_account_sig" bytea NOT NULL,
  "adv_account_sig_key" bytea NOT NULL,
  "adv_device_sig" bytea NOT NULL,
  "platform" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
  "business_name" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
  "push_name" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
  "lid_migration_ts" int8 NOT NULL DEFAULT 0
)
;

-- ----------------------------
-- Table structure for whatsmeow_event_buffer
-- ----------------------------
DROP TABLE IF EXISTS "public"."whatsmeow_event_buffer";
CREATE TABLE "public"."whatsmeow_event_buffer" (
  "our_jid" text COLLATE "pg_catalog"."default" NOT NULL,
  "ciphertext_hash" bytea NOT NULL,
  "plaintext" bytea,
  "server_timestamp" int8 NOT NULL,
  "insert_timestamp" int8 NOT NULL
)
;

-- ----------------------------
-- Table structure for whatsmeow_identity_keys
-- ----------------------------
DROP TABLE IF EXISTS "public"."whatsmeow_identity_keys";
CREATE TABLE "public"."whatsmeow_identity_keys" (
  "our_jid" text COLLATE "pg_catalog"."default" NOT NULL,
  "their_id" text COLLATE "pg_catalog"."default" NOT NULL,
  "identity" bytea NOT NULL
)
;

-- ----------------------------
-- Table structure for whatsmeow_lid_map
-- ----------------------------
DROP TABLE IF EXISTS "public"."whatsmeow_lid_map";
CREATE TABLE "public"."whatsmeow_lid_map" (
  "lid" text COLLATE "pg_catalog"."default" NOT NULL,
  "pn" text COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Table structure for whatsmeow_message_secrets
-- ----------------------------
DROP TABLE IF EXISTS "public"."whatsmeow_message_secrets";
CREATE TABLE "public"."whatsmeow_message_secrets" (
  "our_jid" text COLLATE "pg_catalog"."default" NOT NULL,
  "chat_jid" text COLLATE "pg_catalog"."default" NOT NULL,
  "sender_jid" text COLLATE "pg_catalog"."default" NOT NULL,
  "message_id" text COLLATE "pg_catalog"."default" NOT NULL,
  "key" bytea NOT NULL
)
;

-- ----------------------------
-- Table structure for whatsmeow_pre_keys
-- ----------------------------
DROP TABLE IF EXISTS "public"."whatsmeow_pre_keys";
CREATE TABLE "public"."whatsmeow_pre_keys" (
  "jid" text COLLATE "pg_catalog"."default" NOT NULL,
  "key_id" int4 NOT NULL,
  "key" bytea NOT NULL,
  "uploaded" bool NOT NULL
)
;

-- ----------------------------
-- Table structure for whatsmeow_privacy_tokens
-- ----------------------------
DROP TABLE IF EXISTS "public"."whatsmeow_privacy_tokens";
CREATE TABLE "public"."whatsmeow_privacy_tokens" (
  "our_jid" text COLLATE "pg_catalog"."default" NOT NULL,
  "their_jid" text COLLATE "pg_catalog"."default" NOT NULL,
  "token" bytea NOT NULL,
  "timestamp" int8 NOT NULL,
  "sender_timestamp" int8
)
;

-- ----------------------------
-- Table structure for whatsmeow_retry_buffer
-- ----------------------------
DROP TABLE IF EXISTS "public"."whatsmeow_retry_buffer";
CREATE TABLE "public"."whatsmeow_retry_buffer" (
  "our_jid" text COLLATE "pg_catalog"."default" NOT NULL,
  "chat_jid" text COLLATE "pg_catalog"."default" NOT NULL,
  "message_id" text COLLATE "pg_catalog"."default" NOT NULL,
  "format" text COLLATE "pg_catalog"."default" NOT NULL,
  "plaintext" bytea NOT NULL,
  "timestamp" int8 NOT NULL
)
;

-- ----------------------------
-- Table structure for whatsmeow_sender_keys
-- ----------------------------
DROP TABLE IF EXISTS "public"."whatsmeow_sender_keys";
CREATE TABLE "public"."whatsmeow_sender_keys" (
  "our_jid" text COLLATE "pg_catalog"."default" NOT NULL,
  "chat_id" text COLLATE "pg_catalog"."default" NOT NULL,
  "sender_id" text COLLATE "pg_catalog"."default" NOT NULL,
  "sender_key" bytea NOT NULL
)
;

-- ----------------------------
-- Table structure for whatsmeow_sessions
-- ----------------------------
DROP TABLE IF EXISTS "public"."whatsmeow_sessions";
CREATE TABLE "public"."whatsmeow_sessions" (
  "our_jid" text COLLATE "pg_catalog"."default" NOT NULL,
  "their_id" text COLLATE "pg_catalog"."default" NOT NULL,
  "session" bytea
)
;

-- ----------------------------
-- Table structure for whatsmeow_version
-- ----------------------------
DROP TABLE IF EXISTS "public"."whatsmeow_version";
CREATE TABLE "public"."whatsmeow_version" (
  "version" int4,
  "compat" int4
)
;

-- ----------------------------
-- Function structure for armor
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."armor"(bytea, _text, _text);
CREATE FUNCTION "public"."armor"(bytea, _text, _text)
  RETURNS "pg_catalog"."text" AS '$libdir/pgcrypto', 'pg_armor'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for armor
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."armor"(bytea);
CREATE FUNCTION "public"."armor"(bytea)
  RETURNS "pg_catalog"."text" AS '$libdir/pgcrypto', 'pg_armor'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for crypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."crypt"(text, text);
CREATE FUNCTION "public"."crypt"(text, text)
  RETURNS "pg_catalog"."text" AS '$libdir/pgcrypto', 'pg_crypt'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for dearmor
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."dearmor"(text);
CREATE FUNCTION "public"."dearmor"(text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pg_dearmor'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for decrypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."decrypt"(bytea, bytea, text);
CREATE FUNCTION "public"."decrypt"(bytea, bytea, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pg_decrypt'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for decrypt_iv
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."decrypt_iv"(bytea, bytea, bytea, text);
CREATE FUNCTION "public"."decrypt_iv"(bytea, bytea, bytea, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pg_decrypt_iv'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for digest
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."digest"(text, text);
CREATE FUNCTION "public"."digest"(text, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pg_digest'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for digest
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."digest"(bytea, text);
CREATE FUNCTION "public"."digest"(bytea, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pg_digest'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for encrypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."encrypt"(bytea, bytea, text);
CREATE FUNCTION "public"."encrypt"(bytea, bytea, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pg_encrypt'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for encrypt_iv
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."encrypt_iv"(bytea, bytea, bytea, text);
CREATE FUNCTION "public"."encrypt_iv"(bytea, bytea, bytea, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pg_encrypt_iv'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for fips_mode
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."fips_mode"();
CREATE FUNCTION "public"."fips_mode"()
  RETURNS "pg_catalog"."bool" AS '$libdir/pgcrypto', 'pg_check_fipsmode'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for gen_random_bytes
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gen_random_bytes"(int4);
CREATE FUNCTION "public"."gen_random_bytes"(int4)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pg_random_bytes'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for gen_random_uuid
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gen_random_uuid"();
CREATE FUNCTION "public"."gen_random_uuid"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/pgcrypto', 'pg_random_uuid'
  LANGUAGE c VOLATILE
  COST 1;

-- ----------------------------
-- Function structure for gen_salt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gen_salt"(text, int4);
CREATE FUNCTION "public"."gen_salt"(text, int4)
  RETURNS "pg_catalog"."text" AS '$libdir/pgcrypto', 'pg_gen_salt_rounds'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for gen_salt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gen_salt"(text);
CREATE FUNCTION "public"."gen_salt"(text)
  RETURNS "pg_catalog"."text" AS '$libdir/pgcrypto', 'pg_gen_salt'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for hmac
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."hmac"(text, text, text);
CREATE FUNCTION "public"."hmac"(text, text, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pg_hmac'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for hmac
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."hmac"(bytea, bytea, text);
CREATE FUNCTION "public"."hmac"(bytea, bytea, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pg_hmac'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_armor_headers
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_armor_headers"(text, OUT "key" text, OUT "value" text);
CREATE FUNCTION "public"."pgp_armor_headers"(IN text, OUT "key" text, OUT "value" text)
  RETURNS SETOF "pg_catalog"."record" AS '$libdir/pgcrypto', 'pgp_armor_headers'
  LANGUAGE c IMMUTABLE STRICT
  COST 1
  ROWS 1000;

-- ----------------------------
-- Function structure for pgp_key_id
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_key_id"(bytea);
CREATE FUNCTION "public"."pgp_key_id"(bytea)
  RETURNS "pg_catalog"."text" AS '$libdir/pgcrypto', 'pgp_key_id_w'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_pub_decrypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_pub_decrypt"(bytea, bytea, text);
CREATE FUNCTION "public"."pgp_pub_decrypt"(bytea, bytea, text)
  RETURNS "pg_catalog"."text" AS '$libdir/pgcrypto', 'pgp_pub_decrypt_text'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_pub_decrypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_pub_decrypt"(bytea, bytea);
CREATE FUNCTION "public"."pgp_pub_decrypt"(bytea, bytea)
  RETURNS "pg_catalog"."text" AS '$libdir/pgcrypto', 'pgp_pub_decrypt_text'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_pub_decrypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_pub_decrypt"(bytea, bytea, text, text);
CREATE FUNCTION "public"."pgp_pub_decrypt"(bytea, bytea, text, text)
  RETURNS "pg_catalog"."text" AS '$libdir/pgcrypto', 'pgp_pub_decrypt_text'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_pub_decrypt_bytea
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_pub_decrypt_bytea"(bytea, bytea, text);
CREATE FUNCTION "public"."pgp_pub_decrypt_bytea"(bytea, bytea, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_pub_decrypt_bytea'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_pub_decrypt_bytea
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_pub_decrypt_bytea"(bytea, bytea, text, text);
CREATE FUNCTION "public"."pgp_pub_decrypt_bytea"(bytea, bytea, text, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_pub_decrypt_bytea'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_pub_decrypt_bytea
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_pub_decrypt_bytea"(bytea, bytea);
CREATE FUNCTION "public"."pgp_pub_decrypt_bytea"(bytea, bytea)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_pub_decrypt_bytea'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_pub_encrypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_pub_encrypt"(text, bytea, text);
CREATE FUNCTION "public"."pgp_pub_encrypt"(text, bytea, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_pub_encrypt_text'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_pub_encrypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_pub_encrypt"(text, bytea);
CREATE FUNCTION "public"."pgp_pub_encrypt"(text, bytea)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_pub_encrypt_text'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_pub_encrypt_bytea
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_pub_encrypt_bytea"(bytea, bytea);
CREATE FUNCTION "public"."pgp_pub_encrypt_bytea"(bytea, bytea)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_pub_encrypt_bytea'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_pub_encrypt_bytea
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_pub_encrypt_bytea"(bytea, bytea, text);
CREATE FUNCTION "public"."pgp_pub_encrypt_bytea"(bytea, bytea, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_pub_encrypt_bytea'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_sym_decrypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_sym_decrypt"(bytea, text);
CREATE FUNCTION "public"."pgp_sym_decrypt"(bytea, text)
  RETURNS "pg_catalog"."text" AS '$libdir/pgcrypto', 'pgp_sym_decrypt_text'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_sym_decrypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_sym_decrypt"(bytea, text, text);
CREATE FUNCTION "public"."pgp_sym_decrypt"(bytea, text, text)
  RETURNS "pg_catalog"."text" AS '$libdir/pgcrypto', 'pgp_sym_decrypt_text'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_sym_decrypt_bytea
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_sym_decrypt_bytea"(bytea, text, text);
CREATE FUNCTION "public"."pgp_sym_decrypt_bytea"(bytea, text, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_sym_decrypt_bytea'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_sym_decrypt_bytea
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_sym_decrypt_bytea"(bytea, text);
CREATE FUNCTION "public"."pgp_sym_decrypt_bytea"(bytea, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_sym_decrypt_bytea'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_sym_encrypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_sym_encrypt"(text, text, text);
CREATE FUNCTION "public"."pgp_sym_encrypt"(text, text, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_sym_encrypt_text'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_sym_encrypt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_sym_encrypt"(text, text);
CREATE FUNCTION "public"."pgp_sym_encrypt"(text, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_sym_encrypt_text'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_sym_encrypt_bytea
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_sym_encrypt_bytea"(bytea, text);
CREATE FUNCTION "public"."pgp_sym_encrypt_bytea"(bytea, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_sym_encrypt_bytea'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for pgp_sym_encrypt_bytea
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."pgp_sym_encrypt_bytea"(bytea, text, text);
CREATE FUNCTION "public"."pgp_sym_encrypt_bytea"(bytea, text, text)
  RETURNS "pg_catalog"."bytea" AS '$libdir/pgcrypto', 'pgp_sym_encrypt_bytea'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."lookup_id_seq"
OWNED BY "public"."lookup"."id";
SELECT setval('"public"."lookup_id_seq"', 12, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."migrations_id_seq"
OWNED BY "public"."migrations"."id";
SELECT setval('"public"."migrations_id_seq"', 23, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."system_settings_id_seq"
OWNED BY "public"."system_settings"."id";
SELECT setval('"public"."system_settings_id_seq"', 1, false);

-- ----------------------------
-- Indexes structure for table api_keys
-- ----------------------------
CREATE INDEX "idx_api_keys_hash" ON "public"."api_keys" USING btree (
  "key_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table api_keys
-- ----------------------------
ALTER TABLE "public"."api_keys" ADD CONSTRAINT "api_keys_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table api_logs
-- ----------------------------
CREATE INDEX "idx_api_logs_created_at" ON "public"."api_logs" USING btree (
  "created_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_api_logs_endpoint" ON "public"."api_logs" USING btree (
  "endpoint" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_api_logs_user_id" ON "public"."api_logs" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table api_logs
-- ----------------------------
ALTER TABLE "public"."api_logs" ADD CONSTRAINT "api_logs_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table audit_logs
-- ----------------------------
CREATE INDEX "idx_audit_logs_resource" ON "public"."audit_logs" USING btree (
  "resource_type" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "resource_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
) WHERE resource_id IS NOT NULL;
CREATE INDEX "idx_audit_logs_user_created" ON "public"."audit_logs" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
) WHERE user_id IS NOT NULL;

-- ----------------------------
-- Primary Key structure for table audit_logs
-- ----------------------------
ALTER TABLE "public"."audit_logs" ADD CONSTRAINT "audit_logs_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table auto_response
-- ----------------------------
CREATE INDEX "idx_auto_response_is_active" ON "public"."auto_response" USING btree (
  "is_active" "pg_catalog"."bool_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table auto_response
-- ----------------------------
ALTER TABLE "public"."auto_response" ADD CONSTRAINT "auto_response_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table auto_response_keywords
-- ----------------------------
CREATE INDEX "idx_auto_response_device_id" ON "public"."auto_response_keywords" USING btree (
  "device_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
) WHERE device_id IS NOT NULL;
CREATE INDEX "idx_auto_response_user_id" ON "public"."auto_response_keywords" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table auto_response_keywords
-- ----------------------------
ALTER TABLE "public"."auto_response_keywords" ADD CONSTRAINT "auto_response_keywords_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table auto_response_logs
-- ----------------------------
ALTER TABLE "public"."auto_response_logs" ADD CONSTRAINT "auto_response_logs_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table billing_plans
-- ----------------------------
CREATE INDEX "idx_billing_plans_active_price" ON "public"."billing_plans" USING btree (
  "is_active" "pg_catalog"."bool_ops" ASC NULLS LAST,
  "price" "pg_catalog"."numeric_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table billing_plans
-- ----------------------------
ALTER TABLE "public"."billing_plans" ADD CONSTRAINT "billing_plans_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table blacklists
-- ----------------------------
CREATE INDEX "idx_blacklists_user_id" ON "public"."blacklists" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table blacklists
-- ----------------------------
ALTER TABLE "public"."blacklists" ADD CONSTRAINT "blacklists_user_id_phone_number_key" UNIQUE ("user_id", "phone_number");

-- ----------------------------
-- Primary Key structure for table blacklists
-- ----------------------------
ALTER TABLE "public"."blacklists" ADD CONSTRAINT "blacklists_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table broadcast_campaigns
-- ----------------------------
CREATE INDEX "idx_broadcast_campaigns_deleted_at" ON "public"."broadcast_campaigns" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_broadcast_campaigns_device_id" ON "public"."broadcast_campaigns" USING btree (
  "device_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);
CREATE INDEX "idx_broadcast_campaigns_status" ON "public"."broadcast_campaigns" USING btree (
  "status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_broadcast_campaigns_user_id" ON "public"."broadcast_campaigns" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table broadcast_campaigns
-- ----------------------------
ALTER TABLE "public"."broadcast_campaigns" ADD CONSTRAINT "broadcast_campaigns_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table broadcast_messages
-- ----------------------------
ALTER TABLE "public"."broadcast_messages" ADD CONSTRAINT "broadcast_messages_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table broadcast_recipients
-- ----------------------------
CREATE INDEX "idx_broadcast_recipients_campaign_id" ON "public"."broadcast_recipients" USING btree (
  "campaign_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);
CREATE INDEX "idx_broadcast_recipients_status" ON "public"."broadcast_recipients" USING btree (
  "status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table broadcast_recipients
-- ----------------------------
ALTER TABLE "public"."broadcast_recipients" ADD CONSTRAINT "broadcast_recipients_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table contact
-- ----------------------------
CREATE INDEX "idx_contact_deleted_at" ON "public"."contact" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_contact_group_id" ON "public"."contact" USING btree (
  "group_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);
CREATE INDEX "idx_contact_phone" ON "public"."contact" USING btree (
  "phone" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table contact
-- ----------------------------
ALTER TABLE "public"."contact" ADD CONSTRAINT "contact_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Uniques structure for table contact_group_members
-- ----------------------------
ALTER TABLE "public"."contact_group_members" ADD CONSTRAINT "contact_group_members_group_id_contact_id_key" UNIQUE ("group_id", "contact_id");

-- ----------------------------
-- Primary Key structure for table contact_group_members
-- ----------------------------
ALTER TABLE "public"."contact_group_members" ADD CONSTRAINT "contact_group_members_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table contact_groups
-- ----------------------------
CREATE INDEX "idx_contact_groups_user_id" ON "public"."contact_groups" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table contact_groups
-- ----------------------------
ALTER TABLE "public"."contact_groups" ADD CONSTRAINT "contact_groups_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table contact_labels
-- ----------------------------
ALTER TABLE "public"."contact_labels" ADD CONSTRAINT "contact_labels_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table contacts
-- ----------------------------
CREATE INDEX "idx_contacts_label_id" ON "public"."contacts" USING btree (
  "label_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
) WHERE label_id IS NOT NULL;
CREATE INDEX "idx_contacts_user_name" ON "public"."contacts" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST,
  "name" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "idx_contacts_user_phone_active" ON "public"."contacts" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST,
  "phone_number" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
) WHERE deleted_at IS NULL;

-- ----------------------------
-- Primary Key structure for table contacts
-- ----------------------------
ALTER TABLE "public"."contacts" ADD CONSTRAINT "contacts_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table daily_message_stats
-- ----------------------------
CREATE INDEX "idx_daily_stats_user_date" ON "public"."daily_message_stats" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST,
  "stat_date" "pg_catalog"."date_ops" DESC NULLS FIRST
);
CREATE UNIQUE INDEX "idx_daily_stats_user_device" ON "public"."daily_message_stats" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST,
  "device_id" "pg_catalog"."uuid_ops" ASC NULLS LAST,
  "stat_date" "pg_catalog"."date_ops" ASC NULLS LAST
) WHERE device_id IS NOT NULL;
CREATE UNIQUE INDEX "idx_daily_stats_user_no_device" ON "public"."daily_message_stats" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST,
  "stat_date" "pg_catalog"."date_ops" ASC NULLS LAST
) WHERE device_id IS NULL;

-- ----------------------------
-- Uniques structure for table daily_message_stats
-- ----------------------------
ALTER TABLE "public"."daily_message_stats" ADD CONSTRAINT "daily_message_stats_unique_idx" UNIQUE ("user_id", "device_id", "stat_date");

-- ----------------------------
-- Checks structure for table daily_message_stats
-- ----------------------------
ALTER TABLE "public"."daily_message_stats" ADD CONSTRAINT "daily_message_stats_counts_check" CHECK (sent_count >= 0 AND failed_count >= 0 AND delivered_count >= 0 AND received_count >= 0);

-- ----------------------------
-- Primary Key structure for table daily_message_stats
-- ----------------------------
ALTER TABLE "public"."daily_message_stats" ADD CONSTRAINT "daily_message_stats_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table device_metrics
-- ----------------------------
CREATE INDEX "idx_device_metrics_device_id" ON "public"."device_metrics" USING btree (
  "device_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);

-- ----------------------------
-- Checks structure for table device_metrics
-- ----------------------------
ALTER TABLE "public"."device_metrics" ADD CONSTRAINT "device_metrics_success_rate_check" CHECK (success_rate IS NULL OR success_rate >= 0::numeric AND success_rate <= 100::numeric);

-- ----------------------------
-- Primary Key structure for table device_metrics
-- ----------------------------
ALTER TABLE "public"."device_metrics" ADD CONSTRAINT "device_metrics_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table device_qr_codes
-- ----------------------------
CREATE INDEX "idx_device_qr_device_id" ON "public"."device_qr_codes" USING btree (
  "device_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
) WHERE device_id IS NOT NULL;
CREATE INDEX "idx_device_qr_user_id" ON "public"."device_qr_codes" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table device_qr_codes
-- ----------------------------
ALTER TABLE "public"."device_qr_codes" ADD CONSTRAINT "device_qr_codes_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table device_sessions
-- ----------------------------
CREATE INDEX "idx_device_sessions_device_id" ON "public"."device_sessions" USING btree (
  "device_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);

-- ----------------------------
-- Checks structure for table device_sessions
-- ----------------------------
ALTER TABLE "public"."device_sessions" ADD CONSTRAINT "device_sessions_restart_count_check" CHECK (restart_count >= 0);

-- ----------------------------
-- Primary Key structure for table device_sessions
-- ----------------------------
ALTER TABLE "public"."device_sessions" ADD CONSTRAINT "device_sessions_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table devices
-- ----------------------------
CREATE INDEX "idx_devices_last_seen" ON "public"."devices" USING btree (
  "last_seen_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_devices_phone" ON "public"."devices" USING btree (
  "phone_number" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_devices_status" ON "public"."devices" USING btree (
  "status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_devices_user_id" ON "public"."devices" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table devices
-- ----------------------------
ALTER TABLE "public"."devices" ADD CONSTRAINT "devices_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table failure_records
-- ----------------------------
CREATE INDEX "idx_failure_records_device" ON "public"."failure_records" USING btree (
  "device_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
) WHERE device_id IS NOT NULL;
CREATE INDEX "idx_failure_records_user" ON "public"."failure_records" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST,
  "occurred_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);

-- ----------------------------
-- Primary Key structure for table failure_records
-- ----------------------------
ALTER TABLE "public"."failure_records" ADD CONSTRAINT "failure_records_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table groups
-- ----------------------------
CREATE INDEX "idx_groups_deleted_at" ON "public"."groups" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_groups_user_id" ON "public"."groups" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table groups
-- ----------------------------
ALTER TABLE "public"."groups" ADD CONSTRAINT "groups_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Checks structure for table invoice_items
-- ----------------------------
ALTER TABLE "public"."invoice_items" ADD CONSTRAINT "invoice_items_unit_price_check" CHECK (unit_price >= 0::numeric);
ALTER TABLE "public"."invoice_items" ADD CONSTRAINT "invoice_items_qty_check" CHECK (qty > 0);

-- ----------------------------
-- Primary Key structure for table invoice_items
-- ----------------------------
ALTER TABLE "public"."invoice_items" ADD CONSTRAINT "invoice_items_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table invoices
-- ----------------------------
CREATE INDEX "idx_invoices_status" ON "public"."invoices" USING btree (
  "status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_invoices_user_id" ON "public"."invoices" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table invoices
-- ----------------------------
ALTER TABLE "public"."invoices" ADD CONSTRAINT "invoices_invoice_number_key" UNIQUE ("invoice_number");

-- ----------------------------
-- Checks structure for table invoices
-- ----------------------------
ALTER TABLE "public"."invoices" ADD CONSTRAINT "invoices_amount_check" CHECK (amount >= 0::numeric);

-- ----------------------------
-- Primary Key structure for table invoices
-- ----------------------------
ALTER TABLE "public"."invoices" ADD CONSTRAINT "invoices_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table lookup
-- ----------------------------
ALTER TABLE "public"."lookup" ADD CONSTRAINT "lookup_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table message_templates
-- ----------------------------
CREATE INDEX "idx_message_templates_user_id" ON "public"."message_templates" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);

-- ----------------------------
-- Checks structure for table message_templates
-- ----------------------------
ALTER TABLE "public"."message_templates" ADD CONSTRAINT "message_templates_used_count_check" CHECK (used_count >= 0);

-- ----------------------------
-- Primary Key structure for table message_templates
-- ----------------------------
ALTER TABLE "public"."message_templates" ADD CONSTRAINT "message_templates_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table messages
-- ----------------------------
CREATE INDEX "idx_messages_created_at" ON "public"."messages" USING btree (
  "created_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_messages_device_id" ON "public"."messages" USING btree (
  "device_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);
CREATE INDEX "idx_messages_device_status_dir" ON "public"."messages" USING btree (
  "device_id" "pg_catalog"."uuid_ops" ASC NULLS LAST,
  "status_message" "pg_catalog"."int4_ops" ASC NULLS LAST,
  "direction" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_messages_direction" ON "public"."messages" USING btree (
  "direction" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_messages_receipt" ON "public"."messages" USING btree (
  "receipt_number" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_messages_scheduled_msg_id" ON "public"."messages" USING btree (
  "scheduled_message_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);
CREATE INDEX "idx_messages_status" ON "public"."messages" USING btree (
  "status_message" "pg_catalog"."int4_ops" ASC NULLS LAST
);
CREATE INDEX "idx_messages_whatsapp_message_id" ON "public"."messages" USING btree (
  "whatsapp_message_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
) WHERE whatsapp_message_id IS NOT NULL;

-- ----------------------------
-- Primary Key structure for table messages
-- ----------------------------
ALTER TABLE "public"."messages" ADD CONSTRAINT "messages_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Uniques structure for table migrations
-- ----------------------------
ALTER TABLE "public"."migrations" ADD CONSTRAINT "migrations_version_key" UNIQUE ("version");

-- ----------------------------
-- Primary Key structure for table migrations
-- ----------------------------
ALTER TABLE "public"."migrations" ADD CONSTRAINT "migrations_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Uniques structure for table notification_settings
-- ----------------------------
ALTER TABLE "public"."notification_settings" ADD CONSTRAINT "notification_settings_user_id_event_key_key" UNIQUE ("user_id", "event_key");

-- ----------------------------
-- Primary Key structure for table notification_settings
-- ----------------------------
ALTER TABLE "public"."notification_settings" ADD CONSTRAINT "notification_settings_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table notifications
-- ----------------------------
CREATE INDEX "idx_notifications_user_created" ON "public"."notifications" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);
CREATE INDEX "idx_notifications_user_unread" ON "public"."notifications" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST,
  "is_read" "pg_catalog"."bool_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table notifications
-- ----------------------------
ALTER TABLE "public"."notifications" ADD CONSTRAINT "notifications_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Uniques structure for table onboarding_progress
-- ----------------------------
ALTER TABLE "public"."onboarding_progress" ADD CONSTRAINT "onboarding_progress_user_id_step_key_key" UNIQUE ("user_id", "step_key");

-- ----------------------------
-- Primary Key structure for table onboarding_progress
-- ----------------------------
ALTER TABLE "public"."onboarding_progress" ADD CONSTRAINT "onboarding_progress_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table otp_verifications
-- ----------------------------
CREATE INDEX "idx_otp_expires_at" ON "public"."otp_verifications" USING btree (
  "expires_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_otp_phone_number" ON "public"."otp_verifications" USING btree (
  "phone_number" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_otp_verifications_expires_at" ON "public"."otp_verifications" USING btree (
  "expires_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_otp_verifications_phone_number" ON "public"."otp_verifications" USING btree (
  "phone_number" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_otp_verifications_user_id" ON "public"."otp_verifications" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);

-- ----------------------------
-- Checks structure for table otp_verifications
-- ----------------------------
ALTER TABLE "public"."otp_verifications" ADD CONSTRAINT "otp_verifications_attempt_count_check" CHECK (attempt_count >= 0);

-- ----------------------------
-- Primary Key structure for table otp_verifications
-- ----------------------------
ALTER TABLE "public"."otp_verifications" ADD CONSTRAINT "otp_verifications_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table resource_usage_metrics
-- ----------------------------
ALTER TABLE "public"."resource_usage_metrics" ADD CONSTRAINT "resource_usage_metrics_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Uniques structure for table scheduled_message_recipients
-- ----------------------------
ALTER TABLE "public"."scheduled_message_recipients" ADD CONSTRAINT "scheduled_message_recipients_unique" UNIQUE ("scheduled_message_id", "phone_number");

-- ----------------------------
-- Primary Key structure for table scheduled_message_recipients
-- ----------------------------
ALTER TABLE "public"."scheduled_message_recipients" ADD CONSTRAINT "scheduled_message_recipients_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table scheduled_messages
-- ----------------------------
CREATE INDEX "idx_scheduled_messages_status" ON "public"."scheduled_messages" USING btree (
  "status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_scheduled_messages_user" ON "public"."scheduled_messages" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST,
  "scheduled_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table scheduled_messages
-- ----------------------------
ALTER TABLE "public"."scheduled_messages" ADD CONSTRAINT "scheduled_messages_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table service_health_checks
-- ----------------------------
ALTER TABLE "public"."service_health_checks" ADD CONSTRAINT "service_health_checks_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table settings
-- ----------------------------
ALTER TABLE "public"."settings" ADD CONSTRAINT "settings_pkey" PRIMARY KEY ("key");

-- ----------------------------
-- Indexes structure for table subscriptions
-- ----------------------------
CREATE INDEX "idx_subscriptions_plan_id" ON "public"."subscriptions" USING btree (
  "plan_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);
CREATE INDEX "idx_subscriptions_status" ON "public"."subscriptions" USING btree (
  "status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_subscriptions_user_id" ON "public"."subscriptions" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);
CREATE INDEX "idx_subscriptions_user_status" ON "public"."subscriptions" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST,
  "status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table subscriptions
-- ----------------------------
ALTER TABLE "public"."subscriptions" ADD CONSTRAINT "subscriptions_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table system_settings
-- ----------------------------
ALTER TABLE "public"."system_settings" ADD CONSTRAINT "system_settings_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table usage_quotas
-- ----------------------------
CREATE INDEX "idx_usage_quotas_user_period" ON "public"."usage_quotas" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST,
  "period_key" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table usage_quotas
-- ----------------------------
ALTER TABLE "public"."usage_quotas" ADD CONSTRAINT "usage_quotas_user_id_period_key" UNIQUE ("user_id", "subscription_id", "period_key");

-- ----------------------------
-- Checks structure for table usage_quotas
-- ----------------------------
ALTER TABLE "public"."usage_quotas" ADD CONSTRAINT "usage_quotas_messages_used_check" CHECK (messages_used >= 0);
ALTER TABLE "public"."usage_quotas" ADD CONSTRAINT "usage_quotas_devices_used_check" CHECK (devices_used >= 0);

-- ----------------------------
-- Primary Key structure for table usage_quotas
-- ----------------------------
ALTER TABLE "public"."usage_quotas" ADD CONSTRAINT "usage_quotas_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table user_sessions
-- ----------------------------
CREATE INDEX "idx_user_sessions_expires_at" ON "public"."user_sessions" USING btree (
  "expires_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_user_sessions_refresh_token_hash" ON "public"."user_sessions" USING btree (
  "refresh_token_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "idx_user_sessions_session_token_hash" ON "public"."user_sessions" USING btree (
  "session_token_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_user_sessions_user_id" ON "public"."user_sessions" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table user_sessions
-- ----------------------------
ALTER TABLE "public"."user_sessions" ADD CONSTRAINT "user_sessions_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table users
-- ----------------------------
CREATE INDEX "idx_users_is_ban" ON "public"."users" USING btree (
  "is_banned" "pg_catalog"."bool_ops" ASC NULLS LAST
);
CREATE INDEX "idx_users_is_verify" ON "public"."users" USING btree (
  "is_verified" "pg_catalog"."bool_ops" ASC NULLS LAST
);
CREATE INDEX "idx_users_phone" ON "public"."users" USING btree (
  "phone_number" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "idx_users_phone_number" ON "public"."users" USING btree (
  "phone_number" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table users
-- ----------------------------
ALTER TABLE "public"."users" ADD CONSTRAINT "users_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table warming_pool
-- ----------------------------
CREATE INDEX "idx_warming_pool_device_id" ON "public"."warming_pool" USING btree (
  "device_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);
CREATE INDEX "idx_warming_pool_is_active" ON "public"."warming_pool" USING btree (
  "is_active" "pg_catalog"."bool_ops" ASC NULLS LAST
);
CREATE INDEX "idx_warming_pool_next_action" ON "public"."warming_pool" USING btree (
  "next_action_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table warming_pool
-- ----------------------------
ALTER TABLE "public"."warming_pool" ADD CONSTRAINT "warming_pool_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table warming_sessions
-- ----------------------------
CREATE INDEX "idx_warming_sessions_device_id" ON "public"."warming_sessions" USING btree (
  "device_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);
CREATE INDEX "idx_warming_sessions_status" ON "public"."warming_sessions" USING btree (
  "status" "pg_catalog"."int4_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table warming_sessions
-- ----------------------------
ALTER TABLE "public"."warming_sessions" ADD CONSTRAINT "warming_sessions_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table webhook_deliveries
-- ----------------------------
CREATE INDEX "idx_webhook_deliveries_webhook" ON "public"."webhook_deliveries" USING btree (
  "webhook_id" "pg_catalog"."uuid_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);

-- ----------------------------
-- Checks structure for table webhook_deliveries
-- ----------------------------
ALTER TABLE "public"."webhook_deliveries" ADD CONSTRAINT "webhook_deliveries_attempt_check" CHECK (attempt >= 1);

-- ----------------------------
-- Primary Key structure for table webhook_deliveries
-- ----------------------------
ALTER TABLE "public"."webhook_deliveries" ADD CONSTRAINT "webhook_deliveries_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Uniques structure for table webhook_event_subscriptions
-- ----------------------------
ALTER TABLE "public"."webhook_event_subscriptions" ADD CONSTRAINT "webhook_event_subscriptions_webhook_id_event_key_key" UNIQUE ("webhook_id", "event_key");

-- ----------------------------
-- Primary Key structure for table webhook_event_subscriptions
-- ----------------------------
ALTER TABLE "public"."webhook_event_subscriptions" ADD CONSTRAINT "webhook_event_subscriptions_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table webhook_settings
-- ----------------------------
ALTER TABLE "public"."webhook_settings" ADD CONSTRAINT "webhook_settings_pkey" PRIMARY KEY ("user_id");

-- ----------------------------
-- Indexes structure for table webhooks
-- ----------------------------
CREATE INDEX "idx_webhooks_device_id" ON "public"."webhooks" USING btree (
  "device_id" "pg_catalog"."uuid_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table webhooks
-- ----------------------------
ALTER TABLE "public"."webhooks" ADD CONSTRAINT "webhooks_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Checks structure for table whatsmeow_app_state_mutation_macs
-- ----------------------------
ALTER TABLE "public"."whatsmeow_app_state_mutation_macs" ADD CONSTRAINT "whatsmeow_app_state_mutation_macs_value_mac_check" CHECK (length(value_mac) = 32);
ALTER TABLE "public"."whatsmeow_app_state_mutation_macs" ADD CONSTRAINT "whatsmeow_app_state_mutation_macs_index_mac_check" CHECK (length(index_mac) = 32);

-- ----------------------------
-- Primary Key structure for table whatsmeow_app_state_mutation_macs
-- ----------------------------
ALTER TABLE "public"."whatsmeow_app_state_mutation_macs" ADD CONSTRAINT "whatsmeow_app_state_mutation_macs_pkey" PRIMARY KEY ("jid", "name", "version", "index_mac");

-- ----------------------------
-- Primary Key structure for table whatsmeow_app_state_sync_keys
-- ----------------------------
ALTER TABLE "public"."whatsmeow_app_state_sync_keys" ADD CONSTRAINT "whatsmeow_app_state_sync_keys_pkey" PRIMARY KEY ("jid", "key_id");

-- ----------------------------
-- Checks structure for table whatsmeow_app_state_version
-- ----------------------------
ALTER TABLE "public"."whatsmeow_app_state_version" ADD CONSTRAINT "whatsmeow_app_state_version_hash_check" CHECK (length(hash) = 128);

-- ----------------------------
-- Primary Key structure for table whatsmeow_app_state_version
-- ----------------------------
ALTER TABLE "public"."whatsmeow_app_state_version" ADD CONSTRAINT "whatsmeow_app_state_version_pkey" PRIMARY KEY ("jid", "name");

-- ----------------------------
-- Primary Key structure for table whatsmeow_chat_settings
-- ----------------------------
ALTER TABLE "public"."whatsmeow_chat_settings" ADD CONSTRAINT "whatsmeow_chat_settings_pkey" PRIMARY KEY ("our_jid", "chat_jid");

-- ----------------------------
-- Primary Key structure for table whatsmeow_contacts
-- ----------------------------
ALTER TABLE "public"."whatsmeow_contacts" ADD CONSTRAINT "whatsmeow_contacts_pkey" PRIMARY KEY ("our_jid", "their_jid");

-- ----------------------------
-- Checks structure for table whatsmeow_device
-- ----------------------------
ALTER TABLE "public"."whatsmeow_device" ADD CONSTRAINT "whatsmeow_device_noise_key_check" CHECK (length(noise_key) = 32);
ALTER TABLE "public"."whatsmeow_device" ADD CONSTRAINT "whatsmeow_device_identity_key_check" CHECK (length(identity_key) = 32);
ALTER TABLE "public"."whatsmeow_device" ADD CONSTRAINT "whatsmeow_device_signed_pre_key_check" CHECK (length(signed_pre_key) = 32);
ALTER TABLE "public"."whatsmeow_device" ADD CONSTRAINT "whatsmeow_device_signed_pre_key_id_check" CHECK (signed_pre_key_id >= 0 AND signed_pre_key_id < 16777216);
ALTER TABLE "public"."whatsmeow_device" ADD CONSTRAINT "whatsmeow_device_signed_pre_key_sig_check" CHECK (length(signed_pre_key_sig) = 64);
ALTER TABLE "public"."whatsmeow_device" ADD CONSTRAINT "whatsmeow_device_adv_account_sig_check" CHECK (length(adv_account_sig) = 64);
ALTER TABLE "public"."whatsmeow_device" ADD CONSTRAINT "whatsmeow_device_adv_account_sig_key_check" CHECK (length(adv_account_sig_key) = 32);
ALTER TABLE "public"."whatsmeow_device" ADD CONSTRAINT "whatsmeow_device_adv_device_sig_check" CHECK (length(adv_device_sig) = 64);
ALTER TABLE "public"."whatsmeow_device" ADD CONSTRAINT "whatsmeow_device_registration_id_check" CHECK (registration_id >= 0 AND registration_id < '4294967296'::bigint);

-- ----------------------------
-- Primary Key structure for table whatsmeow_device
-- ----------------------------
ALTER TABLE "public"."whatsmeow_device" ADD CONSTRAINT "whatsmeow_device_pkey" PRIMARY KEY ("jid");

-- ----------------------------
-- Checks structure for table whatsmeow_event_buffer
-- ----------------------------
ALTER TABLE "public"."whatsmeow_event_buffer" ADD CONSTRAINT "whatsmeow_event_buffer_ciphertext_hash_check" CHECK (length(ciphertext_hash) = 32);

-- ----------------------------
-- Primary Key structure for table whatsmeow_event_buffer
-- ----------------------------
ALTER TABLE "public"."whatsmeow_event_buffer" ADD CONSTRAINT "whatsmeow_event_buffer_pkey" PRIMARY KEY ("our_jid", "ciphertext_hash");

-- ----------------------------
-- Checks structure for table whatsmeow_identity_keys
-- ----------------------------
ALTER TABLE "public"."whatsmeow_identity_keys" ADD CONSTRAINT "whatsmeow_identity_keys_identity_check" CHECK (length(identity) = 32);

-- ----------------------------
-- Primary Key structure for table whatsmeow_identity_keys
-- ----------------------------
ALTER TABLE "public"."whatsmeow_identity_keys" ADD CONSTRAINT "whatsmeow_identity_keys_pkey" PRIMARY KEY ("our_jid", "their_id");

-- ----------------------------
-- Uniques structure for table whatsmeow_lid_map
-- ----------------------------
ALTER TABLE "public"."whatsmeow_lid_map" ADD CONSTRAINT "whatsmeow_lid_map_pn_key" UNIQUE ("pn");

-- ----------------------------
-- Primary Key structure for table whatsmeow_lid_map
-- ----------------------------
ALTER TABLE "public"."whatsmeow_lid_map" ADD CONSTRAINT "whatsmeow_lid_map_pkey" PRIMARY KEY ("lid");

-- ----------------------------
-- Primary Key structure for table whatsmeow_message_secrets
-- ----------------------------
ALTER TABLE "public"."whatsmeow_message_secrets" ADD CONSTRAINT "whatsmeow_message_secrets_pkey" PRIMARY KEY ("our_jid", "chat_jid", "sender_jid", "message_id");

-- ----------------------------
-- Checks structure for table whatsmeow_pre_keys
-- ----------------------------
ALTER TABLE "public"."whatsmeow_pre_keys" ADD CONSTRAINT "whatsmeow_pre_keys_key_check" CHECK (length(key) = 32);
ALTER TABLE "public"."whatsmeow_pre_keys" ADD CONSTRAINT "whatsmeow_pre_keys_key_id_check" CHECK (key_id >= 0 AND key_id < 16777216);

-- ----------------------------
-- Primary Key structure for table whatsmeow_pre_keys
-- ----------------------------
ALTER TABLE "public"."whatsmeow_pre_keys" ADD CONSTRAINT "whatsmeow_pre_keys_pkey" PRIMARY KEY ("jid", "key_id");

-- ----------------------------
-- Indexes structure for table whatsmeow_privacy_tokens
-- ----------------------------
CREATE INDEX "idx_whatsmeow_privacy_tokens_our_jid_timestamp" ON "public"."whatsmeow_privacy_tokens" USING btree (
  "our_jid" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "timestamp" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table whatsmeow_privacy_tokens
-- ----------------------------
ALTER TABLE "public"."whatsmeow_privacy_tokens" ADD CONSTRAINT "whatsmeow_privacy_tokens_pkey" PRIMARY KEY ("our_jid", "their_jid");

-- ----------------------------
-- Indexes structure for table whatsmeow_retry_buffer
-- ----------------------------
CREATE INDEX "whatsmeow_retry_buffer_timestamp_idx" ON "public"."whatsmeow_retry_buffer" USING btree (
  "our_jid" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "timestamp" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table whatsmeow_retry_buffer
-- ----------------------------
ALTER TABLE "public"."whatsmeow_retry_buffer" ADD CONSTRAINT "whatsmeow_retry_buffer_pkey" PRIMARY KEY ("our_jid", "chat_jid", "message_id");

-- ----------------------------
-- Primary Key structure for table whatsmeow_sender_keys
-- ----------------------------
ALTER TABLE "public"."whatsmeow_sender_keys" ADD CONSTRAINT "whatsmeow_sender_keys_pkey" PRIMARY KEY ("our_jid", "chat_id", "sender_id");

-- ----------------------------
-- Primary Key structure for table whatsmeow_sessions
-- ----------------------------
ALTER TABLE "public"."whatsmeow_sessions" ADD CONSTRAINT "whatsmeow_sessions_pkey" PRIMARY KEY ("our_jid", "their_id");

-- ----------------------------
-- Foreign Keys structure for table api_logs
-- ----------------------------
ALTER TABLE "public"."api_logs" ADD CONSTRAINT "api_logs_device_id_fkey" FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."api_logs" ADD CONSTRAINT "api_logs_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table auto_response
-- ----------------------------
ALTER TABLE "public"."auto_response" ADD CONSTRAINT "auto_response_device_id_fkey" FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table auto_response_logs
-- ----------------------------
ALTER TABLE "public"."auto_response_logs" ADD CONSTRAINT "auto_response_logs_keyword_id_fkey" FOREIGN KEY ("keyword_id") REFERENCES "public"."auto_response_keywords" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table broadcast_campaigns
-- ----------------------------
ALTER TABLE "public"."broadcast_campaigns" ADD CONSTRAINT "broadcast_campaigns_device_id_fkey" FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."broadcast_campaigns" ADD CONSTRAINT "broadcast_campaigns_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table broadcast_messages
-- ----------------------------
ALTER TABLE "public"."broadcast_messages" ADD CONSTRAINT "broadcast_messages_campaign_id_fkey" FOREIGN KEY ("campaign_id") REFERENCES "public"."broadcast_campaigns" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table broadcast_recipients
-- ----------------------------
ALTER TABLE "public"."broadcast_recipients" ADD CONSTRAINT "broadcast_recipients_campaign_id_fkey" FOREIGN KEY ("campaign_id") REFERENCES "public"."broadcast_campaigns" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."broadcast_recipients" ADD CONSTRAINT "broadcast_recipients_contact_id_fkey" FOREIGN KEY ("contact_id") REFERENCES "public"."contact" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."broadcast_recipients" ADD CONSTRAINT "broadcast_recipients_groups_id_fkey" FOREIGN KEY ("groups_id") REFERENCES "public"."groups" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table contact
-- ----------------------------
ALTER TABLE "public"."contact" ADD CONSTRAINT "contact_group_id_fkey" FOREIGN KEY ("group_id") REFERENCES "public"."groups" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table contact_group_members
-- ----------------------------
ALTER TABLE "public"."contact_group_members" ADD CONSTRAINT "contact_group_members_contact_id_fkey" FOREIGN KEY ("contact_id") REFERENCES "public"."contacts" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."contact_group_members" ADD CONSTRAINT "contact_group_members_group_id_fkey" FOREIGN KEY ("group_id") REFERENCES "public"."contact_groups" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table contacts
-- ----------------------------
ALTER TABLE "public"."contacts" ADD CONSTRAINT "contacts_label_id_fkey" FOREIGN KEY ("label_id") REFERENCES "public"."contact_labels" ("id") ON DELETE SET NULL ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table devices
-- ----------------------------
ALTER TABLE "public"."devices" ADD CONSTRAINT "devices_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table groups
-- ----------------------------
ALTER TABLE "public"."groups" ADD CONSTRAINT "groups_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table invoice_items
-- ----------------------------
ALTER TABLE "public"."invoice_items" ADD CONSTRAINT "invoice_items_invoice_id_fkey" FOREIGN KEY ("invoice_id") REFERENCES "public"."invoices" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table messages
-- ----------------------------
ALTER TABLE "public"."messages" ADD CONSTRAINT "messages_device_id_fkey" FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table scheduled_message_recipients
-- ----------------------------
ALTER TABLE "public"."scheduled_message_recipients" ADD CONSTRAINT "scheduled_message_recipients_contact_id_fkey" FOREIGN KEY ("contact_id") REFERENCES "public"."contacts" ("id") ON DELETE SET NULL ON UPDATE NO ACTION;
ALTER TABLE "public"."scheduled_message_recipients" ADD CONSTRAINT "scheduled_message_recipients_group_id_fkey" FOREIGN KEY ("group_id") REFERENCES "public"."contact_groups" ("id") ON DELETE SET NULL ON UPDATE NO ACTION;
ALTER TABLE "public"."scheduled_message_recipients" ADD CONSTRAINT "scheduled_message_recipients_scheduled_message_id_fkey" FOREIGN KEY ("scheduled_message_id") REFERENCES "public"."scheduled_messages" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table scheduled_messages
-- ----------------------------
ALTER TABLE "public"."scheduled_messages" ADD CONSTRAINT "scheduled_messages_group_id_fkey" FOREIGN KEY ("group_id") REFERENCES "public"."contact_groups" ("id") ON DELETE SET NULL ON UPDATE NO ACTION;
ALTER TABLE "public"."scheduled_messages" ADD CONSTRAINT "scheduled_messages_template_id_fkey" FOREIGN KEY ("template_id") REFERENCES "public"."message_templates" ("id") ON DELETE SET NULL ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table subscriptions
-- ----------------------------
ALTER TABLE "public"."subscriptions" ADD CONSTRAINT "subscriptions_plan_id_fkey" FOREIGN KEY ("plan_id") REFERENCES "public"."billing_plans" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."subscriptions" ADD CONSTRAINT "subscriptions_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table warming_pool
-- ----------------------------
ALTER TABLE "public"."warming_pool" ADD CONSTRAINT "warming_pool_device_id_fkey" FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table warming_sessions
-- ----------------------------
ALTER TABLE "public"."warming_sessions" ADD CONSTRAINT "warming_sessions_device_id_fkey" FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table webhooks
-- ----------------------------
ALTER TABLE "public"."webhooks" ADD CONSTRAINT "webhooks_device_id_fkey" FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table whatsmeow_app_state_mutation_macs
-- ----------------------------
ALTER TABLE "public"."whatsmeow_app_state_mutation_macs" ADD CONSTRAINT "whatsmeow_app_state_mutation_macs_jid_name_fkey" FOREIGN KEY ("jid", "name") REFERENCES "public"."whatsmeow_app_state_version" ("jid", "name") ON DELETE CASCADE ON UPDATE CASCADE;

-- ----------------------------
-- Foreign Keys structure for table whatsmeow_app_state_sync_keys
-- ----------------------------
ALTER TABLE "public"."whatsmeow_app_state_sync_keys" ADD CONSTRAINT "whatsmeow_app_state_sync_keys_jid_fkey" FOREIGN KEY ("jid") REFERENCES "public"."whatsmeow_device" ("jid") ON DELETE CASCADE ON UPDATE CASCADE;

-- ----------------------------
-- Foreign Keys structure for table whatsmeow_app_state_version
-- ----------------------------
ALTER TABLE "public"."whatsmeow_app_state_version" ADD CONSTRAINT "whatsmeow_app_state_version_jid_fkey" FOREIGN KEY ("jid") REFERENCES "public"."whatsmeow_device" ("jid") ON DELETE CASCADE ON UPDATE CASCADE;

-- ----------------------------
-- Foreign Keys structure for table whatsmeow_chat_settings
-- ----------------------------
ALTER TABLE "public"."whatsmeow_chat_settings" ADD CONSTRAINT "whatsmeow_chat_settings_our_jid_fkey" FOREIGN KEY ("our_jid") REFERENCES "public"."whatsmeow_device" ("jid") ON DELETE CASCADE ON UPDATE CASCADE;

-- ----------------------------
-- Foreign Keys structure for table whatsmeow_contacts
-- ----------------------------
ALTER TABLE "public"."whatsmeow_contacts" ADD CONSTRAINT "whatsmeow_contacts_our_jid_fkey" FOREIGN KEY ("our_jid") REFERENCES "public"."whatsmeow_device" ("jid") ON DELETE CASCADE ON UPDATE CASCADE;

-- ----------------------------
-- Foreign Keys structure for table whatsmeow_event_buffer
-- ----------------------------
ALTER TABLE "public"."whatsmeow_event_buffer" ADD CONSTRAINT "whatsmeow_event_buffer_our_jid_fkey" FOREIGN KEY ("our_jid") REFERENCES "public"."whatsmeow_device" ("jid") ON DELETE CASCADE ON UPDATE CASCADE;

-- ----------------------------
-- Foreign Keys structure for table whatsmeow_identity_keys
-- ----------------------------
ALTER TABLE "public"."whatsmeow_identity_keys" ADD CONSTRAINT "whatsmeow_identity_keys_our_jid_fkey" FOREIGN KEY ("our_jid") REFERENCES "public"."whatsmeow_device" ("jid") ON DELETE CASCADE ON UPDATE CASCADE;

-- ----------------------------
-- Foreign Keys structure for table whatsmeow_message_secrets
-- ----------------------------
ALTER TABLE "public"."whatsmeow_message_secrets" ADD CONSTRAINT "whatsmeow_message_secrets_our_jid_fkey" FOREIGN KEY ("our_jid") REFERENCES "public"."whatsmeow_device" ("jid") ON DELETE CASCADE ON UPDATE CASCADE;

-- ----------------------------
-- Foreign Keys structure for table whatsmeow_pre_keys
-- ----------------------------
ALTER TABLE "public"."whatsmeow_pre_keys" ADD CONSTRAINT "whatsmeow_pre_keys_jid_fkey" FOREIGN KEY ("jid") REFERENCES "public"."whatsmeow_device" ("jid") ON DELETE CASCADE ON UPDATE CASCADE;

-- ----------------------------
-- Foreign Keys structure for table whatsmeow_retry_buffer
-- ----------------------------
ALTER TABLE "public"."whatsmeow_retry_buffer" ADD CONSTRAINT "whatsmeow_retry_buffer_our_jid_fkey" FOREIGN KEY ("our_jid") REFERENCES "public"."whatsmeow_device" ("jid") ON DELETE CASCADE ON UPDATE CASCADE;

-- ----------------------------
-- Foreign Keys structure for table whatsmeow_sender_keys
-- ----------------------------
ALTER TABLE "public"."whatsmeow_sender_keys" ADD CONSTRAINT "whatsmeow_sender_keys_our_jid_fkey" FOREIGN KEY ("our_jid") REFERENCES "public"."whatsmeow_device" ("jid") ON DELETE CASCADE ON UPDATE CASCADE;

-- ----------------------------
-- Foreign Keys structure for table whatsmeow_sessions
-- ----------------------------
ALTER TABLE "public"."whatsmeow_sessions" ADD CONSTRAINT "whatsmeow_sessions_our_jid_fkey" FOREIGN KEY ("our_jid") REFERENCES "public"."whatsmeow_device" ("jid") ON DELETE CASCADE ON UPDATE CASCADE;
