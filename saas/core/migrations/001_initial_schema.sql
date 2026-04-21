-- Migration: 001_initial_schema
-- Description: Create initial database schema for WACAST

-- ----------------------------
-- Type structure for enum_direction
-- ----------------------------
DROP TYPE IF EXISTS "public"."enum_direction" CASCADE;
CREATE TYPE "public"."enum_direction" AS ENUM (
  'IN',
  'OUT'
);

-- ----------------------------
-- Table structure for billing_plans
-- ----------------------------
DROP TABLE IF EXISTS "public"."billing_plans" CASCADE;
CREATE TABLE "public"."billing_plans" (
  "id" uuid NOT NULL,
  "name" varchar(100),
  "price" numeric(10,2),
  "max_device" int4,
  "max_messages_day" int4,
  "features" jsonb,
  PRIMARY KEY ("id")
);

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS "public"."users" CASCADE;
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL,
  "phone" varchar(30),
  "nama_lengkap" varchar(50),
  "is_verify" bool,
  "otp_code" varchar(10),
  "otp_expired" timestamptz(6),
  "id_subscribed" uuid,
  "max_device" int4,
  "is_ban" bool,
  "is_api" bool,
  PRIMARY KEY ("id")
);

-- ----------------------------
-- Table structure for subscriptions
-- ----------------------------
DROP TABLE IF EXISTS "public"."subscriptions" CASCADE;
CREATE TABLE "public"."subscriptions" (
  "id" uuid NOT NULL,
  "user_id" uuid,
  "plan_id" uuid,
  "status" int4,
  "created_at" timestamptz(6),
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id"),
  FOREIGN KEY ("plan_id") REFERENCES "public"."billing_plans" ("id")
);

-- ----------------------------
-- Table structure for devices
-- ----------------------------
DROP TABLE IF EXISTS "public"."devices" CASCADE;
CREATE TABLE "public"."devices" (
  "id" uuid NOT NULL,
  "user_id" uuid,
  "unique_name" varchar(100),
  "name_device" varchar(100),
  "phone" varchar(30),
  "status" int4,
  "last_seen" timestamptz(6),
  "session_data" bytea,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id")
);

-- ----------------------------
-- Table structure for messages
-- ----------------------------
DROP TABLE IF EXISTS "public"."messages" CASCADE;
CREATE TABLE "public"."messages" (
  "id" uuid NOT NULL,
  "device_id" uuid,
  "direction" varchar(10),
  "receipt_number" varchar(30),
  "message_type" int4,
  "content" varchar(500),
  "status_message" int4,
  "error_log" varchar(255),
  "created_at" timestamptz(6),
  PRIMARY KEY ("id"),
  FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id")
);

-- ----------------------------
-- Table structure for groups
-- ----------------------------
DROP TABLE IF EXISTS "public"."groups" CASCADE;
CREATE TABLE "public"."groups" (
  "id" uuid NOT NULL,
  "user_id" uuid,
  "group_name" varchar(100),
  "created_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id")
);

-- ----------------------------
-- Table structure for contact
-- ----------------------------
DROP TABLE IF EXISTS "public"."contact" CASCADE;
CREATE TABLE "public"."contact" (
  "id" uuid NOT NULL,
  "group_id" uuid,
  "name" varchar(100),
  "phone" varchar(30),
  "additional_data" jsonb,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  PRIMARY KEY ("id"),
  FOREIGN KEY ("group_id") REFERENCES "public"."groups" ("id")
);

-- ----------------------------
-- Table structure for broadcast_campaigns
-- ----------------------------
DROP TABLE IF EXISTS "public"."broadcast_campaigns" CASCADE;
CREATE TABLE "public"."broadcast_campaigns" (
  "id" uuid NOT NULL,
  "user_id" uuid,
  "device_id" uuid,
  "name_broadcast" varchar(255),
  "total_recipients" int4,
  "processed_count" int4,
  "scheduled_at" timestamptz(6),
  "status" int4,
  "created_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id"),
  FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id")
);

-- ----------------------------
-- Table structure for broadcast_messages
-- ----------------------------
DROP TABLE IF EXISTS "public"."broadcast_messages" CASCADE;
CREATE TABLE "public"."broadcast_messages" (
  "id" uuid NOT NULL,
  "campaign_id" uuid,
  "message_type" int4,
  "message_text" varchar(255),
  "media_url" varchar(255),
  "button_data" jsonb,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("campaign_id") REFERENCES "public"."broadcast_campaigns" ("id")
);

