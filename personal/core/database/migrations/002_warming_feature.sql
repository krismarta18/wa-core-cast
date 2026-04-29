-- Add warming columns to devices table
ALTER TABLE "public"."devices" ADD COLUMN "is_warming" bool DEFAULT false;
ALTER TABLE "public"."devices" ADD COLUMN "warming_until" timestamptz(6);
