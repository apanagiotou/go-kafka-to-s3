# go-kafka-to-s3

go-kafka-to-s3 is a kafka consumer writen in Go. It reads messages from kafka, it saves them in a file and when the file reaches the pre-defined threshold size is uploaded to AWS S3.

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
docker run --rm --name go-kafka-to-s3 -it --env BOOTSTRAP_SERVERS=ip:port --env KAFKA_TOPIC=topic --env S3_BUCKET=your-bucket-name --env S3_REGION=us-east-1 go-kafka-to-s3
```

There are a couple of configuration variables that you can also change. Here are the default values for them:   
- KAFKA_CONSUMER_GROUP=go-kafka-to-s3   
- FILE_SIZE_THRESHOLD_MB=10   

Example:

```
docker run --rm --name go-kafka-to-s3 -it --env BOOTSTRAP_SERVERS=10.1.0.1:9092 --env KAFKA_TOPIC=bookings --env S3_BUCKET=apanagiotou-bookings --env S3_REGION=eu-west-1 go-kafka-to-s3
```

## Running the tests

```
go test ./...
```

## Authors

* **Alexandros Panagiotou** - *[@alexapanagiotou](https://twitter.com/alexapanagiotou)*