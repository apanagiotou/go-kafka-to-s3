package s3

import (
	"errors"
)

type MockUploader struct {
	Error string
}

type UploadOutput struct {
	Location string
}

type UploadInput struct {
	Location string
}

func (u *MockUploader) Upload(input *UploadInput, options ...func(*MockUploader)) (*UploadOutput, error) {
	if u.Error != "" {
		return nil, errors.New(u.Error)
	}

	return &UploadOutput{"fakelocation"}, nil
}
