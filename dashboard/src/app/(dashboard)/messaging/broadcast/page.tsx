"use client";

import { useState, useEffect } from "react";
import { Send, Upload, Plus, Trash2, ChevronDown, Info, Users, Phone, CheckSquare, Square, CheckCheck, Smile, Loader2, Image as ImageIcon, X, Search, User } from "lucide-react";
import { useToast } from "@/components/ui/toast";
import { sessionsApi, groupsApi, broadcastApi, messagesApi, contactsApi } from "@/lib/api";
import { Device, ContactGroup, Contact, BroadcastStatus } from "@/lib/types";

function WABubble({ message, mediaPreview }: { message: string, mediaPreview?: string }) {
  const now = new Date();
  const timeStr = now.getHours().toString().padStart(2, "0") + ":" + now.getMinutes().toString().padStart(2, "0");

  return (
    <div className="flex h-full flex-col rounded-xl border border-gray-200 bg-white shadow-sm overflow-hidden">
      <div className="flex items-center gap-3 bg-[#075E54] px-4 py-3">
        <div className="flex h-9 w-9 items-center justify-center rounded-full bg-white/20 text-sm font-bold text-white">B</div>
        <div>
          <p className="text-sm font-semibold text-white">Penerima</p>
          <p className="text-xs text-green-200">broadcast</p>
        </div>
      </div>
      <div
        className="flex-1 p-4 overflow-y-auto"
        style={{
          backgroundImage: `url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23128c7e' fill-opacity='0.04'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")`,
          backgroundColor: "#ECE5DD",
        }}
      >
        {(message.trim() || mediaPreview) ? (
          <div className="flex justify-end">
            <div className="max-w-[85%] rounded-lg rounded-br-sm bg-[#DCF8C6] p-1.5 shadow-sm">
              {mediaPreview && (
                <div className="mb-1 overflow-hidden rounded-md bg-black/5">
                  <img src={mediaPreview} alt="Preview" className="w-full object-contain" />
                </div>
              )}
              {message.trim() && (
                <p className="whitespace-pre-wrap break-words px-1.5 py-1 text-sm text-gray-800">{message}</p>
              )}
              <div className="mt-1 flex items-center justify-end gap-1 px-1">
                <span className="text-[10px] text-gray-400">{timeStr}</span>
                <CheckCheck className="h-3 w-3 text-[#53BDEB]" />
              </div>
            </div>
          </div>
        ) : (
          <div className="flex h-full items-center justify-center">
            <div className="text-center">
              <Smile className="mx-auto h-8 w-8 text-gray-300" />
              <p className="mt-2 text-xs text-gray-400">Preview pesan broadcast</p>
            </div>
          </div>
        )}
      </div>
      <div className="flex items-center gap-2 bg-[#F0F2F5] px-3 py-2">
        <div className="flex-1 rounded-full bg-white px-4 py-2 text-xs text-gray-400 shadow-sm">Ketik pesan...</div>
        <div className="flex h-9 w-9 items-center justify-center rounded-full bg-[#128C7E]">
          <Send className="h-3.5 w-3.5 text-white" />
        </div>
      </div>
    </div>
  );
}

