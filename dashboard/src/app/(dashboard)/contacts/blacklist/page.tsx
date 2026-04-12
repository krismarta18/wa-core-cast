"use client";

import { useState } from "react";
import { Ban, Plus, Trash2, X, ShieldAlert } from "lucide-react";
import { useToast } from "@/components/ui/toast";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";

interface BlockedNumber { id: number; phone: string; reason: string; blockedAt: string; }

const INITIAL_BLOCKED: BlockedNumber[] = [
  { id: 1, phone: "628111999001", reason: "Spam / pesan tidak diinginkan", blockedAt: "10 Apr 2026" },
  { id: 2, phone: "628111999002", reason: "Nomor palsu", blockedAt: "09 Apr 2026" },
  { id: 3, phone: "628111999003", reason: "Permintaan pengguna", blockedAt: "07 Apr 2026" },
];

const REASONS = [
  "Spam / pesan tidak diinginkan",
  "Nomor palsu",
  "Permintaan pengguna",
  "Konten berbahaya",
  "Lainnya",
];

export default function BlacklistPage() {
  const { success, info } = useToast();
  const [blocked, setBlocked] = useState<BlockedNumber[]>(INITIAL_BLOCKED);
  const [showAdd, setShowAdd] = useState(false);
  const [form, setForm] = useState({ phone: "", reason: REASONS[0] });
  const [confirmDelete, setConfirmDelete] = useState<number | null>(null);

  function addBlock() {
    if (!form.phone) return;
    setBlocked([{ id: Date.now(), phone: form.phone, reason: form.reason, blockedAt: "Baru saja" }, ...blocked]);
    setForm({ phone: "", reason: REASONS[0] });
    setShowAdd(false);
    success("Nomor Diblokir!", `${form.phone} berhasil ditambahkan ke blacklist.`);
  }

  function removeBlock(id: number) {
    setBlocked(blocked.filter((b) => b.id !== id));
    setConfirmDelete(null);
    info("Blokir Dihapus", "Nomor berhasil dihapus dari blacklist.");
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold text-gray-900">Blacklist / Block</h1>
            <p className="text-sm text-gray-500">Nomor yang diblokir tidak akan menerima pesan apapun</p>
          </div>
          <button onClick={() => setShowAdd(true)} className="flex items-center gap-2 rounded-lg bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700">
            <Plus className="h-4 w-4" /> Tambah Blokir
          </button>
        </div>
      </div>

      <div className="p-6 space-y-4">
        {/* Info banner */}
        <div className="flex items-start gap-3 rounded-xl border border-red-100 bg-red-50 p-4">
          <ShieldAlert className="mt-0.5 h-5 w-5 flex-shrink-0 text-red-500" />
          <p className="text-sm text-red-700">
            Pesan ke nomor yang ada di daftar blacklist akan <strong>otomatis dibatalkan</strong> tanpa notifikasi error. Hapus dari daftar untuk mengizinkan pengiriman kembali.
          </p>
        </div>

        {/* Summary */}
        <div className="inline-flex items-center gap-2 rounded-lg border border-gray-200 bg-white px-4 py-2.5 shadow-sm">
          <Ban className="h-4 w-4 text-red-400" />
          <span className="text-sm font-medium text-gray-700">{blocked.length} nomor diblokir</span>
        </div>

        {/* Table */}
        <div className="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
          {blocked.length === 0 ? (
            <div className="p-12 text-center">
              <Ban className="mx-auto h-10 w-10 text-gray-200" />
              <p className="mt-3 text-sm text-gray-400">Tidak ada nomor yang diblokir</p>
            </div>
          ) : (
            <div className="overflow-x-auto">
            <table className="w-full min-w-[480px] text-sm">
              <thead>
                <tr className="border-b border-gray-100 bg-gray-50 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">
                  <th className="px-5 py-3">Nomor</th>
                  <th className="px-5 py-3">Alasan</th>
                  <th className="px-5 py-3">Diblokir</th>
                  <th className="px-5 py-3 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-50">
                {blocked.map((b) => (
                  <tr key={b.id} className="hover:bg-gray-50">
                    <td className="px-5 py-3">
                      <div className="flex items-center gap-2">
                        <Ban className="h-3.5 w-3.5 text-red-400" />
                        <span className="font-mono text-gray-700">{b.phone}</span>
                      </div>
                    </td>
                    <td className="px-5 py-3 text-gray-500">{b.reason}</td>
                    <td className="px-5 py-3 text-gray-400">{b.blockedAt}</td>
                    <td className="px-5 py-3 text-right">
                      <button
                        onClick={() => setConfirmDelete(b.id)}
                        className="flex items-center gap-1.5 rounded-lg border border-gray-200 px-2.5 py-1.5 text-xs text-gray-500 hover:border-green-400 hover:text-green-600"
                      >
                        <Trash2 className="h-3 w-3" /> Hapus Blokir
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            </div>
          )}
        </div>
      </div>

      <ConfirmDialog
        open={confirmDelete !== null}
        title="Hapus dari Blacklist?"
        description="Nomor ini akan dihapus dari blacklist. Pesan ke nomor ini akan diizinkan kembali."
        confirmLabel="Ya, Hapus"
        onConfirm={() => confirmDelete !== null && removeBlock(confirmDelete)}
        onCancel={() => setConfirmDelete(null)}
      />

      {/* Add modal */}
      {showAdd && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/30 backdrop-blur-sm">
          <div className="w-full max-w-sm rounded-2xl bg-white p-6 shadow-xl">
            <div className="mb-5 flex items-center justify-between">
              <h2 className="text-lg font-bold text-gray-900">Blokir Nomor</h2>
              <button onClick={() => setShowAdd(false)}><X className="h-5 w-5 text-gray-400" /></button>
            </div>
            <div className="space-y-3">
              <div>
                <label className="mb-1 block text-sm font-medium text-gray-700">Nomor WA</label>
                <input
                  value={form.phone}
                  onChange={(e) => setForm({ ...form, phone: e.target.value })}
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-red-400 focus:outline-none"
                  placeholder="628xxxxxxxxxx"
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium text-gray-700">Alasan</label>
                <select
                  value={form.reason}
                  onChange={(e) => setForm({ ...form, reason: e.target.value })}
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-red-400 focus:outline-none"
                >
                  {REASONS.map((r) => <option key={r}>{r}</option>)}
                </select>
              </div>
            </div>
            <div className="mt-5 flex gap-3">
              <button onClick={() => setShowAdd(false)} className="flex-1 rounded-lg border border-gray-200 py-2 text-sm text-gray-600 hover:bg-gray-50">Batal</button>
              <button onClick={addBlock} className="flex-1 rounded-lg bg-red-600 py-2 text-sm font-medium text-white hover:bg-red-700">Blokir</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
