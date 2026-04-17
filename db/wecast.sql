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

 Date: 18/04/2026 01:25:51
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
  "key_prefix" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "key_hash" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "last_used_at" timestamptz(6),
  "expires_at" timestamptz(6),
  "is_active" bool NOT NULL DEFAULT true,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "revoked_at" timestamptz(6)
)
;

-- ----------------------------
-- Records of api_keys
-- ----------------------------

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
-- Records of api_logs
-- ----------------------------

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
-- Records of audit_logs
-- ----------------------------

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
-- Records of auto_response
-- ----------------------------

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
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;

-- ----------------------------
-- Records of auto_response_keywords
-- ----------------------------
INSERT INTO "public"."auto_response_keywords" VALUES ('fa0f0c73-2167-4f71-a7a0-d715d510ca8e', 'a9072252-9914-439c-a5c5-c2daaceffd52', NULL, 'hai, halo, selamat pagi, selamat siang, helo, hihi', 'Siapp Bos', 't', '2026-04-15 16:48:33.387646+07', '2026-04-15 17:27:17.774191+07');

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
-- Records of auto_response_logs
-- ----------------------------

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
-- Records of billing_plans
-- ----------------------------
INSERT INTO "public"."billing_plans" VALUES ('c2220d6e-b6c9-4fbf-9bd6-9d398a1e0ac6', 'Enterprise', 799000.00, 100, 100000, '["100 device", "100.000 pesan per hari", "Dedicated support"]', 'monthly', 't', '2026-04-15 09:33:58.359167+07', '2026-04-15 09:33:58.359167+07');
INSERT INTO "public"."billing_plans" VALUES ('527a5991-33f8-4897-8313-cf4f7223f28d', 'Business Pro', 299000.00, 10, 10000, '["10 device", "10.000 pesan per hari", "Priority support"]', 'monthly', 't', '2026-04-15 09:33:58.359167+07', '2026-04-15 09:33:58.359167+07');
INSERT INTO "public"."billing_plans" VALUES ('192600a7-e17f-41fa-8fba-c37296387dda', 'Starter', 99000.00, 2, 2000, '["2 device", "2.000 pesan per hari", "Basic analytics"]', 'monthly', 't', '2026-04-15 09:33:58.359167+07', '2026-04-15 09:33:58.359167+07');

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
-- Records of blacklists
-- ----------------------------

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
-- Records of broadcast_campaigns
-- ----------------------------
INSERT INTO "public"."broadcast_campaigns" VALUES ('b320e563-bb9e-4ecf-b9fc-4eb880eabfd6', 'a9072252-9914-439c-a5c5-c2daaceffd52', '1a47a884-5b7b-4977-9813-0617d9a1780c', NULL, 1, NULL, NULL, 'sending', '2026-04-15 15:16:07.26954+07', NULL, NULL, 'Te', 'Hllo', 5, 0, 0, NULL, NULL, '2026-04-15 15:16:07.3213+07');
INSERT INTO "public"."broadcast_campaigns" VALUES ('ec4f5ee6-33d5-4d2c-afce-5c44ed514ea7', 'a9072252-9914-439c-a5c5-c2daaceffd52', '1a47a884-5b7b-4977-9813-0617d9a1780c', NULL, 1, NULL, NULL, 'sending', '2026-04-15 15:17:29.00245+07', NULL, NULL, 'Test', 'Halo Test Broadcast', 5, 0, 0, NULL, NULL, '2026-04-15 15:17:29.039653+07');
INSERT INTO "public"."broadcast_campaigns" VALUES ('e6e50ea6-761d-4427-82ea-63cbb580e51a', 'a9072252-9914-439c-a5c5-c2daaceffd52', '1a47a884-5b7b-4977-9813-0617d9a1780c', NULL, 1, NULL, NULL, 'sending', '2026-04-15 15:19:20.114018+07', NULL, NULL, 'Ciiiihuyy', 'Cihuy', 5, 0, 0, NULL, NULL, '2026-04-15 15:19:20.135633+07');
INSERT INTO "public"."broadcast_campaigns" VALUES ('f31ec096-0562-4908-9cbf-d0d578670d98', 'a9072252-9914-439c-a5c5-c2daaceffd52', '1a47a884-5b7b-4977-9813-0617d9a1780c', NULL, 1, NULL, NULL, 'sending', '2026-04-15 15:21:46.581843+07', NULL, NULL, 's', 'ok', 5, 0, 0, NULL, NULL, '2026-04-15 15:21:46.617041+07');

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
-- Records of broadcast_messages
-- ----------------------------

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
-- Records of broadcast_recipients
-- ----------------------------
INSERT INTO "public"."broadcast_recipients" VALUES ('64d97da6-726d-4185-9475-e3d1b107d7cd', 'b320e563-bb9e-4ecf-b9fc-4eb880eabfd6', NULL, NULL, 'failed', NULL, NULL, 0, '2026-04-15 15:16:07.285599+07', NULL, '6285887373722', '2026-04-15 15:16:07.327109+07', 'failed to queue scheduled message: failed to enqueue message: pq: column "broadcast_id" of relation "messages" does not exist at position 5:19 (42703)');
INSERT INTO "public"."broadcast_recipients" VALUES ('36166f18-1097-420f-aafe-d4e038d4d448', 'ec4f5ee6-33d5-4d2c-afce-5c44ed514ea7', NULL, NULL, 'failed', NULL, NULL, 0, '2026-04-15 15:17:29.009494+07', NULL, '6285887373722', '2026-04-15 15:17:29.04583+07', 'failed to queue scheduled message: failed to enqueue message: pq: column "broadcast_id" of relation "messages" does not exist at position 5:19 (42703)');
INSERT INTO "public"."broadcast_recipients" VALUES ('c0b22bad-3ba5-455a-a9a5-4e117d5444f3', 'e6e50ea6-761d-4427-82ea-63cbb580e51a', NULL, NULL, 'failed', NULL, NULL, 0, '2026-04-15 15:19:20.116366+07', NULL, '6285887373722', '2026-04-15 15:19:20.137164+07', 'failed to queue scheduled message: failed to enqueue message: pq: column "broadcast_id" of relation "messages" does not exist at position 5:19 (42703)');
INSERT INTO "public"."broadcast_recipients" VALUES ('22285b0a-2cdf-410e-9814-6ab5c3893926', 'f31ec096-0562-4908-9cbf-d0d578670d98', NULL, NULL, 'pending', NULL, NULL, 0, '2026-04-15 15:21:46.587792+07', NULL, '6285887373722', NULL, NULL);

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
-- Records of contact
-- ----------------------------

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
-- Records of contact_group_members
-- ----------------------------
INSERT INTO "public"."contact_group_members" VALUES ('a36bd370-17b7-4300-b152-bdedd18cd6aa', '62b2931f-9414-48a6-970b-8b2153fa0cf5', 'a62539e4-5006-47d0-87dc-caab48488150', '2026-04-15 11:56:45.591015+07');

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
-- Records of contact_groups
-- ----------------------------
INSERT INTO "public"."contact_groups" VALUES ('62b2931f-9414-48a6-970b-8b2153fa0cf5', 'a9072252-9914-439c-a5c5-c2daaceffd52', 'Pelanggan VIP', 'Vip', '2026-04-15 11:56:39.302024+07', '2026-04-15 11:56:39.302024+07', NULL);

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
-- Records of contact_labels
-- ----------------------------

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
-- Records of contacts
-- ----------------------------
INSERT INTO "public"."contacts" VALUES ('a62539e4-5006-47d0-87dc-caab48488150', 'a9072252-9914-439c-a5c5-c2daaceffd52', NULL, 'Kristiatno', '6285887373722', 'null', 'o', '2026-04-15 11:56:16.98449+07', '2026-04-15 11:56:16.98449+07', NULL);

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
-- Records of daily_message_stats
-- ----------------------------
INSERT INTO "public"."daily_message_stats" VALUES ('bcef2db8-7182-4af2-87c3-de1ff58c8a68', 'a9072252-9914-439c-a5c5-c2daaceffd52', '1a47a884-5b7b-4977-9813-0617d9a1780c', '2026-04-15', 7, 0, 14, 0, 100.00, '2026-04-15 13:58:28.677279+07');
INSERT INTO "public"."daily_message_stats" VALUES ('b457c9f0-557a-48b6-b610-1454b3e93b6c', 'a9072252-9914-439c-a5c5-c2daaceffd52', '1a47a884-5b7b-4977-9813-0617d9a1780c', '2026-04-16', 0, 0, 1, 0, 100.00, '2026-04-16 10:18:31.188742+07');

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
-- Records of device_metrics
-- ----------------------------

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
-- Records of device_qr_codes
-- ----------------------------

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
-- Records of device_sessions
-- ----------------------------

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
-- Records of devices
-- ----------------------------
INSERT INTO "public"."devices" VALUES ('a347253a-4b92-4798-a56a-46e6be1df44e', 'a9072252-9914-439c-a5c5-c2daaceffd52', 'a347253a-4b92-4798-a56a-46e6be1df44e', 'CS 10', '62895321576222', 'banned', '2026-04-15 10:42:49.069182+07', NULL, '2026-04-15 09:35:09.501345+07', '2026-04-16 10:53:00.504741+07', NULL, NULL, NULL, NULL);
INSERT INTO "public"."devices" VALUES ('5562062c-bdd6-43c9-a0bc-a67d2579ba02', 'a9072252-9914-439c-a5c5-c2daaceffd52', '5562062c-bdd6-43c9-a0bc-a67d2579ba02', 'test', '6285887373722', 'banned', '2026-04-16 11:06:50.693959+07', NULL, '2026-04-16 11:03:31.74155+07', '2026-04-16 11:37:45.028449+07', NULL, NULL, NULL, NULL);
INSERT INTO "public"."devices" VALUES ('1a47a884-5b7b-4977-9813-0617d9a1780c', 'a9072252-9914-439c-a5c5-c2daaceffd52', '1a47a884-5b7b-4977-9813-0617d9a1780c', 'Customer Service 1', '62895321576222', 'banned', '2026-04-16 10:58:17.517636+07', NULL, '2026-04-15 10:57:50.403033+07', '2026-04-16 11:37:52.98229+07', NULL, NULL, NULL, NULL);

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
-- Records of failure_records
-- ----------------------------

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
-- Records of groups
-- ----------------------------

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
-- Records of invoice_items
-- ----------------------------

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
-- Records of invoices
-- ----------------------------
INSERT INTO "public"."invoices" VALUES ('f337457e-a195-41e0-be7e-da19c2edeaec', 'a9072252-9914-439c-a5c5-c2daaceffd52', '7027ea87-e6ce-4c71-b25e-f4dce16f445c', 'INV-2026-04-7027ea87', '2026-04-15', '2026-04-15', '2026-04-15 00:00:00+07', 99000.00, 'Rupiah', 'paid', 'dummy', '2026-04-15 00:00:00+07');

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
-- Records of lookup
-- ----------------------------
INSERT INTO "public"."lookup" VALUES (1, 'message_status.0', 'pending');
INSERT INTO "public"."lookup" VALUES (2, 'message_status.1', 'sent');
INSERT INTO "public"."lookup" VALUES (3, 'message_status.2', 'delivered');
INSERT INTO "public"."lookup" VALUES (4, 'message_status.3', 'read');
INSERT INTO "public"."lookup" VALUES (5, 'message_status.4', 'failed');
INSERT INTO "public"."lookup" VALUES (6, 'message_direction.OUT', 'outgoing');
INSERT INTO "public"."lookup" VALUES (7, 'message_direction.IN', 'incoming');
INSERT INTO "public"."lookup" VALUES (8, 'message_type.0', 'text');
INSERT INTO "public"."lookup" VALUES (9, 'message_type.1', 'image');
INSERT INTO "public"."lookup" VALUES (10, 'message_type.2', 'document');
INSERT INTO "public"."lookup" VALUES (11, 'message_type.3', 'audio');
INSERT INTO "public"."lookup" VALUES (12, 'message_type.4', 'video');

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
-- Records of message_templates
-- ----------------------------
INSERT INTO "public"."message_templates" VALUES ('4353acc8-73df-48cc-9e96-db446dd96ba9', 'a9072252-9914-439c-a5c5-c2daaceffd52', 'Test_tempalte', 'Umum', 'HaLLO', 0, '2026-04-15 16:48:07.268264+07', '2026-04-15 16:48:07.268264+07');

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
-- Records of messages
-- ----------------------------
INSERT INTO "public"."messages" VALUES ('5508c063-a6bb-498f-ac2c-652190c9d9f2', 'a347253a-4b92-4798-a56a-46e6be1df44e', 'OUT', '6285887373722', 0, 'Hello, this is a test message', 1, NULL, '2026-04-15 09:36:22.356392+07', NULL, '2026-04-15 09:36:22.440602+07', 3, 0, 3, NULL, NULL, NULL, NULL, NULL, NULL);
INSERT INTO "public"."messages" VALUES ('e2683f61-73dc-4cbf-ba3a-1a679f40adbf', 'a347253a-4b92-4798-a56a-46e6be1df44e', 'OUT', '6285887373722', 0, 'Hello, this is a test message', 1, NULL, '2026-04-15 09:40:36.729173+07', NULL, '2026-04-15 09:40:39.517138+07', 3, 0, 3, NULL, NULL, NULL, NULL, NULL, NULL);
INSERT INTO "public"."messages" VALUES ('af3c780c-8e94-45ea-8c5f-e77f681073f4', 'a347253a-4b92-4798-a56a-46e6be1df44e', 'OUT', '6285887373722', 0, 'Hello, this is a test message', 2, NULL, '2026-04-15 09:47:36.264508+07', NULL, '2026-04-15 09:47:58.822578+07', 3, 0, 3, NULL, '3EB0FEAE17564CD6BAC2B3', NULL, NULL, NULL, NULL);
INSERT INTO "public"."messages" VALUES ('2659839f-5f4a-4e76-8022-e295cab78fb0', 'a347253a-4b92-4798-a56a-46e6be1df44e', 'OUT', '6285887373722', 0, 'Hello, this is a test message', 2, NULL, '2026-04-15 09:47:37.067088+07', NULL, '2026-04-15 09:47:58.835986+07', 3, 0, 3, NULL, '3EB0A2968A51D120D25B92', NULL, NULL, NULL, NULL);
INSERT INTO "public"."messages" VALUES ('98c2af77-7e0f-4da2-9d66-dfc591a9fce0', 'a347253a-4b92-4798-a56a-46e6be1df44e', 'OUT', '6285887373722', 0, 'Hello, this is a test message', 2, NULL, '2026-04-15 09:51:20.678851+07', NULL, '2026-04-15 09:51:26.120026+07', 3, 0, 3, NULL, '3EB04BD95BC9288D9BB1B3', NULL, NULL, NULL, NULL);
INSERT INTO "public"."messages" VALUES ('50eb3ef0-8e33-4c64-9c1a-ec1ef98de5fa', 'a347253a-4b92-4798-a56a-46e6be1df44e', 'OUT', '6285887373722', 1, '', 1, NULL, '2026-04-15 09:57:39.845+07', NULL, '2026-04-15 09:57:42.325315+07', 2, 0, 3, NULL, '3EB00804891D3E2D27077D', 'https://lh3.googleusercontent.com/rd-gg/AEir0wJ2fc1CxnMkuJ0m6jpzsTiIcpG8s9pHdZ2-FgzxYJzFLuSDOo_hdv7pVLIaL3dV6iBt6uQQ_0Hrwgk3fzUCBHRuXN8-9-3ERNLCVCNAJmQFWL0RfSBaCLRc9hAVoyBXHDB5M_knW84TbFQavno31Q_jEgSi3bmVdsJPCMnhuzNHssIwq8y7l77UFJeq-_8dl1FVguRTexyyKO_HndFQybbupvUmDiBQOvFXFe0zyomcDzKf5QKyYAHWJ8sP2aSv_jRutjZSy2Qm_eJNNDIBOKUEMaTTlm7oUl4GGzzndR1LOSuRgdl-H_TnO4aA5IJwJFdP0eydXli00YwWxroW4Lm1I1Tt-g5m_46F95CzPWos4H3xHiYuz8dmjrMVCLHsaaIOTn0jcSp7z0-fSVVHo09ih2yK2DQQ22qk3tpGtG_tf5Kxfmm5zy3-_G6fW_Rj-I8GcrcKpR6GbtrAPuvy2Urqus5IbeFx5cE50CFGtVnGahSSyv02OwmU37G_5Vbq6hdQCu9-NL6lVG0huFPUH8cmv2t0OPkyfWW-8wm_EnzIRIt9mLZ0_ckf1u1iNU4xWuX6qdqN_OudUW7JEF0EJWocHEqRV7LjeR6ZfxF5gFk2RKjAWIni6g4q2sPUUJ8xNCKN1cYSMuqq_T3LbI4gfAtnzhuLBqkePaZtByJAGuceett3b3xzf5ucN02dKEmy-eUbK_PdgPpWLR9ua02Y-GTo6D0_y5m6kiRs8QaFOix7Ax4IXS54NkYZQU1QuR00ULyDWkrRUHDIYgGW8leuvsJcXEVdbol60C3jTvwQNglQNDXeeyD7MzBAL7i-eziQ8C4Nl_kVcEO-VxCe8E8FibcQCsa8MUp7RrSnIHRIXlnmWADuhd_BmwYxfA8GC4y2ke5XdIvduXByb3_6OL8EEfDjy6uPpqeuUhj0uHb3_bFb-peodb0hLUqKorUsSH_XdP_7gFecNCE0DgrrS-NbSm_KJ46BWQJ-0HQExGqh0dUBiPiwYk_6BWXXdF7JdGTkKSSkwARyaiBI4P-TLSPXKlkSv4L5bLUlZ3FP8CEVy-niNC-bi4a1XvEc5huprfE97QuoUoHGdX6z8I4LYHcl3910AmNSo8gPzVpQbWiRDqcdnTGNFk2lSiKpYs4bZKjRuKFoeiwlb_sNROWPbzou-GkZMsB07w2mLHgIh0i_dsIl7A5WRQAtn4qdiorrkLn-QCgSdZdbMWr2liCLe_VClxGpK3SYZwy7-R_uTtw8HbeKWi-4jHlxfnOqr_C_fN-1Eu1iH85hZSY1HocVEAOzITiQBgfRh3e5JsKtkF7AnpJXK-Aa7IkNBT9OeWQD5ThSAqkIjv0yDyKRrGatdPmhN-YbP05nkW59rDPfZ8cUSGBzCq6ji0ucj4BuCIaT_dQtRJte3h53t5R3nNCBfHKpxIygU-M2b7fwqVb3p5ptd4AMKor_tNOIR372oPMYAtGi0zIhNtCHtLk5PWccujitJuSSUA9Ayk1hFqBzSHRhadiw0Poo_WvtKR-FXB8ohpnPtg=s1600', 'Check this out!', NULL, NULL);
INSERT INTO "public"."messages" VALUES ('d83353b6-fc6c-4147-9a44-fe5cbe5def3c', 'a347253a-4b92-4798-a56a-46e6be1df44e', 'OUT', '6285887373722', 1, '', 2, NULL, '2026-04-15 09:57:28.49088+07', NULL, '2026-04-15 09:57:45.971074+07', 2, 0, 3, NULL, '3EB0C4994E97B91B70FBFA', 'https://lh3.googleusercontent.com/rd-gg/AEir0wJ2fc1CxnMkuJ0m6jpzsTiIcpG8s9pHdZ2-FgzxYJzFLuSDOo_hdv7pVLIaL3dV6iBt6uQQ_0Hrwgk3fzUCBHRuXN8-9-3ERNLCVCNAJmQFWL0RfSBaCLRc9hAVoyBXHDB5M_knW84TbFQavno31Q_jEgSi3bmVdsJPCMnhuzNHssIwq8y7l77UFJeq-_8dl1FVguRTexyyKO_HndFQybbupvUmDiBQOvFXFe0zyomcDzKf5QKyYAHWJ8sP2aSv_jRutjZSy2Qm_eJNNDIBOKUEMaTTlm7oUl4GGzzndR1LOSuRgdl-H_TnO4aA5IJwJFdP0eydXli00YwWxroW4Lm1I1Tt-g5m_46F95CzPWos4H3xHiYuz8dmjrMVCLHsaaIOTn0jcSp7z0-fSVVHo09ih2yK2DQQ22qk3tpGtG_tf5Kxfmm5zy3-_G6fW_Rj-I8GcrcKpR6GbtrAPuvy2Urqus5IbeFx5cE50CFGtVnGahSSyv02OwmU37G_5Vbq6hdQCu9-NL6lVG0huFPUH8cmv2t0OPkyfWW-8wm_EnzIRIt9mLZ0_ckf1u1iNU4xWuX6qdqN_OudUW7JEF0EJWocHEqRV7LjeR6ZfxF5gFk2RKjAWIni6g4q2sPUUJ8xNCKN1cYSMuqq_T3LbI4gfAtnzhuLBqkePaZtByJAGuceett3b3xzf5ucN02dKEmy-eUbK_PdgPpWLR9ua02Y-GTo6D0_y5m6kiRs8QaFOix7Ax4IXS54NkYZQU1QuR00ULyDWkrRUHDIYgGW8leuvsJcXEVdbol60C3jTvwQNglQNDXeeyD7MzBAL7i-eziQ8C4Nl_kVcEO-VxCe8E8FibcQCsa8MUp7RrSnIHRIXlnmWADuhd_BmwYxfA8GC4y2ke5XdIvduXByb3_6OL8EEfDjy6uPpqeuUhj0uHb3_bFb-peodb0hLUqKorUsSH_XdP_7gFecNCE0DgrrS-NbSm_KJ46BWQJ-0HQExGqh0dUBiPiwYk_6BWXXdF7JdGTkKSSkwARyaiBI4P-TLSPXKlkSv4L5bLUlZ3FP8CEVy-niNC-bi4a1XvEc5huprfE97QuoUoHGdX6z8I4LYHcl3910AmNSo8gPzVpQbWiRDqcdnTGNFk2lSiKpYs4bZKjRuKFoeiwlb_sNROWPbzou-GkZMsB07w2mLHgIh0i_dsIl7A5WRQAtn4qdiorrkLn-QCgSdZdbMWr2liCLe_VClxGpK3SYZwy7-R_uTtw8HbeKWi-4jHlxfnOqr_C_fN-1Eu1iH85hZSY1HocVEAOzITiQBgfRh3e5JsKtkF7AnpJXK-Aa7IkNBT9OeWQD5ThSAqkIjv0yDyKRrGatdPmhN-YbP05nkW59rDPfZ8cUSGBzCq6ji0ucj4BuCIaT_dQtRJte3h53t5R3nNCBfHKpxIygU-M2b7fwqVb3p5ptd4AMKor_tNOIR372oPMYAtGi0zIhNtCHtLk5PWccujitJuSSUA9Ayk1hFqBzSHRhadiw0Poo_WvtKR-FXB8ohpnPtg=s1600', 'Check this out!', NULL, NULL);
INSERT INTO "public"."messages" VALUES ('cdb54135-3c14-4c38-8cbf-424e94b4f878', 'a347253a-4b92-4798-a56a-46e6be1df44e', 'OUT', '6285887373722', 1, '', 1, NULL, '2026-04-15 09:58:09.228518+07', NULL, '2026-04-15 09:58:12.581599+07', 2, 0, 3, NULL, '3EB01422F9ADC78F4B6B78', 'https://lh3.googleusercontent.com/rd-gg/AEir0wJ2fc1CxnMkuJ0m6jpzsTiIcpG8s9pHdZ2-FgzxYJzFLuSDOo_hdv7pVLIaL3dV6iBt6uQQ_0Hrwgk3fzUCBHRuXN8-9-3ERNLCVCNAJmQFWL0RfSBaCLRc9hAVoyBXHDB5M_knW84TbFQavno31Q_jEgSi3bmVdsJPCMnhuzNHssIwq8y7l77UFJeq-_8dl1FVguRTexyyKO_HndFQybbupvUmDiBQOvFXFe0zyomcDzKf5QKyYAHWJ8sP2aSv_jRutjZSy2Qm_eJNNDIBOKUEMaTTlm7oUl4GGzzndR1LOSuRgdl-H_TnO4aA5IJwJFdP0eydXli00YwWxroW4Lm1I1Tt-g5m_46F95CzPWos4H3xHiYuz8dmjrMVCLHsaaIOTn0jcSp7z0-fSVVHo09ih2yK2DQQ22qk3tpGtG_tf5Kxfmm5zy3-_G6fW_Rj-I8GcrcKpR6GbtrAPuvy2Urqus5IbeFx5cE50CFGtVnGahSSyv02OwmU37G_5Vbq6hdQCu9-NL6lVG0huFPUH8cmv2t0OPkyfWW-8wm_EnzIRIt9mLZ0_ckf1u1iNU4xWuX6qdqN_OudUW7JEF0EJWocHEqRV7LjeR6ZfxF5gFk2RKjAWIni6g4q2sPUUJ8xNCKN1cYSMuqq_T3LbI4gfAtnzhuLBqkePaZtByJAGuceett3b3xzf5ucN02dKEmy-eUbK_PdgPpWLR9ua02Y-GTo6D0_y5m6kiRs8QaFOix7Ax4IXS54NkYZQU1QuR00ULyDWkrRUHDIYgGW8leuvsJcXEVdbol60C3jTvwQNglQNDXeeyD7MzBAL7i-eziQ8C4Nl_kVcEO-VxCe8E8FibcQCsa8MUp7RrSnIHRIXlnmWADuhd_BmwYxfA8GC4y2ke5XdIvduXByb3_6OL8EEfDjy6uPpqeuUhj0uHb3_bFb-peodb0hLUqKorUsSH_XdP_7gFecNCE0DgrrS-NbSm_KJ46BWQJ-0HQExGqh0dUBiPiwYk_6BWXXdF7JdGTkKSSkwARyaiBI4P-TLSPXKlkSv4L5bLUlZ3FP8CEVy-niNC-bi4a1XvEc5huprfE97QuoUoHGdX6z8I4LYHcl3910AmNSo8gPzVpQbWiRDqcdnTGNFk2lSiKpYs4bZKjRuKFoeiwlb_sNROWPbzou-GkZMsB07w2mLHgIh0i_dsIl7A5WRQAtn4qdiorrkLn-QCgSdZdbMWr2liCLe_VClxGpK3SYZwy7-R_uTtw8HbeKWi-4jHlxfnOqr_C_fN-1Eu1iH85hZSY1HocVEAOzITiQBgfRh3e5JsKtkF7AnpJXK-Aa7IkNBT9OeWQD5ThSAqkIjv0yDyKRrGatdPmhN-YbP05nkW59rDPfZ8cUSGBzCq6ji0ucj4BuCIaT_dQtRJte3h53t5R3nNCBfHKpxIygU-M2b7fwqVb3p5ptd4AMKor_tNOIR372oPMYAtGi0zIhNtCHtLk5PWccujitJuSSUA9Ayk1hFqBzSHRhadiw0Poo_WvtKR-FXB8ohpnPtg=s1600', 'Check this out!', NULL, NULL);
INSERT INTO "public"."messages" VALUES ('3fdfd7e7-5104-45e5-a019-b3ea7b466381', 'a347253a-4b92-4798-a56a-46e6be1df44e', 'OUT', '6285887373722', 1, '', 2, NULL, '2026-04-15 09:59:11.332061+07', NULL, '2026-04-15 09:59:24.018162+07', 2, 0, 3, NULL, '3EB01E2C19A6A919573912', 'https://lh3.googleusercontent.com/rd-gg/AEir0wJNpnGsJ2KjmQ6nazmfm8KlpJZxnwlkQv0B6X8cCVYpGcryrOgEe1Kd7zr6tVVgiaIoWeiCSKPWH511x7eaEifT7HuojGBLzVFSbkI4p3aTAaA9BqsV4E9IxpAAhvtCpylJMI_F-Z2gDQmmZDtVdf6gWdHEkTa_nskU5287Qkd_6AMkf2gcWeyexr-BnNARzgA06wuTZqy-r5gHiILG0BVgKjNIVF2tJv_m-aZHZeB9fi2K5MFVGLa6ZFeU846Onkc9pJkCs8UNDLk8GV9eijC7UQBDRULuaHSoUExsm_VU34OIuHDAXgWhlEFPm3Loxqfmr4CyHXptknv0A2ddGSNJPZYoI1S2sbtcE0nGnVUXHDbwi0QTgQNzJjqlYOlkg94zFUkaGKfFaDQk_bLAm5io3ulblGOj2kg-fBmepbUnq2auRGO0i3RNzCMnw5Qi5Xw48aUQfgVn-KZx6hJmugzx2poyhMVaB0I6gUjHmF70f7KirbMXm5iCaOoz5gfJmc7z_2NZ5JC8mXX-8dJJgu_f63IxqLcG0LJE4noUw66DAoXqdspMbnFCkWVYPIAF5j_H3qDbvv9m8Mgpf_BSgrasVi-eU4697Vnee335np_Vz6VJrPyWLWOdHYhhYUpPiqp9yzT1LPVxhjDtpkp7kRYSdPAj3G4u_EXz2a2k_PVLQopnnAca_JGyVI4-_TCIBd9JiYIVfrz7p3aal-KERZnVwiX9MmEBd50VA4gpt1x_fTx6pqAsSiq5chyqGjTvcrrxZ4l3wmVflxgqOJ9NX1675ioGNxIVl0uBfk0Vh5MQlyGs_eKA_RgBYxN--hXxpnerLtGPM8eu-X0sl87lI4R7c69xoEEuZgHbfjONg18t3xQfQt9MB7Q1iM4X27bIR5akP-UIQnQ0DF2yFJM9mEGayB6I6imlzkLOtcZAZPbCJlO4fZRDMROsexXuQ2WznX3KtkdLWffeA0bt6ZRmQpBopvKhT6_PzgfMr0iOXBssq-gKK3NJLdoQsfneEu63HJux_59DtwljojTE2gbBfOoJ2SSaqG6J1FHf38oQhSN7nkm76znkgk2Y3r-Qg0ZYYAiXNeN5ptOHy90__xYpNV9d4Y1qqGTKNk9YFuAB8yA781QwmaQEhyYyGHFM4FC4brjOPUkbbssS1L2-SPAlck2tQ61z6w_9YHNZdFsj9eDbA8maCOUBeg3HqP36JSVysoj4SOzWgzYgF_tR3BVVwW56R23EDwZaCsqrEFoxnF3lqqxGsSzjrF_1Vko4j_OmiKyOnPeBH8-N_ZkpUJvbz7elM05i6WAI4SHcBv69cIcAO866vO294Etl7JbQAaRJFOYW-m8etOLqkpvmQH8wqdl66TunJjuCbf-zIJafspelYkUbGOzT4fO7-4p91wpPHAVYBIe2mYad6kEprNaohtIMGKnf_2lElQ-k5x8R4tS6GVnSjb0Cm72z7xgCRCzCcuYE5EklbFTof2oLyhhLchl0o62MPlL2VRPM6xJ3B2GretnBlIfBBHruGBT4mOx1dQ=s1600', 'Check this out!Cuhi', NULL, NULL);
INSERT INTO "public"."messages" VALUES ('24e9514d-bd5b-4539-8b7a-b5ac227c9a68', '1a47a884-5b7b-4977-9813-0617d9a1780c', 'OUT', '6285887373722', 0, 'test', 2, NULL, '2026-04-15 13:46:22.540006+07', NULL, '2026-04-15 13:46:26.461897+07', 3, 0, 3, NULL, '3EB0E55C229436027C35B5', NULL, NULL, NULL, NULL);
INSERT INTO "public"."messages" VALUES ('a59ba91c-8090-4fdb-a1d0-ec08f9c887e2', 'a347253a-4b92-4798-a56a-46e6be1df44e', 'OUT', '6285887373722', 1, '', 1, NULL, '2026-04-15 10:02:31.132281+07', NULL, '2026-04-15 10:02:34.954184+07', 2, 0, 3, NULL, '3EB0BD052363883DF43E96', 'https://lh3.googleusercontent.com/rd-gg/AEir0wJNpnGsJ2KjmQ6nazmfm8KlpJZxnwlkQv0B6X8cCVYpGcryrOgEe1Kd7zr6tVVgiaIoWeiCSKPWH511x7eaEifT7HuojGBLzVFSbkI4p3aTAaA9BqsV4E9IxpAAhvtCpylJMI_F-Z2gDQmmZDtVdf6gWdHEkTa_nskU5287Qkd_6AMkf2gcWeyexr-BnNARzgA06wuTZqy-r5gHiILG0BVgKjNIVF2tJv_m-aZHZeB9fi2K5MFVGLa6ZFeU846Onkc9pJkCs8UNDLk8GV9eijC7UQBDRULuaHSoUExsm_VU34OIuHDAXgWhlEFPm3Loxqfmr4CyHXptknv0A2ddGSNJPZYoI1S2sbtcE0nGnVUXHDbwi0QTgQNzJjqlYOlkg94zFUkaGKfFaDQk_bLAm5io3ulblGOj2kg-fBmepbUnq2auRGO0i3RNzCMnw5Qi5Xw48aUQfgVn-KZx6hJmugzx2poyhMVaB0I6gUjHmF70f7KirbMXm5iCaOoz5gfJmc7z_2NZ5JC8mXX-8dJJgu_f63IxqLcG0LJE4noUw66DAoXqdspMbnFCkWVYPIAF5j_H3qDbvv9m8Mgpf_BSgrasVi-eU4697Vnee335np_Vz6VJrPyWLWOdHYhhYUpPiqp9yzT1LPVxhjDtpkp7kRYSdPAj3G4u_EXz2a2k_PVLQopnnAca_JGyVI4-_TCIBd9JiYIVfrz7p3aal-KERZnVwiX9MmEBd50VA4gpt1x_fTx6pqAsSiq5chyqGjTvcrrxZ4l3wmVflxgqOJ9NX1675ioGNxIVl0uBfk0Vh5MQlyGs_eKA_RgBYxN--hXxpnerLtGPM8eu-X0sl87lI4R7c69xoEEuZgHbfjONg18t3xQfQt9MB7Q1iM4X27bIR5akP-UIQnQ0DF2yFJM9mEGayB6I6imlzkLOtcZAZPbCJlO4fZRDMROsexXuQ2WznX3KtkdLWffeA0bt6ZRmQpBopvKhT6_PzgfMr0iOXBssq-gKK3NJLdoQsfneEu63HJux_59DtwljojTE2gbBfOoJ2SSaqG6J1FHf38oQhSN7nkm76znkgk2Y3r-Qg0ZYYAiXNeN5ptOHy90__xYpNV9d4Y1qqGTKNk9YFuAB8yA781QwmaQEhyYyGHFM4FC4brjOPUkbbssS1L2-SPAlck2tQ61z6w_9YHNZdFsj9eDbA8maCOUBeg3HqP36JSVysoj4SOzWgzYgF_tR3BVVwW56R23EDwZaCsqrEFoxnF3lqqxGsSzjrF_1Vko4j_OmiKyOnPeBH8-N_ZkpUJvbz7elM05i6WAI4SHcBv69cIcAO866vO294Etl7JbQAaRJFOYW-m8etOLqkpvmQH8wqdl66TunJjuCbf-zIJafspelYkUbGOzT4fO7-4p91wpPHAVYBIe2mYad6kEprNaohtIMGKnf_2lElQ-k5x8R4tS6GVnSjb0Cm72z7xgCRCzCcuYE5EklbFTof2oLyhhLchl0o62MPlL2VRPM6xJ3B2GretnBlIfBBHruGBT4mOx1dQ=s1600', 'Check this out!Cuhi', NULL, NULL);
INSERT INTO "public"."messages" VALUES ('1ca7de9d-eb42-4f3b-b76d-2dba67e42724', '1a47a884-5b7b-4977-9813-0617d9a1780c', 'OUT', '6285887373722', 0, 'ok', 2, NULL, '2026-04-15 15:21:46.619976+07', NULL, '2026-04-15 15:22:04.131467+07', 1, 0, 3, '2026-04-15 15:21:46.618868+07', '3EB0CB76C84A011B860DFB', NULL, NULL, 'f31ec096-0562-4908-9cbf-d0d578670d98', NULL);
INSERT INTO "public"."messages" VALUES ('2ef1fb37-7199-4e07-9bd0-f6d49189e33b', '1a47a884-5b7b-4977-9813-0617d9a1780c', 'inbound', '6285887373722:52@s.whatsapp.net', 0, 'halo', 2, NULL, '2026-04-15 17:17:55+07', NULL, '2026-04-15 17:18:00.541089+07', 3, 0, 3, NULL, '3EB097DA44E5225B61C38A', NULL, NULL, NULL, NULL);
INSERT INTO "public"."messages" VALUES ('ffe41fd2-4081-4359-a91d-6e1e72f4f729', 'a347253a-4b92-4798-a56a-46e6be1df44e', 'OUT', '6285887373722', 1, '', 1, NULL, '2026-04-15 10:02:51.159442+07', NULL, '2026-04-15 10:02:54.398968+07', 2, 0, 3, NULL, '3EB088D59731BB3ADB9D80', 'https://lh3.googleusercontent.com/rd-gg/AEir0wJNpnGsJ2KjmQ6nazmfm8KlpJZxnwlkQv0B6X8cCVYpGcryrOgEe1Kd7zr6tVVgiaIoWeiCSKPWH511x7eaEifT7HuojGBLzVFSbkI4p3aTAaA9BqsV4E9IxpAAhvtCpylJMI_F-Z2gDQmmZDtVdf6gWdHEkTa_nskU5287Qkd_6AMkf2gcWeyexr-BnNARzgA06wuTZqy-r5gHiILG0BVgKjNIVF2tJv_m-aZHZeB9fi2K5MFVGLa6ZFeU846Onkc9pJkCs8UNDLk8GV9eijC7UQBDRULuaHSoUExsm_VU34OIuHDAXgWhlEFPm3Loxqfmr4CyHXptknv0A2ddGSNJPZYoI1S2sbtcE0nGnVUXHDbwi0QTgQNzJjqlYOlkg94zFUkaGKfFaDQk_bLAm5io3ulblGOj2kg-fBmepbUnq2auRGO0i3RNzCMnw5Qi5Xw48aUQfgVn-KZx6hJmugzx2poyhMVaB0I6gUjHmF70f7KirbMXm5iCaOoz5gfJmc7z_2NZ5JC8mXX-8dJJgu_f63IxqLcG0LJE4noUw66DAoXqdspMbnFCkWVYPIAF5j_H3qDbvv9m8Mgpf_BSgrasVi-eU4697Vnee335np_Vz6VJrPyWLWOdHYhhYUpPiqp9yzT1LPVxhjDtpkp7kRYSdPAj3G4u_EXz2a2k_PVLQopnnAca_JGyVI4-_TCIBd9JiYIVfrz7p3aal-KERZnVwiX9MmEBd50VA4gpt1x_fTx6pqAsSiq5chyqGjTvcrrxZ4l3wmVflxgqOJ9NX1675ioGNxIVl0uBfk0Vh5MQlyGs_eKA_RgBYxN--hXxpnerLtGPM8eu-X0sl87lI4R7c69xoEEuZgHbfjONg18t3xQfQt9MB7Q1iM4X27bIR5akP-UIQnQ0DF2yFJM9mEGayB6I6imlzkLOtcZAZPbCJlO4fZRDMROsexXuQ2WznX3KtkdLWffeA0bt6ZRmQpBopvKhT6_PzgfMr0iOXBssq-gKK3NJLdoQsfneEu63HJux_59DtwljojTE2gbBfOoJ2SSaqG6J1FHf38oQhSN7nkm76znkgk2Y3r-Qg0ZYYAiXNeN5ptOHy90__xYpNV9d4Y1qqGTKNk9YFuAB8yA781QwmaQEhyYyGHFM4FC4brjOPUkbbssS1L2-SPAlck2tQ61z6w_9YHNZdFsj9eDbA8maCOUBeg3HqP36JSVysoj4SOzWgzYgF_tR3BVVwW56R23EDwZaCsqrEFoxnF3lqqxGsSzjrF_1Vko4j_OmiKyOnPeBH8-N_ZkpUJvbz7elM05i6WAI4SHcBv69cIcAO866vO294Etl7JbQAaRJFOYW-m8etOLqkpvmQH8wqdl66TunJjuCbf-zIJafspelYkUbGOzT4fO7-4p91wpPHAVYBIe2mYad6kEprNaohtIMGKnf_2lElQ-k5x8R4tS6GVnSjb0Cm72z7xgCRCzCcuYE5EklbFTof2oLyhhLchl0o62MPlL2VRPM6xJ3B2GretnBlIfBBHruGBT4mOx1dQ=s1600', 'Check this out!Cuhi', NULL, NULL);
INSERT INTO "public"."messages" VALUES ('b108b1a2-c050-4eae-aa24-6e43835633a7', '1a47a884-5b7b-4977-9813-0617d9a1780c', 'OUT', '6285887373722', 0, 'Krers', 2, NULL, '2026-04-15 13:58:24.402338+07', NULL, '2026-04-15 13:58:30.390722+07', 3, 0, 3, NULL, '3EB08A829F67E8C1D7E9C2', NULL, NULL, NULL, NULL);
INSERT INTO "public"."messages" VALUES ('788c3fac-219c-4cdd-bae4-19ddff41dc7a', '1a47a884-5b7b-4977-9813-0617d9a1780c', 'OUT', '6285887373722', 0, 'HIHIHIHI Test Lpg', 2, NULL, '2026-04-15 16:06:15.175531+07', NULL, '2026-04-15 16:07:07.95724+07', 1, 0, 3, '2026-04-15 16:07:00+07', '3EB0900A67C6B85929D955', NULL, NULL, NULL, NULL);
INSERT INTO "public"."messages" VALUES ('e9de0aaa-d64a-4de7-9fdb-d00f9a3f255e', '1a47a884-5b7b-4977-9813-0617d9a1780c', 'inbound', '6285887373722:52@s.whatsapp.net', 0, 'selamat pagi', 2, NULL, '2026-04-15 17:18:44+07', NULL, '2026-04-15 17:18:44.185568+07', 3, 0, 3, NULL, '3EB0B36C994280666664E2', NULL, NULL, NULL, NULL);
INSERT INTO "public"."messages" VALUES ('1c3fa94a-f9e0-4bc9-8aef-2aea57de76fa', 'a347253a-4b92-4798-a56a-46e6be1df44e', 'OUT', '6285887373722', 1, '', 2, NULL, '2026-04-15 10:03:24.097781+07', NULL, '2026-04-15 10:03:40.319552+07', 2, 0, 3, NULL, '3EB0FE299BE120173FBB44', 'https://lh3.googleusercontent.com/rd-gg/AEir0wIJc8r-NNeWw3vdDZgyx4IoXIfTtKngfPB-Nz7glChEhxgkeTJ7R2Q0kusEZkr-8TKX6GTqzeqGffFknJmRNgA7RYoCSsiOwRZhpLJ1TEQ866zrUPb_UgSYuuP1NElgK5g8y_o794zplIMVPRkIYRmHjz_zKqzffl8AUSTR09hdRbuePbhvsOClsBsHShh-dBoOslR7mVYsLhTbUW4oyS_peCLFsSCn5MhkmQsbUYvh-ms38Qp9ZE2fWYpAi2hrxGN-xs5EXAPb3yR7ABBLFGEdeqStbfvmQFq-xPSsPlMK0W5OjN11_VBt0M6X4zYsgtthCQJuvF_cp7UAfQ0ElnyaBf8rsO_trEn7Ub1qqRlZw3nNCkWjL2PbXL78baEJpTW8qJM3xilGfd23gxDro3UZW5zRnyfo-KJQkhrB00eOnYlkKg8Bt47WFfUbqluyBoYPrcaL8rPS0WOkeatOe4PLvtP0mHYFSRnrfmNRof9z7hNqtsbjGTGS8Tu9TKMzLaRZOSDkYjPdRgh1Wovr7mQhJ4Gzgn3c1bS0tlaPcEiOHmeH-lEKngAP92pPU-ch-ARlN0KZlqdr1xvni6JMAYW6eIYEYqS_jmMyWDck8Lm-Hmiovk4DewlrfF3VcaH-_0x6XwZhD6qTxx_O9OvwUBWTqx0rChQAUAX3Ot1juakjcULRwrgjTCQxVc9O-qvfCip4sH_MD-0wDk9H08nKz-mpv0MkwuojeOLl899hvfjP_Tc9Vnb2NCWA-0WLCJKNQFSyDhvy98Ub7zYYfuVe9zoxmCrMAre4Od0_FlohuSh2JWZMptwgGGIOca-_WSCJ-M8Jn9E1q_q-M_zMasm1MOtlpV-f9DTbk2oHEccdEfVdyY26l446XrLdCDbT_1JUUpuCSAGHyHwuS1sj3l1eJnAGdA7u8GF94kCg_zHzeuUl8lcTILK1HMm4TzIJ1_Mcdc4sx67TBjTHRb449W-e6Ph4yUWhd_GpTp1K3_88FfJWow2S2YS20pxOWIhEoAkBzemr2BJ7OrMHaiRcsU64A4LgfRGyJNeS2ZUMo6k7dALnDVNk7zKGjX8or5zUxYTYY6K7XodnjXYCwIA61wV3XYNMSPZDj1uAuvUC4sm2aLUKlVaKJ1LHSqA98t9AjzkCdl_eqJ6NKznhxsKNZ8PhRho8qSIz5KYOCaPJoXGF3QZ28mSFUbrd8lIZNdtYh-cwZNoR2RgM3vRBA4X7v6eNaTuJ16Z48YKsJoiZE1tFu7Vb2LL174_jz7Da5q5qDHEeo1nJJ6RcHOBr6la2ddyO9bMd9kGeAeczzFnZPKkAg_q9704REBoCveYkmRiszXNBS_RiDH0e_0fPxbEXr4JXdVTSQDjgID3-wJLQUDIE6uRmFzfbQ-c5n8hgpi2-lZtDIj4pT6OKLpZB0v_-NneDhXbcj55va6g5hyjS74Frzn95Z2Hqgj8PuYWQ5TqxXNZVZ1nJtyvNL8ElnJr4F92w2eQxDbqAdunJqUqPlR_SiVXwgYxJrajCIEthFeo4odsIig=s1600', 'Check this out!Cuhi', NULL, NULL);
INSERT INTO "public"."messages" VALUES ('c0acbeea-3739-4e00-9e6d-90c2752ea97a', 'a347253a-4b92-4798-a56a-46e6be1df44e', 'OUT', '6285887373722', 0, 'hehe', 2, NULL, '2026-04-15 10:23:13.335371+07', NULL, '2026-04-15 10:23:19.359047+07', 3, 0, 3, NULL, '3EB038B6024E8637C15A80', NULL, NULL, NULL, NULL);
INSERT INTO "public"."messages" VALUES ('387655c2-ed9f-4a7f-8dc3-5d4010f678f7', '1a47a884-5b7b-4977-9813-0617d9a1780c', 'OUT', '6285887373722', 0, 'HUU', 2, NULL, '2026-04-15 13:59:46.000368+07', NULL, '2026-04-15 14:00:05.457559+07', 1, 0, 3, '2026-04-15 14:00:00+07', '3EB0D3CCFAE5ECD9743809', NULL, NULL, NULL, NULL);
INSERT INTO "public"."messages" VALUES ('ff756da8-d1bb-4712-a7eb-42c5b37012a3', '1a47a884-5b7b-4977-9813-0617d9a1780c', 'inbound', '6285887373722:52@s.whatsapp.net', 0, 'halo', 3, NULL, '2026-04-15 17:14:10+07', NULL, '2026-04-15 17:14:11.089977+07', 3, 0, 3, NULL, '3EB017CAAB2797D213EA00', NULL, NULL, NULL, NULL);
INSERT INTO "public"."messages" VALUES ('bf46d25c-e6f3-40ff-9458-3700929c88f9', 'a347253a-4b92-4798-a56a-46e6be1df44e', 'OUT', '6285887373722', 2, '', 2, NULL, '2026-04-15 10:24:15.156493+07', NULL, '2026-04-15 10:24:23.916448+07', 2, 0, 3, NULL, '3EB0E764A4FCFD115F9E98', 'http://localhost:8080/uploads/1776223455154867400.pdf', 'XIxi', NULL, NULL);
INSERT INTO "public"."messages" VALUES ('2bb7dd54-253d-4e09-a0ba-5bcbe7fa0c6f', '1a47a884-5b7b-4977-9813-0617d9a1780c', 'inbound', '6285887373722:52@s.whatsapp.net', 0, 'cuek', 2, NULL, '2026-04-15 17:14:50+07', NULL, '2026-04-15 17:14:50.5756+07', 3, 0, 3, NULL, '3EB067B2AAE62EBAE0E60C', NULL, NULL, NULL, NULL);
INSERT INTO "public"."messages" VALUES ('39895ae9-28a6-41ed-86c2-0e74c9eeae8a', '1a47a884-5b7b-4977-9813-0617d9a1780c', 'inbound', '122045857812586@lid', 0, 'Halo', 2, NULL, '2026-04-15 17:20:25+07', NULL, '2026-04-15 17:20:25.635715+07', 3, 0, 3, NULL, '3A8AC07E653866DE2FE4', NULL, NULL, NULL, NULL);
INSERT INTO "public"."messages" VALUES ('f4767992-df3e-491a-9a50-9dd4384acb34', '1a47a884-5b7b-4977-9813-0617d9a1780c', 'outbound', '122045857812586@lid', 0, 'Siapp Bos', 2, NULL, '2026-04-15 17:20:25.64383+07', NULL, '2026-04-15 17:20:29.372042+07', 3, 0, 3, NULL, '3EB0D0F172E4D5A6CBF718', NULL, NULL, NULL, NULL);
INSERT INTO "public"."messages" VALUES ('675a2475-066d-446f-a433-ad2a317cd2d4', '1a47a884-5b7b-4977-9813-0617d9a1780c', 'inbound', '122045857812586@lid', 0, 'Halo', 2, NULL, '2026-04-15 17:20:46+07', NULL, '2026-04-15 17:20:46.15561+07', 3, 0, 3, NULL, '3A4FD9A0865D73005159', NULL, NULL, NULL, NULL);
INSERT INTO "public"."messages" VALUES ('e6fa37ce-ce25-4727-a88b-c4508d2682d1', '1a47a884-5b7b-4977-9813-0617d9a1780c', 'OUT', '628588737322', 0, 'Isi Deh', 4, 'Failed after 3 retries: failed to send message: no LID found for 628588737322@s.whatsapp.net from server', '2026-04-15 11:01:55.993552+07', NULL, '2026-04-15 11:03:07.804153+07', 1, 2, 3, '2026-04-15 11:03:00+07', NULL, NULL, NULL, NULL, NULL);
INSERT INTO "public"."messages" VALUES ('8bf24cc6-eac0-4212-afc6-c0c82ce5e401', '1a47a884-5b7b-4977-9813-0617d9a1780c', 'inbound', '122045857812586@lid', 0, 'Halo', 2, NULL, '2026-04-15 17:21:01+07', NULL, '2026-04-15 17:21:01.526928+07', 3, 0, 3, NULL, '3AFDBDF95B3107063CA0', NULL, NULL, NULL, NULL);
INSERT INTO "public"."messages" VALUES ('5c159f0e-6be4-4382-93c3-e5df5d0575c4', '1a47a884-5b7b-4977-9813-0617d9a1780c', 'OUT', '6285887373722', 0, 'Test Schedule', 2, NULL, '2026-04-15 11:11:00.357526+07', NULL, '2026-04-15 11:13:06.556808+07', 1, 0, 3, '2026-04-15 11:13:00+07', '3EB06F98427DFB3A109B67', NULL, NULL, NULL, NULL);
INSERT INTO "public"."messages" VALUES ('671d53ef-126f-4f5a-a96d-853871eb4cbf', '1a47a884-5b7b-4977-9813-0617d9a1780c', 'OUT', '6285887373722', 1, 'Test Gambar', 2, NULL, '2026-04-15 11:19:25.991196+07', NULL, '2026-04-15 11:21:18.472105+07', 1, 0, 3, '2026-04-15 11:21:00+07', '3EB0DBF68666E00DFEB286', 'http://localhost:8080/uploads/1776226765987181400.jpg', 'Test Gambar', NULL, NULL);
INSERT INTO "public"."messages" VALUES ('060312bb-a130-4475-8d52-9933d113fd27', '1a47a884-5b7b-4977-9813-0617d9a1780c', 'outbound', '122045857812586@lid', 0, 'Siapp Bos', 2, NULL, '2026-04-15 17:21:01.529814+07', NULL, '2026-04-16 10:18:31.181189+07', 3, 0, 3, NULL, '3EB06529338A535E657335', NULL, NULL, NULL, NULL);
INSERT INTO "public"."messages" VALUES ('df8b3028-5769-4ac3-86a8-307f22cd6589', '1a47a884-5b7b-4977-9813-0617d9a1780c', 'OUT', '6285887373722', 0, 'Test', 2, NULL, '2026-04-15 13:27:59.499054+07', NULL, '2026-04-15 13:28:03.5247+07', 3, 0, 3, NULL, '3EB08BA277A09749C981E5', NULL, NULL, NULL, NULL);
INSERT INTO "public"."messages" VALUES ('ee6504e5-10c1-42e2-8d6d-7e8f4ee01e9b', '1a47a884-5b7b-4977-9813-0617d9a1780c', 'OUT', '628588737322', 0, 'eehe', 4, 'Failed after 3 retries: failed to send message: no LID found for 628588737322@s.whatsapp.net from server', '2026-04-15 13:44:09.462686+07', NULL, '2026-04-15 13:44:18.644393+07', 3, 2, 3, NULL, NULL, NULL, NULL, NULL, NULL);

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
-- Records of migrations
-- ----------------------------
INSERT INTO "public"."migrations" VALUES (1, '001', '001_initial_schema', '2026-04-15 09:33:58.085006+07');
INSERT INTO "public"."migrations" VALUES (2, '002', '002_add_indexes', '2026-04-15 09:33:58.19282+07');
INSERT INTO "public"."migrations" VALUES (3, '003', '003_auth_schema_alignment', '2026-04-15 09:33:58.261283+07');
INSERT INTO "public"."migrations" VALUES (4, '004', '004_store_plaintext_otp', '2026-04-15 09:33:58.298849+07');
INSERT INTO "public"."migrations" VALUES (5, '005', '005_create_user_sessions', '2026-04-15 09:33:58.303424+07');
INSERT INTO "public"."migrations" VALUES (6, '006', '006_add_refresh_expiry_to_user_sessions', '2026-04-15 09:33:58.308984+07');
INSERT INTO "public"."migrations" VALUES (7, '007', '007_align_user_sessions_schema', '2026-04-15 09:33:58.313512+07');
INSERT INTO "public"."migrations" VALUES (8, '008', '008_align_billing_schema', '2026-04-15 09:33:58.357663+07');
INSERT INTO "public"."migrations" VALUES (9, '009', '009_seed_billing_plans', '2026-04-15 09:33:58.361452+07');
INSERT INTO "public"."migrations" VALUES (10, '010', '010_align_device_schema', '2026-04-15 09:33:58.402768+07');
INSERT INTO "public"."migrations" VALUES (11, '011', '011_align_messages_schema', '2026-04-15 09:33:58.426182+07');
INSERT INTO "public"."migrations" VALUES (12, '012', '012_add_whatsapp_message_id', '2026-04-15 09:46:37.076271+07');
INSERT INTO "public"."migrations" VALUES (13, '013', '013_seed_message_status_lookup', '2026-04-15 09:51:13.261245+07');
INSERT INTO "public"."migrations" VALUES (14, '014', '014_add_media_columns_to_messages', '2026-04-15 09:56:26.996745+07');
INSERT INTO "public"."migrations" VALUES (15, '015', '015_create_broadcast_tables', '2026-04-15 14:55:24.414119+07');
INSERT INTO "public"."migrations" VALUES (16, '016', '016_fix_broadcast_schema', '2026-04-15 14:58:19.757312+07');
INSERT INTO "public"."migrations" VALUES (17, '017', '017_fix_status_type', '2026-04-15 15:02:36.859312+07');
INSERT INTO "public"."migrations" VALUES (18, '018', '018_add_broadcast_id_to_messages', '2026-04-15 15:21:18.935158+07');
INSERT INTO "public"."migrations" VALUES (19, '019', '019_add_scheduled_message_id_to_messages', '2026-04-15 15:57:38.291699+07');
INSERT INTO "public"."migrations" VALUES (20, '020', '020_expand_message_columns', '2026-04-15 17:14:02.493512+07');
INSERT INTO "public"."migrations" VALUES (21, '021', '021_add_limits_to_subscriptions', '2026-04-16 11:51:09.933908+07');

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
-- Records of notification_settings
-- ----------------------------

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
-- Records of notifications
-- ----------------------------

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
-- Records of onboarding_progress
-- ----------------------------

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
-- Records of otp_verifications
-- ----------------------------
INSERT INTO "public"."otp_verifications" VALUES ('e82e6d94-7d56-4322-b7a1-73b216d320a9', 'a9072252-9914-439c-a5c5-c2daaceffd52', '6285887373722', 'register', '905596', 1, '2026-04-15 09:39:35.950724+07', '2026-04-15 09:34:58.75761+07', '2026-04-15 09:34:35.951265+07');
INSERT INTO "public"."otp_verifications" VALUES ('e897fafe-e3a6-4ccb-b39f-e85704e60358', 'a9072252-9914-439c-a5c5-c2daaceffd52', '6285887373722', 'login', '772789', 1, '2026-04-15 11:53:59.236588+07', '2026-04-15 11:49:11.068065+07', '2026-04-15 11:48:59.236876+07');
INSERT INTO "public"."otp_verifications" VALUES ('df58a500-b404-4072-941d-008599c3c345', 'a9072252-9914-439c-a5c5-c2daaceffd52', '6285887373722', 'login', '709622', 1, '2026-04-16 10:42:40.888707+07', '2026-04-16 10:39:04.880008+07', '2026-04-16 10:37:40.890674+07');

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
-- Records of resource_usage_metrics
-- ----------------------------

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
-- Records of scheduled_message_recipients
-- ----------------------------

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
-- Records of scheduled_messages
-- ----------------------------

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
-- Records of service_health_checks
-- ----------------------------

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
-- Records of subscriptions
-- ----------------------------
INSERT INTO "public"."subscriptions" VALUES ('7027ea87-e6ce-4c71-b25e-f4dce16f445c', 'a9072252-9914-439c-a5c5-c2daaceffd52', '192600a7-e17f-41fa-8fba-c37296387dda', 'active', '2026-04-15 17:27:41.354505+07', '2026-04-15 17:27:41.354505+07', '2026-05-15 17:27:41.354505+07', '2026-05-15 17:27:41.354505+07', 't', '2026-04-15 17:27:41.354505+07', 2, 2000);

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
-- Records of system_settings
-- ----------------------------

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
-- Records of usage_quotas
-- ----------------------------

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
-- Records of user_sessions
-- ----------------------------
INSERT INTO "public"."user_sessions" VALUES ('b8944e8a-06cf-4b47-a4cf-21e6c0df6227', 'a9072252-9914-439c-a5c5-c2daaceffd52', 'b66b40749d838e4002f20ef2f188349e20622274fccdefb13dcfa07d590b0a3b', '8b2816d588e91452ed27543de4910fbfc34402b1d67768527c545ccc2f22bdab', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:115.0) Gecko/20100101 Firefox/115.0', '2026-04-16 11:48:51.785716+07', '2026-04-15 11:48:54.132232+07', '2026-04-15 09:34:58.759751+07', '2026-04-22 11:48:51.785716+07', '2026-04-15 11:48:54.132232+07');
INSERT INTO "public"."user_sessions" VALUES ('f2a33c42-96ba-4981-9e88-d7ba5aaeffba', 'a9072252-9914-439c-a5c5-c2daaceffd52', '83daea490f82a3b5d457fd617ab9eba4e4821620fd570bde6d6d35713bbc9a09', '87174c27113d7c129b799435a8fdad65702b162b2bc07af75d17e634a673ce2c', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:115.0) Gecko/20100101 Firefox/115.0', '2026-04-16 11:55:11.00686+07', NULL, '2026-04-15 11:49:11.069222+07', '2026-04-22 11:55:11.00686+07', '2026-04-15 17:28:12.423655+07');
INSERT INTO "public"."user_sessions" VALUES ('beb86c1d-8a9d-4f92-964d-ddefd6b020a3', 'a9072252-9914-439c-a5c5-c2daaceffd52', 'f27ed601f2c3b2549ec60aa1e28d22ea8d7286405bfe437b5b36a6051fbd7a61', '9ca990a689318b4f6e5966a5b7aeb24f79ad61c39e9296374f74461e0857f042', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:115.0) Gecko/20100101 Firefox/115.0', '2026-04-17 10:39:04.885434+07', NULL, '2026-04-16 10:39:04.886565+07', '2026-04-23 10:39:04.885434+07', '2026-04-16 11:51:49.103748+07');

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
-- Records of users
-- ----------------------------
INSERT INTO "public"."users" VALUES ('a9072252-9914-439c-a5c5-c2daaceffd52', '6285887373722', 'Kristianto', 't', NULL, NULL, NULL, NULL, 'f', 'f', NULL, NULL, 'Asia/Jakarta', '2026-04-15 09:34:35.947594+07', '2026-04-16 10:39:04.881064+07', '2026-04-16 10:39:04.881064+07');

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
-- Records of warming_pool
-- ----------------------------

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
-- Records of warming_sessions
-- ----------------------------

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
-- Records of webhook_deliveries
-- ----------------------------

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
-- Records of webhook_event_subscriptions
-- ----------------------------

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
-- Records of webhooks
-- ----------------------------

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
-- Records of whatsmeow_app_state_mutation_macs
-- ----------------------------

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
-- Records of whatsmeow_app_state_sync_keys
-- ----------------------------

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
-- Records of whatsmeow_app_state_version
-- ----------------------------

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
-- Records of whatsmeow_chat_settings
-- ----------------------------

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
-- Records of whatsmeow_contacts
-- ----------------------------

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
-- Records of whatsmeow_device
-- ----------------------------

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
-- Records of whatsmeow_event_buffer
-- ----------------------------

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
-- Records of whatsmeow_identity_keys
-- ----------------------------

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
-- Records of whatsmeow_lid_map
-- ----------------------------
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('26496156590293', '62895321576222');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('235540318277727', '6283865006206');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('45745830908159', '628176326489');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('46295720943727', '6281328554961');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('245126702076157', '628562687387');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('250010264453291', '6281322127789');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('149151648116868', '6289501295270');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('200076806889619', '6281211498500');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('112777284812899', '628111717202');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('242687529713686', '6285213255275');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('11360456966381', '6287783992797');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('78744366641275', '6281219635731');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('60830292922578', '6285123602324');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('55865092653160', '6281400502839');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('86732737347741', '62895602537024');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('232701378461792', '6285157228461');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('56182685315316', '6282111452184');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('120229103407118', '6281375533315');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('4462940831942', '6285271234739');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('77159758598385', '6281316088009');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('175453893607605', '62895412881036');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('77859536302258', '6281283200655');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('238319547998428', '6281218205540');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('186895686467832', '6287828800501');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('215530116001905', '62895330913795');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('74964476653760', '6289613511110');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('159210276946132', '6282114531064');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('160292541608128', '6285742089627');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('48855370477785', '6281296422870');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('68904277815546', '6282113643390');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('75905108045909', '6285921966775');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('151157263605924', '6287770229889');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('249473225781285', '6285894285075');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('196018146689112', '6285600333606');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('78353239412772', '6281219255107');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('279898740535523', '6287741298160');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('44724081700944', '6281113804661');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('146119216631916', '6281387399808');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('203641730375929', '6281380531940');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('114048645492827', '6281286454359');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('53850367082641', '6281299367819');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('105055067877584', '6281295458055');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('210591474028772', '6289688229899');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('193458329411600', '6281388228005');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('53614127116340', '6282183283278');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('159455039717479', '6281294269955');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('212394856943696', '6285691641806');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('123192664400078', '6281384748377');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('175488303653071', '6287817138824');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('53627112702110', '628987455622');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('41837544890433', '6281389557992');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('18008932106423', '6287878316036');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('58059971899433', '6281229263150');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('193222022291615', '6281287300052');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('226272147992782', '62881024285065');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('76034393247862', '6281385564252');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('181961154932968', '6285923224848');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('27530824314943', '628128287097');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('143610855047231', '6281386195411');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('199531362807906', '6285883759042');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('213804059217933', '628974108999');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('167023292223645', '6285555558778');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('159266010837022', '628561910444');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('245998630731975', '6282126999865');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('113550496387232', '6281567987661');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('115062190649503', '6281905926712');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('202031067308106', '6285323699595');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('231730950762554', '6281210667031');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('45896171581651', '6281213656699');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('57578499371178', '6281919190199');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('66155549036705', '6281223361836');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('32513036669048', '6282125837169');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('114594391527428', '6281312022365');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('67822164127927', '6287777888998');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('14439797534921', '6281380807673');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('166709239533647', '6282142888852');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('83786272370934', '6281990888883');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('56865417343084', '6281331759016');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('226649736011781', '6282112071006');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('174152468164682', '628982811889');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('33165754269787', '628176626292');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('52128135569492', '6285713929007');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('130086036578415', '6281298361741');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('185723311382615', '6285813943587');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('148232357310529', '628972379451');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('180969101340905', '6282225257716');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('3363445993706', '6282289980024');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('101842751082719', '6282280006892');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('251367339913347', '6281296424900');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('144087713894528', '6281218649955');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('187376672468995', '6287891177899');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('187664301047947', '6285740040080');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('88725602181208', '6281295103636');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('10127683915936', '6281284628692');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('88643997855894', '6281329262606');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('24163569901702', '6281382470626');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('199570118156382', '6285810901351');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('228015921483777', '6281260099260');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('31275968630927', '62818889882');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('225094991376592', '628174167101');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('28849865756825', '6285716570935');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('281367669690533', '6281286510611');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('63724614369449', '6281585688088');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('89335487590448', '6282130083004');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('190262890475741', '62818999001');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('24906867716118', '6281210109489');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('258127735861374', '6281392194713');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('1915605745770', '6281808692255');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('134913982492674', '6285179860543');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('76321619214505', '6281388435563');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('265863223578669', '6281528344441');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('189653340680332', '6289525228259');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('214168661713013', '6285945583975');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('35218748649667', '6282113389066');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('86660125573209', '6285930952491');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('151067018944584', '6285602321005');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('43594001961078', '6282111721287');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('49035960434852', '628988361021');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('35742835306691', '6281295258056');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('145049836875852', '628157775300');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('146140708229164', '6287784462734');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('101090813063245', '6281287621446');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('154077807808645', '6283815551444');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('268491408039999', '628176467999');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('48743667757133', '6281226890403');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('89361441939645', '6289625884189');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('33737035227323', '6281398592490');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('246157544558845', '6283842748840');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('75626404933704', '6285285467218');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('150242402005213', '6285740332287');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('74367476203662', '6285710767026');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('258703697682589', '6281211336692');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('69196553699548', '6289523122979');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('106949131640939', '6281314930087');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('63264935432312', '6287882100244');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('155701590638592', '6285281488764');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('62642249044022', '6282114597990');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('62621109760186', '6283801812818');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('100682774401202', '6281288819242');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('47339230224606', '6289509144156');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('128617493344431', '6281248176375');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('185718932566057', '6281110512088');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('129235599528125', '628998847460');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('36271116292118', '6281197900067');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('120057740980410', '6287716736492');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('249851182903436', '6281290349227');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('211518582992915', '6288802660478');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('127835456942163', '6281389049390');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('46708172013761', '6288980526861');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('81239709106312', '6281295716507');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('33518444933292', '6285777007078');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('175075768705104', '6285172389182');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('219069840117761', '6281235551045');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('249851417727200', '6288975300744');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('92797180858396', '6287778334168');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('274538772361357', '628112652355');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('74733118877818', '628991626008');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('58179778031635', '628998988090');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('156843783487510', '6285163061296');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('74389001355338', '6285600406676');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('258857746079873', '6282122129122');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('192783985987770', '6283873262523');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('36528764014695', '6283892860230');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('123987266941178', '6281281166676');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('135403675877465', '628983619578');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('138366817411245', '6282112113539');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('157149061722138', '62859106980723');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('271081273356477', '6281385493235');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('247742488162538', '6289507191514');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('135875669319853', '6285719355633');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('69015896637449', '6281808747747');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('245977172725962', '6285775063022');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('242682999889947', '628119000845');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('139938741878917', '6289501117889');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('35648295702680', '6285966250097');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('133900034686993', '6289636390362');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('174234106102000', '6282310371152');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('167358198997002', '6282249606683');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('226250421440736', '6285652116202');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('103350133625075', '6287848741001');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('238796088049758', '6287770703770');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('31499424424099', '628999135281');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('176523155865735', '6281283061385');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('30799311171816', '6285312532599');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('151367733760230', '6285779493668');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('69561357455571', '6285216517893');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('247145990983858', '6285285977834');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('112940577460431', '6285880406770');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('34587472306179', '6285171108439');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('171296533029092', '6282123438360');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('257960685101133', '6281234502471');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('227479016976596', '6285176700934');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('9040973308074', '62811230666');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('7898461642797', '6281276833560');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('149546567000256', '6281314304460');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('190580734869740', '6281283182648');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('203319607861446', '6285775913000');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('53674424410284', '6287731111100');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('265416144326672', '6281287739677');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('167336506089678', '6285933483127');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('208645300195512', '6281808285021');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('43744577499263', '6281288822183');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('90276337025135', '6287865558111');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('148975302791319', '6282122912200');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('184456430276730', '6285876298892');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('129609295253606', '6283896176680');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('265712933330964', '6287722631243');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('153326087852224', '6281289665896');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('234702917111988', '6281399393969');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('222209410859097', '6282345827812');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('98492290724008', '6281223272741');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('189670335992029', '6287708779612');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('91504429277307', '6281299151942');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('204926177280136', '62881024109046');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('2190466879496', '628990301310');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('16419676794916', '6282110381848');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('68612555587616', '6285946927124');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('168706348974158', '6285311467860');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('77738975211767', '6281210307764');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('59884930363513', '6281389363669');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('144964273098976', '6282123433101');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('217081018630327', '6285251522827');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('260567596044391', '62895358073131');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('89498830532841', '628157942086');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('195893659713707', '6281249360453');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('184619286716497', '6285774049354');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('229995649736722', '6288214217963');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('68410658545679', '6285819945812');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('175200473755790', '6285814745820');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('17536603164826', '6287708967717');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('26968116453405', '6285377236807');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('225292626985159', '6285555558888');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('271356167979118', '6281280190035');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('72434690617396', '6283895778343');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('158948317507653', '628558855716');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('55615816790095', '6281188811195');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('193862022779008', '6289682501911');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('138061874770091', '62895413282783');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('77562998952134', '6282114240724');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('271798432231522', '62811841487');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('148344479432902', '6281994000229');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('177923315232863', '62895385899499');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('45006962311351', '6287833700297');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('150976824639692', '6282125372401');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('202636674519096', '6281218667722');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('109728176804067', '6287732004751');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('265669731979352', '628194040617');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('79448892268785', '6288226862236');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('160125557948467', '6282124171679');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('275496449364043', '6282113432547');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('280195210674421', '6281330405070');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('133015204315153', '6281227079687');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('196224523231459', '6281398035878');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('18249433514217', '6282122926648');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('81793659191498', '6281399936426');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('112434274693313', '6282246608889');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('28729170469096', '6281313995600');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('121487730192508', '6287743067101');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('89236636176417', '6287887647999');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('232938172117067', '6282146043046');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('200012332040422', '6287888828447');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('37456493686817', '6281296946003');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('166923987853334', '6281297164086');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('32517247754488', '62817459696');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('279095749414967', '6281290655599');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('131911448006873', '6287809474635');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('251182589173768', '6281386160914');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('156676179152932', '6287767711177');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('136893576572994', '6281281076427');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('207159325360382', '6285775092964');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('208288918589618', '6285742027386');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('182420263448826', '628111617758');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('59837585055770', '6283871374046');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('163939556048996', '6281318444198');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('154185165230113', '6281310156810');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('112639929745419', '6285282000540');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('89159326769284', '6281210191242');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('186994353258615', '6281222880101');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('90834464702531', '6281318388699');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('89731060736201', '6281220796588');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('274113553780794', '628992552786');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('206566888312872', '6281210004330');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('19082657161433', '6281380761403');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('204427877158948', '6281398176746');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('33320490561592', '628988833445');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('142477302440012', '628567921412');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('361012154479', '6281219080914');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('31928820441304', '6281283039995');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('134419625017457', '62895622601047');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('110565762552062', '62817181175');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('225846694572125', '6283879127951');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('228664092455034', '6285697942215');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('157424073863209', '6281234866928');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('99626715721969', '6285117487007');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('271940300312637', '6289610778254');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('77069195190360', '6281399913440');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('54438794399989', '6287862246402');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('37302176858333', '6282210231819');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('29309259501688', '6285942894883');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('70175554551916', '6285741647005');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('261683784216589', '6289515126984');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('18455524860083', '6285693995139');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('137683800215690', '6281391233165');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('206927262941373', '628888177888');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('276213759217778', '6285640588173');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('74560850415869', '6282112862828');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('91792192102453', '6281288527656');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('242485548851330', '6285217869596');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('7310655115341', '6285939771927');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('255610733989958', '6285777000446');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('109286063661071', '6287761072038');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('231056473108679', '6281289762911');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('8152049270839', '6281315024226');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('24356776358128', '6285718683556');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('65099171680380', '628170078883');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('700079730790', '6289518864610');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('109732287238217', '6281119164777');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('122711829426351', '6285210209860');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('258754666868776', '6285945435577');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('146046118301910', '6281389903241');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('273331852996840', '6281314430134');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('12245018878121', '6289651131819');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('59575625629832', '6281996257999');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('222054104227922', '6281287874427');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('58561929465980', '6282321614207');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('106085759369336', '6285810007177');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('186083870539991', '6289610900030');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('269741746819322', '6285123933227');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('12214937321625', '6287781221045');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('159979025735911', '628995058079');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('259416075014222', '6288226782127');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('140359648686157', '6287882082482');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('270024694587568', '6285795584378');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('170063759982688', '6287819999580');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('25662547087360', '6282111530495');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('115397433028763', '6281999691031');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('241510373167269', '6281283493237');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('83258511503531', '6283867509596');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('97379961298999', '6285775963605');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('240234801434799', '6285716005094');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('145462355034351', '6281387771783');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('277377745686715', '6283879391142');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('89151189856280', '6285137471552');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('238134830837910', '6287810327255');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('107919676850307', '6281282039988');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('12228141052124', '6281290008359');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('191327992053915', '6281282618800');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('108487065518283', '6289527670426');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('55912421146706', '6282249812779');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('45410806681839', '6281318557082');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('227208283058393', '62895608573030');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('205046100840567', '6287779845625');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('106446721171484', '6285651316176');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('261486786138272', '6287773322550');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('127066674573435', '6282112209688');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('58596356276371', '6285701802001');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('209384034574340', '6281358800989');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('27590920282269', '6281318723834');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('32036815413424', '6283834882902');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('142743372333238', '6282246124401');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('104544436547758', '6281315153969');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('172499459395669', '6281558333356');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('37409836282088', '6288215316530');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('200261742162013', '6281399977231');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('101262460752121', '6282360926395');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('71270771241134', '6281575060315');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('170764326170713', '628179166777');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('267499320934467', '6281286033686');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('44509333307445', '6287822175357');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('10711682048062', '6285883099658');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('133303034224705', '6281288781966');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('212893123485786', '6283125593160');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('160752187011166', '6285159548144');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('183790307676321', '6281283902909');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('274079194067020', '6285641031211');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('39827701498085', '62895366985599');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('51909427785897', '6281292887825');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('140600250740875', '6289653980828');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('26010959532050', '6285788886211');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('188618421346309', '6287832374528');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('39105777955019', '6285640173690');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('38831235535051', '6285213930688');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('196525120602322', '6281188807129');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('98698801483829', '6282210828697');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('89988171587692', '628562540781');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('174788022649025', '6281315611951');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('245882817609765', '6281389300953');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('21436332777618', '6281904863055');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('149374717943941', '628979132286');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('212162861641971', '6282114136198');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('225756416336054', '6282111412172');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('10063158710320', '6281295984003');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('226374941995087', '6287832000151');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('176093826928762', '6282125606512');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('66902823026698', '6281329006234');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('273314572447926', '6281398269773');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('52471699378202', '6285784195233');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('226782913519870', '6285777433433');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('4604272095484', '6282298868408');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('236008486482010', '6281211552714');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('124378427695332', '6281295223581');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('34660486783068', '6285211117999');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('238250895638742', '6281775222426');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('66843096129791', '6287878600876');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('212678559715414', '6281327142277');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('205844964724916', '6281221741953');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('108229048701176', '62895393346461');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('27913697104050', '6285117499947');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('67602986598566', '62817142288');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('249069482024980', '6281323323388');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('95928748904458', '6285921931653');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('265484930916411', '6285786959771');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('130794773299264', '6281386490044');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('172370224484564', '628987110808');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('255834189742172', '6281219924143');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('110170256457784', '6281932326788');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('122479800496193', '6281908192823');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('276278083096739', '6288213333889');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('100141558169660', '628111313977');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('161233089106103', '6285975389433');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('49006197624895', '62811270063');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('268585930891435', '6289527385080');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('223012016156854', '6281393847298');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('65197620334660', '6281385559739');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('94197860290749', '6287821699195');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('132547539378314', '6285213185681');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('280384323457047', '6281331845002');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('65799385550964', '6282333558808');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('249490388873274', '6282114794025');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('280122112364686', '6281213123477');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('179895242084384', '6285198885781');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('117751309926577', '6285777322143');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('185340774080569', '6287871558884');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('138611597021401', '628985944901');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('59012884262953', '628119677705');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('13843065475092', '6283878955787');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('177485564092506', '62895391936609');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('168835751645356', '6281299937098');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('119761069449307', '6281387915114');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('32225324220650', '6281382449705');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('76068249718808', '6287711121984');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('16287103205448', '6287832372863');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('74848495792177', '6285876303245');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('53228284698799', '6285113285123');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('261095474352144', '6285163231553');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('158115026706452', '6289616038998');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('65927865479176', '6283890299393');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('17772692164784', '6285773784316');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('74861330329610', '628985957257');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('113988918575334', '6281290808797');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('31821630804054', '6287780807999');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('5330524221556', '62816950727');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('85238373978303', '6289519007939');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('270501452746804', '6287871110565');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('192496290263041', '6281280138988');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('112876203311275', '6282140004030');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('149396226347222', '6281388084089');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('224562566414556', '6281779919815');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('107541820371107', '628111110821');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('42499070501081', '6282127771241');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('65103349153968', '6281213800793');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('273104236470478', '6282112288111');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('135287325888644', '6283893492428');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('189919293083725', '6285939492829');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('121775274897412', '6285600516002');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('9453457961155', '6287819995565');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('255770167898331', '6282125883697');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('214211644919925', '6282328312213');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('108160480198704', '6281286200781');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('102250387136572', '6281381265616');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('109044975030407', '6282310741899');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('97439939854354', '62816797511');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('102821668089895', '6282112068773');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('142833667313913', '6281213672128');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('148782012489822', '6287885726856');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('26311020023846', '6289502250582');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('183077393457324', '6281315726227');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('26757730168973', '62817808635');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('99561687240878', '6285811700721');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('67036034150598', '628111522646');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('148915273904371', '6281932133898');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('192874230595832', '6281905502332');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('106618419208416', '6282123726688');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('69196419432647', '6285780790855');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('121801195671799', '6287731924856');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('245474795765762', '6281366423724');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('279448054157422', '6285925378788');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('10862324625656', '6285280475365');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('236236170076346', '6282130000852');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('264334383046771', '6283199517989');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('27225864814810', '6281264960018');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('225464476008661', '6281285083458');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('172877315842132', '6285175190966');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('274409822654569', '6285959094477');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('206106974511259', '628118129817');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('147987477098661', '6287890909727');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('210358975377564', '6281903113557');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('17029679558899', '6287809999525');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('128149442564248', '628174913001');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('116350865387768', '6285664468767');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('75707589902370', '6281546027094');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('112459591463055', '6281517339422');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('162955488395354', '628981576577');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('49543152451770', '6287771115551');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('188089553162444', '6287796984192');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('171189192421376', '628885383838');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('113013508075577', '6285229337190');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('18489901318251', '6287780297899');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('161852252250355', '6285880555748');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('207163502928124', '62895342080706');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('104170724040949', '6289630505730');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('189636395671799', '6285939920250');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('114426686525632', '6281289740410');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('178919881883690', '628161960688');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('184237084958912', '62895345408726');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('40012083126383', '6289636071615');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('101769417904306', '6281298999212');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('17218691698755', '6285727611225');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('224841974202509', '6281295710072');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('50152702214252', '6289524081118');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('15251613466858', '6281384337126');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('53030095462493', '6285156780816');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('84620049686628', '6281234700393');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('157170553385029', '6285777511012');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('85397254209691', '6285179746108');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('71884934762703', '6282213807744');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('198169975603343', '6282117741199');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('141321822044391', '6281245605550');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('210457759621236', '6287743168350');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('209603010814102', '6287773847999');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('239109469692006', '6281285655969');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('122045857812586', '6285887373722');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('112412581712122', '6281282167087');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('165640363102209', '6289510224043');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('20603478216912', '62895322531726');
INSERT INTO "public"."whatsmeow_lid_map" VALUES ('95868267061402', '6282220072989');

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
-- Records of whatsmeow_message_secrets
-- ----------------------------

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
-- Records of whatsmeow_pre_keys
-- ----------------------------

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
-- Records of whatsmeow_privacy_tokens
-- ----------------------------
INSERT INTO "public"."whatsmeow_privacy_tokens" VALUES ('62895321576222:41@s.whatsapp.net', '122045857812586@lid', E'\\004\\001(\\367\\360\\330\\0225\\203q\\256', 1776246543, NULL);

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
-- Records of whatsmeow_retry_buffer
-- ----------------------------

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
-- Records of whatsmeow_sender_keys
-- ----------------------------

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
-- Records of whatsmeow_sessions
-- ----------------------------

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
-- Records of whatsmeow_version
-- ----------------------------
INSERT INTO "public"."whatsmeow_version" VALUES (13, 8);

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
SELECT setval('"public"."migrations_id_seq"', 21, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."system_settings_id_seq"
OWNED BY "public"."system_settings"."id";
SELECT setval('"public"."system_settings_id_seq"', 1, false);

-- ----------------------------
-- Indexes structure for table api_keys
-- ----------------------------
CREATE INDEX "idx_api_keys_user_active" ON "public"."api_keys" USING btree (
  "user_id" "pg_catalog"."uuid_ops" ASC NULLS LAST,
  "is_active" "pg_catalog"."bool_ops" ASC NULLS LAST
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
