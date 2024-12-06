FROM alpine:3.20
WORKDIR /app

COPY dist/armadan .
COPY web/static ./web/static

RUN chmod +x armadan

RUN adduser -D -u 10001 appuser
USER appuser

ARG BUILD_VERSION
ENV BUILD_VERSION=${BUILD_VERSION}

ENTRYPOINT ["./armadan"]
