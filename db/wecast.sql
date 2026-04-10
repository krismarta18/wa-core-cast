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

 Date: 09/04/2026 16:17:58
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
-- Table structure for api_logs
-- ----------------------------
DROP TABLE IF EXISTS "public"."api_logs";
CREATE TABLE "public"."api_logs" (
  "id" uuid NOT NULL,
  "userId" uuid,
  "endpoint" varchar(255) COLLATE "pg_catalog"."default",
  "reqBody" jsonb,
  "responseBody" jsonb,
  "created_at" varchar(255) COLLATE "pg_catalog"."default",
  "ipAddress" varchar(255) COLLATE "pg_catalog"."default",
  "deviceId" uuid
)
;

-- ----------------------------
-- Table structure for auto_response
-- ----------------------------
DROP TABLE IF EXISTS "public"."auto_response";
CREATE TABLE "public"."auto_response" (
  "id" uuid NOT NULL,
  "deviceId" uuid,
  "keyword" varchar(255) COLLATE "pg_catalog"."default",
  "response_text" varchar(255) COLLATE "pg_catalog"."default",
  "isActive" bool
)
;

-- ----------------------------
-- Table structure for billing_plans
-- ----------------------------
DROP TABLE IF EXISTS "public"."billing_plans";
CREATE TABLE "public"."billing_plans" (
  "id" uuid NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default",
  "price" numeric(10,2),
  "max_device" int4,
  "max_messages_day" int4,
  "features" jsonb
)
;

-- ----------------------------
-- Table structure for broadcast_campaigns
-- ----------------------------
DROP TABLE IF EXISTS "public"."broadcast_campaigns";
CREATE TABLE "public"."broadcast_campaigns" (
  "id" uuid NOT NULL,
  "userId" uuid,
  "deviceId" uuid,
  "nameBroadcast" varchar(255) COLLATE "pg_catalog"."default",
  "totalRecipients" int4,
  "processedCount" int4,
  "scheduled_at" timestamptz(6),
  "status" int4,
  "created_at" timestamptz(6),
  "deleted_at" timestamptz(6)
)
;

-- ----------------------------
-- Table structure for broadcast_messages
-- ----------------------------
DROP TABLE IF EXISTS "public"."broadcast_messages";
CREATE TABLE "public"."broadcast_messages" (
  "id" uuid NOT NULL,
  "campaignId" uuid,
  "messageType" int4,
  "messageText" varchar(255) COLLATE "pg_catalog"."default",
  "mediaUrl" varchar(255) COLLATE "pg_catalog"."default",
  "buttonData" jsonb
)
;

-- ----------------------------
-- Table structure for broadcast_recipients
-- ----------------------------
DROP TABLE IF EXISTS "public"."broadcast_recipients";
CREATE TABLE "public"."broadcast_recipients" (
  "id" uuid NOT NULL,
  "campaignId" uuid,
  "groupsId" uuid,
  "contactId" uuid,
  "status" int4,
  "sentAt" timestamptz(6),
  "errorMessages" varchar(255) COLLATE "pg_catalog"."default",
  "retryCount" int4
)
;

-- ----------------------------
-- Table structure for contact
-- ----------------------------
DROP TABLE IF EXISTS "public"."contact";
CREATE TABLE "public"."contact" (
  "id" uuid NOT NULL,
  "groupId" uuid,
  "name" varchar(100) COLLATE "pg_catalog"."default",
  "phone" varchar(30) COLLATE "pg_catalog"."default",
  "additional_data" jsonb,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6)
)
;

-- ----------------------------
-- Table structure for devices
-- ----------------------------
DROP TABLE IF EXISTS "public"."devices";
CREATE TABLE "public"."devices" (
  "id" uuid NOT NULL,
  "userId" uuid,
  "unique_name" varchar(100) COLLATE "pg_catalog"."default",
  "name_device" varchar(100) COLLATE "pg_catalog"."default",
  "phone" varchar(30) COLLATE "pg_catalog"."default",
  "status" int4,
  "last_seen" timestamptz(6),
  "session_data" bytea
)
;

-- ----------------------------
-- Table structure for groups
-- ----------------------------
DROP TABLE IF EXISTS "public"."groups";
CREATE TABLE "public"."groups" (
  "id" uuid NOT NULL,
  "userId" uuid,
  "groupName" varchar(100) COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6),
  "deleted_at" timestamptz(6)
)
;

-- ----------------------------
-- Table structure for lookup
-- ----------------------------
DROP TABLE IF EXISTS "public"."lookup";
CREATE TABLE "public"."lookup" (
  "id" int4 NOT NULL GENERATED ALWAYS AS IDENTITY (
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1
),
  "keys" varchar(255) COLLATE "pg_catalog"."default",
  "values" varchar(100) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for messages
-- ----------------------------
DROP TABLE IF EXISTS "public"."messages";
CREATE TABLE "public"."messages" (
  "id" uuid NOT NULL,
  "deviceId" uuid,
  "direction" bit(1),
  "receipt_number" varchar(30) COLLATE "pg_catalog"."default",
  "message_type" int4,
  "content" varchar(500) COLLATE "pg_catalog"."default",
  "status_message" int4,
  "error_log" varchar(255) COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6)
)
;

