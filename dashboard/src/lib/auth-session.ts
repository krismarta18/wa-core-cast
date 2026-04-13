import type { AuthUser } from "./types";

const LOCAL_KEY = "wacast.auth.session";
const SESSION_KEY = "wacast.auth.session.temp";
const PENDING_KEY = "wacast.auth.pending";

export interface StoredAuthSession {
  accessToken: string;
  refreshToken: string;
  expiresAt: number;
  refreshExpiresAt: number;
  user: AuthUser;
  persistent: boolean;
}

export interface PendingAuthState {
  phoneNumber: string;
  context: "login" | "register";
  fullName?: string;
  rememberMe?: boolean;
}

type Listener = () => void;

const listeners = new Set<Listener>();

function emitChange() {
  listeners.forEach((listener) => listener());
}

function canUseStorage() {
  return typeof window !== "undefined";
}

function parseJSON<T>(value: string | null): T | null {
  if (!value) {
    return null;
  }

  try {
    return JSON.parse(value) as T;
  } catch {
    return null;
  }
}

function getLocalStorage() {
  return canUseStorage() ? window.localStorage : null;
}

function getSessionStorage() {
  return canUseStorage() ? window.sessionStorage : null;
}

export function subscribeAuthSession(listener: Listener) {
  listeners.add(listener);
  return () => listeners.delete(listener);
}

export function getStoredAuthSession(): StoredAuthSession | null {
  const localSession = parseJSON<StoredAuthSession>(getLocalStorage()?.getItem(LOCAL_KEY) ?? null);
  if (localSession) {
    return { ...localSession, persistent: true };
  }

  const session = parseJSON<StoredAuthSession>(getSessionStorage()?.getItem(SESSION_KEY) ?? null);
  if (session) {
    return { ...session, persistent: false };
  }

  return null;
}

export function saveAuthSession(session: StoredAuthSession) {
  const payload = JSON.stringify(session);

  if (session.persistent) {
    getLocalStorage()?.setItem(LOCAL_KEY, payload);
    getSessionStorage()?.removeItem(SESSION_KEY);
  } else {
    getSessionStorage()?.setItem(SESSION_KEY, payload);
    getLocalStorage()?.removeItem(LOCAL_KEY);
  }

  emitChange();
}

export function clearAuthSession() {
  getLocalStorage()?.removeItem(LOCAL_KEY);
  getSessionStorage()?.removeItem(SESSION_KEY);
  emitChange();
}

export function savePendingAuthState(state: PendingAuthState) {
  getSessionStorage()?.setItem(PENDING_KEY, JSON.stringify(state));
}

export function getPendingAuthState(): PendingAuthState | null {
  return parseJSON<PendingAuthState>(getSessionStorage()?.getItem(PENDING_KEY) ?? null);
}

export function clearPendingAuthState() {
  getSessionStorage()?.removeItem(PENDING_KEY);
}

export function isSessionExpired(session: StoredAuthSession) {
  return session.expiresAt <= Date.now();
}

export function isRefreshExpired(session: StoredAuthSession) {
  return session.refreshExpiresAt <= Date.now();
}