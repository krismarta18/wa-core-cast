-- Migration: 008_align_billing_schema
-- Description: Align billing plans and subscriptions schema with the current billing model

CREATE EXTENSION IF NOT EXISTS pgcrypto;

DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'billing_plans' AND column_name = 'max_device'
  ) AND NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'billing_plans' AND column_name = 'max_devices'
  ) THEN
    ALTER TABLE public.billing_plans RENAME COLUMN max_device TO max_devices;
  END IF;

  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'billing_plans' AND column_name = 'max_messages_day'
  ) AND NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'billing_plans' AND column_name = 'max_messages_per_day'
  ) THEN
    ALTER TABLE public.billing_plans RENAME COLUMN max_messages_day TO max_messages_per_day;
  END IF;
END $$;

ALTER TABLE IF EXISTS public.billing_plans
  ADD COLUMN IF NOT EXISTS billing_cycle varchar(20),
  ADD COLUMN IF NOT EXISTS is_active boolean,
  ADD COLUMN IF NOT EXISTS created_at timestamptz(6),
  ADD COLUMN IF NOT EXISTS updated_at timestamptz(6);

UPDATE public.billing_plans
SET
  billing_cycle = COALESCE(NULLIF(billing_cycle, ''), 'monthly'),
  is_active = COALESCE(is_active, true),
  created_at = COALESCE(created_at, NOW()),
  updated_at = COALESCE(updated_at, NOW())
WHERE billing_cycle IS NULL
   OR is_active IS NULL
   OR created_at IS NULL
   OR updated_at IS NULL;

ALTER TABLE IF EXISTS public.billing_plans
  ALTER COLUMN id SET DEFAULT gen_random_uuid(),
  ALTER COLUMN name SET NOT NULL,
  ALTER COLUMN price SET NOT NULL,
  ALTER COLUMN max_devices SET NOT NULL,
  ALTER COLUMN max_messages_per_day SET NOT NULL,
  ALTER COLUMN billing_cycle SET DEFAULT 'monthly',
  ALTER COLUMN billing_cycle SET NOT NULL,
  ALTER COLUMN is_active SET DEFAULT true,
  ALTER COLUMN is_active SET NOT NULL,
  ALTER COLUMN created_at SET DEFAULT NOW(),
  ALTER COLUMN created_at SET NOT NULL,
  ALTER COLUMN updated_at SET DEFAULT NOW(),
  ALTER COLUMN updated_at SET NOT NULL;

DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'subscriptions' AND column_name = 'status'
  ) THEN
    ALTER TABLE public.subscriptions
      ALTER COLUMN status TYPE varchar(20)
      USING CASE
        WHEN status::text = '1' THEN 'active'
        WHEN status::text = '2' THEN 'expired'
        WHEN status::text = '3' THEN 'cancelled'
        WHEN status::text = '0' THEN 'inactive'
        ELSE COALESCE(NULLIF(status::text, ''), 'inactive')
      END;
  END IF;
END $$;

ALTER TABLE IF EXISTS public.subscriptions
  ADD COLUMN IF NOT EXISTS start_date timestamptz(6),
  ADD COLUMN IF NOT EXISTS end_date timestamptz(6),
  ADD COLUMN IF NOT EXISTS renewal_date timestamptz(6),
  ADD COLUMN IF NOT EXISTS auto_renew boolean,
  ADD COLUMN IF NOT EXISTS updated_at timestamptz(6);

UPDATE public.subscriptions
SET
  status = COALESCE(NULLIF(status, ''), 'inactive'),
  start_date = COALESCE(start_date, created_at, NOW()),
  renewal_date = COALESCE(renewal_date, end_date, start_date, created_at, NOW()),
  auto_renew = COALESCE(auto_renew, true),
  updated_at = COALESCE(updated_at, created_at, NOW())
WHERE status IS NULL
   OR start_date IS NULL
   OR renewal_date IS NULL
   OR auto_renew IS NULL
   OR updated_at IS NULL;

ALTER TABLE IF EXISTS public.subscriptions
  ALTER COLUMN id SET DEFAULT gen_random_uuid(),
  ALTER COLUMN status SET DEFAULT 'inactive',
  ALTER COLUMN status SET NOT NULL,
  ALTER COLUMN auto_renew SET DEFAULT true,
  ALTER COLUMN auto_renew SET NOT NULL,
  ALTER COLUMN updated_at SET DEFAULT NOW(),
  ALTER COLUMN updated_at SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_subscriptions_user_status
  ON public.subscriptions (user_id, status);

CREATE INDEX IF NOT EXISTS idx_billing_plans_active_price
  ON public.billing_plans (is_active, price);