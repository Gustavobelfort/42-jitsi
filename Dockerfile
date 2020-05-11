FROM golang:1.13-alpine AS builder

COPY . /build
WORKDIR /build

RUN apk update && apk add git ca-certificates tzdata && update-ca-certificates

RUN mkdir -p ./bin
RUN go mod download && go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ./bin/ ./cmd/...

RUN adduser -D -g '' 42jitsi

FROM scratch
LABEL maintainers="Charles Labourier <pistache@42madrid.com>, Gustavo Belfort <gustavo@42sp.org.br>"

COPY --from=builder /build/bin/* /bin/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

EXPOSE 5000

USER 42jitsi
