// ─── Session / Device ──────────────────────────────────────────────────────────

export type DeviceStatus = 0 | 1 | 2; // 0=inactive, 1=active, 2=pending/QR

export interface Device {
  device_id: string;
  status: DeviceStatus;
  is_active: boolean;
}

export interface SessionListResponse {
  count: number;
  sessions: Device[];
}

export interface InitiateSessionRequest {
  device_id: string;
  user_id: string;
  phone: string;
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
  direction: MessageDirection;
  status: MessageStatus;
  to_jid: string;
  content: string;
  created_at: string;
  updated_at: string;
}

export interface SendMessageRequest {
  to: string;
  content: string;
}

export interface SendScheduledMessageRequest extends SendMessageRequest {
  scheduled_at: string; // ISO 8601
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

// ─── Generic API ─────────────────────────────────────────────────────────────

export interface ApiError {
  error: string;
}
