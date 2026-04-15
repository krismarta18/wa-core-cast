"use client";

import { useState, useEffect } from "react";
import { Plus, Pencil, Trash2, Tag, ToggleLeft, ToggleRight, Loader2 } from "lucide-react";
import { useToast } from "@/components/ui/toast";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";
import { autoResponseApi } from "@/lib/api";
import type { AutoResponseKeyword } from "@/lib/types";

export default function KeywordsPage() {
  const { success, error, info } = useToast();
  const [keywords, setKeywords] = useState<AutoResponseKeyword[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ keyword: "", response_text: "" });
  const [editingId, setEditingId] = useState<string | null>(null);
  const [confirmDelete, setConfirmDelete] = useState<string | null>(null);
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    fetchKeywords();
  }, []);

  const fetchKeywords = async () => {
    try {
      setLoading(true);
      const res = await autoResponseApi.getKeywords();
      setKeywords(res.keywords || []);
    } catch (err: any) {
      error("Gagal", err?.response?.data?.error?.message || "Gagal mengambil data keywords");
    } finally {
      setLoading(false);
    }
  };

  const toggle = async (id: string) => {
    try {
      const updated = await autoResponseApi.toggleKeyword(id);
      setKeywords((ks) => ks.map((k) => (k.id === id ? updated : k)));
    } catch (err: any) {
      error("Gagal", "Gagal mengubah status keyword");
    }
  };

  const remove = async (id: string) => {
    try {
      await autoResponseApi.deleteKeyword(id);
      setKeywords((ks) => ks.filter((k) => k.id !== id));
      setConfirmDelete(null);
      info("Keyword Dihapus", "Auto-reply keyword berhasil dihapus.");
    } catch (err: any) {
      error("Gagal", "Gagal menghapus keyword");
    }
  };

  const saveKeyword = async () => {
    if (!form.keyword || !form.response_text) return;

    try {
      setIsSaving(true);
      if (editingId) {
        const updated = await autoResponseApi.updateKeyword(editingId, {
          keyword: form.keyword,
          response_text: form.response_text,
        });
        setKeywords((ks) => ks.map((k) => (k.id === editingId ? updated : k)));
        success("Berhasil", `Keyword "${form.keyword}" berhasil diperbarui.`);
      } else {
        const created = await autoResponseApi.createKeyword({
          keyword: form.keyword,
          response_text: form.response_text,
        });
        setKeywords((ks) => [created, ...ks]);
        success("Berhasil", `Keyword "${form.keyword}" berhasil ditambahkan.`);
      }
      
      setForm({ keyword: "", response_text: "" });
      setEditingId(null);
      setShowForm(false);
    } catch (err: any) {
      error("Gagal", "Gagal menyimpan keyword");
    } finally {
      setIsSaving(false);
    }
  };

  const editKeyword = (kw: AutoResponseKeyword) => {
    setForm({ keyword: kw.keyword, response_text: kw.response_text });
    setEditingId(kw.id);
    setShowForm(true);
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
              onClick={() => {
                setForm({ keyword: "", response_text: "" });
                setEditingId(null);
                setShowForm(true);
              }}
              className="inline-flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-sm font-semibold text-white hover:bg-green-700 shadow-sm"
              disabled={showForm}
            >
              <Plus className="h-4 w-4" /> Tambah Keyword
            </button>
          </div>
        </div>

        <div className="p-6 space-y-3">
          {loading ? (
            <div className="flex h-40 items-center justify-center">
              <Loader2 className="h-8 w-8 animate-spin text-green-500" />
            </div>
          ) : (
            <>
              {/* Add/Edit form */}
              {showForm && (
                <div className="rounded-xl border border-green-200 bg-green-50 p-5 shadow-sm">
                  <h3 className="mb-3 font-semibold text-green-800">
                    {editingId ? "Edit Keyword" : "Keyword Baru"}
                  </h3>
                  <div className="space-y-3">
                    <div>
                      <label className="mb-1 block text-sm font-medium text-gray-700">Kata Kunci</label>
                      <input
                        type="text"
                        placeholder="cth: harga, promo, info..."
                        value={form.keyword}
                        onChange={(e) => setForm({ ...form, keyword: e.target.value.toLowerCase() })}
                        className="h-10 w-full rounded-lg border border-gray-300 px-3 text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
                        disabled={isSaving}
                      />
                    </div>
                    <div>
                      <label className="mb-1 block text-sm font-medium text-gray-700">Pesan Balasan</label>
                      <textarea
                        rows={3}
                        placeholder="Pesan yang akan dikirim saat keyword terdeteksi..."
                        value={form.response_text}
                        onChange={(e) => setForm({ ...form, response_text: e.target.value })}
                        className="w-full resize-none rounded-lg border border-gray-300 p-3 text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
                        disabled={isSaving}
                      />
                    </div>
                    <div className="flex gap-2">
                      <button
                        onClick={saveKeyword}
                        disabled={!form.keyword || !form.response_text || isSaving}
                        className="flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-sm font-semibold text-white hover:bg-green-700 disabled:opacity-50"
                      >
                        {isSaving && <Loader2 className="h-4 w-4 animate-spin" />}
                        Simpan
                      </button>
                      <button
                        onClick={() => {
                          setShowForm(false);
                          setForm({ keyword: "", response_text: "" });
                          setEditingId(null);
                        }}
                        disabled={isSaving}
                        className="rounded-lg border border-gray-200 px-4 py-2 text-sm font-semibold text-gray-600 hover:bg-gray-50"
                      >
                        Batal
                      </button>
                    </div>
                  </div>
                </div>
              )}

              {keywords.length === 0 && !showForm && (
                <div className="rounded-xl border border-dashed border-gray-300 bg-white p-12 text-center shadow-sm">
                  <Tag className="mx-auto h-8 w-8 text-gray-400" />
                  <h3 className="mt-4 text-sm font-semibold text-gray-900">Belum ada Keyword</h3>
                  <p className="mt-2 text-sm text-gray-500">
                    Kamu belum memiliki keyword auto-response aktif.
                  </p>
                </div>
              )}

              {keywords.map((k) => (
                <div
                  key={k.id}
                  className={`rounded-xl border bg-white p-5 shadow-sm transition-opacity ${
                    k.is_active ? "border-gray-200" : "border-gray-100 opacity-60"
                  }`}
                >
                  <div className="flex flex-wrap items-start justify-between gap-3">
                    <div className="flex items-center gap-2">
                      <Tag className="h-4 w-4 text-green-600 flex-shrink-0" />
                      <code className="rounded-md bg-gray-100 px-2 py-0.5 text-sm font-bold text-gray-800">
                        {k.keyword}
                      </code>
                      <span
                        className={`rounded-full px-2 py-0.5 text-xs font-medium ${
                          k.is_active ? "bg-green-50 text-green-700" : "bg-gray-100 text-gray-500"
                        }`}
                      >
                        {k.is_active ? "Aktif" : "Nonaktif"}
                      </span>
                    </div>
                    <div className="flex items-center gap-1">
                      <button onClick={() => toggle(k.id)}>
                        {k.is_active ? (
                          <ToggleRight className="h-6 w-6 text-green-600" />
                        ) : (
                          <ToggleLeft className="h-6 w-6 text-gray-400" />
                        )}
                      </button>
                      <button 
                        onClick={() => editKeyword(k)}
                        className="rounded-lg p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600"
                      >
                        <Pencil className="h-4 w-4" />
                      </button>
                      <button
                        onClick={() => setConfirmDelete(k.id)}
                        className="rounded-lg p-1.5 text-gray-400 hover:bg-red-50 hover:text-red-500"
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </div>
                  </div>
                  <p className="mt-3 whitespace-pre-wrap rounded-lg bg-gray-50 px-4 py-3 text-sm text-gray-700 leading-relaxed">
                    {k.response_text}
                  </p>
                </div>
              ))}
            </>
          )}
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
