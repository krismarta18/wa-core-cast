"use client";

import { useState } from "react";
import { Send, Upload, Plus, Trash2, ChevronDown, Info, Users, Phone, CheckSquare, Square, CheckCheck, Smile, Loader2 } from "lucide-react";
import { useToast } from "@/components/ui/toast";

const DEVICES = ["Device Utama", "Device CS", "Device Marketing", "Device Backup"];
const TEMPLATES = ["(Tidak pakai template)", "Promo Bulanan", "Konfirmasi Pesanan", "Pengingat Pembayaran"];

const GROUPS = [
  {
    id: 1, name: "Pelanggan VIP", description: "Pelanggan tier premium", count: 2,
    members: [
      { name: "Budi Santoso", phone: "628111000001" },
      { name: "Dewi Lestari", phone: "628111000004" },
    ],
  },
  {
    id: 2, name: "Tim Internal", description: "Staff dan karyawan", count: 1,
    members: [{ name: "Hendra Kurniawan", phone: "628111000005" }],
  },
  {
    id: 3, name: "Prospek Aktif", description: "Calon pelanggan yang sedang difollow up", count: 2,
    members: [
      { name: "Ahmad Rizki", phone: "628111000003" },
      { name: "Rina Widiastuti", phone: "628111000006" },
    ],
  },
];

type RecipientMode = "group" | "manual";

