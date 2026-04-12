import axios from "axios";
import type {
  Device,
  SessionListResponse,
  InitiateSessionRequest,
  QRCodeResponse,
  SendMessageRequest,
  SendScheduledMessageRequest,
  Message,
  QueueStats,
  ServerStatus,
  ServerStats,
  HealthResponse,
} from "./types";

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080",
  headers: { "Content-Type": "application/json" },
  timeout: 10_000,
});

// ─── Sessions ─────────────────────────────────────────────────────────────────

export const sessionsApi = {
  list: () =>
    api.get<SessionListResponse>("/api/v1/sessions").then((r) => r.data),

  get: (deviceId: string) =>
    api.get<Device>(`/api/v1/sessions/${deviceId}`).then((r) => r.data),

  qr: (deviceId: string) =>
    api
      .get<QRCodeResponse>(`/api/v1/sessions/${deviceId}/qr`)
      .then((r) => r.data),

  initiate: (body: InitiateSessionRequest) =>
    api.post<{ message: string; device_id: string }>(
      "/api/v1/sessions/initiate",
      body
    ).then((r) => r.data),

  stop: (deviceId: string) =>
    api
      .post<{ message: string }>(`/api/v1/sessions/${deviceId}/stop`)
      .then((r) => r.data),
};

// ─── Messages ─────────────────────────────────────────────────────────────────

export const messagesApi = {
  send: (deviceId: string, body: SendMessageRequest) =>
    api
      .post<Message>(`/api/v1/devices/${deviceId}/messages`, body)
      .then((r) => r.data),

  sendMedia: (deviceId: string, formData: FormData) =>
    api
      .post<Message>(`/api/v1/devices/${deviceId}/messages/media`, formData, {
        headers: { "Content-Type": "multipart/form-data" },
      })
      .then((r) => r.data),

  schedule: (deviceId: string, body: SendScheduledMessageRequest) =>
    api
      .post<Message>(`/api/v1/devices/${deviceId}/messages/scheduled`, body)
      .then((r) => r.data),

  getStatus: (messageId: string) =>
    api
      .get<Message>(`/api/v1/messages/${messageId}/status`)
      .then((r) => r.data),

  queueStats: () =>
    api.get<QueueStats>("/api/v1/messages/stats").then((r) => r.data),

  failed: () =>
    api.get<Message[]>("/api/v1/messages/failed").then((r) => r.data),
};

// ─── Server Info ──────────────────────────────────────────────────────────────

export const infoApi = {
  status: () =>
    api.get<ServerStatus>("/api/v1/info/status").then((r) => r.data),

  stats: () =>
    api.get<ServerStats>("/api/v1/info/stats").then((r) => r.data),

  health: () => api.get<HealthResponse>("/health").then((r) => r.data),
};

export default api;
