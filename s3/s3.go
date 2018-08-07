package s3

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/pkg/errors"
)

// Uploader uploads a file to AWS S3
type Uploader struct {
	s3Bucket  string
	s3manager s3manageriface.UploaderAPI
}

// New makes a new S3Uploader. Return nil if error occurs during setup.
func New(s3Bucket string, s3manager *s3manager.Uploader) *Uploader {
	u := &Uploader{s3Bucket: s3Bucket, s3manager: s3manager}
	return u
}

// Upload will upload a single file to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func (u *Uploader) Upload(filename string) error {

	file, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err, "OS File Open")
	}
	defer file.Close()

	log.Infof("Uploading %s to S3...\n", filename)
	result, err := u.s3manager.Upload(&s3manager.UploadInput{
		Bucket: aws.String(u.s3Bucket),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		return errors.Wrap(err, "S3 upload")
	}
	log.Infof("Successfully uploaded %s to %s\n", filename, result.Location)
	return nil
}
