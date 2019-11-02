FROM golang:1.13 as builder

# modules
WORKDIR $GOPATH/src/github.com/shanehowearth/nine/REST
ADD . $GOPATH/src/github.com/shanehowearth/nine/REST

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# build time
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /go/bin/REST main.go routes.go

FROM scratch
WORKDIR /root/
COPY --from=0 /go/bin/REST .

# run options
ENTRYPOINT ["./REST"]
