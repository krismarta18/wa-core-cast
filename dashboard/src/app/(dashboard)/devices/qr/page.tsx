"use client";

import { useState } from "react";
import { QrCode, Smartphone, CheckCircle2, RefreshCw, Info } from "lucide-react";

const STEPS = [
  "Buka WhatsApp di HP kamu",
  "Ketuk ikon titik tiga (⋮) di pojok kanan atas",
  "Pilih \"Perangkat Tertaut\" lalu \"Tautkan Perangkat\"",
  "Arahkan kamera ke QR code di bawah ini",
];

export default function QRScannerPage() {
  const [deviceName, setDeviceName] = useState("");
  const [showQR, setShowQR] = useState(false);
  const [connected, setConnected] = useState(false);

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <h1 className="text-xl font-bold text-gray-900">Connect New Device</h1>
        <p className="text-sm text-gray-500">Hubungkan WhatsApp baru dengan scan QR code</p>
      </div>

      <div className="p-6">
        <div className="mx-auto max-w-xl">
          {!connected ? (
            <div className="space-y-5">
              {/* Step guide */}
              <div className="rounded-xl border border-blue-100 bg-blue-50 p-5">
                <div className="mb-3 flex items-center gap-2 text-sm font-semibold text-blue-700">
                  <Info className="h-4 w-4" /> Cara Menghubungkan Device
                </div>
                <ol className="space-y-2">
                  {STEPS.map((s, i) => (
                    <li key={i} className="flex items-start gap-3 text-sm text-blue-800">
                      <span className="flex h-5 w-5 flex-shrink-0 items-center justify-center rounded-full bg-blue-200 text-xs font-bold">
                        {i + 1}
                      </span>
                      {s}
                    </li>
                  ))}
                </ol>
              </div>

              {/* Form */}
              <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
                <label className="mb-1 block text-sm font-medium text-gray-700">
                  Nama Device <span className="text-red-500">*</span>
                </label>
                <input
                  type="text"
                  placeholder="cth: Device Utama, Device CS, Backup..."
                  value={deviceName}
                  onChange={(e) => setDeviceName(e.target.value)}
                  className="h-10 w-full rounded-lg border border-gray-300 px-3 text-sm placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-green-500"
                />

                <button
                  disabled={!deviceName.trim()}
                  onClick={() => setShowQR(true)}
                  className="mt-4 w-full rounded-lg bg-green-600 py-2.5 text-sm font-semibold text-white hover:bg-green-700 disabled:cursor-not-allowed disabled:bg-gray-200 disabled:text-gray-400 transition-colors"
                >
                  <QrCode className="mr-2 inline h-4 w-4" />
                  Generate QR Code
                </button>
              </div>

              {/* QR Display */}
              {showQR && (
                <div className="rounded-xl border border-gray-200 bg-white p-8 text-center shadow-sm">
                  <p className="mb-1 text-sm font-medium text-gray-700">
                    Scan QR untuk <span className="font-bold text-green-700">{deviceName}</span>
                  </p>
                  <p className="mb-6 text-xs text-gray-400">QR code berlaku selama 60 detik</p>

                  {/* Placeholder QR — akan diganti WebSocket real QR */}
                  <div className="mx-auto flex h-52 w-52 items-center justify-center rounded-xl border-2 border-dashed border-gray-300 bg-gray-50">
                    <div className="text-center">
                      <QrCode className="mx-auto h-16 w-16 text-gray-300" />
                      <p className="mt-2 text-xs text-gray-400">QR akan muncul di sini</p>
                    </div>
                  </div>

                  <button className="mt-5 inline-flex items-center gap-2 text-sm font-medium text-green-600 hover:underline">
                    <RefreshCw className="h-4 w-4" /> Refresh QR
                  </button>

                  {/* Dev shortcut */}
                  <div className="mt-6">
                    <button
                      onClick={() => setConnected(true)}
                      className="rounded-lg border border-gray-200 px-4 py-2 text-xs text-gray-400 hover:bg-gray-50"
                    >
                      Simulasi: Device Terhubung
                    </button>
                  </div>
                </div>
              )}
            </div>
          ) : (
            /* Success state */
            <div className="rounded-xl border border-green-200 bg-green-50 p-10 text-center shadow-sm">
              <CheckCircle2 className="mx-auto h-14 w-14 text-green-500" />
              <h2 className="mt-4 text-lg font-bold text-green-800">Device Berhasil Terhubung!</h2>
              <p className="mt-1 text-sm text-gray-600">
                <strong>{deviceName}</strong> sudah aktif dan siap digunakan.
              </p>
              <div className="mt-6 flex flex-col gap-3 sm:flex-row sm:justify-center">
                <a
                  href="/devices/status"
                  className="rounded-lg bg-green-600 px-5 py-2.5 text-sm font-semibold text-white hover:bg-green-700"
                >
                  <Smartphone className="mr-2 inline h-4 w-4" />
                  Lihat Semua Device
                </a>
                <button
                  onClick={() => { setConnected(false); setShowQR(false); setDeviceName(""); }}
                  className="rounded-lg border border-gray-300 px-5 py-2.5 text-sm font-semibold text-gray-700 hover:bg-gray-50"
                >
                  Tambah Device Lain
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
