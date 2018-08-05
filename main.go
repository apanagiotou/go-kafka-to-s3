package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/apanagiotou/go-kafka-to-s3/file"
	"github.com/apanagiotou/go-kafka-to-s3/kafka"
	"github.com/apanagiotou/go-kafka-to-s3/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {

	// Configuration values are set in env variables
	bootstrapServers := getEnv("BOOTSTRAP_SERVERS", "")
	kafkaTopic := getEnv("KAFKA_TOPIC", "")
	kafkaConsumerGroup := getEnv("KAFKA_CONSUMER_GROUP", "")
	s3Region := getEnv("S3_REGION", "")
	s3Bucket := getEnv("S3_BUCKET", "")

	fileRotateSize := int64(10000000) // 10mb

	// Initialize kafka and file
	kafkaConsumer := kafka.New(bootstrapServers, kafkaTopic, kafkaConsumerGroup, kafkaTopic)
	positionFile, _ := file.New("driver_position.log", fileRotateSize)

	// Initialize S3 Uploader used to upload the rotated file to S3
	sess, err := session.NewSession(&aws.Config{Region: aws.String(s3Region)})
	if err != nil {
		log.Fatal(err, "S3 connection")
	}
	s3Manager := s3manager.NewUploader(sess)
	s3Uploader := s3.New(s3Bucket, s3Manager)

	// This channel is used to store kafka messages
	position := make(chan string, 1000)

	// Writes kafka messages to the file
	go func() {
		for pos := range position {
			_, err := positionFile.Write(pos)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	// Rotate and upload the file to S3 if it has reached fileRotateSize
	go func() {
		for {
			time.Sleep(5)
			if rotate, err := positionFile.Rotateable(); rotate == true {
				if err != nil {
					log.Println(err)
				}
				rotatedFile, err := positionFile.Rotate()
				if err != nil {
					log.Println(err)
				}
				go func() {
					compressed, err := file.Compress(rotatedFile)
					if err != nil {
						log.Println(err)
					}
					err = s3Uploader.Upload(compressed)
					if err != nil {
						log.Println(err)
					}
					err = os.Remove(compressed)
					if err != nil {
						log.Println(err)
					}
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
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
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