function WABubble({ message }: { message: string }) {
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
        className="flex-1 p-4"
        style={{
          backgroundImage: `url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23128c7e' fill-opacity='0.04'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")`,
          backgroundColor: "#ECE5DD",
        }}
      >
        {message.trim() ? (
          <div className="flex justify-end">
            <div className="max-w-[85%] rounded-lg rounded-br-sm bg-[#DCF8C6] px-3 py-2 shadow-sm">
              <p className="whitespace-pre-wrap break-words text-sm text-gray-800">{message}</p>
              <div className="mt-1 flex items-center justify-end gap-1">
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
  const { success, error } = useToast();
  const [recipientMode, setRecipientMode] = useState<RecipientMode>("manual");
  const [recipients, setRecipients] = useState<string[]>([""]);
  const [selectedGroups, setSelectedGroups] = useState<number[]>([]);
  const [message, setMessage] = useState("");
  const [device, setDevice] = useState(DEVICES[0]);
  const [template, setTemplate] = useState(TEMPLATES[0]);
  const [delay, setDelay] = useState("3");
  const [sending, setSending] = useState(false);

  const addRecipient = () => setRecipients((r) => [...r, ""]);
  const removeRecipient = (i: number) => setRecipients((r) => r.filter((_, idx) => idx !== i));
  const updateRecipient = (i: number, val: string) =>
    setRecipients((r) => r.map((v, idx) => (idx === i ? val : v)));
  const validCount = recipients.filter((r) => /^[0-9]{9,15}$/.test(r)).length;

  const toggleGroup = (id: number) =>
    setSelectedGroups((prev) =>
      prev.includes(id) ? prev.filter((g) => g !== id) : [...prev, id]
    );
  const groupRecipientCount = GROUPS.filter((g) => selectedGroups.includes(g.id)).reduce(
    (sum, g) => sum + g.count, 0
  );

  const totalCount = recipientMode === "manual" ? validCount : groupRecipientCount;
  const canSend = totalCount > 0 && message.trim().length > 0;

  function handleSend() {
    if (!canSend || sending) return;
    setSending(true);
    setTimeout(() => {
      setSending(false);
      success("Broadcast Dikirim!", `Pesan berhasil dikirim ke ${totalCount} penerima via ${device}.`);
      setMessage("");
      setRecipients([""]);
      setSelectedGroups([]);
    }, 1800);
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <h1 className="text-xl font-bold text-gray-900">Broadcast / Bulk Sender</h1>
        <p className="text-sm text-gray-500">Kirim pesan ke banyak nomor sekaligus</p>
      </div>

      <div className="p-6">
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-5">
          {/* Form col */}
          <div className="space-y-5 lg:col-span-3">

            {/* Device selector */}
            <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
              <label className="mb-1 block text-sm font-medium text-gray-700">Kirim dari Device</label>
              <div className="relative">
                <select
                  value={device}
                  onChange={(e) => setDevice(e.target.value)}
                  className="h-10 w-full appearance-none rounded-lg border border-gray-300 px-3 pr-8 text-sm text-gray-700 focus:outline-none focus:ring-2 focus:ring-green-500"
                >
                  {DEVICES.map((d) => <option key={d}>{d}</option>)}
                </select>
                <ChevronDown className="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
              </div>
            </div>

            {/* Recipients */}
            <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
              <div className="mb-4 flex items-center gap-1 rounded-lg bg-gray-100 p-1">
                <button
                  onClick={() => setRecipientMode("manual")}
                  className={`flex flex-1 items-center justify-center gap-2 rounded-md py-2 text-sm font-medium transition-all ${
                    recipientMode === "manual" ? "bg-white text-gray-900 shadow-sm" : "text-gray-500 hover:text-gray-700"
                  }`}
                >
                  <Phone className="h-4 w-4" /> By Nomor
                </button>
                <button
                  onClick={() => setRecipientMode("group")}
                  className={`flex flex-1 items-center justify-center gap-2 rounded-md py-2 text-sm font-medium transition-all ${
                    recipientMode === "group" ? "bg-white text-gray-900 shadow-sm" : "text-gray-500 hover:text-gray-700"
                  }`}
                >
                  <Users className="h-4 w-4" /> By Group
                </button>
              </div>

              {recipientMode === "manual" && (
                <>
                  <div className="mb-3 flex items-center justify-between">
                    <span className="text-sm font-medium text-gray-700">
                      Daftar Nomor
                      {validCount > 0 && (
                        <span className="ml-2 rounded-full bg-green-100 px-2 py-0.5 text-xs font-semibold text-green-700">{validCount} valid</span>
                      )}
                    </span>
                    <button className="inline-flex items-center gap-1.5 text-xs font-medium text-green-600 hover:underline">
                      <Upload className="h-3.5 w-3.5" /> Import CSV
                    </button>
                  </div>
                  <div className="max-h-52 space-y-2 overflow-y-auto pr-1">
                    {recipients.map((r, i) => (
                      <div key={i} className="flex gap-2">
                        <input
                          type="tel" inputMode="numeric" placeholder="628xxxxxxxxx"
                          value={r}
                          onChange={(e) => updateRecipient(i, e.target.value.replace(/\D/g, ""))}
                          className="h-9 flex-1 rounded-lg border border-gray-300 px-3 text-sm font-mono placeholder:text-gray-300 focus:outline-none focus:ring-2 focus:ring-green-500"
                        />
                        {recipients.length > 1 && (
                          <button onClick={() => removeRecipient(i)} className="rounded-lg p-2 text-gray-400 hover:bg-red-50 hover:text-red-500">
                            <Trash2 className="h-4 w-4" />
                          </button>
                        )}
                      </div>
                    ))}
                  </div>
                  <button onClick={addRecipient} className="mt-3 inline-flex items-center gap-1.5 text-sm font-medium text-green-600 hover:underline">
                    <Plus className="h-4 w-4" /> Tambah Nomor
                  </button>
                </>
              )}

              {recipientMode === "group" && (
                <>
                  <div className="mb-3 flex items-center justify-between">
                    <span className="text-sm font-medium text-gray-700">
                      Pilih Group
                      {groupRecipientCount > 0 && (
                        <span className="ml-2 rounded-full bg-green-100 px-2 py-0.5 text-xs font-semibold text-green-700">{groupRecipientCount} kontak terpilih</span>
                      )}
                    </span>
                    {selectedGroups.length > 0 && (
                      <button onClick={() => setSelectedGroups([])} className="text-xs text-gray-400 hover:text-red-500">Reset</button>
                    )}
                  </div>
                  <div className="space-y-2">
                    {GROUPS.map((g) => {
                      const selected = selectedGroups.includes(g.id);
                      return (
                        <button key={g.id} type="button" onClick={() => toggleGroup(g.id)}
                          className={`flex w-full items-center gap-3 rounded-xl border-2 px-4 py-3 text-left transition-all ${
                            selected ? "border-green-400 bg-green-50" : "border-gray-100 bg-white hover:border-gray-200"
                          }`}
                        >
                          {selected ? <CheckSquare className="h-5 w-5 flex-shrink-0 text-green-600" /> : <Square className="h-5 w-5 flex-shrink-0 text-gray-300" />}
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
                  </div>
                  {selectedGroups.length > 0 && (
                    <div className="mt-3 rounded-lg bg-gray-50 p-3">
                      <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-gray-400">Preview Penerima</p>
                      <div className="space-y-1 max-h-32 overflow-y-auto">
                        {GROUPS.filter((g) => selectedGroups.includes(g.id)).flatMap((g) => g.members).map((m, i) => (
                          <div key={i} className="flex items-center justify-between text-xs">
                            <span className="text-gray-700">{m.name}</span>
                            <span className="font-mono text-gray-400">{m.phone}</span>
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
                </>
              )}
            </div>

            {/* Message */}
            <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
              <div className="mb-3 flex items-center justify-between">
                <label className="text-sm font-medium text-gray-700">Isi Pesan</label>
                <div className="relative">
                  <select
                    value={template}
                    onChange={(e) => setTemplate(e.target.value)}
                    className="h-8 appearance-none rounded-lg border border-gray-200 bg-gray-50 px-3 pr-7 text-xs text-gray-600 focus:outline-none focus:ring-2 focus:ring-green-500"
                  >
                    {TEMPLATES.map((t) => <option key={t}>{t}</option>)}
                  </select>
                  <ChevronDown className="pointer-events-none absolute right-2 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-gray-400" />
                </div>
              </div>
              <textarea
                rows={5}
                placeholder={"Tulis pesan broadcast di sini...\n\nGunakan {{name}} untuk nama penerima."}
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                className="w-full resize-none rounded-lg border border-gray-300 p-3 text-sm placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-green-500"
              />
              <p className="mt-1 text-right text-xs text-gray-400">{message.length} karakter</p>
            </div>

            {/* Delay */}
            <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
              <label className="mb-1 block text-sm font-medium text-gray-700">Delay Antar Pesan (detik)</label>
              <div className="flex items-center gap-3">
                <input
                  type="number" min="1" max="60" value={delay}
                  onChange={(e) => setDelay(e.target.value)}
                  className="h-10 w-24 rounded-lg border border-gray-300 px-3 text-sm focus:outline-none focus:ring-2 focus:ring-green-500"
                />
                <p className="flex items-center gap-1 text-xs text-gray-500">
                  <Info className="h-3.5 w-3.5" /> Disarankan 3–10 detik untuk menghindari blokir
                </p>
              </div>
            </div>

            <button
              onClick={handleSend}
              disabled={!canSend || sending}
              className="w-full rounded-xl bg-green-600 py-3 text-sm font-semibold text-white shadow-sm hover:bg-green-700 disabled:cursor-not-allowed disabled:bg-gray-200 disabled:text-gray-400 transition-colors"
            >
              {sending ? (
                <><Loader2 className="mr-2 inline h-4 w-4 animate-spin" />Mengirim ke {totalCount} penerima...</>
              ) : (
                <><Send className="mr-2 inline h-4 w-4" />Kirim ke {totalCount} Penerima</>
              )}
            </button>
          </div>

          {/* WA Preview col */}
          <div className="lg:col-span-2">
            <p className="mb-2 text-xs font-medium uppercase tracking-wider text-gray-400">Preview Pesan</p>
            <div className="h-[420px] sticky top-6">
              <WABubble message={message} />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
