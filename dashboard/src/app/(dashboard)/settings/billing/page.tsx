"use client";

import { Crown, TrendingUp, Calendar, CreditCard, ArrowUpRight, Download, Receipt } from "lucide-react";

const PLAN = {
  name: "Business Pro",
  price: "Rp 299.000",
  cycle: "per bulan",
  renewDate: "12 Mei 2026",
  quotaUsed: 8_420,
  quotaLimit: 10_000,
  deviceUsed: 4,
  deviceMax: 10,
};

const USAGE_HISTORY = [
  { date: "12 Apr", sent: 1243, failed: 17 },
  { date: "11 Apr", sent: 980, failed: 8 },
  { date: "10 Apr", sent: 1102, failed: 22 },
  { date: "09 Apr", sent: 876, failed: 5 },
  { date: "08 Apr", sent: 1340, failed: 11 },
  { date: "07 Apr", sent: 490, failed: 3 },
  { date: "06 Apr", sent: 389, failed: 2 },
];

const PLANS = [
  { name: "Starter", price: "Rp 99.000", quota: "2.000 pesan", devices: "2 device", current: false },
  { name: "Business Pro", price: "Rp 299.000", quota: "10.000 pesan", devices: "10 device", current: true },
  { name: "Enterprise", price: "Rp 799.000", quota: "Unlimited", devices: "Unlimited", current: false },
];

const INVOICES = [
  { id: "INV-2026-04", date: "01 Apr 2026", plan: "Business Pro", amount: "Rp 299.000" },
  { id: "INV-2026-03", date: "01 Mar 2026", plan: "Business Pro", amount: "Rp 299.000" },
  { id: "INV-2026-02", date: "01 Feb 2026", plan: "Business Pro", amount: "Rp 299.000" },
  { id: "INV-2026-01", date: "01 Jan 2026", plan: "Starter",      amount: "Rp 99.000" },
  { id: "INV-2025-12", date: "01 Des 2025", plan: "Starter",      amount: "Rp 99.000" },
];

const maxSent = Math.max(...USAGE_HISTORY.map((h) => h.sent));

