# syntax=docker/dockerfile:1

# ---- Build stage ----
# Run the compiler on the NATIVE build arch (BUILDPLATFORM) so it never runs
# under cross-arch emulation, then cross-compile to the target arch via GOARCH.
# This avoids the QEMU pointer-packing crash (fatal error: lfstack.push) that
# hits the Go runtime when the amd64 compiler is emulated on an arm64 host.
FROM --platform=$BUILDPLATFORM golang:1.22 AS build

WORKDIR /src

# Cache module downloads first (only re-runs when go.mod/go.sum change).
COPY go.mod ./
RUN go mod download

# Copy source and build a static binary.
COPY . .

# VERSION is injected at build time and baked into the binary via -ldflags.
# TARGETOS/TARGETARCH are provided automatically by BuildKit/Buildah from the
# --platform flag, so the binary targets the requested platform.
ARG VERSION=dev
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build \
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
