FROM golang:1.4.2

RUN mkdir -p /opt/docker-proxy
RUN go get github.com/tools/godep

COPY ./Godeps /go/src/github.com/xbudex/docker-proxy/Godeps
WORKDIR /go/src/github.com/xbudex/docker-proxy
RUN godep restore

EXPOSE 8080

COPY ./docker /go/src/github.com/xbudex/docker-proxy/docker
COPY ./proxy /go/src/github.com/xbudex/docker-proxy/proxy
COPY ./cmd /go/src/github.com/xbudex/docker-proxy/cmd

RUN go test ./...
RUN go build -o /go/bin/docker-proxy \
    github.com/xbudex/docker-proxy/cmd/docker-proxy

ENTRYPOINT ["docker-proxy"]
