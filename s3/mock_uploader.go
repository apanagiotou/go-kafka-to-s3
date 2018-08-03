package s3

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type MockUploaderAPI struct {
	Error string
}

func (u MockUploaderAPI) Upload(*s3manager.UploadInput, ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	if u.Error != "" {
		return nil, errors.New(u.Error)
	}

	return &s3manager.UploadOutput{}, nil
}

func (u MockUploaderAPI) UploadWithContext(aws.Context, *s3manager.UploadInput, ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	if u.Error != "" {
		return nil, errors.New(u.Error)
	}

	return &s3manager.UploadOutput{}, nil
}
