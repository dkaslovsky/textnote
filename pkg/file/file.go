package file

import (
	"io"
	"os"
)

// ReadWriteable is the interface on which file operations are executed
type ReadWriteable interface {
	Load(io.Reader) error
	Write(io.Writer) error
	GetFilePath() string
}

// ReadWriter executes file operations
type ReadWriter struct{}

// NewReadWriter constructs a new ReadWriter
func NewReadWriter() *ReadWriter {
	return &ReadWriter{}
}

// Read reads from file
func (rw *ReadWriter) Read(rwable ReadWriteable) error {
	r, err := os.Open(rwable.GetFilePath())
	if err != nil {
		return err
	}
	defer r.Close()
	return rwable.Load(r)
}

// Overwrite writes a template to a file, overwriting existing file contents if any
func (rw *ReadWriter) Overwrite(rwable ReadWriteable) error {
	f, err := os.OpenFile(rwable.GetFilePath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	err = rwable.Write(f)
	if err != nil {
		return err
	}
	return nil
}

// Exists evaluates if a file exists
func (rw *ReadWriter) Exists(rwable ReadWriteable) bool {
	fileName := rwable.GetFilePath()
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}
