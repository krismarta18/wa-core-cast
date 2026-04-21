"use client";

import {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { useRouter } from "next/navigation";
import { authApi, clearDashboardSession, infoApi, licenseApi, persistAuthSession, setLicenseHandler, setUnauthorizedHandler } from "@/lib/api";
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
  isLicensed: boolean;
  appConfig: Record<string, string>;
  isConfigLoading: boolean;
  completeAuth: (response: AuthSessionResponse, persistent: boolean) => void;
  logout: () => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const [session, setSession] = useState<StoredAuthSession | null>(null);
  const [hydrated, setHydrated] = useState(false);
  const [isLicensed, setIsLicensed] = useState(true);
  const [appConfig, setAppConfig] = useState<Record<string, string>>({});
  const [isConfigLoading, setIsConfigLoading] = useState(true);

  useEffect(() => {
    const syncSession = () => setSession(getStoredAuthSession());

    const checkLicense = async () => {
      try {
        const res = await licenseApi.getStatus();
        const active = res.data.is_active && !res.data.is_expired;
        setIsLicensed(active);

        const path = window.location.pathname;

        if (!active) {
          // Force to setup if not licensed
          if (path !== "/setup") {
            router.replace("/setup");
          }
        } else {
          // Block setup if already licensed
          if (path === "/setup") {
            router.replace("/");
          }
        }
      } catch (error) {
        console.error("License check failed", error);
      }
    };

    const loadConfig = async () => {
      try {
        const res = await infoApi.getConfig();
        if (res.success) {
          setAppConfig(res.config);
        }
      } catch (err) {
        console.error("Failed to load app config", err);
      } finally {
        setIsConfigLoading(false);
      }
    };

    syncSession();
    checkLicense();
    loadConfig();
    setHydrated(true);

    const unsubscribe = subscribeAuthSession(syncSession);

    setUnauthorizedHandler(() => {
      clearAuthSession();
      router.replace("/login");
    });

    setLicenseHandler(() => {
      setIsLicensed(false);
      router.replace("/setup");
    });

    return () => {
      unsubscribe();
      setUnauthorizedHandler(null);
      setLicenseHandler(null);
    };
  }, [router]);

  const value = useMemo<AuthContextValue>(
    () => ({
      session,
      hydrated,
      isLicensed,
      appConfig,
      isConfigLoading,
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
    [hydrated, isLicensed, router, session]
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