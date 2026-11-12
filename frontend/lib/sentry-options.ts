import type { init } from "@sentry/nextjs";

type InitOptions = NonNullable<Parameters<typeof init>[0]>;

export const DEFAULT_TRACES_SAMPLE_RATE = 0.2;

/**
 * SDK v11 collects everything by default; this is the equivalent of `send_default_pii` off:
 * no user/IP inference, cookies, bodies, query strings or local variables, and only a few
 * harmless request headers.
 */
export const PRIVATE_DATA_COLLECTION: InitOptions["dataCollection"] = {
  userInfo: false,
  cookies: false,
  httpHeaders: { request: { allow: ["user-agent", "content-type", "accept", "host"] }, response: false },
  httpBodies: [],
  urlQueryParams: false,
  databaseQueryData: false,
  stackFrameVariables: false,
};

export function parseTracesSampleRate(value: string | undefined) {
  if (value === undefined || value.trim() === "") return DEFAULT_TRACES_SAMPLE_RATE;
  const parsed = Number(value);
  return Number.isFinite(parsed) && parsed >= 0 && parsed <= 1 ? parsed : DEFAULT_TRACES_SAMPLE_RATE;
}

/** Shared GlitchTip options for the browser and the Next.js server; undefined without a DSN. */
export function buildSentryOptions(config: {
  dsn: string | undefined;
  environment?: string;
  release?: string;
  tracesSampleRate?: string;
  service: string;
}) {
  const dsn = config.dsn?.trim();
  if (!dsn) return undefined;

  return {
    dsn,
    environment: config.environment?.trim() || "production",
    release: config.release?.trim() || undefined,
    tracesSampleRate: parseTracesSampleRate(config.tracesSampleRate),
    // GlitchTip ingests transaction events but ignores v11's default streamed `span` items.
    traceLifecycle: "static",
    dataCollection: PRIVATE_DATA_COLLECTION,
    initialScope: { tags: { service: config.service } },
    beforeSendLog: (log) => ({ ...log, attributes: { "service.name": config.service, ...log.attributes } }),
  } satisfies InitOptions;
}
