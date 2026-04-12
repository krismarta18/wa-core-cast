"use client";

import { useState } from "react";
import { User, Save } from "lucide-react";
import { useToast } from "@/components/ui/toast";

export default function ProfilePage() {
  const { success } = useToast();
  const [profile, setProfile] = useState({
    name: "Admin Seewash",
    email: "admin@seewash.id",
    phone: "628111000000",
    timezone: "Asia/Jakarta",
    company: "Seewash Indonesia",
  });

  const [savedProfile, setSavedProfile] = useState(false);

  function saveProfile() {
    success("Profil Disimpan!", "Informasi akun kamu berhasil diperbarui.");
  }

  const TIMEZONES = ["Asia/Jakarta", "Asia/Makassar", "Asia/Jayapura", "UTC"];

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <h1 className="text-xl font-bold text-gray-900">Profil & Akun</h1>
        <p className="text-sm text-gray-500">Kelola informasi akun dan keamanan</p>
      </div>

      <div className="mx-auto max-w-2xl p-6 space-y-5">
        {/* Avatar & name header */}
        <div className="flex items-center gap-4 rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
          <div className="flex h-16 w-16 items-center justify-center rounded-full bg-green-100 text-2xl font-bold text-green-700">
            {profile.name.charAt(0)}
          </div>
          <div>
            <p className="text-lg font-semibold text-gray-900">{profile.name}</p>
            <p className="text-sm text-gray-500">{profile.email}</p>
            <span className="mt-1 inline-block rounded-full bg-green-50 px-2.5 py-0.5 text-xs font-medium text-green-700">
              Business Pro
            </span>
          </div>
        </div>

        {/* Profile form */}
        <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
          <div className="mb-5 flex items-center gap-2">
            <User className="h-5 w-5 text-gray-400" />
            <h2 className="font-semibold text-gray-900">Informasi Akun</h2>
          </div>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
              <label className="mb-1 block text-sm font-medium text-gray-700">Nama Lengkap</label>
              <input
                value={profile.name}
                onChange={(e) => setProfile({ ...profile, name: e.target.value })}
                className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-green-500 focus:outline-none"
              />
            </div>
            <div>
              <label className="mb-1 block text-sm font-medium text-gray-700">Nama Perusahaan</label>
              <input
                value={profile.company}
                onChange={(e) => setProfile({ ...profile, company: e.target.value })}
                className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-green-500 focus:outline-none"
              />
            </div>
            <div>
              <label className="mb-1 block text-sm font-medium text-gray-700">Email</label>
              <input
                type="email"
                value={profile.email}
                onChange={(e) => setProfile({ ...profile, email: e.target.value })}
                className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-green-500 focus:outline-none"
              />
            </div>
            <div>
              <label className="mb-1 block text-sm font-medium text-gray-700">Nomor WA</label>
              <input
                value={profile.phone}
                onChange={(e) => setProfile({ ...profile, phone: e.target.value })}
                className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-green-500 focus:outline-none"
              />
            </div>
            <div className="sm:col-span-2">
              <label className="mb-1 block text-sm font-medium text-gray-700">Zona Waktu</label>
              <select
                value={profile.timezone}
                onChange={(e) => setProfile({ ...profile, timezone: e.target.value })}
                className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-green-500 focus:outline-none"
              >
                {TIMEZONES.map((tz) => <option key={tz}>{tz}</option>)}
              </select>
            </div>
          </div>
          <button
            onClick={saveProfile}
            className="mt-5 flex items-center gap-2 rounded-lg bg-green-600 px-5 py-2 text-sm font-medium text-white hover:bg-green-700"
          >
            <Save className="h-4 w-4" /> Simpan Perubahan
          </button>
        </div>
      </div>
    </div>
  );
}
