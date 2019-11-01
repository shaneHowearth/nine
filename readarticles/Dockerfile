FROM golang:1.13 as builder

# modules
ENV GO111MODULE=on
WORKDIR $GOPATH/src/readarticles
ADD . $GOPATH/src/readarticles

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

# build time
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /go/bin/readarticle cmd/main.go

# stage 2
FROM redis
WORKDIR /root/
COPY --from=0 /go/bin/readarticle .
COPY . .
COPY readarticle.sh .
RUN chmod 777 readarticle.sh

# install and run redis (in the background)
# RUN apt-get update && apt-get -y install redis-server

ENTRYPOINT ["./readarticle.sh"]
