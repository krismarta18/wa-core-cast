"use client";

import { AlertTriangle, AlertCircle, Clock } from "lucide-react";

const FAILURES = [
  { id: 1, to: "628111222333", device: "Device Marketing", reason: "Number not on WhatsApp", time: "12 Apr 2026, 10:03", type: "not_registered" },
  { id: 2, to: "628222333444", device: "Device Backup", reason: "Session disconnected", time: "12 Apr 2026, 09:55", type: "session_error" },
  { id: 3, to: "628333444555", device: "Device Utama", reason: "Rate limit exceeded", time: "12 Apr 2026, 09:40", type: "rate_limit" },
  { id: 4, to: "628444555666", device: "Device Marketing", reason: "Number blocked by WhatsApp", time: "12 Apr 2026, 09:22", type: "blocked" },
  { id: 5, to: "628555666777", device: "Device CS", reason: "Network timeout", time: "12 Apr 2026, 08:59", type: "timeout" },
  { id: 6, to: "628666777888", device: "Device Backup", reason: "Number not on WhatsApp", time: "12 Apr 2026, 08:44", type: "not_registered" },
];

const REASON_STATS = [
  { reason: "Number not on WhatsApp", count: 42, pct: 48 },
  { reason: "Session disconnected", count: 21, pct: 24 },
  { reason: "Rate limit exceeded", count: 15, pct: 17 },
  { reason: "Network timeout", count: 8, pct: 9 },
  { reason: "Number blocked", count: 2, pct: 2 },
];

const TYPE_COLOR: Record<string, string> = {
  not_registered: "bg-gray-100 text-gray-600",
  session_error: "bg-red-50 text-red-700",
  rate_limit: "bg-orange-50 text-orange-700",
  blocked: "bg-red-100 text-red-800",
  timeout: "bg-yellow-50 text-yellow-700",
};

export default function FailureRatePage() {
  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <h1 className="text-xl font-bold text-gray-900">Failure Rate</h1>
        <p className="text-sm text-gray-500">Analisis kegagalan pengiriman pesan</p>
      </div>

      <div className="p-6 space-y-5">
        {/* Summary */}
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
          {[
            { label: "Total Gagal (7 hari)", value: "88", icon: AlertCircle, color: "text-red-500", bg: "bg-red-50" },
            { label: "Failure Rate", value: "1.4%", icon: AlertTriangle, color: "text-orange-500", bg: "bg-orange-50" },
            { label: "Avg. Retry Time", value: "2.3s", icon: Clock, color: "text-blue-600", bg: "bg-blue-50" },
          ].map((s) => (
            <div key={s.label} className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
              <div className="flex items-center justify-between">
                <p className="text-sm font-medium text-gray-500">{s.label}</p>
                <div className={`flex h-9 w-9 items-center justify-center rounded-lg ${s.bg}`}>
                  <s.icon className={`h-5 w-5 ${s.color}`} />
                </div>
              </div>
              <p className={`mt-2 text-3xl font-bold ${s.color}`}>{s.value}</p>
            </div>
          ))}
        </div>

        {/* Reason breakdown */}
        <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
          <h2 className="mb-4 font-semibold text-gray-900">Penyebab Kegagalan</h2>
          <div className="space-y-3">
            {REASON_STATS.map((r) => (
              <div key={r.reason}>
                <div className="flex justify-between text-sm mb-1">
                  <span className="text-gray-700">{r.reason}</span>
                  <div className="flex items-center gap-2">
                    <span className="text-gray-500">{r.count} kali</span>
                    <span className="font-semibold text-red-600">{r.pct}%</span>
                  </div>
                </div>
                <div className="h-2 w-full overflow-hidden rounded-full bg-gray-100">
                  <div className="h-full rounded-full bg-red-400" style={{ width: `${r.pct}%` }} />
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Failure log */}
        <div className="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
          <div className="border-b border-gray-100 px-5 py-4">
            <h2 className="font-semibold text-gray-900">Log Kegagalan Terbaru</h2>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full min-w-[500px] text-sm">
              <thead>
                <tr className="border-b border-gray-100 bg-gray-50 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">
                  <th className="px-5 py-3">Tujuan</th>
                  <th className="px-5 py-3">Device</th>
                  <th className="px-5 py-3">Alasan</th>
                  <th className="px-5 py-3">Waktu</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-50">
                {FAILURES.map((f) => (
                  <tr key={f.id} className="hover:bg-gray-50">
                    <td className="px-5 py-3 font-mono text-gray-600">{f.to}</td>
                    <td className="px-5 py-3 text-gray-600">{f.device}</td>
                    <td className="px-5 py-3">
                      <span className={`inline-flex items-center gap-1 rounded-full px-2.5 py-1 text-xs font-medium ${TYPE_COLOR[f.type]}`}>
                        <AlertCircle className="h-3 w-3" /> {f.reason}
                      </span>
                    </td>
                    <td className="px-5 py-3 text-gray-500">{f.time}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
}
