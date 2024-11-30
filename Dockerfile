FROM alpine:3.20 AS base
WORKDIR /app

FROM base AS production
COPY dist/armadan .
COPY web/static ./web/static

RUN chmod +x armadan

RUN adduser -D -u 10001 appuser
USER appuser

ENTRYPOINT ["./armadan"]
