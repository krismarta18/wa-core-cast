"use client";

import { useState } from "react";
import { KeyRound, Plus, Copy, Eye, EyeOff, Trash2, Shield } from "lucide-react";
import { useToast } from "@/components/ui/toast";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";

const KEYS = [
  { id: 1, name: "Production Key", key: "wck_live_a8f3d2e1b9c4f7a2d5e8b1c4", created: "01 Mar 2026", lastUsed: "Baru saja", active: true },
  { id: 2, name: "Staging Key", key: "wck_test_b3e7f1a4d8c2e5b9a3f6d1e4", created: "15 Feb 2026", lastUsed: "3 hari lalu", active: true },
  { id: 3, name: "Old Key (deprecated)", key: "wck_live_c1d4e8b2f5a7d3e9b6c2f1a5", created: "10 Jan 2026", lastUsed: "30 hari lalu", active: false },
];

function maskKey(key: string) {
  return key.slice(0, 12) + "••••••••••••••••" + key.slice(-4);
}

export default function APIKeysPage() {
  const { success, info } = useToast();
  const [visible, setVisible] = useState<Record<number, boolean>>({});
  const [copied, setCopied] = useState<number | null>(null);
  const [keys, setKeys] = useState(KEYS);
  const [confirmDelete, setConfirmDelete] = useState<number | null>(null);

  const toggleVisible = (id: number) =>
    setVisible((v) => ({ ...v, [id]: !v[id] }));

  const copyKey = (id: number, key: string) => {
    navigator.clipboard.writeText(key);
    setCopied(id);
    setTimeout(() => setCopied(null), 2000);
    info("Disalin!", "API key berhasil disalin ke clipboard.");
  };

  const createKey = () => {
    const newKey = { id: Date.now(), name: `Key ${keys.length + 1}`, key: `wck_live_${Math.random().toString(36).slice(2, 26)}`, created: "12 Apr 2026", lastUsed: "Belum pernah", active: true };
    setKeys([...keys, newKey]);
    success("API Key Dibuat!", "Key baru berhasil dibuat dan siap digunakan.");
  };

  const deleteKey = (id: number) => {
    setKeys(keys.filter((k) => k.id !== id));
    setConfirmDelete(null);
    info("Key Dihapus", "API key berhasil dihapus.");
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <ConfirmDialog
        open={confirmDelete !== null}
        title="Hapus API Key?"
        description="Key ini akan dihapus permanen. Integrasi yang menggunakan key ini akan berhenti berfungsi."
        confirmLabel="Ya, Hapus"
        onConfirm={() => confirmDelete !== null && deleteKey(confirmDelete)}
        onCancel={() => setConfirmDelete(null)}
      />
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold text-gray-900">API Keys</h1>
            <p className="text-sm text-gray-500">Kelola kunci API untuk integrasi eksternal</p>
          </div>
          <button onClick={createKey} className="inline-flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-sm font-semibold text-white hover:bg-green-700 shadow-sm">
            <Plus className="h-4 w-4" /> Buat Key Baru
          </button>
        </div>
      </div>

      <div className="p-6 space-y-4">
        {/* Info box */}
        <div className="rounded-xl border border-blue-100 bg-blue-50 p-4 text-sm text-blue-800">
          <div className="flex items-start gap-2">
            <Shield className="mt-0.5 h-4 w-4 flex-shrink-0" />
            <span>Simpan API key dengan aman. Jangan bagikan key kepada siapapun. Gunakan environment variable di aplikasi Anda.</span>
          </div>
        </div>

        {/* Keys list */}
        {keys.map((k) => (
          <div key={k.id} className={`rounded-xl border bg-white p-5 shadow-sm transition-opacity ${k.active ? "border-gray-200" : "opacity-60 border-gray-100"}`}>
            <div className="flex flex-wrap items-start justify-between gap-3">
              <div className="flex items-center gap-3">
                <div className={`flex h-9 w-9 items-center justify-center rounded-lg ${k.active ? "bg-green-50" : "bg-gray-100"}`}>
                  <KeyRound className={`h-5 w-5 ${k.active ? "text-green-600" : "text-gray-400"}`} />
                </div>
                <div>
                  <div className="flex items-center gap-2">
                    <p className="font-semibold text-gray-900">{k.name}</p>
                    <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${k.active ? "bg-green-50 text-green-700" : "bg-gray-100 text-gray-500"}`}>
                      {k.active ? "Aktif" : "Nonaktif"}
                    </span>
                  </div>
                  <p className="mt-0.5 text-xs text-gray-500">Dibuat {k.created} · Terakhir digunakan: {k.lastUsed}</p>
                </div>
              </div>
              <button onClick={() => setConfirmDelete(k.id)} className="rounded-lg p-1.5 text-gray-400 hover:bg-red-50 hover:text-red-500">
                <Trash2 className="h-4 w-4" />
              </button>
            </div>

            {/* Key display */}
            <div className="mt-3 flex items-center gap-2">
              <code className="flex-1 overflow-hidden rounded-lg border border-gray-200 bg-gray-50 px-3 py-2 font-mono text-xs text-gray-700 break-all">
                {visible[k.id] ? k.key : maskKey(k.key)}
              </code>
              <button
                onClick={() => toggleVisible(k.id)}
                className="rounded-lg border border-gray-200 p-2 text-gray-400 hover:bg-gray-50 hover:text-gray-600"
              >
                {visible[k.id] ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
              </button>
              <button
                onClick={() => copyKey(k.id, k.key)}
                className="rounded-lg border border-gray-200 p-2 text-gray-400 hover:bg-gray-50 hover:text-gray-600"
              >
                {copied === k.id ? (
                  <span className="text-xs font-medium text-green-600 px-1">✓</span>
                ) : (
                  <Copy className="h-4 w-4" />
                )}
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
