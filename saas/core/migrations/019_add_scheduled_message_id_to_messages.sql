-- Migration 019: Add scheduled_message_id column to messages table
-- This allows tracking messages launched by the internal scheduler
ALTER TABLE "public"."messages" ADD COLUMN IF NOT EXISTS "scheduled_message_id" uuid NULL;

-- Add index for better filtering by scheduled job
CREATE INDEX IF NOT EXISTS idx_messages_scheduled_msg_id ON "public"."messages" ("scheduled_message_id");
