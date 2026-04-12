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
} from "lucide-react";

// ─── Dummy data ───────────────────────────────────────────────────────────────
const PLAN = {
  name: "Business Pro",
  quotaUsed: 8_420,
  quotaLimit: 10_000,
  deviceUsed: 4,
  deviceMax: 10,
};

const RECENT_MESSAGES = [
  { id: 1, to: "628111222333", preview: "Halo, pesanan kamu sudah dikirim!", status: "sent", time: "2 mnt lalu" },
  { id: 2, to: "628222333444", preview: "Konfirmasi pembayaran berhasil ✅", status: "sent", time: "5 mnt lalu" },
  { id: 3, to: "628333444555", preview: "Promo akhir bulan, diskon 30%!", status: "failed", time: "12 mnt lalu" },
  { id: 4, to: "628444555666", preview: "Terima kasih sudah berbelanja 🛍️", status: "pending", time: "18 mnt lalu" },
  { id: 5, to: "628555666777", preview: "Kode OTP kamu: 847291", status: "sent", time: "25 mnt lalu" },
];

const DEVICES = [
  { id: 1, name: "Device Utama", phone: "628112345678", status: "connected" },
  { id: 2, name: "Device CS", phone: "628223456789", status: "connected" },
  { id: 3, name: "Device Marketing", phone: "628334567890", status: "disconnected" },
  { id: 4, name: "Device Backup", phone: "628445678901", status: "connected" },
];
// ─────────────────────────────────────────────────────────────────────────────

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

function StatusDot({ status }: { status: string }) {
  return (
    <span
      className={`inline-block h-2 w-2 rounded-full flex-shrink-0 ${
        status === "connected" ? "bg-green-500" : "bg-gray-300"
      }`}
    />
  );
}

function MessageStatusBadge({ status }: { status: string }) {
  if (status === "sent")
    return (
      <span className="inline-flex items-center gap-1 rounded-full bg-green-50 px-2 py-0.5 text-xs font-medium text-green-700">
        <CheckCircle2 className="h-3 w-3" /> Terkirim
      </span>
    );
  if (status === "failed")
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
  const connectedDevices = DEVICES.filter((d) => d.status === "connected").length;

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
            <p className="mt-2 text-2xl font-bold text-gray-900">{PLAN.name}</p>
            <p className="mt-1 text-xs font-medium text-green-600">● Aktif hingga 30 Apr 2026</p>
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
              {(PLAN.quotaLimit - PLAN.quotaUsed).toLocaleString()}
              <span className="ml-1 text-sm font-normal text-gray-400">sisa</span>
            </p>
            <QuotaBar used={PLAN.quotaUsed} max={PLAN.quotaLimit} />
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
              {PLAN.deviceUsed}
              <span className="text-lg font-normal text-gray-400"> / {PLAN.deviceMax}</span>
            </p>
            <QuotaBar used={PLAN.deviceUsed} max={PLAN.deviceMax} />
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
              <span className="text-lg font-normal text-gray-400"> / {PLAN.deviceUsed}</span>
            </p>
            <p className="mt-2 text-xs text-gray-500">
              {PLAN.deviceUsed - connectedDevices} device terputus
            </p>
          </div>
        </div>

        {/* ── Bottom row ───────────────────────────────────────────────────── */}
        <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
          {/* Device list */}
          <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
            <div className="flex items-center justify-between border-b border-gray-100 px-5 py-4">
              <h2 className="font-semibold text-gray-900">Status Device</h2>
              <a href="/devices/status" className="text-xs font-medium text-green-600 hover:underline">
                Lihat semua →
              </a>
            </div>
            <ul className="divide-y divide-gray-50">
              {DEVICES.map((d) => (
                <li key={d.id} className="flex items-center gap-3 px-5 py-3">
                  <StatusDot status={d.status} />
                  <div className="min-w-0 flex-1">
                    <p className="truncate text-sm font-medium text-gray-900">{d.name}</p>
                    <p className="text-xs text-gray-400">{d.phone}</p>
                  </div>
                  <span
                    className={`rounded-full px-2 py-0.5 text-xs font-medium ${
                      d.status === "connected"
                        ? "bg-green-50 text-green-700"
                        : "bg-gray-100 text-gray-500"
                    }`}
                  >
                    {d.status === "connected" ? "Terhubung" : "Putus"}
                  </span>
                </li>
              ))}
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
            <ul className="divide-y divide-gray-50">
              {RECENT_MESSAGES.map((m) => (
                <li key={m.id} className="flex items-start gap-3 px-5 py-3">
                  <div className="min-w-0 flex-1">
                    <p className="font-mono text-xs text-gray-400">{m.to}</p>
                    <p className="truncate text-sm text-gray-700">{m.preview}</p>
                  </div>
                  <div className="flex flex-shrink-0 flex-col items-end gap-1">
                    <MessageStatusBadge status={m.status} />
                    <span className="text-xs text-gray-400">{m.time}</span>
                  </div>
                </li>
              ))}
            </ul>
          </div>
        </div>

        {/* ── Quick stats bar ─────────────────────────────────────────────── */}
        <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
          <div className="mb-4 flex items-center gap-2">
            <TrendingUp className="h-4 w-4 text-green-600" />
            <h2 className="font-semibold text-gray-900">Statistik Hari Ini</h2>
          </div>
          <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
            {[
              { label: "Pesan Terkirim", value: "1.243", color: "text-green-600" },
              { label: "Pesan Gagal", value: "17", color: "text-red-500" },
              { label: "Broadcast Aktif", value: "3", color: "text-blue-600" },
              { label: "Response Rate", value: "98.6%", color: "text-purple-600" },
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
