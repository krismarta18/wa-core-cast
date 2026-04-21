"use client";

import { useState, useEffect } from "react";
import { KeyRound, Plus, Copy, Eye, EyeOff, Trash2, Shield, Loader2 } from "lucide-react";
import { useToast } from "@/components/ui/toast";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";
import { integrationApi } from "@/lib/api";
import type { APIKey } from "@/lib/types";

function maskKey(prefix: string) {
  return prefix + "••••••••••••••••";
}

export default function APIKeysPage() {
  const { success, error, info } = useToast();
  const [loading, setLoading] = useState(true);
  const [keys, setKeys] = useState<APIKey[]>([]);
  
  const [visible, setVisible] = useState<Record<string, boolean>>({});
  const [copied, setCopied] = useState<string | null>(null);
  const [confirmDelete, setConfirmDelete] = useState<string | null>(null);
  const [isCreating, setIsCreating] = useState(false);
  const [newKeyLabel, setNewKeyLabel] = useState("");
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [createdKey, setCreatedKey] = useState<string | null>(null);

  useEffect(() => {
    fetchKeys();
  }, []);

  const fetchKeys = async () => {
    try {
      setLoading(true);
      const res = await integrationApi.getKeys();
      setKeys(res.keys || []);
    } catch (err: any) {
      const msg = err.response?.data?.details || err.response?.data?.error || "Gagal mengambil daftar API Key";
      error("Gagal", msg);
    } finally {
      setLoading(false);
    }
  };

  const toggleVisible = (id: string) =>
    setVisible((v) => ({ ...v, [id]: !v[id] }));

  const copyToClipboard = (id: string, text: string) => {
    navigator.clipboard.writeText(text);
    setCopied(id);
    setTimeout(() => setCopied(null), 2000);
    info("Disalin!", "API key berhasil disalin ke clipboard.");
  };

  const handleCreate = async () => {
    if (!newKeyLabel) return;
    try {
      setIsCreating(true);
      const res = await integrationApi.createKey(newKeyLabel);
      setCreatedKey(res.key);
      setNewKeyLabel("");
      fetchKeys();
    } catch (err: any) {
      const msg = err.response?.data?.details || err.response?.data?.error || "Gagal membuat API Key";
      error("Gagal", msg);
    } finally {
      setIsCreating(false);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await integrationApi.deleteKey(id);
      setKeys(keys.filter((k) => k.id !== id));
      setConfirmDelete(null);
      success("Dihapus", "API key berhasil dihapus.");
    } catch (err: any) {
      error("Gagal", "Gagal menghapus API Key");
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <ConfirmDialog
        open={confirmDelete !== null}
        title="Hapus API Key?"
        description="Integrasi yang menggunakan key ini akan langsung berhenti berfungsi."
        confirmLabel="Ya, Hapus"
        onConfirm={() => confirmDelete !== null && handleDelete(confirmDelete)}
        onCancel={() => setConfirmDelete(null)}
      />

      {/* New Key Revealed Modal */}
      {createdKey !== null && (
        <div className="fixed inset-0 z-[210] flex items-center justify-center p-4">
          <div className="absolute inset-0 bg-black/60 backdrop-blur-sm" onClick={() => setCreatedKey(null)} />
          <div className="relative w-full max-w-md rounded-2xl border border-gray-200 bg-white p-6 shadow-2xl animate-in zoom-in-95 duration-200">
            <h3 className="text-lg font-bold text-gray-900">API Key Berhasil Dibuat</h3>
            <p className="mt-2 text-sm text-gray-600">
              Salin key ini <span className="font-bold text-red-600 underline">SEKARANG</span>. 
              Anda tidak akan bisa melihatnya lagi setelah menutup jendela ini demi keamanan.
            </p>
            
            <div className="mt-4 rounded-xl bg-green-50 p-4 border border-green-100">
               <div className="flex items-center gap-3">
                  <code className="text-sm font-mono font-bold text-green-800 break-all flex-1 select-all">{createdKey}</code>
                  <button 
                    onClick={() => createdKey && copyToClipboard("new", createdKey)} 
                    className="flex-shrink-0 p-2 hover:bg-green-100 rounded-lg text-green-600 transition-colors"
                    title="Copy to clipboard"
                  >
                    <Copy className="h-5 w-5" />
                  </button>
               </div>
            </div>

            <div className="mt-6">
              <button 
                onClick={() => setCreatedKey(null)}
                className="w-full rounded-xl bg-gray-900 py-3 text-sm font-bold text-white hover:bg-gray-800 transition-colors"
              >
                Saya Sudah Simpan Key Ini
              </button>
            </div>
          </div>
        </div>
      )}

      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold text-gray-900">API Keys</h1>
            <p className="text-sm text-gray-500">Kelola kunci API untuk integrasi sistem eksternal</p>
          </div>
          <button 
            onClick={() => setShowCreateModal(true)} 
            className="inline-flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-sm font-semibold text-white hover:bg-green-700 shadow-sm"
          >
            <Plus className="h-4 w-4" /> Buat Key Baru
          </button>
        </div>
      </div>

      <div className="p-6 space-y-4">
        {showCreateModal && !createdKey && (
           <div className="rounded-xl border border-green-200 bg-white p-5 shadow-sm animate-in fade-in zoom-in-95">
              <h3 className="font-semibold text-gray-900 mb-3">Buat API Key Baru</h3>
              <div className="flex gap-2">
                <input 
                  autoFocus
                  type="text" 
                  placeholder="Nama Key (misal: CRM Integration)" 
                  value={newKeyLabel}
                  onChange={e => setNewKeyLabel(e.target.value)}
                  className="h-10 flex-1 rounded-lg border border-gray-300 px-3 text-sm focus:ring-2 focus:ring-green-500 outline-none"
                />
                <button 
                  onClick={handleCreate}
                  disabled={!newKeyLabel || isCreating}
                  className="rounded-lg bg-green-600 px-4 py-2 text-sm font-bold text-white hover:bg-green-700 disabled:opacity-50 flex items-center gap-2"
                >
                  {isCreating && <Loader2 className="h-3 w-3 animate-spin"/>}
                  Generate
                </button>
                <button onClick={() => setShowCreateModal(false)} className="px-4 text-sm text-gray-500 font-medium">Batal</button>
              </div>
           </div>
        )}

        <div className="rounded-xl border border-blue-100 bg-blue-50 p-4 text-sm text-blue-800">
          <div className="flex items-start gap-2">
            <Shield className="mt-0.5 h-4 w-4 flex-shrink-0" />
            <span>Gunakan API Key ini di Header HTTP: <code className="bg-blue-100 px-1 rounded">X-API-Key: [KEY_ANDA]</code> untuk mengakses API pengiriman pesan secara eksternal.</span>
          </div>
        </div>

        {loading ? (
          <div className="flex justify-center py-12"><Loader2 className="h-8 w-8 animate-spin text-green-600"/></div>
        ) : keys.length === 0 ? (
          <div className="rounded-xl border border-dashed border-gray-300 p-12 text-center">
            <p className="text-gray-500">Belum ada API Key. Silakan buat key baru untuk mulai berintegrasi.</p>
          </div>
        ) : (
          keys.map((k) => (
            <div key={k.id} className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm transition-all hover:shadow-md">
              <div className="flex flex-wrap items-start justify-between gap-3">
                <div className="flex items-center gap-3">
                  <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-green-50 text-green-600">
                    <KeyRound className="h-5 w-5" />
                  </div>
                  <div>
                    <p className="font-semibold text-gray-900">{k.name}</p>
                    <p className="mt-0.5 text-xs text-gray-500">
                      Dibuat {new Date(k.created_at).toLocaleDateString()} · 
                      Terakhir dipakai: {k.last_used_at ? new Date(k.last_used_at).toLocaleString() : "Belum pernah"}
                    </p>
                  </div>
                </div>
                <button onClick={() => setConfirmDelete(k.id)} className="rounded-lg p-1.5 text-gray-400 hover:bg-red-50 hover:text-red-500 outline-none">
                  <Trash2 className="h-4 w-4" />
                </button>
              </div>

              <div className="mt-3 flex items-center gap-2">
                <code className="flex-1 overflow-hidden rounded-lg border border-gray-200 bg-gray-50 px-3 py-2 font-mono text-xs text-gray-700 break-all">
                  {maskKey(k.prefix)}
                </code>
                <button
                  onClick={() => copyToClipboard(k.id, k.prefix)}
                  className="rounded-lg border border-gray-200 p-2 text-gray-400 hover:bg-gray-50 hover:text-gray-600"
                  title="Copy prefix only"
                >
                  {copied === k.id ? (
                    <span className="text-xs font-medium text-green-600 px-1">✓</span>
                  ) : (
                    <Copy className="h-4 w-4" />
                  )}
                </button>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
