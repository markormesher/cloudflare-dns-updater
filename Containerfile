FROM docker.io/golang:1.25.6@sha256:ce63a16e0f7063787ebb4eb28e72d477b00b4726f79874b3205a965ffd797ab2 as builder
WORKDIR /app

ARG CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd

RUN go build -o ./build/main ./cmd/...

# ---

FROM ghcr.io/markormesher/scratch:v0.4.10@sha256:50e90f252c2c5282a4e4895274089ce3b349fb10e77a517fd05721ca4ae1bbe2
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