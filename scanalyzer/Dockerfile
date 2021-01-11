FROM golang:1.15 AS builder

COPY . src/scanalyzer/
RUN cd bin && go build scanalyzer


FROM alpine:3.12

COPY --from=builder /go/bin/scanalyzer /usr/local/bin/

VOLUME /scan

WORKDIR /scan
CMD ["/bin/sh", "-c", "while /usr/local/bin/scanalyzer; do /bin/true; done"]
