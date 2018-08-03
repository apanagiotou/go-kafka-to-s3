package s3

import (
	"os"
	"testing"
)

func TestUpload(t *testing.T) {

	tests := []struct {
		uploadererror    string
		uploaderresponce error
	}{
		{"", nil},
		{"iserror", nil},
	}

	for _, test := range tests {
		mu := MockUploader{Error: test.uploadererror}
		u := &Uploader{s3Region: "testregion", s3Bucket: "testbucket", s3manager: mu}
		os.Create("testfile.log")

		err := u.Upload("testfile.log")
		if err != test.uploaderresponce {
			t.Errorf("Upload did not returned expected error: %t, returned: %t.", test.uploaderresponce, err)
		}
		os.Remove("testfile.log")
	}
}
