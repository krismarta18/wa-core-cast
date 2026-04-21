"use client";

import { useState } from "react";
import {
  QrCode,
  Send,
  KeyRound,
  CheckCircle2,
  ChevronRight,
  ArrowRight,
  Smartphone,
  Webhook,
  Book,
} from "lucide-react";
import Link from "next/link";

interface Step {
  id: number;
  title: string;
  description: string;
  icon: React.ElementType;
  cta: string;
  href: string;
  detail: string[];
  color: string;
}

const STEPS: Step[] = [
  {
    id: 1,
    title: "Hubungkan Device WhatsApp",
    description: "Scan QR code untuk menghubungkan nomor WA pertama kamu ke WACAST.",
    icon: QrCode,
    cta: "Buka Halaman QR",
    href: "/devices/qr",
    detail: [
      "Buka WhatsApp di HP kamu",
      "Pilih Perangkat Tertaut → Tautkan Perangkat",
      "Arahkan kamera ke QR code yang muncul",
    ],
    color: "green",
  },
  {
    id: 2,
    title: "Kirim Pesan Test",
    description: "Pastikan koneksi berjalan dengan mengirim pesan sederhana ke nomor kamu sendiri.",
    icon: Send,
    cta: "Kirim Pesan Test",
    href: "/messaging/new",
    detail: [
      "Pilih device yang sudah terhubung",
      "Masukkan nomor tujuan (bisa nomor sendiri)",
      "Tulis pesan singkat dan tekan Kirim",
    ],
    color: "blue",
  },
  {
    id: 3,
    title: "Buat API Key",
    description: "Generate API key untuk mengintegrasikan WACAST ke aplikasi atau sistem kamu.",
    icon: KeyRound,
    cta: "Buat API Key",
    href: "/api-integration/keys",
    detail: [
      "Klik tombol + Buat API Key",
      "Beri nama yang mudah diingat (mis: Production)",
      "Salin dan simpan key dengan aman",
    ],
    color: "purple",
  },
  {
    id: 4,
    title: "Atur Webhook (Opsional)",
    description: "Terima notifikasi real-time ke server kamu setiap kali ada pesan masuk atau status berubah.",
    icon: Webhook,
    cta: "Konfigurasi Webhook",
    href: "/api-integration/webhooks",
    detail: [
      "Masukkan URL endpoint server kamu",
      "Pilih event yang ingin kamu terima",
      "Simpan dan lakukan uji koneksi",
    ],
    color: "orange",
  },
];

const COLOR_MAP: Record<string, { bg: string; text: string; ring: string; badge: string; btn: string }> = {
  green:  { bg: "bg-green-50",  text: "text-green-700",  ring: "ring-green-200",  badge: "bg-green-100 text-green-700",  btn: "bg-green-600 hover:bg-green-700" },
  blue:   { bg: "bg-blue-50",   text: "text-blue-700",   ring: "ring-blue-200",   badge: "bg-blue-100 text-blue-700",    btn: "bg-blue-600 hover:bg-blue-700"   },
  purple: { bg: "bg-purple-50", text: "text-purple-700", ring: "ring-purple-200", badge: "bg-purple-100 text-purple-700",btn: "bg-purple-600 hover:bg-purple-700"},
  orange: { bg: "bg-orange-50", text: "text-orange-700", ring: "ring-orange-200", badge: "bg-orange-100 text-orange-700",btn: "bg-orange-600 hover:bg-orange-700"},
};

