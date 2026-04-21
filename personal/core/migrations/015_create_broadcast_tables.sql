-- Migration 015: Create Broadcast Tables
-- This ensures broadcast_campaigns and broadcast_recipients exist with correct columns

CREATE TABLE IF NOT EXISTS "public"."broadcast_campaigns" (
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
  "status"            varchar(20)     NOT NULL DEFAULT 'draft',
  "created_at"        timestamptz     NOT NULL DEFAULT now(),
  "updated_at"        timestamptz     NOT NULL DEFAULT now(),
  "deleted_at"        timestamptz     NULL,
  CONSTRAINT "broadcast_campaigns_pkey" PRIMARY KEY ("id")
);

-- For existing tables missing template_id, we add it safely
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='broadcast_campaigns' AND column_name='template_id') THEN
        ALTER TABLE broadcast_campaigns ADD COLUMN template_id uuid NULL;
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS "public"."broadcast_recipients" (
  "id"            uuid            NOT NULL DEFAULT gen_random_uuid(),
  "campaign_id"   uuid            NOT NULL,
  "contact_id"    uuid            NULL,
  "group_id"      uuid            NULL,
  "phone_number"  varchar(30)     NOT NULL,
  "status"        varchar(20)     NOT NULL DEFAULT 'pending',
  "sent_at"       timestamptz     NULL,
  "failed_at"     timestamptz     NULL,
  "error_message" text            NULL,
  "retry_count"   int             NOT NULL DEFAULT 0,
  "created_at"    timestamptz     NOT NULL DEFAULT now(),
  CONSTRAINT "broadcast_recipients_pkey" PRIMARY KEY ("id")
);
