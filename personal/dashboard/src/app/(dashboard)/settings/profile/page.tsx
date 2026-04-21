"use client";

import { useState, useEffect } from "react";
import { User, Save, Loader2, ShieldCheck, Zap, AlertTriangle, Info } from "lucide-react";
import { useToast } from "@/components/ui/toast";
import { authApi, settingsApi } from "@/lib/api";
import type { SystemSetting } from "@/lib/types";

export default function ProfilePage() {
  const { success, error, info } = useToast();
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  
  // Profile state
  const [profile, setProfile] = useState({
    name: "",
    email: "",
    phone: "",
    timezone: "Asia/Jakarta",
  });

  // Anti-Bot settings state
  const [antiBot, setAntiBot] = useState({
    enabled: true,
    suffixLength: 4,
  });
  const [settingsLoading, setSettingsLoading] = useState(true);
  const [settingsSaving, setSettingsSaving] = useState(false);

  // Password state
  const [passwords, setPasswords] = useState({
    old: "",
    new: "",
    confirm: "",
  });
  const [passwordLoading, setPasswordLoading] = useState(false);

  useEffect(() => {
    // Fetch Profile
    authApi.me()
      .then((res) => {
        if (res.success && res.user) {
          setProfile({
            name: res.user.full_name || "",
            email: res.user.email || "",
            phone: res.user.phone_number || "",
            timezone: res.user.timezone || "Asia/Jakarta",
          });
        }
      })
      .catch(console.error)
      .finally(() => setLoading(false));

    // Fetch Anti-Bot Settings
    settingsApi.get()
      .then((res) => {
        const enabled = res.settings.find(s => s.key === "anti_bot_enabled");
        const length = res.settings.find(s => s.key === "anti_bot_suffix_length");
        
        setAntiBot({
          enabled: enabled ? enabled.value === "true" : true,
          suffixLength: length ? parseInt(length.value) || 4 : 4,
        });
      })
      .catch(console.error)
      .finally(() => setSettingsLoading(false));
  }, []);

  async function saveProfile() {
    setSaving(true);
    try {
      await authApi.updateProfile({
        full_name: profile.name,
        email: profile.email,
        timezone: profile.timezone,
      });
      success("Profil Disimpan", "Informasi akun kamu berhasil diperbarui.");
    } catch(err) {
      error("Gagal", "Gagal memperbarui profil.");
    } finally {
      setSaving(false);
    }
  }

  async function saveAntiBot() {
    setSettingsSaving(true);
    try {
      await settingsApi.update("anti_bot_enabled", antiBot.enabled ? "true" : "false");
      await settingsApi.update("anti_bot_suffix_length", antiBot.suffixLength.toString());
      success("Pengaturan Disimpan", "Fitur Anti-Bot Personal Pro berhasil dikonfigurasi.");
    } catch(err) {
      error("Gagal", "Gagal memperbarui pengaturan Anti-Bot.");
    } finally {
      setSettingsSaving(false);
    }
  }

  async function updatePassword() {
    if (!passwords.old || !passwords.new) {
      error("Gagal", "Password lama dan baru harus diisi.");
      return;
    }
    if (passwords.new !== passwords.confirm) {
      error("Gagal", "Konfirmasi password tidak cocok.");
      return;
    }
    
    setPasswordLoading(true);
    try {
      await authApi.changePassword({
        old_password: passwords.old,
        new_password: passwords.new,
      });
      success("Password Berhasil!", "Password admin telah diperbarui.");
      setPasswords({ old: "", new: "", confirm: "" });
    } catch (err: any) {
      error("Gagal", err.response?.data?.error?.message || "Gagal ganti password.");
    } finally {
      setPasswordLoading(false);
    }
  }

  const TIMEZONES = ["Asia/Jakarta", "Asia/Makassar", "Asia/Jayapura", "UTC"];

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <h1 className="text-xl font-bold text-gray-900">Profil & Akun</h1>
        <p className="text-sm text-gray-500">Kelola informasi diri dan fitur Personal Pro</p>
      </div>

      <div className="mx-auto max-w-3xl p-6 grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="space-y-6">
           {/* Header Card */}
           <div className="flex items-center gap-4 rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
            <div className="flex h-14 w-14 items-center justify-center rounded-full bg-indigo-100 text-xl font-bold text-indigo-700">
              {profile.name.charAt(0) || "U"}
            </div>
            <div>
               <p className="text-lg font-bold text-gray-900 leading-tight">{profile.name}</p>
               <span className="mt-1 inline-flex items-center gap-1.5 rounded-full bg-green-50 px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider text-green-700 border border-green-200">
                  <ShieldCheck className="h-3 w-3" /> Personal Pro Edition
               </span>
            </div>
          </div>

          {/* Profile Form */}
          <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
            <h2 className="mb-4 flex items-center gap-2 font-semibold text-gray-900">
              <User className="h-4 w-4 text-gray-400" /> Informasi Dasar
            </h2>
            <div className="space-y-3">
              <div>
                <label className="mb-1 block text-xs font-medium text-gray-500">Nama Lengkap</label>
                <input
                  value={profile.name}
                  onChange={(e) => setProfile({ ...profile, name: e.target.value })}
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none"
                />
              </div>
              <div>
                <label className="mb-1 block text-xs font-medium text-gray-500">Email</label>
                <input
                  type="email"
                  value={profile.email}
                  onChange={(e) => setProfile({ ...profile, email: e.target.value })}
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none"
                />
              </div>
              <div>
                <label className="mb-1 block text-xs font-medium text-gray-500">Zona Waktu</label>
                <select
                  value={profile.timezone}
                  onChange={(e) => setProfile({ ...profile, timezone: e.target.value })}
                  className="w-full rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none"
                >
                  {TIMEZONES.map((tz) => <option key={tz} value={tz}>{tz}</option>)}
                </select>
              </div>
            </div>
            <button
              onClick={saveProfile}
              disabled={saving || loading}
              className="mt-5 w-full flex items-center justify-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-semibold text-white hover:bg-indigo-700 disabled:opacity-50"
            >
              {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
              Update Profil
            </button>
          </div>

          {/* Password Form */}
          <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
            <h2 className="mb-4 flex items-center gap-2 font-semibold text-gray-900">
              <ShieldCheck className="h-4 w-4 text-gray-400" /> Keamanan Admin
            </h2>
            <div className="space-y-3">
              <input
                type="password"
                placeholder="Password Lama"
                value={passwords.old}
                onChange={(e) => setPasswords({ ...passwords, old: e.target.value })}
                className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none"
              />
              <input
                type="password"
                placeholder="Password Baru"
                value={passwords.new}
                onChange={(e) => setPasswords({ ...passwords, new: e.target.value })}
                className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none"
              />
              <input
                type="password"
                placeholder="Konfirmasi Password Baru"
                value={passwords.confirm}
                onChange={(e) => setPasswords({ ...passwords, confirm: e.target.value })}
                className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none"
              />
            </div>
            <button
              onClick={updatePassword}
              disabled={passwordLoading}
              className="mt-5 w-full flex items-center justify-center gap-2 rounded-lg bg-gray-900 px-4 py-2 text-sm font-semibold text-white hover:bg-black disabled:opacity-50"
            >
              {passwordLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
              Ganti Password
            </button>
          </div>
        </div>

        <div className="space-y-6">
          {/* Anti-Bot Pro Settings Card */}
          <div className="rounded-xl border border-indigo-200 bg-white p-6 shadow-md overflow-hidden relative">
            <div className="absolute top-0 right-0 p-1">
               <Zap className="h-20 w-20 text-indigo-500 opacity-5" />
            </div>
            <h2 className="mb-4 flex items-center gap-2 font-bold text-indigo-900">
              <ShieldCheck className="h-5 w-5 text-indigo-600" /> Fitur Anti-Bot (Pro)
            </h2>
            
            <div className="space-y-5">
              <div className="flex items-center justify-between gap-4 p-4 rounded-lg bg-indigo-50 border border-indigo-100">
                <div>
                  <p className="text-sm font-bold text-indigo-900">Detect Anti-Bot</p>
                  <p className="text-[10px] text-indigo-600">Sisipkan karakter acak untuk meniru manusia.</p>
                </div>
                <label className="relative inline-flex cursor-pointer items-center">
                  <input
                    type="checkbox"
                    checked={antiBot.enabled}
                    onChange={(e) => setAntiBot({ ...antiBot, enabled: e.target.checked })}
                    className="peer sr-only"
                  />
                  <div className="h-6 w-11 rounded-full bg-gray-200 after:absolute after:left-[2px] after:top-[2px] after:h-5 after:w-5 after:rounded-full after:border after:border-gray-300 after:bg-white after:transition-all after:content-[''] peer-checked:bg-indigo-600 peer-checked:after:translate-x-full peer-checked:after:border-white peer-focus:outline-none"></div>
                </label>
              </div>

              <div>
                <label className="mb-2 block text-sm font-bold text-gray-700">Panjang Suffix Acak</label>
                <input
                  type="range"
                  min="2"
                  max="12"
                  value={antiBot.suffixLength}
                  onChange={(e) => setAntiBot({ ...antiBot, suffixLength: parseInt(e.target.value) })}
                  className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-indigo-600"
                />
                <div className="mt-1 flex justify-between text-[10px] text-gray-400 font-bold">
                   <span>Minimal (2)</span>
                   <span className="text-indigo-600 text-sm">{antiBot.suffixLength} Karakter</span>
                   <span>Maksimal (12)</span>
                </div>
              </div>

              <div className="rounded-lg border border-amber-200 bg-amber-50 p-3 flex gap-2">
                 <AlertTriangle className="h-4 w-4 text-amber-500 shrink-0 mt-0.5" />
                 <p className="text-[10px] text-amber-700">
                    Semakin panjang suffix, semakin sulit dideteksi sistem bot WhatsApp, namun pesan akan terlihat memiliki kode unik di bawahnya.
                 </p>
              </div>

              <button
                onClick={saveAntiBot}
                disabled={settingsSaving || settingsLoading}
                className="w-full flex items-center justify-center gap-2 rounded-lg bg-indigo-600 py-3 text-sm font-bold text-white shadow-lg hover:bg-indigo-700 disabled:opacity-50 transition-all hover:scale-[1.02] active:scale-95"
              >
                {settingsSaving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
                Aktifkan Konfigurasi Pro
              </button>
            </div>
          </div>

          <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
             <div className="flex items-center gap-3 mb-2">
                <Info className="h-5 w-5 text-gray-400" />
                <h3 className="font-semibold text-gray-900">Edisi Personal</h3>
             </div>
             <p className="text-xs text-gray-500 leading-relaxed">
                Anda menggunakan <strong>WACAST Personal Edition</strong>. 
                Fitur terbatas pada satu instance PC dan tidak memerlukan langganan bulanan. 
                Seluruh data disimpan secara lokal pada komputer Anda.
             </p>
          </div>
        </div>
      </div>
    </div>
  );
}
