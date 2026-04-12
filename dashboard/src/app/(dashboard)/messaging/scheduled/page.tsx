"use client";

import { useState } from "react";
import { CalendarClock, Plus, X, Clock, Send, Phone, Users, CheckSquare, Square, Trash2 } from "lucide-react";

interface ScheduledMsg {
  id: number;
  to: string;
  toLabel: string;
  message: string;
  device: string;
  scheduledAt: string;
  status: "pending" | "sent" | "cancelled";
}

const DEVICES = ["Device Utama", "Device Marketing", "Device CS", "Device Backup"];

const GROUPS = [
  {
    id: 1,
    name: "Pelanggan VIP",
    description: "Pelanggan tier premium",
    count: 2,
    members: [
      { name: "Budi Santoso", phone: "628111000001" },
      { name: "Dewi Lestari", phone: "628111000004" },
    ],
  },
  {
    id: 2,
    name: "Tim Internal",
    description: "Staff dan karyawan",
    count: 1,
    members: [{ name: "Hendra Kurniawan", phone: "628111000005" }],
  },
  {
    id: 3,
    name: "Prospek Aktif",
    description: "Calon pelanggan yang sedang difollow up",
    count: 2,
    members: [
      { name: "Ahmad Rizki", phone: "628111000003" },
      { name: "Rina Widiastuti", phone: "628111000006" },
    ],
  },
];

const INITIAL: ScheduledMsg[] = [
  { id: 1, to: "628111222333", toLabel: "628111222333", message: "Halo! Pesanan kamu siap dikirim hari ini.", device: "Device Utama", scheduledAt: "2026-04-12 14:00", status: "pending" },
  { id: 2, to: "Pelanggan VIP", toLabel: "2 kontak (group)", message: "Promo akhir bulan: diskon 30% berlaku s.d besok!", device: "Device Marketing", scheduledAt: "2026-04-12 18:00", status: "pending" },
  { id: 3, to: "628333444555", toLabel: "628333444555", message: "Selamat ulang tahun! 🎂 Ada hadiah spesial untukmu.", device: "Device CS", scheduledAt: "2026-04-13 07:00", status: "pending" },
  { id: 4, to: "628444555666", toLabel: "628444555666", message: "Tagihan bulan April sudah bisa dicek.", device: "Device Utama", scheduledAt: "2026-04-11 09:00", status: "sent" },
  { id: 5, to: "Prospek Aktif", toLabel: "2 kontak (group)", message: "Flash sale dimulai jam 12 siang!", device: "Device Marketing", scheduledAt: "2026-04-10 11:30", status: "cancelled" },
];

const STATUS_STYLE: Record<string, string> = {
  pending: "bg-yellow-50 text-yellow-700",
  sent: "bg-green-50 text-green-700",
  cancelled: "bg-gray-100 text-gray-500",
};

const STATUS_LABEL: Record<string, string> = {
  pending: "Terjadwal",
  sent: "Terkirim",
  cancelled: "Dibatalkan",
};

type RecipientMode = "manual" | "group";

