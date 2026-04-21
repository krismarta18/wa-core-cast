"use client";

import { useState } from "react";
import { 
  FileText, 
  KeyRound, 
  Send, 
  Webhook, 
  Terminal, 
  Code2, 
  CheckCircle2, 
  AlertCircle,
  Copy,
  ChevronRight,
  ShieldCheck,
  Zap
} from "lucide-react";
import { useToast } from "@/components/ui/toast";

interface CodeBlockProps {
  code: string;
  language?: string;
  title?: string;
}

function CodeBlock({ code, title }: CodeBlockProps) {
  const { info } = useToast();
  
  const copy = () => {
    navigator.clipboard.writeText(code);
    info("Disalin!", "Contoh kode telah disalin ke clipboard.");
  };

  return (
    <div className="group relative my-4 overflow-hidden rounded-xl border border-gray-200 bg-gray-900 shadow-sm">
      {title && (
        <div className="flex items-center justify-between border-b border-gray-800 bg-gray-800/50 px-4 py-2">
          <span className="text-xs font-mono text-gray-400">{title}</span>
          <button 
            onClick={copy}
            className="flex items-center gap-1.5 text-xs text-gray-400 transition-colors hover:text-white"
          >
            <Copy className="h-3 w-3" />
            Salin
          </button>
        </div>
      )}
      <pre className="overflow-x-auto p-4 text-sm leading-relaxed text-gray-300">
        <code>{code}</code>
      </pre>
    </div>
  );
}

function Endpoint({ method, path }: { method: string; path: string }) {
  const colorMap: Record<string, string> = {
    GET: "bg-blue-100 text-blue-700 border-blue-200",
    POST: "bg-green-100 text-green-700 border-green-200",
    PUT: "bg-amber-100 text-amber-700 border-amber-200",
    DELETE: "bg-red-100 text-red-700 border-red-200",
  };

  return (
    <div className="inline-flex items-center gap-3 font-mono text-sm">
      <span className={`rounded-md border px-2 py-0.5 font-bold ${colorMap[method] || "bg-gray-100"}`}>
        {method}
      </span>
      <span className="text-gray-600">{path}</span>
    </div>
  );
}

