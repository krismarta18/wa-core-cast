// ─── Session / Device ──────────────────────────────────────────────────────────

export type DeviceStatus = 0 | 1 | 2; // 0=inactive, 1=active, 2=pending/QR

export interface Device {
  device_id: string;
  status: DeviceStatus;
  is_active: boolean;
  phone?: string;
  display_name?: string;
}

export interface SessionListResponse {
  count: number;
  sessions: Device[];
}

export interface InitiateSessionRequest {
  device_id: string;
  user_id: string;
  phone: string;
  display_name?: string;
}

export interface QRCodeResponse {
  device_id: string;
  qr_code_string: string;
  qr_code_image: {
    base64_png: string;
    png_bytes: number;
    url_format: string;
    direct_url: string;
  };
  status: DeviceStatus;
  message: string;
}

// ─── Messages ─────────────────────────────────────────────────────────────────

export type MessageStatus = "pending" | "sent" | "delivered" | "read" | "failed";
export type MessageDirection = "in" | "out";

export interface Message {
  id: string;
  device_id: string;
  target_jid: string;
  content: string;
  content_type: string;
  status: string;
  scheduled_for?: string;
  sent_at?: string;
  delivered_at?: string;
  read_at?: string;
  failed_at?: string;
  error_log?: string;
  broadcast_id?: string;
  scheduled_message_id?: string;
  created_at: string;
  updated_at: string;
}

export interface SendMessageRequest {
  target_jid: string;
  content: string;
  group_id?: string;
  priority?: number;
}

export interface SendScheduledMessageRequest extends SendMessageRequest {
  scheduled_for: string; // ISO 8601
  media_url?: string;
  content_type?: string;
  caption?: string;
}

export interface QueueStats {
  total: number;
  pending: number;
  processing: number;
  failed: number;
}

// ─── Server Info ──────────────────────────────────────────────────────────────

export interface ServerStatus {
  status: string;
  version: string;
  environment: string;
  uptime_seconds: number;
  active_sessions: number;
  total_sessions: number;
}

export interface ServerStats {
  db_open_connections: number;
  db_idle_connections: number;
  sessions_count: number;
}

// ─── Health ───────────────────────────────────────────────────────────────────

export interface HealthResponse {
  status: string;
  database: boolean;
  uptime?: number;
}

// ─── Auth ─────────────────────────────────────────────────────────────────────

export interface AuthUser {
  id: string;
  phone_number: string;
  full_name: string;
  email?: string | null;
  company_name?: string | null;
  timezone: string;
  is_verified: boolean;
  is_banned: boolean;
  is_api_enabled: boolean;
  created_at: string;
  last_login_at?: string | null;
}

export interface AuthRegisterRequest {
  phone_number: string;
  full_name: string;
}

export interface AuthOTPRequest {
  phone_number: string;
}

export interface AuthVerifyOTPRequest {
  phone_number: string;
  otp_code: string;
}

export interface AuthRefreshTokenRequest {
  refresh_token: string;
}

export interface AuthGenericSuccessResponse {
  success: boolean;
  message?: string;
  phone_number?: string;
}

export interface AuthSessionResponse {
  success: boolean;
  message?: string;
  access_token: string;
  refresh_token: string;
  token_type: string;
  expires_in: number;
  refresh_expires_in: number;
  user: AuthUser;
}

export interface AuthMeResponse {
  success: boolean;
  user: AuthUser;
}

export interface AuthUpdateProfileRequest {
  full_name?: string;
  email?: string;
  company_name?: string;
  timezone?: string;
}

// ─── Billing ──────────────────────────────────────────────────────────────────

export interface BillingUsagePoint {
  date: string;
  sent: number;
  failed: number;
}

export interface BillingPlanSummary {
  id: string;
  name: string;
  price: number;
  quota_limit: number;
  device_max: number;
  current: boolean;
  is_active: boolean;
}

export interface BillingInvoiceSummary {
  id: string;
  subscription_id: string;
  date: string;
  plan_name: string;
  amount: number;
  status: string;
}

export interface BillingCurrentPlan {
  subscription_id: string;
  plan_id: string;
  name: string;
  price: number;
  billing_cycle: string;
  renewal_date?: string | null;
  quota_used: number;
  quota_limit: number;
  device_used: number;
  device_max: number;
  auto_renew: boolean;
  status: string;
  features?: string[] | Record<string, unknown> | null;
}

export interface BillingOverview {
  current_plan?: BillingCurrentPlan | null;
  usage_history: BillingUsagePoint[];
  plans: BillingPlanSummary[];
  invoices: BillingInvoiceSummary[];
}

