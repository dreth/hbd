import { withSentryConfig } from '@sentry/nextjs/config';

/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  // Browser GlitchTip settings are inlined at build time; default them from the build's SENTRY_* values.
  env: {
    NEXT_PUBLIC_SENTRY_ENVIRONMENT:
      process.env.NEXT_PUBLIC_SENTRY_ENVIRONMENT || process.env.SENTRY_ENVIRONMENT || '',
    NEXT_PUBLIC_SENTRY_RELEASE: process.env.NEXT_PUBLIC_SENTRY_RELEASE || process.env.SENTRY_RELEASE || '',
    NEXT_PUBLIC_SENTRY_TRACES_SAMPLE_RATE:
      process.env.NEXT_PUBLIC_SENTRY_TRACES_SAMPLE_RATE || process.env.SENTRY_TRACES_SAMPLE_RATE || '',
  },
};

// Self-hosted GlitchTip: no build telemetry to sentry.io, no release API calls, no source map upload.
export default withSentryConfig(nextConfig, {
  telemetry: false,
  silent: true,
  sourcemaps: { disable: true },
  release: { create: false, finalize: false },
});