export default function OnboardingPage() {
  const [completed, setCompleted] = useState<Set<number>>(new Set());

  function toggle(id: number) {
    setCompleted((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  }

  const progress = Math.round((completed.size / STEPS.length) * 100);
  const allDone = completed.size === STEPS.length;

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <h1 className="text-xl font-bold text-gray-900">Quick Start</h1>
        <p className="text-sm text-gray-500">Ikuti langkah-langkah berikut untuk mulai menggunakan WACAST</p>
      </div>

      <div className="p-6 max-w-3xl mx-auto space-y-6">
        {/* Progress card */}
        <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
          <div className="mb-3 flex items-center justify-between">
            <div>
              <p className="text-sm font-semibold text-gray-700">Progres Setup</p>
              <p className="text-xs text-gray-400">{completed.size} dari {STEPS.length} langkah selesai</p>
            </div>
            <span className="text-2xl font-bold text-green-600">{progress}%</span>
          </div>
          <div className="h-3 w-full overflow-hidden rounded-full bg-gray-100">
            <div
              className="h-full rounded-full bg-green-500 transition-all duration-500"
              style={{ width: `${progress}%` }}
            />
          </div>
          {allDone && (
            <div className="mt-4 flex items-center gap-3 rounded-lg bg-green-50 px-4 py-3">
              <CheckCircle2 className="h-5 w-5 flex-shrink-0 text-green-600" />
              <p className="text-sm font-medium text-green-800">
                Selamat! Setup kamu sudah lengkap. WACAST siap digunakan.
              </p>
            </div>
          )}
        </div>

        {/* Steps */}
        {STEPS.map((step, i) => {
          const c = COLOR_MAP[step.color];
          const done = completed.has(step.id);
          return (
            <div
              key={step.id}
              className={`rounded-xl border bg-white shadow-sm transition-all ${
                done ? "border-green-200 opacity-70" : "border-gray-200"
              }`}
            >
              <div className="p-5">
                <div className="flex items-start gap-4">
                  {/* Step number / check */}
                  <div
                    className={`mt-0.5 flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-xl ${
                      done ? "bg-green-100" : c.bg
                    }`}
                  >
                    {done ? (
                      <CheckCircle2 className="h-5 w-5 text-green-600" />
                    ) : (
                      <step.icon className={`h-5 w-5 ${c.text}`} />
                    )}
                  </div>

                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 flex-wrap">
                      <span className={`rounded-full px-2 py-0.5 text-[11px] font-semibold ${c.badge}`}>
                        Langkah {step.id}
                      </span>
                      {i === 0 && !done && (
                        <span className="rounded-full bg-orange-100 px-2 py-0.5 text-[11px] font-semibold text-orange-700">
                          Mulai di sini
                        </span>
                      )}
                    </div>
                    <h3 className={`mt-1.5 text-base font-semibold ${done ? "line-through text-gray-400" : "text-gray-900"}`}>
                      {step.title}
                    </h3>
                    <p className="mt-0.5 text-sm text-gray-500">{step.description}</p>

                    {/* Sub-steps */}
                    {!done && (
                      <ol className="mt-3 space-y-1.5 pl-1">
                        {step.detail.map((d, di) => (
                          <li key={di} className="flex items-start gap-2 text-xs text-gray-500">
                            <span className="mt-0.5 flex h-4 w-4 flex-shrink-0 items-center justify-center rounded-full bg-gray-100 text-[10px] font-bold text-gray-500">
                              {di + 1}
                            </span>
                            {d}
                          </li>
                        ))}
                      </ol>
                    )}
                  </div>
                </div>

                {/* Actions */}
                <div className="mt-4 flex items-center gap-3 pt-3 border-t border-gray-50">
                  {!done ? (
                    <>
                      <Link
                        href={step.href}
                        className={`inline-flex items-center gap-2 rounded-lg px-4 py-2 text-sm font-semibold text-white transition-colors ${c.btn}`}
                      >
                        {step.cta}
                        <ArrowRight className="h-3.5 w-3.5" />
                      </Link>
                      <button
                        onClick={() => toggle(step.id)}
                        className="text-sm text-gray-400 hover:text-gray-600 hover:underline"
                      >
                        Tandai selesai
                      </button>
                    </>
                  ) : (
                    <button
                      onClick={() => toggle(step.id)}
                      className="text-sm text-gray-400 hover:text-gray-600 hover:underline"
                    >
                      Batalkan
                    </button>
                  )}
                </div>
              </div>
            </div>
          );
        })}

        {/* Docs link */}
        <div className="rounded-xl border border-dashed border-gray-200 bg-white p-5">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-gray-100">
              <Book className="h-5 w-5 text-gray-500" />
            </div>
            <div className="flex-1">
              <p className="text-sm font-semibold text-gray-800">Butuh bantuan lebih?</p>
              <p className="text-xs text-gray-500">Baca dokumentasi lengkap atau hubungi support kami.</p>
            </div>
            <a
              href="https://docs.wacast.id"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-1.5 rounded-lg border border-gray-200 px-3 py-1.5 text-sm font-medium text-gray-600 hover:bg-gray-50"
            >
              Dokumentasi <ChevronRight className="h-3.5 w-3.5" />
            </a>
          </div>
        </div>
      </div>
    </div>
  );
}
