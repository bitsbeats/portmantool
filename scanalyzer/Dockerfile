FROM golang:1.17 AS builder

COPY . src/scanalyzer/
RUN cd src/scanalyzer && go build -o ../../bin/scanalyzer ./cmd/scanalyzer


FROM alpine:3.14

COPY --from=builder /go/bin/scanalyzer /usr/local/bin/

RUN mkdir -p /scan && chown nobody:nobody /scan

EXPOSE 4280

USER nobody
VOLUME /scan
WORKDIR /scan

CMD ["/usr/local/bin/scanalyzer"]
