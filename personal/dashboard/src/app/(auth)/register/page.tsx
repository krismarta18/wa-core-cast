"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { MessageSquareMore } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { PhoneInput } from "@/components/ui/phone-input";
import { authApi } from "@/lib/api";
import { getApiErrorMessage } from "@/lib/api-error";
import { savePendingAuthState } from "@/lib/auth-session";
import { useAuth } from "@/providers/auth-provider";

export default function RegisterPage() {
  const router = useRouter();
  const { session, hydrated, appConfig } = useAuth();
  const [dialCode, setDialCode] = useState("+62");
  const [phone, setPhone] = useState("");
  const [fullName, setFullName] = useState("");
  const [agree, setAgree] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [errors, setErrors] = useState<{ fullName?: string; phone?: string; agree?: string }>({});

  useEffect(() => {
    if (hydrated && session) {
      router.replace("/");
    }
  }, [hydrated, router, session]);

  const validate = () => {
    const e: { fullName?: string; phone?: string; agree?: string } = {};
    if (!fullName.trim()) e.fullName = "Nama lengkap wajib diisi";
    else if (fullName.trim().length < 2) e.fullName = "Nama terlalu pendek";
    if (!phone) e.phone = "Nomor HP wajib diisi";
    else if (phone.length < 7) e.phone = "Nomor tidak valid";
    if (!agree) e.agree = "Anda harus menyetujui syarat & ketentuan";
    setErrors(e);
    return Object.keys(e).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate()) return;

    const fullPhone = dialCode.replace("+", "") + phone;

    setSubmitting(true);

    try {
      await authApi.register({
        phone_number: fullPhone,
        full_name: fullName.trim(),
      });
      savePendingAuthState({
        phoneNumber: fullPhone,
        context: "register",
        fullName: fullName.trim(),
        rememberMe: true,
      });
      toast.success("Kode OTP registrasi berhasil dikirim");
      router.push(`/otp?phone=${encodeURIComponent(fullPhone)}&context=register`);
    } catch (error) {
      setErrors({ phone: getApiErrorMessage(error, "Gagal memulai registrasi") });
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="flex min-h-screen">
      {/* Left branding panel */}
      <div className="hidden lg:flex lg:w-1/2 flex-col items-center justify-center bg-gradient-to-br from-green-600 to-green-800 px-12 text-white">
        <div className="max-w-sm text-center">
          <div className="mb-6 flex justify-center">
            <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-white/10 backdrop-blur p-2">
              <img src="/icon.png" alt="Logo" className="h-full w-full object-contain" />
            </div>
          </div>
          <h1 className="text-3xl font-bold">{appConfig.app_name || "WACAST"}</h1>
          <p className="mt-3 text-base text-green-100 leading-relaxed">
            Mulai kelola WhatsApp Gateway Anda dalam hitungan menit. Gratis untuk
            dicoba.
          </p>

          <div className="mt-10 grid grid-cols-2 gap-4">
            {[
              { num: "25+", label: "Sesi WhatsApp" },
              { num: "99%", label: "Uptime" },
              { num: "Tak terbatas", label: "Pesan/hari" },
              { num: "24/7", label: "Support" },
            ].map(({ num, label }) => (
              <div
                key={label}
                className="rounded-xl bg-white/10 px-4 py-3 text-center"
              >
                <p className="text-xl font-bold">{num}</p>
                <p className="text-xs text-green-100">{label}</p>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Right form panel */}
      <div className="flex w-full flex-col items-center justify-center bg-gray-50 px-6 py-12 lg:w-1/2">
        {/* Mobile logo */}
        <div className="mb-8 flex items-center gap-3 lg:hidden">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-white shadow-sm border border-gray-100 p-1.5">
            <img src="/icon.png" alt="Logo" className="h-full w-full object-contain" />
          </div>
          <span className="text-xl font-bold text-gray-900">{appConfig.app_name || "WACAST"}</span>
        </div>

        <div className="w-full max-w-sm">
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-gray-900">Buat Akun</h2>
            <p className="mt-1 text-sm text-gray-500">
              Masukkan nomor HP Anda untuk mendaftar
            </p>
          </div>

          <form onSubmit={handleSubmit} noValidate className="space-y-4">
            <div className="flex flex-col gap-1">
              <label htmlFor="full-name" className="text-sm font-medium text-gray-700">
                Nama Lengkap
              </label>
              <input
                id="full-name"
                type="text"
                value={fullName}
                onChange={(event) => setFullName(event.target.value)}
                className="h-11 rounded-xl border border-gray-200 bg-white px-4 text-sm text-gray-900 outline-none transition focus:border-green-500"
                placeholder="Masukkan nama lengkap"
              />
              {errors.fullName && <p className="text-xs text-red-600">{errors.fullName}</p>}
            </div>

            {/* Phone */}
            <div className="flex flex-col gap-1">
              <label htmlFor="reg-phone" className="text-sm font-medium text-gray-700">
                Nomor HP
              </label>
              <PhoneInput
                id="reg-phone"
                dialCode={dialCode}
                onDialCodeChange={setDialCode}
                value={phone}
                onChange={setPhone}
                error={errors.phone}
              />
              {errors.phone && <p className="text-xs text-red-600">{errors.phone}</p>}
            </div>

            {/* Terms */}
            <div className="flex flex-col gap-1">
              <label className="flex items-start gap-2 cursor-pointer select-none">
                <input
                  type="checkbox"
                  checked={agree}
                  onChange={(e) => setAgree(e.target.checked)}
                  className="mt-0.5 h-4 w-4 flex-shrink-0 rounded border-gray-300 accent-green-600"
                />
                <span className="text-sm text-gray-600">
                  Saya menyetujui{" "}
                  <a href="#" className="font-medium text-green-600 hover:underline">
                    Syarat &amp; Ketentuan
                  </a>{" "}
                  serta{" "}
                  <a href="#" className="font-medium text-green-600 hover:underline">
                    Kebijakan Privasi
                  </a>
                </span>
              </label>
              {errors.agree && (
                <p className="text-xs text-red-600">{errors.agree}</p>
              )}
            </div>

            <Button type="submit" size="lg" className="w-full mt-2" disabled={submitting}>
              Kirim Kode OTP
            </Button>
          </form>

          <p className="mt-6 text-center text-sm text-gray-500">
            Sudah punya akun?{" "}
            <Link
              href="/login"
              className="font-semibold text-green-600 hover:text-green-700"
            >
              Masuk di sini
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