export default function APIDocsPage() {
  const [activeTab, setActiveTab] = useState("auth");

  const sections = [
    { id: "auth", label: "Autentikasi", icon: KeyRound },
    { id: "messaging", label: "Messaging API", icon: Send },
    { id: "status", label: "Status & Tracking", icon: FileText },
    { id: "webhooks", label: "Webhooks", icon: Webhook },
  ];

  return (
    <div className="mx-auto max-w-6xl py-8">
      {/* Header */}
      <div className="mb-10 text-center">
        <h1 className="mb-3 text-3xl font-extrabold tracking-tight text-gray-900 md:text-4xl">
          API & Integration <span className="text-green-600">Documentation</span>
        </h1>
        <p className="mx-auto max-w-2xl text-lg text-gray-600">
          Panduan lengkap untuk mengintegrasikan sistem Anda dengan WACAST Core.
          Kirim pesan, terima notifikasi, dan kelola sesi secara otomatis.
        </p>
      </div>

      <div className="flex flex-col gap-8 lg:flex-row">
        {/* Sidebar Nav */}
        <div className="w-full lg:w-64">
          <nav className="sticky top-8 space-y-1">
            {sections.map((s) => (
              <button
                key={s.id}
                onClick={() => setActiveTab(s.id)}
                className={`flex w-full items-center gap-3 rounded-lg px-4 py-3 text-sm font-medium transition-all ${
                  activeTab === s.id
                    ? "bg-green-600 text-white shadow-lg shadow-green-100"
                    : "text-gray-600 hover:bg-gray-100 hover:text-gray-900"
                }`}
              >
                <s.icon className="h-4 w-4" />
                {s.label}
                <ChevronRight className={`ml-auto h-4 w-4 transition-transform ${activeTab === s.id ? "rotate-90" : ""}`} />
              </button>
            ))}
          </nav>

          <div className="mt-8 rounded-2xl bg-blue-50 p-6 text-blue-800">
            <div className="mb-2 flex items-center gap-2 font-bold">
              <Zap className="h-4 w-4" />
              Base URL
            </div>
            <code className="block rounded bg-blue-100 px-2 py-1 font-mono text-sm border border-blue-200">
              http://localhost:8080
            </code>
          </div>
        </div>

        {/* Content Area */}
        <div className="flex-1 space-y-12 pb-20">
          {activeTab === "auth" && (
            <section className="animate-in fade-in slide-in-from-bottom-4 duration-500">
              <h2 className="mb-6 flex items-center gap-3 text-2xl font-bold text-gray-900">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-orange-100 text-orange-600">
                  <ShieldCheck className="h-6 w-6" />
                </div>
                Autentikasi (API Key)
              </h2>
              <div className="prose prose-green max-w-none text-gray-600">
                <p>
                  Seluruh endpoint API dilindungi oleh autentikasi menggunakan <strong>API Key</strong>. 
                  Anda dapat membuat dan mengelola kunci ini di halaman <a href="/api-integration/keys" className="font-semibold text-green-600 underline">API Keys</a>.
                </p>
                <div className="rounded-xl border border-gray-100 bg-gray-50 p-6 my-6">
                  <h4 className="mt-0 font-bold text-gray-900">Petunjuk Penggunaan Header</h4>
                  <p className="mb-4">Sisipkan API Key Anda ke dalam header HTTP berikut:</p>
                  <div className="grid gap-4 sm:grid-cols-2">
                    <div className="rounded-lg bg-white p-4 border border-gray-200">
                      <span className="text-xs font-bold text-gray-400 uppercase">Header Name</span>
                      <p className="font-mono font-bold text-green-700">X-API-Key</p>
                    </div>
                    <div className="rounded-lg bg-white p-4 border border-gray-200">
                      <span className="text-xs font-bold text-gray-400 uppercase">Format</span>
                      <p className="font-mono text-gray-700 text-sm">wck_live_xxx...</p>
                    </div>
                  </div>
                </div>
                <div className="flex items-start gap-4 rounded-xl bg-amber-50 p-5 border border-amber-100">
                  <AlertCircle className="mt-1 h-5 w-5 flex-shrink-0 text-amber-500" />
                  <p className="text-sm text-amber-800">
                    <strong>Peringatan Keamanan:</strong> Jangan memberikan API Key Anda kepada siapapun. 
                    Jika kunci Anda bocor, segera hapus dan buat yang baru di dashboard.
                  </p>
                </div>
              </div>
            </section>
          )}

          {activeTab === "messaging" && (
            <section className="animate-in fade-in slide-in-from-bottom-4 duration-500">
              <h2 className="mb-6 flex items-center gap-3 text-2xl font-bold text-gray-900">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-green-100 text-green-600">
                  <Send className="h-6 w-6" />
                </div>
                Messaging API
              </h2>
              <div className="space-y-10">
                {/* Send Text */}
                <div className="rounded-2xl border border-gray-200 bg-white p-6 shadow-sm underline-offset-4">
                  <h3 className="mb-4 text-lg font-bold">1. Mengirim Pesan Teks</h3>
                  <Endpoint method="POST" path="/api/v1/devices/{device_id}/messages" />
                  <div className="mt-6">
                    <h4 className="mb-3 text-sm font-bold text-gray-400 uppercase">Request Body (JSON)</h4>
                    <CodeBlock 
                      title="Payload JSON"
                      code={`{
  "target_jid": "628123456789@s.whatsapp.net",
  "content": "Halo! Ini adalah pesan otomatis dari WACAST Core."
}`} 
                    />
                  </div>
                  <div className="mt-6">
                    <h4 className="mb-3 text-sm font-bold text-gray-400 uppercase">Contoh Permintaan (cURL)</h4>
                    <CodeBlock 
                      code={`curl -X POST http://localhost:8080/api/v1/devices/your-device-id/messages \\
  -H "X-API-Key: wck_live_your_key_here" \\
  -H "Content-Type: application/json" \\
  -d '{"target_jid": "628xxx@s.whatsapp.net", "content": "Hello World"}'`} 
                    />
                  </div>
                </div>

                {/* Send Media */}
                <div className="rounded-2xl border border-gray-200 bg-white p-6 shadow-sm">
                  <h3 className="mb-4 text-lg font-bold">2. Mengirim Pesan Media (Gambar/File)</h3>
                  <Endpoint method="POST" path="/api/v1/devices/{device_id}/messages/media" />
                  <p className="mt-4 text-sm text-gray-600 italic">Gunakan format <code>multipart/form-data</code> untuk mengirim file langsung.</p>
                  <div className="mt-6 overflow-hidden rounded-lg border border-gray-200 bg-gray-50">
                    <table className="w-full text-left text-sm">
                      <thead className="bg-gray-100/50 text-gray-500">
                        <tr>
                          <th className="px-4 py-2 font-bold uppercase tracking-wider">Parameter</th>
                          <th className="px-4 py-2 font-bold uppercase tracking-wider">Tipe</th>
                          <th className="px-4 py-2 font-bold uppercase tracking-wider">Deskripsi</th>
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-gray-200">
                        <tr>
                          <td className="px-4 py-3 font-mono text-green-700">target_jid</td>
                          <td className="px-4 py-3 text-gray-500">string</td>
                          <td className="px-4 py-3">JID tujuan (misal 628...@s.whatsapp.net)</td>
                        </tr>
                        <tr>
                          <td className="px-4 py-3 font-mono text-green-700">content_type</td>
                          <td className="px-4 py-3 text-gray-500">string</td>
                          <td className="px-4 py-3">image | document | audio | video</td>
                        </tr>
                        <tr>
                          <td className="px-4 py-3 font-mono text-green-700">file</td>
                          <td className="px-4 py-3 text-gray-500">binary</td>
                          <td className="px-4 py-3">File fisik yang akan dikirim (Max 50MB)</td>
                        </tr>
                        <tr>
                          <td className="px-4 py-3 font-mono text-green-700">caption</td>
                          <td className="px-4 py-3 text-gray-500">string</td>
                          <td className="px-4 py-3">Optional. Keterangan pada media</td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>
            </section>
          )}

          {activeTab === "status" && (
            <section className="animate-in fade-in slide-in-from-bottom-4 duration-500">
              <h2 className="mb-6 flex items-center gap-3 text-2xl font-bold text-gray-900">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-blue-100 text-blue-600">
                  <FileText className="h-6 w-6" />
                </div>
                Status & Tracking
              </h2>
              <div className="prose prose-blue max-w-none text-gray-600">
                <p>
                  Setiap pengiriman pesan akan mengembalikan <code>message_id</code>. Anda dapat menggunakan ID tersebut untuk menelusuri status pengiriman di server.
                </p>
                <div className="mt-8 rounded-2xl border border-gray-200 bg-white p-6 shadow-sm">
                  <h3 className="mb-4 text-lg font-bold italic">Pengecekan Status Pesan</h3>
                  <Endpoint method="GET" path="/api/v1/messages/{message_id}/status" />
                  <CodeBlock 
                    title="Respon Berhasil"
                    code={`{
  "message_id": "WACAST-1712345678",
  "status": "delivered",
  "timestamp": 1713695420
}`} 
                  />
                  <div className="mt-4 grid grid-cols-2 gap-3 sm:grid-cols-5">
                    {["pending", "sent", "delivered", "read", "failed"].map((s) => (
                      <div key={s} className="flex flex-col items-center justify-center rounded-lg border border-gray-100 bg-gray-50 py-3">
                        <CheckCircle2 className={`h-4 w-4 mb-1 ${s === "failed" ? "text-red-500" : "text-green-500"}`} />
                        <span className="text-[10px] font-bold uppercase text-gray-500">{s}</span>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </section>
          )}

          {activeTab === "webhooks" && (
            <section className="animate-in fade-in slide-in-from-bottom-4 duration-500">
              <h2 className="mb-6 flex items-center gap-3 text-2xl font-bold text-gray-900">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-indigo-100 text-indigo-600">
                  <Webhook className="h-6 w-6" />
                </div>
                Webhook Configuration
              </h2>
              <div className="space-y-8 prose prose-indigo max-w-none text-gray-600">
                <p>
                  Webhook memungkinkan aplikasi Anda menerima notifikasi secara <em>real-time</em> saat pesan diterima atau status berubah.
                </p>

                <div className="rounded-xl border border-indigo-100 bg-indigo-50/50 p-6">
                  <h4 className="mt-0 font-bold text-indigo-900 flex items-center gap-2">
                    <Terminal className="h-4 w-4" />
                    WACAST Signature Verification
                  </h4>
                  <p className="text-sm">
                    Untuk menjamin data berasal dari WACAST, verifikasi header <code>X-WACAST-Signature</code>. 
                    Signature dibuat menggunakan <strong>HMAC-SHA256</strong> dengan menggunakan <strong>Webhook Secret</strong> Anda sebagai kuncinya.
                  </p>
                  <CodeBlock 
                    title="Contoh Verifikasi (Node.js)"
                    code={`const crypto = require('crypto');

const secret = 'YOUR_WEBHOOK_SECRET';
const signature = req.headers['x-wacast-signature'];
const body = JSON.stringify(req.body);

const expected = crypto
  .createHmac('sha256', secret)
  .update(body)
  .digest('hex');

if (signature === expected) {
  console.log('Valid Payload');
}`} 
                  />
                </div>

                <div>
                  <h3 className="font-bold text-gray-900">Format Payload Webhook</h3>
                  <CodeBlock 
                    title="Event: message.received"
                    code={`{
  "event": "message.received",
  "timestamp": 1713695420,
  "data": {
    "message_id": "WA-88223344",
    "target_jid": "628... @s.whatsapp.net",
    "content": "Terima kasih infonya!",
    "status": "read"
  }
}`} 
                  />
                </div>
              </div>
            </section>
          )}
        </div>
      </div>
    </div>
  );
}
