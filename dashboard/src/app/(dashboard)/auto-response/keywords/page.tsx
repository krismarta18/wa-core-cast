"use client";

import { useState } from "react";
import { Plus, Pencil, Trash2, Tag, ToggleLeft, ToggleRight } from "lucide-react";
import { useToast } from "@/components/ui/toast";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";

const KEYWORDS = [
  { id: 1, keyword: "harga", response: "Halo! Untuk info harga silakan kunjungi website kami di wacast.id/pricing 😊", active: true },
  { id: 2, keyword: "promo", response: "Promo terbaru kami tersedia di wacast.id/promo. Jangan sampai ketinggalan! 🎉", active: true },
  { id: 3, keyword: "cs", response: "Tim Customer Service kami siap membantu!\nHub: 628111000222\nJam: 08.00 - 17.00 WIB", active: true },
  { id: 4, keyword: "lokasi", response: "Alamat kami: Jl. Sudirman No. 45, Jakarta Pusat. Buka Senin–Sabtu 09.00–18.00", active: false },
  { id: 5, keyword: "tracking", response: "Untuk cek status pengiriman, reply dengan format: CEKRESI [nomor resi]", active: true },
];

export default function KeywordsPage() {
  const { success, info } = useToast();
  const [keywords, setKeywords] = useState(KEYWORDS);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ keyword: "", response: "" });
  const [confirmDelete, setConfirmDelete] = useState<number | null>(null);

  const toggle = (id: number) =>
    setKeywords((ks) => ks.map((k) => k.id === id ? { ...k, active: !k.active } : k));

  const remove = (id: number) => {
    setKeywords((ks) => ks.filter((k) => k.id !== id));
    setConfirmDelete(null);
    info("Keyword Dihapus", "Auto-reply keyword berhasil dihapus.");
  };

  return (
    <>
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold text-gray-900">Auto Response — Keyword</h1>
            <p className="text-sm text-gray-500">Balas otomatis berdasarkan kata kunci yang dikirim pengguna</p>
          </div>
          <button
            onClick={() => setShowForm(true)}
            className="inline-flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-sm font-semibold text-white hover:bg-green-700 shadow-sm"
          >
            <Plus className="h-4 w-4" /> Tambah Keyword
          </button>
        </div>
      </div>

      <div className="p-6 space-y-3">
        {/* Add form */}
        {showForm && (
          <div className="rounded-xl border border-green-200 bg-green-50 p-5 shadow-sm">
            <h3 className="mb-3 font-semibold text-green-800">Keyword Baru</h3>
            <div className="space-y-3">
              <div>
                <label className="mb-1 block text-sm font-medium text-gray-700">Kata Kunci</label>
                <input
                  type="text"
                  placeholder="cth: harga, promo, info..."
                  value={form.keyword}
                  onChange={(e) => setForm({ ...form, keyword: e.target.value.toLowerCase() })}
                  className="h-10 w-full rounded-lg border border-gray-300 px-3 text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium text-gray-700">Pesan Balasan</label>
                <textarea
                  rows={3}
                  placeholder="Pesan yang akan dikirim saat keyword terdeteksi..."
                  value={form.response}
                  onChange={(e) => setForm({ ...form, response: e.target.value })}
                  className="w-full resize-none rounded-lg border border-gray-300 p-3 text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
                />
              </div>
              <div className="flex gap-2">
                <button
                  onClick={() => {
                    if (form.keyword && form.response) {
                      setKeywords((ks) => [...ks, { id: Date.now(), ...form, active: true }]);
                      setForm({ keyword: "", response: "" });
                      setShowForm(false);
                      success("Keyword Ditambahkan!", `Keyword "${form.keyword}" berhasil disimpan.`);
                    }
                  }}
                  className="rounded-lg bg-green-600 px-4 py-2 text-sm font-semibold text-white hover:bg-green-700"
                >
                  Simpan
                </button>
                <button
                  onClick={() => { setShowForm(false); setForm({ keyword: "", response: "" }); }}
                  className="rounded-lg border border-gray-200 px-4 py-2 text-sm font-semibold text-gray-600 hover:bg-gray-50"
                >
                  Batal
                </button>
              </div>
            </div>
          </div>
        )}

        {keywords.map((k) => (
          <div key={k.id} className={`rounded-xl border bg-white p-5 shadow-sm transition-opacity ${k.active ? "border-gray-200" : "border-gray-100 opacity-60"}`}>
            <div className="flex flex-wrap items-start justify-between gap-3">
              <div className="flex items-center gap-2">
                <Tag className="h-4 w-4 text-green-600 flex-shrink-0" />
                <code className="rounded-md bg-gray-100 px-2 py-0.5 text-sm font-bold text-gray-800">
                  {k.keyword}
                </code>
                <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${k.active ? "bg-green-50 text-green-700" : "bg-gray-100 text-gray-500"}`}>
                  {k.active ? "Aktif" : "Nonaktif"}
                </span>
              </div>
              <div className="flex items-center gap-1">
                <button onClick={() => toggle(k.id)}>
                  {k.active
                    ? <ToggleRight className="h-6 w-6 text-green-600" />
                    : <ToggleLeft className="h-6 w-6 text-gray-400" />}
                </button>
                <button className="rounded-lg p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600">
                  <Pencil className="h-4 w-4" />
                </button>
                <button onClick={() => setConfirmDelete(k.id)} className="rounded-lg p-1.5 text-gray-400 hover:bg-red-50 hover:text-red-500">
                  <Trash2 className="h-4 w-4" />
                </button>
              </div>
            </div>
            <p className="mt-3 whitespace-pre-wrap rounded-lg bg-gray-50 px-4 py-3 text-sm text-gray-700 leading-relaxed">
              {k.response}
            </p>
          </div>
        ))}
      </div>
    </div>
    <ConfirmDialog
      open={confirmDelete !== null}
      title="Hapus Keyword?"
      description="Keyword auto-reply ini akan dihapus permanen."
      confirmLabel="Ya, Hapus"
      onConfirm={() => confirmDelete !== null && remove(confirmDelete)}
      onCancel={() => setConfirmDelete(null)}
    />
    </>
  );
}
