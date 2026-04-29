"use client";

import { useState, useEffect } from "react";
import { Wifi, WifiOff, RefreshCw, QrCode, Trash2, LogOut, Loader2, AlertTriangle } from "lucide-react";
import { useRouter } from "next/navigation";
import { sessionsApi } from "@/lib/api";
import { Device } from "@/lib/types";
import { toast } from "sonner";

type ActionType = "disconnect" | "delete" | null;

export default function SessionManagementPage() {
  const router = useRouter();
  const [devices, setDevices] = useState<Device[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [reconnecting, setReconnecting] = useState<string | null>(null);
  const [isProcessing, setIsProcessing] = useState<string | null>(null);
  
  // Custom Modal State
  const [confirm, setConfirm] = useState<{ deviceId: string; action: ActionType } | null>(null);

  const fetchDevices = () => {
    setIsLoading(true);
    sessionsApi.list()
      .then((res) => {
        setDevices(res.sessions || []);
      })
      .catch((err) => {
        toast.error("Gagal memuat daftar sesi.");
        console.error(err);
      })
      .finally(() => {
        setIsLoading(false);
      });
  };

  useEffect(() => {
    fetchDevices();
  }, []);

  const deviceToHandle = devices.find((d) => d.device_id === confirm?.deviceId);

  async function handleReconnect(id: string, name: string) {
    setReconnecting(id);
    try {
      const res = await sessionsApi.reconnect(id);
      if (res.success) {
        toast.success(res.message);
        // Refresh list after a bit
        setTimeout(fetchDevices, 3000);
      }
    } catch (err: any) {
      toast.info("Sesi lama tidak ditemukan atau kedaluwarsa. Membuka QR...");
      setTimeout(() => {
        router.push("/devices/qr");
      }, 1500);
    } finally {
      setReconnecting(null);
    }
  }

  async function executeAction() {
    if (!confirm || !confirm.deviceId) return;
    const { deviceId, action } = confirm;
    
    setIsProcessing(deviceId);
    try {
      if (action === "disconnect") {
        await sessionsApi.stop(deviceId);
        toast.success("Koneksi berhasil diputuskan.");
      } else if (action === "delete") {
        await sessionsApi.delete(deviceId);
        toast.success("Perangkat berhasil dihapus.");
      }
      setConfirm(null);
      fetchDevices();
    } catch (err: any) {
      toast.error(err.response?.data?.error || "Gagal memproses aksi.");
      console.error(err);
    } finally {
      setIsProcessing(null);
    }
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
            <h1 className="text-xl font-bold text-gray-900">Session Management</h1>
            <p className="text-sm text-gray-500">Monitor dan kelola semua sesi WhatsApp Anda</p>
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
          <div className="flex items-center gap-2 rounded-full bg-green-50 px-4 py-2 text-sm font-medium text-green-700 shadow-sm border border-green-100">
            <Wifi className="h-4 w-4" />
            {devices.filter((d) => d.status === 1).length} Terhubung
          </div>
          <div className="flex items-center gap-2 rounded-full bg-red-50 px-4 py-2 text-sm font-medium text-red-600 shadow-sm border border-red-100">
            <WifiOff className="h-4 w-4" />
            {devices.filter((d) => d.status === 0).length} Terputus
          </div>
        </div>

        {/* Table/List */}
        <div className="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
          <div className="overflow-x-auto">
            <table className="w-full min-w-[700px] text-sm text-left">
              <thead>
                <tr className="border-b border-gray-100 bg-gray-50 text-xs font-semibold uppercase tracking-wider text-gray-500">
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
                       <p className="mt-2 text-sm text-gray-400 font-medium">Memuat data device...</p>
                    </td>
                  </tr>
                ) : devices.length === 0 ? (
                  <tr>
                    <td colSpan={5} className="py-20 text-center">
                       <p className="text-sm text-gray-400">Tidak ada device yang terdaftar.</p>
                       <button 
                         onClick={() => router.push("/devices/qr")}
                         className="mt-4 text-sm font-semibold text-green-600 hover:text-green-700 hover:underline"
                       >
                         Tautkan Perangkat Sekarang
                       </button>
                    </td>
                  </tr>
                ) : (
                  devices.map((d) => (
                    <tr key={d.device_id} className="hover:bg-gray-50 transition-colors">
                      <td className="px-5 py-4 font-semibold text-gray-900">
                        {d.display_name || "Tanpa Nama"}
                      </td>
                      <td className="px-5 py-4 font-mono text-gray-500 text-xs">
                        {d.phone || "-"}
                      </td>
                      <td className="px-5 py-4">
                        {getStatusLabel(d.status)}
                      </td>
                      <td className="px-5 py-4 text-[10px] font-mono text-gray-400">
                        {d.device_id}
                      </td>
                      <td className="px-5 py-4 text-right">
                        <div className="flex items-center justify-end gap-2">
                          {d.status === 1 ? (
                            <button
                              onClick={() => setConfirm({ deviceId: d.device_id, action: "disconnect" })}
                              disabled={!!isProcessing}
                              className="inline-flex items-center gap-1.5 rounded-lg border border-red-100 bg-red-50 px-3 py-1.5 text-xs font-bold text-red-600 hover:bg-red-100 transition-colors shadow-sm"
                            >
                              {isProcessing === d.device_id ? (
                                <Loader2 className="h-3 w-3 animate-spin" />
                              ) : (
                                <LogOut className="h-3.5 w-3.5" />
                              )}
                              Putuskan Sesi
                            </button>
                          ) : d.status === 0 ? (
                            <>
                              <button
                                onClick={() => handleReconnect(d.device_id, d.display_name || d.device_id)}
                                disabled={!!reconnecting || !!isProcessing}
                                className="inline-flex items-center gap-1.5 rounded-lg bg-green-600 px-3 py-1.5 text-xs font-bold text-white hover:bg-green-700 transition-all shadow-sm shadow-green-100 active:scale-95"
                              >
                                <QrCode className="h-3.5 w-3.5" /> Reconnect
                              </button>
                              <button
                                onClick={() => setConfirm({ deviceId: d.device_id, action: "delete" })}
                                disabled={!!isProcessing}
                                className="inline-flex items-center gap-1.5 rounded-lg border border-red-200 bg-red-600 px-3 py-1.5 text-xs font-bold text-white hover:bg-red-700 transition-all shadow-sm shadow-red-100 active:scale-95"
                              >
                                {isProcessing === d.device_id ? (
                                  <Loader2 className="h-3 w-3 animate-spin" />
                                ) : (
                                  <Trash2 className="h-3.5 w-3.5" />
                                )}
                                Hapus
                              </button>
                            </>
                          ) : (
                            <button 
                              onClick={() => setConfirm({ deviceId: d.device_id, action: "disconnect" })}
                              disabled={!!isProcessing}
                              className="rounded-lg p-2 text-gray-400 hover:bg-gray-100 hover:text-red-500 transition-colors"
                            >
                              <LogOut className="h-4 w-4" />
                            </button>
                          )}
                        </div>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      {/* Custom Modal Confirmation */}
      {confirm && deviceToHandle && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in duration-200">
          <div className="w-full max-w-sm rounded-3xl bg-white p-6 shadow-2xl animate-in zoom-in-95 duration-200">
            <div className={`mb-4 flex h-14 w-14 items-center justify-center rounded-2xl ${confirm.action === 'delete' ? 'bg-red-100 text-red-600' : 'bg-orange-100 text-orange-600'}`}>
              <AlertTriangle className="h-7 w-7" />
            </div>
            
            <h3 className="text-xl font-bold text-gray-900">
              {confirm.action === "delete" ? "Hapus Perangkat?" : "Putuskan Sesi?"}
            </h3>
            
            <p className="mt-2 text-sm text-gray-600 leading-relaxed">
              {confirm.action === "delete" 
                ? `Anda akan menghapus data perangkat **${deviceToHandle.display_name || deviceToHandle.device_id}** secara permanen. Tindakan ini tidak dapat dibatalkan.`
                : `Anda akan memutus koneksi WhatsApp pada **${deviceToHandle.display_name || deviceToHandle.device_id}**. Anda perlu melakukan scan QR ulang untuk menghubungkannya kembali.`
              }
            </p>

            <div className="mt-8 flex gap-3">
              <button
                disabled={!!isProcessing}
                onClick={() => setConfirm(null)}
                className="flex-1 rounded-2xl border border-gray-200 py-3 text-sm font-bold text-gray-600 hover:bg-gray-50 transition-colors"
              >
                Batal
              </button>
              <button
                disabled={!!isProcessing}
                onClick={executeAction}
                className={`flex-1 rounded-2xl py-3 text-sm font-bold text-white shadow-lg transition-all active:scale-95 flex items-center justify-center gap-2 ${
                  confirm.action === "delete" ? "bg-red-600 hover:bg-red-700 shadow-red-100" : "bg-orange-500 hover:bg-orange-600 shadow-orange-100"
                }`}
              >
                {isProcessing ? (
                  <Loader2 className="h-4 w-4 animate-spin" />
                ) : (
                  confirm.action === "delete" ? "Ya, Hapus" : "Ya, Putuskan"
                )}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
