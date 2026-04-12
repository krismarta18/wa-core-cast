import { Sidebar } from "@/components/layout/sidebar";
import { Breadcrumb } from "@/components/layout/breadcrumb";
import { MobileSidebarTrigger } from "@/components/layout/mobile-sidebar-trigger";
import { MobileFAB } from "@/components/layout/mobile-fab";
import { ToastProvider } from "@/components/ui/toast";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <ToastProvider>
      <div className="flex h-full min-h-screen">
        {/* Desktop sidebar */}
        <div className="hidden lg:flex lg:flex-shrink-0">
          <Sidebar />
        </div>

        {/* Main content */}
        <div className="flex flex-1 flex-col overflow-hidden">
          {/* Mobile top bar */}
          <div className="flex items-center gap-2 border-b border-gray-200 bg-white px-4 py-2 lg:hidden">
            <MobileSidebarTrigger />
            <span className="text-sm font-semibold text-gray-800">WACAST</span>
          </div>
          <Breadcrumb />
          <main className="flex-1 overflow-y-auto bg-gray-50">{children}</main>
          <MobileFAB />
        </div>
      </div>
    </ToastProvider>
  );
}