export default function BroadcastPage() {
  const { success, error: errorToast } = useToast();
  
  // State for data
  const [devices, setDevices] = useState<Device[]>([]);
  const [groups, setGroups] = useState<ContactGroup[]>([]);
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [groupMembers, setGroupMembers] = useState<Record<string, string[]>>({}); // groupId -> recipientPhones
  
  // State for form
  const [name, setName] = useState("");
  const [recipientMode, setRecipientMode] = useState<"manual" | "phonebook" | "group">("manual");
  const [recipients, setRecipients] = useState<string[]>([""]);
  const [selectedContacts, setSelectedContacts] = useState<string[]>([]); // Array of phone numbers
  const [selectedGroups, setSelectedGroups] = useState<string[]>([]);
  const [contactSearch, setContactSearch] = useState("");
  
  const [message, setMessage] = useState("");
  const [selectedDevice, setSelectedDevice] = useState("");
  const [delay, setDelay] = useState("5");
  const [mediaFile, setMediaFile] = useState<File | null>(null);
  const [mediaPreview, setMediaPreview] = useState<string | undefined>();
  
  const [loading, setLoading] = useState(true);
  const [sending, setSending] = useState(false);

  useEffect(() => {
    fetchInitialData();
  }, []);

  const fetchInitialData = async () => {
    try {
      const [sessionsRes, groupsRes, contactsRes] = await Promise.all([
        sessionsApi.list(),
        groupsApi.list(),
        contactsApi.list()
      ]);
      
      const activeDevices = sessionsRes.sessions.filter(s => s.status === 1);
      setDevices(activeDevices);
      if (activeDevices.length > 0) setSelectedDevice(activeDevices[0].device_id);
      
      setGroups(groupsRes.groups);
      setContacts(contactsRes.contacts);
      
      // Fetch members for each group to get phone numbers
      const memberPromises = groupsRes.groups.map(async (g) => {
        try {
          const mems = await groupsApi.listMembers(g.id);
          return { id: g.id, phones: mems.members.map(m => m.phone_number) };
        } catch (e) {
          return { id: g.id, phones: [] };
        }
      });
      
      const results = await Promise.all(memberPromises);
      const memberMap: Record<string, string[]> = {};
      results.forEach(r => memberMap[r.id] = r.phones);
      setGroupMembers(memberMap);
      
    } catch (err) {
      errorToast("Gagal memuat data", "Pastikan backend berjalan.");
    } finally {
      setLoading(false);
    }
  };

  const handleMediaChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setMediaFile(file);
      const reader = new FileReader();
      reader.onloadend = () => setMediaPreview(reader.result as string);
      reader.readAsDataURL(file);
    }
  };

  const removeMedia = () => {
    setMediaFile(null);
    setMediaPreview(undefined);
  };

  const addRecipient = () => setRecipients((r) => [...r, ""]);
  const removeRecipient = (i: number) => setRecipients((r) => r.filter((_, idx) => idx !== i));
  const updateRecipient = (i: number, val: string) =>
    setRecipients((r) => r.map((v, idx) => (idx === i ? val : v)));
  
  const validManualCount = recipients.filter((r) => /^[0-9]{9,15}$/.test(r)).length;
  
  const toggleContact = (phone: string) =>
    setSelectedContacts((prev) =>
      prev.includes(phone) ? prev.filter((p) => p !== phone) : [...prev, phone]
    );

  const toggleGroup = (id: string) =>
    setSelectedGroups((prev) =>
      prev.includes(id) ? prev.filter((g) => g !== id) : [...prev, id]
    );

  const getGroupRecipientPhones = () => {
    const phones = new Set<string>();
    selectedGroups.forEach(gid => {
      groupMembers[gid]?.forEach(p => phones.add(p));
    });
    return Array.from(phones);
  };

  const totalCount = 
    recipientMode === "manual" ? validManualCount : 
    recipientMode === "phonebook" ? selectedContacts.length : 
    getGroupRecipientPhones().length;

  const canSend = name.trim().length > 0 && totalCount > 0 && selectedDevice;

  async function handleSend() {
    if (!canSend || sending) return;
    setSending(true);
    
    try {
      let finalRecipients: string[] = [];
      if (recipientMode === "manual") {
        finalRecipients = recipients.filter(r => /^[0-9]{9,15}$/.test(r));
      } else if (recipientMode === "phonebook") {
        finalRecipients = selectedContacts;
      } else {
        finalRecipients = getGroupRecipientPhones();
      }

      let mediaUrl: string | undefined;

      // 1. Upload media if exists
      if (mediaFile) {
        const formData = new FormData();
        formData.append("file", mediaFile);
        const uploadRes = await messagesApi.uploadMedia(formData);
        mediaUrl = uploadRes.url;
      }

      // 2. Create Campaign
      const campaign = await broadcastApi.create({
        device_id: selectedDevice,
        name: name,
        message_content: message,
        delay_seconds: parseInt(delay),
        recipients: finalRecipients,
      });

      // 3. Start Campaign
      await broadcastApi.start(campaign.id);

      success("Broadcast Dimulai!", `Campaign ${name} telah didaftarkan dan mulai dikirim ke ${totalCount} penerima.`);
      
      // Reset form
      setName("");
      setMessage("");
      setRecipients([""]);
      setSelectedContacts([]);
      setSelectedGroups([]);
      removeMedia();
      
    } catch (err: any) {
      errorToast("Gagal mengirim broadcast", err.response?.data?.error || err.message);
    } finally {
      setSending(false);
    }
  }

  const filteredContacts = contacts.filter(c => 
    c.name?.toLowerCase().includes(contactSearch.toLowerCase()) || 
    c.phone_number.includes(contactSearch)
  );

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50">
        <div className="text-center">
          <Loader2 className="mx-auto h-8 w-8 animate-spin text-green-500" />
          <p className="mt-4 text-sm text-gray-500 font-medium">Memuat data...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 pb-12">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold text-gray-900">Broadcast / Bulk Sender</h1>
            <p className="text-sm text-gray-500">Kirim pesan ke banyak nomor sekaligus</p>
          </div>
          <div className="hidden sm:flex items-center gap-3">
            <div className={`flex items-center gap-2 rounded-full px-3 py-1 text-xs font-semibold ${devices.length > 0 ? "bg-green-100 text-green-700" : "bg-red-100 text-red-700"}`}>
              <div className={`h-2 w-2 rounded-full ${devices.length > 0 ? "bg-green-500" : "bg-red-500"}`} />
              {devices.length} Device Aktif
            </div>
          </div>
        </div>
      </div>

      <div className="p-6">
        <div className="mx-auto max-w-7xl">
          <div className="grid grid-cols-1 gap-8 lg:grid-cols-12">
            
            {/* Form Column */}
            <div className="space-y-6 lg:col-span-8">
              
              {/* Campaign Name & Device */}
              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
                  <label className="mb-1.5 block text-sm font-semibold text-gray-700">Nama Campaign</label>
                  <input
                    type="text"
                    placeholder="Contoh: Promo Ramadhan 2024"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    className="h-10 w-full rounded-lg border border-gray-300 px-3 text-sm focus:border-green-500 focus:outline-none focus:ring-1 focus:ring-green-500"
                  />
                </div>
                <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
                  <label className="mb-1.5 block text-sm font-semibold text-gray-700">Pilih Device Pengirim</label>
                  <div className="relative">
                    <select
                      value={selectedDevice}
                      onChange={(e) => setSelectedDevice(e.target.value)}
                      className="h-10 w-full appearance-none rounded-lg border border-gray-300 px-3 pr-8 text-sm focus:border-green-500 focus:outline-none focus:ring-1 focus:ring-green-500"
                    >
                      {devices.length === 0 && <option value="">Tidak ada device aktif</option>}
                      {devices.map((d) => (
                        <option key={d.device_id} value={d.device_id}>
                          {d.display_name || d.phone || d.device_id}
                        </option>
                      ))}
                    </select>
                    <ChevronDown className="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
                  </div>
                </div>
              </div>

              {/* Recipients */}
              <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
                <div className="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                  <h2 className="text-sm font-bold uppercase tracking-wider text-gray-400">Daftar Penerima</h2>
                  <div className="flex gap-1 rounded-lg bg-gray-100 p-0.5">
                    <button
                      onClick={() => setRecipientMode("manual")}
                      className={`px-4 py-1.5 text-xs font-semibold rounded-md transition-all ${recipientMode === "manual" ? "bg-white text-gray-900 shadow-sm" : "text-gray-500 hover:text-gray-700"}`}
                    >
                      Manual Phone
                    </button>
                    <button
                      onClick={() => setRecipientMode("phonebook")}
                      className={`px-4 py-1.5 text-xs font-semibold rounded-md transition-all ${recipientMode === "phonebook" ? "bg-white text-gray-900 shadow-sm" : "text-gray-500 hover:text-gray-700"}`}
                    >
                      Phone Book
                    </button>
                    <button
                      onClick={() => setRecipientMode("group")}
                      className={`px-4 py-1.5 text-xs font-semibold rounded-md transition-all ${recipientMode === "group" ? "bg-white text-gray-900 shadow-sm" : "text-gray-500 hover:text-gray-700"}`}
                    >
                      Group Contact
                    </button>
                  </div>
                </div>

                {recipientMode === "manual" ? (
                  <div className="space-y-4">
                    <div className="flex items-center justify-between text-sm">
                      <span className="text-gray-600">Masukkan nomor WhatsApp (diawali 62)</span>
                      <span className="font-semibold text-green-600">{validManualCount} Valid</span>
                    </div>
                    <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 md:grid-cols-3 max-h-[300px] overflow-y-auto pr-2 pb-2">
                      {recipients.map((r, i) => (
                        <div key={i} className="group relative">
                          <input
                            type="tel"
                            placeholder="628..."
                            value={r}
                            onChange={(e) => updateRecipient(i, e.target.value.replace(/\D/g, ""))}
                            className="h-9 w-full rounded-lg border border-gray-300 px-3 pl-8 text-xs font-mono focus:border-green-500 focus:outline-none focus:ring-1 focus:ring-green-500"
                          />
                          <Phone className="absolute left-2.5 top-2.5 h-3.5 w-3.5 text-gray-400" />
                          {recipients.length > 1 && (
                            <button 
                              onClick={() => removeRecipient(i)} 
                              className="absolute right-1 top-1 rounded-md p-1.5 text-gray-300 opacity-0 group-hover:opacity-100 hover:bg-red-50 hover:text-red-500 transition-all"
                            >
                              <Trash2 className="h-3 w-3" />
                            </button>
                          )}
                        </div>
                      ))}
                    </div>
                    <button onClick={addRecipient} className="flex h-9 items-center justify-center gap-2 rounded-lg border-2 border-dashed border-gray-200 px-4 text-xs font-semibold text-gray-500 hover:border-green-300 hover:text-green-600 transition-all w-full">
                      <Plus className="h-4 w-4" /> Tambah Penerima
                    </button>
                  </div>
                ) : recipientMode === "phonebook" ? (
                  <div className="space-y-4">
                    <div className="relative">
                      <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
                      <input 
                        type="text" 
                        placeholder="Cari nama atau nomor..." 
                        value={contactSearch}
                        onChange={(e) => setContactSearch(e.target.value)}
                        className="h-10 w-full rounded-lg border border-gray-300 pl-10 pr-4 text-sm focus:border-green-500 focus:outline-none focus:ring-1 focus:ring-green-500"
                      />
                    </div>
                    <div className="grid grid-cols-1 gap-2 sm:grid-cols-2 max-h-[300px] overflow-y-auto pr-2 pb-2">
                      {filteredContacts.length === 0 && (
                        <p className="col-span-full py-8 text-center text-xs text-gray-400 italic">Kontak tidak ditemukan.</p>
                      )}
                      {filteredContacts.map((c) => {
                        const isSelected = selectedContacts.includes(c.phone_number);
                        return (
                          <div 
                            key={c.id} 
                            onClick={() => toggleContact(c.phone_number)}
                            className={`flex cursor-pointer items-center justify-between rounded-xl border p-3 transition-all ${isSelected ? "border-green-500 bg-green-50" : "border-gray-100 bg-white hover:border-gray-200"}`}
                          >
                            <div className="flex items-center gap-3">
                              <div className={`flex h-8 w-8 items-center justify-center rounded-full ${isSelected ? "bg-green-500 text-white" : "bg-gray-100 text-gray-400"}`}>
                                <User className="h-4 w-4" />
                              </div>
                              <div className="min-w-0">
                                <p className="truncate text-xs font-bold text-gray-900">{c.name || "Tanpa Nama"}</p>
                                <p className="text-[10px] text-gray-500">{c.phone_number}</p>
                              </div>
                            </div>
                            {isSelected ? <CheckSquare className="h-4 w-4 text-green-500" /> : <Square className="h-4 w-4 text-gray-200" />}
                          </div>
                        );
                      })}
                    </div>
                    <div className="flex items-center justify-between rounded-lg bg-gray-50 px-4 py-2 text-xs font-medium text-gray-500">
                      <span>{selectedContacts.length} Kontak Terpilih</span>
                      {selectedContacts.length > 0 && (
                        <button onClick={() => setSelectedContacts([])} className="text-red-500 hover:underline">Hapus Semua</button>
                      )}
                    </div>
                  </div>
                ) : (
                  <div className="space-y-4">
                    <div className="text-sm text-gray-600 mb-2">Pilih satu atau lebih group contact:</div>
                    <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 max-h-[300px] overflow-y-auto pr-2 pb-2">
                      {groups.length === 0 && <p className="text-xs text-gray-400 italic">Belum ada group kontak.</p>}
                      {groups.map((g) => {
                        const isSelected = selectedGroups.includes(g.id);
                        return (
                          <div 
                            key={g.id} 
                            onClick={() => toggleGroup(g.id)}
                            className={`flex cursor-pointer items-center justify-between rounded-xl border p-4 transition-all ${isSelected ? "border-green-500 bg-green-50 ring-1 ring-green-500" : "border-gray-200 bg-white hover:border-gray-300"}`}
                          >
                            <div className="flex items-center gap-3">
                              <div className={`flex h-8 w-8 items-center justify-center rounded-lg ${isSelected ? "bg-green-500 text-white" : "bg-gray-100 text-gray-500"}`}>
                                <Users className="h-4 w-4" />
                              </div>
                              <div>
                                <p className="text-sm font-bold text-gray-900">{g.name}</p>
                                <p className="text-[10px] text-gray-500">{groupMembers[g.id]?.length || 0} Kontak</p>
                              </div>
                            </div>
                            {isSelected ? <CheckSquare className="h-5 w-5 text-green-500" /> : <Square className="h-5 w-5 text-gray-300" />}
                          </div>
                        );
                      })}
                    </div>
                  </div>
                )}
              </div>

              {/* Message Content */}
              <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
                <h2 className="mb-4 text-sm font-bold uppercase tracking-wider text-gray-400">Isi Pesan</h2>
                <div className="space-y-4">
                  <textarea
                    rows={6}
                    placeholder="Halo! Ada promo menarik hari ini..."
                    value={message}
                    onChange={(e) => setMessage(e.target.value)}
                    className="w-full resize-none rounded-xl border border-gray-300 p-4 text-sm focus:border-green-500 focus:outline-none focus:ring-1 focus:ring-green-500"
                  />
                  
                  <div className="flex flex-wrap items-center justify-between gap-4">
                    <div className="flex gap-2">
                      <label className="flex cursor-pointer items-center gap-2 rounded-lg border border-gray-300 px-4 py-2 hover:bg-gray-50 transition-colors">
                        <ImageIcon className="h-4 w-4 text-gray-500" />
                        <span className="text-xs font-semibold text-gray-700">Lampirkan Media</span>
                        <input type="file" className="hidden" accept="image/*,video/*" onChange={handleMediaChange} />
                      </label>
                      {mediaPreview && (
                        <div className="relative h-10 w-10 overflow-hidden rounded-md border border-gray-300">
                          <img src={mediaPreview} alt="Media" className="h-full w-full object-cover" />
                          <button onClick={removeMedia} className="absolute -right-1 -top-1 rounded-full bg-red-500 p-0.5 text-white shadow-sm">
                            <X className="h-3 w-3" />
                          </button>
                        </div>
                      )}
                    </div>
                    <div className="flex items-center gap-3">
                      <span className="text-xs text-gray-400">Delay:</span>
                      <div className="flex items-center gap-1">
                        <input 
                          type="number"
                          value={delay}
                          onChange={(e) => setDelay(e.target.value)}
                          className="w-16 rounded-md border border-gray-300 px-2 py-1 text-center text-xs focus:outline-none focus:ring-1 focus:ring-green-500" 
                        />
                        <span className="text-[10px] font-medium text-gray-500">detik / pesan</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              {/* Action */}
              <button
                onClick={handleSend}
                disabled={!canSend || sending}
                className={`flex h-12 w-full items-center justify-center gap-3 rounded-xl text-sm font-bold text-white shadow-lg transition-all ${canSend && !sending ? "bg-green-600 hover:bg-green-700 hover:scale-[1.01] active:scale-95 shadow-green-200" : "bg-gray-300 cursor-not-allowed"}`}
              >
                {sending ? (
                  <>
                    <Loader2 className="h-5 w-5 animate-spin" />
                    Memproses Campaign...
                  </>
                ) : (
                  <>
                    <Send className="h-4 w-4" />
                    Kirim ke {totalCount} Penerima Sekarang
                  </>
                )}
              </button>
            </div>

            {/* Preview Column */}
            <div className="lg:col-span-4">
              <div className="sticky top-6">
                <h2 className="mb-4 text-sm font-bold uppercase tracking-wider text-gray-400">WhatsApp Preview</h2>
                <div className="h-[550px]">
                  <WABubble message={message} mediaPreview={mediaPreview} />
                </div>
                <div className="mt-6 rounded-xl bg-blue-50 p-4 border border-blue-100">
                  <div className="flex gap-3">
                    <Info className="h-5 w-5 flex-shrink-0 text-blue-500" />
                    <div>
                      <p className="text-xs font-bold text-blue-900">Tips Anti-Blokir</p>
                      <p className="mt-1 text-[10px] text-blue-700 leading-relaxed">
                        WhatsApp memantau aktivitas pengiriman cepat. Gunakan delay minimal 5-10 detik dan hindari mengirimkan link yang terlalu spammy ke nomor yang belum menyimpan kontak Anda.
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>

          </div>
        </div>
      </div>
    </div>
  );
}
