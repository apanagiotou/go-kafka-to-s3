package file

import (
	"os"
	"testing"
)

// TestWrite writes 9 bytes to the file (testline) and expects Write to return 9
func TestWrite(t *testing.T) {
	wr, _ := New("testfile.log", 1)
	line, _ := wr.Write("12345678")
	if line != 9 {
		t.Errorf("Write did not returned expected int: %d, returned: %d.", 9, line)
	}
	os.Remove("testfile.log")
}

// TestRotateable writes a string to the file and checks if is a candidate for rotation
func TestRotateable(t *testing.T) {

	tests := []struct {
		linestring string
		maxsize    int64
		rotateable bool
	}{
		{"1", 2, false},
		{"12345678", 2, true},
	}

	for _, test := range tests {
		wr, _ := New("testfile.log", test.maxsize)
		wr.Write(test.linestring)
		rotateable, _ := wr.Rotateable()
		if rotateable != test.rotateable {
			t.Errorf("Rotateable did not returned expected bool: %t, returned: %t.", test.rotateable, rotateable)
		}
		os.Remove("testfile.log")
	}
}

// TestRotate rotates a file and checks if the filename is not empty
func TestRotate(t *testing.T) {

	wr, _ := New("testfile.log", 1)
	wr.Write("12345678")
	rotatedfile, _ := wr.Rotate()
	if rotatedfile == "" {
		t.Errorf("Rotate failed to rotate the file")
	}
	os.Remove("testfile.log")
	os.Remove(rotatedfile)
}
