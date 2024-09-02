# syntax=docker/dockerfile:1

FROM golang:1.23 AS builder
LABEL description="mawinter-API server"

WORKDIR /app

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,target=. \
    go mod download -x

RUN --mount=type=bind,target=. \
    go build -a -tags "netgo" -installsuffix netgo  -ldflags="-s -w -extldflags \"-static\" \
    -X main.version=$(git describe --tag --abbrev=0) \
    -X main.revision=$(git rev-list -1 HEAD) \
    -X main.build=$(git describe --tags)" \
    -o /bin/ ./...

FROM gcr.io/distroless/base-debian12
ENV TZ=Asia/Tokyo
COPY --from=builder /bin/mawinter /usr/local/bin/mawinter
ENTRYPOINT ["/usr/local/bin/mawinter", "start"]

EXPOSE 8080
