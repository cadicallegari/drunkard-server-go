FROM alpine:3.7

RUN apk --no-cache update && \
    apk --no-cache add ca-certificates && \
    rm -rf /var/cache/apk/*

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY ./cmd/drunkard/drunkard /app/drunkard
COPY ./templates /templates

ENTRYPOINT ["/app/drunkard"]
