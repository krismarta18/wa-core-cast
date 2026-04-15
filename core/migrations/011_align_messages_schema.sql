-- Migration: 011_align_messages_schema
-- Description: Align messages table to support outgoing queue and incoming message storage

DO $$
BEGIN
  -- Add recipient_jid column (renamed from receipt_number for clarity)
  -- Keep receipt_number as is since it exists

  -- Add target_jid column for outgoing messages if not exists
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'messages' AND column_name = 'target_jid'
  ) THEN
    ALTER TABLE public.messages ADD COLUMN target_jid varchar(100);
  END IF;

  -- Add updated_at column if not exists
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'messages' AND column_name = 'updated_at'
  ) THEN
    ALTER TABLE public.messages ADD COLUMN updated_at timestamptz DEFAULT now();
  END IF;

  -- Add priority column if not exists
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'messages' AND column_name = 'priority'
  ) THEN
    ALTER TABLE public.messages ADD COLUMN priority int4 DEFAULT 3;
  END IF;

  -- Add retry_count column if not exists
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'messages' AND column_name = 'retry_count'
  ) THEN
    ALTER TABLE public.messages ADD COLUMN retry_count int4 DEFAULT 0;
  END IF;

  -- Add max_retries column if not exists
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'messages' AND column_name = 'max_retries'
  ) THEN
    ALTER TABLE public.messages ADD COLUMN max_retries int4 DEFAULT 3;
  END IF;

  -- Add scheduled_for column if not exists
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'messages' AND column_name = 'scheduled_for'
  ) THEN
    ALTER TABLE public.messages ADD COLUMN scheduled_for timestamptz;
  END IF;

END $$;

-- Add index for queue processing
CREATE INDEX IF NOT EXISTS idx_messages_device_status_dir
  ON public.messages (device_id, status_message, direction);

CREATE INDEX IF NOT EXISTS idx_messages_created_at
  ON public.messages (created_at);
