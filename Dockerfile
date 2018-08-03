FROM golang:1.10.3
WORKDIR /go/src/github.com/docker-services/kafka-s3-connector/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kafka-s3-connector .

FROM scratch
COPY --from=0 /go/src/github.com/docker-services/kafka-s3-connector .
ENTRYPOINT ["./kafka-s3-connector"]
