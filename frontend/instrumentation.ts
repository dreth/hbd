import * as Sentry from "@sentry/nextjs";

import { buildSentryOptions } from "./lib/sentry-options";

// Next.js server GlitchTip reporting. It uses the frontend DSN baked in at build time, because the
// container's runtime SENTRY_DSN belongs to the Go API (which may report to its own project).
export function register() {
  const options = buildSentryOptions({
    dsn: process.env.NEXT_PUBLIC_SENTRY_DSN,
    environment: process.env.SENTRY_ENVIRONMENT || process.env.NODE_ENV,
    release: process.env.SENTRY_RELEASE || process.env.NEXT_PUBLIC_SENTRY_RELEASE,
    tracesSampleRate: process.env.SENTRY_TRACES_SAMPLE_RATE,
    service: "hbd-web",
  });
  if (options) Sentry.init(options);
}

export const onRequestError = Sentry.captureRequestError;
