-- Migration to add limit columns to subscriptions table
-- This freezing of limits ensures users keep their plan's terms even if the master plan is updated later.

ALTER TABLE subscriptions ADD COLUMN IF NOT EXISTS max_devices INTEGER DEFAULT 0;
ALTER TABLE subscriptions ADD COLUMN IF NOT EXISTS max_messages_per_day INTEGER DEFAULT 0;

-- Update existing subscriptions from their respective billing plans
UPDATE subscriptions s
SET 
    max_devices = p.max_devices,
    max_messages_per_day = p.max_messages_per_day
FROM billing_plans p
WHERE s.plan_id = p.id;
