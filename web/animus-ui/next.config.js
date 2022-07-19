module.exports = {
  reactStrictMode: true,
  env: {
    NEXT_API_URL: process.env.NEXT_API_URL,
    NEXT_IPFS_DEFAULT_GATEWAY: process.env.NEXT_IPFS_DEFAULT_GATEWAY,
  },
  eslint: {
    // Warning: This allows production builds to successfully complete even if
    // your project has ESLint errors.
    ignoreDuringBuilds: true,
  },
};
