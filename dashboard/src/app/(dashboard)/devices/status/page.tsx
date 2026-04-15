"use client";

import { useState, useEffect } from "react";
import { Wifi, WifiOff, RefreshCw, QrCode, MoreVertical, Loader2 } from "lucide-react";
import { useRouter } from "next/navigation";
import { sessionsApi } from "@/lib/api";
import { Device } from "@/lib/types";
import { toast } from "sonner";

export default function ConnectionStatusPage() {
  const router = useRouter();
  const [devices, setDevices] = useState<Device[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [reconnecting, setReconnecting] = useState<string | null>(null);

  const fetchDevices = () => {
    setIsLoading(true);
    sessionsApi.list()
      .then((res) => {
        setDevices(res.sessions || []);
      })
      .catch((err) => {
        toast.error("Gagal memuat status koneksi.");
        console.error(err);
      })
      .finally(() => {
        setIsLoading(false);
      });
  };

  useEffect(() => {
    fetchDevices();
  }, []);

  function handleReconnect(id: string, name: string) {
    setReconnecting(id);
    toast.info(`Membuka QR untuk ${name || id}...`);
    setTimeout(() => {
      router.push("/devices/qr");
    }, 1000);
  }

  const getStatusLabel = (status: number) => {
    switch (status) {
      case 1:
        return (
          <span className="inline-flex items-center gap-1.5 rounded-full bg-green-50 px-2.5 py-1 text-xs font-medium text-green-700">
            <span className="h-1.5 w-1.5 rounded-full bg-green-500" /> Terhubung
          </span>
        );
      case 2:
        return (
          <span className="inline-flex items-center gap-1.5 rounded-full bg-yellow-50 px-2.5 py-1 text-xs font-medium text-yellow-700">
            <span className="h-1.5 w-1.5 rounded-full bg-yellow-500" /> Menunggu QR
          </span>
        );
      default:
        return (
          <span className="inline-flex items-center gap-1.5 rounded-full bg-red-50 px-2.5 py-1 text-xs font-medium text-red-600">
            <span className="h-1.5 w-1.5 rounded-full bg-red-400" /> Terputus
          </span>
        );
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold text-gray-900">Connection Status</h1>
            <p className="text-sm text-gray-500">Monitor status koneksi semua device WhatsApp</p>
          </div>
          <button
            onClick={fetchDevices}
            disabled={isLoading}
            className="inline-flex items-center gap-2 rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 disabled:opacity-50"
          >
            {isLoading ? <RefreshCw className="h-4 w-4 animate-spin" /> : <RefreshCw className="h-4 w-4" />} Refresh
          </button>
        </div>
      </div>

      <div className="p-6">
        {/* Summary pills */}
        <div className="mb-6 flex flex-wrap gap-3">
          <div className="flex items-center gap-2 rounded-full bg-green-50 px-4 py-2 text-sm font-medium text-green-700">
            <Wifi className="h-4 w-4" />
            {devices.filter((d) => d.status === 1).length} Terhubung
          </div>
          <div className="flex items-center gap-2 rounded-full bg-red-50 px-4 py-2 text-sm font-medium text-red-600">
            <WifiOff className="h-4 w-4" />
            {devices.filter((d) => d.status === 0).length} Terputus
          </div>
        </div>

        {/* Table/List */}
        <div className="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
          <div className="overflow-x-auto">
            <table className="w-full min-w-[700px] text-sm">
              <thead>
                <tr className="border-b border-gray-100 bg-gray-50 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">
                  <th className="px-5 py-3">Device Name</th>
                  <th className="px-5 py-3">WhatsApp Number</th>
                  <th className="px-5 py-3">Status</th>
                  <th className="px-5 py-3">Device ID</th>
                  <th className="px-5 py-3 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-50">
                {isLoading && devices.length === 0 ? (
                  <tr>
                    <td colSpan={5} className="py-20 text-center">
                       <Loader2 className="mx-auto h-8 w-8 animate-spin text-gray-300" />
                       <p className="mt-2 text-sm text-gray-400">Memuat data device...</p>
                    </td>
                  </tr>
                ) : devices.length === 0 ? (
                  <tr>
                    <td colSpan={5} className="py-20 text-center">
                       <p className="text-sm text-gray-400">Tidak ada device yang terdaftar.</p>
                       <button 
                         onClick={() => router.push("/devices/qr")}
                         className="mt-4 text-sm font-semibold text-green-600 hover:underline"
                       >
                         Tautkan Perangkat Sekarang
                       </button>
                    </td>
                  </tr>
                ) : (
                  devices.map((d) => (
                    <tr key={d.device_id} className="hover:bg-gray-50 transition-colors">
                      <td className="px-5 py-4 font-medium text-gray-900">
                        {d.display_name || "Tanpa Nama"}
                      </td>
                      <td className="px-5 py-4 font-mono text-gray-500">
                        {d.phone || "-"}
                      </td>
                      <td className="px-5 py-4">
                        {getStatusLabel(d.status)}
                      </td>
                      <td className="px-5 py-4 text-xs font-mono text-gray-400">
                        {d.device_id}
                      </td>
                      <td className="px-5 py-4 text-right">
                        {d.status === 0 ? (
                          <button
                            onClick={() => handleReconnect(d.device_id, d.display_name || d.device_id)}
                            disabled={!!reconnecting}
                            className="inline-flex items-center gap-1.5 rounded-lg bg-green-600 px-3 py-1.5 text-xs font-semibold text-white hover:bg-green-700 transition-colors"
                          >
                            <QrCode className="h-3.5 w-3.5" /> Reconnect
                          </button>
                        ) : (
                          <button className="rounded-lg p-2 text-gray-400 hover:bg-gray-100 hover:text-gray-600 transition-colors">
                            <MoreVertical className="h-4 w-4" />
                          </button>
                        )}
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
}
