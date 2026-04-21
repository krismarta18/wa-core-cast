"use client";

import { useState, useEffect, useRef } from "react";
import { Send, Paperclip, ChevronDown, CheckCheck, Smile, X, File, FileText, Image as ImageIcon, Video, Loader2, Upload, Download, Info, BookText, Users, Search } from "lucide-react";
import { toast } from "sonner";
import { sessionsApi, messagesApi, autoResponseApi, contactsApi } from "@/lib/api";
import type { Device, MessageTemplate, Contact } from "@/lib/types";

function WABubble({ from, message, mediaUrl, mediaType }: { from: string; message: string; mediaUrl?: string | null; mediaType?: string }) {
  const now = new Date();
  const timeStr = now.getHours().toString().padStart(2, "0") + ":" + now.getMinutes().toString().padStart(2, "0");

  return (
    <div className="flex h-full flex-col overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
      {/* WA chat header */}
      <div className="flex items-center gap-3 bg-[#075E54] px-4 py-3">
        <div className="flex h-9 w-9 items-center justify-center rounded-full bg-white/20 text-sm font-bold text-white">
          {from ? from.charAt(0).toUpperCase() : "N"}
        </div>
        <div>
          <p className="text-sm font-semibold text-white">{from || "Nomor Tujuan"}</p>
          <p className="text-xs text-green-200">online</p>
        </div>
      </div>

      {/* Chat area */}
      <div
        className="flex-1 p-4"
        style={{
          backgroundImage: `url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23128c7e' fill-opacity='0.04'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")`,
          backgroundColor: "#ECE5DD",
        }}
      >
        {message.trim() || mediaUrl ? (
          <div className="flex justify-end">
            <div className="max-w-[85%] rounded-lg rounded-br-sm bg-[#DCF8C6] p-1 shadow-sm">
              {mediaUrl && (
                <div className="mb-1 flex items-center justify-center overflow-hidden rounded-sm bg-black/5 p-1">
                  {mediaType === "image" ? (
                    <img src={mediaUrl} className="max-h-48 max-w-full rounded-sm object-cover" alt="Preview" />
                  ) : mediaType === "video" ? (
                    <video src={mediaUrl} className="max-h-48 max-w-full rounded-sm object-cover" controls />
                  ) : (
                    <div className="flex w-full min-w-32 items-center gap-2 rounded-sm bg-black/10 p-3">
                      <FileText className="h-6 w-6 text-gray-600" />
                      <span className="truncate text-sm font-medium text-gray-700">Document</span>
                    </div>
                  )}
                </div>
              )}
              {message.trim() && (
                <div className="px-2 pb-1 pt-1">
                  <p className="whitespace-pre-wrap break-words text-sm text-gray-800">{message}</p>
                </div>
              )}
              <div className="mt-1 flex items-center justify-end gap-1 px-2 pb-1">
                <span className="text-[10px] text-gray-400">{timeStr}</span>
                <CheckCheck className="h-3 w-3 text-[#53BDEB]" />
              </div>
            </div>
          </div>
        ) : (
          <div className="flex h-full items-center justify-center">
            <div className="text-center">
              <Smile className="mx-auto h-8 w-8 text-gray-300" />
              <p className="mt-2 text-xs text-gray-400">Preview pesan akan muncul di sini</p>
            </div>
          </div>
        )}
      </div>

      {/* WA input bar */}
      <div className="flex items-center gap-2 bg-[#F0F2F5] px-3 py-2">
        <div className="flex-1 rounded-full bg-white px-4 py-2 text-xs text-gray-400 shadow-sm">
          Ketik pesan...
        </div>
        <div className="flex h-9 w-9 items-center justify-center rounded-full bg-[#128C7E]">
          <Send className="h-3.5 w-3.5 text-white" />
        </div>
      </div>
    </div>
  );
}

