FROM golang:1.24-alpine AS builder

RUN apk update && apk add --no-cache sqlite-dev build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1
RUN go build -o armadan ./cmd/armadan/main.go

FROM alpine:3.21

RUN apk add --no-cache sqlite-libs

WORKDIR /app

COPY --from=builder /app/armadan .
COPY --from=builder /app/web/static ./web/static

RUN adduser -D -u 10001 appuser
USER appuser

ARG BUILD_VERSION
ENV BUILD_VERSION=${BUILD_VERSION}

ENTRYPOINT ["./armadan"]
