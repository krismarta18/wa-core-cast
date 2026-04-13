"use client";

import {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { useRouter } from "next/navigation";
import { authApi, clearDashboardSession, persistAuthSession, setUnauthorizedHandler } from "@/lib/api";
import {
  clearAuthSession,
  getStoredAuthSession,
  subscribeAuthSession,
  type StoredAuthSession,
} from "@/lib/auth-session";
import type { AuthSessionResponse } from "@/lib/types";

type AuthContextValue = {
  session: StoredAuthSession | null;
  hydrated: boolean;
  completeAuth: (response: AuthSessionResponse, persistent: boolean) => void;
  logout: () => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const [session, setSession] = useState<StoredAuthSession | null>(null);
  const [hydrated, setHydrated] = useState(false);

  useEffect(() => {
    const syncSession = () => setSession(getStoredAuthSession());

    syncSession();
    setHydrated(true);

    const unsubscribe = subscribeAuthSession(syncSession);

    setUnauthorizedHandler(() => {
      clearAuthSession();
      router.replace("/login");
    });

    return () => {
      unsubscribe();
      setUnauthorizedHandler(null);
    };
  }, [router]);

  const value = useMemo<AuthContextValue>(
    () => ({
      session,
      hydrated,
      completeAuth: (response, persistent) => {
        const nextSession = persistAuthSession(response, persistent);
        setSession(nextSession);
      },
      logout: async () => {
        try {
          await authApi.logout();
        } finally {
          clearDashboardSession();
          setSession(null);
          router.replace("/login");
        }
      },
    }),
    [hydrated, router, session]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }

  return context;
}