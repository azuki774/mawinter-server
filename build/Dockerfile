FROM golang:1.19.1-alpine3.16 as builder
LABEL description="mawinter-API server"

COPY cmd/ /app/cmd/
COPY internal/ /app/internal/
COPY go.mod /app/go.mod
COPY .git /app/.git
WORKDIR /app
RUN go mod tidy && apk add git
RUN go build -a -tags "netgo" -installsuffix netgo -ldflags="-s -w -extldflags \"-static\" \
    -X main.version=$(git describe --tag --abbrev=0) \
    -X main.revision=$(git rev-list -1 HEAD) \
    -X main.build=$(git describe --tags)" \
    -o bin/ ./...

FROM alpine:3.16
COPY --from=builder /app/bin/mawinter /usr/local/bin/mawinter
RUN apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
ENTRYPOINT ["/usr/local/bin/mawinter", "start"]

EXPOSE 8080
