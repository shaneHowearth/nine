FROM golang:1.13 as builder

# modules
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

# run the readarticle script
ENTRYPOINT ["./readarticle.sh"]
