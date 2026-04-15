"use client";

import { useState, useEffect } from "react";
import { Plus, Pencil, Trash2, Copy, Bot, Loader2 } from "lucide-react";
import { useToast } from "@/components/ui/toast";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";
import { autoResponseApi } from "@/lib/api";
import type { MessageTemplate } from "@/lib/types";

const CATEGORY_COLORS: Record<string, string> = {
  "Transaksi": "bg-blue-50 text-blue-700",
  "Keuangan": "bg-yellow-50 text-yellow-700",
  "Umum": "bg-gray-100 text-gray-600",
  "Marketing": "bg-purple-50 text-purple-700",
};

export default function TemplatesPage() {
  const { success, error, info } = useToast();
  const [templates, setTemplates] = useState<MessageTemplate[]>([]);
  const [loading, setLoading] = useState(true);
  
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ name: "", category: "Umum", content: "" });
  const [editingId, setEditingId] = useState<string | null>(null);
  
  const [copied, setCopied] = useState<string | null>(null);
  const [confirmDelete, setConfirmDelete] = useState<string | null>(null);
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    fetchTemplates();
  }, []);

  const fetchTemplates = async () => {
    try {
      setLoading(true);
      const res = await autoResponseApi.getTemplates();
      setTemplates(res.templates || []);
    } catch (err: any) {
      error("Gagal", "Gagal mengambil data templates");
    } finally {
      setLoading(false);
    }
  };

  const copyContent = (id: string, content: string) => {
    navigator.clipboard.writeText(content);
    setCopied(id);
    setTimeout(() => setCopied(null), 2000);
    info("Disalin!", "Isi template berhasil disalin ke clipboard.");
  };

  const saveTemplate = async () => {
    if (!form.name || !form.content) return;

    try {
      setIsSaving(true);
      if (editingId) {
        const updated = await autoResponseApi.updateTemplate(editingId, form);
        setTemplates((ts) => ts.map((t) => (t.id === editingId ? updated : t)));
        success("Berhasil", `Template "${form.name}" diperbarui.`);
      } else {
        const created = await autoResponseApi.createTemplate(form);
        setTemplates((ts) => [created, ...ts]);
        success("Berhasil", "Template baru berhasil ditambahkan.");
      }
      
      setForm({ name: "", category: "Umum", content: "" });
      setEditingId(null);
      setShowForm(false);
    } catch (err: any) {
      error("Gagal", "Gagal menyimpan template");
    } finally {
      setIsSaving(false);
    }
  };

  const editTemplate = (t: MessageTemplate) => {
    setForm({ name: t.name, category: t.category, content: t.content });
    setEditingId(t.id);
    setShowForm(true);
  };

  const remove = async (id: string) => {
    try {
      await autoResponseApi.deleteTemplate(id);
      setTemplates((ts) => ts.filter((t) => t.id !== id));
      setConfirmDelete(null);
      info("Template Dihapus", "Template berhasil dihapus permanen.");
    } catch (err: any) {
      error("Gagal", "Gagal menghapus template");
    }
  };

  return (
    <>
      <div className="min-h-screen bg-gray-50">
        <div className="border-b border-gray-200 bg-white px-6 py-4">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-xl font-bold text-gray-900">Message Template</h1>
              <p className="text-sm text-gray-500">Template pesan siap pakai dengan variabel dinamis</p>
            </div>
            <button 
              onClick={() => {
                setForm({ name: "", category: "Umum", content: "" });
                setEditingId(null);
                setShowForm(true);
              }}
              className="inline-flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-sm font-semibold text-white hover:bg-green-700 shadow-sm"
              disabled={showForm}  
            >
              <Plus className="h-4 w-4" /> Buat Template
            </button>
          </div>
        </div>

        <div className="p-6">
          <div className="mb-5 rounded-xl border border-blue-100 bg-blue-50 p-4 text-sm text-blue-800">
            <Bot className="mr-2 inline h-4 w-4" />
            Gunakan variabel seperti <code className="rounded bg-blue-100 px-1">{"{{name}}"}</code>, <code className="rounded bg-blue-100 px-1">{"{{amount}}"}</code> — template disiapkan untuk auto-replace variabel (future feature).
          </div>

          {loading ? (
            <div className="flex h-40 items-center justify-center">
              <Loader2 className="h-8 w-8 animate-spin text-green-500" />
            </div>
          ) : (
            <div className="grid gap-4 sm:grid-cols-2">
              {/* Form Input Container as a Card */}
              {showForm && (
                <div className="rounded-xl border border-green-200 bg-green-50 p-5 shadow-sm flex flex-col sm:col-span-2 md:col-span-1 lg:col-span-2 xl:col-span-2">
                  <h3 className="mb-3 font-semibold text-green-800">
                    {editingId ? "Edit Template" : "Template Baru"}
                  </h3>
                  <div className="space-y-3">
                    <div className="grid grid-cols-2 gap-3">
                      <div>
                        <label className="mb-1 block text-sm font-medium text-gray-700">Nama Template</label>
                        <input
                          type="text"
                          value={form.name}
                          onChange={(e) => setForm({ ...form, name: e.target.value })}
                          className="h-10 w-full rounded-lg border border-gray-300 px-3 text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
                          disabled={isSaving}
                        />
                      </div>
                      <div>
                        <label className="mb-1 block text-sm font-medium text-gray-700">Kategori</label>
                        <select
                          value={form.category}
                          onChange={(e) => setForm({ ...form, category: e.target.value })}
                          className="h-10 w-full rounded-lg border border-gray-300 bg-white px-3 text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
                          disabled={isSaving}
                        >
                          <option value="Umum">Umum</option>
                          <option value="Transaksi">Transaksi</option>
                          <option value="Keuangan">Keuangan</option>
                          <option value="Marketing">Marketing</option>
                        </select>
                      </div>
                    </div>
                    <div>
                      <label className="mb-1 block text-sm font-medium text-gray-700">Isi Pesan</label>
                      <textarea
                        rows={4}
                        value={form.content}
                        onChange={(e) => setForm({ ...form, content: e.target.value })}
                        className="w-full resize-none rounded-lg border border-gray-300 p-3 text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
                        disabled={isSaving}
                      />
                    </div>
                    <div className="flex gap-2 pt-2">
                      <button
                        onClick={saveTemplate}
                        disabled={!form.name || !form.content || isSaving}
                        className="flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-sm font-semibold text-white hover:bg-green-700 disabled:opacity-50"
                      >
                        {isSaving && <Loader2 className="h-4 w-4 animate-spin" />}
                        Simpan
                      </button>
                      <button
                        onClick={() => {
                          setShowForm(false);
                          setForm({ name: "", category: "Umum", content: "" });
                          setEditingId(null);
                        }}
                        disabled={isSaving}
                        className="rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-semibold text-gray-600 hover:bg-gray-50"
                      >
                        Batal
                      </button>
                    </div>
                  </div>
                </div>
              )}

              {/* Empty state */}
              {templates.length === 0 && !showForm && (
                <div className="sm:col-span-2 rounded-xl border border-dashed border-gray-300 bg-white p-12 text-center shadow-sm">
                  <h3 className="text-sm font-semibold text-gray-900">Belum ada Template</h3>
                  <p className="mt-2 text-sm text-gray-500">Anda belum membuat message template sama sekali.</p>
                </div>
              )}

              {/* Template Cards */}
              {templates.map((t) => (
                <div key={t.id} className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm flex flex-col">
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <div className="flex items-center gap-2 flex-wrap">
                        <h3 className="font-semibold text-gray-900">{t.name}</h3>
                        <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${CATEGORY_COLORS[t.category] ?? "bg-gray-100 text-gray-600"}`}>
                          {t.category}
                        </span>
                      </div>
                      <p className="mt-0.5 text-xs text-gray-400">Digunakan {t.used_count} kali</p>
                    </div>
                    <div className="flex gap-1 flex-shrink-0">
                      <button
                        onClick={() => copyContent(t.id, t.content)}
                        className="rounded-lg p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600"
                        title="Copy text"
                      >
                        {copied === t.id ? <span className="text-xs font-medium text-green-600">✓</span> : <Copy className="h-4 w-4" />}
                      </button>
                      <button 
                        onClick={() => editTemplate(t)}
                        className="rounded-lg p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600"
                      >
                        <Pencil className="h-4 w-4" />
                      </button>
                      <button
                        onClick={() => setConfirmDelete(t.id)}
                        className="rounded-lg p-1.5 text-gray-400 hover:bg-red-50 hover:text-red-500"
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </div>
                  </div>
                  <div className="mt-3 flex-1 overflow-hidden rounded-lg bg-gray-50 px-4 py-3">
                    <p className="whitespace-pre-wrap text-sm text-gray-700 leading-relaxed line-clamp-4">
                      {t.content}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
      <ConfirmDialog
        open={confirmDelete !== null}
        title="Hapus Template?"
        description="Template ini akan dihapus permanen dan tidak dapat dipulihkan."
        confirmLabel="Ya, Hapus"
        onConfirm={() => confirmDelete !== null && remove(confirmDelete)}
        onCancel={() => setConfirmDelete(null)}
      />
    </>
  );
}
