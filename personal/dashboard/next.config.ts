import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: 'export',
  distDir: 'dashboard_out',
  images: {
    unoptimized: true,
  },
};

export default nextConfig;
