FROM golang:1.13 as builder

# modules
WORKDIR $GOPATH/src/github.com/shanehowearth/nine/createarticle
ADD . $GOPATH/src/github.com/shanehowearth/nine/createarticle

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# build time
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /go/bin/createarticle cmd/main.go

FROM scratch
WORKDIR /root/
COPY --from=0 /go/bin/createarticle .
# run options
ENV MQURI="amqp://guest:guest@rabbitmq-server:5672/"
ENV DBURI="postgres://articlewriter:hu8jmn3@articledb/article_postgres_db?sslmode=disable"
ENV PORT_NUM=5100
EXPOSE 5100
ENTRYPOINT ["./createarticle"]