export default function BillingPage() {
  const quotaPct = Math.round((PLAN.quotaUsed / PLAN.quotaLimit) * 100);
  const devicePct = Math.round((PLAN.deviceUsed / PLAN.deviceMax) * 100);

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <h1 className="text-xl font-bold text-gray-900">Billing & Kuota</h1>
        <p className="text-sm text-gray-500">Detail paket, penggunaan, dan riwayat tagihan</p>
      </div>

      <div className="p-6 space-y-5">
        {/* Current plan */}
        <div className="rounded-xl border border-green-200 bg-gradient-to-br from-green-50 to-white p-6 shadow-sm">
          <div className="flex flex-wrap items-start justify-between gap-4">
            <div className="flex items-center gap-3">
              <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-green-600">
                <Crown className="h-6 w-6 text-white" />
              </div>
              <div>
                <p className="text-xs font-medium text-green-600 uppercase tracking-wide">Paket Aktif</p>
                <p className="text-xl font-bold text-gray-900">{PLAN.name}</p>
              </div>
            </div>
            <div className="text-right">
              <p className="text-2xl font-bold text-gray-900">{PLAN.price}</p>
              <p className="text-sm text-gray-500">{PLAN.cycle}</p>
            </div>
          </div>
          <div className="mt-4 flex items-center gap-2 text-sm text-gray-500">
            <Calendar className="h-4 w-4" />
            Perpanjang otomatis: <span className="font-medium text-gray-800">{PLAN.renewDate}</span>
          </div>
        </div>

        {/* Quota bars */}
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
            <div className="flex items-center justify-between mb-2">
              <p className="text-sm font-medium text-gray-700">Kuota Pesan</p>
              <p className="text-sm font-bold text-gray-900">{PLAN.quotaUsed.toLocaleString("id")} / {PLAN.quotaLimit.toLocaleString("id")}</p>
            </div>
            <div className="h-2.5 w-full overflow-hidden rounded-full bg-gray-100">
              <div
                className={`h-full rounded-full transition-all ${quotaPct >= 90 ? "bg-red-500" : quotaPct >= 75 ? "bg-yellow-500" : "bg-green-500"}`}
                style={{ width: `${quotaPct}%` }}
              />
            </div>
            <p className={`mt-1.5 text-xs ${quotaPct >= 90 ? "text-red-500 font-semibold" : "text-gray-400"}`}>
              {quotaPct}% terpakai {quotaPct >= 90 && "— hampir habis!"}
            </p>
          </div>

          <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
            <div className="flex items-center justify-between mb-2">
              <p className="text-sm font-medium text-gray-700">Slot Device</p>
              <p className="text-sm font-bold text-gray-900">{PLAN.deviceUsed} / {PLAN.deviceMax}</p>
            </div>
            <div className="h-2.5 w-full overflow-hidden rounded-full bg-gray-100">
              <div
                className="h-full rounded-full bg-blue-500 transition-all"
                style={{ width: `${devicePct}%` }}
              />
            </div>
            <p className="mt-1.5 text-xs text-gray-400">{devicePct}% terpakai</p>
          </div>
        </div>

        {/* 7-day usage chart */}
        <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
          <div className="mb-4 flex items-center gap-2">
            <TrendingUp className="h-5 w-5 text-gray-400" />
            <h2 className="font-semibold text-gray-900">Penggunaan 7 Hari Terakhir</h2>
          </div>
          <div className="flex items-end gap-2 h-32">
            {USAGE_HISTORY.map((h) => (
              <div key={h.date} className="flex flex-1 flex-col items-center gap-1">
                <div className="flex w-full flex-col-reverse gap-0.5">
                  <div
                    className="w-full rounded-t bg-green-400"
                    style={{ height: `${Math.round((h.sent / maxSent) * 100)}px` }}
                    title={`${h.sent} terkirim`}
                  />
                </div>
                <p className="text-xs text-gray-400">{h.date}</p>
              </div>
            ))}
          </div>
          <div className="mt-3 flex items-center gap-4 text-xs text-gray-500">
            <span className="flex items-center gap-1"><span className="inline-block h-2 w-2 rounded-sm bg-green-400" /> Terkirim</span>
          </div>
        </div>

        {/* Plan comparison */}
        <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
          <div className="mb-4 flex items-center gap-2">
            <CreditCard className="h-5 w-5 text-gray-400" />
            <h2 className="font-semibold text-gray-900">Pilihan Paket</h2>
          </div>
          <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
            {PLANS.map((p) => (
              <div
                key={p.name}
                className={`rounded-xl border-2 p-4 ${p.current ? "border-green-400 bg-green-50" : "border-gray-100"}`}
              >
                {p.current && (
                  <span className="mb-2 inline-block rounded-full bg-green-600 px-2.5 py-0.5 text-xs font-semibold text-white">
                    Paket Kamu
                  </span>
                )}
                <p className="font-bold text-gray-900">{p.name}</p>
                <p className="mt-1 text-xl font-bold text-gray-800">{p.price}<span className="text-sm font-normal text-gray-400">/bln</span></p>
                <ul className="mt-3 space-y-1 text-xs text-gray-500">
                  <li>✓ {p.quota}</li>
                  <li>✓ {p.devices}</li>
                </ul>
                {!p.current && (
                  <button className="mt-4 flex w-full items-center justify-center gap-1 rounded-lg border border-gray-200 py-1.5 text-xs text-gray-600 hover:border-green-400 hover:text-green-600">
                    Upgrade <ArrowUpRight className="h-3 w-3" />
                  </button>
                )}
              </div>
            ))}
          </div>
        </div>

        {/* Invoice history */}
        <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
          <div className="flex items-center gap-2 border-b border-gray-100 px-5 py-4">
            <Receipt className="h-5 w-5 text-gray-400" />
            <h2 className="font-semibold text-gray-900">Riwayat Invoice</h2>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full min-w-[520px] text-sm">
              <thead>
                <tr className="border-b border-gray-100 bg-gray-50 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">
                  <th className="px-5 py-3">ID Invoice</th>
                  <th className="px-5 py-3">Tanggal</th>
                  <th className="px-5 py-3">Paket</th>
                  <th className="px-5 py-3">Nominal</th>
                  <th className="px-5 py-3">Status</th>
                  <th className="px-5 py-3 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-50">
                {INVOICES.map((inv) => (
                  <tr key={inv.id} className="hover:bg-gray-50">
                    <td className="px-5 py-3 font-mono text-xs font-medium text-gray-700">{inv.id}</td>
                    <td className="px-5 py-3 text-gray-600">{inv.date}</td>
                    <td className="px-5 py-3 text-gray-600">{inv.plan}</td>
                    <td className="px-5 py-3 font-semibold text-gray-900">{inv.amount}</td>
                    <td className="px-5 py-3">
                      <span className="inline-flex items-center gap-1 rounded-full bg-green-50 px-2.5 py-1 text-xs font-medium text-green-700">
                        ✓ Lunas
                      </span>
                    </td>
                    <td className="px-5 py-3 text-right">
                      <button className="inline-flex items-center gap-1 rounded-lg border border-gray-200 px-2.5 py-1.5 text-xs text-gray-500 hover:border-green-400 hover:text-green-600">
                        <Download className="h-3 w-3" /> PDF
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

      </div>
    </div>
  );
}
