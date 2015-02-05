# docker-proxy

HTTP Proxy to route multiple hosts on a single port to multiple
docker containers

# WARNING: proof of concept

This is currently a proof of concept. Expect the API to change

## Building with Docker

```bash
$ # Build docker proxy in docker
$ docker build --tag docker-proxy .
$ # Run docker-proxy from inside a docker container
$ docker run --publish 8888:8080 --volume /var/run/docker.sock:/var/run/docker.sock docker-proxy
```
