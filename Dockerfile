FROM golang:1.10.3-alpine as builder
WORKDIR /go/src/github.com/apanagiotou/go-kafka-to-s3/
RUN apk add librdkafka-dev build-base
COPY . .
RUN GOOS=linux go build -a -o go-kafka-to-s3 .
ENTRYPOINT ["./go-kafka-to-s3"]
