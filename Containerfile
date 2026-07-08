FROM docker.io/golang:1.26.5@sha256:63f132d58c1f589f0dcda584933a9bb44bfda1150f1506377f5a902f34d86033 as builder
WORKDIR /app

ARG CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd

RUN go build -o ./build/main ./cmd/...

# ---

FROM ghcr.io/markormesher/scratch:v0.4.23@sha256:2953f9c23784e8f1ce42bd711363f205226949df120a1d16811fcb23e4a59139
WORKDIR /app

COPY --from=builder /app/build/main /usr/local/bin/cloudflare-dns-updater

CMD ["/usr/local/bin/cloudflare-dns-updater"]

LABEL image.name=markormesher/cloudflare-dns-updater
LABEL image.registry=ghcr.io
LABEL org.opencontainers.image.description=""
LABEL org.opencontainers.image.documentation=""
LABEL org.opencontainers.image.title="cloudflare-dns-updater"
LABEL org.opencontainers.image.url=""
LABEL org.opencontainers.image.vendor=""
LABEL org.opencontainers.image.version=""
