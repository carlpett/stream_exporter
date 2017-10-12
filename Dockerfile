FROM golang:1.8
WORKDIR /go/src/github.com/carlpett/stream_exporter/
RUN go get -u github.com/kardianos/govendor
COPY . .
RUN govendor build +p

FROM centos:centos7
WORKDIR /usr/bin/
COPY --from=0 /go/src/github.com/carlpett/stream_exporter/stream_exporter /usr/bin/stream_exporter
RUN adduser -u 1000 stream_exporter
USER stream_exporter
ENTRYPOINT ["/usr/bin/stream_exporter"]
EXPOSE 9178