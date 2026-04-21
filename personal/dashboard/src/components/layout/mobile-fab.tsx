"use client";

import { MessageSquarePlus } from "lucide-react";
import { useRouter } from "next/navigation";

export function MobileFAB() {
  const router = useRouter();
  return (
    <button
      onClick={() => router.push("/messaging/new")}
      aria-label="Kirim Pesan"
      className="fixed bottom-6 right-6 z-30 flex h-14 w-14 items-center justify-center rounded-full bg-green-600 shadow-xl hover:bg-green-700 active:scale-95 transition-transform lg:hidden"
    >
      <MessageSquarePlus className="h-6 w-6 text-white" />
    </button>
  );
}
