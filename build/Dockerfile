FROM golang:1.17.5-alpine3.14 as builder
LABEL description="mawinter-API server"

COPY cmd/ /app/cmd/
COPY internal/ /app/internal/
COPY go.mod /app/go.mod
WORKDIR /app
RUN go mod tidy && go build -a -tags "netgo" -installsuffix netgo  -ldflags="-s -w -extldflags \"-static\"" -o /app/cmd/mawinter/mawinter /app/cmd/mawinter/main.go

FROM alpine:3.14
COPY --from=builder /app/cmd/mawinter/mawinter /mawinter
RUN apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
ENTRYPOINT ["/mawinter"]

EXPOSE 80
