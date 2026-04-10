-- Migration: 002_add_indexes
-- Description: Add performance indexes

-- Indexes for users table
CREATE INDEX IF NOT EXISTS idx_users_phone ON "public"."users" ("phone");
CREATE INDEX IF NOT EXISTS idx_users_is_ban ON "public"."users" ("is_ban");
CREATE INDEX IF NOT EXISTS idx_users_is_verify ON "public"."users" ("is_verify");

-- Indexes for devices table
CREATE INDEX IF NOT EXISTS idx_devices_user_id ON "public"."devices" ("user_id");
CREATE INDEX IF NOT EXISTS idx_devices_phone ON "public"."devices" ("phone");
CREATE INDEX IF NOT EXISTS idx_devices_status ON "public"."devices" ("status");
CREATE INDEX IF NOT EXISTS idx_devices_last_seen ON "public"."devices" ("last_seen");

-- Indexes for messages table
CREATE INDEX IF NOT EXISTS idx_messages_device_id ON "public"."messages" ("device_id");
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON "public"."messages" ("created_at");
CREATE INDEX IF NOT EXISTS idx_messages_status ON "public"."messages" ("status_message");
CREATE INDEX IF NOT EXISTS idx_messages_direction ON "public"."messages" ("direction");
CREATE INDEX IF NOT EXISTS idx_messages_receipt ON "public"."messages" ("receipt_number");

-- Indexes for groups table
CREATE INDEX IF NOT EXISTS idx_groups_user_id ON "public"."groups" ("user_id");
CREATE INDEX IF NOT EXISTS idx_groups_deleted_at ON "public"."groups" ("deleted_at");

-- Indexes for contacts table
CREATE INDEX IF NOT EXISTS idx_contact_group_id ON "public"."contact" ("group_id");
CREATE INDEX IF NOT EXISTS idx_contact_phone ON "public"."contact" ("phone");
CREATE INDEX IF NOT EXISTS idx_contact_deleted_at ON "public"."contact" ("deleted_at");

-- Indexes for broadcast_campaigns table
CREATE INDEX IF NOT EXISTS idx_broadcast_campaigns_user_id ON "public"."broadcast_campaigns" ("user_id");
CREATE INDEX IF NOT EXISTS idx_broadcast_campaigns_device_id ON "public"."broadcast_campaigns" ("device_id");
CREATE INDEX IF NOT EXISTS idx_broadcast_campaigns_status ON "public"."broadcast_campaigns" ("status");
CREATE INDEX IF NOT EXISTS idx_broadcast_campaigns_deleted_at ON "public"."broadcast_campaigns" ("deleted_at");

-- Indexes for broadcast_recipients table
CREATE INDEX IF NOT EXISTS idx_broadcast_recipients_campaign_id ON "public"."broadcast_recipients" ("campaign_id");
CREATE INDEX IF NOT EXISTS idx_broadcast_recipients_status ON "public"."broadcast_recipients" ("status");

-- Indexes for auto_response table
CREATE INDEX IF NOT EXISTS idx_auto_response_device_id ON "public"."auto_response" ("device_id");
CREATE INDEX IF NOT EXISTS idx_auto_response_is_active ON "public"."auto_response" ("is_active");

-- Indexes for warming_pool table
CREATE INDEX IF NOT EXISTS idx_warming_pool_device_id ON "public"."warming_pool" ("device_id");
CREATE INDEX IF NOT EXISTS idx_warming_pool_is_active ON "public"."warming_pool" ("is_active");
CREATE INDEX IF NOT EXISTS idx_warming_pool_next_action ON "public"."warming_pool" ("next_action_at");

-- Indexes for warming_sessions table
CREATE INDEX IF NOT EXISTS idx_warming_sessions_device_id ON "public"."warming_sessions" ("device_id");
CREATE INDEX IF NOT EXISTS idx_warming_sessions_status ON "public"."warming_sessions" ("status");

-- Indexes for webhooks table
CREATE INDEX IF NOT EXISTS idx_webhooks_device_id ON "public"."webhooks" ("device_id");

-- Indexes for api_logs table
CREATE INDEX IF NOT EXISTS idx_api_logs_user_id ON "public"."api_logs" ("user_id");
CREATE INDEX IF NOT EXISTS idx_api_logs_created_at ON "public"."api_logs" ("created_at");
CREATE INDEX IF NOT EXISTS idx_api_logs_endpoint ON "public"."api_logs" ("endpoint");

-- Indexes for subscriptions table
CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON "public"."subscriptions" ("user_id");
CREATE INDEX IF NOT EXISTS idx_subscriptions_plan_id ON "public"."subscriptions" ("plan_id");
CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON "public"."subscriptions" ("status");

-- All indexes created successfully
