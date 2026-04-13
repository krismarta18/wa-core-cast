"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import { navigation } from "@/lib/navigation";
import { MessageSquareMore, ExternalLink, LogOut } from "lucide-react";
import { useAuth } from "@/providers/auth-provider";

export function Sidebar({ onNavigate }: { onNavigate?: () => void }) {
  const pathname = usePathname();
  const { session, logout } = useAuth();

  return (
    <aside className="flex h-screen w-64 flex-col border-r border-gray-200 bg-white">
      {/* Logo */}
      <div className="flex h-16 items-center gap-2 border-b border-gray-200 px-6">
        <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-green-600">
          <MessageSquareMore className="h-4 w-4 text-white" />
        </div>
        <span className="text-lg font-bold text-gray-900">WACAST</span>
      </div>

      {/* Navigation */}
      <nav className="flex-1 overflow-y-auto px-3 py-4">
        {navigation.map((group, i) => (
          <div key={i} className="mb-5">
            {group.title && (
              <p className="mb-1 px-3 text-xs font-semibold uppercase tracking-wider text-gray-400">
                {group.title}
              </p>
            )}
            <ul className="space-y-0.5">
              {group.items.map((item) => {
                const isActive = item.external
                  ? false
                  : item.href === "/"
                  ? pathname === "/"
                  : pathname.startsWith(item.href);

                if (item.external) {
                  return (
                    <li key={item.href}>
                      <a
                        href={item.href}
                        target="_blank"
                        rel="noopener noreferrer"
                        onClick={onNavigate}
                        className="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900"
                      >
                        <item.icon className="h-4 w-4 flex-shrink-0 text-gray-400" />
                        {item.label}
                        <ExternalLink className="ml-auto h-3 w-3 text-gray-300" />
                      </a>
                    </li>
                  );
                }

                return (
                  <li key={item.href}>
                    <Link
                      href={item.href}
                      onClick={onNavigate}
                      className={cn(
                        "flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                        isActive
                          ? "bg-green-50 text-green-700"
                          : "text-gray-600 hover:bg-gray-100 hover:text-gray-900"
                      )}
                    >
                      <item.icon
                        className={cn(
                          "h-4 w-4 flex-shrink-0",
                          isActive ? "text-green-600" : "text-gray-400"
                        )}
                      />
                      {item.label}
                      {item.badge && (
                        <span className="ml-auto rounded-full bg-green-100 px-2 py-0.5 text-xs font-semibold text-green-700">
                          {item.badge}
                        </span>
                      )}
                    </Link>
                  </li>
                );
              })}
            </ul>
          </div>
        ))}
      </nav>

      {/* Footer */}
      <div className="border-t border-gray-200 px-6 py-4">
        {session && (
          <div className="mb-3 rounded-lg border border-gray-100 bg-gray-50 px-3 py-2">
            <p className="text-sm font-semibold text-gray-900">{session.user.full_name}</p>
            <p className="text-xs text-gray-500">{session.user.phone_number}</p>
          </div>
        )}
        <button
          type="button"
          onClick={() => {
            onNavigate?.();
            void logout();
          }}
          className="mb-3 flex w-full items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900"
        >
          <LogOut className="h-4 w-4 text-gray-400" />
          Keluar
        </button>
        <p className="text-xs text-gray-400">WACAST Core v1.0.0</p>
      </div>
    </aside>
  );
}
