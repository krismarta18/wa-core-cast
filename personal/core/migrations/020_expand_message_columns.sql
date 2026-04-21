-- Migration: 020_expand_message_columns
-- Description: Expand receipt_number, content, and error_log columns to TEXT to avoid "value too long" errors.

DO $$
BEGIN
    -- Alter receipt_number to TEXT
    ALTER TABLE public.messages ALTER COLUMN receipt_number TYPE text;

    -- Alter content to TEXT
    ALTER TABLE public.messages ALTER COLUMN content TYPE text;

    -- Alter error_log to TEXT
    ALTER TABLE public.messages ALTER COLUMN error_log TYPE text;
END $$;
