import {
  LayoutDashboard,
  Smartphone,
  Wifi,
  QrCode,
  Info,
  Power,
  MessageSquare,
  Send,
  FileClock,
  PlusSquare,
  CalendarClock,
  KeyRound,
  Webhook,
  FileText,
  Bot,
  Tag,
  BarChart2,
  AlertTriangle,
  BookUser,
  Users,
  ShieldOff,
  UserCircle,
  CreditCard,
  Bell,
  Rocket,
  Flame,
  type LucideIcon,
} from "lucide-react";

export interface NavItem {
  label: string;
  href: string;
  icon: LucideIcon;
  badge?: string;
  external?: boolean;
}

export interface NavGroup {
  title?: string;
  items: NavItem[];
}

export const navigation: NavGroup[] = [
  {
    items: [
      { label: "Dashboard", href: "/", icon: LayoutDashboard },
    ],
  },
  {
    title: "Device Management",
    items: [
      { label: "Session Management", href: "/devices/status", icon: Wifi },
      { label: "Connect New Device", href: "/devices/qr", icon: QrCode },
      { label: "Account Warming", href: "/warming", icon: Flame, badge: "Pro" },
      { label: "Multi Device Info", href: "/devices/info", icon: Info },
    ],
  },
  {
    title: "Messaging",
    items: [
      { label: "New Message", href: "/messaging/new", icon: PlusSquare },
      { label: "Broadcast", href: "/messaging/broadcast", icon: Send },
      { label: "Scheduled", href: "/messaging/scheduled", icon: CalendarClock },
      { label: "Message Logs", href: "/messaging/logs", icon: FileClock },
    ],
  },
  {
    title: "API & Integration",
    items: [
      { label: "API Keys", href: "/api-integration/keys", icon: KeyRound },
      { label: "Webhook Settings", href: "/api-integration/webhooks", icon: Webhook },
      { label: "Documentation", href: "/api-integration/docs", icon: FileText },
    ],
  },
  {
    title: "Auto Response & Template",
    items: [
      { label: "Keyword", href: "/auto-response/keywords", icon: Tag },
      { label: "Message Template", href: "/auto-response/templates", icon: Bot },
    ],
  },
  {
    title: "Monitoring & Analytics",
    items: [
      { label: "Usage Statistics", href: "/monitoring/usage", icon: BarChart2 },
      { label: "Failure Rate", href: "/monitoring/failure", icon: AlertTriangle },
    ],
  },
  {
    title: "Contact Management",
    items: [
      { label: "Phone Book", href: "/contacts/phonebook", icon: BookUser },
      { label: "Group Contact", href: "/contacts/groups", icon: Users },
      { label: "Blacklist / Block", href: "/contacts/blacklist", icon: ShieldOff },
    ],
  },
  {
    title: "Settings",
    items: [
      { label: "Quick Start", href: "/settings/onboarding", icon: Rocket, badge: "New" },
      { label: "Profil & Akun", href: "/settings/profile", icon: UserCircle },
    ],
  },
];

