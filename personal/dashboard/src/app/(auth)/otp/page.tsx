"use client";

import { useState, useRef, useEffect, useCallback, Suspense } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { MessageSquareMore, RotateCcw, ArrowLeft } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { authApi } from "@/lib/api";
import { getApiErrorMessage } from "@/lib/api-error";
import { clearPendingAuthState, getPendingAuthState } from "@/lib/auth-session";
import { useAuth } from "@/providers/auth-provider";

const OTP_LENGTH = 6;
const RESEND_COOLDOWN = 60; // seconds

function OTPPageInner() {
  const router = useRouter();
  const params = useSearchParams();
  const { completeAuth, session, hydrated, appConfig } = useAuth();
  const pendingAuth = getPendingAuthState();
  const phone = params.get("phone") ?? pendingAuth?.phoneNumber ?? "";
  const context = (params.get("context") ?? pendingAuth?.context ?? "login") as "login" | "register";

  const [otp, setOtp] = useState<string[]>(Array(OTP_LENGTH).fill(""));
  const [error, setError] = useState("");
  const [countdown, setCountdown] = useState(RESEND_COOLDOWN);
  const [submitting, setSubmitting] = useState(false);
  const [resending, setResending] = useState(false);

  const inputRefs = useRef<(HTMLInputElement | null)[]>([]);

  useEffect(() => {
    if (hydrated && session) {
      router.replace("/");
      return;
    }

    if (!pendingAuth || !phone) {
      router.replace(context === "register" ? "/register" : "/login");
    }
  }, [context, hydrated, pendingAuth, phone, router, session]);

  // Countdown timer
  useEffect(() => {
    if (countdown <= 0) return;
    const t = setInterval(() => setCountdown((c) => c - 1), 1000);
    return () => clearInterval(t);
  }, [countdown]);

  const focusInput = (index: number) => {
    inputRefs.current[index]?.focus();
  };

  const handleChange = (index: number, value: string) => {
    // Only allow single digit
    const digit = value.replace(/\D/g, "").slice(-1);
    const next = [...otp];
    next[index] = digit;
    setOtp(next);
    setError("");

    // Auto-advance
    if (digit && index < OTP_LENGTH - 1) {
      focusInput(index + 1);
    }
  };

  const handleKeyDown = (index: number, e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Backspace") {
      if (otp[index]) {
        const next = [...otp];
        next[index] = "";
        setOtp(next);
      } else if (index > 0) {
        focusInput(index - 1);
      }
    } else if (e.key === "ArrowLeft" && index > 0) {
      focusInput(index - 1);
    } else if (e.key === "ArrowRight" && index < OTP_LENGTH - 1) {
      focusInput(index + 1);
    }
  };

  const handlePaste = (e: React.ClipboardEvent) => {
    e.preventDefault();
    const pasted = e.clipboardData.getData("text").replace(/\D/g, "").slice(0, OTP_LENGTH);
    if (!pasted) return;
    const next = Array(OTP_LENGTH).fill("");
    pasted.split("").forEach((ch, i) => { next[i] = ch; });
    setOtp(next);
    focusInput(Math.min(pasted.length, OTP_LENGTH - 1));
  };

  const handleSubmit = useCallback(
    async (e?: React.FormEvent) => {
      e?.preventDefault();
      const code = otp.join("");
      if (code.length < OTP_LENGTH) {
        setError("Masukkan 6 digit kode OTP");
        return;
      }

      setSubmitting(true);

      try {
        const result = await authApi.verifyOTP({
          phone_number: phone,
          otp_code: code,
        });
        completeAuth(result, pendingAuth?.rememberMe ?? false);
        clearPendingAuthState();
        toast.success("Login berhasil");
        router.replace("/");
      } catch (submitError) {
        setError(getApiErrorMessage(submitError, "Verifikasi OTP gagal"));
      } finally {
        setSubmitting(false);
      }
    },
    [completeAuth, otp, pendingAuth?.rememberMe, phone, router]
  );

  // Auto-submit when all digits filled
  useEffect(() => {
    if (otp.every((d) => d !== "")) {
      handleSubmit();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [otp]);

  const handleResend = async () => {
    if (countdown > 0) return;

    setResending(true);

    try {
      if (context === "register") {
        await authApi.register({
          phone_number: phone,
          full_name: pendingAuth?.fullName ?? "User WACAST",
        });
      } else {
        await authApi.requestOTP({ phone_number: phone });
      }

      setOtp(Array(OTP_LENGTH).fill(""));
      setError("");
      setCountdown(RESEND_COOLDOWN);
      focusInput(0);
      toast.success("Kode OTP baru berhasil dikirim");
    } catch (resendError) {
      setError(getApiErrorMessage(resendError, "Gagal mengirim ulang OTP"));
    } finally {
      setResending(false);
    }
  };

  const maskedPhone = phone
    ? phone.slice(0, 4) + "****" + phone.slice(-3)
    : "nomor Anda";

  const backHref = context === "register" ? "/register" : "/login";

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
            Verifikasi identitas Anda untuk menjaga keamanan akun.
          </p>

          <div className="mt-10 rounded-2xl bg-white/10 p-6 text-left">
            <p className="mb-3 text-sm font-semibold text-green-100">
              Kenapa perlu OTP?
            </p>
            <ul className="space-y-2 text-sm text-green-100">
              <li className="flex items-start gap-2">
                <span className="mt-0.5 text-green-300">✓</span>
                Memastikan hanya pemilik nomor yang bisa masuk
              </li>
              <li className="flex items-start gap-2">
                <span className="mt-0.5 text-green-300">✓</span>
                Melindungi akun dari akses tidak sah
              </li>
              <li className="flex items-start gap-2">
                <span className="mt-0.5 text-green-300">✓</span>
                Kode hanya berlaku 60 detik
              </li>
            </ul>
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
          {/* Back link */}
          <Link
            href={backHref}
            className="mb-6 inline-flex items-center gap-1.5 text-sm text-gray-500 hover:text-gray-700"
          >
            <ArrowLeft className="h-4 w-4" />
            Kembali
          </Link>

          {/* Header */}
          <div className="mb-8">
            {/* OTP icon */}
            <div className="mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-green-50">
              <MessageSquareMore className="h-7 w-7 text-green-600" />
            </div>
            <h2 className="text-2xl font-bold text-gray-900">Verifikasi OTP</h2>
            <p className="mt-2 text-sm text-gray-500 leading-relaxed">
              Kode OTP 6 digit telah dikirim ke WhatsApp{" "}
              <span className="font-semibold text-gray-700">{maskedPhone}</span>
            </p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-6">
            {/* OTP input boxes */}
            <div className="flex justify-between gap-2" onPaste={handlePaste}>
              {otp.map((digit, i) => (
                <input
                  key={i}
                  ref={(el) => { inputRefs.current[i] = el; }}
                  type="text"
                  inputMode="numeric"
                  maxLength={1}
                  value={digit}
                  onChange={(e) => handleChange(i, e.target.value)}
                  onKeyDown={(e) => handleKeyDown(i, e)}
                  autoFocus={i === 0}
                  className={`h-13 w-full max-w-[52px] rounded-xl border-2 text-center text-xl font-bold transition-all focus:outline-none focus:ring-0 ${
                    error
                      ? "border-red-400 bg-red-50 text-red-600"
                      : digit
                      ? "border-green-500 bg-green-50 text-green-700"
                      : "border-gray-200 bg-white text-gray-900 focus:border-green-500"
                  }`}
                  style={{ height: "52px" }}
                />
              ))}
            </div>

            {/* Error message */}
            {error && (
              <p className="text-sm text-red-600 text-center">{error}</p>
            )}

            <Button
              type="submit"
              size="lg"
              className="w-full"
              disabled={otp.some((d) => !d) || submitting}
            >
              Verifikasi
            </Button>
          </form>

          {/* Resend section */}
          <div className="mt-6 text-center">
            <p className="text-sm text-gray-500">Tidak menerima kode?</p>
            <button
              type="button"
              onClick={handleResend}
              disabled={countdown > 0 || resending}
              className="mt-1 inline-flex items-center gap-1.5 text-sm font-semibold text-green-600 hover:text-green-700 disabled:cursor-not-allowed disabled:opacity-50"
            >
              <RotateCcw className="h-3.5 w-3.5" />
              {countdown > 0 ? `Kirim ulang dalam ${countdown}s` : "Kirim ulang kode"}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default function OTPPage() {
  return (
    <Suspense>
      <OTPPageInner />
    </Suspense>
  );
}
