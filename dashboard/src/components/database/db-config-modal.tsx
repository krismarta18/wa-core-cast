"use client";

import { useState, useEffect, useRef } from "react";
import { X, Database, CheckCircle2, XCircle, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { configApi, type DbConfig } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

interface DbConfigModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export default function DbConfigModal({ open, onOpenChange }: DbConfigModalProps) {
  const [loading, setLoading] = useState(false);
  const [testing, setTesting] = useState(false);
  const [status, setStatus] = useState<"connected" | "disconnected" | "unknown">("unknown");
  const modalRef = useRef<HTMLDivElement>(null);

  const [form, setForm] = useState<DbConfig>({
    Host: "localhost",
    Port: 5432,
    User: "postgres",
    Password: "",
    DBName: "wacast",
    SSLMode: "disable",
  });

  useEffect(() => {
    if (open) {
      loadCurrentConfig();
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "unset";
    }
    return () => { document.body.style.overflow = "unset"; };
  }, [open]);

  // Close on Escape
  useEffect(() => {
    function handleKey(e: KeyboardEvent) {
      if (e.key === "Escape") onOpenChange(false);
    }
    if (open) document.addEventListener("keydown", handleKey);
    return () => document.removeEventListener("keydown", handleKey);
  }, [open, onOpenChange]);

  const loadCurrentConfig = async () => {
    try {
      const res = await configApi.getDbStatus();
      if (res.success) {
        setForm({
          ...form,
          Host: res.config.host,
          Port: res.config.port,
          User: res.config.user,
          DBName: res.config.database,
          SSLMode: res.config.ssl_mode,
        });
        setStatus(res.status as any);
      }
    } catch (error) {
      console.error("Failed to load DB status", error);
    }
  };

  const handleTest = async () => {
    setTesting(true);
    try {
      const res = await configApi.testDbConnection(form);
      if (res.success) {
        toast.success(res.message);
      } else {
        toast.error(res.message);
      }
    } catch (error) {
      toast.error("Failed to test connection. Backend might be unreachable.");
    } finally {
      setTesting(false);
    }
  };

  const handleSave = async () => {
    setLoading(true);
    try {
      const res = await configApi.saveDbConnection(form);
      if (res.success) {
        toast.success(res.message);
        setStatus("connected");
        onOpenChange(false);
      } else {
        toast.error(res.message);
      }
    } catch (error) {
      toast.error("Failed to save configuration.");
    } finally {
      setLoading(false);
    }
  };

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-[200] flex items-center justify-center p-4">
      {/* Backdrop */}
      <div 
        className="absolute inset-0 bg-black/40 backdrop-blur-sm" 
        onClick={() => onOpenChange(false)} 
      />

      {/* Modal Content */}
      <div 
        ref={modalRef}
        className="relative w-full max-w-lg rounded-2xl border border-gray-200 bg-white p-6 shadow-2xl animate-in fade-in zoom-in duration-200"
      >
        <button
          onClick={() => onOpenChange(false)}
          className="absolute right-4 top-4 rounded-lg p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-600 transition-colors"
        >
          <X className="h-5 w-5" />
        </button>

        <div className="mb-6">
          <div className="flex items-center gap-2 mb-1">
            <Database className="h-5 w-5 text-green-600" />
            <h2 className="text-xl font-bold text-gray-900">Database Configuration</h2>
          </div>
          <p className="text-sm text-gray-500">
            Configure your PostgreSQL connection settings below.
          </p>
        </div>

        <div className="space-y-4">
          <div className="flex items-center gap-2 px-3 py-2 rounded-lg bg-gray-50 border border-gray-100 text-sm">
            <span className="text-gray-500 font-medium">System Status:</span>
            {status === "connected" ? (
              <span className="flex items-center gap-1 text-green-600 font-semibold">
                <CheckCircle2 className="h-4 w-4" /> Connected
              </span>
            ) : (
              <span className="flex items-center gap-1 text-red-600 font-semibold">
                <XCircle className="h-4 w-4" /> Disconnected
              </span>
            )}
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <Input
              label="Host"
              placeholder="localhost"
              value={form.Host}
              onChange={(e) => setForm({ ...form, Host: e.target.value })}
            />
            <Input
              label="Port"
              type="number"
              placeholder="5432"
              value={form.Port}
              onChange={(e) => setForm({ ...form, Port: parseInt(e.target.value) || 5432 })}
            />
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <Input
              label="User"
              placeholder="postgres"
              value={form.User}
              onChange={(e) => setForm({ ...form, User: e.target.value })}
            />
            <Input
              label="Password"
              type="password"
              placeholder="••••••••"
              value={form.Password}
              onChange={(e) => setForm({ ...form, Password: e.target.value })}
            />
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <Input
              label="Database Name"
              placeholder="wacast"
              value={form.DBName}
              onChange={(e) => setForm({ ...form, DBName: e.target.value })}
            />
            <div className="flex flex-col gap-1">
              <label className="text-sm font-medium text-gray-700">SSL Mode</label>
              <select
                value={form.SSLMode}
                onChange={(e) => setForm({ ...form, SSLMode: e.target.value })}
                className="h-9 rounded-lg border border-gray-300 bg-white px-3 text-sm focus:outline-none focus:ring-2 focus:ring-green-500 appearance-none bg-no-repeat bg-[right_0.5rem_center]"
              >
                <option value="disable">Disable</option>
                <option value="require">Require</option>
                <option value="verify-ca">Verify CA</option>
                <option value="verify-full">Verify Full</option>
              </select>
            </div>
          </div>
        </div>

        <div className="mt-8 flex flex-col sm:flex-row justify-end gap-3 border-t pt-6">
          <Button
            type="button"
            variant="outline"
            onClick={handleTest}
            disabled={testing || loading}
            className="w-full sm:w-auto"
          >
            {testing ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Testing...
              </>
            ) : (
              "Test Connection"
            )}
          </Button>
          <Button
            type="button"
            onClick={handleSave}
            disabled={testing || loading}
            className="w-full sm:w-auto bg-green-600 hover:bg-green-700 text-white font-semibold"
          >
            {loading ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Saving...
              </>
            ) : (
              "Save Configuration"
            )}
          </Button>
        </div>
      </div>
    </div>
  );
}
