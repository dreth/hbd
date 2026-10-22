# Security upgrade, 23 October 2026

The application uses Go 1.27.3, Node 24.22.1, npm 12.2.1, Next 16.4.2, React 19.3.1 and S6 3.2.3.2. Both Dockerfiles pin base images and downloaded artifacts. Go modules and npm dependencies are frozen. The image build runs backend tests against a synthetic SQLite schema, vet, the frontend dependency audit, lint and production build. Runtime frontend files use Next standalone output; build-only dependencies stay in the build stage. Existing API/frontend ports and the local Nginx proxy remain available.

## Migrations assessed

Node 20 is unsupported. Node 24 is the supported LTS deployment branch. Next 14 and React 18 move together to current stable Next 16/React 19. Webpack is retained. DayPicker 8 to 10 and date-fns 3 to 4 require the updated navigation/class mapping. Tailwind stays on compatible version 3; its selector parser is overridden to the patched compatible API. ESLint 10 uses the current Next flat config; upstream React/A11y plugins still emit older peer-range warnings, but the complete project lint and build pass. These build tools are absent from the standalone runtime application.

AWS SDK for Go v1 is unsupported. SDK v2 preserves the existing four S3 environment variables, bucket/object names, content type, body bytes and content length. Errors propagate and requests have a 60-second timeout. Unit tests exercise byte/metadata preservation and failures without connecting to a provider. Real S3 uploads and real Telegram delivery are not test actions.

SQLBoiler moves to the maintained aarondl module, version 4.19.7. The SQLite schema and nullable public model ID types stay unchanged. A checked-in generator template corrects a pre-existing scalar-key panic: the local key's type determines whether it can be nil, independently of the referenced primary key's nullable representation. Regenerate with `cd backend && ./security/regenerate-models.sh`. This uses a new temporary database, never the deployed database. Generated relationship tests pass.

## Braces security backport

Published braces 3.0.3 has no upstream fix for GHSA-vfj7-8cjw-p6xm. The postinstall script verifies the exact original parser SHA-256, applies a depth limit of 100 to braces and parentheses, then verifies the patched SHA-256. The audit checks every installed copy and ordinary/escaped/quoted and boundary/deep parsing behavior regardless of feed contents. It accepts advisory metadata only when every dependency path resolves solely to this attested backport; all other findings fail CI.

## Verification and rollback

The 23 October server rehearsal restored the full application volume with outbound networking disabled. All original account and birthday row hashes were preserved. Native registration/password login, rejected wrong credentials and invalid JWT, credential encryption/decryption, timezones, birthday CRUD and frontend routes passed. Chrome 156 exercised password login, leap-day creation, the migrated sorting control, theme switching and logout/protected redirect with no uncaught JavaScript errors. No real provider notifications or uploads were made. The image scanner found no reported high/critical or other rated findings in the HBD runtime image (one unknown metadata finding).

Before deployment retain the exact previous image/configuration and a coherent full data-volume backup. Backups, environment secrets and test fixtures stay in protected server directories and are never committed. This release adds no database migration. Keep the current deployment's host networking/proxy, environment, mounts and persistent data configuration. The working server's pre-existing Compose overrides must not be reset by a Git checkout. CI deploys only a successful exact main commit and preserves local Compose configuration with a fast-forward merge.
