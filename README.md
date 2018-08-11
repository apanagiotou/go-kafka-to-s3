# go-kafka-to-s3

go-kafka-to-s3 is a kafka consumer written in Go. Reads messages from a topic, gathers them in batches (files) of a pre-defined size and stores them in S3. It's useful if you want to have your data in a persistent storage for later process.

## Getting Started


These instructions will get you a copy of the project up and running.

### Installing

A step by step series of examples that tell you how to get a development env running

First build the docker image

```
git clone https://github.com/apanagiotou/go-kafka-to-s3
cd go-kafka-to-s3
docker build -t go-kafka-to-s3 .
```

The run the docker container. You need (at least) one kafka broker, an AWS key and secret and, an S3 bucket

```
docker run --rm --name go-kafka-to-s3 -it \
    -e BOOTSTRAP_SERVERS=ip:port \
    -e KAFKA_TOPIC=topic \
    -e S3_BUCKET=your-bucket-name \
    -e S3_REGION=eu-west-1 \
    -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
    -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
    go-kafka-to-s3
```

There are a couple of configuration variables that you can also change. Here are the default values for them:   
- KAFKA_CONSUMER_GROUP=go-kafka-to-s3   
- S3_BUCKET_SUBPATH=your-bucket-folder
- FILE_SIZE_THRESHOLD_MB=10   

Example:

```
docker run --rm --name go-kafka-to-s3 -it \
    -e BOOTSTRAP_SERVERS=10.1.0.1:9092 \
    -e KAFKA_TOPIC=bookings \
    -e S3_BUCKET=apanagiotou-bookings \
    -e S3_BUCKET_SUBPATH=test-folder \
    -e S3_REGION=eu-west-1 \
    -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
    -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
    go-kafka-to-s3
```

## Running the tests

```
go test ./...
```

## Authors

* **Alexandros Panagiotou** - *[apanagiotou.com](https://apanagiotou.com)* / Twitter: *[@alexapanagiotou](https://twitter.com/alexapanagiotou)*
