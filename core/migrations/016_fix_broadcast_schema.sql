-- Migration 016: Robust Broadcast Schema Alignment
-- This ensures ALL required columns exist in broadcast_campaigns and broadcast_recipients

-- Ensure broadcast_campaigns columns
ALTER TABLE "public"."broadcast_campaigns" ADD COLUMN IF NOT EXISTS "user_id" uuid;
ALTER TABLE "public"."broadcast_campaigns" ADD COLUMN IF NOT EXISTS "device_id" uuid;
ALTER TABLE "public"."broadcast_campaigns" ADD COLUMN IF NOT EXISTS "template_id" uuid;
ALTER TABLE "public"."broadcast_campaigns" ADD COLUMN IF NOT EXISTS "name" varchar(255);
ALTER TABLE "public"."broadcast_campaigns" ADD COLUMN IF NOT EXISTS "message_content" text;
ALTER TABLE "public"."broadcast_campaigns" ADD COLUMN IF NOT EXISTS "delay_seconds" int NOT NULL DEFAULT 5;
ALTER TABLE "public"."broadcast_campaigns" ADD COLUMN IF NOT EXISTS "total_recipients" int NOT NULL DEFAULT 0;
ALTER TABLE "public"."broadcast_campaigns" ADD COLUMN IF NOT EXISTS "processed_count" int NOT NULL DEFAULT 0;
ALTER TABLE "public"."broadcast_campaigns" ADD COLUMN IF NOT EXISTS "success_count" int NOT NULL DEFAULT 0;
ALTER TABLE "public"."broadcast_campaigns" ADD COLUMN IF NOT EXISTS "failed_count" int NOT NULL DEFAULT 0;
ALTER TABLE "public"."broadcast_campaigns" ADD COLUMN IF NOT EXISTS "scheduled_at" timestamptz;
ALTER TABLE "public"."broadcast_campaigns" ADD COLUMN IF NOT EXISTS "started_at" timestamptz;
ALTER TABLE "public"."broadcast_campaigns" ADD COLUMN IF NOT EXISTS "completed_at" timestamptz;
ALTER TABLE "public"."broadcast_campaigns" ADD COLUMN IF NOT EXISTS "status" varchar(20) NOT NULL DEFAULT 'draft';
ALTER TABLE "public"."broadcast_campaigns" ADD COLUMN IF NOT EXISTS "created_at" timestamptz NOT NULL DEFAULT now();
ALTER TABLE "public"."broadcast_campaigns" ADD COLUMN IF NOT EXISTS "updated_at" timestamptz NOT NULL DEFAULT now();
ALTER TABLE "public"."broadcast_campaigns" ADD COLUMN IF NOT EXISTS "deleted_at" timestamptz;

-- Ensure broadcast_recipients columns
ALTER TABLE "public"."broadcast_recipients" ADD COLUMN IF NOT EXISTS "campaign_id" uuid;
ALTER TABLE "public"."broadcast_recipients" ADD COLUMN IF NOT EXISTS "contact_id" uuid;
ALTER TABLE "public"."broadcast_recipients" ADD COLUMN IF NOT EXISTS "group_id" uuid;
ALTER TABLE "public"."broadcast_recipients" ADD COLUMN IF NOT EXISTS "phone_number" varchar(30);
ALTER TABLE "public"."broadcast_recipients" ADD COLUMN IF NOT EXISTS "status" varchar(20) NOT NULL DEFAULT 'pending';
ALTER TABLE "public"."broadcast_recipients" ADD COLUMN IF NOT EXISTS "sent_at" timestamptz;
ALTER TABLE "public"."broadcast_recipients" ADD COLUMN IF NOT EXISTS "failed_at" timestamptz;
ALTER TABLE "public"."broadcast_recipients" ADD COLUMN IF NOT EXISTS "error_message" text;
ALTER TABLE "public"."broadcast_recipients" ADD COLUMN IF NOT EXISTS "retry_count" int NOT NULL DEFAULT 0;
ALTER TABLE "public"."broadcast_recipients" ADD COLUMN IF NOT EXISTS "created_at" timestamptz NOT NULL DEFAULT now();
