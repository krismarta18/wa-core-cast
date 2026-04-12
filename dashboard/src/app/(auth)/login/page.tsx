"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { MessageSquareMore } from "lucide-react";
import { Button } from "@/components/ui/button";
import { PhoneInput } from "@/components/ui/phone-input";

export default function LoginPage() {
  const router = useRouter();
  const [dialCode, setDialCode] = useState("+62");
  const [phone, setPhone] = useState("");
  const [rememberMe, setRememberMe] = useState(false);
  const [errors, setErrors] = useState<{ phone?: string }>({});

  const validate = () => {
    const e: { phone?: string } = {};
    if (!phone) e.phone = "Nomor HP wajib diisi";
    else if (phone.length < 7) e.phone = "Nomor tidak valid";
    setErrors(e);
    return Object.keys(e).length === 0;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate()) return;
    const fullPhone = dialCode.replace("+", "") + phone;
    router.push(`/otp?phone=${encodeURIComponent(fullPhone)}&context=login`);
  };

  return (
    <div className="flex min-h-screen">
      {/* Left panel — branding */}
      <div className="hidden lg:flex lg:w-1/2 flex-col items-center justify-center bg-gradient-to-br from-green-600 to-green-800 px-12 text-white">
        <div className="max-w-sm text-center">
          <div className="mb-6 flex justify-center">
            <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-white/20 backdrop-blur">
              <MessageSquareMore className="h-8 w-8 text-white" />
            </div>
          </div>
          <h1 className="text-3xl font-bold">WACAST</h1>
          <p className="mt-3 text-base text-green-100 leading-relaxed">
            Platform manajemen WhatsApp Gateway yang handal dan mudah digunakan
            untuk bisnis Anda.
          </p>

          <div className="mt-10 space-y-4 text-left">
            {[
              "Multi-device WhatsApp dalam satu dashboard",
              "Kirim pesan massal dan terjadwal",
              "Real-time monitoring dan analytics",
            ].map((f) => (
              <div key={f} className="flex items-center gap-3">
                <div className="flex h-5 w-5 flex-shrink-0 items-center justify-center rounded-full bg-white/25">
                  <span className="text-xs">✓</span>
                </div>
                <span className="text-sm text-green-100">{f}</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Right panel — form */}
      <div className="flex w-full flex-col items-center justify-center bg-gray-50 px-6 py-12 lg:w-1/2">
        {/* Mobile logo */}
        <div className="mb-8 flex items-center gap-2 lg:hidden">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-green-600">
            <MessageSquareMore className="h-5 w-5 text-white" />
          </div>
          <span className="text-xl font-bold text-gray-900">WACAST</span>
        </div>

        <div className="w-full max-w-sm">
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-gray-900">Selamat Datang</h2>
            <p className="mt-1 text-sm text-gray-500">
              Masukkan nomor HP untuk masuk ke WACAST
            </p>
          </div>

          <form onSubmit={handleSubmit} noValidate className="space-y-4">
            {/* Phone with country code */}
            <div className="flex flex-col gap-1">
              <label htmlFor="phone" className="text-sm font-medium text-gray-700">
                Nomor HP
              </label>
              <PhoneInput
                id="phone"
                dialCode={dialCode}
                onDialCodeChange={setDialCode}
                value={phone}
                onChange={setPhone}
                error={errors.phone}
              />
              {errors.phone && (
                <p className="text-xs text-red-600">{errors.phone}</p>
              )}
            </div>

            {/* Remember me */}
            <div className="flex items-center text-sm">
              <label className="flex items-center gap-2 text-gray-600 cursor-pointer select-none">
                <input
                  type="checkbox"
                  checked={rememberMe}
                  onChange={(e) => setRememberMe(e.target.checked)}
                  className="h-4 w-4 rounded border-gray-300 accent-green-600"
                />
                Ingat saya
              </label>
            </div>

            <Button type="submit" size="lg" className="w-full mt-2">
              Kirim Kode OTP
            </Button>
          </form>

          <p className="mt-6 text-center text-sm text-gray-500">
            Belum punya akun?{" "}
            <Link
              href="/register"
              className="font-semibold text-green-600 hover:text-green-700"
            >
              Daftar sekarang
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
