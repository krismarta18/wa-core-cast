-- Migration: 003_auth_schema_alignment
-- Description: Align auth-related tables with the current service schema

CREATE EXTENSION IF NOT EXISTS pgcrypto;

DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'users' AND column_name = 'phone'
  ) AND NOT EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'users' AND column_name = 'phone_number'
  ) THEN
    ALTER TABLE public.users RENAME COLUMN phone TO phone_number;
  END IF;

  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'users' AND column_name = 'nama_lengkap'
  ) AND NOT EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'users' AND column_name = 'full_name'
  ) THEN
    ALTER TABLE public.users RENAME COLUMN nama_lengkap TO full_name;
  END IF;

  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'users' AND column_name = 'is_verify'
  ) AND NOT EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'users' AND column_name = 'is_verified'
  ) THEN
    ALTER TABLE public.users RENAME COLUMN is_verify TO is_verified;
  END IF;

  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'users' AND column_name = 'is_ban'
  ) AND NOT EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'users' AND column_name = 'is_banned'
  ) THEN
    ALTER TABLE public.users RENAME COLUMN is_ban TO is_banned;
  END IF;

  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'users' AND column_name = 'is_api'
  ) AND NOT EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'users' AND column_name = 'is_api_enabled'
  ) THEN
    ALTER TABLE public.users RENAME COLUMN is_api TO is_api_enabled;
  END IF;
END $$;

ALTER TABLE IF EXISTS public.users
  ADD COLUMN IF NOT EXISTS email varchar(255),
  ADD COLUMN IF NOT EXISTS company_name varchar(255),
  ADD COLUMN IF NOT EXISTS timezone varchar(100),
  ADD COLUMN IF NOT EXISTS created_at timestamptz(6),
  ADD COLUMN IF NOT EXISTS updated_at timestamptz(6),
  ADD COLUMN IF NOT EXISTS last_login_at timestamptz(6);

ALTER TABLE IF EXISTS public.users
  ALTER COLUMN id SET DEFAULT gen_random_uuid(),
  ALTER COLUMN phone_number SET NOT NULL,
  ALTER COLUMN full_name SET NOT NULL,
  ALTER COLUMN is_verified SET DEFAULT false,
  ALTER COLUMN is_banned SET DEFAULT false,
  ALTER COLUMN is_api_enabled SET DEFAULT false,
  ALTER COLUMN timezone SET DEFAULT 'Asia/Jakarta',
  ALTER COLUMN created_at SET DEFAULT NOW(),
  ALTER COLUMN updated_at SET DEFAULT NOW();

UPDATE public.users
SET
  phone_number = COALESCE(phone_number, ''),
  full_name = COALESCE(full_name, ''),
  is_verified = COALESCE(is_verified, false),
  is_banned = COALESCE(is_banned, false),
  is_api_enabled = COALESCE(is_api_enabled, false),
  timezone = COALESCE(NULLIF(timezone, ''), 'Asia/Jakarta'),
  created_at = COALESCE(created_at, NOW()),
  updated_at = COALESCE(updated_at, NOW());

ALTER TABLE IF EXISTS public.users
  ALTER COLUMN phone_number SET NOT NULL,
  ALTER COLUMN full_name SET NOT NULL,
  ALTER COLUMN is_verified SET NOT NULL,
  ALTER COLUMN is_banned SET NOT NULL,
  ALTER COLUMN is_api_enabled SET NOT NULL,
  ALTER COLUMN timezone SET NOT NULL,
  ALTER COLUMN created_at SET NOT NULL,
  ALTER COLUMN updated_at SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_phone_number ON public.users (phone_number);

CREATE TABLE IF NOT EXISTS public.otp_verifications (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NULL REFERENCES public.users(id) ON DELETE CASCADE,
  phone_number varchar(30) NOT NULL,
  context varchar(20) NOT NULL,
  otp_code varchar(10) NOT NULL,
  attempt_count int4 NOT NULL DEFAULT 0,
  expires_at timestamptz(6) NOT NULL,
  verified_at timestamptz(6),
  created_at timestamptz(6) NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_otp_verifications_phone_number
  ON public.otp_verifications (phone_number);

CREATE INDEX IF NOT EXISTS idx_otp_verifications_user_id
  ON public.otp_verifications (user_id);

CREATE INDEX IF NOT EXISTS idx_otp_verifications_expires_at
  ON public.otp_verifications (expires_at);