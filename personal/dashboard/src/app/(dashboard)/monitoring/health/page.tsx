"use client";

import { Server, Database, Wifi, HardDrive, CheckCircle2, AlertTriangle, XCircle, RefreshCw } from "lucide-react";

const SERVICES = [
  { name: "WhatsApp Gateway", status: "healthy", latency: "12ms", uptime: "99.9%", since: "10 hari" },
  { name: "Database (PostgreSQL)", status: "healthy", latency: "4ms", uptime: "100%", since: "30 hari" },
  { name: "Webhook Dispatcher", status: "warning", latency: "142ms", uptime: "97.2%", since: "—" },
  { name: "API Server", status: "healthy", latency: "8ms", uptime: "99.8%", since: "15 hari" },
];

const METRICS = [
  { label: "CPU Usage", value: 34, color: "bg-green-500", unit: "%" },
  { label: "Memory", value: 61, color: "bg-blue-500", unit: "%" },
  { label: "Disk", value: 45, color: "bg-purple-500", unit: "%" },
];

const DEVICES_POOL = [
  { name: "Device Utama", number: "+62 811-0000-001", status: "connected" },
  { name: "Device Marketing", number: "+62 811-0000-002", status: "connected" },
  { name: "Device CS", number: "+62 811-0000-003", status: "connected" },
  { name: "Device Backup", number: "+62 811-0000-004", status: "disconnected" },
];

const STATUS_ICON = {
  healthy: <CheckCircle2 className="h-4 w-4 text-green-500" />,
  warning: <AlertTriangle className="h-4 w-4 text-yellow-500" />,
  error: <XCircle className="h-4 w-4 text-red-500" />,
};
const STATUS_BADGE: Record<string, string> = {
  healthy: "bg-green-50 text-green-700",
  warning: "bg-yellow-50 text-yellow-700",
  error: "bg-red-50 text-red-700",
};

export default function SystemHealthPage() {
  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold text-gray-900">System Health</h1>
            <p className="text-sm text-gray-500">Status server dan komponen gateway</p>
          </div>
          <button className="flex items-center gap-2 rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm text-gray-600 hover:bg-gray-50">
            <RefreshCw className="h-4 w-4" /> Refresh
          </button>
        </div>
      </div>

      <div className="p-6 space-y-5">
        {/* Resource gauges */}
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
          {METRICS.map((m) => (
            <div key={m.label} className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
              <div className="flex items-center justify-between mb-3">
                <p className="text-sm font-medium text-gray-700">{m.label}</p>
                <span className="text-2xl font-bold text-gray-900">{m.value}{m.unit}</span>
              </div>
              <div className="h-3 w-full overflow-hidden rounded-full bg-gray-100">
                <div className={`h-full rounded-full ${m.color} transition-all`} style={{ width: `${m.value}%` }} />
              </div>
              <p className="mt-1 text-xs text-gray-400">{m.value < 70 ? "Normal" : m.value < 85 ? "Moderate" : "High"}</p>
            </div>
          ))}
        </div>

        {/* Services table */}
        <div className="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
          <div className="border-b border-gray-100 px-5 py-4">
            <h2 className="font-semibold text-gray-900">Status Layanan</h2>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full min-w-[600px] text-sm">
              <thead>
                <tr className="border-b border-gray-100 bg-gray-50 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">
                  <th className="px-5 py-3">Layanan</th>
                  <th className="px-5 py-3">Status</th>
                  <th className="px-5 py-3">Latency</th>
                  <th className="px-5 py-3">Uptime</th>
                  <th className="px-5 py-3">Uptime Sejak</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-50">
                {SERVICES.map((s) => (
                  <tr key={s.name} className="hover:bg-gray-50">
                    <td className="px-5 py-3 font-medium text-gray-800">{s.name}</td>
                    <td className="px-5 py-3">
                      <span className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-medium ${STATUS_BADGE[s.status]}`}>
                        {STATUS_ICON[s.status as keyof typeof STATUS_ICON]} {s.status.charAt(0).toUpperCase() + s.status.slice(1)}
                      </span>
                    </td>
                    <td className="px-5 py-3 font-mono text-gray-600">{s.latency}</td>
                    <td className="px-5 py-3 font-mono text-gray-600">{s.uptime}</td>
                    <td className="px-5 py-3 text-gray-500">{s.since}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* WhatsApp connection pool */}
        <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
          <div className="border-b border-gray-100 px-5 py-4">
            <h2 className="font-semibold text-gray-900">WhatsApp Connection Pool</h2>
          </div>
          <div className="divide-y divide-gray-50">
            {DEVICES_POOL.map((d) => (
              <div key={d.name} className="flex items-center justify-between px-5 py-3">
                <div className="flex items-center gap-3">
                  <Wifi className={`h-4 w-4 ${d.status === "connected" ? "text-green-500" : "text-gray-300"}`} />
                  <div>
                    <p className="text-sm font-medium text-gray-800">{d.name}</p>
                    <p className="text-xs text-gray-400">{d.number}</p>
                  </div>
                </div>
                <span className={`text-xs font-medium rounded-full px-2.5 py-1 ${d.status === "connected" ? "bg-green-50 text-green-700" : "bg-gray-100 text-gray-500"}`}>
                  {d.status === "connected" ? "Connected" : "Disconnected"}
                </span>
              </div>
            ))}
          </div>
        </div>

        {/* Last checked */}
        <p className="text-center text-xs text-gray-400">Terakhir dicek: 12 Apr 2026, 10:05:32 WIB</p>
      </div>
    </div>
  );
}
