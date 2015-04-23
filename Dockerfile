FROM golang:1.4.2

RUN go get github.com/golang/lint/golint
RUN go get github.com/tools/godep
RUN go get golang.org/x/tools/cmd/cover
RUN go get golang.org/x/tools/cmd/godoc
RUN go get golang.org/x/tools/cmd/vet

RUN mkdir -p /opt/docker-proxy
COPY ./Godeps /go/src/github.com/xbudex/docker-proxy/Godeps
WORKDIR /go/src/github.com/xbudex/docker-proxy
RUN godep restore

COPY ./docker /go/src/github.com/xbudex/docker-proxy/docker
COPY ./proxy /go/src/github.com/xbudex/docker-proxy/proxy
COPY ./cmd /go/src/github.com/xbudex/docker-proxy/cmd

#RUN go test ./...
RUN go vet ./...
RUN golint ./...
RUN go build -o /go/bin/docker-proxy \
    github.com/xbudex/docker-proxy/cmd/docker-proxy

EXPOSE 8080
ENTRYPOINT ["docker-proxy"]
