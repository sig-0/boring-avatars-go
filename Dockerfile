# Build
ARG GO_VERSION=1.24
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-bookworm AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

# Copy the rest
COPY . .

# Build the binary
ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=0
RUN GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -trimpath -ldflags="-s -w" -o /build/boring-avatars ./cmd

# Runtime
FROM debian:bookworm-slim

WORKDIR /app

# Only the binary
COPY --from=build /build/boring-avatars /app/boring-avatars

# Server listens on 8545 by default
EXPOSE 8545

ENTRYPOINT ["./boring-avatars", "serve"]