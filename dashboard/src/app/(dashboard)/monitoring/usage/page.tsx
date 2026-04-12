"use client";

import { BarChart2, TrendingUp, TrendingDown } from "lucide-react";

const DAILY = [
  { day: "Sen", sent: 1420, failed: 12 },
  { day: "Sel", sent: 1850, failed: 8 },
  { day: "Rab", sent: 2100, failed: 22 },
  { day: "Kam", sent: 1760, failed: 15 },
  { day: "Jum", sent: 2340, failed: 9 },
  { day: "Sab", sent: 980, failed: 5 },
  { day: "Min", sent: 420, failed: 3 },
];

const maxSent = Math.max(...DAILY.map((d) => d.sent));

const DEVICE_STATS = [
  { name: "Device Utama", sent: 4321, success: 98.7 },
  { name: "Device CS", sent: 2108, success: 99.1 },
  { name: "Device Marketing", sent: 980, success: 94.3 },
  { name: "Device Backup", sent: 523, success: 96.8 },
];

export default function UsageStatisticsPage() {
  const totalSent = DAILY.reduce((s, d) => s + d.sent, 0);
  const totalFailed = DAILY.reduce((s, d) => s + d.failed, 0);
  const successRate = (((totalSent - totalFailed) / totalSent) * 100).toFixed(1);

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <h1 className="text-xl font-bold text-gray-900">Usage Statistics</h1>
        <p className="text-sm text-gray-500">Statistik penggunaan pesan 7 hari terakhir</p>
      </div>

      <div className="p-6 space-y-6">
        {/* Summary */}
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
          {[
            { label: "Total Terkirim", value: totalSent.toLocaleString(), icon: TrendingUp, color: "text-green-600", bg: "bg-green-50" },
            { label: "Total Gagal", value: totalFailed.toString(), icon: TrendingDown, color: "text-red-500", bg: "bg-red-50" },
            { label: "Success Rate", value: `${successRate}%`, icon: BarChart2, color: "text-blue-600", bg: "bg-blue-50" },
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

        {/* Bar chart */}
        <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
          <h2 className="mb-5 font-semibold text-gray-900">Pesan Terkirim per Hari</h2>
          <div className="flex items-end justify-between gap-3 h-40">
            {DAILY.map((d) => (
              <div key={d.day} className="flex flex-1 flex-col items-center gap-1.5">
                <div
                  className="w-full rounded-t-md bg-green-500 transition-all"
                  style={{ height: `${(d.sent / maxSent) * 128}px` }}
                  title={`${d.sent} pesan`}
                />
                <span className="text-xs text-gray-500">{d.day}</span>
              </div>
            ))}
          </div>
        </div>

        {/* Per device */}
        <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
          <h2 className="mb-4 font-semibold text-gray-900">Statistik per Device</h2>
          <div className="space-y-4">
            {DEVICE_STATS.map((d) => (
              <div key={d.name}>
                <div className="flex items-center justify-between text-sm mb-1">
                  <span className="font-medium text-gray-900">{d.name}</span>
                  <div className="flex items-center gap-3">
                    <span className="text-gray-500">{d.sent.toLocaleString()} pesan</span>
                    <span className={`font-semibold ${d.success >= 98 ? "text-green-600" : d.success >= 95 ? "text-yellow-600" : "text-red-500"}`}>
                      {d.success}%
                    </span>
                  </div>
                </div>
                <div className="h-2 w-full overflow-hidden rounded-full bg-gray-100">
                  <div
                    className={`h-full rounded-full ${d.success >= 98 ? "bg-green-500" : d.success >= 95 ? "bg-yellow-500" : "bg-red-500"}`}
                    style={{ width: `${d.success}%` }}
                  />
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
