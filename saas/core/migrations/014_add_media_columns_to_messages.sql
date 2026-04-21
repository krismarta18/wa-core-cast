-- Migration: 014_add_media_columns_to_messages
-- Description: Add media_url and caption columns to messages table for media messaging

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'messages' AND column_name = 'media_url'
  ) THEN
    ALTER TABLE public.messages ADD COLUMN media_url text;
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'messages' AND column_name = 'caption'
  ) THEN
    ALTER TABLE public.messages ADD COLUMN caption text;
  END IF;
END $$;
