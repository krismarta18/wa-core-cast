"use client";

import { useState, useEffect } from "react";
import { Webhook, Save, ToggleLeft, ToggleRight, Info, Loader2 } from "lucide-react";
import { useToast } from "@/components/ui/toast";
import { integrationApi } from "@/lib/api";

const EVENTS = [
  { key: "message.sent", label: "Pesan Terkirim", desc: "Trigger saat pesan berhasil dikirim" },
  { key: "message.failed", label: "Pesan Gagal", desc: "Trigger saat pesan gagal dikirim" },
  { key: "message.received", label: "Pesan Diterima", desc: "Trigger saat pesan masuk diterima" },
  { key: "device.connected", label: "Device Terhubung", desc: "Trigger saat device reconnect" },
  { key: "device.disconnected", label: "Device Terputus", desc: "Trigger saat device disconnect" },
];

export default function WebhookSettingsPage() {
  const { success, error } = useToast();
  const [loading, setLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  
  const [url, setUrl] = useState("");
  const [secret, setSecret] = useState("");
  const [isActive, setIsActive] = useState(true);
  const [enabledEvents, setEnabledEvents] = useState<string[]>([]);

  useEffect(() => {
    fetchSettings();
  }, []);

  const fetchSettings = async () => {
    try {
      setLoading(true);
      const res = await integrationApi.getWebhook();
      setUrl(res.url || "");
      setSecret(res.secret || "");
      setIsActive(res.is_active);
      setEnabledEvents(res.enabled_events || []);
    } catch (err: any) {
      error("Gagal", "Gagal mengambil pengaturan webhook");
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    try {
      setIsSaving(true);
      await integrationApi.updateWebhook({
        url,
        secret,
        is_active: isActive,
        enabled_events: enabledEvents
      });
      success("Tersimpan", "Pengaturan webhook berhasil diperbarui.");
    } catch (err: any) {
      error("Gagal", "Gagal menyimpan pengaturan webhook");
    } finally {
      setIsSaving(false);
    }
  };

  const toggleEvent = (key: string) => {
    if (enabledEvents.includes(key)) {
      setEnabledEvents(enabledEvents.filter(e => e !== key));
    } else {
      setEnabledEvents([...enabledEvents, key]);
    }
  };

  if (loading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-green-600" />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold text-gray-900">Webhook Settings</h1>
            <p className="text-sm text-gray-500">Konfigurasi endpoint yang menerima notifikasi event real-time</p>
          </div>
          <div className="flex items-center gap-2">
            <span className="text-xs font-semibold text-gray-400">{isActive ? "AKTIF" : "NONAKTIF"}</span>
            <button onClick={() => setIsActive(!isActive)}>
              {isActive ? <ToggleRight className="h-7 w-7 text-green-600" /> : <ToggleLeft className="h-7 w-7 text-gray-400" />}
            </button>
          </div>
        </div>
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
                  Webhook Secret (HMAC Signature)
                </label>
                <input
                  type="text"
                  value={secret}
                  onChange={(e) => setSecret(e.target.value)}
                  placeholder="whsec_..."
                  className="h-10 w-full rounded-lg border border-gray-300 px-3 font-mono text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
                />
              </div>
            </div>
          </div>

          {/* Events */}
          <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
            <h2 className="mb-4 font-semibold text-gray-900">Event yang Ingin Diterima</h2>
            <div className="space-y-3">
              {EVENTS.map((ev) => (
                <div key={ev.key} className="flex items-center justify-between gap-4 rounded-lg bg-gray-50 px-4 py-3">
                  <div>
                    <p className="text-sm font-medium text-gray-900">{ev.label}</p>
                    <p className="text-xs text-gray-500">{ev.desc}</p>
                  </div>
                  <button
                    onClick={() => toggleEvent(ev.key)}
                    className="flex-shrink-0"
                  >
                    {enabledEvents.includes(ev.key) ? (
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
            <span>
              WACAST akan mengirimkan POST request ke URL Anda. Gunakan secret di atas untuk memverifikasi header 
              <code className="mx-1 rounded bg-yellow-100 px-1 font-bold">X-WACAST-Signature</code>.
            </span>
          </div>

          <button
            onClick={handleSave}
            disabled={isSaving}
            className="flex w-full items-center justify-center gap-2 rounded-xl bg-green-600 py-3 text-sm font-semibold text-white shadow-sm hover:bg-green-700 transition-colors disabled:opacity-50"
          >
            {isSaving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
            Simpan Pengaturan Webhook
          </button>
        </div>
      </div>
    </div>
  );
}
