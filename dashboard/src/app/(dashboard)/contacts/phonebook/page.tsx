"use client";

import { useState } from "react";
import { Search, UserPlus, Tag, Phone, Trash2, Edit2, X } from "lucide-react";
import { useToast } from "@/components/ui/toast";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";

const LABELS = ["VIP", "Pelanggan", "Prospek", "Supplier", "Internal"];

const INITIAL_CONTACTS = [
  { id: 1, name: "Budi Santoso", phone: "628111000001", label: "VIP", note: "Pelanggan setia" },
  { id: 2, name: "Siti Aminah", phone: "628111000002", label: "Pelanggan", note: "" },
  { id: 3, name: "Ahmad Rizki", phone: "628111000003", label: "Prospek", note: "Follow up Senin" },
  { id: 4, name: "Dewi Lestari", phone: "628111000004", label: "Supplier", note: "" },
  { id: 5, name: "Hendra Kurniawan", phone: "628111000005", label: "Internal", note: "Tim CS" },
  { id: 6, name: "Rina Widiastuti", phone: "628111000006", label: "Pelanggan", note: "" },
];

interface Contact { id: number; name: string; phone: string; label: string; note: string; }

export default function PhonebookPage() {
  const { success, info } = useToast();
  const [contacts, setContacts] = useState<Contact[]>(INITIAL_CONTACTS);
  const [search, setSearch] = useState("");
  const [filterLabel, setFilterLabel] = useState("all");
  const [showAdd, setShowAdd] = useState(false);
  const [form, setForm] = useState({ name: "", phone: "", label: "Pelanggan", note: "" });
  const [confirmDelete, setConfirmDelete] = useState<number | null>(null);

  const filtered = contacts.filter((c) => {
    const matchSearch = c.name.toLowerCase().includes(search.toLowerCase()) || c.phone.includes(search);
    const matchLabel = filterLabel === "all" || c.label === filterLabel;
    return matchSearch && matchLabel;
  });

  const LABEL_COLOR: Record<string, string> = {
    VIP: "bg-yellow-50 text-yellow-700",
    Pelanggan: "bg-blue-50 text-blue-700",
    Prospek: "bg-purple-50 text-purple-700",
    Supplier: "bg-orange-50 text-orange-700",
    Internal: "bg-green-50 text-green-700",
  };

  function addContact() {
    if (!form.name || !form.phone) return;
    setContacts([...contacts, { id: Date.now(), ...form }]);
    setForm({ name: "", phone: "", label: "Pelanggan", note: "" });
    setShowAdd(false);
    success("Kontak Ditambahkan!", `${form.name} berhasil disimpan ke phonebook.`);
  }

  function deleteContact(id: number) {
    setContacts(contacts.filter((c) => c.id !== id));
    setConfirmDelete(null);
    info("Kontak Dihapus", "Kontak berhasil dihapus dari phonebook.");
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold text-gray-900">Phone Book</h1>
            <p className="text-sm text-gray-500">Kelola kontak dan label pengiriman</p>
          </div>
          <button onClick={() => setShowAdd(true)} className="flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700">
            <UserPlus className="h-4 w-4" /> Tambah Kontak
          </button>
        </div>
      </div>

      <div className="p-6 space-y-4">
        {/* Filters */}
        <div className="flex flex-wrap gap-3">
          <div className="relative flex-1 min-w-[200px]">
            <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
            <input
              type="text"
              placeholder="Cari nama atau nomor..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full rounded-lg border border-gray-200 py-2 pl-9 pr-4 text-sm focus:border-green-500 focus:outline-none"
            />
          </div>
          <select
            value={filterLabel}
            onChange={(e) => setFilterLabel(e.target.value)}
            className="rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-green-500 focus:outline-none"
          >
            <option value="all">Semua Label</option>
            {LABELS.map((l) => <option key={l} value={l}>{l}</option>)}
          </select>
          <div className="flex items-center gap-2 text-sm text-gray-500">
            <Phone className="h-4 w-4" /> <span>{filtered.length} kontak</span>
          </div>
        </div>

        {/* Table */}
        <div className="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
          <div className="overflow-x-auto">
          <table className="w-full min-w-[600px] text-sm">
            <thead>
              <tr className="border-b border-gray-100 bg-gray-50 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">
                <th className="px-5 py-3">Nama</th>
                <th className="px-5 py-3">Nomor</th>
                <th className="px-5 py-3">Label</th>
                <th className="px-5 py-3">Catatan</th>
                <th className="px-5 py-3 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-50">
              {filtered.map((c) => (
                <tr key={c.id} className="hover:bg-gray-50">
                  <td className="px-5 py-3 font-medium text-gray-800">{c.name}</td>
                  <td className="px-5 py-3 font-mono text-gray-600">{c.phone}</td>
                  <td className="px-5 py-3">
                    <span className={`inline-flex items-center gap-1 rounded-full px-2.5 py-1 text-xs font-medium ${LABEL_COLOR[c.label] ?? "bg-gray-100 text-gray-600"}`}>
                      <Tag className="h-3 w-3" /> {c.label}
                    </span>
                  </td>
                  <td className="px-5 py-3 text-gray-500">{c.note || "—"}</td>
                  <td className="px-5 py-3 text-right">
                    <div className="flex items-center justify-end gap-1">
                      <button className="rounded p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600">
                        <Edit2 className="h-3.5 w-3.5" />
                      </button>
                      <button onClick={() => setConfirmDelete(c.id)} className="rounded p-1.5 text-gray-400 hover:bg-red-50 hover:text-red-500">
                        <Trash2 className="h-3.5 w-3.5" />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
              {filtered.length === 0 && (
                <tr><td colSpan={5} className="px-5 py-10 text-center text-gray-400">Tidak ada kontak ditemukan</td></tr>
              )}
            </tbody>
          </table>
          </div>
        </div>
      </div>

      <ConfirmDialog
        open={confirmDelete !== null}
        title="Hapus Kontak?"
        description="Kontak ini akan dihapus dari phonebook secara permanen."
        confirmLabel="Ya, Hapus"
        onConfirm={() => confirmDelete !== null && deleteContact(confirmDelete)}
        onCancel={() => setConfirmDelete(null)}
      />

      {/* Add modal */}
      {showAdd && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/30 backdrop-blur-sm">
          <div className="w-full max-w-md rounded-2xl bg-white p-6 shadow-xl">
            <div className="mb-5 flex items-center justify-between">
              <h2 className="text-lg font-bold text-gray-900">Tambah Kontak</h2>
              <button onClick={() => setShowAdd(false)}><X className="h-5 w-5 text-gray-400" /></button>
            </div>
            <div className="space-y-3">
              <div>
                <label className="mb-1 block text-sm font-medium text-gray-700">Nama</label>
                <input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-green-500 focus:outline-none" placeholder="Nama lengkap" />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium text-gray-700">Nomor WA</label>
                <input value={form.phone} onChange={(e) => setForm({ ...form, phone: e.target.value })} className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-green-500 focus:outline-none" placeholder="628xxxxxxxxxx" />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium text-gray-700">Label</label>
                <select value={form.label} onChange={(e) => setForm({ ...form, label: e.target.value })} className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-green-500 focus:outline-none">
                  {LABELS.map((l) => <option key={l}>{l}</option>)}
                </select>
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium text-gray-700">Catatan</label>
                <input value={form.note} onChange={(e) => setForm({ ...form, note: e.target.value })} className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-green-500 focus:outline-none" placeholder="Opsional" />
              </div>
            </div>
            <div className="mt-5 flex gap-3">
              <button onClick={() => setShowAdd(false)} className="flex-1 rounded-lg border border-gray-200 py-2 text-sm text-gray-600 hover:bg-gray-50">Batal</button>
              <button onClick={addContact} className="flex-1 rounded-lg bg-green-600 py-2 text-sm font-medium text-white hover:bg-green-700">Simpan</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
