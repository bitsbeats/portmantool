FROM alpine:3.12

COPY . /build

RUN apk add --no-cache nmap

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

VOLUME /scan

WORKDIR /scan
CMD ["/usr/local/bin/runner", "30s", "nmap", "127.0.0.1/32"]
