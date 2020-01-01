FROM golang as builder
WORKDIR /go/src/github.com/samos123/k8s-qos
ADD pkg pkg
ADD cmd cmd
RUN go get -d ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o k8s-qos cmd/k8s-qos/main.go
ADD https://download.docker.com/linux/static/stable/x86_64/docker-18.09.9.tgz /tmp/
RUN tar -xzf /tmp/docker-18.09.9.tgz --directory /tmp/
RUN ls /tmp/ && ls /tmp/docker/
FROM alpine:3.11
COPY tools/getveth.sh /usr/local/bin/
COPY --from=builder /go/src/github.com/samos123/k8s-qos/k8s-qos /usr/local/bin/
COPY --from=builder /tmp/docker/docker /usr/local/bin/
RUN apk add --no-cache iproute2 util-linux bash
CMD ["/usr/local/bin/k8s-qos"]
