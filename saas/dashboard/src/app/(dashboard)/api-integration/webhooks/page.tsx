"use client";

import { useState } from "react";
import { Webhook, Save, ToggleLeft, ToggleRight, Info } from "lucide-react";
import { useToast } from "@/components/ui/toast";

const EVENTS = [
  { key: "message.sent", label: "Pesan Terkirim", desc: "Trigger saat pesan berhasil dikirim" },
  { key: "message.failed", label: "Pesan Gagal", desc: "Trigger saat pesan gagal dikirim" },
  { key: "message.received", label: "Pesan Diterima", desc: "Trigger saat pesan masuk diterima" },
  { key: "device.connected", label: "Device Terhubung", desc: "Trigger saat device reconnect" },
  { key: "device.disconnected", label: "Device Terputus", desc: "Trigger saat device disconnect" },
  { key: "qr.generated", label: "QR Code Generated", desc: "Trigger saat QR code baru dibuat" },
];

export default function WebhookSettingsPage() {
  const { success } = useToast();
  const [url, setUrl] = useState("https://myapp.com/webhook/wacast");
  const [secret, setSecret] = useState("whsec_••••••••••••••••");
  const [enabled, setEnabled] = useState<Record<string, boolean>>(
    Object.fromEntries(EVENTS.map((e) => [e.key, true]))
  );

  const handleSave = () => {
    success("Webhook Disimpan!", "Konfigurasi endpoint berhasil diperbarui.");
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <h1 className="text-xl font-bold text-gray-900">Webhook Settings</h1>
        <p className="text-sm text-gray-500">Konfigurasi endpoint yang menerima notifikasi event</p>
      </div>

      <div className="p-6">
        <div className="mx-auto max-w-2xl space-y-5">

          {/* Endpoint */}
          <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
            <h2 className="mb-4 font-semibold text-gray-900">Endpoint URL</h2>

            <div className="space-y-4">
              <div>
                <label className="mb-1 block text-sm font-medium text-gray-700">
                  URL Webhook <span className="text-red-500">*</span>
                </label>
                <div className="relative">
                  <Webhook className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
                  <input
                    type="url"
                    value={url}
                    onChange={(e) => setUrl(e.target.value)}
                    placeholder="https://yourapp.com/webhook"
                    className="h-10 w-full rounded-lg border border-gray-300 pl-9 pr-3 text-sm placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-green-500"
                  />
                </div>
              </div>

              <div>
                <label className="mb-1 block text-sm font-medium text-gray-700">
                  Secret Key (untuk verifikasi signature)
                </label>
                <input
                  type="text"
                  value={secret}
                  onChange={(e) => setSecret(e.target.value)}
                  className="h-10 w-full rounded-lg border border-gray-300 px-3 font-mono text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
                />
              </div>
            </div>
          </div>

          {/* Events */}
          <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
            <h2 className="mb-4 font-semibold text-gray-900">Event yang Dikirim</h2>
            <div className="space-y-3">
              {EVENTS.map((ev) => (
                <div key={ev.key} className="flex items-center justify-between gap-4 rounded-lg bg-gray-50 px-4 py-3">
                  <div>
                    <p className="text-sm font-medium text-gray-900">{ev.label}</p>
                    <p className="text-xs text-gray-500">{ev.desc}</p>
                  </div>
                  <button
                    onClick={() => setEnabled((e) => ({ ...e, [ev.key]: !e[ev.key] }))}
                    className="flex-shrink-0"
                  >
                    {enabled[ev.key] ? (
                      <ToggleRight className="h-7 w-7 text-green-600" />
                    ) : (
                      <ToggleLeft className="h-7 w-7 text-gray-400" />
                    )}
                  </button>
                </div>
              ))}
            </div>
          </div>

          {/* Info */}
          <div className="flex items-start gap-2 rounded-xl border border-yellow-100 bg-yellow-50 p-4 text-sm text-yellow-800">
            <Info className="mt-0.5 h-4 w-4 flex-shrink-0" />
            Webhook akan dikirim dengan method POST. Header <code className="rounded bg-yellow-100 px-1">X-WACAST-Signature</code> berisi HMAC-SHA256 dari body menggunakan secret key.
          </div>

          <button
            onClick={handleSave}
            className="w-full rounded-xl bg-green-600 py-3 text-sm font-semibold text-white shadow-sm hover:bg-green-700 transition-colors"
          >
            <Save className="mr-2 inline h-4 w-4" />Simpan Pengaturan
          </button>
        </div>
      </div>
    </div>
  );
}
