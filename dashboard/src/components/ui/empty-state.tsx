"use client";

import { FileX, Search, WifiOff, Inbox } from "lucide-react";

interface EmptyStateProps {
  variant?: "no-data" | "no-results" | "offline" | "inbox";
  title?: string;
  description?: string;
  action?: React.ReactNode;
}

const VARIANTS = {
  "no-data": { icon: FileX, defaultTitle: "Belum ada data", defaultDesc: "Data akan muncul di sini setelah ada aktivitas." },
  "no-results": { icon: Search, defaultTitle: "Tidak ada hasil", defaultDesc: "Coba ubah kata kunci atau filter pencarian." },
  offline: { icon: WifiOff, defaultTitle: "Tidak dapat memuat", defaultDesc: "Periksa koneksi internet kamu lalu coba lagi." },
  inbox: { icon: Inbox, defaultTitle: "Kosong", defaultDesc: "Tidak ada item untuk ditampilkan." },
};

export function EmptyState({ variant = "no-data", title, description, action }: EmptyStateProps) {
  const { icon: Icon, defaultTitle, defaultDesc } = VARIANTS[variant];
  return (
    <div className="flex flex-col items-center justify-center py-16 text-center">
      <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-gray-100">
        <Icon className="h-8 w-8 text-gray-300" />
      </div>
      <p className="mt-4 text-base font-semibold text-gray-700">{title ?? defaultTitle}</p>
      <p className="mt-1 text-sm text-gray-400">{description ?? defaultDesc}</p>
      {action && <div className="mt-5">{action}</div>}
    </div>
  );
}
