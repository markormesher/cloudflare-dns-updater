FROM docker.io/golang:1.24.1@sha256:52ff1b35ff8de185bf9fd26c70077190cd0bed1e9f16a2d498ce907e5c421268 as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd

RUN go build -o ./build/main ./cmd/...

# ---

FROM gcr.io/distroless/base-debian12@sha256:125eb09bbd8e818da4f9eac0dfc373892ca75bec4630aa642d315ecf35c1afb7
WORKDIR /app

LABEL image.registry=ghcr.io
LABEL image.name=markormesher/cloudflare-dns-updater

COPY --from=builder /app/build/main /usr/local/bin/cloudflare-dns-updater

CMD ["/usr/local/bin/cloudflare-dns-updater"]
