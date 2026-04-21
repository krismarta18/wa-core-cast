-- Migration: 012_add_whatsapp_message_id
-- Description: Add whatsapp_message_id column to messages for receipt tracking

DO $$
BEGIN
  -- Add whatsapp_message_id to track the WA-assigned message ID for receipt matching
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'messages' AND column_name = 'whatsapp_message_id'
  ) THEN
    ALTER TABLE public.messages ADD COLUMN whatsapp_message_id varchar(255);
  END IF;
END $$;

-- Index for fast receipt lookup
CREATE INDEX IF NOT EXISTS idx_messages_whatsapp_message_id
  ON public.messages (whatsapp_message_id)
  WHERE whatsapp_message_id IS NOT NULL;
