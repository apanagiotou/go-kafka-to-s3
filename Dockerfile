FROM golang:1.10.3
WORKDIR /go/src/github.com/docker-services/go-kafka-to-s3/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-kafka-to-s3 .

FROM scratch
COPY --from=0 /go/src/github.com/docker-services/go-kafka-to-s3 .
ENTRYPOINT ["./go-kafka-to-s3"]
