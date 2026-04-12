"use client";

import { useState } from "react";
import { Plus, Pencil, Trash2, Copy, Bot } from "lucide-react";
import { useToast } from "@/components/ui/toast";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";

const TEMPLATES = [
  { id: 1, name: "Konfirmasi Pesanan", category: "Transaksi", content: "Halo {{name}} 👋\n\nTerima kasih telah berbelanja! Pesanan #{{order_id}} kamu sudah kami terima.\nEstimasi pengiriman: {{delivery_date}}\n\nInfo lebih lanjut hubungi CS kami. 😊", usedCount: 412 },
  { id: 2, name: "Pengingat Pembayaran", category: "Keuangan", content: "Halo {{name}},\n\nKami mengingatkan bahwa tagihan Anda sebesar Rp {{amount}} jatuh tempo pada {{due_date}}.\n\nSilakan lakukan pembayaran sebelum tanggal tersebut. Terima kasih! 🙏", usedCount: 238 },
  { id: 3, name: "Welcome Message", category: "Umum", content: "Selamat datang di {{brand_name}}! 🎉\n\nKami senang bisa melayani kamu. Ketik *menu* untuk melihat layanan kami atau *cs* untuk berbicara dengan tim kami.", usedCount: 187 },
  { id: 4, name: "Promo Flash Sale", category: "Marketing", content: "🔥 FLASH SALE HARI INI! 🔥\n\nDiskon hingga {{discount}}% untuk semua produk!\nBerlaku sampai {{end_time}} hari ini.\n\nShop now: {{link}}", usedCount: 95 },
];

const CATEGORY_COLORS: Record<string, string> = {
  "Transaksi": "bg-blue-50 text-blue-700",
  "Keuangan": "bg-yellow-50 text-yellow-700",
  "Umum": "bg-gray-100 text-gray-600",
  "Marketing": "bg-purple-50 text-purple-700",
};

export default function TemplatesPage() {
  const { success, info } = useToast();
  const [templates, setTemplates] = useState(TEMPLATES);
  const [copied, setCopied] = useState<number | null>(null);
  const [confirmDelete, setConfirmDelete] = useState<number | null>(null);

  const copyContent = (id: number, content: string) => {
    navigator.clipboard.writeText(content);
    setCopied(id);
    setTimeout(() => setCopied(null), 2000);
    info("Disalin!", "Isi template berhasil disalin ke clipboard.");
  };

  const createTemplate = () => {
    success("Template Dibuat!", "Template baru berhasil ditambahkan.");
  };

  const deleteTemplate = (id: number) => {
    setTemplates(templates.filter((t) => t.id !== id));
    setConfirmDelete(null);
    info("Template Dihapus", "Template berhasil dihapus.");
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold text-gray-900">Message Template</h1>
            <p className="text-sm text-gray-500">Template pesan siap pakai dengan variabel dinamis</p>
          </div>
          <button onClick={createTemplate} className="inline-flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-sm font-semibold text-white hover:bg-green-700 shadow-sm">
            <Plus className="h-4 w-4" /> Buat Template
          </button>
        </div>
      </div>

      <div className="p-6">
        {/* Info */}
        <div className="mb-5 rounded-xl border border-blue-100 bg-blue-50 p-4 text-sm text-blue-800">
          <Bot className="mr-2 inline h-4 w-4" />
          Gunakan variabel seperti <code className="rounded bg-blue-100 px-1">{"{{name}}"}</code>, <code className="rounded bg-blue-100 px-1">{"{{amount}}"}</code> — akan digantikan data nyata saat pesan dikirim.
        </div>

        <div className="grid gap-4 sm:grid-cols-2">
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
                  <p className="mt-0.5 text-xs text-gray-400">Digunakan {t.usedCount} kali</p>
                </div>
                <div className="flex gap-1 flex-shrink-0">
                  <button
                    onClick={() => copyContent(t.id, t.content)}
                    className="rounded-lg p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600"
                    title="Copy"
                  >
                    {copied === t.id ? <span className="text-xs font-medium text-green-600">✓</span> : <Copy className="h-4 w-4" />}
                  </button>
                  <button className="rounded-lg p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600">
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
              <button className="mt-3 w-full rounded-lg border border-green-200 py-2 text-xs font-semibold text-green-700 hover:bg-green-50 transition-colors">
                Gunakan Template Ini
              </button>
            </div>
          ))}
        </div>
      </div>
      <ConfirmDialog
        open={confirmDelete !== null}
        title="Hapus Template?"
        description="Template ini akan dihapus permanen dan tidak dapat dipulihkan."
        confirmLabel="Ya, Hapus"
        onConfirm={() => confirmDelete !== null && deleteTemplate(confirmDelete)}
        onCancel={() => setConfirmDelete(null)}
      />
    </div>
  );
}
