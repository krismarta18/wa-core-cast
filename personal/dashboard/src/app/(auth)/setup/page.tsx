"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { 
  Database, 
  ShieldCheck, 
  Settings2, 
  Copy, 
  CheckCircle2, 
  XCircle, 
  Loader2, 
  Server,
  Key
} from "lucide-react";
import { toast } from "sonner";
import { configApi, licenseApi, type DbConfig, type LicenseStatus } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

export default function SetupPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [licenseStatus, setLicenseStatus] = useState<LicenseStatus | null>(null);
  
  // License State
  const [serialKey, setSerialKey] = useState("");
  
  useEffect(() => {
    checkInitialStatus();
  }, []);

  const checkInitialStatus = async () => {
    try {
      const lic = await licenseApi.getStatus();
      setLicenseStatus(lic.data);
    } catch (error) {
      console.error("Failed to fetch setup status", error);
    }
  };

  const copyHWID = () => {
    if (licenseStatus?.hwid) {
      navigator.clipboard.writeText(licenseStatus.hwid);
      toast.success("Hardware ID copied to clipboard");
    }
  };

  const handleActivate = async () => {
    if (!serialKey) return toast.error("Please enter a serial key");
    setLoading(true);
    try {
      const res = await licenseApi.activate(serialKey);
      if (res.success) {
        toast.success(res.message);
        router.push("/login"); // Direct to login after activation
      } else {
        toast.error(res.message);
      }
    } catch (error: any) {
      toast.error(error.response?.data?.message || "Failed to activate license");
    } finally {
      setLoading(false);
    }
  };

  const waNumber = '6285887373722';
  const waLink = `https://wa.me/${waNumber}?text=${encodeURIComponent(
    `Halo Admin, saya ingin aktivasi WACAST PRO.\n\nHWID saya: ${licenseStatus?.hwid || "..."}`
  )}`;

  return (
    <div className="min-h-screen bg-slate-950 flex flex-col items-center justify-center p-4">
      <div className="w-full max-w-md">
        <div className="flex items-center justify-center gap-3 mb-10">
          <div className="bg-blue-600 p-2.5 rounded-2xl shadow-lg shadow-blue-500/20">
            <ShieldCheck className="h-8 w-8 text-white" />
          </div>
          <h1 className="text-3xl font-bold text-white tracking-tight">System Activation</h1>
        </div>

        <Card className="border-slate-800 bg-slate-900 shadow-2xl">
          <CardHeader className="pb-4">
            <CardTitle className="text-xl text-white">License Activation</CardTitle>
            <CardDescription className="text-slate-400">
              Enter your serial key to unlock WACAST PRO features.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="p-4 bg-slate-950 border border-slate-800 rounded-2xl space-y-3">
              <label className="text-[10px] font-black text-slate-500 uppercase tracking-widest">
                Hardware ID (HWID)
              </label>
              <div className="flex items-center gap-2">
                <code className="flex-1 bg-slate-900 border border-slate-800 px-3 py-2.5 rounded-lg font-mono text-sm text-blue-400 font-bold overflow-hidden text-ellipsis">
                  {licenseStatus?.hwid || "..."}
                </code>
                <Button variant="secondary" size="icon" onClick={copyHWID} className="h-10 w-10 shrink-0">
                  <Copy className="h-4 w-4" />
                </Button>
              </div>
            </div>

            <div className="space-y-3">
              <label className="text-xs font-bold text-slate-400 uppercase tracking-wider">
                Serial Key
              </label>
              <Input 
                placeholder="XXXX-XXXX-XXXX-XXXX"
                value={serialKey}
                onChange={(e) => setSerialKey(e.target.value)}
                className="h-12 bg-slate-950 border-slate-800 text-white font-mono text-center tracking-widest uppercase focus:border-blue-500 focus:ring-blue-500/20"
              />
            </div>

            <Button 
              onClick={handleActivate} 
              disabled={loading}
              className="w-full h-12 text-sm font-bold bg-blue-600 hover:bg-blue-700 shadow-lg shadow-blue-500/10 transition-all active:scale-[0.98]"
            >
              {loading ? <Loader2 className="animate-spin h-5 w-5" /> : "Activate System"}
            </Button>

            <div className="pt-4 border-t border-slate-800">
              <div className="text-center space-y-4">
                <p className="text-xs text-slate-500">Need a license key?</p>
                <a 
                  href={waLink}
                  target="_blank" 
                  className="flex items-center justify-center gap-2 w-full py-3 rounded-xl bg-emerald-600/10 border border-emerald-600/20 text-emerald-400 hover:bg-emerald-600/20 transition-all font-semibold text-sm"
                >
                  <svg width="18" height="18" fill="currentColor" viewBox="0 0 24 24"><path d="M.057 24l1.687-6.163c-1.041-1.804-1.588-3.849-1.587-5.946.003-6.556 5.338-11.891 11.893-11.891 3.181.001 6.167 1.24 8.413 3.488 2.245 2.248 3.481 5.236 3.48 8.414-.003 6.557-5.338 11.892-11.893 11.892-1.99-.001-3.951-.5-5.688-1.448l-6.305 1.654zm6.597-3.807c1.676.995 3.276 1.591 5.392 1.592 5.448 0 9.886-4.438 9.889-9.885.002-5.462-4.415-9.89-9.881-9.892-5.452 0-9.887 4.434-9.889 9.884-.001 2.225.651 3.891 1.746 5.634l-.999 3.648 3.742-.981zm11.387-5.464c-.074-.124-.272-.198-.57-.347-.297-.149-1.758-.868-2.031-.967-.272-.099-.47-.149-.669.149-.198.297-.768.967-.941 1.165-.173.198-.347.223-.644.074-.297-.149-1.255-.462-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.297-.347.446-.521.151-.172.2-.296.3-.495.099-.198.05-.372-.025-.521-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.501-.669-.51l-.57-.01c-.198 0-.52.074-.792.372s-1.04 1.016-1.04 2.479 1.065 2.876 1.213 3.074c.149.198 2.095 3.2 5.076 4.487.709.306 1.263.489 1.694.626.712.226 1.36.194 1.872.118.571-.085 1.758-.719 2.006-1.413.248-.695.248-1.29.173-1.414z"/></svg>
                  WhatsApp Admin
                </a>
                <p className="text-[10px] text-slate-600 font-medium tracking-tight">0858-8737-3722</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <p className="mt-8 text-center text-[10px] text-slate-600 uppercase tracking-[0.2em]">
          Powered by WACAST Engine v1.0
        </p>
      </div>
    </div>
  );
}
