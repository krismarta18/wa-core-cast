import axios from "axios";
import type {
  AuthGenericSuccessResponse,
  AuthMeResponse,
  AuthOTPRequest,
  AuthRefreshTokenRequest,
  AuthRegisterRequest,
  AuthSessionResponse,
  AuthVerifyOTPRequest,
  BillingCheckoutRequest,
  BillingCheckoutResponse,
  BillingOverviewResponse,
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
import {
  clearAuthSession,
  getStoredAuthSession,
  isRefreshExpired,
  saveAuthSession,
  type StoredAuthSession,
} from "./auth-session";

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080",
  headers: { "Content-Type": "application/json" },
  timeout: 10_000,
});

const authClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080",
  headers: { "Content-Type": "application/json" },
  timeout: 10_000,
});

let refreshPromise: Promise<StoredAuthSession | null> | null = null;
let unauthorizedHandler: (() => void) | null = null;

function notifyUnauthorized() {
  unauthorizedHandler?.();
}

function toStoredSession(response: AuthSessionResponse, persistent: boolean): StoredAuthSession {
  const now = Date.now();

  return {
    accessToken: response.access_token,
    refreshToken: response.refresh_token,
    expiresAt: now + response.expires_in * 1000,
    refreshExpiresAt: now + response.refresh_expires_in * 1000,
    user: response.user,
    persistent,
  };
}

export function setUnauthorizedHandler(handler: (() => void) | null) {
  unauthorizedHandler = handler;
}

export function persistAuthSession(response: AuthSessionResponse, persistent: boolean) {
  const session = toStoredSession(response, persistent);
  saveAuthSession(session);
  return session;
}

export function clearDashboardSession() {
  clearAuthSession();
}

async function refreshAccessToken() {
  const currentSession = getStoredAuthSession();
  if (!currentSession || isRefreshExpired(currentSession)) {
    clearAuthSession();
    notifyUnauthorized();
    return null;
  }

  if (!refreshPromise) {
    refreshPromise = authApi
      .refresh({ refresh_token: currentSession.refreshToken })
      .then((response) => persistAuthSession(response, currentSession.persistent))
      .catch(() => {
        clearAuthSession();
        notifyUnauthorized();
        return null;
      })
      .finally(() => {
        refreshPromise = null;
      });
  }

  return refreshPromise;
}

api.interceptors.request.use((config) => {
  const session = getStoredAuthSession();
  if (session?.accessToken) {
    config.headers.Authorization = `Bearer ${session.accessToken}`;
  }

  return config;
});

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config as typeof error.config & { _retry?: boolean };

    if (error.response?.status !== 401 || originalRequest?._retry) {
      return Promise.reject(error);
    }

    originalRequest._retry = true;

    const refreshedSession = await refreshAccessToken();
    if (!refreshedSession) {
      return Promise.reject(error);
    }

    originalRequest.headers = originalRequest.headers ?? {};
    originalRequest.headers.Authorization = `Bearer ${refreshedSession.accessToken}`;
    return api(originalRequest);
  }
);

// ─── Auth ─────────────────────────────────────────────────────────────────────

export const authApi = {
  register: (body: AuthRegisterRequest) =>
    authClient
      .post<AuthGenericSuccessResponse>("/api/v1/auth/register", body)
      .then((r) => r.data),

  requestOTP: (body: AuthOTPRequest) =>
    authClient
      .post<AuthGenericSuccessResponse>("/api/v1/auth/request-otp", body)
      .then((r) => r.data),

  verifyOTP: (body: AuthVerifyOTPRequest) =>
    authClient
      .post<AuthSessionResponse>("/api/v1/auth/verify-otp", body)
      .then((r) => r.data),

  refresh: (body: AuthRefreshTokenRequest) =>
    authClient
      .post<AuthSessionResponse>("/api/v1/auth/refresh", body)
      .then((r) => r.data),

  me: () => api.get<AuthMeResponse>("/api/v1/auth/me").then((r) => r.data),

  logout: () => api.post<AuthGenericSuccessResponse>("/api/v1/auth/logout").then((r) => r.data),
};

// ─── Billing ──────────────────────────────────────────────────────────────────

export const billingApi = {
  overview: () =>
    api.get<BillingOverviewResponse>("/api/v1/billing/overview").then((r) => r.data),

  checkout: (body: BillingCheckoutRequest) =>
    api.post<BillingCheckoutResponse>("/api/v1/billing/checkout", body).then((r) => r.data),
};

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
