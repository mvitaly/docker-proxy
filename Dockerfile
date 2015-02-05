FROM golang:1.4.1

RUN mkdir -p /opt/docker-proxy
RUN go get github.com/tools/godep
COPY . /opt/docker-proxy
WORKDIR /opt/docker-proxy
RUN godep go build -o /go/bin/docker-proxy .
RUN godep go test

EXPOSE 8080
VOLUME /var/run/docker.sock
ENTRYPOINT ["docker-proxy"]
