"use client";

import { useState, useEffect } from "react";
import { Users, UserPlus, ChevronDown, ChevronRight, Trash2, Plus, X, Loader2 } from "lucide-react";
import { useToast } from "@/components/ui/toast";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";
import { groupsApi, contactsApi } from "@/lib/api";
import { Contact, ContactGroup } from "@/lib/types";

interface GroupWithMembers extends ContactGroup {
  members: Contact[];
}

export default function ContactGroupsPage() {
  const { success, info, error: showError } = useToast();
  const [groups, setGroups] = useState<GroupWithMembers[]>([]);
  const [loading, setLoading] = useState(true);
  const [expanded, setExpanded] = useState<string | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [newGroup, setNewGroup] = useState({ name: "", description: "" });
  const [confirmDelete, setConfirmDelete] = useState<string | null>(null);
  const [showAddMember, setShowAddMember] = useState<string | null>(null);
  const [allContacts, setAllContacts] = useState<Contact[]>([]);

  useEffect(() => {
    fetchGroups();
  }, []);

  async function fetchGroups() {
    setLoading(true);
    try {
      const res = await groupsApi.list();
      const groupsWithEmptyMembers = (res.groups || []).map(g => ({ ...g, members: [] }));
      setGroups(groupsWithEmptyMembers);
    } catch (err) {
      showError("Gagal mengambil grup", "Silakan coba lagi nanti.");
    } finally {
      setLoading(false);
    }
  }

  async function fetchMembers(groupId: string) {
    try {
      const res = await groupsApi.listMembers(groupId);
      setGroups(prev => prev.map(g => g.id === groupId ? { ...g, members: res.members || [] } : g));
    } catch (err) {
      showError("Gagal mengambil anggota", "Sistem tidak dapat memuat anggota grup.");
    }
  }

  async function toggleExpand(id: string) {
    if (expanded !== id) {
      await fetchMembers(id);
      setExpanded(id);
    } else {
      setExpanded(null);
    }
  }

  async function createGroup() {
    if (!newGroup.name) return;
    try {
      await groupsApi.create(newGroup);
      setNewGroup({ name: "", description: "" });
      setShowCreate(false);
      success("Grup Dibuat!", `Grup "${newGroup.name}" berhasil ditambahkan.`);
      fetchGroups();
    } catch (err) {
      showError("Gagal Membuat Grup", "Terjadi kesalahan di server.");
    }
  }

  async function deleteGroup(id: string) {
    try {
      await groupsApi.delete(id);
      if (expanded === id) setExpanded(null);
      setConfirmDelete(null);
      info("Grup Dihapus", "Grup kontak berhasil dihapus.");
      fetchGroups();
    } catch (err) {
      showError("Gagal Menghapus", "Grup tidak dapat dihapus saat ini.");
    }
  }

  async function openAddMember(groupId: string) {
    setShowAddMember(groupId);
    try {
      const res = await contactsApi.list();
      setAllContacts(res.contacts || []);
    } catch (err) {
      showError("Gagal mengambil kontak", "Tidak bisa memuat daftar kontak.");
    }
  }

  async function addMemberToGroup(contactId: string) {
    if (!showAddMember) return;
    try {
      await groupsApi.addMember(showAddMember, contactId);
      success("Anggota Ditambahkan", "Kontak berhasil dimasukkan ke dalam grup.");
      fetchMembers(showAddMember);
      setShowAddMember(null);
    } catch (err) {
      showError("Gagal Menambah", "Mungkin kontak sudah ada di dalam grup.");
    }
  }

  async function removeMember(groupId: string, contactId: string) {
    try {
      await groupsApi.removeMember(groupId, contactId);
      fetchMembers(groupId);
      info("Anggota Dihapus", "Kontak telah dihapus dari grup.");
    } catch (err) {
      showError("Gagal Menghapus", "Terjadi kesalahan saat menghapus anggota.");
    }
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
            <p className="text-xs text-gray-500">Total Member Terdaftar</p>
            <p className="mt-1 text-2xl font-bold text-gray-900">{groups.reduce((a, g) => a + g.members.length, 0)}</p>
          </div>
        </div>

        {/* Group list */}
        {loading ? (
           <div className="py-12 text-center">
             <Loader2 className="mx-auto h-8 w-8 animate-spin text-green-600" />
             <p className="mt-2 text-gray-500 text-sm">Memuat daftar grup...</p>
           </div>
        ) : groups.map((g) => (
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
                  <p className="text-xs text-gray-400">{g.description || "Tanpa deskripsi"}</p>
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
                  <p className="py-4 text-center text-sm text-gray-400">Belum ada anggota di grup ini</p>
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
                          <td className="py-2 font-mono text-gray-500">{m.phone_number}</td>
                          <td className="py-2 text-right">
                            <button onClick={(e) => { e.stopPropagation(); removeMember(g.id, m.id); }} className="rounded p-1 text-gray-400 hover:bg-red-50 hover:text-red-500">
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
                  <button onClick={(e) => { e.stopPropagation(); openAddMember(g.id); }} className="flex items-center gap-1.5 rounded-lg border border-dashed border-gray-300 px-3 py-1.5 text-xs text-gray-500 hover:border-green-400 hover:text-green-600">
                    <UserPlus className="h-3.5 w-3.5" /> Tambah Anggota
                  </button>
                  <button onClick={(e) => { e.stopPropagation(); setConfirmDelete(g.id); }} className="ml-auto flex items-center gap-1.5 rounded-lg border border-dashed border-red-200 px-3 py-1.5 text-xs text-red-400 hover:border-red-400 hover:text-red-600">
                    <Trash2 className="h-3.5 w-3.5" /> Hapus Grup
                  </button>
                </div>
              </div>
            )}
          </div>
        ))}

        {!loading && groups.length === 0 && (
          <div className="rounded-xl border border-dashed border-gray-300 bg-white p-12 text-center">
            <Users className="mx-auto h-10 w-10 text-gray-300" />
            <p className="mt-3 text-gray-500">Belum ada grup kontak</p>
          </div>
        )}
      </div>

      <ConfirmDialog
        open={confirmDelete !== null}
        title="Hapus Grup?"
        description="Grup ini akan dihapus secara permanen. Anggota (kontak) di dalamnya tidak akan terhapus dari phonebook."
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

      {/* Add Member modal */}
      {showAddMember && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/30 backdrop-blur-sm">
          <div className="w-full max-w-md rounded-2xl bg-white p-6 shadow-xl max-h-[80vh] flex flex-col">
            <div className="mb-5 flex items-center justify-between">
              <h2 className="text-lg font-bold text-gray-900">Pilih Kontak</h2>
              <button onClick={() => setShowAddMember(null)}><X className="h-5 w-5 text-gray-400" /></button>
            </div>
            <div className="overflow-y-auto pr-2">
              <div className="space-y-2">
                {allContacts.length === 0 ? (
                  <p className="text-center py-4 text-gray-500 text-sm">Tidak ada kontak tersedia</p>
                ) : allContacts.map(c => (
                  <div key={c.id} className="flex items-center justify-between p-3 rounded-lg border border-gray-100 hover:bg-gray-50">
                    <div>
                      <p className="text-sm font-medium text-gray-900">{c.name}</p>
                      <p className="text-xs text-gray-500 font-mono">{c.phone_number}</p>
                    </div>
                    <button 
                      onClick={() => addMemberToGroup(c.id)}
                      className="text-xs font-medium text-green-600 hover:underline"
                    >
                      Pilih
                    </button>
                  </div>
                ))}
              </div>
            </div>
            <div className="mt-5 border-t pt-4">
              <button onClick={() => setShowAddMember(null)} className="w-full rounded-lg border border-gray-200 py-2 text-sm text-gray-600 hover:bg-gray-50">Tutup</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
