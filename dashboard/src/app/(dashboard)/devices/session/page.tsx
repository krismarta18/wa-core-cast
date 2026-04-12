"use client";

import { useState } from "react";
import { Power, RotateCcw, LogOut, AlertTriangle } from "lucide-react";

const DEVICES = [
  { id: 1, name: "Device Utama", phone: "628112345678", status: "connected" },
  { id: 2, name: "Device CS", phone: "628223456789", status: "connected" },
  { id: 3, name: "Device Marketing", phone: "628334567890", status: "disconnected" },
  { id: 4, name: "Device Backup", phone: "628445678901", status: "connected" },
];

type Action = "logout" | "restart" | null;

export default function DeviceSessionPage() {
  const [confirm, setConfirm] = useState<{ deviceId: number; action: Action } | null>(null);

  const device = DEVICES.find((d) => d.id === confirm?.deviceId);

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <h1 className="text-xl font-bold text-gray-900">Logout / Restart Device</h1>
        <p className="text-sm text-gray-500">Kelola sesi WhatsApp — logout atau restart koneksi device</p>
      </div>

      <div className="p-6">
        <div className="mb-4 rounded-xl border border-yellow-100 bg-yellow-50 p-4 text-sm text-yellow-800">
          <div className="flex items-start gap-2">
            <AlertTriangle className="mt-0.5 h-4 w-4 flex-shrink-0" />
            <span>
              <strong>Perhatian:</strong> Logout akan memutus WhatsApp dari device dan harus scan QR ulang.
              Restart hanya akan reconnect ulang tanpa logout.
            </span>
          </div>
        </div>

        <div className="space-y-3">
          {DEVICES.map((d) => (
            <div key={d.id} className="flex flex-wrap items-center justify-between gap-4 rounded-xl border border-gray-200 bg-white px-5 py-4 shadow-sm">
              <div className="flex items-center gap-3">
                <span
                  className={`h-2.5 w-2.5 rounded-full flex-shrink-0 ${
                    d.status === "connected" ? "bg-green-500" : "bg-gray-300"
                  }`}
                />
                <div>
                  <p className="font-semibold text-gray-900">{d.name}</p>
                  <p className="font-mono text-xs text-gray-400">{d.phone}</p>
                </div>
              </div>

              <div className="flex gap-2">
                <button
                  onClick={() => setConfirm({ deviceId: d.id, action: "restart" })}
                  disabled={d.status === "disconnected"}
                  className="inline-flex items-center gap-1.5 rounded-lg border border-blue-200 bg-blue-50 px-3 py-1.5 text-xs font-semibold text-blue-700 hover:bg-blue-100 disabled:cursor-not-allowed disabled:opacity-40 transition-colors"
                >
                  <RotateCcw className="h-3.5 w-3.5" /> Restart
                </button>
                <button
                  onClick={() => setConfirm({ deviceId: d.id, action: "logout" })}
                  className="inline-flex items-center gap-1.5 rounded-lg border border-red-200 bg-red-50 px-3 py-1.5 text-xs font-semibold text-red-700 hover:bg-red-100 transition-colors"
                >
                  <LogOut className="h-3.5 w-3.5" /> Logout
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Confirm dialog */}
      {confirm && device && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4">
          <div className="w-full max-w-sm rounded-2xl bg-white p-6 shadow-xl">
            <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-red-50">
              <Power className="h-6 w-6 text-red-500" />
            </div>
            <h3 className="text-lg font-bold text-gray-900">
              {confirm.action === "logout" ? "Logout Device?" : "Restart Device?"}
            </h3>
            <p className="mt-2 text-sm text-gray-600">
              {confirm.action === "logout"
                ? `WhatsApp di ${device.name} akan diputus dan harus scan QR ulang.`
                : `Koneksi ${device.name} akan direstart. Proses ini membutuhkan beberapa detik.`}
            </p>
            <div className="mt-5 flex gap-3">
              <button
                onClick={() => setConfirm(null)}
                className="flex-1 rounded-lg border border-gray-200 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-50"
              >
                Batal
              </button>
              <button
                onClick={() => setConfirm(null)}
                className={`flex-1 rounded-lg py-2 text-sm font-semibold text-white transition-colors ${
                  confirm.action === "logout" ? "bg-red-600 hover:bg-red-700" : "bg-blue-600 hover:bg-blue-700"
                }`}
              >
                {confirm.action === "logout" ? "Ya, Logout" : "Ya, Restart"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
