# Builder
FROM  golang:1.22.6-alpine3.20 AS builder

COPY ../ /go/src

RUN CGO_ENABLED=0 go build -C /go/src/cmd/ -o /cmd/shadowlog

FROM busybox:musl AS busybox

FROM debian:stable-slim AS certs
RUN apt update && apt install -y ca-certificates

FROM scratch AS base

COPY --from=busybox /etc/passwd /etc/passwd
COPY --from=busybox /etc/group /etc/group
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /app/ssl/certs/

RUN --mount=from=busybox,dst=/usr/ ["busybox", "sh", "-c", "mkdir -p /app && chmod 777 /app"]
RUN --mount=from=busybox,dst=/usr/ ["busybox", "sh", "-c", "addgroup -S go -g 1000 && adduser -S go -u 1000 --ingroup go --disabled-password"]

ENV HOME=/app
ENV USER=go
ENV PATH=/usr/local/bin:/app
ENV SSL_CERT_DIR=/app/ssl/certs

FROM base AS server

USER go

COPY --from=builder /cmd/shadowlog /app/shadowlog

ENTRYPOINT [ "/app/shadowlog" ]

CMD ["run"]
