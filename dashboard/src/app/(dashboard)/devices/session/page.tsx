"use client";

import { useState, useEffect } from "react";
import { Power, RotateCcw, LogOut, AlertTriangle, Loader2, Smartphone, Trash2 } from "lucide-react";
import { sessionsApi } from "@/lib/api";
import { Device } from "@/lib/types";
import { toast } from "sonner";

type Action = "logout" | "restart" | null;

export default function DeviceSessionPage() {
  const [devices, setDevices] = useState<Device[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [confirm, setConfirm] = useState<{ deviceId: string; action: Action } | null>(null);
  const [isProcessing, setIsProcessing] = useState(false);

  const fetchDevices = () => {
    setIsLoading(true);
    sessionsApi.list()
      .then((res) => {
        setDevices(res.sessions || []);
      })
      .catch((err) => {
        toast.error("Gagal memuat sesi.");
        console.error(err);
      })
      .finally(() => {
        setIsLoading(false);
      });
  };

  useEffect(() => {
    fetchDevices();
  }, []);

  const device = devices.find((d) => d.device_id === confirm?.deviceId);

  const handleAction = async () => {
    if (!confirm || !device) return;
    
    setIsProcessing(true);
    try {
      if (confirm.action === "logout") {
        await sessionsApi.stop(confirm.deviceId);
        toast.success(`Sesi ${device.display_name || device.device_id} berhasil dihentikan.`);
      } else {
        // Restart logic (currently fallback to stop/start or just notify not supported yet)
        toast.info("Fitur Restart sedang dalam pengembangan. Silakan gunakan Logout lalu hubungkan kembali.");
      }
      setConfirm(null);
      fetchDevices(); // Refresh list
    } catch (err: any) {
      toast.error(err.response?.data?.error || "Gagal memproses aksi.");
    } finally {
      setIsProcessing(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <h1 className="text-xl font-bold text-gray-900">Session Management</h1>
        <p className="text-sm text-gray-500">Kelola sesi aktif dan hentikan koneksi perangkat</p>
      </div>

      <div className="p-6">
        <div className="mb-6 rounded-xl border border-yellow-100 bg-yellow-50 p-4 text-sm text-yellow-800">
          <div className="flex items-start gap-2">
            <AlertTriangle className="mt-0.5 h-4 w-4 flex-shrink-0" />
            <span>
              <strong>Peringatan:</strong> Menghentikan sesi (Logout) akan memutus koneksi WhatsApp. 
              Beberapa data lokal mungkin akan dihapus untuk memastikan keamanan privasi.
            </span>
          </div>
        </div>

        {isLoading && devices.length === 0 ? (
          <div className="flex justify-center py-20">
            <Loader2 className="h-8 w-8 animate-spin text-gray-300" />
          </div>
        ) : devices.length === 0 ? (
          <div className="bg-white rounded-xl border border-gray-200 p-12 text-center">
            <Power className="mx-auto h-12 w-12 text-gray-200 mb-3" />
            <p className="text-gray-400">Tidak ada sesi aktif yang ditemukan</p>
          </div>
        ) : (
          <div className="grid gap-4">
            {devices.map((d) => (
              <div key={d.device_id} className="flex flex-wrap items-center justify-between gap-4 rounded-xl border border-gray-200 bg-white px-6 py-5 shadow-sm transition-all hover:border-gray-300">
                <div className="flex items-center gap-4">
                  <div className={`flex h-10 w-10 items-center justify-center rounded-full ${d.status === 1 ? "bg-green-100 text-green-600" : "bg-gray-100 text-gray-400"}`}>
                    <Smartphone className="h-5 w-5" />
                  </div>
                  <div>
                    <p className="font-bold text-gray-900">{d.display_name || "Device WhatsApp"}</p>
                    <div className="flex items-center gap-2 mt-0.5">
                       <span className={`h-1.5 w-1.5 rounded-full ${d.status === 1 ? "bg-green-500" : "bg-gray-300"}`} />
                       <p className="font-mono text-[10px] text-gray-400 uppercase tracking-wider">
                          {d.phone || d.device_id}
                       </p>
                    </div>
                  </div>
                </div>

                <div className="flex gap-2">
                  <button
                    onClick={() => setConfirm({ deviceId: d.device_id, action: "logout" })}
                    className="inline-flex items-center gap-2 rounded-lg border border-red-100 bg-red-50 px-4 py-2 text-xs font-bold text-red-600 hover:bg-red-100 hover:text-red-700 transition-colors"
                  >
                    <LogOut className="h-3.5 w-3.5" /> Stop Session
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Confirm dialog */}
      {confirm && device && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm px-4">
          <div className="w-full max-w-sm rounded-2xl bg-white p-6 shadow-2xl animate-in zoom-in-95 duration-200">
            <div className="mb-4 flex h-14 w-14 items-center justify-center rounded-full bg-red-100">
              <AlertTriangle className="h-7 w-7 text-red-600" />
            </div>
            <h3 className="text-xl font-bold text-gray-900">
              Hentikan Sesi?
            </h3>
            <p className="mt-2 text-sm text-gray-600 leading-relaxed">
              Anda akan memutus koneksi <strong>{device.display_name || device.device_id}</strong>. 
              Sesi ini akan dihapus dari server dan Anda perlu melakukan scan QR ulang untuk menghubungkannya kembali.
            </p>
            <div className="mt-6 flex gap-3">
              <button
                disabled={isProcessing}
                onClick={() => setConfirm(null)}
                className="flex-1 rounded-xl border border-gray-200 py-3 text-sm font-bold text-gray-700 hover:bg-gray-50 disabled:opacity-50 transition-colors"
              >
                Batal
              </button>
              <button
                disabled={isProcessing}
                onClick={handleAction}
                className="flex-1 rounded-xl bg-red-600 py-3 text-sm font-bold text-white hover:bg-red-700 shadow-lg shadow-red-200 disabled:opacity-50 transition-all"
              >
                {isProcessing ? (
                  <Loader2 className="mx-auto h-5 w-5 animate-spin" />
                ) : (
                  "Ya, Hentikan"
                )}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
