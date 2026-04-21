"use client";

import { useState, useRef, useEffect } from "react";
import { Bell, RefreshCw, Wifi, AlertTriangle, CheckCircle2, Info, Search, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useQueryClient } from "@tanstack/react-query";
import { MobileSidebarTrigger } from "./mobile-sidebar-trigger";
import { useRouter } from "next/navigation";

interface HeaderProps {
  title: string;
  description?: string;
}

const NOTIFS = [
  { id: 1, type: "warning", icon: Wifi, title: "Device Terputus", body: "Device Backup kehilangan koneksi WA.", time: "2 mnt lalu", read: false },
  { id: 2, type: "alert", icon: AlertTriangle, title: "Kuota 85%", body: "Kuota pesan sudah terpakai 8.500 / 10.000.", time: "1 jam lalu", read: false },
  { id: 3, type: "success", icon: CheckCircle2, title: "Broadcast Selesai", body: "Broadcast ke 340 kontak berhasil terkirim.", time: "3 jam lalu", read: true },
  { id: 4, type: "info", icon: Info, title: "Perpanjangan Otomatis", body: "Paket Business Pro diperbarui 12 Mei 2026.", time: "Kemarin", read: true },
];

const TYPE_STYLE: Record<string, string> = {
  warning: "bg-yellow-50 text-yellow-600",
  alert: "bg-red-50 text-red-600",
  success: "bg-green-50 text-green-600",
  info: "bg-blue-50 text-blue-600",
};

const SEARCH_INDEX = [
  { label: "Connection Status", href: "/devices/status", group: "Device" },
  { label: "Connect New Device", href: "/devices/qr", group: "Device" },
  { label: "Multi Device Info", href: "/devices/info", group: "Device" },
  { label: "Session Management", href: "/devices/session", group: "Device" },
  { label: "New Message", href: "/messaging/new", group: "Messaging" },
  { label: "Broadcast", href: "/messaging/broadcast", group: "Messaging" },
  { label: "Scheduled Messages", href: "/messaging/scheduled", group: "Messaging" },
  { label: "Message Logs", href: "/messaging/logs", group: "Messaging" },
  { label: "API Keys", href: "/api-integration/keys", group: "API" },
  { label: "Webhook Settings", href: "/api-integration/webhooks", group: "API" },
  { label: "Keyword Auto-Reply", href: "/auto-response/keywords", group: "Auto Response" },
  { label: "Message Templates", href: "/auto-response/templates", group: "Auto Response" },
  { label: "Usage Statistics", href: "/monitoring/usage", group: "Monitoring" },
  { label: "Failure Rate", href: "/monitoring/failure", group: "Monitoring" },
  { label: "Phone Book", href: "/contacts/phonebook", group: "Contacts" },
  { label: "Group Contacts", href: "/contacts/groups", group: "Contacts" },
  { label: "Blacklist", href: "/contacts/blacklist", group: "Contacts" },
  { label: "Profile & Account", href: "/settings/profile", group: "Settings" },
  { label: "Billing & Quota", href: "/settings/billing", group: "Settings" },
  { label: "Notifications", href: "/settings/notifications", group: "Settings" },
];

