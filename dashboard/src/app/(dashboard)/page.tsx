"use client";

import {
  Smartphone,
  Wifi,
  MessageSquare,
  Crown,
  AlertCircle,
  CheckCircle2,
  Clock,
  TrendingUp,
  Loader2,
} from "lucide-react";

import { useEffect, useState } from "react";
import { 
  messagesApi, 
  billingApi, 
  sessionsApi, 
  analyticsApi 
} from "@/lib/api";
import { 
  BillingCurrentPlan, 
  Device, 
  Message, 
  UsageStatsResponse 
} from "@/lib/types";

function QuotaBar({ used, max }: { used: number; max: number }) {
  const pct = Math.min((used / max) * 100, 100);
  const color = pct >= 90 ? "bg-red-500" : pct >= 70 ? "bg-yellow-500" : "bg-green-500";
  return (
    <div className="mt-2">
      <div className="h-2 w-full overflow-hidden rounded-full bg-gray-100">
        <div className={`h-full rounded-full transition-all ${color}`} style={{ width: `${pct}%` }} />
      </div>
      <p className="mt-1 text-xs text-gray-500">
        {used.toLocaleString()} / {max.toLocaleString()} ({pct.toFixed(0)}%)
      </p>
    </div>
  );
}

function StatusDot({ status }: { status: number | string }) {
  // 1 = active (connected)
  const isConnected = status === 1 || status === "connected";
  return (
    <span
      className={`inline-block h-2 w-2 rounded-full flex-shrink-0 ${
        isConnected ? "bg-green-500" : "bg-gray-300"
      }`}
    />
  );
}

function MessageStatusBadge({ status }: { status: number | string }) {
  // Map integers to string if needed
  const s = typeof status === "number" ? 
    { 0: "pending", 1: "sent", 2: "delivered", 3: "read", 4: "failed" }[status] || "pending" 
    : status;

  if (s === "sent" || s === "delivered" || s === "read")
    return (
      <span className="inline-flex items-center gap-1 rounded-full bg-green-50 px-2 py-0.5 text-xs font-medium text-green-700">
        <CheckCircle2 className="h-3 w-3" /> {s === "read" ? "Dibaca" : "Terkirim"}
      </span>
    );
  if (s === "failed")
    return (
      <span className="inline-flex items-center gap-1 rounded-full bg-red-50 px-2 py-0.5 text-xs font-medium text-red-700">
        <AlertCircle className="h-3 w-3" /> Gagal
      </span>
    );
  return (
    <span className="inline-flex items-center gap-1 rounded-full bg-yellow-50 px-2 py-0.5 text-xs font-medium text-yellow-700">
      <Clock className="h-3 w-3" /> Pending
    </span>
  );
}

