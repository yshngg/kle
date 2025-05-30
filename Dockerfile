# ---------- Builder Stage ----------
FROM --platform=${BUILDPLATFORM} golang:1.24-alpine AS builder

LABEL maintainer="Yusheng Guo <yshngg@outlook.com>"
LABEL org.opencontainers.image.source="https://github.com/yshngg/kle"
LABEL org.opencontainers.image.description="A Kubernetes Leader Election Demo."
LABEL org.opencontainers.image.licenses=MIT

ARG TARGETOS
ARG TARGETARCH

# Disable cgo for static build
ENV CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    LDFLAG_LOCATION=github.com/yshngg/kle/pkg/version

# Use BuildKit cache mounts for Go modules and build cache
WORKDIR /app

# install git for metadata and clean up
RUN apk add --no-cache git \
    && rm -rf /var/cache/apk/*

# Copy only go.mod & go.sum to leverage layer caching
COPY go.mod go.sum ./

# Download dependencies with cache
RUN go mod download

# copy sources
COPY . .

# Embed version metadata and build binary
ARG VERSION="$(git describe --tags --always)"
RUN set -eux; \
    VERSION="${VERSION}"; \
    BRANCH="$(git rev-parse --abbrev-ref HEAD)"; \
    SHA1="$(git rev-parse HEAD)"; \
    BUILD="$(date -u +%FT%T%z)"; \
    go build -a -o /bin/kle \
    -ldflags "-s -w \
    -X ${LDFLAG_LOCATION}.version=${VERSION} \
    -X ${LDFLAG_LOCATION}.buildDate=${BUILD} \
    -X ${LDFLAG_LOCATION}.gitbranch=${BRANCH} \
    -X ${LDFLAG_LOCATION}.gitsha1=${SHA1}" \
    github.com/yshngg/kle/cmd

# ---------- Runtime Stage ----------
FROM alpine:3.21 AS runner

COPY --from=builder /bin/kle /kle

# Use non-root UID/GID
USER nobody:nogroup
WORKDIR /

EXPOSE 2190
ENTRYPOINT ["/kle"]