export function Header({ title, description }: HeaderProps) {
  const queryClient = useQueryClient();
  const router = useRouter();
  const [notifOpen, setNotifOpen] = useState(false);
  const [notifs, setNotifs] = useState(NOTIFS);
  const [searchOpen, setSearchOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const notifRef = useRef<HTMLDivElement>(null);
  const searchRef = useRef<HTMLDivElement>(null);
  const searchInputRef = useRef<HTMLInputElement>(null);

  const unread = notifs.filter((n) => !n.read).length;

  const searchResults =
    searchQuery.length >= 1
      ? SEARCH_INDEX.filter(
          (item) =>
            item.label.toLowerCase().includes(searchQuery.toLowerCase()) ||
            item.group.toLowerCase().includes(searchQuery.toLowerCase())
        ).slice(0, 8)
      : [];

  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (notifRef.current && !notifRef.current.contains(e.target as Node))
        setNotifOpen(false);
      if (
        searchRef.current &&
        !searchRef.current.contains(e.target as Node)
      ) {
        setSearchOpen(false);
        setSearchQuery("");
      }
    }
    document.addEventListener("mousedown", handleClick);
    return () => document.removeEventListener("mousedown", handleClick);
  }, []);

  useEffect(() => {
    function handleKey(e: KeyboardEvent) {
      if ((e.metaKey || e.ctrlKey) && e.key === "k") {
        e.preventDefault();
        setSearchOpen(true);
        setTimeout(() => searchInputRef.current?.focus(), 50);
      }
      if (e.key === "Escape") {
        setSearchOpen(false);
        setSearchQuery("");
      }
    }
    document.addEventListener("keydown", handleKey);
    return () => document.removeEventListener("keydown", handleKey);
  }, []);

  function markAllRead() {
    setNotifs(notifs.map((n) => ({ ...n, read: true })));
  }

  function goTo(href: string) {
    router.push(href);
    setSearchOpen(false);
    setSearchQuery("");
  }

  return (
    <header className="flex h-16 items-center justify-between border-b border-gray-200 bg-white px-6 gap-4">
      <div className="flex items-center gap-3 min-w-0">
        <MobileSidebarTrigger />
        <div className="min-w-0">
          <h1 className="truncate text-base font-semibold text-gray-900">{title}</h1>
          {description && (
            <p className="truncate text-xs text-gray-500">{description}</p>
          )}
        </div>
      </div>

      <div className="flex items-center gap-2 flex-shrink-0">
        {/* Global search */}
        <div ref={searchRef} className="relative">
          <button
            onClick={() => {
              setSearchOpen(true);
              setTimeout(() => searchInputRef.current?.focus(), 50);
            }}
            className="flex items-center gap-2 rounded-lg border border-gray-200 bg-gray-50 px-3 py-1.5 text-sm text-gray-400 hover:bg-gray-100 transition-colors"
          >
            <Search className="h-3.5 w-3.5" />
            <span className="hidden sm:inline">Cari halaman...</span>
            <kbd className="hidden sm:inline-flex items-center rounded border border-gray-200 bg-white px-1.5 py-0.5 text-[10px] font-mono text-gray-400">
              Ctrl+K
            </kbd>
          </button>

          {searchOpen && (
            <div className="absolute right-0 top-11 z-50 w-80 rounded-xl border border-gray-200 bg-white shadow-xl overflow-hidden">
              <div className="flex items-center gap-2 border-b border-gray-100 px-3 py-2.5">
                <Search className="h-4 w-4 flex-shrink-0 text-gray-400" />
                <input
                  ref={searchInputRef}
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  placeholder="Cari halaman atau fitur..."
                  className="flex-1 text-sm text-gray-700 placeholder:text-gray-400 focus:outline-none"
                />
                {searchQuery && (
                  <button onClick={() => setSearchQuery("")}>
                    <X className="h-3.5 w-3.5 text-gray-400" />
                  </button>
                )}
              </div>
              {searchQuery.length === 0 ? (
                <div className="px-4 py-6 text-center text-xs text-gray-400">
                  Ketik untuk mencari halaman atau fitur
                </div>
              ) : searchResults.length === 0 ? (
                <div className="px-4 py-6 text-center text-xs text-gray-400">
                  Tidak ada hasil untuk &ldquo;{searchQuery}&rdquo;
                </div>
              ) : (
                <ul className="max-h-64 overflow-y-auto py-1">
                  {searchResults.map((item) => (
                    <li key={item.href}>
                      <button
                        onClick={() => goTo(item.href)}
                        className="flex w-full items-center justify-between px-4 py-2.5 text-left hover:bg-gray-50"
                      >
                        <span className="text-sm font-medium text-gray-800">
                          {item.label}
                        </span>
                        <span className="text-xs text-gray-400">{item.group}</span>
                      </button>
                    </li>
                  ))}
                </ul>
              )}
            </div>
          )}
        </div>

        <Button
          variant="ghost"
          size="icon"
          onClick={() => queryClient.invalidateQueries()}
          title="Refresh all data"
        >
          <RefreshCw className="h-4 w-4" />
        </Button>

        {/* Notification bell */}
        <div ref={notifRef} className="relative">
          <button
            onClick={() => setNotifOpen((v) => !v)}
            className="relative flex h-9 w-9 items-center justify-center rounded-lg hover:bg-gray-100"
          >
            <Bell className="h-4 w-4 text-gray-600" />
            {unread > 0 && (
              <span className="absolute right-1.5 top-1.5 flex h-4 w-4 items-center justify-center rounded-full bg-red-500 text-[10px] font-bold text-white">
                {unread}
              </span>
            )}
          </button>

          {notifOpen && (
            <div className="absolute right-0 top-11 z-50 w-80 rounded-xl border border-gray-200 bg-white shadow-xl">
              <div className="flex items-center justify-between border-b border-gray-100 px-4 py-3">
                <p className="font-semibold text-gray-900">Notifikasi</p>
                {unread > 0 && (
                  <button
                    onClick={markAllRead}
                    className="text-xs text-green-600 hover:underline"
                  >
                    Tandai semua dibaca
                  </button>
                )}
              </div>
              <ul className="max-h-80 overflow-y-auto divide-y divide-gray-50">
                {notifs.map((n) => (
                  <li
                    key={n.id}
                    className={`flex gap-3 px-4 py-3 hover:bg-gray-50 cursor-pointer ${
                      !n.read ? "bg-blue-50/40" : ""
                    }`}
                    onClick={() =>
                      setNotifs(
                        notifs.map((x) =>
                          x.id === n.id ? { ...x, read: true } : x
                        )
                      )
                    }
                  >
                    <div
                      className={`mt-0.5 flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg ${
                        TYPE_STYLE[n.type]
                      }`}
                    >
                      <n.icon className="h-4 w-4" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <p
                        className={`text-sm font-medium ${
                          !n.read ? "text-gray-900" : "text-gray-600"
                        }`}
                      >
                        {n.title}
                      </p>
                      <p className="text-xs text-gray-400 leading-snug">{n.body}</p>
                      <p className="mt-0.5 text-xs text-gray-300">{n.time}</p>
                    </div>
                    {!n.read && (
                      <span className="mt-1.5 h-2 w-2 flex-shrink-0 rounded-full bg-blue-500" />
                    )}
                  </li>
                ))}
              </ul>
              <div className="border-t border-gray-100 px-4 py-2.5 text-center">
                <button className="text-xs text-gray-400 hover:text-gray-600">
                  Lihat semua notifikasi
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}