-- ----------------------------
-- Table structure for subscriptions
-- ----------------------------
DROP TABLE IF EXISTS "public"."subscriptions";
CREATE TABLE "public"."subscriptions" (
  "id" uuid NOT NULL,
  "userId" uuid,
  "planId" uuid,
  "status" int4,
  "created_at" timestamptz(6)
)
;

-- ----------------------------
-- Table structure for system_settings
-- ----------------------------
DROP TABLE IF EXISTS "public"."system_settings";
CREATE TABLE "public"."system_settings" (
  "id" int4 NOT NULL GENERATED ALWAYS AS IDENTITY (
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1
),
  "keys" varchar(100) COLLATE "pg_catalog"."default",
  "value" varchar(255) COLLATE "pg_catalog"."default",
  "description" varchar(255) COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6)
)
;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS "public"."users";
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL,
  "phone" varchar(30) COLLATE "pg_catalog"."default",
  "nama_lengkap" varchar(50) COLLATE "pg_catalog"."default",
  "is_verify" bool,
  "otp_code" varchar(10) COLLATE "pg_catalog"."default",
  "otp_expired" timestamptz(6),
  "idSubscribed" uuid,
  "max_device" int4,
  "is_ban" bool,
  "is_api" bool
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
  "dailyLimit" int4,
  "message_send_today" int4,
  "isActive" bool,
  "next_action_at" timestamptz(6)
)
;

-- ----------------------------
-- Table structure for warming_sessions
-- ----------------------------
DROP TABLE IF EXISTS "public"."warming_sessions";
CREATE TABLE "public"."warming_sessions" (
  "id" uuid NOT NULL,
  "deviceId" uuid,
  "target_phone" varchar(30) COLLATE "pg_catalog"."default",
  "message_sent" varchar(255) COLLATE "pg_catalog"."default",
  "response_received" varchar(255) COLLATE "pg_catalog"."default",
  "status" int4
)
;

-- ----------------------------
-- Table structure for webhooks
-- ----------------------------
DROP TABLE IF EXISTS "public"."webhooks";
CREATE TABLE "public"."webhooks" (
  "id" uuid NOT NULL,
  "deviceId" uuid,
  "webhookUrl" varchar(255) COLLATE "pg_catalog"."default",
  "secretKey" varchar(255) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."lookup_id_seq"
OWNED BY "public"."lookup"."id";
SELECT setval('"public"."lookup_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."system_settings_id_seq"
OWNED BY "public"."system_settings"."id";
SELECT setval('"public"."system_settings_id_seq"', 1, false);

-- ----------------------------
-- Primary Key structure for table api_logs
-- ----------------------------
ALTER TABLE "public"."api_logs" ADD CONSTRAINT "api_logs_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table auto_response
-- ----------------------------
ALTER TABLE "public"."auto_response" ADD CONSTRAINT "auto_response_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table billing_plans
-- ----------------------------
ALTER TABLE "public"."billing_plans" ADD CONSTRAINT "billing_plans_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table broadcast_campaigns
-- ----------------------------
ALTER TABLE "public"."broadcast_campaigns" ADD CONSTRAINT "broadcast_campaigns_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table broadcast_messages
-- ----------------------------
ALTER TABLE "public"."broadcast_messages" ADD CONSTRAINT "broadcast_messages_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table broadcast_recipients
-- ----------------------------
ALTER TABLE "public"."broadcast_recipients" ADD CONSTRAINT "broadast_recipients_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table contact
-- ----------------------------
ALTER TABLE "public"."contact" ADD CONSTRAINT "contact_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table devices
-- ----------------------------
ALTER TABLE "public"."devices" ADD CONSTRAINT "devices_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table groups
-- ----------------------------
ALTER TABLE "public"."groups" ADD CONSTRAINT "groups_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table lookup
-- ----------------------------
ALTER TABLE "public"."lookup" ADD CONSTRAINT "lookup_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table messages
-- ----------------------------
ALTER TABLE "public"."messages" ADD CONSTRAINT "messages_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table subscriptions
-- ----------------------------
ALTER TABLE "public"."subscriptions" ADD CONSTRAINT "subscriptions_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table system_settings
-- ----------------------------
ALTER TABLE "public"."system_settings" ADD CONSTRAINT "system_settings_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table users
-- ----------------------------
ALTER TABLE "public"."users" ADD CONSTRAINT "users_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table warming_pool
-- ----------------------------
ALTER TABLE "public"."warming_pool" ADD CONSTRAINT "warming_pool_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table warming_sessions
-- ----------------------------
ALTER TABLE "public"."warming_sessions" ADD CONSTRAINT "warming_sessions_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table webhooks
-- ----------------------------
ALTER TABLE "public"."webhooks" ADD CONSTRAINT "webhooks_pkey" PRIMARY KEY ("id");