export default function NewMessagePage() {
  const [activeTab, setActiveTab] = useState<"single" | "bulk">("single");
  const [devices, setDevices] = useState<Device[]>([]);
  const [device, setDevice] = useState("");
  const [to, setTo] = useState("");
  const [message, setMessage] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isLoadingDevices, setIsLoadingDevices] = useState(true);
  const [templates, setTemplates] = useState<MessageTemplate[]>([]);
  const [showTemplates, setShowTemplates] = useState(false);
  
  // Contact picker states
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [showContacts, setShowContacts] = useState(false);
  const [contactSearch, setContactSearch] = useState("");
  const [isLoadingContacts, setIsLoadingContacts] = useState(false);

  // File states for Single Send
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [attachedFile, setAttachedFile] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);

  // File states for Bulk Send
  const csvInputRef = useRef<HTMLInputElement>(null);
  const [csvFile, setCsvFile] = useState<File | null>(null);
  const [csvPreview, setCsvPreview] = useState<{phone: string, content: string}[]>([]);
  const [csvHeaders, setCsvHeaders] = useState<string[]>([]);

  useEffect(() => {
    sessionsApi.list()
      .then((res) => {
        const activeDevices = (res.sessions || []).filter((d: Device) => d.status === 1 || d.is_active);
        setDevices(activeDevices);
        if (activeDevices.length > 0) {
          setDevice(activeDevices[0].device_id);
        }
      })
      .catch((err) => {
        toast.error("Gagal memuat daftar perangkat aktif.");
        console.error(err);
      })
      .finally(() => {
        setIsLoadingDevices(false);
      });

    // Fetch templates
    autoResponseApi.getTemplates().then(res => setTemplates(res.templates || []));
    
    // Fetch contacts
    setIsLoadingContacts(true);
    contactsApi.list()
      .then(res => setContacts(res.contacts || []))
      .finally(() => setIsLoadingContacts(false));
  }, []);

  useEffect(() => {
    if (attachedFile) {
      const url = URL.createObjectURL(attachedFile);
      setPreviewUrl(url);
      return () => URL.revokeObjectURL(url);
    } else {
      setPreviewUrl(null);
    }
  }, [attachedFile]);

  const isValidSingle = /^[0-9]{9,15}$/.test(to) && (message.trim().length > 0 || attachedFile !== null) && device !== "";
  const isValidBulk = csvFile !== null && device !== "";

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      if (file.size > 50 * 1024 * 1024) { // 50MB
        toast.error("File terlalu besar (maksimal 50MB)");
        return;
      }
      setAttachedFile(file);
    }
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  const handleCsvChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      if (!file.name.endsWith(".csv")) {
        toast.error("Hanya file .csv yang didukung");
        return;
      }
      setCsvFile(file);
      
      // Generate Preview
      const reader = new FileReader();
      reader.onload = (event) => {
        const text = event.target?.result as string;
        const lines = text.split("\n").map(l => l.trim()).filter(l => l !== "");
        if (lines.length < 2) return;
        
        const headers = lines[0].split(",").map(h => h.trim().replace(/^"|"$/g, ''));
        setCsvHeaders(headers);
        
        const phoneIdx = headers.findIndex(h => ["phone", "nomor", "no_hp"].includes(h.toLowerCase()));
        const msgIdx = headers.findIndex(h => ["message", "pesan", "content"].includes(h.toLowerCase()));
        
        const previewRows: {phone: string, content: string}[] = [];
        const maxPreview = Math.min(6, lines.length); // 1 header + 5 rows
        
        for (let i = 1; i < maxPreview; i++) {
          const row = lines[i].split(",").map(c => c.trim().replace(/^"|"$/g, ''));
          if (row.length < headers.length) continue;
          
          let phone = phoneIdx !== -1 ? row[phoneIdx] : "N/A";
          let content = msgIdx !== -1 ? row[msgIdx] : "";
          
          // Apply variable replacement for preview
          headers.forEach((h, idx) => {
            const placeholder = `[${h}]`;
            const regex = new RegExp(`\\[${h}\\]`, 'gi');
            content = content.replace(regex, row[idx] || "");
          });
          
          previewRows.push({ phone, content });
        }
        setCsvPreview(previewRows);
      };
      reader.readAsText(file);
    }
  };

  const determineContentType = (file: File) => {
    if (file.type.startsWith("image/")) return "image";
    if (file.type.startsWith("video/")) return "video";
    if (file.type.startsWith("audio/")) return "audio";
    return "document";
  };

  const applyContactToMessage = (contact: Contact) => {
    setTo(contact.phone_number);
    
    let newMessage = message;
    // Replace [Nama]
    newMessage = newMessage.replace(/\[Nama\]/gi, contact.name);
    // Replace [Phone]
    newMessage = newMessage.replace(/\[Phone\]/gi, contact.phone_number);
    
    // Replace additional data if exists
    if (contact.additional_data) {
      try {
        const data = typeof contact.additional_data === 'string' ? JSON.parse(contact.additional_data) : contact.additional_data;
        Object.keys(data).forEach(key => {
          const regex = new RegExp(`\\[${key}\\]`, 'gi');
          newMessage = newMessage.replace(regex, data[key]);
        });
      } catch (e) {}
    }
    
    setMessage(newMessage);
    setShowContacts(false);
    toast.success(`Kontak ${contact.name} terpilih`);
  };

  const filteredContacts = contacts.filter(c => 
    c.name.toLowerCase().includes(contactSearch.toLowerCase()) || 
    c.phone_number.includes(contactSearch)
  );

  const handleSendSingle = async () => {
    if (!isValidSingle) return;
    setIsSubmitting(true);

    try {
      if (attachedFile) {
        toast.info("Mengunggah media dan mengirim pesan...");
        const formData = new FormData();
        formData.append("file", attachedFile);
        formData.append("target_jid", to);
        formData.append("content_type", determineContentType(attachedFile));
        if (message.trim().length > 0) {
          formData.append("caption", message.trim());
        }

        await messagesApi.sendMedia(device, formData);
      } else {
        await messagesApi.send(device, {
          target_jid: to,
          content: message,
          priority: 3
        });
      }

      toast.success("Pesan berhasil dimasukkan ke antrian!");
      setTo(""); 
      setMessage(""); 
      setAttachedFile(null);
    } catch (err: any) {
      toast.error(err.response?.data?.error || "Gagal mengirim pesan");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleSendBulk = async () => {
    if (!isValidBulk || !csvFile) return;
    setIsSubmitting(true);

    try {
      toast.info("Memproses file CSV...");
      const formData = new FormData();
      formData.append("file", csvFile);
      
      const res = await messagesApi.bulkSend(device, formData);
      toast.success(`Berhasil! ${res.count} pesan telah dimasukkan ke antrian.`);
      setCsvFile(null);
      if (csvInputRef.current) csvInputRef.current.value = "";
    } catch (err: any) {
      toast.error(err.response?.data?.error || "Gagal mengirim pesan massal");
    } finally {
      setIsSubmitting(false);
    }
  };

  const downloadCsvTemplate = () => {
    const content = "phone,message\n62812xxx,Halo [Name] apa kabar?\n62813xxx,Segera lunasi tagihan Anda.";
    const blob = new Blob([content], { type: "text/csv" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "template_bulk_wacast.csv";
    a.click();
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <h1 className="text-xl font-bold text-gray-900">Kirim Pesan</h1>
        <p className="text-sm text-gray-500">Pilih metode pengiriman pesan WhatsApp Anda</p>
      </div>

      <div className="p-6">
        {/* Tabs */}
        <div className="mb-6 flex gap-1 rounded-lg bg-gray-100 p-1 w-full max-w-md">
          <button
            onClick={() => setActiveTab("single")}
            className={`flex-1 rounded-md px-4 py-2 text-sm font-medium transition-colors ${
              activeTab === "single" ? "bg-white text-green-600 shadow-sm" : "text-gray-500 hover:text-gray-700"
            }`}
          >
            Kirim Tunggal
          </button>
          <button
            onClick={() => setActiveTab("bulk")}
            className={`flex-1 rounded-md px-4 py-2 text-sm font-medium transition-colors ${
              activeTab === "bulk" ? "bg-white text-green-600 shadow-sm" : "text-gray-500 hover:text-gray-700"
            }`}
          >
            Masal (CSV)
          </button>
        </div>

        <div className="grid grid-cols-1 gap-6 lg:grid-cols-5">
          <div className="space-y-4 lg:col-span-3">
            {/* Device Selector (Global for both tabs) */}
            <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
              <label className="mb-1 block text-sm font-medium text-gray-700">Kirim dari Device</label>
              <div className="relative">
                <select
                  value={device}
                  onChange={(e) => setDevice(e.target.value)}
                  disabled={isLoadingDevices || devices.length === 0}
                  className="h-10 w-full appearance-none rounded-lg border border-gray-300 px-3 pr-8 text-sm text-gray-700 focus:outline-none focus:ring-2 focus:ring-green-500 disabled:bg-gray-100"
                >
                  {isLoadingDevices && <option>Memuat perangkat...</option>}
                  {!isLoadingDevices && devices.length === 0 && <option value="">(Tidak ada device aktif)</option>}
                  {devices.map((d: any) => (
                    <option key={d.device_id || d.id} value={d.device_id || d.id}>
                      {d.display_name || d.phone || d.device_id}
                    </option>
                  ))}
                </select>
                <ChevronDown className="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
              </div>
            </div>

            {activeTab === "single" ? (
              <>
                <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
                  <div className="mb-1 flex items-center justify-between">
                    <label className="text-sm font-medium text-gray-700">Nomor Tujuan <span className="text-red-500">*</span></label>
                    <button 
                      type="button"
                      onClick={() => setShowContacts(!showContacts)}
                      className="flex items-center gap-1 text-xs font-semibold text-green-600 hover:text-green-700 transition-colors"
                    >
                      <Users className="h-3 w-3" />
                      Cari dari Buku Telepon
                    </button>
                  </div>
                  <div className="relative">
                    <input
                      type="tel"
                      placeholder="628xxxxxxxxx"
                      value={to}
                      onChange={(e) => setTo(e.target.value.replace(/\D/g, ""))}
                      className="h-10 w-full rounded-lg border border-gray-300 px-3 font-mono text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
                    />
                    {showContacts && (
                      <div className="absolute left-0 right-0 top-11 z-20 overflow-hidden rounded-xl border border-gray-200 bg-white shadow-2xl animate-in fade-in slide-in-from-top-2">
                        <div className="flex items-center gap-2 border-b border-gray-100 px-3 py-2 bg-gray-50">
                          <Search className="h-4 w-4 text-gray-400" />
                          <input
                            autoFocus
                            type="text"
                            placeholder="Cari nama atau nomor..."
                            value={contactSearch}
                            onChange={(e) => setContactSearch(e.target.value)}
                            className="bg-transparent text-sm outline-none w-full"
                          />
                          <button onClick={() => setShowContacts(false)}><X className="h-4 w-4 text-gray-400" /></button>
                        </div>
                        <div className="max-h-60 overflow-y-auto p-1">
                          {isLoadingContacts ? (
                            <div className="flex justify-center p-4"><Loader2 className="h-5 w-5 animate-spin text-green-500" /></div>
                          ) : filteredContacts.length === 0 ? (
                            <p className="p-4 text-center text-xs text-gray-400">Kontak tidak ditemukan</p>
                          ) : (
                            filteredContacts.map((c) => (
                              <button
                                key={c.id}
                                onClick={() => applyContactToMessage(c)}
                                className="flex w-full flex-col px-4 py-2 text-left hover:bg-green-50 transition-colors rounded-lg group"
                              >
                                <span className="text-sm font-medium text-gray-700 group-hover:text-green-700">{c.name}</span>
                                <span className="text-xs text-gray-400 font-mono">{c.phone_number}</span>
                              </button>
                            ))
                          )}
                        </div>
                      </div>
                    )}
                  </div>
                </div>

                <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
                  <div className="relative">
                    {attachedFile && (
                      <div className="mb-3 flex items-center justify-between rounded-lg border border-blue-200 bg-blue-50 p-2 pr-3">
                        <div className="flex items-center gap-3 overflow-hidden">
                          <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded bg-blue-100">
                             {attachedFile.type.startsWith("image/") ? <ImageIcon className="h-5 w-5 text-blue-600" /> : <File className="h-5 w-5 text-blue-600" />}
                          </div>
                          <div className="overflow-hidden">
                            <p className="truncate text-sm font-medium text-blue-900">{attachedFile.name}</p>
                            <p className="text-xs text-blue-600">{(attachedFile.size / 1024 / 1024).toFixed(2)} MB</p>
                          </div>
                        </div>
                        <button onClick={() => setAttachedFile(null)} className="rounded-full p-1 text-blue-400 hover:bg-blue-100"><X className="h-4 w-4" /></button>
                      </div>
                    )}
                    <textarea
                      rows={attachedFile ? 3 : 6}
                      placeholder={attachedFile ? "Tulis caption..." : "Tulis pesan..."}
                      value={message}
                      onChange={(e) => setMessage(e.target.value)}
                      className="w-full resize-none rounded-lg border border-gray-300 p-3 pb-12 text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
                    />
                    <div className="absolute bottom-2 left-2 flex items-center gap-1">
                       <button 
                        type="button" 
                        onClick={() => fileInputRef.current?.click()} 
                        className="p-1.5 text-gray-400 hover:text-green-600 transition-colors"
                        title="Lampirkan File"
                       >
                         <Paperclip className="h-4 w-4" />
                       </button>
                       <div className="h-4 w-[1px] bg-gray-200 mx-1" />
                       <button
                        type="button"
                        onClick={() => setShowTemplates(!showTemplates)}
                        className={`flex items-center gap-1.5 rounded-md px-2 py-1 text-xs font-medium transition-colors ${showTemplates ? "bg-green-100 text-green-700" : "text-gray-400 hover:text-green-600"}`}
                       >
                         <BookText className="h-3.5 w-3.5" />
                         Pilih Template
                       </button>
                       <input type="file" ref={fileInputRef} className="hidden" onChange={handleFileChange} />
                    </div>

                    {/* Template Picker Overlay */}
                    {showTemplates && (
                      <div className="absolute bottom-12 left-2 z-10 w-64 max-h-60 overflow-y-auto rounded-xl border border-gray-200 bg-white p-2 shadow-xl animate-in fade-in slide-in-from-bottom-2">
                        <div className="mb-2 px-2 py-1 flex items-center justify-between border-b border-gray-50">
                          <span className="text-[10px] font-bold uppercase tracking-wider text-gray-400">Gudang Pesan</span>
                          <button onClick={() => setShowTemplates(false)}><X className="h-3 w-3 text-gray-300 hover:text-gray-600"/></button>
                        </div>
                        {templates.length === 0 ? (
                          <p className="p-4 text-center text-xs text-gray-400">Belum ada template</p>
                        ) : (
                          <div className="space-y-1">
                            {templates.map((t) => (
                              <button
                                key={t.id}
                                onClick={() => {
                                  setMessage(t.content);
                                  setShowTemplates(false);
                                  toast.success("Template dimuat");
                                }}
                                className="w-full rounded-lg px-3 py-2 text-left hover:bg-green-50 group transition-colors"
                              >
                                <p className="truncate text-sm font-medium text-gray-700 group-hover:text-green-700">{t.name}</p>
                                <p className="truncate text-[10px] text-gray-400">{t.content}</p>
                              </button>
                            ))}
                          </div>
                        )}
                      </div>
                    )}
                  </div>
                </div>

                <button
                  onClick={handleSendSingle}
                  disabled={!isValidSingle || isSubmitting}
                  className="w-full flex justify-center items-center rounded-xl bg-green-600 py-3 text-sm font-semibold text-white shadow-sm hover:bg-green-700 disabled:bg-gray-200"
                >
                  {isSubmitting ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <Send className="mr-2 h-4 w-4" />}
                  Kirim Pesan {attachedFile && "Media"}
                </button>
              </>
            ) : (
              // Bulk Send Tab UI
              <div className="space-y-4">
                <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm text-center">
                  <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-green-50">
                    <Upload className="h-8 w-8 text-green-600" />
                  </div>
                  <h3 className="text-lg font-bold text-gray-900">Unggah File CSV</h3>
                  <p className="mb-6 text-sm text-gray-500">Pastikan file CSV memiliki kolom <strong>phone</strong> dan <strong>message</strong>.</p>
                  
                  <div 
                    onClick={() => csvInputRef.current?.click()}
                    className={`cursor-pointer rounded-xl border-2 border-dashed p-8 transition-colors ${
                      csvFile ? "border-green-400 bg-green-50" : "border-gray-200 hover:border-green-300 hover:bg-gray-50"
                    }`}
                  >
                    {csvFile ? (
                      <div className="flex flex-col items-center">
                        <CheckCheck className="h-8 w-8 text-green-600" />
                        <p className="mt-2 text-sm font-semibold text-green-800">{csvFile.name}</p>
                        <p className="text-xs text-green-600">{(csvFile.size / 1024).toFixed(2)} KB</p>
                        <button 
                          onClick={(e) => { e.stopPropagation(); setCsvFile(null); setCsvPreview([]); }}
                          className="mt-3 text-xs font-bold text-red-500 hover:underline"
                        >
                          Ganti File
                        </button>
                      </div>
                    ) : (
                      <>
                        <p className="text-sm font-medium text-gray-700">Klik di sini untuk memilih file</p>
                        <p className="mt-1 text-xs text-gray-400">Hanya CSV (Maks 10MB)</p>
                      </>
                    )}
                    <input type="file" ref={csvInputRef} className="hidden" accept=".csv" onChange={handleCsvChange} />
                  </div>

                  {/* CSV Preview Table */}
                  {csvFile && csvPreview.length > 0 && (
                    <div className="mt-6 overflow-hidden rounded-xl border border-gray-200">
                      <div className="bg-gray-50 px-4 py-2 border-b border-gray-200 flex items-center justify-between">
                        <span className="text-[10px] font-bold uppercase tracking-wider text-gray-400">Pratinjau Data (5 Baris Pertama)</span>
                        <span className="text-[10px] font-medium text-green-600 bg-green-50 px-1.5 py-0.5 rounded">Format OK</span>
                      </div>
                      <div className="overflow-x-auto">
                        <table className="w-full text-left text-xs">
                          <thead className="bg-gray-50 text-gray-500">
                            <tr>
                              <th className="px-4 py-2 font-semibold">Penerima</th>
                              <th className="px-4 py-2 font-semibold">Pesan Final</th>
                            </tr>
                          </thead>
                          <tbody className="divide-y divide-gray-100">
                            {csvPreview.map((row, i) => (
                              <tr key={i}>
                                <td className="px-4 py-3 font-mono text-gray-600">{row.phone}</td>
                                <td className="px-4 py-3 text-gray-700 max-w-xs truncate">{row.content}</td>
                              </tr>
                            ))}
                          </tbody>
                        </table>
                      </div>
                    </div>
                  )}

                  <div className="mt-6 flex flex-wrap justify-center gap-3">
                    <button
                      onClick={downloadCsvTemplate}
                      className="flex items-center gap-2 rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50"
                    >
                      <Download className="h-4 w-4" /> Unduh Template
                    </button>
                    <button
                      onClick={handleSendBulk}
                      disabled={!isValidBulk || isSubmitting}
                      className="flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700 disabled:opacity-50"
                    >
                      {isSubmitting ? <Loader2 className="h-4 w-4 animate-spin" /> : <Send className="h-4 w-4" />}
                      Kirim Pesan Massal
                    </button>
                  </div>
                </div>

                <div className="rounded-xl border border-blue-200 bg-blue-50 p-5 shadow-sm">
                  <div className="flex gap-3">
                    <Info className="h-5 w-5 text-blue-500 shrink-0" />
                    <div className="text-sm text-blue-700">
                      <p className="font-bold">Tips Personalisasi:</p>
                      <ul className="mt-1 list-inside list-disc space-y-1">
                        <li>Gunakan <code>[Nama]</code> untuk menyapa pelanggan (jika ada kolom "Nama" di CSV).</li>
                        <li>Gunakan <code>[Tagihan]</code> untuk info nominal (jika ada kolom "Tagihan").</li>
                        <li>Dukungan placeholder bersifat fleksibel sesuai dengan header di file CSV Anda!</li>
                      </ul>
                    </div>
                  </div>
                </div>
              </div>
            )}
          </div>

          <div className="lg:col-span-2">
            <p className="mb-2 text-xs font-medium uppercase tracking-wider text-gray-400">Preview (WA Style)</p>
            <div className="h-[480px]">
              <WABubble from={to || "Penerima"} message={activeTab === "bulk" ? "Pesan dari CSV akan diproses otomatis..." : message} mediaUrl={previewUrl} mediaType={attachedFile ? determineContentType(attachedFile) : undefined} />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