-- ----------------------------
-- Table structure for broadcast_recipients
-- ----------------------------
DROP TABLE IF EXISTS "public"."broadcast_recipients" CASCADE;
CREATE TABLE "public"."broadcast_recipients" (
  "id" uuid NOT NULL,
  "campaign_id" uuid,
  "groups_id" uuid,
  "contact_id" uuid,
  "status" int4,
  "sent_at" timestamptz(6),
  "error_messages" varchar(255),
  "retry_count" int4,
  "created_at" timestamptz(6),
  PRIMARY KEY ("id"),
  FOREIGN KEY ("campaign_id") REFERENCES "public"."broadcast_campaigns" ("id"),
  FOREIGN KEY ("groups_id") REFERENCES "public"."groups" ("id"),
  FOREIGN KEY ("contact_id") REFERENCES "public"."contact" ("id")
);

-- ----------------------------
-- Table structure for auto_response
-- ----------------------------
DROP TABLE IF EXISTS "public"."auto_response" CASCADE;
CREATE TABLE "public"."auto_response" (
  "id" uuid NOT NULL,
  "device_id" uuid,
  "keyword" varchar(255),
  "response_text" varchar(255),
  "is_active" bool,
  "created_at" timestamptz(6),
  PRIMARY KEY ("id"),
  FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id")
);

-- ----------------------------
-- Table structure for warming_pool
-- ----------------------------
DROP TABLE IF EXISTS "public"."warming_pool" CASCADE;
CREATE TABLE "public"."warming_pool" (
  "id" uuid NOT NULL,
  "device_id" uuid,
  "intensity" int4,
  "daily_limit" int4,
  "message_send_today" int4,
  "is_active" bool,
  "next_action_at" timestamptz(6),
  "created_at" timestamptz(6),
  PRIMARY KEY ("id"),
  FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id")
);

-- ----------------------------
-- Table structure for warming_sessions
-- ----------------------------
DROP TABLE IF EXISTS "public"."warming_sessions" CASCADE;
CREATE TABLE "public"."warming_sessions" (
  "id" uuid NOT NULL,
  "device_id" uuid,
  "target_phone" varchar(30),
  "message_sent" varchar(255),
  "response_received" varchar(255),
  "status" int4,
  "created_at" timestamptz(6),
  PRIMARY KEY ("id"),
  FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id")
);

-- ----------------------------
-- Table structure for webhooks
-- ----------------------------
DROP TABLE IF EXISTS "public"."webhooks" CASCADE;
CREATE TABLE "public"."webhooks" (
  "id" uuid NOT NULL,
  "device_id" uuid,
  "webhook_url" varchar(255),
  "secret_key" varchar(255),
  "created_at" timestamptz(6),
  PRIMARY KEY ("id"),
  FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id")
);

-- ----------------------------
-- Table structure for api_logs
-- ----------------------------
DROP TABLE IF EXISTS "public"."api_logs" CASCADE;
CREATE TABLE "public"."api_logs" (
  "id" uuid NOT NULL,
  "user_id" uuid,
  "endpoint" varchar(255),
  "req_body" jsonb,
  "response_body" jsonb,
  "created_at" timestamptz(6),
  "ip_address" varchar(255),
  "device_id" uuid,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id"),
  FOREIGN KEY ("device_id") REFERENCES "public"."devices" ("id")
);

-- ----------------------------
-- Table structure for lookup
-- ----------------------------
DROP TABLE IF EXISTS "public"."lookup" CASCADE;
CREATE TABLE "public"."lookup" (
  "id" SERIAL PRIMARY KEY,
  "keys" varchar(255),
  "values" varchar(100)
);

-- ----------------------------
-- Table structure for system_settings
-- ----------------------------
DROP TABLE IF EXISTS "public"."system_settings" CASCADE;
CREATE TABLE "public"."system_settings" (
  "id" SERIAL PRIMARY KEY,
  "keys" varchar(100),
  "value" varchar(255),
  "description" varchar(255),
  "created_at" timestamptz(6)
);

-- All tables created successfully
