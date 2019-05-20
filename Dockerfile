FROM golang:1.12
WORKDIR /go/src/github.com/carlpett/stream_exporter/
COPY . .
RUN make build

FROM busybox:glibc
WORKDIR /usr/bin/
COPY --from=0 /go/src/github.com/carlpett/stream_exporter/stream_exporter /usr/bin/stream_exporter
USER 1000
ENTRYPOINT ["/usr/bin/stream_exporter"]
EXPOSE 9178