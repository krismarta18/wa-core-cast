"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { MessageSquareMore, Database } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { authApi, persistAuthSession } from "@/lib/api";
import { getApiErrorMessage } from "@/lib/api-error";
import { useAuth } from "@/providers/auth-provider";
import DbConfigModal from "@/components/database/db-config-modal";

export default function LoginPage() {
  const router = useRouter();
  const { session, hydrated, appConfig } = useAuth();
  const [password, setPassword] = useState("");
  const [rememberMe, setRememberMe] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [errors, setErrors] = useState<{ password?: string }>({});
  const [dbModalOpen, setDbModalOpen] = useState(false);

  useEffect(() => {
    if (hydrated && session) {
      router.replace("/");
    }
  }, [hydrated, router, session]);

  const validate = () => {
    const e: { password?: string } = {};
    if (!password) e.password = "Password wajib diisi";
    setErrors(e);
    return Object.keys(e).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate()) return;

    setSubmitting(true);

    try {
      const response = await authApi.loginWithPassword({ password });
      // Use existing persistAuthSession logic
      persistAuthSession(response, rememberMe);
      toast.success("Login berhasil");
      router.push("/");
    } catch (error) {
      setErrors({
        password: getApiErrorMessage(error, "Password salah atau terjadi kesalahan"),
      });
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="flex min-h-screen">
      {/* Left panel — branding */}
      <div className="hidden lg:flex lg:w-1/2 flex-col items-center justify-center bg-gradient-to-br from-green-600 to-green-800 px-12 text-white">
        <div className="max-w-sm text-center">
          <div className="mb-6 flex justify-center">
            <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-white/10 backdrop-blur p-2">
              <img src="/icon.png" alt="Logo" className="h-full w-full object-contain" />
            </div>
          </div>
          <h1 className="text-3xl font-bold">{appConfig.app_name || "WACAST"}</h1>
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
        <div className="mb-8 flex items-center gap-3 lg:hidden">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-white shadow-sm border border-gray-100 p-1.5">
            <img src="/icon.png" alt="Logo" className="h-full w-full object-contain" />
          </div>
          <span className="text-xl font-bold text-gray-900">{appConfig.app_name || "WACAST"}</span>
        </div>

        <div className="w-full max-w-sm">
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-gray-900">Admin Login</h2>
            <p className="mt-1 text-sm text-gray-500">
              Masukkan password untuk masuk ke dashboard
            </p>
          </div>

          <form onSubmit={handleSubmit} noValidate className="space-y-4">
            <div className="flex flex-col gap-1">
              <label htmlFor="password" className="text-sm font-medium text-gray-700">
                Password Aplikasi
              </label>
              <input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Masukkan password..."
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                disabled={submitting}
              />
              {errors.password && (
                <p className="text-xs text-red-600">{errors.password}</p>
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

            <Button type="submit" size="lg" className="w-full mt-2" disabled={submitting}>
              Masuk Dashboard
            </Button>
          </form>

          <div className="mt-12 flex justify-center border-t pt-6">
            <button
              onClick={() => setDbModalOpen(true)}
              className="flex items-center gap-2 text-sm font-medium text-gray-500 hover:text-green-600 transition-colors"
            >
              <Database className="h-4 w-4" />
              Database Connection
            </button>
          </div>
        </div>
      </div>

      <DbConfigModal open={dbModalOpen} onOpenChange={setDbModalOpen} />
    </div>
  );
}
