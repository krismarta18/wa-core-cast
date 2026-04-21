"use client";

import { useEffect, useState } from "react";
import { AlertTriangle, AlertCircle, Clock, Loader2 } from "lucide-react";
import { analyticsApi } from "@/lib/api";
import { FailureRateResponse } from "@/lib/types";

const TYPE_COLOR: Record<string, string> = {
  not_registered: "bg-gray-100 text-gray-600",
  session_error: "bg-red-50 text-red-700",
  rate_limit: "bg-orange-50 text-orange-700",
  blocked: "bg-red-100 text-red-800",
  timeout: "bg-yellow-50 text-yellow-700",
  send_failed: "bg-red-50 text-red-600",
};

export default function FailureRatePage() {
  const [data, setData] = useState<FailureRateResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchData() {
      try {
        setLoading(true);
        const res = await analyticsApi.failures();
        setData(res);
      } catch (err) {
        console.error("Failed to fetch failure analytics", err);
        setError("Gagal mengambil data analisis kegagalan");
      } finally {
        setLoading(false);
      }
    }
    fetchData();
  }, []);

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50">
        <Loader2 className="h-8 w-8 animate-spin text-green-600" />
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50">
        <p className="text-red-500">{error || "Terjadi kesalahan"}</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 pb-10">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <h1 className="text-xl font-bold text-gray-900">Failure Rate</h1>
        <p className="text-sm text-gray-500">Analisis kegagalan pengiriman pesan</p>
      </div>

      <div className="p-6 space-y-5">
        {/* Summary */}
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
          {[
            { label: "Total Gagal (7 hari)", value: data.total_failed.toString(), icon: AlertCircle, color: "text-red-500", bg: "bg-red-50" },
            { label: "Failure Rate", value: `${data.failure_rate}%`, icon: AlertTriangle, color: "text-orange-500", bg: "bg-orange-50" },
            { label: "Avg. Retry Time", value: data.avg_retry_time, icon: Clock, color: "text-blue-600", bg: "bg-blue-50" },
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
            {data.reason_stats.length > 0 ? (
              data.reason_stats.map((r) => (
                <div key={r.reason}>
                  <div className="flex justify-between text-sm mb-1">
                    <span className="text-gray-700 capitalize">{r.reason.replace(/_/g, " ")}</span>
                    <div className="flex items-center gap-2">
                      <span className="text-gray-500">{r.count} kali</span>
                      <span className="font-semibold text-red-600">{r.pct.toFixed(0)}%</span>
                    </div>
                  </div>
                  <div className="h-2 w-full overflow-hidden rounded-full bg-gray-100">
                    <div className="h-full rounded-full bg-red-400" style={{ width: `${r.pct}%` }} />
                  </div>
                </div>
              ))
            ) : (
              <p className="text-center py-4 text-sm text-gray-400">Belum ada data penyebab kegagalan</p>
            )}
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
                {data.latest_logs.length > 0 ? (
                  data.latest_logs.map((f) => (
                    <tr key={f.id} className="hover:bg-gray-50">
                      <td className="px-5 py-3 font-mono text-gray-600">{f.to}</td>
                      <td className="px-5 py-3 text-gray-600">{f.device}</td>
                      <td className="px-5 py-3">
                        <span className={`inline-flex items-center gap-1 rounded-full px-2.5 py-1 text-xs font-medium ${TYPE_COLOR[f.type] || "bg-gray-100"}`}>
                          <AlertCircle className="h-3 w-3" /> {f.reason || f.type}
                        </span>
                      </td>
                      <td className="px-5 py-3 text-gray-500">{f.time}</td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan={4} className="px-5 py-10 text-center text-gray-400">
                      Tidak ada log kegagalan ditemukan
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
}
