"use client";

import { useState, useEffect } from "react";
import { CalendarClock, Plus, X, Clock, Send, Phone, Users, CheckSquare, Square, Trash2, Loader2, Smartphone, AlertCircle } from "lucide-react";
import { sessionsApi, messagesApi } from "@/lib/api";
import { Device, Message } from "@/lib/types";
import { toast } from "sonner";
import { format } from "date-fns";
import { id } from "date-fns/locale";

type RecipientMode = "manual" | "group";

export default function ScheduledPage() {
  const [devices, setDevices] = useState<Device[]>([]);
  const [pendingMessages, setPendingMessages] = useState<Message[]>([]);
  const [historyMessages, setHistoryMessages] = useState<Message[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Modal state
  const [recipientMode, setRecipientMode] = useState<RecipientMode>("manual");
  const [manualTo, setManualTo] = useState("");
  const [formMessage, setFormMessage] = useState("");
  const [formDevice, setFormDevice] = useState("");
  const [formScheduledAt, setFormScheduledAt] = useState("");
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  // Status mapping helper
  const getStatusInfo = (status: string) => {
    switch (status) {
      case "pending": return { label: "Pending", color: "bg-yellow-50 text-yellow-600", dot: "bg-yellow-500" };
      case "sent": return { label: "Sent", color: "bg-blue-50 text-blue-600", dot: "bg-blue-500" };
      case "delivered": return { label: "Delivered", color: "bg-green-50 text-green-600", dot: "bg-green-500" };
      case "read": return { label: "Read", color: "bg-green-100 text-green-700", dot: "bg-green-600" };
      case "failed": return { label: "Failed", color: "bg-red-50 text-red-600", dot: "bg-red-500" };
      default: return { label: status || "Unknown", color: "bg-gray-50 text-gray-500", dot: "bg-gray-400" };
    }
  };

  const fetchData = async () => {
    setIsLoading(true);
    try {
      const sessRes = await sessionsApi.list();
      const activeDevices = (sessRes.sessions || []).filter(d => d.status === 1);
      setDevices(activeDevices);
      if (activeDevices.length > 0 && !formDevice) {
        setFormDevice(activeDevices[0].device_id);
      }

      if (activeDevices.length > 0) {
        // Fetch from first device as default view, or all if you prefer
        // For simplicity, we fetch based on the first active device if any
        const devId = activeDevices[0].device_id;
        const [schRes, histRes] = await Promise.all([
          messagesApi.listScheduled(devId),
          messagesApi.listHistory(devId)
        ]);
        setPendingMessages(schRes.messages || []);
        setHistoryMessages(histRes.messages || []);
      }
    } catch (err) {
      console.error(err);
      toast.error("Gagal memuat data pesan terjadwal.");
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const canSubmit =
    (formMessage.trim() || selectedFile) &&
    formScheduledAt &&
    formDevice &&
    (recipientMode === "manual" ? /^[0-9]{9,15}$/.test(manualTo) : false);

  function resetModal() {
    setRecipientMode("manual");
    setManualTo("");
    setFormMessage("");
    setFormScheduledAt("");
    setSelectedFile(null);
    setShowModal(false);
  }

  async function addScheduled() {
    if (!canSubmit) return;
    setIsSubmitting(true);
    try {
      if (selectedFile) {
        const formData = new FormData();
        formData.append("target_jid", manualTo);
        formData.append("content", formMessage);
        formData.append("scheduled_for", new Date(formScheduledAt).toISOString());
        formData.append("file", selectedFile);
        
        const isImage = selectedFile.type.startsWith('image/');
        formData.append("content_type", isImage ? "image" : "document");
        formData.append("caption", formMessage);

        await messagesApi.schedule(formDevice, formData);
      } else {
        await messagesApi.schedule(formDevice, {
          target_jid: manualTo,
          content: formMessage,
          scheduled_for: new Date(formScheduledAt).toISOString()
        });
      }

      toast.success("Pesan berhasil dijadwalkan.");
      resetModal();
      fetchData();
    } catch (err: any) {
      if (err.response?.status === 403) {
        toast.error("Batas Pesan Tercapai", {
          description: "Anda telah mencapai kuota pesan harian. Silakan upgrade paket Anda untuk menjadwalkan lebih banyak pesan."
        });
      } else {
        toast.error(err.response?.data?.error || "Gagal menjadwalkan pesan.");
      }
    } finally {
      setIsSubmitting(false);
    }
  }

  async function cancelMsg(id: string) {
    try {
      await messagesApi.cancelScheduled(id);
      toast.success("Jadwal pesan dibatalkan.");
      fetchData();
    } catch (err) {
      toast.error("Gagal membatalkan jadwal.");
    }
  }

  if (isLoading && devices.length === 0) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50">
        <Loader2 className="h-8 w-8 animate-spin text-green-600" />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 pb-20">
      <div className="border-b border-gray-200 bg-white px-6 py-4 sticky top-0 z-10">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold text-gray-900">Scheduled Message</h1>
            <p className="text-sm text-gray-500">Jadwalkan pengiriman pesan otomatis</p>
          </div>
          <button
            onClick={() => setShowModal(true)}
            disabled={devices.length === 0}
            className="flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-sm font-bold text-white hover:bg-green-700 shadow-lg shadow-green-100 disabled:bg-gray-200 disabled:shadow-none transition-all"
          >
            <Plus className="h-4 w-4" /> Jadwalkan Pesan
          </button>
        </div>
      </div>

      <div className="p-6 space-y-6">
        {devices.length === 0 && (
          <div className="flex items-center gap-3 rounded-xl border border-amber-200 bg-amber-50 p-4 text-amber-800">
            <AlertCircle className="h-5 w-5 flex-shrink-0" />
            <p className="text-sm">Tidak ada perangkat aktif. Harap hubungkan perangkat terlebih dahulu di menu <b>Connection Status</b>.</p>
          </div>
        )}

        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-yellow-50 text-yellow-600">
                <Clock className="h-5 w-5" />
              </div>
              <div>
                <p className="text-xs font-medium text-gray-500 uppercase tracking-wider">Antrian Aktif</p>
                <p className="text-2xl font-bold text-gray-900">{pendingMessages.length}</p>
              </div>
            </div>
          </div>
          <div className="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-green-50 text-green-600">
                <Send className="h-5 w-5" />
              </div>
              <div>
                <p className="text-xs font-medium text-gray-500 uppercase tracking-wider">Pesan Terkirim</p>
                <p className="text-2xl font-bold text-gray-900">
                  {historyMessages.filter(m => {
                    return m.status === "sent" || m.status === "delivered" || m.status === "read";
                  }).length}
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Pending queue */}
        <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-sm">
          <div className="flex items-center gap-2 border-b border-gray-100 px-6 py-4 bg-gray-50/50">
            <Clock className="h-4 w-4 text-yellow-500" />
            <h2 className="font-bold text-gray-900">Antrian Terjadwal</h2>
          </div>
          {pendingMessages.length === 0 ? (
            <div className="p-16 text-center">
              <CalendarClock className="mx-auto h-12 w-12 text-gray-100 mb-2" />
              <p className="text-sm text-gray-400">Tidak ada pesan yang sedang dijadwalkan</p>
            </div>
          ) : (
            <div className="divide-y divide-gray-50">
              {pendingMessages.map((m) => {
                const isMedia = m.content_type && m.content_type !== "text";
                return (
                <div key={m.id} className="group flex items-center justify-between gap-4 px-6 py-5 hover:bg-gray-50/80 transition-colors">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-3">
                      <span className="font-mono text-sm font-bold text-gray-900">{m.target_jid}</span>
                      <span className="inline-flex items-center gap-1 rounded-full bg-blue-50 px-2 py-0.5 text-[10px] font-bold text-blue-600 uppercase">
                        <Smartphone className="h-2.5 w-2.5" /> {devices.find(d => d.device_id === m.device_id)?.display_name || "Device"}
                      </span>
                      {isMedia && (
                        <span className="inline-flex items-center gap-1 rounded-full bg-purple-50 px-2 py-0.5 text-[10px] font-bold text-purple-600 uppercase">
                          <Smartphone className="h-2.5 w-2.5" /> {m.content_type}
                        </span>
                      )}
                    </div>
                    <p className="mt-1 line-clamp-1 text-sm text-gray-500">{m.content || "(No text content)"}</p>
                  </div>
                  <div className="flex flex-shrink-0 items-center gap-4">
                    <div className="text-right">
                      <p className="text-sm font-bold text-yellow-600">
                        {m.scheduled_for ? format(new Date(m.scheduled_for), "HH:mm, dd MMM", { locale: id }) : "-"}
                      </p>
                      <p className="text-[10px] text-gray-400 uppercase font-medium tracking-tighter">Terjadwal</p>
                    </div>
                    <button
                      onClick={() => cancelMsg(m.id)}
                      className="opacity-0 group-hover:opacity-100 flex h-9 w-9 items-center justify-center rounded-xl bg-red-50 text-red-500 hover:bg-red-100 transition-all"
                      title="Batalkan"
                    >
                      <Trash2 className="h-4 w-4" />
                    </button>
                  </div>
                </div>
                );
              })}
            </div>
          )}
        </div>

        {/* History */}
        <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-sm">
          <div className="flex items-center gap-2 border-b border-gray-100 px-6 py-4 bg-gray-50/50">
            <Send className="h-4 w-4 text-gray-400" />
            <h2 className="font-bold text-gray-900">Riwayat Pengiriman</h2>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-gray-50 bg-gray-50/30 text-left text-[10px] font-bold uppercase tracking-widest text-gray-400">
                  <th className="px-6 py-3">Penerima</th>
                  <th className="px-6 py-3">Pesan</th>
                  <th className="px-6 py-3">Status</th>
                  <th className="px-6 py-3">Waktu</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-50">
                {historyMessages.map((m) => {
                  const statusInfo = getStatusInfo(m.status);
                  const isMedia = m.content_type && m.content_type !== "text";
                  return (
                    <tr key={m.id} className="hover:bg-gray-50/50">
                      <td className="px-6 py-4">
                        <p className="font-mono font-medium text-gray-700">{m.target_jid}</p>
                        <p className="text-[10px] text-gray-400 font-bold uppercase tracking-tighter">
                          Via {devices.find(d => d.device_id === m.device_id)?.display_name || "Device"}
                        </p>
                      </td>
                      <td className="px-6 py-4">
                        <div className="flex items-center gap-2">
                          {isMedia && <Smartphone className="h-3.5 w-3.5 text-purple-500" />}
                          <p className="max-w-xs truncate text-gray-500">{m.content}</p>
                        </div>
                        {m.status === "failed" && m.error_log && (
                          <p className="text-[10px] text-red-500 mt-1 italic max-w-xs truncate" title={m.error_log}>
                            Error: {m.error_log}
                          </p>
                        )}
                      </td>
                      <td className="px-6 py-4">
                        <span className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-[10px] font-bold uppercase tracking-wider ${statusInfo.color}`}>
                          <span className={`h-1 w-1 rounded-full ${statusInfo.dot}`} />
                          {statusInfo.label}
                        </span>
                      </td>
                      <td className="px-6 py-4 text-xs text-gray-400">
                        {m.sent_at ? format(new Date(m.sent_at), "dd/MM/yy HH:mm") : format(new Date(m.created_at), "dd/MM/yy HH:mm")}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      {/* Modal */}
      {showModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in duration-200">
          <div className="w-full max-w-md rounded-3xl bg-white shadow-2xl overflow-hidden animate-in zoom-in-95 duration-200">
            {/* Modal header */}
            <div className="flex items-center justify-between border-b border-gray-100 px-6 py-5">
              <h2 className="text-xl font-bold text-gray-900">Jadwalkan Pesan</h2>
              <button 
                onClick={resetModal}
                className="rounded-full p-2 hover:bg-gray-100 transition-colors"
                disabled={isSubmitting}
              >
                <X className="h-5 w-5 text-gray-400" />
              </button>
            </div>

            <div className="max-h-[70vh] overflow-y-auto px-6 py-6 space-y-5">
              {/* Recipient Input */}
              <div>
                <label className="mb-2 block text-sm font-bold text-gray-700">Nomor WhatsApp Penerima</label>
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
                    <Phone className="h-4 w-4 text-gray-400" />
                  </div>
                  <input
                    value={manualTo}
                    onChange={(e) => setManualTo(e.target.value.replace(/\D/g, ""))}
                    placeholder="628xxxxxxxxxx"
                    disabled={isSubmitting}
                    className="block w-full rounded-xl border border-gray-200 pl-10 pr-3 py-3 text-sm font-mono focus:border-green-500 focus:ring-1 focus:ring-green-500 focus:outline-none bg-gray-50/50"
                  />
                </div>
                <p className="mt-1.5 text-[10px] text-gray-400 font-medium">Contoh: 628123456789 (format internasional)</p>
              </div>

              {/* Device */}
              <div>
                <label className="mb-2 block text-sm font-bold text-gray-700">Device Pengirim</label>
                <select
                  value={formDevice}
                  onChange={(e) => setFormDevice(e.target.value)}
                  disabled={isSubmitting}
                  className="w-full rounded-xl border border-gray-200 px-3 py-3 text-sm focus:border-green-500 focus:outline-none bg-gray-50/50"
                >
                  {devices.map((d) => (
                    <option key={d.device_id} value={d.device_id}>
                      {d.display_name || d.device_id} ({d.phone || "No Phone"})
                    </option>
                  ))}
                </select>
              </div>

              {/* Scheduled at */}
              <div>
                <label className="mb-2 block text-sm font-bold text-gray-700">Waktu Kirim (Server Time)</label>
                <input
                  type="datetime-local"
                  value={formScheduledAt}
                  onChange={(e) => setFormScheduledAt(e.target.value)}
                  disabled={isSubmitting}
                  className="w-full rounded-xl border border-gray-200 px-3 py-3 text-sm focus:border-green-500 focus:outline-none bg-gray-50/50"
                />
              </div>

              {/* File Attachment */}
              <div className="space-y-2">
                <label className="block text-sm font-bold text-gray-700">Lampiran Media (Opsional)</label>
                <div 
                  className={`group relative flex flex-col items-center justify-center rounded-2xl border-2 border-dashed p-6 transition-all hover:bg-gray-50 ${
                    selectedFile ? "border-green-500 bg-green-50" : "border-gray-200 bg-gray-50/30"
                  }`}
                >
                  <input
                    type="file"
                    className="absolute inset-0 z-10 cursor-pointer opacity-0"
                    onChange={(e) => {
                      const file = e.target.files?.[0];
                      if (file) setSelectedFile(file);
                    }}
                    disabled={isSubmitting}
                  />
                  {selectedFile ? (
                    <div className="flex w-full items-center justify-between">
                      <div className="flex items-center gap-3">
                        <div className="rounded-lg bg-green-100 p-2 text-green-600">
                          <Plus className="h-5 w-5" />
                        </div>
                        <div className="min-w-0">
                          <p className="truncate text-sm font-bold text-gray-900">{selectedFile.name}</p>
                          <p className="text-[10px] text-gray-500 uppercase font-bold tracking-widest">
                            {(selectedFile.size / 1024).toFixed(1)} KB
                          </p>
                        </div>
                      </div>
                      <button 
                        onClick={(e) => {
                          e.stopPropagation();
                          setSelectedFile(null);
                        }}
                        className="rounded-full p-2 text-gray-400 hover:bg-red-50 hover:text-red-500 transition-colors"
                      >
                        <X className="h-4 w-4" />
                      </button>
                    </div>
                  ) : (
                    <div className="text-center">
                      <div className="mx-auto mb-2 flex h-10 w-10 items-center justify-center rounded-xl bg-white text-gray-400 shadow-sm transition-transform group-hover:scale-110">
                        <Plus className="h-5 w-5" />
                      </div>
                      <p className="text-sm font-bold text-gray-600">Klik atau drag file ke sini</p>
                      <p className="text-[10px] text-gray-400 font-medium">PNG, JPG, PDF (Max 50MB)</p>
                    </div>
                  )}
                </div>
              </div>

              {/* Message / Caption */}
              <div>
                <label className="mb-2 block text-sm font-bold text-gray-700">
                  {selectedFile ? "Keterangan (Caption)" : "Isi Pesan"}
                </label>
                <textarea
                  value={formMessage}
                  onChange={(e) => setFormMessage(e.target.value)}
                  rows={selectedFile ? 2 : 4}
                  placeholder={selectedFile ? "Tambahkan keterangan untuk media..." : "Tulis pesan yang ingin dijadwalkan..."}
                  disabled={isSubmitting}
                  className="w-full resize-none rounded-xl border border-gray-200 px-4 py-3 text-sm focus:border-green-500 focus:outline-none bg-gray-50/50"
                />
              </div>
            </div>

            {/* Modal footer */}
            <div className="flex gap-3 border-t border-gray-100 px-6 py-5 bg-gray-50/30">
              <button
                onClick={resetModal}
                disabled={isSubmitting}
                className="flex-1 rounded-xl border border-gray-200 py-3 text-sm font-bold text-gray-600 hover:bg-gray-100 transition-colors"
              >
                Batal
              </button>
              <button
                onClick={addScheduled}
                disabled={!canSubmit || isSubmitting}
                className="flex-1 rounded-xl bg-green-600 py-3 text-sm font-bold text-white hover:bg-green-700 shadow-lg shadow-green-100 disabled:bg-gray-200 disabled:shadow-none transition-all flex items-center justify-center gap-2"
              >
                {isSubmitting ? (
                  <>
                    <Loader2 className="h-4 w-4 animate-spin" /> Sedang Memproses...
                  </>
                ) : (
                  "Simpan Jadwal"
                )}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
