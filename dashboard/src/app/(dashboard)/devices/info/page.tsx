"use client";

import { Smartphone, MessageSquare, Clock, Wifi, WifiOff } from "lucide-react";

const DEVICES = [
  { id: 1, name: "Device Utama", phone: "628112345678", status: "connected", msgSent: 4_321, msgReceived: 1_204, connectedSince: "10 Apr 2026, 08:00", platform: "Android 13", waVersion: "2.24.5.78" },
  { id: 2, name: "Device CS", phone: "628223456789", status: "connected", msgSent: 2_108, msgReceived: 876, connectedSince: "11 Apr 2026, 14:30", platform: "iOS 17", waVersion: "2.24.4.90" },
  { id: 3, name: "Device Marketing", phone: "628334567890", status: "disconnected", msgSent: 980, msgReceived: 210, connectedSince: "-", platform: "Android 12", waVersion: "-" },
  { id: 4, name: "Device Backup", phone: "628445678901", status: "connected", msgSent: 523, msgReceived: 98, connectedSince: "09 Apr 2026, 09:15", platform: "Android 14", waVersion: "2.24.5.80" },
];

export default function MultiDeviceInfoPage() {
  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <h1 className="text-xl font-bold text-gray-900">Multi Device Info</h1>
        <p className="text-sm text-gray-500">Detail lengkap semua device yang terdaftar</p>
      </div>

      <div className="p-6 space-y-4">
        {DEVICES.map((d) => (
          <div key={d.id} className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
            <div className="flex flex-wrap items-start justify-between gap-3">
              <div className="flex items-center gap-3">
                <div className={`flex h-10 w-10 items-center justify-center rounded-xl ${d.status === "connected" ? "bg-green-50" : "bg-gray-100"}`}>
                  <Smartphone className={`h-5 w-5 ${d.status === "connected" ? "text-green-600" : "text-gray-400"}`} />
                </div>
                <div>
                  <div className="flex items-center gap-2">
                    <h3 className="font-semibold text-gray-900">{d.name}</h3>
                    {d.status === "connected" ? (
                      <span className="inline-flex items-center gap-1 rounded-full bg-green-50 px-2 py-0.5 text-xs font-medium text-green-700">
                        <Wifi className="h-3 w-3" /> Terhubung
                      </span>
                    ) : (
                      <span className="inline-flex items-center gap-1 rounded-full bg-gray-100 px-2 py-0.5 text-xs font-medium text-gray-500">
                        <WifiOff className="h-3 w-3" /> Terputus
                      </span>
                    )}
                  </div>
                  <p className="font-mono text-sm text-gray-500">{d.phone}</p>
                </div>
              </div>
            </div>

            <div className="mt-4 grid grid-cols-2 gap-3 sm:grid-cols-4">
              <div className="rounded-lg bg-gray-50 px-4 py-3">
                <div className="flex items-center gap-1.5 text-xs text-gray-500 mb-1">
                  <MessageSquare className="h-3 w-3" /> Pesan Terkirim
                </div>
                <p className="text-lg font-bold text-gray-900">{d.msgSent.toLocaleString()}</p>
              </div>
              <div className="rounded-lg bg-gray-50 px-4 py-3">
                <div className="flex items-center gap-1.5 text-xs text-gray-500 mb-1">
                  <MessageSquare className="h-3 w-3" /> Pesan Diterima
                </div>
                <p className="text-lg font-bold text-gray-900">{d.msgReceived.toLocaleString()}</p>
              </div>
              <div className="rounded-lg bg-gray-50 px-4 py-3">
                <div className="flex items-center gap-1.5 text-xs text-gray-500 mb-1">
                  <Clock className="h-3 w-3" /> Terhubung Sejak
                </div>
                <p className="text-sm font-semibold text-gray-900">{d.connectedSince}</p>
              </div>
              <div className="rounded-lg bg-gray-50 px-4 py-3">
                <div className="text-xs text-gray-500 mb-1">Platform</div>
                <p className="text-sm font-semibold text-gray-900">{d.platform}</p>
                <p className="text-xs text-gray-400">WA {d.waVersion}</p>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