export default function ScheduledPage() {
  const [messages, setMessages] = useState<ScheduledMsg[]>(INITIAL);
  const [showModal, setShowModal] = useState(false);

  // Modal state
  const [recipientMode, setRecipientMode] = useState<RecipientMode>("manual");
  const [manualTo, setManualTo] = useState("");
  const [selectedGroups, setSelectedGroups] = useState<number[]>([]);
  const [formMessage, setFormMessage] = useState("");
  const [formDevice, setFormDevice] = useState(DEVICES[0]);
  const [formScheduledAt, setFormScheduledAt] = useState("");

  const pending = messages.filter((m) => m.status === "pending");
  const history = messages.filter((m) => m.status !== "pending");

  const groupCount = GROUPS.filter((g) => selectedGroups.includes(g.id)).reduce((s, g) => s + g.count, 0);
  const canSubmit =
    formMessage.trim() &&
    formScheduledAt &&
    (recipientMode === "manual" ? /^[0-9]{9,15}$/.test(manualTo) : selectedGroups.length > 0);

  function toggleGroup(id: number) {
    setSelectedGroups((prev) => prev.includes(id) ? prev.filter((g) => g !== id) : [...prev, id]);
  }

  function resetModal() {
    setRecipientMode("manual");
    setManualTo("");
    setSelectedGroups([]);
    setFormMessage("");
    setFormDevice(DEVICES[0]);
    setFormScheduledAt("");
    setShowModal(false);
  }

  function addScheduled() {
    if (!canSubmit) return;
    let to: string;
    let toLabel: string;
    if (recipientMode === "manual") {
      to = manualTo;
      toLabel = manualTo;
    } else {
      const names = GROUPS.filter((g) => selectedGroups.includes(g.id)).map((g) => g.name).join(", ");
      to = names;
      toLabel = `${groupCount} kontak (group)`;
    }
    setMessages([
      { id: Date.now(), to, toLabel, message: formMessage, device: formDevice, scheduledAt: formScheduledAt, status: "pending" },
      ...messages,
    ]);
    resetModal();
  }

  function cancelMsg(id: number) {
    setMessages(messages.map((m) => (m.id === id ? { ...m, status: "cancelled" } : m)));
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold text-gray-900">Scheduled Message</h1>
            <p className="text-sm text-gray-500">Jadwalkan pengiriman pesan otomatis</p>
          </div>
          <button
            onClick={() => setShowModal(true)}
            className="flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700"
          >
            <Plus className="h-4 w-4" /> Jadwalkan Pesan
          </button>
        </div>
      </div>

      <div className="p-6 space-y-6">
        {/* Stats */}
        <div className="grid grid-cols-3 gap-4">
          {[
            { label: "Terjadwal", value: messages.filter((m) => m.status === "pending").length, color: "text-yellow-600", bg: "bg-yellow-50" },
            { label: "Terkirim", value: messages.filter((m) => m.status === "sent").length, color: "text-green-600", bg: "bg-green-50" },
            { label: "Dibatalkan", value: messages.filter((m) => m.status === "cancelled").length, color: "text-gray-500", bg: "bg-gray-50" },
          ].map((s) => (
            <div key={s.label} className={`rounded-xl border border-gray-200 ${s.bg} p-4 shadow-sm`}>
              <p className="text-xs font-medium text-gray-500">{s.label}</p>
              <p className={`mt-1 text-3xl font-bold ${s.color}`}>{s.value}</p>
            </div>
          ))}
        </div>

        {/* Pending queue */}
        <div className="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
          <div className="flex items-center gap-2 border-b border-gray-100 px-5 py-4">
            <Clock className="h-4 w-4 text-yellow-500" />
            <h2 className="font-semibold text-gray-900">Antrian Terjadwal</h2>
            <span className="ml-auto rounded-full bg-yellow-100 px-2 py-0.5 text-xs font-semibold text-yellow-700">
              {pending.length} pesan
            </span>
          </div>
          {pending.length === 0 ? (
            <div className="p-10 text-center">
              <CalendarClock className="mx-auto h-8 w-8 text-gray-200" />
              <p className="mt-2 text-sm text-gray-400">Tidak ada pesan terjadwal</p>
            </div>
          ) : (
            <div className="divide-y divide-gray-50">
              {pending.map((m) => (
                <div key={m.id} className="flex items-start justify-between gap-4 px-5 py-4 hover:bg-gray-50">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="font-mono text-sm font-medium text-gray-700">{m.to}</span>
                      {m.toLabel !== m.to && (
                        <span className="rounded-full bg-purple-50 px-2 py-0.5 text-xs text-purple-600">{m.toLabel}</span>
                      )}
                      <span className="rounded-full bg-blue-50 px-2 py-0.5 text-xs text-blue-600">{m.device}</span>
                    </div>
                    <p className="mt-0.5 truncate text-sm text-gray-500">{m.message}</p>
                  </div>
                  <div className="flex flex-shrink-0 items-center gap-3">
                    <div className="text-right">
                      <p className="text-xs font-medium text-yellow-700">{m.scheduledAt}</p>
                      <p className="text-xs text-gray-400">WIB</p>
                    </div>
                    <button
                      onClick={() => cancelMsg(m.id)}
                      className="rounded-lg border border-gray-200 px-2.5 py-1.5 text-xs text-gray-500 hover:border-red-300 hover:text-red-500"
                    >
                      Batalkan
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* History */}
        <div className="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
          <div className="flex items-center gap-2 border-b border-gray-100 px-5 py-4">
            <Send className="h-4 w-4 text-gray-400" />
            <h2 className="font-semibold text-gray-900">Riwayat</h2>
          </div>
          <div className="overflow-x-auto">
          <table className="w-full min-w-[600px] text-sm">
            <thead>
              <tr className="border-b border-gray-50 bg-gray-50 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">
                <th className="px-5 py-3">Tujuan</th>
                <th className="px-5 py-3">Pesan</th>
                <th className="px-5 py-3">Device</th>
                <th className="px-5 py-3">Waktu</th>
                <th className="px-5 py-3">Status</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-50">
              {history.map((m) => (
                <tr key={m.id} className="hover:bg-gray-50">
                  <td className="px-5 py-3">
                    <p className="font-mono text-gray-600">{m.to}</p>
                    {m.toLabel !== m.to && <p className="text-xs text-purple-500">{m.toLabel}</p>}
                  </td>
                  <td className="px-5 py-3 max-w-[200px] truncate text-gray-500">{m.message}</td>
                  <td className="px-5 py-3 text-gray-500">{m.device}</td>
                  <td className="px-5 py-3 text-gray-400">{m.scheduledAt}</td>
                  <td className="px-5 py-3">
                    <span className={`rounded-full px-2.5 py-1 text-xs font-medium ${STATUS_STYLE[m.status]}`}>
                      {STATUS_LABEL[m.status]}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          </div>
        </div>
      </div>

      {/* Modal */}
      {showModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/30 backdrop-blur-sm p-4">
          <div className="w-full max-w-md rounded-2xl bg-white shadow-xl overflow-hidden">
            {/* Modal header */}
            <div className="flex items-center justify-between border-b border-gray-100 px-6 py-4">
              <h2 className="text-lg font-bold text-gray-900">Jadwalkan Pesan</h2>
              <button onClick={resetModal}>
                <X className="h-5 w-5 text-gray-400" />
              </button>
            </div>

            <div className="max-h-[70vh] overflow-y-auto px-6 py-5 space-y-4">
              {/* Recipient mode toggle */}
              <div>
                <label className="mb-2 block text-sm font-medium text-gray-700">Penerima</label>
                <div className="flex items-center gap-1 rounded-lg bg-gray-100 p-1">
                  <button
                    type="button"
                    onClick={() => setRecipientMode("manual")}
                    className={`flex flex-1 items-center justify-center gap-2 rounded-md py-2 text-sm font-medium transition-all ${
                      recipientMode === "manual" ? "bg-white text-gray-900 shadow-sm" : "text-gray-500 hover:text-gray-700"
                    }`}
                  >
                    <Phone className="h-3.5 w-3.5" /> By Nomor
                  </button>
                  <button
                    type="button"
                    onClick={() => setRecipientMode("group")}
                    className={`flex flex-1 items-center justify-center gap-2 rounded-md py-2 text-sm font-medium transition-all ${
                      recipientMode === "group" ? "bg-white text-gray-900 shadow-sm" : "text-gray-500 hover:text-gray-700"
                    }`}
                  >
                    <Users className="h-3.5 w-3.5" /> By Group
                  </button>
                </div>
              </div>

              {/* Manual input */}
              {recipientMode === "manual" && (
                <div>
                  <input
                    value={manualTo}
                    onChange={(e) => setManualTo(e.target.value.replace(/\D/g, ""))}
                    placeholder="628xxxxxxxxxx"
                    inputMode="numeric"
                    className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm font-mono focus:border-green-500 focus:outline-none"
                  />
                  <p className="mt-1 text-xs text-gray-400">Masukkan nomor dalam format internasional</p>
                </div>
              )}

              {/* Group picker */}
              {recipientMode === "group" && (
                <div className="space-y-2">
                  {GROUPS.map((g) => {
                    const selected = selectedGroups.includes(g.id);
                    return (
                      <button
                        key={g.id}
                        type="button"
                        onClick={() => toggleGroup(g.id)}
                        className={`flex w-full items-center gap-3 rounded-xl border-2 px-4 py-3 text-left transition-all ${
                          selected ? "border-green-400 bg-green-50" : "border-gray-100 bg-white hover:border-gray-200"
                        }`}
                      >
                        {selected ? (
                          <CheckSquare className="h-5 w-5 flex-shrink-0 text-green-600" />
                        ) : (
                          <Square className="h-5 w-5 flex-shrink-0 text-gray-300" />
                        )}
                        <div className="flex flex-1 items-center justify-between min-w-0">
                          <div className="min-w-0">
                            <p className={`text-sm font-semibold ${selected ? "text-green-800" : "text-gray-800"}`}>{g.name}</p>
                            <p className="text-xs text-gray-400 truncate">{g.description}</p>
                          </div>
                          <span className={`ml-3 flex-shrink-0 rounded-full px-2.5 py-1 text-xs font-semibold ${selected ? "bg-green-200 text-green-800" : "bg-gray-100 text-gray-500"}`}>
                            {g.count} kontak
                          </span>
                        </div>
                      </button>
                    );
                  })}
                  {selectedGroups.length > 0 && (
                    <p className="text-xs text-green-700 font-medium pt-1">
                      {groupCount} kontak dari {selectedGroups.length} group terpilih
                    </p>
                  )}
                </div>
              )}

              {/* Device */}
              <div>
                <label className="mb-1 block text-sm font-medium text-gray-700">Device Pengirim</label>
                <select
                  value={formDevice}
                  onChange={(e) => setFormDevice(e.target.value)}
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-green-500 focus:outline-none"
                >
                  {DEVICES.map((d) => <option key={d}>{d}</option>)}
                </select>
              </div>

              {/* Scheduled at */}
              <div>
                <label className="mb-1 block text-sm font-medium text-gray-700">Waktu Kirim</label>
                <input
                  type="datetime-local"
                  value={formScheduledAt}
                  onChange={(e) => setFormScheduledAt(e.target.value)}
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-green-500 focus:outline-none"
                />
              </div>

              {/* Message */}
              <div>
                <label className="mb-1 block text-sm font-medium text-gray-700">Pesan</label>
                <textarea
                  value={formMessage}
                  onChange={(e) => setFormMessage(e.target.value)}
                  rows={4}
                  placeholder="Tulis pesan..."
                  className="w-full resize-none rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-green-500 focus:outline-none"
                />
              </div>
            </div>

            {/* Modal footer */}
            <div className="flex gap-3 border-t border-gray-100 px-6 py-4">
              <button
                onClick={resetModal}
                className="flex-1 rounded-lg border border-gray-200 py-2 text-sm text-gray-600 hover:bg-gray-50"
              >
                Batal
              </button>
              <button
                onClick={addScheduled}
                disabled={!canSubmit}
                className="flex-1 rounded-lg bg-green-600 py-2 text-sm font-medium text-white hover:bg-green-700 disabled:bg-gray-200 disabled:text-gray-400 disabled:cursor-not-allowed"
              >
                Jadwalkan
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
