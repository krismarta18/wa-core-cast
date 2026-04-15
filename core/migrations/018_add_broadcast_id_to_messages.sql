-- Migration 018: Add broadcast_id column to messages table
ALTER TABLE "public"."messages" ADD COLUMN IF NOT EXISTS "broadcast_id" uuid NULL;
