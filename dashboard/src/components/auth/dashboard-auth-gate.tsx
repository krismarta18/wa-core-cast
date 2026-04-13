"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/providers/auth-provider";

export function DashboardAuthGate({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const { session, hydrated } = useAuth();

  useEffect(() => {
    if (hydrated && !session) {
      router.replace("/login");
    }
  }, [hydrated, router, session]);

  if (!hydrated || !session) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50 text-sm text-gray-500">
        Memuat sesi dashboard...
      </div>
    );
  }

  return <>{children}</>;
}