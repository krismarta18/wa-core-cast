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

// ─── Generic API ─────────────────────────────────────────────────────────────

export interface ApiError {
  error: string;
}
