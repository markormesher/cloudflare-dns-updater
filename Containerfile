FROM docker.io/golang:1.25.1@sha256:bb979b278ffb8d31c8b07336fd187ef8fafc8766ebeaece524304483ea137e96 as builder
WORKDIR /app

ARG CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd

RUN go build -o ./build/main ./cmd/...

# ---

FROM ghcr.io/markormesher/scratch:v0.3.0@sha256:8813090ae525a56a8649fd646008af7492d8995edbeb6e1bdaa0c92773d27f8f
WORKDIR /app

LABEL image.registry=ghcr.io
LABEL image.name=markormesher/cloudflare-dns-updater

COPY --from=builder /app/build/main /usr/local/bin/cloudflare-dns-updater

CMD ["/usr/local/bin/cloudflare-dns-updater"]
