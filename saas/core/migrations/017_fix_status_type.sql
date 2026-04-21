-- Migration 017: Fix Status Column Types
-- Some existing tables might have 'status' as an integer from previous versions
-- We ensure they are varchar(20) as expected by the Go models

-- Fix broadcast_campaigns
DO $$
BEGIN
    -- Check if status is integer
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='broadcast_campaigns' AND column_name='status' AND data_type='integer'
    ) THEN
        ALTER TABLE broadcast_campaigns ALTER COLUMN status TYPE varchar(20) USING 'draft';
    END IF;
END
$$;

-- Ensure it is varchar(20) regardless
ALTER TABLE "public"."broadcast_campaigns" ALTER COLUMN "status" TYPE varchar(20);
ALTER TABLE "public"."broadcast_campaigns" ALTER COLUMN "status" SET DEFAULT 'draft';

-- Fix broadcast_recipients
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='broadcast_recipients' AND column_name='status' AND data_type='integer'
    ) THEN
        ALTER TABLE broadcast_recipients ALTER COLUMN status TYPE varchar(20) USING 'pending';
    END IF;
END
$$;

ALTER TABLE "public"."broadcast_recipients" ALTER COLUMN "status" TYPE varchar(20);
ALTER TABLE "public"."broadcast_recipients" ALTER COLUMN "status" SET DEFAULT 'pending';
