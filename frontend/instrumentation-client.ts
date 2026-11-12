import * as Sentry from "@sentry/nextjs";

import { buildSentryOptions } from "@/lib/sentry-options";

// Browser GlitchTip reporting. NEXT_PUBLIC_* values are fixed at image build time (see next.config.mjs);
// an empty NEXT_PUBLIC_SENTRY_DSN leaves the SDK uninitialised.
const options = buildSentryOptions({
  dsn: process.env.NEXT_PUBLIC_SENTRY_DSN,
  environment: process.env.NEXT_PUBLIC_SENTRY_ENVIRONMENT || process.env.NODE_ENV,
  release: process.env.NEXT_PUBLIC_SENTRY_RELEASE,
  tracesSampleRate: process.env.NEXT_PUBLIC_SENTRY_TRACES_SAMPLE_RATE,
  service: "hbd-browser",
});

if (options) {
  Sentry.init({
    ...options,
    integrations: [Sentry.consoleLoggingIntegration({ levels: ["warn", "error"] })],
  });
}

export const onRouterTransitionStart = Sentry.captureRouterTransitionStart;