export interface BillingOverviewResponse {
  success: boolean;
  billing: BillingOverview;
}

export interface BillingCheckoutRequest {
  plan_id: string;
}

export interface BillingCheckoutData {
  subscription: BillingCurrentPlan;
  invoice: BillingInvoiceSummary;
  payment_status: string;
  payment_method: string;
}

export interface BillingCheckoutResponse {
  success: boolean;
  message?: string;
  checkout: BillingCheckoutData;
}

// ─── Generic API ─────────────────────────────────────────────────────────────

export interface ApiError {
  error: string;
}

// ─── Contacts ─────────────────────────────────────────────────────────────────

export interface Contact {
  id: string;
  user_id: string;
  label_id?: string;
  name: string;
  phone_number: string;
  note?: string;
  additional_data?: any;
  created_at: string;
}

export interface CreateContactRequest {
  name: string;
  phone_number: string;
  label_id?: string;
  note?: string;
  additional_data?: Record<string, any>;
}

export interface UpdateContactRequest {
  name?: string;
  phone_number?: string;
  label_id?: string;
  note?: string;
  additional_data?: Record<string, any>;
}

export interface ContactGroup {
  id: string;
  user_id: string;
  name: string;
  description?: string;
  created_at: string;
}

export interface CreateContactGroupRequest {
  name: string;
  description?: string;
}

export interface BlacklistEntry {
  id: string;
  phone_number: string;
  reason: string;
  blocked_at: string;
}

// ─── Analytics ────────────────────────────────────────────────────────────────

export interface DailyStatPoint {
  id: string;
  user_id: string;
  device_id: string;
  stat_date: string;
  sent_count: number;
  failed_count: number;
  delivered_count: number;
  received_count: number;
  success_rate: number;
  created_at: string;
}

export interface DeviceAnalytics {
  name: string;
  sent: number;
  success: number;
}

export interface UsageStatsResponse {
  total_sent: number;
  total_failed: number;
  success_rate: number;
  daily: DailyStatPoint[];
  device_stats: DeviceAnalytics[];
}

export interface FailureReasonStat {
  reason: string;
  count: number;
  pct: number;
}

export interface FailureLogItem {
  id: string;
  to: string;
  device: string;
  reason: string;
  time: string;
  type: string;
}

export interface FailureRateResponse {
  total_failed: number;
  failure_rate: number;
  avg_retry_time: string;
  reason_stats: FailureReasonStat[];
  latest_logs: FailureLogItem[];
}
// ─── Broadcasts ───────────────────────────────────────────────────────────────

export type BroadcastStatus = "draft" | "queued" | "sending" | "completed" | "failed" | "cancelled";

export interface BroadcastCampaign {
  id: string;
  user_id: string;
  device_id: string;
  template_id?: string;
  name: string;
  message_content?: string;
  total_recipients: number;
  success_count: number;
  failed_count: number;
  delay_seconds: number;
  scheduled_at?: string;
  started_at?: string;
  completed_at?: string;
  status: BroadcastStatus;
  created_at: string;
  updated_at: string;
}

export interface BroadcastRecipient {
  id: string;
  campaign_id: string;
  group_id?: string;
  contact_id?: string;
  phone_number: string;
  status: "pending" | "sent" | "failed";
  sent_at?: string;
  failed_at?: string;
  error_message?: string;
  retry_count: number;
  created_at: string;
}

export interface CreateBroadcastRequest {
  device_id: string;
  name: string;
  template_id?: string;
  message_content?: string;
  delay_seconds: number;
  scheduled_at?: string;
  recipients: string[];
}

export interface BroadcastCampaignResponse {
  id: string;
  name: string;
  total_recipients: number;
  success_count: number;
  failed_count: number;
  status: BroadcastStatus;
  scheduled_at?: string;
  started_at?: string;
  completed_at?: string;
  created_at: string;
}

 
 
// --- Auto Response ------------------------------------------------------------

export interface AutoResponseKeyword {
  id: string;
  user_id: string;
  device_id?: string;
  keyword: string;
  response_text: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateKeywordRequest {
  device_id?: string;
  keyword: string;
  response_text: string;
}

export interface UpdateKeywordRequest {
  keyword?: string;
  response_text?: string;
  is_active?: boolean;
}

// --- Message Templates ----------------------------------------------------------

export interface MessageTemplate {
  id: string;
  user_id: string;
  name: string;
  category: string;
  content: string;
  used_count: number;
  created_at: string;
  updated_at: string;
}

export interface CreateTemplateRequest {
  name: string;
  category: string;
  content: string;
}

export interface UpdateTemplateRequest {
  name?: string;
  category?: string;
  content?: string;
}