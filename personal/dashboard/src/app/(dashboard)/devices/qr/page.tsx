"use client";

import { useState, useEffect, useRef } from "react";
import { QrCode, Smartphone, CheckCircle2, RefreshCw, Info, Loader2 } from "lucide-react";
import { authApi, sessionsApi } from "@/lib/api";
import { useToast } from "@/components/ui/toast";

const STEPS = [
  "Buka WhatsApp di HP kamu",
  "Ketuk ikon titik tiga (⋮) di pojok kanan atas",
  "Pilih \"Perangkat Tertaut\" lalu \"Tautkan Perangkat\"",
  "Arahkan kamera ke QR code di bawah ini",
];

export default function QRScannerPage() {
  const { error } = useToast();
  const [deviceName, setDeviceName] = useState("");
  const [showQR, setShowQR] = useState(false);
  const [connected, setConnected] = useState(false);
  
  // Real-time integration states
  const [qrSrc, setQrSrc] = useState<string>("");
  const [loading, setLoading] = useState(false);
  const [deviceId, setDeviceId] = useState<string>("");
  const pollingRef = useRef<NodeJS.Timeout | null>(null);

  // Helper to safely slugify device name -> now using UUID
  function generateSafeId() {
    return crypto.randomUUID();
  }

  // Cleanup polling on unmount
  useEffect(() => {
    return () => {
      if (pollingRef.current) clearInterval(pollingRef.current);
    };
  }, []);

  async function handleGenerateQR() {
    setLoading(true);
    try {
      // 1. Get current logged in user context
      const meRes = await authApi.me();
      if (!meRes.success || !meRes.user) {
        throw new Error("Sesi pengguna tidak valid.");
      }

      const newDeviceId = generateSafeId();
      setDeviceId(newDeviceId);

      // 2. Initiate session
      await sessionsApi.initiate({
        device_id: newDeviceId,
        user_id: meRes.user.id,
        phone: meRes.user.phone_number,
        display_name: deviceName,
      });

      setShowQR(true);
      startPollingQR(newDeviceId);

    } catch (err: any) {
      console.error(err);
      if (err.response?.status === 403) {
        error("Batas Device Tercapai", "Paket Anda telah mencapai batas maksimum perangkat. Silakan upgrade paket untuk menambah perangkat baru.");
      } else {
        error("Gagal", err.response?.data?.error || "Gagal membuat sesi.");
      }
    } finally {
      setLoading(false);
    }
  }

  function startPollingQR(currentDeviceId: string) {
    // Clear any existing polling
    if (pollingRef.current) clearInterval(pollingRef.current);

    // Initial fetch
    fetchQrAndStatus(currentDeviceId);

    // Pool every 4 seconds
    pollingRef.current = setInterval(() => {
      fetchQrAndStatus(currentDeviceId);
    }, 4000);
  }

  async function fetchQrAndStatus(currentDeviceId: string) {
    try {
      // Check status
      const statusRes = await sessionsApi.get(currentDeviceId).catch(() => null);
      if (statusRes && statusRes.status === 1) { // 1 = session active
        if (pollingRef.current) clearInterval(pollingRef.current);
        setConnected(true);
        setShowQR(false);
        return;
      }

      // Fetch QR
      const qrRes = await sessionsApi.qr(currentDeviceId).catch(() => null);
      if (qrRes && qrRes.qr_code_image?.base64_png) {
        setQrSrc(qrRes.qr_code_image.base64_png);
      }
    } catch (err) {
      console.error("Kesalahan sinkronisasi dengan server", err);
    }
  }

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
                  disabled={!deviceName.trim() || loading}
                  onClick={handleGenerateQR}
                  className="mt-4 w-full flex items-center justify-center rounded-lg bg-green-600 py-2.5 text-sm font-semibold text-white hover:bg-green-700 disabled:cursor-not-allowed disabled:bg-gray-200 disabled:text-gray-400 transition-colors"
                >
                  {loading ? (
                    <>
                      <Loader2 className="mr-2 inline h-4 w-4 animate-spin" />
                      Membuat Sesi...
                    </>
                  ) : (
                    <>
                      <QrCode className="mr-2 inline h-4 w-4" />
                      Generate QR Code
                    </>
                  )}
                </button>
              </div>

              {/* QR Display */}
              {showQR && (
                <div className="rounded-xl border border-gray-200 bg-white p-8 text-center shadow-sm">
                  <p className="mb-1 text-sm font-medium text-gray-700">
                    Scan QR untuk <span className="font-bold text-green-700">{deviceName}</span>
                  </p>
                  <p className="mb-6 text-xs text-gray-400">
                    ID Perangkat: <span className="font-mono text-gray-600 bg-gray-100 px-1 rounded">{deviceId}</span>
                  </p>

                  <div className="mx-auto flex h-60 w-60 items-center justify-center rounded-xl border-2 border-dashed border-gray-300 bg-gray-50 overflow-hidden">
                    {qrSrc ? (
                      <img src={qrSrc} alt="WhatsApp Web QR Code" className="h-full w-full object-contain p-2" />
                    ) : (
                      <div className="text-center text-gray-400 flex flex-col items-center">
                        <Loader2 className="mb-2 h-8 w-8 animate-spin text-green-500" />
                        <span className="text-xs">Mengambil QR terbaru...</span>
                      </div>
                    )}
                  </div>

                  <button onClick={() => fetchQrAndStatus(deviceId)} className="mt-5 inline-flex items-center gap-2 text-sm font-medium text-green-600 hover:underline">
                    <RefreshCw className="h-4 w-4" />Refresh Sekarang
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
                  onClick={() => { setConnected(false); setShowQR(false); setDeviceName(""); setQrSrc(""); setDeviceId(""); }}
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
