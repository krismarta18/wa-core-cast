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
  const [step, setStep] = useState<"license" | "database">("license");
  const [loading, setLoading] = useState(false);
  const [licenseStatus, setLicenseStatus] = useState<LicenseStatus | null>(null);
  
  // License State
  const [serialKey, setSerialKey] = useState("");
  
  // DB State
  const [dbForm, setDbForm] = useState<DbConfig>({
    Host: "localhost",
    Port: 5432,
    User: "postgres",
    Password: "",
    DBName: "wacast",
    SSLMode: "disable",
  });

  useEffect(() => {
    checkInitialStatus();
  }, []);

  const checkInitialStatus = async () => {
    try {
      const lic = await licenseApi.getStatus();
      setLicenseStatus(lic.data);
      if (lic.data?.is_active && !lic.data?.is_expired) {
        setStep("database");
      }
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
        // Refresh status to update HWID and active state
        await checkInitialStatus();
        setStep("database");
      } else {
        toast.error(res.message);
      }
    } catch (error: any) {
      toast.error(error.response?.data?.message || "Failed to activate license");
    } finally {
      setLoading(false);
    }
  };

  const handleSaveDB = async () => {
    setLoading(true);
    try {
      const res = await configApi.saveDbConnection(dbForm);
      if (res.success) {
        toast.success("Configuration complete!");
        router.push("/login"); // Go to login after full setup
      } else {
        toast.error(res.message);
      }
    } catch (error: any) {
      toast.error(error.response?.data?.message || "Failed to save database config");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col items-center justify-center p-4">
      <div className="w-full max-w-2xl">
        <div className="flex items-center justify-center gap-3 mb-8">
          <div className="bg-green-600 p-2 rounded-xl">
            <Settings2 className="h-8 w-8 text-white" />
          </div>
          <h1 className="text-3xl font-bold text-gray-900 tracking-tight">WACAST Setup</h1>
        </div>

        <div className="flex items-center justify-between mb-8 px-4">
          <div className={`flex flex-col items-center gap-2 ${step === "license" ? "text-green-600" : "text-gray-400"}`}>
            <div className={`h-10 w-10 rounded-full flex items-center justify-center border-2 ${step === "license" ? "border-green-600 bg-green-50 shadow-sm" : "border-gray-200"}`}>
              <Key className="h-5 w-5" />
            </div>
            <span className="text-xs font-semibold uppercase tracking-wider">License</span>
          </div>
          <div className="h-[2px] flex-1 mx-4 bg-gray-200 mt-[-20px]">
            <div className={`h-full bg-green-600 transition-all duration-500 ${step === "database" ? "w-full" : "w-0"}`} />
          </div>
          <div className={`flex flex-col items-center gap-2 ${step === "database" ? "text-green-600" : "text-gray-400"}`}>
            <div className={`h-10 w-10 rounded-full flex items-center justify-center border-2 ${step === "database" ? "border-green-600 bg-green-50 shadow-sm" : "border-gray-200"}`}>
              <Database className="h-5 w-5" />
            </div>
            <span className="text-xs font-semibold uppercase tracking-wider">Database</span>
          </div>
        </div>

        {step === "license" && (
          <Card className="border-none shadow-xl">
            <CardHeader>
              <CardTitle>System Activation</CardTitle>
              <CardDescription>Enter your serial key to unlock the WACAST Core Dashboard.</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="p-4 bg-gray-50 border rounded-xl space-y-3">
                <label className="text-xs font-bold text-gray-400 uppercase">Your Hardware ID (HWID)</label>
                <div className="flex items-center gap-3">
                  <code className="flex-1 bg-white border px-4 py-3 rounded-lg font-mono text-lg text-green-700 font-bold overflow-hidden text-ellipsis">
                    {licenseStatus?.hwid || "FETCHING..."}
                  </code>
                  <Button variant="outline" size="icon" onClick={copyHWID} className="h-12 w-12 rounded-lg">
                    <Copy className="h-5 w-5" />
                  </Button>
                </div>
                <p className="text-[11px] text-gray-500 leading-relaxed italic">
                  *Kirim HWID di atas ke Admin untuk mendapatkan Serial Key produk.
                </p>
              </div>

              <div className="space-y-2">
                <label className="text-sm font-semibold text-gray-700">Enter Serial Key</label>
                <Input 
                  placeholder="Paste your serial key here"
                  value={serialKey}
                  onChange={(e) => setSerialKey(e.target.value)}
                  className="h-14 font-mono text-center tracking-widest uppercase border-2 focus:ring-green-500"
                />
              </div>

              <Button 
                onClick={handleActivate} 
                disabled={loading}
                className="w-full h-14 text-lg font-bold bg-green-600 hover:bg-green-700 shadow-lg shadow-green-200 transition-all active:scale-[0.98]"
              >
                {loading ? <Loader2 className="animate-spin" /> : "Activate Now"}
              </Button>
            </CardContent>
          </Card>
        )}

        {step === "database" && (
          <Card className="border-none shadow-xl">
            <CardHeader>
              <CardTitle>Database Connection</CardTitle>
              <CardDescription>Configure where WACAST will store its data.</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="grid grid-cols-2 gap-4">
                <Input 
                  label="Host" 
                  value={dbForm.Host} 
                  onChange={(e) => setDbForm({...dbForm, Host: e.target.value})}
                />
                <Input 
                  label="Port" 
                  type="number"
                  value={dbForm.Port} 
                  onChange={(e) => setDbForm({...dbForm, Port: parseInt(e.target.value) || 5432})}
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <Input 
                  label="User" 
                  value={dbForm.User} 
                  onChange={(e) => setDbForm({...dbForm, User: e.target.value})}
                />
                <Input 
                  label="Password" 
                  type="password"
                  value={dbForm.Password} 
                  onChange={(e) => setDbForm({...dbForm, Password: e.target.value})}
                />
              </div>
              <Input 
                label="Database Name" 
                value={dbForm.DBName} 
                onChange={(e) => setDbForm({...dbForm, DBName: e.target.value})}
              />

              <Button 
                onClick={handleSaveDB} 
                disabled={loading}
                className="w-full h-14 text-lg font-bold bg-green-600 hover:bg-green-700"
              >
                {loading ? <Loader2 className="animate-spin" /> : "Finish Setup"}
              </Button>
            </CardContent>
          </Card>
        )}

        <div className="mt-8 text-center">
          <p className="text-gray-400 text-xs flex items-center justify-center gap-1">
            <ShieldCheck className="h-3 w-3" /> Secure Hardware-Locked System
          </p>
        </div>
      </div>
    </div>
  );
}
