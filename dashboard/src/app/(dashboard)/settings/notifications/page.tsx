"use client";

import { useState } from "react";
import { Bell, Mail, Smartphone, Save } from "lucide-react";
import { useToast } from "@/components/ui/toast";

interface NotifSetting {
  id: string;
  label: string;
  description: string;
  email: boolean;
  inApp: boolean;
  emailAllowed: boolean;
}

const INITIAL_SETTINGS: NotifSetting[] = [
  { id: "device_disconnect", label: "Device Terputus", description: "Notifikasi saat koneksi WA device terputus", email: true, inApp: true, emailAllowed: true },
  { id: "quota_80", label: "Kuota 80%", description: "Peringatan saat kuota pesan mencapai 80%", email: true, inApp: true, emailAllowed: true },
  { id: "quota_90", label: "Kuota 90%", description: "Peringatan kritis saat kuota pesan mencapai 90%", email: true, inApp: true, emailAllowed: true },
  { id: "failure_spike", label: "Lonjakan Kegagalan", description: "Notifikasi jika failure rate melebihi 5% dalam 1 jam", email: true, inApp: true, emailAllowed: true },
  { id: "scheduled_sent", label: "Pesan Terjadwal Terkirim", description: "Konfirmasi setelah pesan terjadwal berhasil dikirim", email: false, inApp: true, emailAllowed: false },
  { id: "daily_report", label: "Laporan Harian", description: "Ringkasan statistik pengiriman setiap pukul 08.00", email: true, inApp: false, emailAllowed: true },
  { id: "broadcast_complete", label: "Broadcast Selesai", description: "Notifikasi saat seluruh penerima broadcast telah diproses", email: false, inApp: true, emailAllowed: true },
];

function Toggle({ checked, onChange }: { checked: boolean; onChange: () => void }) {
  return (
    <button
      type="button"
      onClick={onChange}
      className={`relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors focus:outline-none ${checked ? "bg-green-500" : "bg-gray-200"}`}
    >
      <span
        className={`inline-block h-5 w-5 transform rounded-full bg-white shadow-md transition-transform ${checked ? "translate-x-5" : "translate-x-0"}`}
      />
    </button>
  );
}

export default function NotificationsPage() {
  const { success } = useToast();
  const [settings, setSettings] = useState<NotifSetting[]>(INITIAL_SETTINGS);

  function toggle(id: string, field: "email" | "inApp") {
    setSettings(settings.map((s) =>
      s.id === id ? { ...s, [field]: !s[field as keyof NotifSetting] } : s
    ));
  }

  function save() {
    success("Pengaturan Disimpan!", "Preferensi notifikasi kamu berhasil diperbarui.");
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <h1 className="text-xl font-bold text-gray-900">Notifikasi</h1>
        <p className="text-sm text-gray-500">Atur jenis notifikasi yang ingin kamu terima</p>
      </div>

      <div className="mx-auto max-w-2xl p-6 space-y-4">
        {/* Legend */}
        <div className="flex items-center gap-6 rounded-xl border border-gray-200 bg-white px-5 py-3 shadow-sm text-sm text-gray-500">
          <div className="flex items-center gap-2">
            <Mail className="h-4 w-4 text-gray-400" />
            <span>Email</span>
          </div>
          <div className="flex items-center gap-2">
            <Smartphone className="h-4 w-4 text-gray-400" />
            <span>In-App</span>
          </div>
        </div>

        {/* Settings rows */}
        <div className="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm divide-y divide-gray-50">
          {settings.map((s) => (
            <div key={s.id} className="flex items-center gap-4 px-5 py-4">
              <div className="flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-lg bg-gray-100">
                <Bell className="h-4 w-4 text-gray-500" />
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-gray-800">{s.label}</p>
                <p className="text-xs text-gray-400">{s.description}</p>
              </div>
              <div className="flex flex-shrink-0 items-center gap-4">
                <div className="flex flex-col items-center gap-1">
                  <Mail className="h-3 w-3 text-gray-300" />
                  <Toggle
                    checked={s.email && s.emailAllowed}
                    onChange={() => s.emailAllowed && toggle(s.id, "email")}
                  />
                </div>
                <div className="flex flex-col items-center gap-1">
                  <Smartphone className="h-3 w-3 text-gray-300" />
                  <Toggle
                    checked={s.inApp}
                    onChange={() => toggle(s.id, "inApp")}
                  />
                </div>
              </div>
            </div>
          ))}
        </div>

        <button
          onClick={save}
          className="flex w-full items-center justify-center gap-2 rounded-xl bg-green-600 py-3 text-sm font-medium text-white hover:bg-green-700"
        >
          <Save className="h-4 w-4" /> Simpan Pengaturan
        </button>
      </div>
    </div>
  );
}
