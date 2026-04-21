-- Migration: 007_align_user_sessions_schema
-- Description: Backfill missing columns for user_sessions on partially migrated databases

ALTER TABLE IF EXISTS public.user_sessions
  ADD COLUMN IF NOT EXISTS refresh_token_hash varchar(64),
  ADD COLUMN IF NOT EXISTS ip_address varchar(64),
  ADD COLUMN IF NOT EXISTS user_agent text,
  ADD COLUMN IF NOT EXISTS last_active_at timestamptz(6),
  ADD COLUMN IF NOT EXISTS refresh_expires_at timestamptz(6),
  ADD COLUMN IF NOT EXISTS revoked_at timestamptz(6),
  ADD COLUMN IF NOT EXISTS created_at timestamptz(6);

UPDATE public.user_sessions
SET
  last_active_at = COALESCE(last_active_at, created_at, NOW()),
  refresh_expires_at = COALESCE(refresh_expires_at, expires_at),
  created_at = COALESCE(created_at, NOW())
WHERE last_active_at IS NULL
   OR refresh_expires_at IS NULL
   OR created_at IS NULL;

ALTER TABLE IF EXISTS public.user_sessions
  ALTER COLUMN last_active_at SET DEFAULT NOW(),
  ALTER COLUMN refresh_expires_at SET DEFAULT NOW(),
  ALTER COLUMN created_at SET DEFAULT NOW();

ALTER TABLE IF EXISTS public.user_sessions
  ALTER COLUMN last_active_at SET NOT NULL,
  ALTER COLUMN refresh_expires_at SET NOT NULL,
  ALTER COLUMN created_at SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_user_sessions_refresh_token_hash
  ON public.user_sessions (refresh_token_hash);