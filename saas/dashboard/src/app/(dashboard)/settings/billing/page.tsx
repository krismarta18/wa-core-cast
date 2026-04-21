"use client";

import { useEffect, useState } from "react";
import { ArrowUpRight, Calendar, CreditCard, Crown, Download, Receipt, TrendingUp } from "lucide-react";

import { getApiErrorMessage } from "@/lib/api-error";
import { billingApi } from "@/lib/api";
import type { BillingOverview } from "@/lib/types";

function formatCurrency(amount: number) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(amount);
}

function formatDate(date?: string | null) {
  if (!date) {
    return "-";
  }

  return new Intl.DateTimeFormat("id-ID", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  }).format(new Date(date));
}

function formatCycle(cycle: string) {
  return cycle === "yearly" ? "per tahun" : "per bulan";
}

function percentage(used: number, limit: number) {
  if (limit <= 0) {
    return 0;
  }

  return Math.min(100, Math.round((used / limit) * 100));
}

export default function BillingPage() {
  const [billing, setBilling] = useState<BillingOverview | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [checkoutMessage, setCheckoutMessage] = useState<string | null>(null);
  const [submittingPlanId, setSubmittingPlanId] = useState<string | null>(null);

  useEffect(() => {
    let active = true;

    async function loadBilling() {
      try {
        setLoading(true);
        setError(null);
        setCheckoutMessage(null);
        const response = await billingApi.overview();
        if (!active) {
          return;
        }
        setBilling(response.billing);
      } catch (err) {
        if (!active) {
          return;
        }
        setError(getApiErrorMessage(err, "Gagal memuat data billing"));
      } finally {
        if (active) {
          setLoading(false);
        }
      }
    }

    void loadBilling();

    return () => {
      active = false;
    };
  }, []);

  async function handleCheckout(planId: string) {
    try {
      setSubmittingPlanId(planId);
      setError(null);
      const checkoutResponse = await billingApi.checkout({ plan_id: planId });
      const overviewResponse = await billingApi.overview();
      setBilling(overviewResponse.billing);
      setCheckoutMessage(checkoutResponse.message ?? "Paket berhasil diaktifkan.");
    } catch (err) {
      setError(getApiErrorMessage(err, "Gagal memproses pembayaran dummy"));
    } finally {
      setSubmittingPlanId(null);
    }
  }

  const plan = billing?.current_plan ?? null;
  const usageHistory = billing?.usage_history ?? [];
  const plans = billing?.plans ?? [];
  const invoices = billing?.invoices ?? [];
  const maxSent = Math.max(1, ...usageHistory.map((entry) => entry.sent));
  const quotaPct = percentage(plan?.quota_used ?? 0, plan?.quota_limit ?? 0);
  const devicePct = percentage(plan?.device_used ?? 0, plan?.device_max ?? 0);

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <h1 className="text-xl font-bold text-gray-900">Billing & Kuota</h1>
        <p className="text-sm text-gray-500">Detail paket, penggunaan, dan riwayat tagihan</p>
      </div>

      <div className="p-6 space-y-5">
        {error ? (
          <div className="rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
            {error}
          </div>
        ) : null}

        {checkoutMessage ? (
          <div className="rounded-xl border border-green-200 bg-green-50 px-4 py-3 text-sm text-green-700">
            {checkoutMessage}
          </div>
        ) : null}

        {/* Current plan */}
        <div className="rounded-xl border border-green-200 bg-gradient-to-br from-green-50 to-white p-6 shadow-sm">
          <div className="flex flex-wrap items-start justify-between gap-4">
            <div className="flex items-center gap-3">
              <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-green-600">
                <Crown className="h-6 w-6 text-white" />
              </div>
              <div>
                <p className="text-xs font-medium text-green-600 uppercase tracking-wide">Paket Aktif</p>
                <p className="text-xl font-bold text-gray-900">{plan?.name ?? (loading ? "Memuat..." : "Belum ada paket aktif")}</p>
              </div>
            </div>
            <div className="text-right">
              <p className="text-2xl font-bold text-gray-900">{plan ? formatCurrency(plan.price) : "-"}</p>
              <p className="text-sm text-gray-500">{plan ? formatCycle(plan.billing_cycle) : "-"}</p>
            </div>
          </div>
          <div className="mt-4 flex items-center gap-2 text-sm text-gray-500">
            <Calendar className="h-4 w-4" />
            Perpanjang otomatis: <span className="font-medium text-gray-800">{formatDate(plan?.renewal_date)}</span>
          </div>
        </div>

        {/* Quota bars */}
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
            <div className="flex items-center justify-between mb-2">
              <p className="text-sm font-medium text-gray-700">Kuota Pesan</p>
              <p className="text-sm font-bold text-gray-900">{(plan?.quota_used ?? 0).toLocaleString("id")} / {(plan?.quota_limit ?? 0).toLocaleString("id")}</p>
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
              <p className="text-sm font-bold text-gray-900">{plan?.device_used ?? 0} / {plan?.device_max ?? 0}</p>
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
            {usageHistory.map((h) => (
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
          {!loading && usageHistory.length === 0 ? (
            <p className="text-sm text-gray-400">Belum ada data penggunaan untuk 7 hari terakhir.</p>
          ) : null}
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
            {plans.map((p) => (
              <div
                key={p.id}
                className={`rounded-xl border-2 p-4 ${p.current ? "border-green-400 bg-green-50" : "border-gray-100"}`}
              >
                {p.current && (
                  <span className="mb-2 inline-block rounded-full bg-green-600 px-2.5 py-0.5 text-xs font-semibold text-white">
                    Paket Kamu
                  </span>
                )}
                <p className="font-bold text-gray-900">{p.name}</p>
                <p className="mt-1 text-xl font-bold text-gray-800">{formatCurrency(p.price)}<span className="text-sm font-normal text-gray-400">/bln</span></p>
                <ul className="mt-3 space-y-1 text-xs text-gray-500">
                  <li>✓ {p.quota_limit > 0 ? `${p.quota_limit.toLocaleString("id")} pesan/hari` : "Kuota fleksibel"}</li>
                  <li>✓ {p.device_max > 0 ? `${p.device_max} device` : "Device fleksibel"}</li>
                </ul>
                {!p.current && (
                  <button
                    className="mt-4 flex w-full items-center justify-center gap-1 rounded-lg border border-gray-200 py-1.5 text-xs text-gray-600 hover:border-green-400 hover:text-green-600 disabled:cursor-not-allowed disabled:opacity-60"
                    onClick={() => void handleCheckout(p.id)}
                    type="button"
                    disabled={submittingPlanId !== null}
                  >
                    {submittingPlanId === p.id ? "Memproses..." : "Bayar Dummy & Aktifkan"} <ArrowUpRight className="h-3 w-3" />
                  </button>
                )}
              </div>
            ))}
          </div>
          {!loading && plans.length === 0 ? (
            <p className="text-sm text-gray-400">Belum ada paket billing yang tersedia.</p>
          ) : null}
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
                {invoices.map((inv) => (
                  <tr key={inv.id} className="hover:bg-gray-50">
                    <td className="px-5 py-3 font-mono text-xs font-medium text-gray-700">{inv.id}</td>
                    <td className="px-5 py-3 text-gray-600">{formatDate(inv.date)}</td>
                    <td className="px-5 py-3 text-gray-600">{inv.plan_name}</td>
                    <td className="px-5 py-3 font-semibold text-gray-900">{formatCurrency(inv.amount)}</td>
                    <td className="px-5 py-3">
                      <span className={`inline-flex items-center gap-1 rounded-full px-2.5 py-1 text-xs font-medium ${inv.status === "active" ? "bg-green-50 text-green-700" : "bg-gray-100 text-gray-600"}`}>
                        {inv.status === "active" ? "✓ Aktif" : inv.status}
                      </span>
                    </td>
                    <td className="px-5 py-3 text-right">
                      <button className="inline-flex items-center gap-1 rounded-lg border border-gray-200 px-2.5 py-1.5 text-xs text-gray-500 hover:border-green-400 hover:text-green-600" type="button">
                        <Download className="h-3 w-3" /> PDF
                      </button>
                    </td>
                  </tr>
                ))}
                {!loading && invoices.length === 0 ? (
                  <tr>
                    <td className="px-5 py-6 text-center text-sm text-gray-400" colSpan={6}>
                      Belum ada riwayat invoice.
                    </td>
                  </tr>
                ) : null}
              </tbody>
            </table>
          </div>
        </div>

      </div>
    </div>
  );
}
