FROM golang:1.4.1

RUN mkdir -p /opt/docker-proxy
RUN go get github.com/tools/godep

COPY ./Godeps /go/src/github.com/xbudex/docker-proxy/Godeps
WORKDIR /go/src/github.com/xbudex/docker-proxy
RUN godep restore

COPY ./docker /go/src/github.com/xbudex/docker-proxy/docker
COPY ./proxy /go/src/github.com/xbudex/docker-proxy/proxy
COPY ./cmd /go/src/github.com/xbudex/docker-proxy/cmd

RUN go test ./...
RUN go build github.com/xbudex/docker-proxy/cmd/docker-proxy

EXPOSE 8080
VOLUME /var/run/docker.sock
ENTRYPOINT ["docker-proxy"]
