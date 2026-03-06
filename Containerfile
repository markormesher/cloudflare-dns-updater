FROM docker.io/golang:1.26.1@sha256:e2ddb153f786ee6210bf8c40f7f35490b3ff7d38be70d1a0d358ba64225f6428 as builder
WORKDIR /app

ARG CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd

RUN go build -o ./build/main ./cmd/...

# ---

FROM ghcr.io/markormesher/scratch:v0.4.14@sha256:6212e9b38cc40ad4520f39a4043ac3f7bd12438449356311b4fa327cabce7db0
WORKDIR /app

COPY --from=builder /app/build/main /usr/local/bin/cloudflare-dns-updater

CMD ["/usr/local/bin/cloudflare-dns-updater"]

LABEL image.name=markormesher/cloudflare-dns-updater
LABEL image.registry=ghcr.io
LABEL org.opencontainers.image.description=""
LABEL org.opencontainers.image.documentation=""
LABEL org.opencontainers.image.title="cloudflare-dns-updater"
LABEL org.opencontainers.image.url="https://github.com/markormesher/cloudflare-dns-updater"
LABEL org.opencontainers.image.vendor=""
LABEL org.opencontainers.image.version=""
