FROM golang:1.26.0 AS builder

WORKDIR /app

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -trimpath -ldflags="-s -w" -o nanoshare ./cmd/nanoshare

FROM alpine:latest 

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -S appuser \
    && adduser -S -G appuser -H -s /sbin/nologin appuser

COPY --from=builder --chown=appuser:appuser /app/nanoshare /app/nanoshare

RUN mkdir -p /app/data && chown appuser:appuser /app/data

USER appuser

EXPOSE 8080

ENTRYPOINT [ "/app/nanoshare" ]

