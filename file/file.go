package file

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

// RotateWriter writes and rotates files
type RotateWriter struct {
	lock     sync.Mutex
	filename string
	fp       *os.File
	maxsize  int64
}

// New makes a new RotateWriter. Return nil if error occurs during setup.
func New(filename string, maxsize int64) (w *RotateWriter, err error) {
	w = &RotateWriter{filename: filename, maxsize: maxsize}
	w.fp, err = os.OpenFile(filename, syscall.O_RDWR|syscall.O_CREAT, 0666)
	if err != nil {
		return nil, nil
	}
	return w, nil
}

// Write writes a  new line string to file
func (w *RotateWriter) Write(line string) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.fp.WriteString(fmt.Sprintf("%s\n", line))
}

// Rotateable checks if the file is ready for rotation. Meaning that it reached the size.
func (w *RotateWriter) Rotateable() (bool, error) {
	fi, err := os.Stat(w.filename)
	if err != nil {
		return false, errors.Wrap(err, "OS File stat")
	}
	size := fi.Size()
	if size > w.maxsize {
		return true, nil
	}
	return false, nil
}

// Rotate performs the actual act of rotating and reopening file. Returns the rotated filename
func (w *RotateWriter) Rotate() (string, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	rotatedFilename := ""

	// Close existing file if open
	if w.fp != nil {
		err := w.fp.Close()
		w.fp = nil
		if err != nil {
			return rotatedFilename, errors.Wrap(err, "OS File close")
		}
	}
	// Rename the dest file if it already exists
	_, err := os.Stat(w.filename)
	if err == nil {
		rotatedFilename = w.filename + "." + time.Now().Format(time.RFC3339)
		err = os.Rename(w.filename, rotatedFilename)
		if err != nil {
			return rotatedFilename, errors.Wrap(err, "OS File rename")
		}
	}

	// Create the file.
	w.fp, err = os.Create(w.filename)
	if err != nil {
		return rotatedFilename, errors.Wrap(err, "OS File create")
	}
	return rotatedFilename, nil
}

// Compress compresses plain files
func Compress(filename string) (string, error) {

	compressed := filename + ".gz"

	// Open file on disk.
	f, err := os.Open(filename)
	if err == nil {
		return "", errors.Wrap(err, "OS File open")
	}

	// Create a Reader and use ReadAll to get all the bytes from the file.
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	// Open file for writing.
	f, err = os.Create(compressed)
	if err == nil {
		return "", errors.Wrap(err, "OS File create")
	}

	// Write compressed data.
	fmt.Printf("Compressing %s\n", filename)
	w := gzip.NewWriter(f)
	w.Write(content)
	w.Close()

	// Remove old
	os.Remove(filename)

	// Done.
	fmt.Printf("File %s compressed\n", filename)
	return compressed, nil
}
