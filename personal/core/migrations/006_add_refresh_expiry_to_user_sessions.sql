-- Migration: 006_add_refresh_expiry_to_user_sessions
-- Description: Add refresh token expiry tracking to persisted user sessions

ALTER TABLE IF EXISTS public.user_sessions
  ADD COLUMN IF NOT EXISTS refresh_expires_at timestamptz(6);

UPDATE public.user_sessions
SET refresh_expires_at = COALESCE(refresh_expires_at, expires_at)
WHERE refresh_expires_at IS NULL;

ALTER TABLE IF EXISTS public.user_sessions
  ALTER COLUMN refresh_expires_at SET NOT NULL;
