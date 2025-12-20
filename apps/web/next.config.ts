import type { NextConfig } from 'next'

const nextConfig: NextConfig = {
  experimental: {
    turbopack: {
      // This tells Next.js exactly where the root is
      root: '.',
    },
  },
  reactCompiler: true,
}

export default nextConfig
