package s3

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manageriface"
	"github.com/pkg/errors"
)

// Uploader uploads a file to AWS S3
type Uploader struct {
	s3Region  string
	s3Bucket  string
	s3manager s3manageriface.UploaderAPI
}

// New makes a new S3Uploader. Return nil if error occurs during setup.
func New(s3Region, s3Bucket string) *Uploader {
	u := &Uploader{s3Region: s3Region, s3Bucket: s3Bucket}
	u.ManageConection()
	return u
}

// ManageConection manages a new session to S3
func (u *Uploader) ManageConection() error {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(u.s3Region)})
	if err != nil {
		return errors.Wrap(err, "S3 connection")
	}
	u.s3manager = s3manager.NewUploader(sess)
	return nil
}

// Upload will upload a single file to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func (u *Uploader) Upload(filename string) error {

	file, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err, "OS File Open")
	}
	defer file.Close()

	fmt.Printf("Uploading %s to S3...\n", filename)
	result, err := u.s3manager.Upload(&s3manager.UploadInput{
		Bucket: aws.String(u.s3Bucket),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		return errors.Wrap(err, "S3 upload")
	}
	fmt.Printf("Successfully uploaded %s to %s\n", filename, result.Location)
	return nil
}
