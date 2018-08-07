package main

import (
	"flag"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/apanagiotou/go-kafka-to-s3/file"
	"github.com/apanagiotou/go-kafka-to-s3/kafka"
	"github.com/apanagiotou/go-kafka-to-s3/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	kafkaBrokers       = flag.String("kafkaBrokers", getEnv("BOOTSTRAP_SERVERS", ""), "The comma separated list of brokers in the Kafka cluster")
	kafkaTopic         = flag.String("kafkaTopic", getEnv("KAFKA_TOPIC", ""), "REQUIRED: the topic to consume")
	kafkaConsumerGroup = flag.String("kafkaConsumerGroup", "go-kafka-to-s3", "The consumer group name")
	offset             = flag.String("offset", "latest", "The offset to start with. Can be `oldest`, `newest`")
	bufferSize         = flag.Int("bufferSize", 1000, "The buffer size of the message channel.")
	s3Bucket           = flag.String("s3Bucket", getEnv("S3_BUCKET", ""), "The S3 bucket to upload the files.")
	s3BucketPath       = flag.String("s3BucketPath", getEnv("S3_BUCKET_SUBPATH", "/"), "The S3 bucket to upload the files.")
	s3Region           = flag.String("s3Region", getEnv("S3_REGION", ""), "The S3 region of the bucket.")
	fileRotateSizeStr  = flag.String("fileRotateSizeStr", getEnv("FILE_SIZE_THRESHOLD_MB", "10"), "The threshold to rotate the files and upload them.")
)

func main() {
	flag.Parse()

	// Get the size in Mebabytes from the env var and convert in int64 bytes
	fileRotateSize, _ := strconv.Atoi(*fileRotateSizeStr)
	fileRotateSizeInBytes := int64(fileRotateSize * 1024 * 1024)

	// Initialize kafka and file
	kafkaConsumer := kafka.New(*kafkaBrokers, *kafkaTopic, *kafkaConsumerGroup, *offset)
	positionFile, err := file.New(*kafkaTopic, fileRotateSizeInBytes)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize S3 Uploader used to upload the rotated file to S3
	sess, err := session.NewSession(&aws.Config{Region: aws.String(*s3Region)})
	if err != nil {
		log.Fatal(err)
	}
	s3Manager := s3manager.NewUploader(sess)
	s3Uploader := s3.New(*s3Bucket, *s3BucketPath, s3Manager)

	// This channel is used to store kafka messages
	position := make(chan string, *bufferSize)

	// Writes kafka messages to the file
	go func() {
		log.Debug("Write goroutine created")
		for pos := range position {
			_, err := positionFile.Write(pos)
			if err != nil {
				log.Error(err)
			}
		}
	}()

	// Rotate and upload the file to S3 if it has reached fileRotateSize
	go func() {
		log.Debug("Rotate/Upload goroutine created")
		for {
			time.Sleep(5)
			if rotate, err := positionFile.Rotateable(); rotate == true {
				log.Debug("The file is rotatable")
				if err != nil {
					log.Error(err)
				}
				rotatedFile, err := positionFile.Rotate()
				log.Debugf("File rotated: %s", rotatedFile)
				if err != nil {
					log.Error(err)
				}
				go func() {
					log.Debug("Compress/Upload routine started")
					compressed, err := file.Compress(rotatedFile)
					log.Debugf("File compressed: %s", compressed)
					if err != nil {
						log.Error(err)
					}
					err = s3Uploader.Upload(compressed)
					if err != nil {
						log.Error(err)
					}
					log.Debugf("File uploaded: %s", compressed)
					err = os.Remove(compressed)
					if err != nil {
						log.Error(err)
					}
					log.Debugf("Rotated file deleted: %s", compressed)
				}()
			}
		}
	}()

	for {
		// ReadMessage automatically commits offsets when using consumer groups.
		msg, err := kafkaConsumer.Consume()

		if err == nil {
			position <- string(msg.Value)
		} else {
			log.Errorf("Consumer error: %v (%v)\n", err, msg)
			break
		}
	}

	kafkaConsumer.Close()
}

// getEnv takes an env and a fallback value. If the env exists it returns it.
// If not and there is a fallback it returns that,
// Else, stops the program
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if fallback != "" {
		return fallback
	}
	log.Fatal(key + " is not set")
	return ""
}
