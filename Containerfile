FROM docker.io/golang:1.25.0@sha256:91e2cd436f7adbfad0a0cbb7bf8502fa863ed8461414ceebe36c6304731e0fd9 as builder
WORKDIR /app

ARG CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd

RUN go build -o ./build/main ./cmd/...

# ---

FROM ghcr.io/markormesher/scratch:v0.1.0
WORKDIR /app

LABEL image.registry=ghcr.io
LABEL image.name=markormesher/cloudflare-dns-updater

COPY --from=builder /app/build/main /usr/local/bin/cloudflare-dns-updater

CMD ["/usr/local/bin/cloudflare-dns-updater"]
