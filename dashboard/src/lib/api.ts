import axios from "axios";
import type {
  AuthGenericSuccessResponse,
  AuthMeResponse,
  AuthOTPRequest,
  AuthRefreshTokenRequest,
  AuthRegisterRequest,
  AuthRegisterRequest,
  AuthSessionResponse,
  AuthUpdateProfileRequest,
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
  Contact,
  CreateContactRequest,
  UpdateContactRequest,
  ContactGroup,
  CreateContactGroupRequest,
  BlacklistEntry,
  UsageStatsResponse,
  FailureRateResponse,
  BroadcastCampaign,
  CreateBroadcastRequest,
  BroadcastStatus,
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

  updateProfile: (body: AuthUpdateProfileRequest) =>
    api.put<AuthMeResponse>("/api/v1/auth/profile", body).then((r) => r.data),

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

  schedule: (deviceId: string, body: SendScheduledMessageRequest | FormData) =>
    api
      .post<Message>(`/api/v1/devices/${deviceId}/messages/scheduled`, body, {
        headers: body instanceof FormData ? { 'Content-Type': 'multipart/form-data' } : undefined
      })
      .then((r) => r.data),

  listScheduled: (deviceId: string) =>
    api
      .get<{ messages: Message[] }>(`/api/v1/devices/${deviceId}/messages/scheduled`)
      .then((r) => r.data),

  listHistory: (deviceId: string) =>
    api
      .get<{ messages: Message[] }>(`/api/v1/devices/${deviceId}/messages/history`)
      .then((r) => r.data),

  cancelScheduled: (messageId: string) =>
    api
      .delete<{ message: string }>(`/api/v1/messages/${messageId}`)
      .then((r) => r.data),

  getStatus: (messageId: string) =>
    api
      .get<Message>(`/api/v1/messages/${messageId}/status`)
      .then((r) => r.data),

  queueStats: () =>
    api.get<QueueStats>("/api/v1/messages/stats").then((r) => r.data),

  failed: () =>
    api.get<Message[]>("/api/v1/messages/failed").then((r) => r.data),

  uploadMedia: (formData: FormData) =>
    api
      .post<{ url: string; filename: string; size: number }>("/api/v1/upload/media", formData, {
        headers: { "Content-Type": "multipart/form-data" },
      })
      .then((r) => r.data),
  
  getLogs: (limit: number = 50, offset: number = 0) =>
    api
      .get<{ messages: Message[]; count: number; limit: number; offset: number }>(
        "/api/v1/messages/logs",
        { params: { limit, offset } }
      )
      .then((r) => r.data),
};

// ─── Server Info ──────────────────────────────────────────────────────────────

export const infoApi = {
  status: () =>
    api.get<ServerStatus>("/api/v1/info/status").then((r) => r.data),

  stats: () =>
    api.get<ServerStats>("/api/v1/info/stats").then((r) => r.data),

  health: () => api.get<HealthResponse>("/health").then((r) => r.data),
};

// ─── Contacts ─────────────────────────────────────────────────────────────────

export const contactsApi = {
  list: () =>
    api.get<{ contacts: Contact[] }>("/api/v1/contacts").then((r) => r.data),

  create: (body: CreateContactRequest) =>
    api.post<Contact>("/api/v1/contacts", body).then((r) => r.data),

  update: (id: string, body: UpdateContactRequest) =>
    api.put<Contact>(`/api/v1/contacts/${id}`, body).then((r) => r.data),

  delete: (id: string) =>
    api.delete<{ message: string }>(`/api/v1/contacts/${id}`).then((r) => r.data),
};

// ─── Contact Groups ───────────────────────────────────────────────────────────

export const groupsApi = {
  list: () =>
    api.get<{ groups: ContactGroup[] }>("/api/v1/contact-groups").then((r) => r.data),

  create: (body: CreateContactGroupRequest) =>
    api.post<ContactGroup>("/api/v1/contact-groups", body).then((r) => r.data),

  delete: (id: string) =>
    api.delete<{ message: string }>(`/api/v1/contact-groups/${id}`).then((r) => r.data),

  listMembers: (id: string) =>
    api.get<{ members: Contact[] }>(`/api/v1/contact-groups/${id}/members`).then((r) => r.data),

  addMember: (id: string, contactId: string) =>
    api.post<{ message: string }>(`/api/v1/contact-groups/${id}/members`, { contact_id: contactId }).then((r) => r.data),

  removeMember: (id: string, contactId: string) =>
    api.delete<{ message: string }>(`/api/v1/contact-groups/${id}/members/${contactId}`).then((r) => r.data),
};

// ─── Blacklist ────────────────────────────────────────────────────────────────

export const blacklistApi = {
  list: () =>
    api.get<{ blacklist: BlacklistEntry[] }>("/api/v1/blacklists").then((r) => r.data),

  block: (body: { phone_number: string; reason: string }) =>
    api.post<{ message: string }>("/api/v1/blacklists", body).then((r) => r.data),

  unblock: (id: string) =>
    api.delete<{ message: string }>(`/api/v1/blacklists/${id}`).then((r) => r.data),
};

// ─── Analytics ────────────────────────────────────────────────────────────────

export const analyticsApi = {
  usage: () =>
    api.get<UsageStatsResponse>("/api/v1/analytics/usage").then((r) => r.data),

  failures: () =>
    api.get<FailureRateResponse>("/api/v1/analytics/failures").then((r) => r.data),
};

// ─── Broadcasts ───────────────────────────────────────────────────────────────

export const broadcastApi = {
  list: () =>
    api.get<{ broadcasts: BroadcastCampaign[] }>("/api/v1/broadcasts").then((r) => r.data),

  create: (body: CreateBroadcastRequest) =>
    api.post<BroadcastCampaign>("/api/v1/broadcasts", body).then((r) => r.data),

  get: (id: string) =>
    api.get<BroadcastCampaign>(`/api/v1/broadcasts/${id}`).then((r) => r.data),

  start: (id: string) =>
    api.post<{ message: string }>(`/api/v1/broadcasts/${id}/start`).then((r) => r.data),
};

export default api;
