FROM docker.io/golang:1.25.7@sha256:011d6e21edbc198b7aeb06d705f17bc1cc219e102c932156ad61db45005c5d31 as builder
WORKDIR /app

ARG CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd

RUN go build -o ./build/main ./cmd/...

# ---

FROM ghcr.io/markormesher/scratch:v0.4.12@sha256:f8ec68ff0857514cedc0cccff097ba234c4a05bff884d0d42de1d0ce630e1829
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