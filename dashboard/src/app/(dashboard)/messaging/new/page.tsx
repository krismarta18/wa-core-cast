"use client";

import { useState, useEffect, useRef } from "react";
import { Send, Paperclip, ChevronDown, CheckCheck, Smile, X, File, FileText, Image as ImageIcon, Video, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { sessionsApi, messagesApi } from "@/lib/api";
import type { Device } from "@/lib/types";

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
  const [devices, setDevices] = useState<Device[]>([]);
  const [device, setDevice] = useState("");
  const [to, setTo] = useState("");
  const [message, setMessage] = useState("");
  const [sent, setSent] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isLoadingDevices, setIsLoadingDevices] = useState(true);

  // File states
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [attachedFile, setAttachedFile] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);

  useEffect(() => {
    sessionsApi.list()
      .then((res) => {
        // sessionsApi.list() returns SessionListResponse which has 'sessions' array
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

  const isValid = /^[0-9]{9,15}$/.test(to) && (message.trim().length > 0 || attachedFile !== null) && device !== "";

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      if (file.size > 50 * 1024 * 1024) { // 50MB
        toast.error("File terlalu besar (maksimal 50MB)");
        return;
      }
      setAttachedFile(file);
    }
    // reset input
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  const determineContentType = (file: File) => {
    if (file.type.startsWith("image/")) return "image";
    if (file.type.startsWith("video/")) return "video";
    if (file.type.startsWith("audio/")) return "audio";
    return "document";
  };

  const handleSend = async () => {
    if (!isValid) return;
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

      setSent(true);
      toast.success("Pesan berhasil dimasukkan ke antrian!");
      
      setTimeout(() => { 
        setSent(false); 
        setTo(""); 
        setMessage(""); 
        setAttachedFile(null);
      }, 2500);
    } catch (err: any) {
      toast.error(err.response?.data?.error || "Gagal mengirim pesan");
      console.error(err);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <h1 className="text-xl font-bold text-gray-900">New Message</h1>
        <p className="text-sm text-gray-500">Kirim pesan WhatsApp ke satu nomor tertentu</p>
      </div>

      <div className="p-6">
        {sent ? (
          <div className="mx-auto max-w-md rounded-xl border border-green-200 bg-green-50 p-10 text-center shadow-sm">
            <Send className="mx-auto h-12 w-12 text-green-500" />
            <h2 className="mt-4 text-lg font-bold text-green-800">Pesan Masuk Antrian!</h2>
            <p className="mt-1 text-sm text-gray-600">Pesan sedang diproses untuk dikirim ke <strong>{to}</strong></p>
          </div>
        ) : (
          <div className="grid grid-cols-1 gap-6 lg:grid-cols-5">
            {/* Form */}
            <div className="space-y-4 lg:col-span-3">
              {/* Device */}
              <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
                <label className="mb-1 block text-sm font-medium text-gray-700">Kirim dari Device</label>
                <div className="relative">
                  <select
                    value={device}
                    onChange={(e) => setDevice(e.target.value)}
                    disabled={isLoadingDevices || devices.length === 0}
                    className="h-10 w-full appearance-none rounded-lg border border-gray-300 px-3 pr-8 text-sm text-gray-700 focus:outline-none focus:ring-2 focus:ring-green-500 disabled:bg-gray-100"
                  >
                    {isLoadingDevices && <option>Loading devices...</option>}
                    {!isLoadingDevices && devices.length === 0 && <option value="">(Tidak ada device aktif)</option>}
                    {devices.map((d: any) => <option key={d.device_id || d.id} value={d.device_id || d.id}>{d.display_name || d.phone || d.device_id || d.id}</option>)}
                  </select>
                  <ChevronDown className="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
                </div>
                {!isLoadingDevices && devices.length === 0 && (
                  <p className="mt-2 text-xs text-red-500">Silakan hubungkan device Anda di menu Devices terlebih dahulu.</p>
                )}
              </div>

              {/* To */}
              <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
                <label className="mb-1 block text-sm font-medium text-gray-700">
                  Nomor Tujuan <span className="text-red-500">*</span>
                </label>
                <input
                  type="tel"
                  inputMode="numeric"
                  placeholder="628xxxxxxxxx"
                  value={to}
                  onChange={(e) => setTo(e.target.value.replace(/\D/g, ""))}
                  className="h-10 w-full rounded-lg border border-gray-300 px-3 font-mono text-sm placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-green-500"
                />
                <p className="mt-1 text-xs text-gray-400">Format: 628xxxxxxxxx (tanpa + atau spasi)</p>
              </div>

              {/* Message */}
              <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
                <div className="flex justify-between items-center mb-1">
                  <label className="block text-sm font-medium text-gray-700">
                    Pesan {attachedFile ? "" : <span className="text-red-500">*</span>}
                  </label>
                  {attachedFile && (
                    <span className="text-xs text-green-600 bg-green-50 px-2 py-0.5 rounded-full font-medium">
                      Mode Lampiran Media Aktif
                    </span>
                  )}
                </div>

                <div className="relative">
                  {attachedFile && (
                    <div className="mb-3 flex items-center justify-between rounded-lg border border-blue-200 bg-blue-50 p-2 pr-3">
                      <div className="flex items-center gap-3 overflow-hidden">
                        <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded bg-blue-100">
                           {attachedFile.type.startsWith("image/") ? <ImageIcon className="h-5 w-5 text-blue-600" /> :
                            attachedFile.type.startsWith("video/") ? <Video className="h-5 w-5 text-blue-600" /> :
                            <File className="h-5 w-5 text-blue-600" />}
                        </div>
                        <div className="overflow-hidden">
                          <p className="truncate text-sm font-medium text-blue-900">{attachedFile.name}</p>
                          <p className="text-xs text-blue-600">{(attachedFile.size / 1024 / 1024).toFixed(2)} MB</p>
                        </div>
                      </div>
                      <button 
                         onClick={() => setAttachedFile(null)}
                         title="Hapus lampiran"
                         className="rounded-full p-1.5 text-blue-400 hover:bg-blue-100 hover:text-blue-600 transition-colors"
                      >
                         <X className="h-4 w-4" />
                      </button>
                    </div>
                  )}

                  <textarea
                    rows={attachedFile ? 3 : 6}
                    placeholder={attachedFile ? "Tulis caption (opsional)..." : "Tulis pesan di sini..."}
                    value={message}
                    onChange={(e) => setMessage(e.target.value)}
                    className="w-full resize-none rounded-lg border border-gray-300 p-3 pb-10 text-sm placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-green-500"
                  />
                  
                  {/* File Upload Button inside text area bounds */}
                  <div className="absolute bottom-2 left-2 flex items-center">
                     <button 
                       onClick={() => fileInputRef.current?.click()}
                       className="inline-flex items-center justify-center rounded-lg p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-700 transition"
                       title="Lampirkan File"
                     >
                       <Paperclip className="h-4 w-4" />
                     </button>
                     <input 
                       type="file" 
                       ref={fileInputRef} 
                       className="hidden" 
                       onChange={handleFileChange}
                       accept="image/*,video/*,audio/*,.pdf,.doc,.docx,.xls,.xlsx" 
                     />
                  </div>
                  
                  <div className="absolute bottom-3 right-3">
                    <span className="text-xs text-gray-400">{message.length} karakter</span>
                  </div>
                </div>
              </div>

              <button
                onClick={handleSend}
                disabled={!isValid || isSubmitting}
                className="w-full flex justify-center items-center rounded-xl bg-green-600 py-3 text-sm font-semibold text-white shadow-sm hover:bg-green-700 disabled:cursor-not-allowed disabled:bg-gray-200 disabled:text-gray-400 transition-colors"
              >
                {isSubmitting ? (
                  <>
                     <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                     Memproses...
                  </>
                ) : (
                  <>
                     <Send className="mr-2 inline h-4 w-4" />
                     Kirim Pesan {attachedFile && "dengan Media"}
                  </>
                )}
              </button>
            </div>

            {/* WA Preview */}
            <div className="lg:col-span-2">
              <p className="mb-2 text-xs font-medium uppercase tracking-wider text-gray-400">Preview Pesan</p>
              <div className="h-[480px]">
                <WABubble 
                   from={to || "Penerima"} 
                   message={message} 
                   mediaUrl={previewUrl} 
                   mediaType={attachedFile ? determineContentType(attachedFile) : undefined} 
                />
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
