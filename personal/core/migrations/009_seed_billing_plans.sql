-- Migration: 009_seed_billing_plans
-- Description: Seed default billing plans for dashboard billing overview

INSERT INTO public.billing_plans (
  id,
  name,
  price,
  billing_cycle,
  max_devices,
  max_messages_per_day,
  features,
  is_active,
  created_at,
  updated_at
)
SELECT
  gen_random_uuid(),
  seed.name,
  seed.price,
  seed.billing_cycle,
  seed.max_devices,
  seed.max_messages_per_day,
  seed.features::jsonb,
  true,
  NOW(),
  NOW()
FROM (
  VALUES
    ('Starter', 99000::numeric, 'monthly', 2, 2000, '["2 device","2.000 pesan per hari","Basic analytics"]'),
    ('Business Pro', 299000::numeric, 'monthly', 10, 10000, '["10 device","10.000 pesan per hari","Priority support"]'),
    ('Enterprise', 799000::numeric, 'monthly', 100, 100000, '["100 device","100.000 pesan per hari","Dedicated support"]')
) AS seed(name, price, billing_cycle, max_devices, max_messages_per_day, features)
WHERE NOT EXISTS (
  SELECT 1 FROM public.billing_plans existing WHERE LOWER(existing.name) = LOWER(seed.name)
);