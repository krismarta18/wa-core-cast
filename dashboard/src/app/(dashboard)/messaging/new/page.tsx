"use client";

import { useState } from "react";
import { Send, Paperclip, ChevronDown, CheckCheck, Smile } from "lucide-react";

const DEVICES = ["Device Utama", "Device CS", "Device Marketing", "Device Backup"];

function WABubble({ from, message }: { from: string; message: string }) {
  const now = new Date();
  const timeStr = now.getHours().toString().padStart(2, "0") + ":" + now.getMinutes().toString().padStart(2, "0");

  return (
    <div className="flex h-full flex-col rounded-xl border border-gray-200 bg-white shadow-sm overflow-hidden">
      {/* WA chat header */}
      <div className="flex items-center gap-3 bg-[#075E54] px-4 py-3">
        <div className="flex h-9 w-9 items-center justify-center rounded-full bg-white/20 text-sm font-bold text-white">
          {from.charAt(0).toUpperCase()}
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
  const [to, setTo] = useState("");
  const [message, setMessage] = useState("");
  const [device, setDevice] = useState(DEVICES[0]);
  const [sent, setSent] = useState(false);

  const isValid = /^[0-9]{9,15}$/.test(to) && message.trim().length > 0;

  const handleSend = () => {
    if (!isValid) return;
    setSent(true);
    setTimeout(() => { setSent(false); setTo(""); setMessage(""); }, 2500);
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
            <h2 className="mt-4 text-lg font-bold text-green-800">Pesan Terkirim!</h2>
            <p className="mt-1 text-sm text-gray-600">Pesan sedang dikirim ke <strong>{to}</strong></p>
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
                    className="h-10 w-full appearance-none rounded-lg border border-gray-300 px-3 pr-8 text-sm text-gray-700 focus:outline-none focus:ring-2 focus:ring-green-500"
                  >
                    {DEVICES.map((d) => <option key={d}>{d}</option>)}
                  </select>
                  <ChevronDown className="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
                </div>
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
                <label className="mb-1 block text-sm font-medium text-gray-700">
                  Pesan <span className="text-red-500">*</span>
                </label>
                <textarea
                  rows={6}
                  placeholder="Tulis pesan di sini..."
                  value={message}
                  onChange={(e) => setMessage(e.target.value)}
                  className="w-full resize-none rounded-lg border border-gray-300 p-3 text-sm placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-green-500"
                />
                <div className="mt-2 flex items-center justify-between">
                  <button className="inline-flex items-center gap-1.5 text-xs text-gray-400 hover:text-gray-600">
                    <Paperclip className="h-3.5 w-3.5" /> Lampiran (segera hadir)
                  </button>
                  <span className="text-xs text-gray-400">{message.length} karakter</span>
                </div>
              </div>

              <button
                onClick={handleSend}
                disabled={!isValid}
                className="w-full rounded-xl bg-green-600 py-3 text-sm font-semibold text-white shadow-sm hover:bg-green-700 disabled:cursor-not-allowed disabled:bg-gray-200 disabled:text-gray-400 transition-colors"
              >
                <Send className="mr-2 inline h-4 w-4" />
                Kirim Pesan
              </button>
            </div>

            {/* WA Preview */}
            <div className="lg:col-span-2">
              <p className="mb-2 text-xs font-medium uppercase tracking-wider text-gray-400">Preview Pesan</p>
              <div className="h-[420px]">
                <WABubble from={to || "Penerima"} message={message} />
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
