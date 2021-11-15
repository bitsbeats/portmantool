FROM alpine:3.14

ENV SCAN_INTERVAL=30s

COPY . /build

RUN set -e; \
	apk add --no-cache --virtual .build-deps \
		gcc \
		libc-dev \
		make \
	; \
	cd /build; \
	make; \
	make install; \
	cd /; \
	rm -r /build; \
	apk del --no-network .build-deps

RUN apk add --no-cache nmap

VOLUME /scan
WORKDIR /scan

ENTRYPOINT ["/bin/sh", "-c", "while true; do run.sh \"$@\" & sleep \"$SCAN_INTERVAL\"; wait \"$!\"; done", "loop"]
CMD ["127.0.0.1/32"]
