"use client";

import { useState, useEffect, useCallback } from "react";
import { 
  Flame, 
  RefreshCw, 
  Play, 
  Square, 
  Clock, 
  CheckCircle2, 
  AlertCircle,
  Smartphone,
  ShieldCheck,
  Zap,
  Loader2
} from "lucide-react";
import { sessionsApi, warmingApi } from "@/lib/api";
import { Device } from "@/lib/types";
import { toast } from "sonner";
import { cn } from "@/lib/utils";

export default function WarmingPage() {
  const [devices, setDevices] = useState<Device[]>([]);
  const [selectedDevices, setSelectedDevices] = useState<string[]>([]);
  const [duration, setDuration] = useState(20);
  const [status, setStatus] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isStarting, setIsStarting] = useState(false);
  const [isStopping, setIsStopping] = useState(false);

  const fetchStatus = useCallback(async () => {
    try {
      const res = await warmingApi.status();
      if (res.success) {
        setStatus(res.data);
      }
    } catch (err) {
      console.error("Failed to fetch warming status", err);
    }
  }, []);

  const fetchDevices = useCallback(async () => {
    setIsLoading(true);
    try {
      const res = await sessionsApi.list();
      setDevices(res.sessions?.filter(d => d.status === 1) || []);
    } catch (err) {
      toast.error("Gagal memuat daftar perangkat aktif.");
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchDevices();
    fetchStatus();
    const interval = setInterval(fetchStatus, 5000);
    return () => clearInterval(interval);
  }, [fetchDevices, fetchStatus]);

  const handleToggleDevice = (id: string) => {
    setSelectedDevices(prev => 
      prev.includes(id) ? prev.filter(d => d !== id) : [...prev, id]
    );
  };

  const handleStart = async () => {
    if (selectedDevices.length < 2) {
      toast.error("Pilih minimal 2 perangkat untuk memulai warming.");
      return;
    }

    setIsStarting(true);
    try {
      const res = await warmingApi.start(selectedDevices, duration);
      if (res.success) {
        toast.success(res.message);
        fetchStatus();
      }
    } catch (err: any) {
      toast.error(err.response?.data?.error || "Gagal memulai warming.");
    } finally {
      setIsStarting(false);
    }
  };

  const handleStop = async () => {
    setIsStopping(true);
    try {
      const res = await warmingApi.stop();
      if (res.success) {
        toast.success(res.message);
        fetchStatus();
      }
    } catch (err: any) {
      toast.error(err.response?.data?.error || "Gagal menghentikan warming.");
    } finally {
      setIsStopping(false);
    }
  };

  const [localTimeRemaining, setLocalTimeRemaining] = useState(0);

  useEffect(() => {
    if (status?.is_active) {
      setLocalTimeRemaining(status.remaining_seconds);
    } else {
      setLocalTimeRemaining(0);
    }
  }, [status]);

  useEffect(() => {
    if (localTimeRemaining > 0) {
      const timer = setInterval(() => {
        setLocalTimeRemaining(prev => Math.max(0, prev - 1));
      }, 1000);
      return () => clearInterval(timer);
    }
  }, [localTimeRemaining]);

  const isWarming = status?.is_active;
  const timeRemaining = localTimeRemaining;
  const progress = status?.total_duration_seconds 
    ? ((status.total_duration_seconds - timeRemaining) / status.total_duration_seconds) * 100 
    : 0;

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  return (
    <div className="min-h-screen bg-gray-50 pb-20">
      {/* Header Section */}
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-6 py-8">
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
            <div className="flex items-center gap-4">
              <div className="h-14 w-14 rounded-2xl bg-orange-500 flex items-center justify-center shadow-lg shadow-orange-100">
                <Flame className="h-8 w-8 text-white animate-pulse" />
              </div>
              <div>
                <h1 className="text-2xl font-bold text-gray-900">WhatsApp Account Warming</h1>
                <p className="text-gray-500">Tingkatkan reputasi nomor Anda untuk menghindari blokir.</p>
              </div>
            </div>
            
            {isWarming && (
              <div className="flex items-center gap-3 bg-orange-50 px-4 py-2 rounded-full border border-orange-100">
                <span className="flex h-2 w-2 rounded-full bg-orange-500 animate-ping" />
                <span className="text-sm font-bold text-orange-700">SESI AKTIF: {formatTime(timeRemaining)}</span>
              </div>
            )}
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-6 py-8 grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Left Column: Configuration */}
        <div className="lg:col-span-2 space-y-8">
          {/* Device Selection Card */}
          <div className="bg-white rounded-3xl border border-gray-200 shadow-sm overflow-hidden">
            <div className="px-6 py-5 border-b border-gray-100 flex items-center justify-between">
              <h2 className="text-lg font-bold text-gray-900 flex items-center gap-2">
                <Smartphone className="h-5 w-5 text-green-600" />
                Pilih Perangkat
              </h2>
              <button 
                onClick={fetchDevices}
                className="text-sm text-green-600 font-semibold hover:underline"
              >
                Refresh
              </button>
            </div>
            
            <div className="p-6">
              {isLoading ? (
                <div className="py-12 flex flex-col items-center">
                  <Loader2 className="h-8 w-8 animate-spin text-green-500" />
                  <p className="mt-2 text-sm text-gray-400">Memuat perangkat...</p>
                </div>
              ) : devices.length === 0 ? (
                <div className="py-12 text-center">
                  <AlertCircle className="h-10 w-10 text-gray-300 mx-auto mb-3" />
                  <p className="text-gray-500">Tidak ada perangkat aktif ditemukan.</p>
                  <p className="text-xs text-gray-400 mt-1">Pastikan minimal 2 perangkat terhubung.</p>
                </div>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {devices.map((device) => (
                    <button
                      key={device.device_id}
                      onClick={() => !isWarming && handleToggleDevice(device.device_id)}
                      disabled={isWarming}
                      className={cn(
                        "flex items-center gap-4 p-4 rounded-2xl border transition-all text-left",
                        selectedDevices.includes(device.device_id)
                          ? "border-green-500 bg-green-50 ring-2 ring-green-500 ring-offset-2"
                          : "border-gray-100 bg-white hover:border-green-200"
                      )}
                    >
                      <div className={cn(
                        "h-10 w-10 rounded-xl flex items-center justify-center",
                        selectedDevices.includes(device.device_id) ? "bg-green-500 text-white" : "bg-gray-100 text-gray-500"
                      )}>
                        <Smartphone className="h-5 w-5" />
                      </div>
                      <div className="flex-1 overflow-hidden">
                        <p className="font-bold text-gray-900 truncate">{device.display_name || "Device"}</p>
                        <p className="text-xs text-gray-500 font-mono">{device.phone}</p>
                      </div>
                      <div className={cn(
                        "h-6 w-6 rounded-full border-2 flex items-center justify-center transition-colors",
                        selectedDevices.includes(device.device_id) ? "bg-green-500 border-green-500" : "border-gray-200"
                      )}>
                        {selectedDevices.includes(device.device_id) && <CheckCircle2 className="h-4 w-4 text-white" />}
                      </div>
                    </button>
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* Settings Card */}
          <div className="bg-white rounded-3xl border border-gray-200 shadow-sm p-6">
            <h2 className="text-lg font-bold text-gray-900 flex items-center gap-2 mb-6">
              <Clock className="h-5 w-5 text-blue-600" />
              Durasi Warming
            </h2>
            
            <div className="space-y-6">
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm font-medium text-gray-600">Menit Penggunaan</span>
                <span className="text-2xl font-black text-blue-600">{duration} <small className="text-xs font-normal text-gray-400">menit</small></span>
              </div>
              <input
                type="range"
                min="5"
                max="60"
                step="5"
                value={duration}
                onChange={(e) => setDuration(parseInt(e.target.value))}
                disabled={isWarming}
                className="w-full h-2 bg-gray-100 rounded-lg appearance-none cursor-pointer accent-blue-600"
              />
              <div className="flex justify-between text-xs text-gray-400 px-1">
                <span>5m</span>
                <span>15m</span>
                <span>30m</span>
                <span>45m</span>
                <span>60m</span>
              </div>
              
              <div className="bg-blue-50 p-4 rounded-2xl flex gap-3">
                <ShieldCheck className="h-5 w-5 text-blue-600 shrink-0" />
                <p className="text-xs text-blue-800 leading-relaxed">
                  Sistem akan secara otomatis melakukan simulasi percakapan natural (*Ping-Pong*) antar perangkat terpilih untuk membangun *Trust Score* di sistem WhatsApp.
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Right Column: Status & Control */}
        <div className="space-y-8">
          <div className="bg-white rounded-3xl border border-gray-200 shadow-xl p-8 sticky top-8">
            <div className="text-center mb-8">
              <div className={cn(
                "h-24 w-24 rounded-full mx-auto flex items-center justify-center mb-4 transition-all duration-500",
                isWarming ? "bg-orange-500 shadow-2xl shadow-orange-200" : "bg-gray-100 shadow-inner"
              )}>
                {isWarming ? (
                  <Zap className="h-12 w-12 text-white animate-bounce" />
                ) : (
                  <Flame className="h-12 w-12 text-gray-300" />
                )}
              </div>
              <h3 className="text-xl font-black text-gray-900">
                {isWarming ? "Sedang Memanaskan..." : "Siap Dimulai"}
              </h3>
              <p className="text-sm text-gray-500 mt-1">
                {isWarming ? "Lockdown mode aktif" : "Pilih perangkat & klik mulai"}
              </p>
            </div>

            {isWarming && (
              <div className="space-y-6 mb-8">
                <div className="h-3 w-full bg-gray-100 rounded-full overflow-hidden">
                  <div 
                    className="h-full bg-gradient-to-r from-orange-400 to-orange-600 transition-all duration-1000 ease-linear"
                    style={{ width: `${progress}%` }}
                  />
                </div>
                
                <div className="grid grid-cols-2 gap-4">
                  <div className="bg-gray-50 p-4 rounded-2xl text-center">
                    <p className="text-[10px] uppercase font-bold text-gray-400">Tersisa</p>
                    <p className="text-lg font-black text-gray-900">{formatTime(timeRemaining)}</p>
                  </div>
                  <div className="bg-gray-50 p-4 rounded-2xl text-center">
                    <p className="text-[10px] uppercase font-bold text-gray-400">Progress</p>
                    <p className="text-lg font-black text-gray-900">{Math.round(progress)}%</p>
                  </div>
                </div>
              </div>
            )}

            <div className="space-y-3">
              {!isWarming ? (
                <button
                  onClick={handleStart}
                  disabled={isStarting || selectedDevices.length < 2}
                  className="w-full bg-green-600 hover:bg-green-700 disabled:bg-gray-200 text-white py-4 rounded-2xl font-bold flex items-center justify-center gap-3 shadow-lg shadow-green-100 transition-all active:scale-95"
                >
                  {isStarting ? <Loader2 className="h-5 w-5 animate-spin" /> : <Play className="h-5 w-5 fill-current" />}
                  Mulai Warming Sekarang
                </button>
              ) : (
                <button
                  onClick={handleStop}
                  disabled={isStopping}
                  className="w-full bg-red-600 hover:bg-red-700 text-white py-4 rounded-2xl font-bold flex items-center justify-center gap-3 shadow-lg shadow-red-100 transition-all active:scale-95"
                >
                  {isStopping ? <Loader2 className="h-5 w-5 animate-spin" /> : <Square className="h-5 w-5 fill-current" />}
                  Hentikan Paksa
                </button>
              )}
            </div>

            <div className="mt-8 border-t border-gray-100 pt-6">
              <div className="flex items-center gap-2 mb-4">
                <div className="h-1.5 w-1.5 rounded-full bg-blue-500" />
                <span className="text-xs font-bold text-gray-400 uppercase">Peraturan Lockdown</span>
              </div>
              <ul className="space-y-3">
                <li className="flex gap-2 text-xs text-gray-500">
                  <CheckCircle2 className="h-4 w-4 text-green-500 shrink-0" />
                  Broadcast & pengiriman pesan manual akan diblokir.
                </li>
                <li className="flex gap-2 text-xs text-gray-500">
                  <CheckCircle2 className="h-4 w-4 text-green-500 shrink-0" />
                  Auto-response tetap berjalan normal.
                </li>
                <li className="flex gap-2 text-xs text-gray-500">
                  <CheckCircle2 className="h-4 w-4 text-green-500 shrink-0" />
                  Sangat disarankan durasi minimal 15 menit.
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
