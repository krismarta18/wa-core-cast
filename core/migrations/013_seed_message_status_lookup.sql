-- Migration: 013_seed_message_status_lookup
-- Description: Seed master data for message status_message values into lookup table
--              keys format: 'message_status.<int_value>'
--              values: label teks status

INSERT INTO public.lookup (keys, values)
SELECT 'message_status.0', 'pending'
WHERE NOT EXISTS (SELECT 1 FROM public.lookup WHERE keys = 'message_status.0');

INSERT INTO public.lookup (keys, values)
SELECT 'message_status.1', 'sent'
WHERE NOT EXISTS (SELECT 1 FROM public.lookup WHERE keys = 'message_status.1');

INSERT INTO public.lookup (keys, values)
SELECT 'message_status.2', 'delivered'
WHERE NOT EXISTS (SELECT 1 FROM public.lookup WHERE keys = 'message_status.2');

INSERT INTO public.lookup (keys, values)
SELECT 'message_status.3', 'read'
WHERE NOT EXISTS (SELECT 1 FROM public.lookup WHERE keys = 'message_status.3');

INSERT INTO public.lookup (keys, values)
SELECT 'message_status.4', 'failed'
WHERE NOT EXISTS (SELECT 1 FROM public.lookup WHERE keys = 'message_status.4');

-- Seed juga untuk direction (arah pesan)
INSERT INTO public.lookup (keys, values)
SELECT 'message_direction.OUT', 'outgoing'
WHERE NOT EXISTS (SELECT 1 FROM public.lookup WHERE keys = 'message_direction.OUT');

INSERT INTO public.lookup (keys, values)
SELECT 'message_direction.IN', 'incoming'
WHERE NOT EXISTS (SELECT 1 FROM public.lookup WHERE keys = 'message_direction.IN');

-- Seed untuk message_type
INSERT INTO public.lookup (keys, values)
SELECT 'message_type.0', 'text'
WHERE NOT EXISTS (SELECT 1 FROM public.lookup WHERE keys = 'message_type.0');

INSERT INTO public.lookup (keys, values)
SELECT 'message_type.1', 'image'
WHERE NOT EXISTS (SELECT 1 FROM public.lookup WHERE keys = 'message_type.1');

INSERT INTO public.lookup (keys, values)
SELECT 'message_type.2', 'document'
WHERE NOT EXISTS (SELECT 1 FROM public.lookup WHERE keys = 'message_type.2');

INSERT INTO public.lookup (keys, values)
SELECT 'message_type.3', 'audio'
WHERE NOT EXISTS (SELECT 1 FROM public.lookup WHERE keys = 'message_type.3');

INSERT INTO public.lookup (keys, values)
SELECT 'message_type.4', 'video'
WHERE NOT EXISTS (SELECT 1 FROM public.lookup WHERE keys = 'message_type.4');