export default function DashboardPage() {
  const [billing, setBilling] = useState<BillingCurrentPlan | null>(null);
  const [devices, setDevices] = useState<Device[]>([]);
  const [recentMessages, setRecentMessages] = useState<Message[]>([]);
  const [stats, setStats] = useState<UsageStatsResponse | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchDashboardData() {
      try {
        setLoading(true);
        const [billingRes, sessionsRes, messagesRes, statsRes] = await Promise.all([
          billingApi.overview(),
          sessionsApi.list(),
          messagesApi.getLogs(5, 0),
          analyticsApi.usage(),
        ]);

        setBilling(billingRes.billing.current_plan || null);
        setDevices(sessionsRes.sessions || []);
        setRecentMessages(messagesRes.messages || []);
        setStats(statsRes);
      } catch (err) {
        console.error("Dashboard data fetch failed", err);
      } finally {
        setLoading(false);
      }
    }
    fetchDashboardData();
  }, []);

  const connectedDevices = devices.filter((d) => d.status === 1 || d.status === "connected" as any).length;
  const quotaLimit = billing?.quota_limit || 0;
  const quotaUsed = billing?.quota_used || 0;
  const deviceMax = billing?.device_max || 0;
  const deviceUsed = billing?.device_used || 0;

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50">
        <Loader2 className="h-8 w-8 animate-spin text-green-600" />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Page header */}
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <h1 className="text-xl font-bold text-gray-900">Dashboard</h1>
        <p className="text-sm text-gray-500">Selamat datang kembali! Berikut ringkasan akun Anda.</p>
      </div>

      <div className="p-6 space-y-6">

        {/* ── Plan & Quota row ─────────────────────────────────────────────── */}
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
          {/* Paket aktif */}
          <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
            <div className="flex items-center justify-between">
              <p className="text-sm font-medium text-gray-500">Paket Aktif</p>
              <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-yellow-50">
                <Crown className="h-5 w-5 text-yellow-500" />
              </div>
            </div>
            <p className="mt-2 text-2xl font-bold text-gray-900">{billing?.name || "Free Plan"}</p>
            <p className="mt-1 text-xs font-medium text-green-600">
              ● {billing?.renewal_date ? `Renewal ${new Date(billing.renewal_date).toLocaleDateString()}` : "Active"}
            </p>
          </div>

          {/* Kuota pesan */}
          <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
            <div className="flex items-center justify-between">
              <p className="text-sm font-medium text-gray-500">Kuota Pesan</p>
              <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-blue-50">
                <MessageSquare className="h-5 w-5 text-blue-500" />
              </div>
            </div>
            <p className="mt-2 text-2xl font-bold text-gray-900">
              {(quotaLimit - quotaUsed).toLocaleString()}
              <span className="ml-1 text-sm font-normal text-gray-400">sisa</span>
            </p>
            <QuotaBar used={quotaUsed} max={quotaLimit} />
          </div>

          {/* Device slot */}
          <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
            <div className="flex items-center justify-between">
              <p className="text-sm font-medium text-gray-500">Slot Device</p>
              <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-purple-50">
                <Smartphone className="h-5 w-5 text-purple-500" />
              </div>
            </div>
            <p className="mt-2 text-2xl font-bold text-gray-900">
              {deviceUsed}
              <span className="text-lg font-normal text-gray-400"> / {deviceMax}</span>
            </p>
            <QuotaBar used={deviceUsed} max={deviceMax} />
          </div>

          {/* Device aktif */}
          <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
            <div className="flex items-center justify-between">
              <p className="text-sm font-medium text-gray-500">Device Aktif</p>
              <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-green-50">
                <Wifi className="h-5 w-5 text-green-500" />
              </div>
            </div>
            <p className="mt-2 text-2xl font-bold text-gray-900">
              {connectedDevices}
              <span className="text-lg font-normal text-gray-400"> / {deviceUsed}</span>
            </p>
            <p className="mt-2 text-xs text-gray-500">
              {deviceUsed - connectedDevices} device terputus
            </p>
          </div>
        </div>

        {/* ── Bottom row ───────────────────────────────────────────────────── */}
        <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
          {/* Device list */}
          <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
            <div className="flex items-center justify-between border-b border-gray-100 px-5 py-4">
              <h2 className="font-semibold text-gray-900">Status Device</h2>
              <a href="/devices/session" className="text-xs font-medium text-green-600 hover:underline">
                Lihat semua →
              </a>
            </div>
            <ul className="divide-y divide-gray-50 max-h-[320px] overflow-y-auto">
              {devices.length > 0 ? devices.map((d) => (
                <li key={d.device_id} className="flex items-center gap-3 px-5 py-3 transition-colors hover:bg-gray-50">
                  <StatusDot status={d.status} />
                  <div className="min-w-0 flex-1">
                    <p className="truncate text-sm font-medium text-gray-900">{d.display_name || "Untitled Device"}</p>
                    <p className="text-xs text-gray-400 font-mono">{d.phone || "No phone linked"}</p>
                  </div>
                  <span
                    className={`rounded-full px-2 py-0.5 text-xs font-medium ${
                      (d.status === 1 || d.status === "connected" as any)
                        ? "bg-green-50 text-green-700"
                        : "bg-gray-100 text-gray-500"
                    }`}
                  >
                    {(d.status === 1 || d.status === "connected" as any) ? "Terhubung" : "Terputus"}
                  </span>
                </li>
              )) : (
                <li className="px-5 py-10 text-center text-sm text-gray-400">Belum ada device ditambahkan</li>
              )}
            </ul>
          </div>

          {/* Recent messages */}
          <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
            <div className="flex items-center justify-between border-b border-gray-100 px-5 py-4">
              <h2 className="font-semibold text-gray-900">Pesan Terbaru</h2>
              <a href="/messaging/logs" className="text-xs font-medium text-green-600 hover:underline">
                Lihat semua →
              </a>
            </div>
            <ul className="divide-y divide-gray-50 max-h-[320px] overflow-y-auto">
              {recentMessages.length > 0 ? recentMessages.map((m) => (
                <li key={m.id} className="flex items-start gap-3 px-5 py-3 transition-colors hover:bg-gray-50">
                  <div className="min-w-0 flex-1">
                    <p className="font-mono text-xs text-gray-400">{m.target_jid.split('@')[0]}</p>
                    <p className="truncate text-sm text-gray-700">{m.content}</p>
                  </div>
                  <div className="flex flex-shrink-0 flex-col items-end gap-1">
                    <MessageStatusBadge status={m.status} />
                    <span className="text-[10px] text-gray-400">
                      {new Date(m.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                    </span>
                  </div>
                </li>
              )) : (
                <li className="px-5 py-10 text-center text-sm text-gray-400">Belum ada aktivitas pesan</li>
              )}
            </ul>
          </div>
        </div>

        {/* ── Quick stats bar ─────────────────────────────────────────────── */}
        <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
          <div className="mb-4 flex items-center gap-2">
            <TrendingUp className="h-4 w-4 text-green-600" />
            <h2 className="font-semibold text-gray-900">Ringkasan Statistik</h2>
          </div>
          <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
            {[
              { label: "Total Terkirim", value: (stats?.total_sent || 0).toLocaleString(), color: "text-green-600" },
              { label: "Total Gagal", value: (stats?.total_failed || 0).toLocaleString(), color: "text-red-500" },
              { label: "Success Rate", value: `${(stats?.success_rate || 100).toFixed(1)}%`, color: "text-blue-600" },
              { label: "Devices Tracking", value: devices.length.toString(), color: "text-purple-600" },
            ].map((s) => (
              <div key={s.label} className="rounded-lg bg-gray-50 px-4 py-3">
                <p className={`text-2xl font-bold ${s.color}`}>{s.value}</p>
                <p className="mt-0.5 text-xs text-gray-500">{s.label}</p>
              </div>
            ))}
          </div>
        </div>

      </div>
    </div>
  );
}
