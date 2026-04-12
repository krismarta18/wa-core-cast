"use client";

import { useState } from "react";
import { Users, UserPlus, ChevronDown, ChevronRight, Trash2, Plus, X } from "lucide-react";
import { useToast } from "@/components/ui/toast";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";

interface Member { id: number; name: string; phone: string; }
interface Group { id: number; name: string; description: string; members: Member[]; }

const INITIAL_GROUPS: Group[] = [
  {
    id: 1,
    name: "Pelanggan VIP",
    description: "Pelanggan tier premium",
    members: [
      { id: 1, name: "Budi Santoso", phone: "628111000001" },
      { id: 2, name: "Dewi Lestari", phone: "628111000004" },
    ],
  },
  {
    id: 2,
    name: "Tim Internal",
    description: "Staff dan karyawan",
    members: [
      { id: 5, name: "Hendra Kurniawan", phone: "628111000005" },
    ],
  },
  {
    id: 3,
    name: "Prospek Aktif",
    description: "Calon pelanggan yang sedang difollow up",
    members: [
      { id: 3, name: "Ahmad Rizki", phone: "628111000003" },
      { id: 6, name: "Rina Widiastuti", phone: "628111000006" },
    ],
  },
];

export default function ContactGroupsPage() {
  const { success, info } = useToast();
  const [groups, setGroups] = useState<Group[]>(INITIAL_GROUPS);
  const [expanded, setExpanded] = useState<number | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [newGroup, setNewGroup] = useState({ name: "", description: "" });
  const [confirmDelete, setConfirmDelete] = useState<number | null>(null);

  function toggleExpand(id: number) {
    setExpanded(expanded === id ? null : id);
  }

  function createGroup() {
    if (!newGroup.name) return;
    setGroups([...groups, { id: Date.now(), ...newGroup, members: [] }]);
    setNewGroup({ name: "", description: "" });
    setShowCreate(false);
    success("Grup Dibuat!", `Grup "${newGroup.name}" berhasil ditambahkan.`);
  }

  function deleteGroup(id: number) {
    setGroups(groups.filter((g) => g.id !== id));
    if (expanded === id) setExpanded(null);
    setConfirmDelete(null);
    info("Grup Dihapus", "Grup kontak berhasil dihapus.");
  }

  function removeMember(groupId: number, memberId: number) {
    setGroups(groups.map((g) =>
      g.id === groupId ? { ...g, members: g.members.filter((m) => m.id !== memberId) } : g
    ));
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold text-gray-900">Group Kontak</h1>
            <p className="text-sm text-gray-500">Kelompokkan kontak untuk broadcast massal</p>
          </div>
          <button onClick={() => setShowCreate(true)} className="flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700">
            <Plus className="h-4 w-4" /> Buat Grup
          </button>
        </div>
      </div>

      <div className="p-6 space-y-3">
        {/* Stats */}
        <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 mb-2">
          <div className="rounded-xl border border-gray-200 bg-white p-4 shadow-sm">
            <p className="text-xs text-gray-500">Total Grup</p>
            <p className="mt-1 text-2xl font-bold text-gray-900">{groups.length}</p>
          </div>
          <div className="rounded-xl border border-gray-200 bg-white p-4 shadow-sm">
            <p className="text-xs text-gray-500">Total Member</p>
            <p className="mt-1 text-2xl font-bold text-gray-900">{groups.reduce((a, g) => a + g.members.length, 0)}</p>
          </div>
        </div>

        {/* Group list */}
        {groups.map((g) => (
          <div key={g.id} className="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
            <div
              className="flex cursor-pointer items-center justify-between px-5 py-4 hover:bg-gray-50"
              onClick={() => toggleExpand(g.id)}
            >
              <div className="flex items-center gap-3">
                <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-green-50">
                  <Users className="h-5 w-5 text-green-600" />
                </div>
                <div>
                  <p className="font-semibold text-gray-900">{g.name}</p>
                  <p className="text-xs text-gray-400">{g.description}</p>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <span className="rounded-full bg-gray-100 px-2.5 py-1 text-xs font-medium text-gray-600">
                  {g.members.length} anggota
                </span>
                {expanded === g.id ? <ChevronDown className="h-4 w-4 text-gray-400" /> : <ChevronRight className="h-4 w-4 text-gray-400" />}
              </div>
            </div>

            {expanded === g.id && (
              <div className="border-t border-gray-100 bg-gray-50 px-5 py-3">
                {g.members.length === 0 ? (
                  <p className="py-4 text-center text-sm text-gray-400">Belum ada anggota</p>
                ) : (
                  <div className="overflow-x-auto">
                  <table className="w-full min-w-[400px] text-sm">
                    <thead>
                      <tr className="text-left text-xs font-semibold uppercase tracking-wider text-gray-400">
                        <th className="pb-2">Nama</th>
                        <th className="pb-2">Nomor</th>
                        <th className="pb-2 text-right">Aksi</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-100">
                      {g.members.map((m) => (
                        <tr key={m.id} className="hover:bg-white">
                          <td className="py-2 text-gray-700">{m.name}</td>
                          <td className="py-2 font-mono text-gray-500">{m.phone}</td>
                          <td className="py-2 text-right">
                            <button onClick={() => removeMember(g.id, m.id)} className="rounded p-1 text-gray-400 hover:bg-red-50 hover:text-red-500">
                              <Trash2 className="h-3.5 w-3.5" />
                            </button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                  </div>
                )}
                <div className="mt-3 flex gap-2">
                  <button className="flex items-center gap-1.5 rounded-lg border border-dashed border-gray-300 px-3 py-1.5 text-xs text-gray-500 hover:border-green-400 hover:text-green-600">
                    <UserPlus className="h-3.5 w-3.5" /> Tambah Anggota
                  </button>
                  <button onClick={() => setConfirmDelete(g.id)} className="ml-auto flex items-center gap-1.5 rounded-lg border border-dashed border-red-200 px-3 py-1.5 text-xs text-red-400 hover:border-red-400 hover:text-red-600">
                    <Trash2 className="h-3.5 w-3.5" /> Hapus Grup
                  </button>
                </div>
              </div>
            )}
          </div>
        ))}

        {groups.length === 0 && (
          <div className="rounded-xl border border-dashed border-gray-300 bg-white p-12 text-center">
            <Users className="mx-auto h-10 w-10 text-gray-300" />
            <p className="mt-3 text-gray-500">Belum ada grup kontak</p>
          </div>
        )}
      </div>

      <ConfirmDialog
        open={confirmDelete !== null}
        title="Hapus Grup?"
        description="Grup dan semua anggotanya akan dihapus permanen."
        confirmLabel="Ya, Hapus"
        onConfirm={() => confirmDelete !== null && deleteGroup(confirmDelete)}
        onCancel={() => setConfirmDelete(null)}
      />

      {/* Create modal */}
      {showCreate && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/30 backdrop-blur-sm">
          <div className="w-full max-w-sm rounded-2xl bg-white p-6 shadow-xl">
            <div className="mb-5 flex items-center justify-between">
              <h2 className="text-lg font-bold text-gray-900">Buat Grup Baru</h2>
              <button onClick={() => setShowCreate(false)}><X className="h-5 w-5 text-gray-400" /></button>
            </div>
            <div className="space-y-3">
              <div>
                <label className="mb-1 block text-sm font-medium text-gray-700">Nama Grup</label>
                <input value={newGroup.name} onChange={(e) => setNewGroup({ ...newGroup, name: e.target.value })} className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-green-500 focus:outline-none" placeholder="Contoh: Pelanggan VIP" />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium text-gray-700">Deskripsi</label>
                <input value={newGroup.description} onChange={(e) => setNewGroup({ ...newGroup, description: e.target.value })} className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-green-500 focus:outline-none" placeholder="Opsional" />
              </div>
            </div>
            <div className="mt-5 flex gap-3">
              <button onClick={() => setShowCreate(false)} className="flex-1 rounded-lg border border-gray-200 py-2 text-sm text-gray-600 hover:bg-gray-50">Batal</button>
              <button onClick={createGroup} className="flex-1 rounded-lg bg-green-600 py-2 text-sm font-medium text-white hover:bg-green-700">Buat</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
