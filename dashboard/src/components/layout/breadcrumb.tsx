"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { ChevronRight, Home } from "lucide-react";

const LABELS: Record<string, string> = {
  devices: "Device Management",
  status: "Connection Status",
  qr: "Connect New Device",
  info: "Multi Device Info",
  session: "Session Management",
  messaging: "Messaging",
  new: "New Message",
  broadcast: "Broadcast",
  scheduled: "Scheduled",
  logs: "Message Logs",
  "api-integration": "API & Integration",
  keys: "API Keys",
  webhooks: "Webhook Settings",
  "auto-response": "Auto Response & Template",
  keywords: "Keyword",
  templates: "Message Template",
  monitoring: "Monitoring & Analytics",
  usage: "Usage Statistics",
  failure: "Failure Rate",
  contacts: "Contact Management",
  phonebook: "Phone Book",
  groups: "Group Contact",
  blacklist: "Blacklist / Block",
  settings: "Settings",
  profile: "Profil & Akun",
  billing: "Billing & Kuota",
  notifications: "Notifikasi",
  onboarding: "Quick Start",
};

// Segments that are group prefixes only — no real page behind them
const NO_PAGE_SEGMENTS = new Set([
  "devices",
  "messaging",
  "api-integration",
  "auto-response",
  "monitoring",
  "contacts",
  "settings",
]);

export function Breadcrumb() {
  const pathname = usePathname();
  if (pathname === "/") return null;

  const segments = pathname.split("/").filter(Boolean);
  const crumbs = segments.map((seg, i) => ({
    label: LABELS[seg] ?? seg,
    href: "/" + segments.slice(0, i + 1).join("/"),
    isLast: i === segments.length - 1,
    hasPage: !NO_PAGE_SEGMENTS.has(seg),
  }));

  return (
    <nav className="flex items-center gap-1 border-b border-gray-100 bg-white px-6 py-2 text-xs text-gray-400">
      <Link href="/" className="flex items-center gap-1 hover:text-gray-600">
        <Home className="h-3 w-3" /> Dashboard
      </Link>
      {crumbs.map((c) => (
        <span key={c.href} className="flex items-center gap-1">
          <ChevronRight className="h-3 w-3" />
          {c.isLast || !c.hasPage ? (
            <span className={c.isLast ? "font-medium text-gray-600" : ""}>{c.label}</span>
          ) : (
            <Link href={c.href} className="hover:text-gray-600">{c.label}</Link>
          )}
        </span>
      ))}
    </nav>
  );
}
