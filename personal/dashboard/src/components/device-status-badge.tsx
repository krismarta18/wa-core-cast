import { Badge } from "@/components/ui/badge";
import type { DeviceStatus } from "@/lib/types";

const STATUS_MAP: Record<
  DeviceStatus,
  { label: string; variant: "success" | "warning" | "default" }
> = {
  1: { label: "Active", variant: "success" },
  2: { label: "Pending QR", variant: "warning" },
  0: { label: "Inactive", variant: "default" },
};

export function DeviceStatusBadge({ status }: { status: DeviceStatus }) {
  const { label, variant } = STATUS_MAP[status] ?? STATUS_MAP[0];
  return <Badge variant={variant}>{label}</Badge>;
}
