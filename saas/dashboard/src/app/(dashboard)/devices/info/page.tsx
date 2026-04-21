"use client";

import { useState, useEffect } from "react";
import { Smartphone, MessageSquare, Clock, Wifi, WifiOff, Loader2, RefreshCw, Hash } from "lucide-react";
import { sessionsApi } from "@/lib/api";
import { Device } from "@/lib/types";
import { toast } from "sonner";

export default function MultiDeviceInfoPage() {
  const [devices, setDevices] = useState<Device[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const fetchData = () => {
    setIsLoading(true);
    sessionsApi.list()
      .then((res) => {
        setDevices(res.sessions || []);
      })
      .catch((err) => {
        toast.error("Gagal memuat detail perangkat.");
        console.error(err);
      })
      .finally(() => {
        setIsLoading(false);
      });
  };

  useEffect(() => {
    fetchData();
  }, []);

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="border-b border-gray-200 bg-white px-6 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold text-gray-900">Multi Device Info</h1>
            <p className="text-sm text-gray-500">Detail teknis lengkap semua device yang terdaftar</p>
          </div>
          <button
            onClick={fetchData}
            disabled={isLoading}
            className="inline-flex items-center gap-2 rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 disabled:opacity-50 transition-all"
          >
            {isLoading ? <RefreshCw className="h-4 w-4 animate-spin" /> : <RefreshCw className="h-4 w-4" />} Refresh
          </button>
        </div>
      </div>

      <div className="p-6">
        {isLoading && devices.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-24 space-y-4">
            <Loader2 className="h-10 w-10 animate-spin text-green-500" />
            <p className="text-sm text-gray-400">Mensinkronisasi data perangkat...</p>
          </div>
        ) : devices.length === 0 ? (
          <div className="rounded-2xl border-2 border-dashed border-gray-200 bg-white p-12 text-center uppercase tracking-tighter">
            <Smartphone className="mx-auto h-12 w-12 text-gray-200 mb-3" />
            <p className="text-gray-400 font-bold">Belum ada perangkat terdaftar</p>
          </div>
        ) : (
          <div className="space-y-5">
            {devices.map((d) => (
              <div key={d.device_id} className="group overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-sm transition-all hover:shadow-md">
                <div className="flex flex-col sm:flex-row divide-y sm:divide-y-0 sm:divide-x divide-gray-100">
                  {/* Left branding */}
                  <div className="p-6 sm:w-1/3 bg-gray-50/50">
                    <div className="flex items-center gap-4 mb-4">
                      <div className={`flex h-12 w-12 items-center justify-center rounded-2xl shadow-sm ${d.status === 1 ? "bg-green-100 text-green-600" : "bg-gray-100 text-gray-400"}`}>
                        <Smartphone className="h-6 w-6" />
                      </div>
                      <div>
                        <h3 className="text-lg font-bold text-gray-900 leading-tight">
                          {d.display_name || "Device WhatsApp"}
                        </h3>
                        <div className="mt-1 flex items-center gap-2">
                          {d.status === 1 ? (
                            <span className="inline-flex items-center gap-1 rounded-full bg-green-500/10 px-2 py-0.5 text-[10px] font-bold text-green-600 uppercase tracking-widest">
                              <Wifi className="h-2.5 w-2.5" /> Online
                            </span>
                          ) : (
                            <span className="inline-flex items-center gap-1 rounded-full bg-red-500/10 px-2 py-0.5 text-[10px] font-bold text-red-500 uppercase tracking-widest">
                              <WifiOff className="h-2.5 w-2.5" /> Offline
                            </span>
                          )}
                        </div>
                      </div>
                    </div>
                    <div className="space-y-2">
                       <div className="flex items-center gap-2 text-xs text-gray-500 bg-white border border-gray-100 rounded-lg p-2">
                          <Hash className="h-3 w-3" />
                          <span className="font-mono">{d.device_id}</span>
                       </div>
                    </div>
                  </div>

                  {/* Right stats */}
                  <div className="flex-1 p-6 grid grid-cols-2 lg:grid-cols-3 gap-6">
                    <div>
                      <p className="text-[10px] font-bold text-gray-400 uppercase tracking-widest mb-1">WhatsApp Number</p>
                      <p className="text-sm font-semibold text-gray-900 font-mono tracking-tight">{d.phone || "Not connected"}</p>
                    </div>
                    <div>
                      <p className="text-[10px] font-bold text-gray-400 uppercase tracking-widest mb-1">Status Code</p>
                      <p className="text-sm font-semibold text-gray-900">
                        {d.status === 1 ? "Active Session" : d.status === 2 ? "Waiting for Scan" : "Disconnected"}
                      </p>
                    </div>
                    <div>
                      <p className="text-[10px] font-bold text-gray-400 uppercase tracking-widest mb-1">Connection Mode</p>
                      <p className="text-sm font-semibold text-gray-900">Multi-Device Web</p>
                    </div>
                    <div>
                      <p className="text-[10px] font-bold text-gray-400 uppercase tracking-widest mb-1">Platform</p>
                      <p className="text-sm font-semibold text-gray-900">Chrome (Desktop)</p>
                    </div>
                    <div>
                      <p className="text-[10px] font-bold text-gray-400 uppercase tracking-widest mb-1">Last Activity</p>
                      <div className="flex items-center gap-1.5 text-sm font-semibold text-gray-900">
                        <Clock className="h-3.5 w-3.5 text-blue-500" />
                        <span>Recent</span>
                      </div>
                    </div>
                    <div className="col-span-1">
                       <p className="text-[10px] font-bold text-gray-400 uppercase tracking-widest mb-1">Data Storage</p>
                       <p className="text-sm font-semibold text-green-600">Persistent (Local)</p>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
