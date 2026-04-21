-- Migration: 004_store_plaintext_otp
-- Description: Rename OTP storage column to plaintext otp_code

DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'otp_verifications' AND column_name = 'otp_code_hash'
  ) AND NOT EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'otp_verifications' AND column_name = 'otp_code'
  ) THEN
    ALTER TABLE public.otp_verifications RENAME COLUMN otp_code_hash TO otp_code;
  END IF;
END $$;

DELETE FROM public.otp_verifications;

ALTER TABLE IF EXISTS public.otp_verifications
  ALTER COLUMN otp_code TYPE varchar(10);