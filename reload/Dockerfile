# Arguments
ARG GOLANG_VERSION
ARG ALPINE_VERSION

# Build part
FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS builder

RUN apk --no-cache add \
    make \
    git && \
    adduser -D -h /tmp/build build

USER build
WORKDIR /tmp/build

COPY --chown=build Makefile Makefile
COPY --chown=build go.mod go.mod
COPY --chown=build go.sum go.sum
COPY --chown=build vendor vendor 

ARG VERSION
ARG NAME

COPY --chown=build main.go main.go

RUN make build

# Exec part
FROM alpine:3.11.3

ARG VERSION
ARG NAME

# Appication Configuration
ENV CONFIG_PATH="/config"

# Copy from builder
COPY --from=builder /tmp/build/${NAME}-${VERSION} /usr/local/bin/app

# Exec
CMD ["app"]
