FROM alpine:3.14

ENV SCAN_INTERVAL=30s

RUN apk add --no-cache curl nmap

VOLUME /scan
WORKDIR /scan

ENTRYPOINT ["/bin/sh", "-c", "while true; do run.sh \"$@\" & sleep \"$SCAN_INTERVAL\"; wait \"$!\"; done", "loop"]
CMD ["127.0.0.1/32"]
