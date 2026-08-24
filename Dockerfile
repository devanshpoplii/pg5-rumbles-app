# syntax=docker/dockerfile:1

# ---- Build stage ----
# Pinned by digest in production builds; tag shown here for readability.
FROM golang:1.22 AS build

WORKDIR /src

# Cache module downloads first (only re-runs when go.mod/go.sum change).
COPY go.mod ./
RUN go mod download

# Copy source and build a static binary.
COPY . .

# VERSION is injected at build time and baked into the binary via -ldflags.
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o /out/rumbles .

# ---- Runtime stage ----
# distroless: no shell, no package manager, minimal attack surface.
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /
COPY --from=build /out/rumbles /rumbles

EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["/rumbles"]
