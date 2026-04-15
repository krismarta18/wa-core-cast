"use client";

import { useEffect, useState } from "react";
import { BarChart2, TrendingUp, TrendingDown, Loader2 } from "lucide-react";
import { analyticsApi } from "@/lib/api";
import { UsageStatsResponse } from "@/lib/types";

export default function UsageStatisticsPage() {
  const [data, setData] = useState<UsageStatsResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchData() {
      try {
        setLoading(true);
        const res = await analyticsApi.usage();
        setData(res);
      } catch (err) {
        console.error("Failed to fetch usage stats", err);
        setError("Gagal mengambil data statistik");
      } finally {
        setLoading(false);
      }
    }
    fetchData();
  }, []);

  const formatDay = (dateStr: string) => {
    const days = ["Min", "Sen", "Sel", "Rab", "Kam", "Jum", "Sab"];
    const d = new Date(dateStr);
    return days[d.getDay()];
  };

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

  const daily = data.daily || [];
  const maxSent = Math.max(...daily.map((d) => d.sent_count), 1);
  const successRate = data.success_rate.toFixed(1);

  return (
    <div className="min-h-screen bg-gray-50 pb-10">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <h1 className="text-xl font-bold text-gray-900">Usage Statistics</h1>
        <p className="text-sm text-gray-500">Statistik penggunaan pesan 7 hari terakhir</p>
      </div>

      <div className="p-6 space-y-6">
        {/* Summary */}
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
          {[
            { label: "Total Terkirim", value: data.total_sent.toLocaleString(), icon: TrendingUp, color: "text-green-600", bg: "bg-green-50" },
            { label: "Total Gagal", value: data.total_failed.toLocaleString(), icon: TrendingDown, color: "text-red-500", bg: "bg-red-50" },
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
          {daily.length > 0 ? (
            <div className="flex items-end justify-between gap-3 h-40">
              {daily.map((d) => (
                <div key={d.id} className="flex flex-1 flex-col items-center gap-1.5">
                  <div
                    className="w-full rounded-t-md bg-green-500 transition-all hover:bg-green-600"
                    style={{ height: `${(d.sent_count / maxSent) * 128}px` }}
                    title={`${d.sent_count} pesan terkirim`}
                  />
                  <span className="text-xs text-gray-500">{formatDay(d.stat_date)}</span>
                </div>
              ))}
            </div>
          ) : (
            <div className="flex h-40 items-center justify-center text-sm text-gray-400">
              Belum ada data statistik untuk 7 hari terakhir
            </div>
          )}
        </div>

        {/* Per device */}
        <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
          <h2 className="mb-4 font-semibold text-gray-900">Statistik per Device</h2>
          <div className="space-y-4">
            {data.device_stats.length > 0 ? (
              data.device_stats.map((d) => (
                <div key={d.name}>
                  <div className="flex items-center justify-between text-sm mb-1">
                    <span className="font-medium text-gray-900">{d.name}</span>
                    <div className="flex items-center gap-3">
                      <span className="text-gray-500">{d.sent.toLocaleString()} pesan</span>
                      <span className={`font-semibold ${d.success >= 98 ? "text-green-600" : d.success >= 95 ? "text-yellow-600" : "text-red-500"}`}>
                        {d.success.toFixed(1)}%
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
              ))
            ) : (
              <p className="text-center py-4 text-sm text-gray-400">Belum ada data device</p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
