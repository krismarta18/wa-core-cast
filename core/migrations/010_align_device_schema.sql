-- Migration: 010_align_device_schema
-- Description: Align devices table with the current service schema

DO $$
BEGIN
  -- Rename name_device to display_name
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'devices' AND column_name = 'name_device'
  ) AND NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'devices' AND column_name = 'display_name'
  ) THEN
    ALTER TABLE public.devices RENAME COLUMN name_device TO display_name;
  END IF;

  -- Rename phone to phone_number
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'devices' AND column_name = 'phone'
  ) AND NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'devices' AND column_name = 'phone_number'
  ) THEN
    ALTER TABLE public.devices RENAME COLUMN phone TO phone_number;
  END IF;

  -- Rename last_seen to last_seen_at
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'devices' AND column_name = 'last_seen'
  ) AND NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'devices' AND column_name = 'last_seen_at'
  ) THEN
    ALTER TABLE public.devices RENAME COLUMN last_seen TO last_seen_at;
  END IF;
END $$;

-- Convert status to varchar(20)
ALTER TABLE IF EXISTS public.devices
  ALTER COLUMN status TYPE varchar(20) USING CASE
    WHEN status = 1 THEN 'connected'
    WHEN status = 0 THEN 'disconnected'
    WHEN status = 2 THEN 'pending_qr'
    ELSE 'pending_qr'
  END;

-- Add new columns that are expected by device_queries.go
ALTER TABLE IF EXISTS public.devices
  ADD COLUMN IF NOT EXISTS connected_since timestamptz(6),
  ADD COLUMN IF NOT EXISTS platform varchar(255),
  ADD COLUMN IF NOT EXISTS wa_version varchar(50),
  ADD COLUMN IF NOT EXISTS battery_level int4;
