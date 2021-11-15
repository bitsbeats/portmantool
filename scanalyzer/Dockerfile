FROM golang:1.17-alpine AS builder

COPY . src/scanalyzer/
RUN cd src/scanalyzer && go build -o ../../bin/scanalyzer ./cmd/scanalyzer


FROM alpine:3.14

COPY --from=builder /go/bin/scanalyzer /usr/local/bin/

EXPOSE 4280

VOLUME /scan
WORKDIR /scan

CMD ["/usr/local/bin/scanalyzer"]
