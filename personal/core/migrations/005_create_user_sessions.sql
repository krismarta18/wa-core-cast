-- Migration: 005_create_user_sessions
-- Description: Create persistent auth sessions for JWT-backed logins

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS public.user_sessions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
  session_token_hash varchar(64) NOT NULL,
  refresh_token_hash varchar(64),
  ip_address varchar(64),
  user_agent text,
  last_active_at timestamptz(6) NOT NULL DEFAULT NOW(),
  expires_at timestamptz(6) NOT NULL,
  refresh_expires_at timestamptz(6) NOT NULL,
  revoked_at timestamptz(6),
  created_at timestamptz(6) NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_sessions_session_token_hash
  ON public.user_sessions (session_token_hash);

CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id
  ON public.user_sessions (user_id);

CREATE INDEX IF NOT EXISTS idx_user_sessions_expires_at
  ON public.user_sessions (expires_at);