"use client";

import { useState } from "react";
import { Wifi, WifiOff, RefreshCw, QrCode, MoreVertical } from "lucide-react";
import { useRouter } from "next/navigation";

const INITIAL_DEVICES = [
  { id: 1, name: "Device Utama",    phone: "628112345678", status: "connected",    uptime: "3h 42m", lastSeen: "Baru saja",  battery: 87 },
  { id: 2, name: "Device CS",       phone: "628223456789", status: "connected",    uptime: "1h 15m", lastSeen: "Baru saja",  battery: 62 },
  { id: 3, name: "Device Marketing",phone: "628334567890", status: "disconnected", uptime: "-",      lastSeen: "2 jam lalu", battery: 0 },
  { id: 4, name: "Device Backup",   phone: "628445678901", status: "connected",    uptime: "5h 8m",  lastSeen: "Baru saja",  battery: 45 },
  { id: 5, name: "Device Support",  phone: "628556789012", status: "disconnected", uptime: "-",      lastSeen: "1 hari lalu",battery: 0 },
];

export default function ConnectionStatusPage() {
  const router = useRouter();
  const [devices, setDevices] = useState(INITIAL_DEVICES);
  const [reconnecting, setReconnecting] = useState<number | null>(null);
  const [toast, setToast] = useState<{ id: number; name: string } | null>(null);

  function handleReconnect(id: number, name: string) {
    setReconnecting(id);
    // Simulate reconnect attempt
    setTimeout(() => {
      setReconnecting(null);
      setToast({ id, name });
      setTimeout(() => setToast(null), 3000);
      router.push("/devices/qr");
    }, 1200);
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Toast */}
      {toast && (
        <div className="fixed bottom-6 right-6 z-50 flex items-center gap-3 rounded-xl border border-green-200 bg-white px-5 py-3 shadow-lg">
          <div className="flex h-8 w-8 items-center justify-center rounded-full bg-green-100">
            <QrCode className="h-4 w-4 text-green-600" />
          </div>
          <p className="text-sm font-medium text-gray-800">
            Membuka QR untuk <span className="text-green-600">{toast.name}</span>...
          </p>
        </div>
      )}

      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold text-gray-900">Connection Status</h1>
            <p className="text-sm text-gray-500">Monitor status koneksi semua device WhatsApp</p>
          </div>
          <button
            onClick={() => setDevices([...INITIAL_DEVICES])}
            className="inline-flex items-center gap-2 rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50"
          >
            <RefreshCw className="h-4 w-4" /> Refresh
          </button>
        </div>
      </div>

      <div className="p-6">
        {/* Summary pills */}
        <div className="mb-6 flex flex-wrap gap-3">
          <div className="flex items-center gap-2 rounded-full bg-green-50 px-4 py-2 text-sm font-medium text-green-700">
            <Wifi className="h-4 w-4" />
            {devices.filter((d) => d.status === "connected").length} Terhubung
          </div>
          <div className="flex items-center gap-2 rounded-full bg-red-50 px-4 py-2 text-sm font-medium text-red-600">
            <WifiOff className="h-4 w-4" />
            {devices.filter((d) => d.status === "disconnected").length} Terputus
          </div>
        </div>

        {/* Table */}
        <div className="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
          <div className="overflow-x-auto">
            <table className="w-full min-w-[700px] text-sm">
              <thead>
                <tr className="border-b border-gray-100 bg-gray-50 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">
                  <th className="px-5 py-3">Device</th>
                  <th className="px-5 py-3">Nomor</th>
                  <th className="px-5 py-3">Status</th>
                  <th className="px-5 py-3">Uptime</th>
                  <th className="px-5 py-3">Terakhir Aktif</th>
                  <th className="px-5 py-3">Baterai</th>
                  <th className="px-5 py-3">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-50">
                {devices.map((d) => (
                  <tr key={d.id} className="hover:bg-gray-50">
                    <td className="px-5 py-3 font-medium text-gray-900">{d.name}</td>
                    <td className="px-5 py-3 font-mono text-gray-500">{d.phone}</td>
                    <td className="px-5 py-3">
                      {d.status === "connected" ? (
                        <span className="inline-flex items-center gap-1.5 rounded-full bg-green-50 px-2.5 py-1 text-xs font-medium text-green-700">
                          <span className="h-1.5 w-1.5 rounded-full bg-green-500" /> Terhubung
                        </span>
                      ) : (
                        <span className="inline-flex items-center gap-1.5 rounded-full bg-red-50 px-2.5 py-1 text-xs font-medium text-red-600">
                          <span className="h-1.5 w-1.5 rounded-full bg-red-400" /> Terputus
                        </span>
                      )}
                    </td>
                    <td className="px-5 py-3 text-gray-600">{d.uptime}</td>
                    <td className="px-5 py-3 text-gray-600">{d.lastSeen}</td>
                    <td className="px-5 py-3">
                      {d.battery > 0 ? (
                        <div className="flex items-center gap-2">
                          <div className="h-2 w-16 overflow-hidden rounded-full bg-gray-100">
                            <div
                              className={`h-full rounded-full ${d.battery >= 50 ? "bg-green-500" : d.battery >= 20 ? "bg-yellow-500" : "bg-red-500"}`}
                              style={{ width: `${d.battery}%` }}
                            />
                          </div>
                          <span className="text-xs text-gray-500">{d.battery}%</span>
                        </div>
                      ) : (
                        <span className="text-xs text-gray-400">-</span>
                      )}
                    </td>
                    <td className="px-5 py-3">
                      {d.status === "disconnected" ? (
                        <button
                          onClick={() => handleReconnect(d.id, d.name)}
                          disabled={reconnecting === d.id}
                          className="inline-flex items-center gap-1.5 rounded-lg bg-green-600 px-3 py-1.5 text-xs font-semibold text-white hover:bg-green-700 disabled:cursor-not-allowed disabled:opacity-60 transition-colors"
                        >
                          {reconnecting === d.id ? (
                            <>
                              <RefreshCw className="h-3.5 w-3.5 animate-spin" /> Menghubungkan...
                            </>
                          ) : (
                            <>
                              <QrCode className="h-3.5 w-3.5" /> Reconnect
                            </>
                          )}
                        </button>
                      ) : (
                        <button className="rounded-lg p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-600">
                          <MoreVertical className="h-4 w-4" />
                        </button>
                      )}
                    </td>
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
