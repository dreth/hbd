FROM golang@sha256:d530cf60723abf1ffae0cd16bb546b379ec0550d2f3815e5b16ed69dbc00fba4 AS go-builder
ENV GOTOOLCHAIN=local GOMAXPROCS=4 GOFLAGS="-p=4 -mod=readonly"
RUN apk upgrade --no-cache && apk add --no-cache gcc libc-dev sqlite-dev musl-dev sqlite
WORKDIR /build
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN set -eu; mkdir /tmp/synthetic-db; \
    for migration in migrations/sqlite/*.up.sql; do sqlite3 /tmp/synthetic-db/hbd.db < "$migration"; done; \
    SQLITE3_DBNAME=/tmp/synthetic-db/hbd.db go test ./...; \
    go vet ./...; \
    CGO_ENABLED=1 go build -tags musl -ldflags "-extldflags -static -s -w" -o /build/main .

FROM node@sha256:06d7adf3e0ef61a13f1d28682726366f97a15aa1888c934d3ecec86c33918b53 AS node-base
RUN apk upgrade --no-cache && apk add --no-cache ca-certificates tzdata sqlite bash xz
ADD --checksum=sha256:228bd64aaef23270dbdbd2cf3d5ef5c467602ef38c4066246cef1984adab9a67 https://registry.npmjs.org/npm/-/npm-12.2.1.tgz /tmp/npm.tgz
RUN mkdir /tmp/npm-security && tar -xzf /tmp/npm.tgz -C /tmp/npm-security && \
    rm -rf /usr/local/lib/node_modules/npm && mv /tmp/npm-security/package /usr/local/lib/node_modules/npm && \
    test "$(npm --version)" = 12.2.1

FROM node-base AS frontend-build
ENV NEXT_TELEMETRY_DISABLED=1
WORKDIR /build/frontend
COPY frontend/package.json frontend/package-lock.json ./
COPY frontend/security ./security
RUN npm ci && npm run security:audit
COPY frontend/ ./
RUN npm run lint && npm run build

FROM node-base AS runtime
ENV NEXT_TELEMETRY_DISABLED=1 NODE_ENV=production
WORKDIR /app
ADD --checksum=sha256:5379750ed30a84bbd2e2dd74847ba6b5bd29cd0b2e3ea2ec58049b57eb2eda12 https://github.com/just-containers/s6-overlay/releases/download/v3.2.3.2/s6-overlay-noarch.tar.xz /tmp/s6-noarch.tar.xz
ADD --checksum=sha256:e6befcc96a437a3831386ecfc51808c5d3e939dc5fe3c02ae9284599e8aa2408 https://github.com/just-containers/s6-overlay/releases/download/v3.2.3.2/s6-overlay-x86_64.tar.xz /tmp/s6-amd64.tar.xz
ADD --checksum=sha256:b17f17a82e7a515c682a91edaf2ffdabb73f891981b6c1fd712115693a2f8b4c https://github.com/just-containers/s6-overlay/releases/download/v3.2.3.2/s6-overlay-aarch64.tar.xz /tmp/s6-arm64.tar.xz
RUN set -eu; case "$(uname -m)" in x86_64) arch=amd64;; aarch64) arch=arm64;; *) exit 1;; esac; \
    tar -C / -Jxpf /tmp/s6-noarch.tar.xz; tar -C / -Jxpf /tmp/s6-${arch}.tar.xz
COPY --from=go-builder /build/main ./backend/main
COPY backend/migrations ./backend/migrations
COPY --from=frontend-build /build/frontend/.next/standalone ./frontend
COPY --from=frontend-build /build/frontend/.next/static ./frontend/.next/static
COPY --from=frontend-build /build/frontend/public ./frontend/public
COPY services/backend/run /etc/services.d/backend/run
COPY services/frontend/run /etc/services.d/frontend/run
COPY services/frontend/finish /etc/services.d/frontend/finish
RUN mkdir -p /app/data && chmod +x /etc/services.d/backend/run /etc/services.d/frontend/run /etc/services.d/frontend/finish
EXPOSE 8418 8417
ENTRYPOINT ["/init"]
